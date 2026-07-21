package core

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/ashish9868/rapidbackend/models"
	"github.com/ashish9868/rapidbackend/utils"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rs/xid"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/extra/bundebug"
	_ "modernc.org/sqlite"
)

type App struct {
	Bun      *bun.DB
	BaseUtil *utils.BaseUtil
	Gin      *gin.Engine
	FeFs     *fs.FS
	Version  *string
}

type ResourceAction struct {
	Handler     func(ctx *gin.Context, app *App)
	Middlewares []gin.HandlerFunc
}
type ResourceHandler struct {
	Index  *ResourceAction
	Show   *ResourceAction
	Create *ResourceAction
	Update *ResourceAction
	Delete *ResourceAction
}

func NewApp(embed embed.FS) *App {

	baseUtil := utils.NewBaseUtil()
	time.Local = time.UTC

	data_dir := "app_data"
	env_path := data_dir + "/.env"
	env_data := []string{
		fmt.Sprintf(`PORT=%d`, utils.DEFAULT_PORT),
		fmt.Sprintf(`DATA_DIR=%s`, "app_data"),
		fmt.Sprintf(`BUNDEBUG=%d`, 1),
		fmt.Sprintf(`GIN_MODE=%s`, "release"),
		fmt.Sprintf(`ENCRYPTION_KEY=%s`, baseUtil.HashPassword(xid.New().String())),
	}
	err := baseUtil.SafeCreateFile(env_path, strings.Join(env_data, "\n"))

	if err != nil {
		println(fmt.Println(err.Error()))
		panic("Unable to create environment settings.")
	}

	err = godotenv.Load(env_path)
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	dsn := "file:" + data_dir + "/app.db?" + url.Values{
		"_pragma": []string{
			"journal_mode(WAL)",
			"synchronous(NORMAL)",
			"foreign_keys(ON)",
			"busy_timeout(10000)", // 10 seconds
			// "temp_store(MEMORY)",
			"cache_size(-20000)",   // ~20MB cache
			"mmap_size(268435456)", // 256MB
		},
	}.Encode()

	// Open database
	sqldb, err := sql.Open("sqlite", dsn)
	if err != nil {
		slog.Error(err.Error())
		panic(err)
	}

	sqldb.SetMaxOpenConns(10)
	sqldb.SetMaxIdleConns(5)

	// Create Bun instance
	db := bun.NewDB(sqldb, sqlitedialect.New())

	db.AddQueryHook(
		bundebug.NewQueryHook(

			bundebug.WithEnabled(false),
			bundebug.FromEnv("BUNDEBUG"),
		),
	)
	gin.SetMode(baseUtil.SafeEnvGet("GIN_MODE", gin.DebugMode))
	engine := gin.Default()

	app := &App{Bun: db, BaseUtil: baseUtil, Gin: engine, FeFs: baseUtil.SubFs(embed, "frontend/dist")}
	app.InitializeSystem()
	fmt.Printf("APP will start on PORT: %s\n\n", baseUtil.SafeEnvGet("PORT", strconv.Itoa(utils.DEFAULT_PORT)))
	return app
}

func (app *App) ResourceRoutes(name string, group *gin.RouterGroup, handler ResourceHandler, middlewares ...gin.HandlerFunc) {
	base := "/" + strings.Trim(name, "/")
	if handler.Index != nil {
		group.Use(append(middlewares, handler.Index.Middlewares...)...)
		group.GET(base, func(ctx *gin.Context) {
			handler.Index.Handler(ctx, app)
		})
	}
	if handler.Show != nil {
		group.Use(append(middlewares, handler.Show.Middlewares...)...)
		group.GET(base+"/:id", func(ctx *gin.Context) {
			handler.Show.Handler(ctx, app)
		})
	}

	if handler.Create != nil {
		group.Use(append(middlewares, handler.Create.Middlewares...)...)
		group.POST(base, func(ctx *gin.Context) {
			handler.Create.Handler(ctx, app)
		})
	}

	if handler.Update != nil {
		group.Use(append(middlewares, handler.Update.Middlewares...)...)
		group.PUT(base+"/:id", func(ctx *gin.Context) {
			handler.Update.Handler(ctx, app)
		})
		group.PATCH(base+"/:id", func(ctx *gin.Context) {
			handler.Update.Handler(ctx, app)
		})
	}

	if handler.Delete != nil {
		group.Use(append(middlewares, handler.Delete.Middlewares...)...)
		group.GET(base+"/:id", func(ctx *gin.Context) {
			handler.Delete.Handler(ctx, app)
		})
	}

}

func (app *App) WithTransaction(ctx context.Context, db *bun.DB, fn func(tx bun.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (app *App) InitializeSystem() {
	app.Bun.NewCreateTable().Model((*models.SuperAdmin)(nil)).IfNotExists().WithForeignKeys().Exec(context.Background())
	app.Bun.NewCreateTable().Model((*models.Project)(nil)).IfNotExists().WithForeignKeys().Exec(context.Background())
	app.Bun.NewCreateTable().Model((*models.ProjectUser)(nil)).IfNotExists().WithForeignKeys().Exec(context.Background())
	app.Bun.NewCreateTable().Model((*models.AccessKeyToken)(nil)).IfNotExists().WithForeignKeys().Exec(context.Background())
	app.Bun.NewCreateTable().Model((*models.ProjectPage)(nil)).IfNotExists().WithForeignKeys().Exec(context.Background())
	app.Bun.NewCreateTable().Model((*models.ProjectCollection)(nil)).IfNotExists().WithForeignKeys().Exec(context.Background())
	app.Bun.NewCreateTable().Model((*models.ProjectCollectionField)(nil)).IfNotExists().WithForeignKeys().Exec(context.Background())
	app.Bun.NewCreateTable().Model((*models.ProjectCollectionRecord)(nil)).IfNotExists().WithForeignKeys().Exec(context.Background())
	app.Bun.NewCreateTable().Model((*models.EmailTemplate)(nil)).IfNotExists().WithForeignKeys().Exec(context.Background())
	app.Bun.NewCreateTable().Model((*models.SystemSetting)(nil)).IfNotExists().WithForeignKeys().Exec(context.Background())
}

func (app *App) ServeStatic() {
	if app.FeFs != nil {
		staticSub := app.BaseUtil.SubFs(*app.FeFs, "assets")
		if staticSub != nil {
			app.Gin.StaticFS("/assets", http.FS(*staticSub))
		}
	}
}

func (app *App) ServeNoRoute() {
	handler := func(ctx *gin.Context) {
		if app.FeFs == nil {
			ctx.Status(http.StatusNotFound)
			return
		}
		app.BaseUtil.PrintFiles(*app.FeFs)
		path := "index.html"
		if !app.BaseUtil.FileExists(*app.FeFs, path) {
			ctx.Status(http.StatusNotFound)
			return
		}
		http.ServeFileFS(ctx.Writer, ctx.Request, *app.FeFs, path)
	}
	app.Gin.NoRoute(handler)
}

func (app *App) ErrorJson(body any, err error) gin.H {
	return gin.H{
		"body":   body,
		"errors": err,
	}
}
