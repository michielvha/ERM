package handlers

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "erm/db"
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