package main

import (
	"Todo-App/internal/config"
	"Todo-App/internal/database"
	"Todo-App/internal/handlers"
	"Todo-App/internal/middleware"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {

	var cfg *config.Config
	var err error
	cfg, err = config.Load()

	if err != nil {
		fmt.Println("Error loading configuration:", err)
		return
	}

	var pool *pgxpool.Pool
	pool, err = database.Connect(cfg.DatabaseURL)

	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return
	}
	defer pool.Close()

	var router *gin.Engine = gin.Default()
	router.SetTrustedProxies(nil)
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Server is running",
			"status":  "success",
		})
	})
	router.GET("/api/live/ws", func(c *gin.Context) {
		c.Status(200)
	})

	// router.GET("/todo", handlers.GetAllTodosHandler(pool))
	// router.POST("/todos", handlers.CreateTodoHandler(pool))
	// router.GET("/todo/:id", handlers.GetToDoByIDHandler(pool))
	// router.PUT("/todo/:id", handlers.UpdateTodo(pool))
	// router.DELETE("/todo/:id", handlers.DeleteTodoHandler(pool))

	router.POST("/auth/register", handlers.CreateUserHandler(pool))
	router.POST("/auth/login", handlers.LoginHandler(pool, cfg))

	protected := router.Group("/todos")
	protected.Use(middleware.AuthMiddleware(cfg))
	{
		protected.POST("", handlers.CreateTodoHandler(pool))
		protected.GET("", handlers.GetAllTodosHandler(pool))
		protected.GET("/:id", handlers.GetToDoByIDHandler(pool))
		protected.PUT("/:id", handlers.UpdateToDoHandler(pool))
		protected.DELETE("/:id", handlers.DeleteToDoHandler(pool))
	}

	// Middleware Test Route
	router.GET("/protected-test", middleware.AuthMiddleware(cfg), handlers.TestProtectedHandler())

	error := router.Run(":" + cfg.Port)
	if error != nil {
		fmt.Println("Error starting server:", error)
	}
}
