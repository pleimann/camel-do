# ADR-0003: HTMX for Server-Side Rendering with Dynamic Updates

## Status
Accepted

## Context
Camel-Do requires a frontend architecture that provides a responsive, interactive user experience while maintaining simplicity and performance. The application needs real-time updates for task management operations (drag & drop, form submissions, dynamic content updates) without the complexity of a full JavaScript framework.

Key requirements:
- Interactive task management interface
- Real-time updates without page refreshes
- Drag & drop functionality for task organization
- Fast page load times
- Minimal JavaScript complexity
- Server-side rendering for SEO and performance
- Progressive enhancement approach

## Decision
We will use HTMX for dynamic interactions combined with server-side rendering using Go's Templ templating system.

Architecture components:
- **Server-Side Rendering**: Templ generates complete HTML on the server
- **HTMX**: Provides AJAX interactions and partial page updates
- **Alpine.js**: Lightweight JavaScript for client-side interactivity
- **Tailwind CSS**: Utility-first CSS framework for styling

The approach provides:
- Complete HTML generated on the server for fast initial loads
- HTMX for seamless partial updates and form submissions
- Alpine.js for complex client-side interactions (drag & drop)
- Progressive enhancement - works without JavaScript

## Consequences

### Positive Consequences
- **Fast Initial Load**: Server-side rendering provides complete HTML immediately
- **Minimal JavaScript**: Reduced client-side complexity and bundle size
- **SEO Friendly**: Full HTML content available for search engines
- **Progressive Enhancement**: Application works with JavaScript disabled
- **Simple State Management**: Server maintains application state
- **Type Safety**: Templ provides compile-time template validation
- **Real-time Feel**: HTMX provides smooth, app-like interactions

### Negative Consequences
- **Server Load**: All rendering happens on server, increases CPU usage
- **Network Dependency**: Most interactions require server round-trips
- **Learning Curve**: Team needs to learn HTMX patterns and concepts
- **Debugging Complexity**: Harder to debug than pure JavaScript applications
- **Limited Offline**: Requires server connectivity for most operations

### Risks
- **HTMX Adoption**: Relatively new technology with smaller community
- **Performance Bottlenecks**: Server rendering could become bottleneck under load
- **Complexity Growth**: May need refactoring if client-side complexity increases

## Alternatives Considered

### Alternative 1: React Single Page Application (SPA)
- Full JavaScript framework with client-side rendering
- Pros: Rich ecosystem, excellent tooling, familiar to many developers
- Cons: Complex build process, large bundle size, SEO challenges
- Why not chosen: Overengineered for task management app, deployment complexity

### Alternative 2: Vue.js with Server-Side Rendering
- Progressive framework with SSR capabilities
- Pros: Gentle learning curve, good performance, flexible architecture
- Cons: Build complexity, requires Node.js toolchain, larger bundle
- Why not chosen: Adds unnecessary complexity for relatively simple UI needs

### Alternative 3: Pure Server-Side Rendering (Traditional Forms)
- Classic form-based web application
- Pros: Very simple, works everywhere, no JavaScript required
- Cons: Poor user experience, full page reloads, no real-time interactions
- Why not chosen: UX requirements demand more interactive experience

### Alternative 4: Stimulus + Turbo (Hotwire)
- Similar approach to HTMX from Ruby on Rails ecosystem
- Pros: Mature approach, good documentation, proven in production
- Cons: Ruby-centric documentation, less Go community adoption
- Why not chosen: HTMX better suited for Go ecosystem, simpler concepts

## Implementation Notes

### Template Structure
```
templates/
├── pages/          -> Full page templates (main.templ, dialogs)
├── blocks/         -> Major UI sections (backlog, tasklist, timeline)
└── components/     -> Reusable UI elements (buttons, forms, modals)
```

### HTMX Integration Patterns
```html
<!-- Form submission with partial update -->
<form hx-post="/tasks" hx-target="#task-list" hx-swap="afterbegin">
    <input name="title" required>
    <button type="submit">Add Task</button>
</form>

<!-- Dynamic content loading -->
<div hx-get="/components/timeline" 
     hx-trigger="every 30s"
     hx-target="this"
     hx-swap="outerHTML">
    Loading...
</div>

<!-- Drag & drop with Alpine.js -->
<div x-data="taskDragDrop()" 
     x-on:drop="handleDrop($event)"
     class="task-container">
    <!-- task elements -->
</div>
```

### Server-Side Rendering with Templ
```go
// Task list component
templ TaskList(tasks []model.Task) {
    <div id="task-list" class="space-y-2">
        for _, task := range tasks {
            @TaskCard(task)
        }
    </div>
}

// HTMX-enabled task card
templ TaskCard(task model.Task) {
    <div class="task-card" 
         hx-get={ fmt.Sprintf("/tasks/%s/edit", task.ID) }
         hx-target="this"
         hx-swap="outerHTML">
        <h3>{ task.Title.String }</h3>
        <p>{ task.Description.String }</p>
    </div>
}
```

### Error Handling
```go
// HTMX-compatible error responses
func customHTTPErrorHandler(err error, c echo.Context) {
    if c.Response().Committed {
        return
    }
    
    code := http.StatusInternalServerError
    if he, ok := err.(*echo.HTTPError); ok {
        code = he.Code
    }
    
    // Render error message as HTMX response
    messageTemplate := components.ErrorMessage(err.Error())
    htmx.NewResponse().
        Retarget("body").
        Reswap("beforeend").
        StatusCode(code).
        RenderTempl(c.Request().Context(), c.Response().Writer, messageTemplate)
}
```

### Performance Optimizations
- **Template Caching**: Compiled templates cached in memory
- **Partial Rendering**: Only render changed portions of UI
- **Efficient Selectors**: Use precise HTMX selectors to minimize DOM manipulation
- **Connection Reuse**: HTTP keep-alive for HTMX requests

### Progressive Enhancement Strategy
1. **Base Functionality**: Works with traditional form submissions
2. **Enhanced Interactions**: HTMX adds seamless updates
3. **Advanced Features**: Alpine.js for drag & drop and complex UI states

## Testing Strategy

### Template Testing
```go
func TestTaskList_Render(t *testing.T) {
    tasks := []model.Task{
        {ID: "1", Title: zero.StringFrom("Test Task")},
    }
    
    component := TaskList(tasks)
    html, err := templ.ToGoHTML(context.Background(), component)
    require.NoError(t, err)
    
    assert.Contains(t, html, "Test Task")
    assert.Contains(t, html, `id="task-list"`)
}
```

### HTMX Integration Testing
```go
func TestTaskHandler_CreateTask_HTMX(t *testing.T) {
    req := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader("title=Test"))
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    req.Header.Set("HX-Request", "true")
    
    rec := httptest.NewRecorder()
    handler.ServeHTTP(rec, req)
    
    assert.Equal(t, http.StatusOK, rec.Code)
    assert.Contains(t, rec.Body.String(), "Test")
}
```

## References
- [HTMX Documentation](https://htmx.org/docs/)
- [Templ Documentation](https://templ.guide/)
- [Alpine.js Documentation](https://alpinejs.dev/)
- [HTMX Examples](https://htmx.org/examples/)

---
*Date: 2025-08-27*
*Authors: Claude Code*
*Reviewers: N/A*