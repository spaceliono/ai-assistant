## Package `database`

MySQL database initialization and queries for assistant records and reference data.

### Public Functions

```go
func InitDB(cfg *config.Config) error
```
Initializes the global `*sql.DB` connection from `config.Config` and starts a background ping.

```go
func GetDB() *sql.DB
```
Returns the initialized global DB.

```go
func CloseDB() error
```
Closes the global DB connection.

```go
func InsertAIAssistantRecord(record *models.AIAssistantRecords) error
```
Persists an assistant-related record for a conversation.

```go
func GetAIAssistantRecordsBySource(sourcePlatform, sourceID string) ([]models.AIAssistantRecords, error)
```
Returns all records for a given `(source_platform, source_id)` pair ordered by `id`.

```go
func GetJobs() ([]string, error)
```
Returns job names from `eb_user_tag`.

### Example: Initialize and Query

```go
cfg, _ := config.LoadConfig("config.yml")
if err := database.InitDB(cfg); err != nil { panic(err) }
defer database.CloseDB()

records, err := database.GetAIAssistantRecordsBySource("web", "user-123")
if err != nil { panic(err) }
for _, r := range records {
    fmt.Println(r.Role, r.Content)
}
```

### Suggested Table Schema (reference)

This example schema matches the fields used by the code. Adjust types/charset as needed.

```sql
CREATE TABLE IF NOT EXISTS AI_Assistant_Records (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  source_platform VARCHAR(64) NOT NULL,
  source_id VARCHAR(128) NOT NULL,
  ai_platform VARCHAR(64) NOT NULL,
  role VARCHAR(32) NOT NULL,
  content TEXT NOT NULL,
  create_time DATETIME NOT NULL
);
```

