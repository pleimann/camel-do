# ADR-0002: BoltDB as Embedded Database

## Status
Accepted

## Context
Camel-Do requires a database solution for storing tasks, projects, and application state. As a single-user desktop application, we need to balance simplicity, performance, and deployment convenience. The application needs to be distributable as a single binary without external database dependencies.

Key requirements:
- Single-user operation (no concurrent access from multiple users)
- Embedded database (no external server required)
- ACID transactions for data consistency
- Good performance for typical task management operations
- Simple deployment and backup strategy
- Cross-platform compatibility

## Decision
We will use BoltDB as the embedded database solution for Camel-Do.

BoltDB is a pure Go key-value database that provides:
- ACID compliance with read-write transactions
- Embedded, serverless architecture
- Single file storage with no external dependencies
- Good performance for read-heavy workloads
- Simple API with bucket-based organization
- Cross-platform support

Data organization:
```
camel-do.db
├── tasks/          -> Task entities (key: task_id, value: GOB-encoded Task)
├── projects/       -> Project entities (key: project_id, value: GOB-encoded Project)
├── oauth/          -> OAuth tokens (key: service_name, value: encrypted tokens)
└── settings/       -> Application settings (key: setting_name, value: setting_value)
```

## Consequences

### Positive Consequences
- **Zero Configuration**: No database server setup or configuration required
- **Single Binary Deployment**: Database embedded in application, simplifying distribution
- **ACID Transactions**: Ensures data consistency for critical operations
- **Good Performance**: Fast read operations suitable for task management workloads
- **Simple Backup**: Database backup is a simple file copy operation
- **Cross-Platform**: Works consistently across all supported operating systems
- **No Network Dependency**: Eliminates network-related database connection issues

### Negative Consequences
- **Single User Limitation**: Cannot support concurrent access from multiple users
- **Limited Query Capabilities**: Key-value store requires application-level indexing
- **Write Scalability**: Write operations require exclusive locks
- **Memory Usage**: Entire database mapped to memory for performance

### Risks
- **Database Corruption** - Mitigation: Regular backups, graceful shutdown handling
- **Storage Limitations** - Mitigation: Monitor database file size, implement cleanup procedures
- **Migration Complexity** - Mitigation: Version-aware data migration procedures

## Alternatives Considered

### Alternative 1: SQLite
- Embedded SQL database with rich query capabilities
- Pros: SQL interface, better query support, wider ecosystem
- Cons: More complex than needed, SQL injection concerns, larger binary size
- Why not chosen: Overkill for simple key-value operations, adds complexity

### Alternative 2: PostgreSQL/MySQL
- Full-featured relational databases
- Pros: Rich feature set, excellent tooling, multi-user support
- Cons: Requires external server, complex setup, overkill for single-user app
- Why not chosen: Violates simplicity requirements, external dependency

### Alternative 3: In-Memory + JSON Files
- Simple file-based persistence with JSON serialization
- Pros: Very simple, human-readable format, easy debugging
- Cons: No transactions, poor performance, risk of data corruption
- Why not chosen: Lack of ACID properties, poor reliability for critical data

### Alternative 4: BadgerDB
- Modern embedded key-value database
- Pros: Better performance than BoltDB, LSM-tree design
- Cons: More complex, less mature ecosystem, larger memory footprint
- Why not chosen: BoltDB sufficient for current needs, prefer stability over cutting-edge performance

## Implementation Notes

### Data Serialization
- Use Go's `encoding/gob` for binary serialization of structs
- Provides good performance and handles complex Go types
- Version-aware serialization for future migrations

### Transaction Management
```go
func (ts *TaskService) AddTask(task *model.Task) error {
    return ts.db.Update(func(tx *bolt.Tx) error {
        bucket, err := tx.CreateBucketIfNotExists([]byte("tasks"))
        if err != nil {
            return err
        }
        
        taskBytes, err := task.Marshal()
        if err != nil {
            return err
        }
        
        return bucket.Put([]byte(task.ID), taskBytes)
    })
}
```

### Database Location
- Development: Local project directory
- Production: User's configuration directory (`~/.config/camel-do/`)
- File permissions: 0600 (user read/write only)

### Migration Strategy
- Store schema version in settings bucket
- Implement version-aware migration functions
- Backup before migration, rollback capability

## References
- [BoltDB Documentation](https://pkg.go.dev/go.etcd.io/bbolt)
- [BoltDB Design Paper](https://github.com/boltdb/bolt)
- [Embedded Database Comparison](https://github.com/dgraph-io/badger#comparison-with-other-databases)

---
*Date: 2025-08-27*
*Authors: Claude Code*
*Reviewers: N/A*