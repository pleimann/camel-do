package main

import (
	"database/sql"
	"embed"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/user"
	"strconv"

	"github.com/angelofallars/htmx-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	gowebly "github.com/gowebly/helpers"
	_ "github.com/mattn/go-sqlite3"
	database "github.com/pleimann/camel-do/db"

	"github.com/pleimann/camel-do/services/home"
	"github.com/pleimann/camel-do/services/project"
	"github.com/pleimann/camel-do/services/task"
	"github.com/pleimann/camel-do/templates/components"
)

func main() {
	var seed bool
	flag.BoolVar(&seed, "seed", false, "seed database with some data")
	flag.Parse()

	slog.SetLogLoggerLevel(slog.LevelDebug)

	usr, err := user.Current()
	if err != nil {
		log.Fatalf("Failed to locate the user's home directory: %s", err)
	}

	slog.Info("current user", "name", usr.Name, "home-dir", usr.HomeDir)

	if err = createDatabase(); err != nil {
		log.Fatalf("Failed to create database service! %s", err)
	}

	defer db.Close()

	taskSyncService, err = task.NewTaskSyncService(db)
	if err != nil {
		log.Fatalf("Failed create sync service! %s", err)
	}

	projectService, err = project.NewService(&project.ProjectServiceConfig{}, db)
	if err != nil {
		log.Fatalf("error creating ProjectService: %s", err)
	}

	taskService, err = task.NewTaskService(&task.TaskServiceConfig{}, db)
	if err != nil {
		log.Fatalf("error creating TaskService: %s", err)
	}

	if seed {
		database.Seed(20, taskService, projectService)
	}

	// Run your server.
	if err := runServer(); err != nil {
		log.Fatalf("Failed to start server!", "details", err.Error())
		os.Exit(1)
	}

	os.Exit(0)
}

//go:embed all:static
var static embed.FS

const databasePath = "./camel-do.db"

var db *sql.DB
var taskService *task.TaskService
var taskSyncService *task.TaskSyncService
var projectService *project.ProjectService

func createDatabase() error {
	var err error

	if db, err = sql.Open("sqlite3", databasePath); err != nil {
		return err
	}

	if err = database.Migrate(db); err != nil {
		log.Fatal(err)
	}

	return nil
}

// runServer runs a new HTTP server with the loaded environment variables.
func runServer() error {
	// Validate environment variables.
	port, err := strconv.Atoi(gowebly.Getenv("PORT", "4000"))
	if err != nil {
		return err
	}

	// Create a new Echo server.
	router := echo.New()
	// router.HTTPErrorHandler = customHTTPErrorHandler

	// Add Echo middlewares.
	// router.Use(middleware.Logger())
	router.Use(middleware.Recover())

	// Handle index page view.
	indexViewHandler := home.NewHomeHandler(taskService, projectService)
	router.GET("/", indexViewHandler.ServeHTTP)

	// Serve embedded static files found at ./static
	router.StaticFS("/static", echo.MustSubFS(static, "static"))

	project.NewProjectHandler(router.Group("/projects"), projectService)
	task.NewTaskHandler(router.Group("/tasks"), taskService, projectService)

	return router.Start(fmt.Sprintf(":%d", port))

	// Create a new server instance with options from environment variables.
	// For more information, see https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/
	// Note: The ReadTimeout and WriteTimeout settings may interfere with SSE (Server-Sent Event) or WS (WebSocket) connections.
	// For SSE or WS, these timeouts can cause the connection to reset after 10 or 5 seconds due to the ReadTimeout and WriteTimeout setting.
	// If you plan to use SSE or WS, consider commenting out or removing the ReadTimeout and WriteTimeout key-value pairs.
	// server := http.Server{
	// 	Addr:         fmt.Sprintf(":%d", port),
	// 	Handler:      router, // handle all Echo routes
	// 	ReadTimeout:  5 * time.Second,
	// 	WriteTimeout: 10 * time.Second,
	// }

	// // Send log message.
	// slog.Info("Starting server...", "port", port)

	// return server.ListenAndServe()
}

func customHTTPErrorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}

	c.Logger().Error(err)

	messageTemplate := components.ErrorMessage(err.Error())

	htmx.NewResponse().
		Retarget("#messages").
		StatusCode(code).
		RenderTempl(c.Request().Context(), c.Response().Writer, messageTemplate)
}
