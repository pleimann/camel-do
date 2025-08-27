# ADR-0004: Google Calendar Integration with OAuth 2.0

## Status
Accepted

## Context
Camel-Do needs to integrate with Google Calendar to provide bi-directional synchronization between tasks and calendar events. This requires secure authentication and authorization to access user's calendar data while maintaining user privacy and security.

Key requirements:
- Secure access to Google Calendar API
- User consent and authorization management
- Token refresh and lifecycle management
- Bi-directional synchronization (tasks ↔ calendar events)
- Minimal required permissions (principle of least privilege)
- Offline operation when calendar sync unavailable

Google provides OAuth 2.0 for API access with the following considerations:
- Authorization Code flow required for web applications
- PKCE (Proof Key for Code Exchange) recommended for security
- Token refresh capability needed for long-term access
- Scope limitations for user privacy

## Decision
We will implement Google Calendar integration using OAuth 2.0 Authorization Code flow with the following architecture:

1. **OAuth 2.0 Flow**: Authorization Code flow with PKCE for enhanced security
2. **Scope Management**: Request minimal required scopes (`https://www.googleapis.com/auth/calendar`)
3. **Token Storage**: Encrypted storage of access and refresh tokens in BoltDB
4. **Sync Service**: Dedicated service for bi-directional synchronization
5. **Error Handling**: Graceful degradation when calendar service unavailable

### Integration Architecture
```
User Authentication Flow:
1. User initiates OAuth flow
2. Redirect to Google OAuth consent screen
3. User grants calendar access permission
4. Google redirects back with authorization code
5. Exchange code for access/refresh tokens
6. Store encrypted tokens locally
7. Use tokens for Calendar API access

Sync Flow:
1. Task created/updated → Create/update calendar event
2. Calendar event changed → Update corresponding task
3. Conflict resolution with last-modified-wins strategy
```

## Consequences

### Positive Consequences
- **Secure Authentication**: Industry-standard OAuth 2.0 provides robust security
- **User Control**: Users explicitly grant permissions and can revoke access
- **Bi-directional Sync**: Seamless integration between task management and calendar
- **Token Management**: Automatic token refresh prevents authentication failures
- **Privacy Focused**: Minimal scope requests reduce privacy concerns
- **Offline Capability**: Local task management continues without calendar connectivity

### Negative Consequences
- **Complexity**: OAuth 2.0 flow adds implementation and testing complexity
- **Network Dependency**: Calendar sync requires internet connectivity
- **Google Dependency**: Relies on Google's service availability and API stability
- **Rate Limiting**: Subject to Google's API quotas and rate limits
- **Token Security**: Requires secure token storage and handling

### Risks
- **Token Compromise**: Mitigation through encryption at rest and secure handling
- **API Changes**: Mitigation through version pinning and update procedures
- **Service Outages**: Mitigation through graceful degradation and retry logic
- **Rate Limiting**: Mitigation through intelligent sync scheduling and backoff

## Alternatives Considered

### Alternative 1: CalDAV Protocol
- Open standard for calendar access
- Pros: Standard protocol, works with multiple calendar providers
- Cons: Complex protocol, requires server setup, limited Google Calendar support
- Why not chosen: Google Calendar's CalDAV support is limited and being deprecated

### Alternative 2: Google Calendar API with API Keys
- Simple API key authentication
- Pros: Simpler implementation, no OAuth complexity
- Cons: Cannot access private calendar data, security limitations
- Why not chosen: Cannot access user's personal calendar data

### Alternative 3: Manual Import/Export
- Users manually import/export calendar data
- Pros: No authentication complexity, works offline
- Cons: Poor user experience, no real-time sync, manual process
- Why not chosen: Doesn't meet user experience requirements for seamless integration

### Alternative 4: Third-Party Calendar Services
- Use calendar service like Calendly or similar
- Pros: Simpler integration, specialized calendar features
- Cons: Additional service dependency, cost, limited Google Calendar integration
- Why not chosen: Users want to use their existing Google Calendar

## Implementation Notes

### OAuth 2.0 Configuration
```go
oauth2Config := &oauth2.Config{
    ClientID:     credentials.ClientID,
    ClientSecret: credentials.ClientSecret,
    Endpoint:     google.Endpoint,
    RedirectURL:  "http://localhost:4000/auth/callback",
    Scopes:       []string{calendar.CalendarScope},
}
```

### Token Storage with Encryption
```go
type TokenStorage struct {
    db  *bolt.DB
    key []byte // Encryption key
}

func (ts *TokenStorage) StoreToken(token *oauth2.Token) error {
    encryptedToken, err := ts.encryptToken(token)
    if err != nil {
        return err
    }
    
    return ts.db.Update(func(tx *bolt.Tx) error {
        bucket, err := tx.CreateBucketIfNotExists([]byte("oauth"))
        if err != nil {
            return err
        }
        return bucket.Put([]byte("google"), encryptedToken)
    })
}
```

### Calendar Event Mapping
```go
// Map Task to Google Calendar Event
func (cs *CalendarService) taskToEvent(task *model.Task) *calendar.Event {
    event := &calendar.Event{
        Summary:     task.Title.String,
        Description: task.Description.String,
    }
    
    if task.StartTime.Valid && task.Duration.Valid {
        startTime := task.StartTime.Time
        endTime := startTime.Add(time.Duration(task.Duration.Int32) * time.Minute)
        
        event.Start = &calendar.EventDateTime{
            DateTime: startTime.Format(time.RFC3339),
            TimeZone: "America/New_York", // TODO: Use user's timezone
        }
        event.End = &calendar.EventDateTime{
            DateTime: endTime.Format(time.RFC3339),
            TimeZone: "America/New_York",
        }
    }
    
    // Store task ID for bi-directional sync
    event.ExtendedProperties = &calendar.EventExtendedProperties{
        Private: map[string]string{
            "camel_do_task_id": task.ID,
        },
    }
    
    return event
}
```

### Sync Service Implementation
```go
type TaskSyncService struct {
    taskService     *TaskService
    calendarService *CalendarService
    lastSyncTime    time.Time
}

func (tss *TaskSyncService) SyncTasks() error {
    // Sync tasks to calendar
    tasks, err := tss.taskService.GetTasksModifiedSince(tss.lastSyncTime)
    if err != nil {
        return err
    }
    
    for _, task := range tasks {
        if err := tss.syncTaskToCalendar(task); err != nil {
            slog.Error("Failed to sync task to calendar", "task_id", task.ID, "error", err)
            continue
        }
    }
    
    // Sync calendar events to tasks
    events, err := tss.calendarService.GetEventsModifiedSince(tss.lastSyncTime)
    if err != nil {
        return err
    }
    
    for _, event := range events {
        if err := tss.syncEventToTask(event); err != nil {
            slog.Error("Failed to sync event to task", "event_id", event.Id, "error", err)
            continue
        }
    }
    
    tss.lastSyncTime = time.Now()
    return nil
}
```

### Error Handling and Retry Logic
```go
func (cs *CalendarService) createEventWithRetry(event *calendar.Event) (*calendar.Event, error) {
    backoff := &backoff.ExponentialBackOff{
        InitialInterval:     time.Second,
        RandomizationFactor: 0.5,
        Multiplier:          2,
        MaxInterval:         30 * time.Second,
        MaxElapsedTime:      5 * time.Minute,
        Clock:              backoff.SystemClock,
    }
    
    var result *calendar.Event
    err := backoff.Retry(func() error {
        var err error
        result, err = cs.service.Events.Insert("primary", event).Do()
        
        // Handle specific Google API errors
        if err != nil {
            if isRateLimitError(err) {
                return err // Retry
            }
            if isAuthError(err) {
                // Try token refresh
                if refreshErr := cs.refreshToken(); refreshErr == nil {
                    return err // Retry with new token
                }
            }
            return backoff.Permanent(err) // Don't retry
        }
        
        return nil
    }, backoff)
    
    return result, err
}
```

### Security Considerations
- **Token Encryption**: All OAuth tokens encrypted at rest using AES-GCM
- **Secure Storage**: Tokens stored in user's protected configuration directory
- **HTTPS Enforcement**: All OAuth flows require HTTPS in production
- **State Parameter**: Cryptographically secure state parameter for CSRF protection
- **Scope Limitation**: Request minimal required permissions
- **Token Rotation**: Regular token refresh to limit exposure window

### Testing Strategy
```go
func TestOAuthFlow(t *testing.T) {
    // Test OAuth flow with mock Google responses
    mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Mock Google OAuth responses
    }))
    defer mockServer.Close()
    
    // Test authorization URL generation
    authURL := oauthService.GetAuthURL()
    assert.Contains(t, authURL, "oauth2/auth")
    assert.Contains(t, authURL, "calendar")
    
    // Test token exchange
    token, err := oauthService.ExchangeCode("mock_code", "mock_state")
    assert.NoError(t, err)
    assert.NotEmpty(t, token.AccessToken)
}
```

## References
- [Google OAuth 2.0 Documentation](https://developers.google.com/identity/protocols/oauth2)
- [Google Calendar API Documentation](https://developers.google.com/calendar/api)
- [OAuth 2.0 Security Best Practices](https://datatracker.ietf.org/doc/html/draft-ietf-oauth-security-topics)
- [PKCE RFC 7636](https://datatracker.ietf.org/doc/html/rfc7636)

---
*Date: 2025-08-27*
*Authors: Claude Code*
*Reviewers: N/A*