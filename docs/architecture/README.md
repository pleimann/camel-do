# Architecture Documentation

This directory contains comprehensive architecture documentation for the Camel-Do task management application.

## Documentation Structure

- [System Context](01-system-context.md) - High-level system overview and external integrations
- [Container Architecture](02-container-architecture.md) - Service architecture and deployment view  
- [Component Architecture](03-component-architecture.md) - Internal module structure and relationships
- [Data Architecture](04-data-architecture.md) - Data models and persistence patterns
- [Security Architecture](05-security-architecture.md) - Security patterns and threat model
- [Quality Attributes](06-quality-attributes.md) - Cross-cutting concerns and non-functional requirements
- [Architecture Decision Records](adrs/) - Historical architectural decisions and rationale

## Architecture Overview

Camel-Do is a modern task management application with Google Calendar integration built using:

- **Backend**: Go 1.24+ with Echo web framework
- **Database**: BoltDB (embedded key-value store) 
- **Frontend**: HTMX + Alpine.js with server-side rendering
- **Templates**: Templ for type-safe HTML generation
- **Styling**: Tailwind CSS with DaisyUI components
- **Build Tools**: Air (live reload), Bun (frontend bundling)

## Key Architectural Principles

1. **Service-Oriented Architecture**: Clear separation of concerns with domain-specific services
2. **Server-Side Rendering**: HTMX for dynamic interactions without heavy JavaScript frameworks
3. **Type Safety**: Leveraging Go's type system and Templ for template safety
4. **Embedded Dependencies**: Self-contained application with embedded assets and database
5. **Event-Driven Sync**: Bi-directional synchronization with Google Calendar

## Documentation Standards

- Diagrams are created using Mermaid syntax for version control compatibility
- Architecture Decision Records follow the template in `adrs/template.md`
- All documentation follows the C4 model hierarchy (Context → Containers → Components → Code)