# ADR-0001: Service-Oriented Architecture

## Status
Accepted

## Context
The Camel-Do task management application requires a clear architectural pattern to organize business logic, maintain separation of concerns, and enable testability. We need to decide how to structure the codebase to support features like task management, project organization, Google Calendar integration, and OAuth authentication.

## Decision
We will implement a service-oriented architecture with the following layers:

1. **Services Layer** (`/services/`) - Domain-specific business logic
2. **Database Layer** (`/model/` and data access) - Data persistence with type-safe operations
3. **Template Layer** (`/templates/`) - Server-side rendered HTML using Templ
4. **Handler Layer** - HTTP request/response handling

Each service encapsulates a specific domain:
- `task/` - Task CRUD operations and Google Calendar sync
- `project/` - Project management and organization  
- `cal/` - Google Calendar integration
- `oauth/` - Authentication and OAuth flows
- `home/` - Dashboard and main page logic

## Consequences

### Positive Consequences
- Clear separation of concerns with domain boundaries
- Services are independently testable and maintainable
- Business logic is isolated from HTTP handling concerns
- Supports dependency injection for better testing
- Enables gradual refactoring and evolution of individual services

### Negative Consequences
- Slightly more complex initial setup compared to monolithic structure
- Requires careful management of service dependencies
- May lead to over-engineering for simple CRUD operations

### Risks
- **Service boundaries may become unclear** - Mitigation: Regular architecture reviews and clear domain definitions
- **Circular dependencies between services** - Mitigation: Strict dependency injection pattern and interface definitions

## Alternatives Considered

### Alternative 1: Monolithic Structure
- Single package with all business logic
- Pros: Simpler initial setup, fewer files
- Cons: Poor separation of concerns, difficult to test, harder to maintain as codebase grows
- Why not chosen: Doesn't scale well and makes testing individual components difficult

### Alternative 2: Hexagonal Architecture (Ports & Adapters)
- Strict separation between core domain and external adapters
- Pros: Very clean architecture, excellent testability
- Cons: Higher complexity overhead for a relatively simple application
- Why not chosen: Over-engineered for current requirements, can evolve to this pattern later

## Implementation Notes
- All services accept their dependencies via constructor injection
- Services expose interfaces for better testability
- Database access is encapsulated within services, not exposed directly to handlers
- Each service has its own test file following Go conventions
- Services are initialized in `main.go` with proper dependency wiring

## References
- [Clean Architecture by Robert Martin](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Go project layout](https://github.com/golang-standards/project-layout)
- [Dependency Injection in Go](https://blog.drewolson.org/dependency-injection-in-go)

---
*Date: 2025-08-27*
*Authors: Claude Code*
*Reviewers: N/A*