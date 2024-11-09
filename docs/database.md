## Troubleshooting

### Auth

Commands for linux environments:

````shell
export PGPASSWORD='secure_password'
psql -h localhost -U erm_user -d erm_database
````

alternatively, if you are on windows:

````powershell
$env:PGPASSWORD = 'secure_password'
psql -h localhost -U erm_user -d erm_database
````

### PostgreSQL commands

To show tables in `psql`, the interactive terminal for PostgreSQL, you can use the `\dt` command. Here’s how:

1. Start `psql` and connect to your database:
   ```bash
   psql -U your_username -d your_database
   ```

2. Once connected, list all tables by running:
   ```sql
   \dt
   ```

### Explanation:
- `\dt` lists all tables in the current schema.
- If you want to list tables in a specific schema, you can use:
  ```sql
  \dt schema_name.*
  ```

### Additional Commands:
- To show all database objects, including tables, views, and indexes, use:
  ```sql
  \d
  ```

- For more detailed information about a specific table, use:
  ```sql
  \d table_name
  ```

This should give you a clear overview of the tables present in your PostgreSQL database.


# DB documentation

To incorporate PostgreSQL into your project, here’s how you can structure it:

**Add PostgreSQL to Docker Compose**:
Add a PostgreSQL service to your existing `docker-compose.yml` to spin up the database alongside your app:

```yaml
services:
  erm-app:
    image: edgeforge/erm
    container_name: erm_app
    ports:
      - "8080:8080" # Maps the container's port 8080 to the host's port 8080
    environment:
      - GIN_MODE=release # Add any necessary environment variables here
    networks:
      - erm_network
    restart: unless-stopped

  db:
    image: postgres:latest
    container_name: erm_db
    environment:
      POSTGRES_USER: erm_user
      POSTGRES_PASSWORD: secure_password
      POSTGRES_DB: erm_database
    ports:
      - "5432:5432" # Maps the container's port 5432 to the host's port 5432
    networks:
      - erm_network
    volumes:
      - db_data:/var/lib/postgresql/data
    restart: unless-stopped

networks:
  erm_network:
    driver: bridge

volumes:
  db_data:
```

**Connect Your Go Application to PostgreSQL**:
Create a separate package or set of functions to handle database interactions. Using a Go library like `pgx` or `gorm` is common for managing database connections.

**Example structure**:
- `db/connection.go`: A package that initializes and provides the database connection.
- `handlers/db_handler.go`: Functions that handle database operations.

**Sample code for `db/connection.go`**:

```go
package db

import (
    "database/sql"
    _ "github.com/lib/pq" // PostgreSQL driver
    "github.com/rs/zerolog/log"
)

func InitDB() (*sql.DB, error) {
    connStr := "postgresql://erm_user:secure_password@db:5432/erm_database?sslmode=disable" // TODO: remove hardcoding from var to env var.
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal().Msgf("Failed to connect to the database:", err)
        return nil, err
    }

    // Ping to ensure connection is active
    if err := db.Ping(); err != nil {
        log.Fatal().Msgf("Failed to ping the database:", err)
        return nil, err
    }

    log.Info().Msg("Database connection established")
    return db, nil
}
```

**Create Handlers for Database Operations**:
You can create specific handlers for database operations or integrate database logic into existing handlers, depending on your project’s complexity.

**Example of a simple handler**:

```go
package handlers

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/MKTHEPLUGG/ERM/db"
)

func GetUsers(c *gin.Context) {
    dbConn, err := db.InitDB()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection failed"})
        return
    }
    defer dbConn.Close()

    rows, err := dbConn.Query("SELECT id, name FROM users")
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query users"})
        return
    }
    defer rows.Close()

    var users []map[string]interface{}
    for rows.Next() {
        var id int
        var name string
        if err := rows.Scan(&id, &name); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan user"})
            return
        }
        users = append(users, map[string]interface{}{"id": id, "name": name})
    }

    c.JSON(http.StatusOK, users)
}
```


### How to Set Up the `migrations` Directory:
1. **Create the Directory**:
   - In your project directory (where `main.go` is located), create a folder named `migrations`:
     ```bash
     mkdir migrations
     ```

2. **Add Migration Files**:
   - Inside the `migrations` folder, create your migration files. These files should be named with a version number and description to ensure they are applied in order.
   - Example:
     ```bash
     migrations/
     ├── 001_create_users_table.up.sql
     ├── 001_create_users_table.down.sql
     ├── 002_add_roles_table.up.sql
     └── 002_add_roles_table.down.sql
     ```

   - The `.up.sql` file contains the SQL statements for applying the migration (e.g., creating tables).
   - The `.down.sql` file contains the SQL statements for rolling back the migration (e.g., dropping tables).


### Example Migration File (`001_create_users_table.up.sql`):
```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Running the Migrations:
1. **Embed the files into binary**
   - as to not have to include a migration directory into every container it seemed logical to encorporate that logic into the binary. We use IOFS provider for that.
   ````bash
   func runMigrations(db *sql.DB) {
      log.Info().Msg("Starting migrations")

      // Create a new iofs driver for the embedded migrations
      sourceDriver, err := iofs.New(migrationFiles, "migrations")
      if err != nil {
          log.Fatal().Msgf("Failed to create iofs source driver: %v", err)
      }

      // Create a database driver
      driver, err := postgres.WithInstance(db, &postgres.Config{})
      if err != nil {
          log.Fatal().Msgf("Failed to create migration driver: %v", err)
      }

      // Create the migration instance using the iofs source and database driver
      m, err := migrate.NewWithInstance("iofs", sourceDriver, "postgres", driver)
      if err != nil {
          log.Fatal().Msgf("Failed to create migration instance: %v", err)
      }

      // Run the migrations
      if err := m.Up(); err != nil && err != migrate.ErrNoChange {
          log.Fatal().Msgf("Migration failed: %v", err)
      } else {
          log.Info().Msg("Migrations applied successfully")
      }
    }
   ````

### Final Tips:
- **Order of Migrations**: Ensure the migration files are numbered sequentially (e.g., `001`, `002`) to apply them in the correct order.
- **Rollback**: The `.down.sql` file should contain the reverse logic of `.up.sql` to undo changes if needed.
