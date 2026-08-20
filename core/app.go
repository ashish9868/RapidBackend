package core

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/a-h/templ"
	"github.com/ashish9868/rapidbackend/constants"
	respository "github.com/ashish9868/rapidbackend/core/repository"
	"github.com/ashish9868/rapidbackend/middlewares"
	"github.com/ashish9868/rapidbackend/models"
	"github.com/ashish9868/rapidbackend/utils"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-playground/form"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/xid"
	"github.com/starfederation/datastar-go/datastar"
)

type App struct {
	RootRouter            *http.ServeMux
	DB                    *sqlx.DB
	FeFs                  *fs.FS
	MigrationFs           *fs.FS
	Version               *string
	BaseRepository        *respository.BaseRepository
	AuthRepository        *respository.AuthRepository
	AccessTokenRepository *respository.AccessTokenRepository
}

type ResourceAction struct {
	Handler     func(w http.ResponseWriter, r *http.Request, app *App)
	Middlewares []middlewares.Middleware
}
type ResourceHandler struct {
	Index      *ResourceAction
	Show       *ResourceAction
	Create     *ResourceAction
	CreateForm *ResourceAction
	Update     *ResourceAction
	UpdateForm *ResourceAction
	Delete     *ResourceAction
}

func NewApp(embed embed.FS) *App {

	time.Local = time.UTC

	data_dir := "app_data"
	env_path := data_dir + "/.env"
	env_data := []string{
		fmt.Sprintf(`PORT=%d`, constants.DEFAULT_PORT),
		fmt.Sprintf(`DATA_DIR=%s`, "app_data"),
		fmt.Sprintf(`DEBUG=%d`, 1),
		fmt.Sprintf(`GIN_MODE=%s`, "release"),
		fmt.Sprintf(`ENCRYPTION_KEY=%s`, utils.HashPassword(xid.New().String())),
	}
	err := utils.SafeCreateFile(env_path, strings.Join(env_data, "\n"))

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
	sqldb, err := sqlx.Connect("sqlite3", dsn)
	if err != nil {
		slog.Error(err.Error())
		panic(err)
	}

	sqldb.SetMaxOpenConns(10)
	sqldb.SetMaxIdleConns(5)

	baseRepository := &respository.BaseRepository{DB: sqldb}
	authRepository := &respository.AuthRepository{BaseRepository: baseRepository}
	accessTokenRepository := &respository.AccessTokenRepository{BaseRepository: baseRepository}
	rootRouter := &http.ServeMux{}
	app := &App{
		RootRouter:            rootRouter,
		DB:                    sqldb,
		FeFs:                  utils.SubFs(embed, "static"),
		MigrationFs:           utils.SubFs(embed, "database/migrations"),
		BaseRepository:        baseRepository,
		AuthRepository:        authRepository,
		AccessTokenRepository: accessTokenRepository,
	}
	app.serveStatic()
	app.serveNoRoute()
	app.DBMigrate()
	fmt.Printf("APP will start on PORT: %s\n\n", utils.SafeEnvGet("PORT", strconv.Itoa(constants.DEFAULT_PORT)))
	return app
}

func (app *App) ResourceRoutes(name string, group *http.ServeMux, handler ResourceHandler, _middlewares ...middlewares.Middleware) {
	base := "/" + strings.Trim(name, "/")

	id_segment := ""
	parts := strings.Split(base, "/")
	if len(parts) > 0 {
		id_segment = utils.Singular(parts[len(parts)-1])
	}

	id_segment = id_segment + "_id"

	if handler.Index != nil {
		route := fmt.Sprintf("GET %s", base)
		utils.LogF("Registering route : %s", route)
		handlers := append(_middlewares, handler.Index.Middlewares...)
		group.Handle(route, middlewares.Chain(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handler.Index.Handler(w, r, app)
		}), handlers...))
	}

	if handler.Show != nil {
		route := fmt.Sprintf("GET %s/:%s", base, id_segment)
		utils.LogF("Registering route : %s", route)
		handlers := append(_middlewares, handler.Show.Middlewares...)
		group.Handle(route, middlewares.Chain(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handler.Show.Handler(w, r, app)
		}), handlers...))

	}

	if handler.Create != nil {
		route := fmt.Sprintf("POST %s", base)
		utils.LogF("Registering route : %s", route)
		handlers := append(_middlewares, handler.Create.Middlewares...)
		group.Handle(route, middlewares.Chain(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handler.Create.Handler(w, r, app)
		}), handlers...))

	}

	if handler.CreateForm != nil {
		route := fmt.Sprintf("GET %s/create", base)
		utils.LogF("Registering route : %s", route)
		handlers := append(_middlewares, handler.CreateForm.Middlewares...)
		group.Handle(route, middlewares.Chain(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handler.CreateForm.Handler(w, r, app)
		}), handlers...))
	}

	if handler.Update != nil {
		route := fmt.Sprintf("PUT %s/:%s", base, id_segment)
		route_patch := fmt.Sprintf("PATCH %s/:%s", base, id_segment)
		utils.LogF("Registering route : %s, %s", route, route_patch)

		handlers := append(_middlewares, handler.Update.Middlewares...)
		group.Handle(route, middlewares.Chain(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handler.Update.Handler(w, r, app)
		}), handlers...))

		group.Handle(route_patch, middlewares.Chain(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handler.Update.Handler(w, r, app)
		}), handlers...))
	}

	if handler.UpdateForm != nil {
		route := fmt.Sprintf("GET %s/:%s/update", base, id_segment)
		utils.LogF("Registering route : %s", route)
		handlers := append(_middlewares, handler.UpdateForm.Middlewares...)
		group.Handle(route, middlewares.Chain(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handler.UpdateForm.Handler(w, r, app)
		}), handlers...))
	}

	if handler.Delete != nil {
		route := fmt.Sprintf("DELETE %s/:%s", base, id_segment)
		utils.LogF("Registering route : %s", route)
		handlers := append(_middlewares, handler.Delete.Middlewares...)
		group.Handle(route, middlewares.Chain(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handler.Delete.Handler(w, r, app)
		}), handlers...))
	}

}

func (app *App) DBMigrate() {
	// execute migrations
	if app.MigrationFs != nil {
		files, _ := utils.ListFiles(*app.MigrationFs, ".sql")
		for _, file := range files {
			data := utils.ReadFsFile(*app.MigrationFs, file)
			if len(data) > 0 {
				utils.LogF("Executing SQL %s", data)
				app.DB.Exec(data)
			}
		}
	}
}

func (app *App) serveStatic() {
	if app.FeFs != nil {
		fs := http.FileServer(http.FS(*app.FeFs))
		app.RootRouter.Handle("/static/", http.StripPrefix("/static", fs))
	}
}

func (app *App) serveNoRoute() {
	app.RootRouter.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/login", http.StatusPermanentRedirect)
	}))
}

func (app *App) BindSafely(w http.ResponseWriter, r *http.Request, obj any) map[string]string {
	contentType := r.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "application/json") {
		r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		err := decoder.Decode(obj)
		if err != nil {
			return app.FormatErrors(err)
		}
	}
	decoder := form.NewDecoder()
	if err := r.ParseForm(); err != nil {
		return app.FormatErrors(err)
	}
	err := decoder.Decode(obj, r.PostForm)

	if err != nil {
		return app.FormatErrors(err)
	}
	return nil
}

func (a *App) GetUserFromRequest(r *http.Request) *models.User {
	user, ok := r.Context().Value(constants.USER_CONTEXT_KEY).(*models.User)
	if ok {
		return user
	}
	return nil
}

func (app *App) GetSSE(w http.ResponseWriter, r *http.Request) *datastar.ServerSentEventGenerator {
	return datastar.NewSSE(w, r)
}

func (app *App) SetAuthCookie(w http.ResponseWriter, value string, ageHours int) {
	http.SetCookie(w, &http.Cookie{
		Name:     "__auth",
		Value:    value,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   60 * ageHours,
	})
}
func (app *App) RenderComponent(w http.ResponseWriter, component templ.Component) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	component.Render(context.Background(), w)
}

func (app *App) FormatErrors(err error) map[string]string {
	result := make(map[string]string)

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
