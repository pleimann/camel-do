# Agent Guidelines for Camel-Do

## Build, Lint, and Test Commands

### Development Server

- `air` - Live reload development server (preferred)
- `gowebly run` - Alternative development server

### Building

- `go build -o ./tmp/camel-do .` - Build binary
- `bun run build` - Build frontend assets (production)
- `bun run dev` - Build frontend assets (development, no optimization)
- `bun run watch` - Watch and rebuild frontend assets

### Code Generation

- `go generate` - Run go generate directives
- `go tool templ generate` - Generate Go code from .templ files

### Testing

- `go test ./...` - Run all tests
- `go test -v ./services/task` - Run tests for specific package
- `go test -run TestSpecificFunction` - Run single test function

### Linting and Formatting

- `golangci-lint run` - Run comprehensive Go linting
- `bun run fmt` - Format frontend code with Prettier
- `go fmt ./...` - Format Go code

## Code Style Guidelines

### Go Code Style

- **Formatting**: 4 spaces indentation, LF line endings
- **Imports**: Grouped as (standard library, third-party, local packages)
- **Naming**: PascalCase for exported, camelCase for unexported
- **Error Handling**: Wrap errors with `fmt.Errorf("context: %w", err)`
- **Logging**: Use `slog` for structured logging with appropriate levels
- **Function Length**: Max 120 lines per function
- **Line Length**: Max 300 characters
- **Types**: Use `any` instead of `interface{}`, prefer concrete types

### Frontend Code Style

- **Formatting**: 2 spaces indentation, single quotes, no semicolons
- **Trailing Commas**: ES5 style in objects/arrays
- **JavaScript**: Alpine.js for reactivity, HTMX for interactions

### Architecture Patterns

- **Service Layer**: Business logic in `/services/` packages
- **Model Layer**: Data structures in `/model/` with JSON marshaling
- **Handler Layer**: HTTP handlers with Echo framework
- **Template Layer**: Server-side rendering with Templ

### Database Patterns

- **BoltDB**: Embedded key-value store with bucket organization
- **Transactions**: Use `db.Update()` for writes, `db.View()` for reads
- **Serialization**: GOB encoding for complex structs

### Testing Patterns

- **Unit Tests**: Test individual functions/methods
- **Integration Tests**: Test service interactions
- **Test Files**: `*_test.go` alongside implementation files

### Commit Conventions

- Use conventional commits: `feat:`, `fix:`, `docs:`, `refactor:`, etc.
- Include scope when relevant: `feat(task): add priority field`

### Security

- Never log or expose secrets/keys
- Validate all user inputs
- Use proper error handling without information leakage

## Known Issues to Address

### Compilation Errors

- `services/task/taskhandlers.go`: Missing methods `GetEventsForDate` and `GetTasksForDate`
- `model/task.go`: Unused variable `endHours`

### Code Quality

- Run `golangci-lint run` to identify and fix linting issues
- Address all compilation errors before committing changes
