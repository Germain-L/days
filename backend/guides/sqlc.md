Below is a concise, up‑to‑date research‑style overview of **sqlc** for Go developers, incorporating the latest information (through April 2025) from the web results you asked about.

---

## 1. What is sqlc?

**sqlc** is a command‑line tool that **generates type‑safe Go code directly from raw SQL**.  
It parses your SQL files, validates them against the target database engine, and emits Go structs and query methods that compile‑time guarantee correct types and placeholders.

Key benefits:

| Benefit | Explanation |
|---------|-------------|
| **Zero‑runtime SQL errors** | Queries are parsed by the database’s own parser during generation, so syntax errors are caught before you run the program. |
| **Full SQL power** | Write plain SQL (CTEs, window functions, JSON operators, etc.) – no ORM abstraction limits. |
| **Fast, lightweight** | Generated code uses the native driver (`pgx/v5` for PostgreSQL, `mysql` driver, etc.) and has no reflection overhead. |
| **Schema‑driven models** | Table definitions become Go structs automatically, with optional overrides for custom types. |
| **Migrations‑aware** | Works well with tools like `golang‑migrate`; you can keep schema and generated code in sync. |

---

## 2. Recent releases (as of 2025)

| Version | Release date | Notable changes |
|---------|--------------|-----------------|
| **1.29.0** | 2025‑04‑14 | • Bug fixes for PostgreSQL `WHERE … UNION …` and MySQL `VECTOR` column.<br>• New `sqlc-gen-from-template` plugin.<br>• Option to wrap query errors with the query name (`wrap_query_error`).<br>• Upgrade to Go 1.24.1 and many dependency bumps. |
| **1.28.0** | 2025‑01‑20 | • Added support for TiDB parser, MySQL `VECTOR` type.<br>• Boolean‑based dynamic filters (`@param::boolean`).<br>• `sqlc.verify` for schema‑change safety.<br>• `sqlc.push` (renamed from `upload`) for CI/CD integration. |
| **1.27.0** | 2024‑08‑05 | • Local managed‑database support (no cloud dependency).<br>• Initialisms configuration for Go naming.<br>• Added C# language support and many ORM‑style improvements. |
| **1.26.0** | 2024‑03‑28 | • Security fix for output plugins.<br>• Breaking change: `SMALLINT` → `int16`, `TINYINT` → `int8` for MySQL. |

*All changelogs are available in the official docs: <https://docs.sqlc.dev/en/latest/reference/changelog.html>.*

---

## 3. Core concepts & workflow

### 3.1 Project layout (recommended)

```
myapp/
├── db/
│   ├── migrations/          # .sql migration files (golang‑migrate)
│   ├── schema.sql           # CREATE TABLE statements (used by sqlc)
│   ├── queries/
│   │   ├── users.sql        # sqlc‑annotated queries
│   │   └── orders.sql
│   └── sqlc.yaml            # configuration (see §3.2)
├── internal/
│   └── repo/                # wrapper around generated code
├── cmd/
│   └── myapp/               # main entry point / CLI
└── go.mod
```

### 3.2 `sqlc.yaml` (v2 format)

```yaml
version: "2"
sql:
  - engine: "postgresql"
    schema: "db/schema.sql"
    queries: "db/queries"
    gen:
      go:
        package: "db"
        out: "db/sqlc"
        sql_package: "pgx/v5"
        emit_json_tags: true
        emit_interface: true          # optional, for test‑ability
        emit_exact_table_names: false
        emit_pointers_for_null_types: true
overrides:
  - db_type: "timestamptz"
    go_type:
      import: "time"
      type: "Time"
```

*Key options*:

| Option | What it does |
|--------|--------------|
| `sql_package` | Choose driver (`pgx/v5`, `database/sql`, `mysql`, etc.). |
| `emit_interface` | Generates an interface (`Querier`) for easier mocking. |
| `emit_pointers_for_null_types` | Nullable columns become `*T` or `pgtype.*`. |
| `overrides` | Map DB types to custom Go types (e.g., `decimal.Decimal`). |

### 3.3 Writing queries

Each query must start with a comment:

```sql
-- name: GetUser :one
SELECT id, name, bio FROM users WHERE id = $1;
```

Supported tags:

| Tag | Return type |
|-----|-------------|
| `:one` | `(Model, error)` – exactly one row |
| `:many` | `([]Model, error)` – zero‑or‑more rows |
| `:exec` | `error` – no rows returned |
| `:execrows` | `(int64, error)` – rows‑affected count |
| `:copyfrom` | `(int64, error)` – bulk COPY import (PostgreSQL) |

### 3.4 Dynamic filters (new in 1.28)

Use boolean parameters to enable/disable predicates without string concatenation:

```sql
-- name: ListOrders :many
SELECT * FROM orders
WHERE
  (@by_customer::boolean AND customer_id = @customer_id) OR
  (@by_status::boolean   AND status = @status)
ORDER BY created_at DESC;
```

In Go you set the booleans:

```go
params := ListOrdersParams{
    ByCustomer: true,
    CustomerID: 42,
    ByStatus:   false, // ignored
}
orders, err := q.ListOrders(ctx, params)
```

### 3.5 Schema‑change verification (`sqlc verify`)

`sqlc verify` checks that a new schema version does **not break existing queries**. It runs the queries against the new schema in a sandbox and reports errors such as ambiguous column references. This is useful in CI pipelines.

```bash
sqlc verify --against v1.27.0
```

---

## 4. Integration with other tools

| Tool | How it fits |
|------|-------------|
| **golang‑migrate** | Apply `*.up.sql` migrations; `sqlc` reads the final schema (`schema.sql`) to generate models. |
| **pgx/v5** | Recommended driver; `sqlc` emits `*pgxpool.Pool`‑compatible code. |
| **tablewriter** | Handy for CLI output of generated structs (e.g., `PrintAuthors`). |
| **godotenv** | Load DB credentials from `.env` before calling `sqlc.New(db)`. |
| **CI/CD** | Use `sqlc push` (formerly `upload`) to push schema & queries to sqlc Cloud for verification. |

---

## 5. Common pitfalls & solutions

| Issue | Cause | Fix |
|-------|-------|-----|
| **Generated `DBTX` interface missing `*Context` methods** (v1.27) | Using `emit_interface: false` with older `sqlc-gen-go` version. | Upgrade to `sqlc-dev/sqlc` ≥ 1.28 or set `emit_interface: true`. |
| **Nullable columns become `interface{}`** | No `emit_pointers_for_null_types` flag. | Set `emit_pointers_for_null_types: true` or use `pgtype` overrides. |
| **`SELECT *` rewrites to explicit columns** | sqlc rewrites for safety; you may see a different query in generated code. | Accept it – it guarantees compile‑time column‑type matching. |
| **Migrations out‑of‑sync with generated code** | Forgetting to run `sqlc generate` after a migration. | Add a `make generate` step after `migrate up` in your workflow. |
| **`sqlc verify` fails on CTEs** | Older version (pre‑1.28) didn’t support some CTE syntax. | Upgrade to ≥ 1.28 or simplify the CTE for verification. |

---

## 6. Example: Minimal Go program using sqlc

```go
package main

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"myapp/db/sqlc" // generated package
	"myapp/helpers" // env & CRUD wrappers
)

func main() {
	ctx := context.Background()
	connStr := helpers.GetDBConnectionString()
	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		log.Fatalf("DB pool: %v", err)
	}
	defer pool.Close()

	q := sqlc.New(pool) // generated Queries struct

	// Create a user
	author, err := helpers.CreateAuthor(ctx, q, "Ada Lovelace", "First programmer")
	if err != nil {
		log.Fatalf("create: %v", err)
	}
	helpers.PrintAuthor(author)

	// List all authors
	authors, _ := helpers.ReadAuthors(ctx, q)
	helpers.PrintAuthors(authors)

	// Update name
	_, _ = helpers.UpdateAuthorName(ctx, q, author.ID, "Ada Byron")
}
```

*The helper functions (`CreateAuthor`, `PrintAuthor`, etc.) are thin wrappers around the generated methods, keeping the main program tidy.*

---

## 7. Where to learn more

| Resource | Link |
|----------|------|
| Official docs & tutorial | <https://docs.sqlc.dev/en/latest/> |
| Getting‑started tutorial (PostgreSQL) | <https://docs.sqlc.dev/en/stable/tutorials/getting-started-postgresql.html> |
| `sqlc verify` blog post (2024) | <https://sqlc.dev/blog/verify/> |
| Community examples (GitHub) | <https://github.com/sqlc-dev/sqlc/tree/main/examples> |
| Advanced patterns (dynamic filters) | Brandur’s “sqlc 2024 check‑in” – <https://brandur.org/fragments/sqlc-2024> |

---

### TL;DR

- **sqlc** turns raw SQL into compile‑time‑checked Go code.
- The latest stable release (1.29.0, Apr 2025) adds better error‑wrapping, boolean‑based dynamic filters, and Go 1.24 support.
- Use a `sqlc.yaml` v2 config, write queries with `-- name: … :tag`, run `sqlc generate`, and you get a clean `Queries` struct that works with `pgx/v5`.
- Pair it with migrations (`golang‑migrate`) and `sqlc verify` for safe schema evolution.

Happy coding! 🚀