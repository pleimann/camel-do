package main

import (
	"embed"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/user"
	"strconv"
	"time"

	"github.com/angelofallars/htmx-go"
	"github.com/gorilla/mux"
	_ "github.com/joho/godotenv/autoload"
	"github.com/pleimann/camel-do/services/oauth"
	"github.com/pleimann/camel-do/services/task"
	"github.com/pleimann/camel-do/templates"
	"github.com/pleimann/camel-do/templates/pages"
	"github.com/pleimann/camel-do/utils"
)

func main() {
	usr, err := user.Current()
	if err != nil {
		log.Fatalf("Failed to locate the user's home directory: %s\n", err)
	}

	slog.Info("current user", "name", usr.Name, "home-dir", usr.HomeDir)

	// Run your server.
	err = runServer()
	if err != nil {
		slog.Error("Failed to start server!", "details", err.Error())
		os.Exit(1)
	}

	os.Exit(0)
}

//go:embed all:static
var static embed.FS

var taskHandler *task.TaskHandler
var taskService *task.TaskService

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

	// Add taskService sub router
	httpClient := oauth.NewGoogleAuth().GetClient()
	taskService, err = task.NewTaskService(&task.Config{}, httpClient)
	if err != nil {
		log.Fatalf("error creating TaskService %s", err.Error())
	}

	taskHandler = task.NewTaskHandler(router.PathPrefix("/tasks/").Subrouter(), taskService)

	// Create a new server instance with options from environment variables.
	// For more information, see https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/
	// Note: The ReadTimeout and WriteTimeout settings may interfere with SSE (Server-Sent Event) or WS (WebSocket) connections.
	// For SSE or WS, these timeouts can cause the connection to reset after 10 or 5 seconds due to the ReadTimeout and WriteTimeout settings.
	// If you plan to use SSE or WS, consider commenting out or removing the ReadTimeout and WriteTimeout key-value pairs.
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      router,
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
		pages.BodyContent(taskService.GetTasks()), // define body content
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
