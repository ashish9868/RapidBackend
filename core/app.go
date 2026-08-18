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
	"github.com/ashish9868/rapidbackend/core/respository"
	"github.com/ashish9868/rapidbackend/middlewares"
	"github.com/ashish9868/rapidbackend/utils"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/xid"
)

var schema = `
CREATE TABLE superadmins (
    id VARCHAR(255) PRIMARY KEY NOT NULL,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255),
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    email_verified_at DATETIME,
    is_active BOOLEAN DEFAULT 0,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE projects (
    id VARCHAR(255) PRIMARY KEY NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    slug VARCHAR(255) NOT NULL UNIQUE,
    settings JSON,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE users (
    id VARCHAR(255) PRIMARY KEY NOT NULL,

    project_id VARCHAR(255),

    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255),
    email VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,

    email_verified_at DATETIME NOT NULL,
    is_active BOOLEAN DEFAULT 0,

    permissions_json JSON,

    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT unq_project_email
        UNIQUE (project_id, email),

    FOREIGN KEY (project_id)
        REFERENCES projects(id)
        ON DELETE CASCADE
);


CREATE TABLE access_key_tokens (
    id VARCHAR(255) PRIMARY KEY NOT NULL,

    collection_id VARCHAR(255) NOT NULL,
    collection VARCHAR(255) NOT NULL,

    access_token VARCHAR(255),

    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE project_collections (
    id VARCHAR(255) PRIMARY KEY NOT NULL,

    project_id VARCHAR(255) NOT NULL,

    name VARCHAR(255),
    type VARCHAR(255) DEFAULT 'base',
    sort_order INTEGER DEFAULT 0,
    required BOOLEAN DEFAULT 0,

    rules JSON,
    options JSON,

    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (project_id)
        REFERENCES projects(id)
        ON DELETE CASCADE
);


CREATE TABLE project_collection_fields (
    id VARCHAR(255) PRIMARY KEY NOT NULL,

    collection_id VARCHAR(255) NOT NULL,

    name VARCHAR(255),
    type VARCHAR(255) DEFAULT 'base',

    is_required BOOLEAN DEFAULT 0,
    is_indexed BOOLEAN DEFAULT 0,
    is_unique BOOLEAN DEFAULT 0,
    is_sortable BOOLEAN DEFAULT 0,
    is_filterable BOOLEAN DEFAULT 0,

    options JSON,

    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (collection_id)
        REFERENCES project_collections(id)
        ON DELETE CASCADE
);


CREATE TABLE project_collection_records (
    id VARCHAR(255) PRIMARY KEY NOT NULL,

    collection_id VARCHAR(255) NOT NULL,

    data JSON,

    version INTEGER DEFAULT 1,

    name VARCHAR(255),
    type VARCHAR(255) DEFAULT 'base',

    is_required BOOLEAN DEFAULT 0,
    is_indexed BOOLEAN DEFAULT 0,
    is_unique BOOLEAN DEFAULT 0,
    is_sortable BOOLEAN DEFAULT 0,
    is_filterable BOOLEAN DEFAULT 0,

    options JSON,

    created_by_id VARCHAR(255) NOT NULL,
    updated_by_id VARCHAR(255) NOT NULL,

    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (collection_id)
        REFERENCES project_collections(id)
        ON DELETE CASCADE,

    FOREIGN KEY (created_by_id)
        REFERENCES users(id),

    FOREIGN KEY (updated_by_id)
        REFERENCES users(id)
);


CREATE TABLE project_pages (
    id VARCHAR(255) PRIMARY KEY NOT NULL,

    project_id VARCHAR(255) NOT NULL,

    title VARCHAR(255) NOT NULL,
    description TEXT,

    slug VARCHAR(255) NOT NULL UNIQUE,

    settings JSON,

    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT unq_project_page
        UNIQUE (project_id, title),

    FOREIGN KEY (project_id)
        REFERENCES projects(id)
        ON DELETE CASCADE
);


CREATE TABLE email_templates (
    id VARCHAR(255) PRIMARY KEY NOT NULL,

    name VARCHAR(255) UNIQUE,
    description TEXT,

    is_system_template BOOLEAN DEFAULT 1,

    html_content TEXT,
    text_content TEXT,

    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE settings (
    id VARCHAR(255) PRIMARY KEY NOT NULL,

    project_id VARCHAR(255),

    value TEXT,

    created_by_id VARCHAR(255) NOT NULL,
    updated_by_id VARCHAR(255) NOT NULL,

    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (project_id)
        REFERENCES projects(id)
        ON DELETE CASCADE,

    FOREIGN KEY (created_by_id)
        REFERENCES users(id),

    FOREIGN KEY (updated_by_id)
        REFERENCES users(id)
);
`

type App struct {
	RootRouter     *http.ServeMux
	DB             *sqlx.DB
	FeFs           *fs.FS
	Version        *string
	BaseRepository *respository.BaseRepository
	AuthRepository *respository.AuthRepository
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
	rootRouter := &http.ServeMux{}
	app := &App{
		RootRouter:     rootRouter,
		DB:             sqldb,
		FeFs:           utils.SubFs(embed, "static"),
		BaseRepository: baseRepository,
		AuthRepository: &respository.AuthRepository{BaseRepository: baseRepository},
	}
	app.serveStatic()
	app.serveNoRoute()
	app.initializeSystem()
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

func (app *App) initializeSystem() {
	app.DB.Exec(schema)
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

func (app *App) BindSafely(w http.ResponseWriter, r *http.Request, obj any) error {
	contentType := r.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "application/json") {
		r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		return decoder.Decode(obj)
	}
	return r.ParseForm()
}

type Response struct {
	Code       int
	View       string
	Data       any
	Error      error
	HxRedirect string
	FormData   any
}

func (app *App) RenderComponent(w http.ResponseWriter, component templ.Component) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	component.Render(context.Background(), w)
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
