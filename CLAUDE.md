# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

### Core Development
- **Start development server**: `gowebly run` or `air` (preferred for live reload)
- **Build frontend assets**: `bun run build` (production) or `bun run dev` (development)
- **Watch frontend assets**: `bun run watch` (continuous development)
- **Generate templates**: `templ generate` (auto-generated with Air)
- **Format code**: `bun run fmt` (runs Prettier)
- **Run tests**: `go test ./...`
- **Build binary**: `go build -o ./tmp/camel-do .`
- **Lint code**: `golangci-lint run`

### Database Operations
- **Run with seeded data**: `go run . -seed`

## Architecture Overview

### Service-Oriented Design
Camel-Do uses a layered service architecture with clear separation of concerns:

- **Services Layer** (`/services/`): Domain-specific business logic
  - `home/`: Dashboard and main page handlers
  - `task/`: Task CRUD operations, sync with Google Calendar
  - `project/`: Project management and organization
  - `cal/`: Google Calendar integration
  - `oauth/`: Authentication and OAuth flows

- **Database Layer**: BoltDB embedded database with direct bucket operations
  - Uses BoltDB (embedded key-value store) with bucket-based storage
  - Tasks stored in "tasks" bucket with gob encoding/decoding
  - No separate ORM or migration system - manual bucket management
  - Domain models in `/model/` directory

- **Template Layer** (`/templates/`): Server-side rendered HTML
  - `pages/`: Full page templates
  - `blocks/`: UI sections (backlog, tasklist, timeline, titlebar)
  - `components/`: Reusable UI components

### Key Technologies
- **Backend**: Go 1.25+ with Echo web framework, BoltDB embedded database
- **Frontend**: HTMX + Alpine.js for reactivity, Tailwind CSS + DaisyUI for styling
- **Templates**: Templ for type-safe HTML generation
- **Build**: Air for live reload, Bun for frontend bundling, Parcel for asset processing

### Google Calendar Integration
- OAuth2 authentication flow in `services/oauth/`
- Bi-directional sync between tasks and calendar events
- Task sync service handles creating/updating calendar events from tasks

### Database Schema
- Tasks stored in BoltDB buckets (embedded key-value store)
- Task model includes Google Calendar event IDs for sync (GTaskID field)
- Projects stored separately for task organization
- Data marshaled/unmarshaled using gob encoding
- Model definitions in `/model/` (task.go, project.go, event.go, etc.)

## Development Workflow

1. **Air Configuration**: `.air.toml` handles live reloading with pre-build steps:
   - Runs `bun run dev` to build frontend assets
   - Runs `go generate` for code generation
   - Runs `go tool templ generate` to generate Go code from templates
   - Excludes test files and generated files from watching

2. **Frontend Development**: 
   - Source files in `/assets/`
   - Built assets output to `/static/`
   - HTMX for dynamic interactions without JavaScript frameworks

3. **Template Development**:
   - Write `.templ` files which generate `_templ.go` files
   - Templates are type-safe and compiled into Go code

4. **Testing**:
   - Unit tests use standard Go testing
   - Example test structure in `templates/components/timepicker_templ_test.go`

## Project Structure Notes

- **Embedded Assets**: Static files and credentials are embedded in the binary (`//go:embed` directives in main.go)
- **Configuration**: Uses environment variables with sensible defaults (PORT=4000, proxy on 4001)
- **Database Path**: BoltDB file stored in user config directory at `~/.config/camel-do/camel-do.db` (or equivalent on Windows/macOS)
- **Service Initialization**: All services initialized in main.go with dependency injection pattern
- **Error Handling**: Custom HTMX error handler (`customHTTPErrorHandler`) for dynamic error display
- **Routing Structure**:
  - Main routes in main.go
  - Service-specific routes registered via handler groups (`/projects`, `/tasks`, `/timeline`, `/components`)
- **Model Enums**: Auto-generated enums for colors and icons using go-enum tool
- **Commit Convention**: Use conventional commits