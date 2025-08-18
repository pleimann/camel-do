# Design Document

## Overview

The Today's Schedule component will be a new timeline-based view that displays all tasks scheduled for the current day in chronological order. This component will extend the existing timeline functionality to provide interactive task management capabilities including completion toggling, hiding tasks, and rescheduling directly from the timeline view.

The component will integrate with the existing task service architecture and leverage the current HTMX-based frontend approach for real-time updates without page refreshes.

## Architecture

### Component Structure
The Today's Schedule component will be built as a new template block under `templates/blocks/schedule/` that extends the existing timeline functionality. It will reuse the existing timeline grid layout but add interactive controls for each task.

### Service Layer Integration
The component will utilize the existing `TaskService` for all data operations and extend it with new methods for task hiding functionality. No changes to the core task model are required as the existing fields support all necessary operations.

### Frontend Architecture
Following the existing HTMX pattern, the component will use:
- Server-side rendered templates with Templ
- HTMX for dynamic interactions and partial page updates
- Alpine.js for client-side state management where needed
- Tailwind CSS with DaisyUI for styling consistency

## Components and Interfaces

### New Template Components

#### 1. Schedule Timeline Block (`templates/blocks/schedule/schedule.templ`)
- Main container for the today's schedule view
- Extends existing timeline grid layout
- Renders time markers and task cards with interactive controls
- Handles responsive layout for different screen sizes

#### 2. Interactive Task Card (`templates/blocks/schedule/task_card.templ`)
- Enhanced version of existing timeline task card
- Includes completion toggle, hide button, and reschedule controls
- Maintains existing visual design with added interaction elements
- Supports touch-friendly controls for mobile devices

#### 3. Reschedule Modal (`templates/components/reschedule_modal.templ`)
- Date/time picker interface for rescheduling tasks
- Integrates with existing modal component system
- Provides quick time slot selection for same-day rescheduling

### Service Layer Extensions

#### TaskService Methods
```go
// Hide a task from timeline view (soft delete approach)
func (t *TaskService) HideTask(id string) error

// Unhide a previously hidden task
func (t *TaskService) UnhideTask(id string) error

// Get today's tasks excluding hidden ones
func (t *TaskService) GetTodaysVisibleTasks() (*model.TaskList, error)

// Reschedule task to new date/time
func (t *TaskService) RescheduleTask(id string, newStartTime time.Time) error
```

#### Handler Extensions
New HTTP endpoints in `TaskHandler`:
- `PUT /tasks/:id/hide` - Hide task from timeline
- `PUT /tasks/:id/unhide` - Unhide task
- `PUT /tasks/:id/reschedule` - Reschedule task to new time

### Database Schema Considerations

The existing task schema supports all required functionality:
- `completed` field for completion status
- `start_time` field for scheduling
- A new `hidden` boolean field will be added to support hiding functionality

## Data Models

### Task Model Extensions
```go
type Task struct {
    // ... existing fields
    Hidden zero.Bool `form:"hidden,default:false"` // New field for hiding tasks
}
```

### Timeline Position
The existing `TimelinePosition` struct will be reused for positioning tasks in the timeline grid.

### Task Actions
```go
type TaskAction struct {
    Type   string // "complete", "hide", "reschedule"
    TaskID string
    Data   map[string]interface{} // Additional data for specific actions
}
```

## Error Handling

### Client-Side Error Handling
- HTMX error responses will display user-friendly error messages
- Failed actions will revert UI state and show appropriate feedback
- Network errors will be handled gracefully with retry options

### Server-Side Error Handling
- Database operation failures will return appropriate HTTP status codes
- Validation errors will provide specific field-level feedback
- Concurrent modification conflicts will be handled with optimistic locking

### Error Recovery
- Failed completion toggles will revert checkbox state
- Failed hide operations will restore task visibility
- Failed reschedule operations will maintain original time

## Testing Strategy

### Unit Tests
- Test all new TaskService methods with various scenarios
- Test task action handlers with valid and invalid inputs
- Test timeline positioning logic with edge cases
- Test responsive layout behavior

### Integration Tests
- Test complete user workflows (complete → hide → reschedule)
- Test HTMX interactions and partial page updates
- Test concurrent user actions on the same tasks
- Test database transaction handling

### End-to-End Tests
- Test full timeline interaction flows
- Test mobile responsiveness and touch interactions
- Test accessibility compliance with screen readers
- Test performance with large numbers of tasks

### Test Data Scenarios
- Tasks with various durations and overlapping times
- Tasks without specific start times
- Tasks scheduled for different time zones
- Edge cases like midnight transitions and daylight saving time

## Implementation Considerations

### Performance Optimization
- Limit timeline to current day only to reduce data load
- Use efficient database queries with proper indexing
- Implement client-side caching for frequently accessed data
- Optimize template rendering for large task lists

### Accessibility
- Ensure all interactive controls are keyboard accessible
- Provide proper ARIA labels for screen readers
- Maintain sufficient color contrast for task indicators
- Support high contrast mode and reduced motion preferences

### Mobile Responsiveness
- Adapt timeline layout for smaller screens
- Provide touch-friendly control sizes (minimum 44px)
- Implement swipe gestures for common actions
- Optimize scrolling behavior for mobile devices

### Browser Compatibility
- Support modern browsers with HTMX and Alpine.js
- Provide graceful degradation for older browsers
- Test across different operating systems and devices
- Ensure consistent behavior across browser engines

## Security Considerations

### Authorization
- Verify user ownership of tasks before allowing modifications
- Implement proper session management for task operations
- Validate all user inputs to prevent injection attacks
- Use CSRF protection for state-changing operations

### Data Validation
- Validate date/time inputs for rescheduling operations
- Sanitize task titles and descriptions for XSS prevention
- Implement rate limiting for rapid task modifications
- Validate task state transitions (e.g., can't complete hidden tasks)

## Integration Points

### Existing Components
- Reuse existing modal system for reschedule interface
- Integrate with existing project color and icon system
- Maintain consistency with existing task card styling
- Leverage existing error handling and notification system

### External Services
- Sync task completion status with Google Calendar integration
- Update external calendar events when tasks are rescheduled
- Handle conflicts between local changes and external sync
- Maintain data consistency across all integrated services