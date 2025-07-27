# Technology Stack

## Backend
- **Language**: Go 1.24+
- **Web Framework**: Echo v4 (HTTP router and middleware)
- **Database**: SQLite with go-sqlite3 driver
- **Database Migrations**: golang-migrate/migrate
- **Database Query Builder**: go-jet/jet for type-safe SQL
- **Template Engine**: Templ for HTML generation
- **Authentication**: OAuth2 with Google API integration

## Frontend
- **Reactivity**: HTMX with Alpine.js for dynamic interactions
- **CSS Framework**: Tailwind CSS v4 with DaisyUI components
- **Icons**: Lucide icons
- **Build Tool**: Parcel for asset bundling
- **Runtime**: Bun for frontend package management

## Development Tools
- **Live Reload**: Air for Go backend hot reloading
- **Code Generation**: Templ CLI for HTML template compilation
- **Linting**: golangci-lint with custom configuration
- **Formatting**: Prettier for frontend code
- **Package Management**: Go modules + Bun

## Common Commands

### Development
```bash
# Start development server with live reload
gowebly run

# Or manually:
air  # Starts backend with hot reload
bun run watch  # Watches frontend assets

# Generate templates
templ generate

# Build frontend assets
bun run build  # Production build
bun run dev    # Development build
```

### Database
```bash
# Run with database seeding
go run . -seed

# Database migrations are handled automatically on startup
```

### Testing
```bash
# Run Go tests
go test ./...

# Run specific service tests
go test ./services/task/...
```

### Deployment
```bash
# Build Docker image
docker-compose up

# Format code
bun run fmt  # Frontend formatting
golangci-lint run  # Go linting
```

## Key Dependencies
- `github.com/labstack/echo/v4` - Web framework
- `github.com/a-h/templ` - Template engine
- `github.com/go-jet/jet/v2` - SQL query builder
- `github.com/angelofallars/htmx-go` - HTMX integration
- `golang.org/x/oauth2` - OAuth2 client
- `google.golang.org/api` - Google API client