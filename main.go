package main

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/angelofallars/htmx-go"
	"github.com/gorilla/mux"
	_ "github.com/joho/godotenv/autoload"
	"github.com/pleimann/camel-do/services/task"
	"github.com/pleimann/camel-do/templates"
	"github.com/pleimann/camel-do/templates/pages"
	"github.com/pleimann/camel-do/utils"
)

func main() {
	// Run your server.
	if server, err := runServer(); err != nil {
		slog.Error("Failed to start server!", "details", err.Error())
		os.Exit(1)

	} else {
		var wait time.Duration
		flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
		flag.Parse()

		c := make(chan os.Signal, 1)
		// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
		// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
		signal.Notify(c, os.Interrupt)

		// Block until we receive our signal.
		<-c

		// Create a deadline to wait for.
		ctx, cancel := context.WithTimeout(context.Background(), wait)
		defer cancel()

		// Doesn't block if no connections, but will otherwise wait
		// until the timeout deadline.
		go server.Shutdown(ctx)

		<-ctx.Done()

		// Optionally, you could run srv.Shutdown in a goroutine and block on
		// <-ctx.Done() if your application should wait for other services
		// to finalize based on context cancellation.
		log.Println("shutting down")

		os.Exit(0)
	}
}

var taskService = task.NewTaskService(&task.Config{})

//go:embed all:static
var static embed.FS

// runServer runs a new HTTP server with the loaded environment variables.
func runServer() (*http.Server, error) {
	port, err := strconv.Atoi(utils.EnvWithDefault("PORT", "4000"))
	if err != nil {
		return nil, err
	}

	// Create a new HTTP router.
	router := mux.NewRouter()

	// This will serve files under http://localhost:8000/static/<filename>
	router.PathPrefix("/static/").Handler(http.FileServer(http.FS(static)))

	// Handle index page view.
	router.HandleFunc("/", indexViewHandler).Methods("GET")

	// Add taskService sub router
	tasksRouter := router.PathPrefix("/tasks/").Subrouter()
	task.Routes(tasksRouter, taskService)

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
	slog.Info("Starting server...", "port", port)

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	return server, nil
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

	// Define template meta tags.
	metaTags := pages.MetaTags(
		"camel-do, todo, tasks", // define meta keywords
		"Welcome to Camel Do! You're here because camels are awesome and you need more of them in your life.", // define meta description
	)

	// Define template body content.
	bodyContent := pages.BodyContent(taskService.GetTasks())

	// Define template layout for index page.
	indexTemplate := templates.Layout(
		"Camel Do ", // define title text
		metaTags,    // define meta tags
		bodyContent, // define body content
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
