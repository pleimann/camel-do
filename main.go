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

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/mattn/go-sqlite3"
	database "github.com/pleimann/camel-do/db"
	"github.com/pleimann/camel-do/services/home"
	"github.com/pleimann/camel-do/services/project"
	"github.com/pleimann/camel-do/services/task"
	"github.com/pleimann/camel-do/utils"
)

var seed bool

func main() {
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

	// Run server
	if err = runServer(); err != nil {
		log.Fatalf("Failed to start server! %s", err)
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
	port, err := strconv.Atoi(utils.EnvWithDefault("PORT", "4000"))
	if err != nil {
		return err
	}

	// Create a new HTTP router.
	router := mux.NewRouter()

	// This will serve files under http://localhost:4000/static/<filename>
	router.PathPrefix("/static/").Handler(http.FileServer(http.FS(static)))

	// Add projectService sub router
	projectService, err = project.NewService(&project.ProjectServiceConfig{}, db)
	if err != nil {
		log.Fatalf("error creating ProjectService %s", err.Error())
	}

	project.NewProjectHandler(router.PathPrefix("/projects").Subrouter(), projectService)

	// Add taskService sub router
	taskService, err = task.NewTaskService(&task.TaskServiceConfig{}, db)
	if err != nil {
		log.Fatalf("error creating TaskService %s", err.Error())
	}

	task.NewTaskHandler(router.PathPrefix("/tasks").Subrouter(), taskService, projectService)

	// Handle index page view.
	indexViewHandler := home.NewHomeHandler(taskService, projectService)
	router.HandleFunc("/", indexViewHandler.ServeHTTP).Methods(http.MethodGet)

	if seed {
		database.Seed(20, taskService, projectService)
	}

	router.Use(handlers.RecoveryHandler(handlers.PrintRecoveryStack(true)))

	// Create a new server instance with options from environment variables.
	// For more information, see https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/
	// Note: The ReadTimeout and WriteTimeout settings may interfere with SSE (Server-Sent Event) or WS (WebSocket) connections.
	// For SSE or WS, these timeouts can cause the connection to reset after 10 or 5 seconds due to the ReadTimeout and WriteTimeout settings.
	// If you plan to use SSE or WS, consider commenting out or removing the ReadTimeout and WriteTimeout key-value pairs.
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router,
	}

	// Send log message.
	slog.Info("Starting server...", "port", server.Addr)

	if err := server.ListenAndServe(); err != nil {
		slog.Error("Error starting server", "err", err)
	}

	return nil
}
