package main

import (
	"context"
	"database/sql"
	"embed"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"os/user"
	"strconv"
	"time"

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

	// slog.SetLogLoggerLevel(slog.LevelDebug)

	usr, err := user.Current()
	if err != nil {
		log.Fatalf("Failed to locate the user's home directory: %s", err)
	}

	log.Printf("current user", "name", usr.Name, "home-dir", usr.HomeDir)

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
		log.Fatal("Failed to start server!", "details", err.Error())
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
	e := echo.New()
	e.HTTPErrorHandler = customHTTPErrorHandler

	// Add Echo middlewares.
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Handle index page view.
	indexViewHandler := home.NewHomeHandler(taskService, projectService)
	e.GET("/", indexViewHandler.ServeHTTP).Name = "root"

	// Serve embedded static files found at ./static
	e.StaticFS("/static", echo.MustSubFS(static, "static")).Name = "static"

	projectsGroup := e.Group("/projects")
	project.NewProjectHandler(projectsGroup, projectService)

	tasksGroup := e.Group("/tasks")
	task.NewTaskHandler(tasksGroup, taskService, projectService)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// Start server
	go func() {
		if err := e.Start(fmt.Sprintf(":%d", port)); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server with a timeout of 10 seconds.
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return e.Shutdown(ctx)
}

func customHTTPErrorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}

	messageTemplate := components.ErrorMessage(err.Error())

	htmx.
		NewResponse().
		StatusCode(code).
		RenderTempl(c.Request().Context(), c.Response().Writer, messageTemplate)
}
