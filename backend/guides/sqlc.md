# Building Type-Safe Data Layers with sqlc, Go, and PostgreSQL

Before diving into the details, here is the single most important takeaway: **sqlc lets you keep writing plain SQL while guaranteeing at compile-time that your Go code and your PostgreSQL schema stay perfectly in sync**[2][11].  This eliminates an entire class of runtime bugs, removes most boiler-plate, and remains faster than heavyweight ORMs[29][42].

---

## 1. Why sqlc?

### 1.1 The traditional pain points
Writing raw SQL in Go through `database/sql` yields full control and optimum performance, but developers must  
1. marshal every parameter,  
2. scan every column manually, and  
3. discover mismatches only after the program starts[78][55].

### 1.2 sqlc’s compile-time contract
sqlc parses your *schema* (`CREATE TABLE …`) and *queries* (`SELECT …`) then emits **type-safe** Go functions, structs, and interfaces matching the exact column types returned by those queries[2][40].  Any change that breaks a query aborts the build with a compiler error instead of a pager at 03:00 AM[11][81].

### 1.3 Focus of this guide
This document explains how to wire sqlc, Go 1.23+, and PostgreSQL 13+ end-to-end: installation, project layout, migrations, generated code, transactions, prepared statements, CI/CD, and operational best practices—all with *text-only* examples suitable for terminal environments.

---

## 2. Installation and First Project

### 2.1 Prerequisites
* **Go:** `brew install go` or the official binary.  
* **PostgreSQL:** local server (`brew install postgresql`) or Docker image `postgres:15`.  
* **sqlc CLI:** any of  
  ```bash
  brew install sqlc                  # macOS[6]
  sudo snap install sqlc             # Ubuntu[6]
  go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest  # cross-platform[6]
  ```
* **pgx driver:**  
  ```bash
  go get github.com/jackc/pgx/v5@latest[39]
  ```

### 2.2 Directory scaffold
```
myapp/
 ├── db/
 │   ├── migrations/     # *.up.sql & *.down.sql (golang-migrate)
 │   ├── queries/        # hand-written *.sql
 │   └── schema.sql      # reference schema for sqlc vet
 ├── internal/           # application code
 ├── sqlc.yaml           # sqlc config
 └── go.mod
```
This structure keeps DDL, DML, and generated code clearly separated[11][59].

### 2.3 Example schema
```sql
-- db/migrations/001_init.up.sql
CREATE TABLE authors (
  id   BIGSERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  bio  TEXT
);[2]
```

### 2.4 Example queries
```sql
-- db/queries/authors.sql
-- name: GetAuthor :one
SELECT * FROM authors WHERE id = $1 LIMIT 1;

-- name: ListAuthors :many
SELECT * FROM authors ORDER BY name;

-- name: CreateAuthor :one
INSERT INTO authors (name,bio)
VALUES ($1,$2) RETURNING *;[2]
```

### 2.5 sqlc configuration
```yaml
# sqlc.yaml
version: "2"
sql:
  - engine: "postgresql"
    schema: "db/migrations"
    queries: "db/queries"
    gen:
      go:
        package: "db"
        out: "internal/db"
        sql_package: "pgx/v5"   # use native pgx types[39]
        emit_json_tags: true
        emit_prepared_queries: true[50]
        emit_interface: true[21]
```
Key options:  
* `sql_package` swaps the standard driver for pgx’s high-performance API[39].  
* `emit_prepared_queries` generates a `Prepare(ctx,db)` helper returning ready-to-use `*sql.Stmt` handles[50].  

### 2.6 Generate code
```bash
sqlc generate
```
The folder `internal/db` now contains `models.go`, `db.go`, and `authors.sql.go`, each with strongly typed methods such as  
```go
func (q *Queries) GetAuthor(ctx context.Context, id int64) (Author, error)
```
compiling against pgx interfaces rather than `database/sql`[39].

---

## 3. Database Migrations

sqlc **does not** apply migrations; it only parses them[30][33].  Use a dedicated tool—popular choices include **golang-migrate** and **Atlas**.  

### 3.1 Using golang-migrate
```bash
brew install golang-migrate[23]
migrate -database \
  "postgres://user:pass@localhost:5432/mydb?sslmode=disable" \
  -path db/migrations up
```
Name files with zero-padded prefixes (`001_init.up.sql`, `002_add_books.up.sql`) so sqlc parses them in order[30].

### 3.2 Declarative migrations with Atlas
Atlas can produce diff-based migrations and keep `schema.hcl` as single source of truth while sqlc keeps type safety[27].

---

## 4. Application Code Integration

### 4.1 Connection pooling
pgx’s `pgxpool.Pool` offers queue-based pooling out-of-the-box[39][63].  
```go
pool, err := pgxpool.New(ctx,
  "postgres://user:pass@localhost:5432/mydb")
defer pool.Close()

queries := db.New(pool) // generated constructor
```
Under the hood `pgxpool.Pool` satisfies sqlc’s `DBTX` interface, so the same value works for single queries, prepared statements, and transactions[63].

### 4.2 Executing queries
```go
author, err := queries.CreateAuthor(ctx,
  db.CreateAuthorParams{Name: "Ada", Bio: pgtype.Text{String: "Pioneer", Valid: true}})
```
The parameter struct and result struct are generated; any missing or extra fields break compilation[2].

### 4.3 Transactions
```go
func TransferAuthor(ctx context.Context, pool *pgxpool.Pool,
    fromID, toID int64) error {

  tx, err := pool.Begin(ctx)
  if err != nil { return err }
  defer tx.Rollback(ctx)

  qtx := queries.WithTx(tx) // generated helper[41]
  a, err := qtx.GetAuthor(ctx, fromID)
  if err != nil { return err }

  if err := qtx.DeleteAuthor(ctx, fromID); err != nil { return err }
  if _, err := qtx.CreateAuthor(ctx,
        db.CreateAuthorParams{Name: a.Name, Bio: a.Bio}); err != nil { return err }

  return tx.Commit(ctx)
}[41][79]
```
Inside a transaction you **must** call the `WithTx` variant so all statements share the same backend connection[41].

### 4.4 Prepared statements
When `emit_prepared_queries` is true, initialize once at startup:  
```go
prepared, err := db.Prepare(ctx, pool) // O(n) prepare[50]
author, err := prepared.GetAuthor(ctx, 1)
```
Prepared statements lower planning overhead for high-throughput APIs[80].

---

## 5. Configuration Nuances

| Option                          | Meaning | Typical usage | Source |
|---------------------------------|---------|---------------|--------|
| `emit_json_tags`                | Add `json:"field"` to struct fields | Marshal query results into JSON responses | [21] |
| `omit_unused_structs`           | Do not generate models unused by any query | Shrink binary size in micro-services | [81] |
| `rules` + `sqlc vet`            | Custom CEL lint rules for queries | Ban `DELETE` without `WHERE`, enforce index usage | [72][46] |
| `sqlc/db-prepare` vet rule      | Prepares every query during CI against a live DB | Detect stale migrations before merge | [62][72] |
| `diff`                          | Fails CI if generated code is out-of-date | Enforces developers to commit generated files | [24][62] |

---

## 6. Testing and Continuous Delivery

1. **Unit tests** call the generated code against a Dockerised PostgreSQL instance; pgx supports `pgxmock` for pure unit tests.  
2. **CI Pipeline**  
   ```yaml
   steps:
     - sqlc diff       # ensure generated code committed[62]
     - sqlc vet        # lint queries + prepare rule[62]
     - go test ./...   # run application tests
   ```  
3. **Deploy** run migrations (`migrate up`) **before** rolling out new binaries to guarantee schema compatibility[23].

---

## 7. Operational Best Practices

* **Connection pool limits**: tune `SetMaxOpenConns` and `SetMaxIdleConns` (or pgxpool equivalents) to stay below PostgreSQL `max_connections`[55][57].  
* **Isolation levels**: default `ReadCommitted` suffices for most CRUD; use `Serializable` only when truly necessary as it increases contention[79].  
* **Retry only safe errors**: detect transient errors like `deadlock_detected` or network resets and retry with exponential back-off[54][58].  
* **Version control everything**: schema migrations, queries, and generated code are committed—sqlc code is deterministic and reviewable[62].  
* **Avoid hidden SQL**: sqlc cannot parse ad-hoc string concatenation; keep every statement in `.sql` files with `-- name:` annotations[40].  
* **Leverage vet**: enforce house rules (e.g., prohibit `SELECT *`) with custom CEL expressions evaluated during CI[46][72].

---

## 8. Advanced Topics

### 8.1 Multiple Packages
sqlc v2 supports several `sql:` blocks, each generating its own Go package—ideal for splitting big monolith schemas into bounded contexts[2].

### 8.2 Custom Types
Map Postgres enums or domains to Go types via `overrides:` in `sqlc.yaml`; for pgx you can also register codecs for composite types[39].

### 8.3 JSON/JSONB
pgx maps `jsonb` directly into `pgtype.JSONB`; generated structs already embed that type, so no manual marshalling is required[39].

### 8.4 Query annotations
Beyond `:one` and `:many`, use `:execrows` to return affected-row counts, `:copyfrom` for bulk loads, and upcoming `@param` type annotations to drop casts[40][43].

---

## 9. Common Pitfalls and Remedies

| Symptom | Root cause | Fix |
|---------|------------|-----|
| `relation "authors" does not exist` at runtime | Migrations ran after code start or on different DB | Enforce migration order in entrypoint and CI[30] |
| Compile error `cannot use pgtype.Int4` | Using `pgx/v5` without `sql_package: pgx/v5` | Set `sql_package` in config[39] |
| `conn busy` panics when nesting queries | Re-using same connection in parallel goroutines | Use pgxpool and avoid long-running scans[63][84] |
| Generated code reappears in `git status` | Colleague forgot `sqlc generate` | Add `sqlc diff` to CI to block merge[62] |

---

## 10. Conclusion

Combining **sqlc** with **Go** and **PostgreSQL** offers a sweet spot between handwritten SQL’s clarity and a modern type system’s safety guarantees.  By shifting query validation from runtime to compile-time, teams gain earlier feedback, superior IDE assistance, leaner binaries, and simpler debugging.  Integrating pgx unlocks even higher throughput without sacrificing that safety.  Once you layer in structured migrations, rigorous `sqlc vet` rules, and an automated CI pipeline, you have a robust data layer that scales from experimental prototypes to production at planetary scale—while still letting every engineer open a `.sql` file and understand exactly which query is running in production[42][57].
