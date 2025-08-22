package main

import (
	"context"
	"database/sql"
	"embed"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path"
	"strconv"
	"time"

	"github.com/angelofallars/htmx-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	gowebly "github.com/gowebly/helpers"
	_ "github.com/mattn/go-sqlite3"
	database "github.com/pleimann/camel-do/db"

	"github.com/pleimann/camel-do/services/cal"
	"github.com/pleimann/camel-do/services/home"
	"github.com/pleimann/camel-do/services/oauth"
	"github.com/pleimann/camel-do/services/project"
	"github.com/pleimann/camel-do/services/task"
	"github.com/pleimann/camel-do/templates/components"
)

//go:embed all:static
var static embed.FS

//go:embed credentials.json
var credentials string

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	slog.SetDefault(logger)

	var seed bool
	flag.BoolVar(&seed, "seed", false, "seed database with some data")
	flag.Parse()

	var err error

	if err = createDatabase(); err != nil {
		log.Fatalf("Failed to create database service! %s", err)
	}

	defer db.Close()

	googleAuth := oauth.NewGoogleAuth(credentials)

	taskSyncService, err = task.NewTaskSyncService(googleAuth, db)
	if err != nil {
		log.Fatalf("Failed create task sync service! %s", err)
	}

	calendarService, err = cal.NewCalendarService(&cal.CalendarServiceConfig{}, googleAuth, db)
	if err != nil {
		log.Fatalf("error creating CalendarService: %s", err)
	}

	projectService, err = project.NewProjectService(&project.ProjectServiceConfig{}, db)
	if err != nil {
		log.Fatalf("error creating ProjectService: %s", err)
	}

	taskService, err = task.NewTaskService(&task.TaskServiceConfig{}, db)
	if err != nil {
		log.Fatalf("error creating TaskService: %s", err)
	}

	if tasks, _ := taskService.GetTodaysTasks(); tasks.IsEmpty() {
		database.Seed(20, taskService, projectService)
	}

	// Run your server.
	if err := runServer(); err != nil {
		log.Fatal("Failed to start server!", "details", err.Error())
		os.Exit(1)
	}

	os.Exit(0)
}

const databaseFileName = "camel-do.db"

var db *sql.DB
var taskService *task.TaskService
var taskSyncService *task.TaskSyncService
var calendarService *cal.CalendarService
var projectService *project.ProjectService

func createDatabase() error {
	var err error

	userConfigDir, err := os.UserConfigDir()

	if err != nil {
		return err
	}

	databasePath := path.Join(userConfigDir, "camel-do", databaseFileName)

	if err := os.MkdirAll(path.Dir(databasePath), 0700); err != nil {
		log.Fatalf("Unable to create directory for db file: %v", err)
	}

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
	e.Logger.SetHeader("time=${time_rfc3339_nano} ${level}")
	e.HTTPErrorHandler = customHTTPErrorHandler

	// Add Echo middlewares.
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:   true,
		LogURI:      true,
		LogError:    true,
		HandleError: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				slog.LogAttrs(context.Background(), slog.LevelInfo, "REQUEST",
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
				)
			} else {
				slog.LogAttrs(context.Background(), slog.LevelError, "REQUEST_ERROR",
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
					slog.String("err", v.Error.Error()),
				)
			}
			return nil
		},
	}))
	e.Use(middleware.Recover())

	// Handle index page view.
	indexViewHandler := home.NewHomeHandler(taskService, calendarService, projectService)
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
