package core

import (
	"context"
	"database/sql"
	"embed"
	"io/fs"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/ashish9868/rapidbackend/models"
	"github.com/ashish9868/rapidbackend/utils"
	"github.com/gin-gonic/gin"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/extra/bundebug"
	_ "modernc.org/sqlite"
)

type App struct {
	Bun      *bun.DB
	BaseUtil *utils.BaseUtil
	Gin      *gin.Engine
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

func NewApp(embeddedFs embed.FS) *App {
	os.Setenv("BUNDEBUG", "1")
	gin.SetMode(gin.ReleaseMode)

	baseUtil := utils.NewBaseUtil()
	time.Local = time.UTC

	db_folder := "app_data"

	stat, err := os.Stat(db_folder)

	createFolder := true
	if err == nil {
		if stat.IsDir() {
			createFolder = false
		}
	}

	if createFolder {
		err := os.Mkdir("app_data", os.ModePerm)
		if err != nil {
			panic("Unable to create data directory")
		}
	}

	dsn := "file:app_data/app.db?" + url.Values{
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

	engine := gin.Default()
	tmplFS, err := fs.Sub(embeddedFs, "templates")
	if err != nil {
		panic(err)
	}

	staticFS, err := fs.Sub(tmplFS, "static")
	if err != nil {
		panic(err)
	}

	engine.StaticFS("/static", http.FS(staticFS))

	engine.GET("/", func(ctx *gin.Context) {
		ctx.FileFromFS("index.htm", http.FS(tmplFS))
	})
	engine.NoRoute(func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		if len(path) < 1 {
			path = "index"
		}
		ctx.FileFromFS(path+".htm", http.FS(tmplFS))
	})

	app := &App{Bun: db, BaseUtil: baseUtil, Gin: engine}
	app.InitializeBaseMigrations()
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

func (app *App) InitializeBaseMigrations() {
	app.Bun.NewCreateTable().Model((*models.SuperAdmin)(nil)).IfNotExists().WithForeignKeys().Exec(context.Background())
	app.Bun.NewCreateTable().Model((*models.Project)(nil)).IfNotExists().WithForeignKeys().Exec(context.Background())
	app.Bun.NewCreateTable().Model((*models.ProjectUser)(nil)).IfNotExists().WithForeignKeys().Exec(context.Background())
	app.Bun.NewCreateTable().Model((*models.Collection)(nil)).IfNotExists().WithForeignKeys().Exec(context.Background())
	app.Bun.NewCreateTable().Model((*models.CollectionField)(nil)).IfNotExists().WithForeignKeys().Exec(context.Background())
	app.Bun.NewCreateTable().Model((*models.CollectionRecord)(nil)).IfNotExists().WithForeignKeys().Exec(context.Background())
	app.Bun.NewCreateTable().Model((*models.Settings)(nil)).IfNotExists().WithForeignKeys().Exec(context.Background())
}
