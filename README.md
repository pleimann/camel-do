# Camel-Do

A comprehensive task management application built with Go that seamlessly integrates with Google Calendar. Camel-Do provides a modern, responsive interface for organizing projects, scheduling tasks, and maintaining productivity workflows with real-time synchronization across platforms.

## Features

### Core Task Management
- **Task Creation & Editing**: Rich task forms with title, description, duration, and completion tracking
- **Task Scheduling**: Interactive date/time picker with calendar interface for precise scheduling
- **Task Status Management**: Mark tasks as completed, hidden, or prioritized with ranking system
- **Task Views**: Multiple display modes including backlog cards and structured task lists

### Project Organization
- **Project-based Grouping**: Organize tasks into customizable projects with unique identifiers
- **Visual Customization**: Choose from predefined colors (Zinc, Red, Orange, Amber, Yellow, Lime, Green, Emerald, Teal, Cyan, Sky, Blue, Indigo, Violet, Purple, Fuchsia, Pink, Rose, Stone, Neutral, Slate, Gray) and icons for project identification
- **Project Management**: Full CRUD operations for creating, editing, and deleting projects

### Google Integration
- **OAuth2 Authentication**: Secure Google account integration with automatic token management
- **Calendar Synchronization**: Bi-directional sync between tasks and Google Calendar events
- **Event Management**: Create, update, and sync calendar events from scheduled tasks
- **Google Tasks Integration**: Connect with Google Tasks API for cross-platform compatibility

### User Interface & Experience
- **Real-time Updates**: Dynamic interface powered by HTMX without full page refreshes
- **Responsive Design**: Modern UI built with Tailwind CSS and DaisyUI components
- **Interactive Components**: Custom date pickers, time selectors, and dialog modals
- **Timeline Visualization**: Visual representation of tasks across time periods
- **Dashboard Overview**: Comprehensive view of today's tasks and project status

### Data Management
- **Local Storage**: Embedded BoltDB for fast, reliable local data persistence
- **Automatic Migrations**: Database schema updates handled automatically
- **Data Seeding**: Optional test data generation for development and demonstration
- **Backup & Sync**: Tasks synchronized with Google services for data redundancy

## Quick Start

This README contains information about:

- [Project overview](#project-overview)
- [Getting started](#getting-started)
- [Development](#development)
- [Project structure](#project-structure)
- [Deployment](#deployment)

## Project overview

Backend:

- Module name in the go.mod file: `github.com/pleimann/camel-do`
- Go web framework/router: `Echo`
- Server port: `7000`

Frontend:

- Package name in the package.json file: `camel-do`
- Reactivity library: `htmx with Alpine.js`
- CSS framework: `Tailwind CSS with Preline UI components`

Tools:

- Air tool to live-reloading: ✓
- Bun as a frontend runtime: ✓
- Templ to generate HTML: ✓
- Config for golangci-lint: ✓

## Project structure

```console
.
├── main.go                    # Application entry point
├── go.mod/go.sum             # Go dependencies
├── package.json              # Frontend dependencies
├── .air.toml                 # Live reload config
├── docker-compose.yml        # Container setup
├── services/                 # Business logic layer
│   ├── home/                 # Dashboard handlers
│   ├── project/              # Project management
│   ├── task/                 # Task operations & sync
│   ├── cal/                  # Calendar integration
│   └── oauth/                # Authentication
├── db/                       # Database layer
│   ├── migrations/           # SQL migrations
│   ├── model/               # Database entities
│   └── table/               # Generated Jet tables
├── model/                    # Domain models
├── templates/                # Templ HTML templates
│   ├── pages/               # Full page templates
│   ├── blocks/              # UI sections
│   └── components/          # Reusable components
├── assets/                   # Source frontend files
│   ├── scripts.js
│   └── styles.css
├── static/                   # Built assets
└── utils/                    # Shared utilities
```

## Getting started

> ❗️ Please make sure that you have installed the executable files for all the necessary tools before starting your project. Exactly:
>
> - `Air`: [https://github.com/air-verse/air](https://github.com/air-verse/air)
> - `Bun`: [https://github.com/oven-sh/bun](https://github.com/oven-sh/bun)
> - `Templ`: [https://github.com/a-h/templ](https://github.com/a-h/templ)
> - `golangci-lint`: [https://github.com/golangci/golangci-lint](https://github.com/golangci/golangci-lint)

To start your project, run the **Gowebly** CLI command in your terminal:

```console
gowebly run
```

## Development

### Architecture

Camel-Do follows a service-oriented architecture:

- **Services Layer**: Domain-specific business logic (task, project, calendar, auth)
- **Database Layer**: SQLite with type-safe Jet queries and migrations
- **Template Layer**: Server-side rendered HTML using Templ
- **Frontend**: HTMX + Alpine.js for dynamic interactions, styled with Tailwind CSS

### Key Technologies

- **Backend**: Go 1.24+ with Echo web framework
- **Database**: SQLite with automatic migrations
- **Templates**: Templ for type-safe HTML generation
- **Frontend**: HTMX, Alpine.js, Tailwind CSS, DaisyUI
- **Build Tools**: Air (live reload), Bun (frontend), Parcel (bundling)

### Development Workflow

1. **Start development server**: `gowebly run` or `air`
2. **Watch frontend assets**: `bun run watch`
3. **Generate templates**: `templ generate` (automatic with Air)
4. **Run tests**: `go test ./...`
5. **Format code**: `bun run fmt` and `golangci-lint run`

## Deploying your project

All deploy settings are located in the `Dockerfile` and `docker-compose.yml` files in your project folder.

To deploy your project to a remote server, follow these steps:

1. Go to your hosting/cloud provider and create a new VDS/VPS.
2. Update all OS packages on the server and install Docker, Docker Compose and Git packages.
3. Use `git clone` command to clone the repository with your project to the server and navigate to its folder.
4. Run the `docker-compose up` command to start your project on your server.

> ❗️ Don't forget to generate Go files from `*.templ` templates before run the `docker-compose up` command.
