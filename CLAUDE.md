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

### Database Operations
- **Run with seeded data**: `go run . -seed`
- **Database migrations**: Automatically applied on startup via `database.Migrate(db)`

## Architecture Overview

### Service-Oriented Design
Camel-Do uses a layered service architecture with clear separation of concerns:

- **Services Layer** (`/services/`): Domain-specific business logic
  - `home/`: Dashboard and main page handlers
  - `task/`: Task CRUD operations, sync with Google Calendar
  - `project/`: Project management and organization
  - `cal/`: Google Calendar integration
  - `oauth/`: Authentication and OAuth flows

- **Database Layer** (`/db/`): Data persistence with type-safe queries
  - SQLite database with automatic migrations
  - Jet ORM for type-safe SQL generation
  - Separate model structs in `/db/model/` and domain models in `/model/`

- **Template Layer** (`/templates/`): Server-side rendered HTML
  - `pages/`: Full page templates
  - `blocks/`: UI sections (backlog, tasklist, timeline, titlebar)
  - `components/`: Reusable UI components

### Key Technologies
- **Backend**: Go 1.24+ with Echo web framework, SQLite database
- **Frontend**: HTMX + Alpine.js for reactivity, Tailwind CSS + DaisyUI for styling
- **Templates**: Templ for type-safe HTML generation
- **Build**: Air for live reload, Bun for frontend bundling, Parcel for asset processing

### Google Calendar Integration
- OAuth2 authentication flow in `services/oauth/`
- Bi-directional sync between tasks and calendar events
- Task sync service handles creating/updating calendar events from tasks

### Database Schema
- Tasks table with Google Calendar event IDs for sync
- Projects table for task organization
- Database migrations in `/db/migrations/`
- Jet-generated table definitions in `/db/table/`

## Development Workflow

1. **Air Configuration**: `.air.toml` handles live reloading with pre-build steps:
   - Runs `bun run dev` to build frontend assets
   - Generates Go code from templates
   - Excludes test files and generated files from watching

2. **Frontend Development**: 
   - Source files in `/assets/`
   - Built assets output to `/static/`
   - HTMX for dynamic interactions without JavaScript frameworks

3. **Template Development**:
   - Write `.templ` files which generate `_templ.go` files
   - Templates are type-safe and compiled into Go code

4. **Testing**:
   - Unit tests use testify for assertions
   - Example test structure in `services/task/taskservice_test.go`

## Project Structure Notes

- **Embedded Assets**: Static files and credentials are embedded in the binary
- **Configuration**: Uses environment variables with sensible defaults (PORT=4000)
- **Database Path**: Local SQLite file in `./camel-do/camel-do.db`
- **Service Initialization**: All services initialized in main.go with dependency injection
- **Error Handling**: Custom HTMX error handler for dynamic error display