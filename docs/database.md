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
    "log"
)

func InitDB() (*sql.DB, error) {
    connStr := "postgresql://erm_user:secure_password@db:5432/erm_database?sslmode=disable"
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal("Failed to connect to the database:", err)
        return nil, err
    }

    // Ping to ensure connection is active
    if err := db.Ping(); err != nil {
        log.Fatal("Failed to ping the database:", err)
        return nil, err
    }

    log.Println("Database connection established")
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
    "your_project/db"
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

Yes, `"file://migrations"` refers to a directory named `migrations` located in the same directory as your `main.go` file or within your project's root directory. You can place your SQL migration scripts in this `migrations` directory, and the `golang-migrate` library will read and apply them.

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
1. **Ensure the `migrations` directory path is correct**:
   - Use the path `"file://migrations"` if the `migrations` directory is in the same folder as `main.go`.
   - Adjust the path accordingly if your `migrations` folder is nested or in a different directory (e.g., `"file://db/migrations"`).

2. **Run Migrations in Code**:
   ```go
   func main() {
       dbConn, err := db.InitDB()
       if err != nil {
           log.Fatal().Msg("Could not establish a database connection")
       }
       defer dbConn.Close()

       // Run database migrations
       runMigrations(dbConn)

       // Start the web server
       r := gin.Default()
       r.POST("/login", handlers.Login)
       protected := r.Group("/protected")
       protected.Use(middleware.JWTAuthMiddleware())
       {
           protected.GET("/", handlers.ProtectedEndpoint)
       }
       log.Info().Msg("Server is starting on port 8080")
       r.Run("0.0.0.0:8080")
   }
   ```

### Final Tips:
- **Order of Migrations**: Ensure the migration files are numbered sequentially (e.g., `001`, `002`) to apply them in the correct order.
- **Rollback**: The `.down.sql` file should contain the reverse logic of `.up.sql` to undo changes if needed.

With this setup, your migrations will be organized and ready to be run automatically when your application starts.
