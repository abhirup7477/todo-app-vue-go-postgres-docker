package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var db *sql.DB

type Task struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed"`
}

type TaskInput struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
}

func root(c *gin.Context) {
	c.JSON(200, "Home Page!")
}

func getAllTodos(c *gin.Context) {
	query := `select id, title, description, completed from tasks`
	rows, err := db.Query(query)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "Database query failed",
		})
		return
	}

	var todos []Task
	var todo Task
	for rows.Next() {
		err := rows.Scan(
			&todo.ID,
			&todo.Title,
			&todo.Description,
			&todo.Completed,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Row scan failed",
			})
			return
		}
		todos = append(todos, todo)
	}

	c.JSON(http.StatusOK, todos)
}

func addTask(c *gin.Context) {
	var input TaskInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	query := `
		insert into tasks (title, description)
		values ($1, $2)
		returning id, completed
	`
	var (
		id        uuid.UUID
		completed bool
	)
	err := db.QueryRow(query, input.Title, input.Description).
		Scan(&id, &completed)

	if err != nil {
		// Unique constraint error
		if strings.Contains(err.Error(), "unique") {
			c.JSON(http.StatusConflict, gin.H{
				"error": "Task with this title already exists",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to insert task",
		})
		return
	}

	c.JSON(http.StatusCreated, Task{
		ID:          id,
		Title:       input.Title,
		Description: input.Description,
		Completed:   completed,
	})
}

func getById(c *gin.Context) {
	id := c.Param("id")
	var todo Task

	query := `
		select id, title, description, completed
		from tasks
		where id = $1
	`
	err := db.QueryRow(query, id).Scan(
		&todo.ID,
		&todo.Title,
		&todo.Description,
		&todo.Completed,
	)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Task not found",
		})
		return
	}
	c.JSON(http.StatusOK, todo)
}

func removeById(c *gin.Context) {
	id := c.Param("id")

	query := `
		DELETE FROM tasks
		WHERE id = $1
	`

	result, err := db.Exec(query, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete task",
		})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to verify delete result",
		})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Task not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Task deleted successfully",
		"id":      id,
	})
}

func update(c *gin.Context) {
	var input Task

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
	}

	query := `
		update tasks set title=$1, description=$2 where id = $3;
	`

	result, err := db.Exec(query, input.Title, input.Description, input.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update task",
		})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to verify update result",
		})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Task not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Task updated successfully",
		"task":    input,
	})
}

func toggleCompleted(c *gin.Context) {
	idStr := c.Param("id")

	// Parse the string into uuid.UUID
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid task ID",
		})
		return
	}

	// Fetch current completed status
	var completed bool
	querySelect := `SELECT completed FROM tasks WHERE id = $1`
	err = db.QueryRow(querySelect, id).Scan(&completed)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Task not found",
		})
		return
	}

	// Toggle status
	newStatus := !completed
	queryUpdate := `UPDATE tasks SET completed = $1 WHERE id = $2`
	_, err = db.Exec(queryUpdate, newStatus, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to toggle task status",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Task status toggled successfully",
		"id":        id,
		"completed": newStatus,
	})
}

func main() {
	var err error
	dbURL := "postgres://" +
		os.Getenv("POSTGRES_USER") + ":" +
		os.Getenv("POSTGRES_PASSWORD") + "@" +
		os.Getenv("POSTGRES_HOST") + ":" +
		os.Getenv("POSTGRES_PORT") + "/" +
		os.Getenv("POSTGRES_DB") +
		"?sslmode=disable"

	db, err = sql.Open("pgx", dbURL)

	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	router := gin.Default()
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})
	// API routes
	r := router.Group("/todos")
	{
		r.GET("", getAllTodos)
		r.POST("/add", addTask)
		r.GET("/:id", getById)
		r.DELETE("/:id", removeById)
		r.PUT("", update)
		r.PATCH("/toggle/:id", toggleCompleted)
	}

	// Serve frontend safely
	router.Static("/frontend", "./frontend/dist")

	router.NoRoute(func(c *gin.Context) {
		c.File("./frontend/dist/index.html")
	})

	port := os.Getenv("BACKEND_PORT")
	if port == "" {
		port = "8080" // fallback for safety
	}

	router.Run(":" + port)
}
