package main

import (
	"database/sql"
	"embed"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"maps"
	"math/rand"
	"net/http"
	"os"
	"os/user"
	"slices"
	"strconv"
	"time"

	"github.com/angelofallars/htmx-go"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/guregu/null/v6/zero"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/mattn/go-sqlite3"
	database "github.com/pleimann/camel-do/db"
	"github.com/pleimann/camel-do/services/oauth"
	"github.com/pleimann/camel-do/services/project"
	"github.com/pleimann/camel-do/services/task"
	"github.com/pleimann/camel-do/templates"
	"github.com/pleimann/camel-do/templates/pages"
	"github.com/pleimann/camel-do/utils"
)

var seed bool

func main() {
	flag.BoolVar(&seed, "seed", false, "seed database with some data")
	flag.Parse()

	slog.SetLogLoggerLevel(slog.LevelDebug)

	usr, err := user.Current()
	if err != nil {
		log.Fatalf("Failed to locate the user's home directory: %s\n", err)
	}

	slog.Info("current user", "name", usr.Name, "home-dir", usr.HomeDir)

	if err = createDatabase(); err != nil {
		slog.Error("Failed to create database service!", "details", err.Error())
		os.Exit(1)
	}

	defer db.Close()

	if err = createSyncService(); err != nil {
		slog.Error("Failed create sync service!", "details", err.Error())
		os.Exit(1)
	}

	// Run your server.
	if err = runServer(); err != nil {
		slog.Error("Failed to start server!", "details", err.Error())
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

func createSyncService() error {
	httpClient := oauth.NewGoogleAuth().GetClient()

	var err error
	taskSyncService, err = task.NewTaskSyncService(httpClient, db)
	if err != nil {
		return err
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

	// Handle index page view.
	router.HandleFunc("/", indexViewHandler).Methods(http.MethodGet)

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

	if seed {
		seedData(20)
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

// indexViewHandler handles a view for the index page.
func indexViewHandler(w http.ResponseWriter, r *http.Request) {
	// Check, if the current URL is '/'.
	if r.URL.Path != "/" {
		// If not, return HTTP 404 error.
		http.NotFound(w, r)
		slog.Error("render page", "method", r.Method, "status", http.StatusNotFound, "path", r.URL.Path)
		return
	}

	// Get backlog and tasks scheduled for today
	backlogTasks, err := taskService.GetBacklogTasks()
	if err != nil {
		slog.Error("get backlog tasks", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	todaysTasks, err := taskService.GetTodaysTasks()
	if err != nil {
		slog.Error("get tasks for today", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	projectIndex, err := projectService.GetProjects()
	if err != nil {
		slog.Error("get all projects", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	weekday := time.Now().Weekday()

	// Define template layout for index page.
	indexTemplate := templates.Layout(
		templates.Config{
			Title:    "Camel Do ", // define title text
			LoginUri: "http://localhost:4000/auth/google/login",
		},
		pages.MetaTags(
			"camel-do, todo, tasks", // define meta keywords
			"Welcome to Camel Do! You're here because camels are awesome and you need more of them in your life.", // define meta description
		),
		pages.BodyContent(backlogTasks, weekday, todaysTasks, projectIndex), // define body content
	)

	// Render index page template.
	if err := htmx.NewResponse().RenderTempl(r.Context(), w, indexTemplate); err != nil {
		// If not, return HTTP 400 error.
		w.WriteHeader(http.StatusInternalServerError)
		slog.Error("render template", "method", r.Method, "status", http.StatusInternalServerError, "path", r.URL.Path)
		return
	}

	// Send log message.
	slog.Info("render page", "method", r.Method, "status", http.StatusOK, "path", r.URL.Path)
}

func seedData(count int) {
	projects, err := project.GenerateRandomProjects()
	if err != nil {
		log.Fatal(err)
	}

	tasks, err := task.GenerateRandomTasks(count)
	if err != nil {
		log.Fatal(err)
	}

	for i := range projects {
		err := projectService.AddProject(projects[i])
		if err != nil {
			return
		}
	}

	projectsIndex, _ := projectService.GetProjects()

	projects = slices.Collect(maps.Values(projectsIndex))

	for _, t := range tasks {
		randProject := projects[rand.Intn(len(projects))]

		t.ProjectID = zero.StringFrom(randProject.ID)

		if err := taskService.AddTask(&t); err != nil {
			log.Fatal(err)
		}
	}
}
