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

	"github.com/ashish9868/rapidbackend/core/services"
	"github.com/ashish9868/rapidbackend/models"
	"github.com/ashish9868/rapidbackend/utils"
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/joho/godotenv"
	"github.com/rs/xid"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/extra/bundebug"
	_ "modernc.org/sqlite"
)

type App struct {
	Bun         *bun.DB
	BaseUtil    *utils.BaseUtil
	Gin         *gin.Engine
	FeFs        *fs.FS
	Version     *string
	AuthService *services.AuthService
	Repository  *Repository
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

func NoCacheMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate")
		c.Header("Pragma", "no-cache")
		c.Header("Expires", time.Unix(0, 0).Format(http.TimeFormat)) // Thu, 01 Jan 1970 00:00:00 GMT
		c.Next()
	}
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
	engine.Use(NoCacheMiddleware())
	app := &App{Bun: db, BaseUtil: baseUtil, Gin: engine, FeFs: baseUtil.SubFs(embed, "static"), AuthService: services.NewAuthService(db, *baseUtil), Repository: NewRepository(db)}
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
	id_segment := ""
	parts := strings.Split(base, "/")
	if len(parts) > 0 {
		id_segment = app.BaseUtil.Singular(parts[len(parts)-1])
	}

	id_segment = id_segment + "_id"
	if handler.Show != nil {
		group.Use(append(middlewares, handler.Show.Middlewares...)...)
		group.GET(base+id_segment, func(ctx *gin.Context) {
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
		group.PUT(base+id_segment, func(ctx *gin.Context) {
			handler.Update.Handler(ctx, app)
		})
		group.PATCH(base+id_segment, func(ctx *gin.Context) {
			handler.Update.Handler(ctx, app)
		})
	}

	if handler.Delete != nil {
		group.Use(append(middlewares, handler.Delete.Middlewares...)...)
		group.DELETE(base+id_segment, func(ctx *gin.Context) {
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
	app.Bun.NewCreateTable().Model((*models.Project)(nil)).IfNotExists().WithForeignKeys().Exec(context.Background())
	app.Bun.NewCreateTable().Model((*models.Superadmin)(nil)).IfNotExists().WithForeignKeys().Exec(context.Background())
	app.Bun.NewCreateTable().Model((*models.User)(nil)).IfNotExists().WithForeignKeys().Exec(context.Background())
	app.Bun.NewCreateTable().Model((*models.AccessKeyToken)(nil)).IfNotExists().WithForeignKeys().Exec(context.Background())
	app.Bun.NewCreateTable().Model((*models.ProjectPage)(nil)).IfNotExists().WithForeignKeys().Exec(context.Background())
	app.Bun.NewCreateTable().Model((*models.ProjectCollection)(nil)).IfNotExists().WithForeignKeys().Exec(context.Background())
	app.Bun.NewCreateTable().Model((*models.ProjectCollectionField)(nil)).IfNotExists().WithForeignKeys().Exec(context.Background())
	app.Bun.NewCreateTable().Model((*models.ProjectCollectionRecord)(nil)).IfNotExists().WithForeignKeys().Exec(context.Background())
	app.Bun.NewCreateTable().Model((*models.EmailTemplate)(nil)).IfNotExists().WithForeignKeys().Exec(context.Background())
	app.Bun.NewCreateTable().Model((*models.SystemSetting)(nil)).IfNotExists().WithForeignKeys().Exec(context.Background())

	t := time.Now()
	app.Bun.NewInsert().Model(&models.User{
		ID:              xid.New().String(),
		FirstName:       "Ashish",
		LastName:        "Kumar",
		Email:           "funappzco@gmail.com",
		Password:        app.BaseUtil.HashPassword("Asdf1234@#$"),
		IsActive:        true,
		EmailVerifiedAt: &t,
	}).Exec(context.Background())
}

func (app *App) ServeStatic() {
	if app.FeFs != nil {
		app.Gin.StaticFS("/static", http.FS(*app.FeFs))
	}
}

func (app *App) ServeNoRoute() {
	app.Gin.NoRoute(app.NewAuthMiddleWare(false, true), func(ctx *gin.Context) {
		ctx.Redirect(http.StatusPermanentRedirect, "/")
	})
}

func (app *App) ErrorJson(body any, err error) gin.H {
	return gin.H{
		"body":   body,
		"errors": err,
	}
}

func (app *App) SetAuthCookie(ctx *gin.Context, value string, maxAge int) {
	if maxAge < 60 {
		token, _ := ctx.Cookie(gin.AuthUserKey)
		if len(token) > 0 {
			app.Bun.NewDelete().Model(&models.AccessKeyToken{}).Where("access_token = ?", token).Exec(context.Background())
		}
	}
	ctx.SetCookie(gin.AuthUserKey, value, maxAge, "/", ctx.Request.Host, false, true)
}

func (app *App) HttpUnauthorized(ctx *gin.Context) {
	val, _ := ctx.Cookie(gin.AuthUserKey)
	app.SetAuthCookie(ctx, "", -1)
	if app.IsHTMX(ctx) {
		ctx.Header("HX-Redirect", "/#/")
		ctx.Status(http.StatusUnauthorized)
	} else {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"Uauthorized": true,
			"Redirect":    len(val) > 0,
		})
	}
	ctx.Abort()
	return
}

func (app *App) IsHTMX(ctx *gin.Context) bool {
	return ctx.GetHeader("HX-Request") == "true"
}

func (app *App) NewAuthMiddleWare(throw401 bool, silent bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := ctx.Cookie(gin.AuthUserKey)
		if err != nil {
			token = ctx.GetHeader("Authorization")
			if !strings.HasPrefix(token, "Bearer ") && !silent {
				app.HttpUnauthorized(ctx)
				return
			}
			token = strings.TrimPrefix(token, "Bearer ")
			if len(token) < 1 && !silent {
				app.HttpUnauthorized(ctx)
				return
			}
		}

		user := app.AuthService.GetUserByToken(token)
		if user != nil {
			ctx.Set(gin.AuthUserKey, user)
		}
		if user == nil && throw401 && !silent {
			app.HttpUnauthorized(ctx)
			return
		}
		ctx.Next()
	}

}

func (app *App) NewPublicMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		_, exists := ctx.Get(gin.AuthUserKey)
		if !exists {
			ctx.Next()
			return
		}
		ctx.Header("HX-Redirect", "/#/dashboard")
		ctx.Status(http.StatusTemporaryRedirect)
		ctx.Abort()
	}
}

func (app *App) BindSafely(ctx *gin.Context, obj any) error {
	switch ctx.ContentType() {
	case "application/json":
		return ctx.ShouldBindJSON(obj)
	default:
		return ctx.ShouldBind(obj)
	}
}

type Response struct {
	Code       int
	View       string
	Data       any
	Error      error
	HxRedirect string
	FormData   any
}

func (app *App) SendResponse(ctx *gin.Context, response Response) {
	app.BaseUtil.PrintFiles(*app.FeFs)
	formData := response.FormData
	if formData == nil {
		formData = map[string]any{}
	}
	templateData := gin.H{
		"data":   response.Data,
		"errors": app.FormatErrors(response.Error),
		"code":   response.Code,
		"form":   formData,
	}
	if app.IsHTMX(ctx) {
		if len(response.HxRedirect) > 0 {
			ctx.Header("HX-Redirect", response.HxRedirect)
			ctx.Status(http.StatusTemporaryRedirect)
			return
		} else {
			ctx.HTML(200, response.View, templateData)
			return
		}
	} else {
		ctx.JSON(response.Code, templateData)
	}
	ctx.Abort()
}

func (app *App) FormatErrors(err error) map[string]any {
	result := make(map[string]any)

	if errs, ok := err.(validation.Errors); ok {
		for field, e := range errs {
			if e != nil {
				result[field] = e.Error()
			}
		}
		return result
	}

	result["global"] = err.Error()
	return result

}
