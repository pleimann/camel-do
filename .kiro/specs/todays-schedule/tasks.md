# Implementation Plan

- [-] 1. Add database support for task hiding functionality
  - Add `hidden` column to tasks table via database migration
  - Update task model to include Hidden field with proper zero value handling
  - _Requirements: 3.3, 3.4, 3.5_

- [ ] 2. Extend TaskService with new methods for schedule functionality
  - Implement `HideTask(id string) error` method to mark tasks as hidden
  - Implement `UnhideTask(id string) error` method to restore hidden tasks
  - Implement `GetTodaysVisibleTasks() (*model.TaskList, error)` method to exclude hidden tasks
  - Implement `RescheduleTask(id string, newStartTime time.Time) error` method for time updates
  - Write comprehensive unit tests for all new service methods
  - _Requirements: 3.1, 3.2, 4.3, 4.4, 4.5_

- [ ] 3. Create schedule timeline template components
  - Create `templates/blocks/schedule/schedule.templ` with timeline grid layout
  - Implement time markers display with proper formatting and current time indicator
  - Add responsive layout handling for mobile and desktop views
  - Include proper HTMX attributes for dynamic updates
  - _Requirements: 1.1, 1.2, 1.3, 5.1, 5.2, 5.3, 6.1, 6.2, 6.3_

- [ ] 4. Create interactive task card component
  - Create `templates/blocks/schedule/task_card.templ` with completion toggle
  - Add hide button with proper HTMX integration
  - Add reschedule button that triggers modal
  - Implement visual states for completed, overdue, and normal tasks
  - Ensure touch-friendly controls for mobile devices
  - _Requirements: 2.1, 2.2, 2.3, 3.1, 4.1, 1.5, 6.4_

- [ ] 5. Create reschedule modal component
  - Create `templates/components/reschedule_modal.templ` with date/time picker
  - Implement quick time slot selection for same-day rescheduling
  - Add form validation for date/time inputs
  - Integrate with existing modal system for consistent behavior
  - _Requirements: 4.2, 4.3_

- [ ] 6. Add new HTTP handlers for schedule actions
  - Add `PUT /tasks/:id/hide` endpoint in TaskHandler
  - Add `PUT /tasks/:id/unhide` endpoint in TaskHandler  
  - Add `PUT /tasks/:id/reschedule` endpoint in TaskHandler
  - Implement proper error handling and HTTP status codes
  - Add request validation and user authorization checks
  - _Requirements: 3.2, 3.6, 4.4, 4.6_

- [ ] 7. Update existing handlers to support hidden tasks
  - Modify `GetTodaysTasks` usage to use `GetTodaysVisibleTasks` where appropriate
  - Update task completion handler to work with schedule timeline
  - Ensure proper HTMX responses for schedule component updates
  - _Requirements: 2.4, 2.5, 3.4_

- [ ] 8. Implement timeline positioning and layout logic
  - Add logic for handling overlapping tasks in timeline display
  - Implement proper spacing and visual grouping for time slots
  - Add current time indicator positioning
  - Handle edge cases like tasks spanning multiple time slots
  - _Requirements: 5.4, 5.5, 1.4_

- [ ] 9. Add schedule component to main application
  - Integrate schedule timeline into main page layout
  - Add navigation or toggle between different views (backlog, schedule, etc.)
  - Ensure proper initialization and data loading
  - Test integration with existing components
  - _Requirements: 1.1_

- [ ] 10. Implement responsive design and mobile optimizations
  - Add CSS media queries for mobile layout adaptations
  - Implement touch gesture support for common actions
  - Optimize scrolling behavior for timeline navigation
  - Test and adjust control sizes for touch accessibility
  - _Requirements: 6.1, 6.2, 6.3, 6.4, 6.5_

- [ ] 11. Add comprehensive error handling and user feedback
  - Implement client-side error display for failed actions
  - Add loading states for async operations
  - Create user-friendly error messages for common scenarios
  - Add confirmation dialogs for destructive actions
  - _Requirements: 2.4, 3.4, 4.6_

- [ ] 12. Write integration tests for schedule functionality
  - Create tests for complete user workflows (complete → hide → reschedule)
  - Test HTMX interactions and partial page updates
  - Test concurrent user actions and data consistency
  - Add tests for edge cases and error scenarios
  - _Requirements: All requirements validation_