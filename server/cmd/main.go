package main

import (
	"log"
	"os"

	"dune-takehome-server/database"
	"dune-takehome-server/handlers"
	"dune-takehome-server/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Connect to MongoDB
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017/dune-form-builder"
	}

	if err := database.Connect(mongoURI); err != nil {
		log.Fatalf("❌ Failed to connect to MongoDB: %v", err)
	}

	defer func() {
		if err := database.Disconnect(); err != nil {
			log.Printf("Error disconnecting from MongoDB: %v", err)
		}
	}()

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	// Health check route
	app.Get("/health", func(c *fiber.Ctx) error {
		// Check MongoDB connection
		mongoStatus := "connected"
		if err := database.Ping(); err != nil {
			mongoStatus = "disconnected"
		}

		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "Dune Security Form Builder API",
			"mongodb": mongoStatus,
			"version": "1.0.0",
		})
	})

	// API routes
	api := app.Group("/api/v1")

	// TODO: Add route handlers
	setupRoutes(api)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Server starting on port %s", port)
	log.Fatal(app.Listen(":" + port))
}

func setupRoutes(api fiber.Router) {
	// Initialize handlers
	userHandler := handlers.NewUserHandler()
	formHandler := handlers.NewFormHandler()

	// Auth routes
	auth := api.Group("/auth")
	auth.Post("/register", userHandler.Register)
	auth.Post("/login", userHandler.Login)
	auth.Get("/profile", middleware.AuthRequired(), userHandler.GetProfile)

	// Form routes
	forms := api.Group("/forms", middleware.AuthRequired())
	forms.Get("/", formHandler.GetUserForms)
	forms.Post("/", formHandler.CreateForm)
	forms.Get("/:id", formHandler.GetFormByID)
	forms.Put("/:id", formHandler.UpdateForm)

	// Public routes (no auth required)
	public := api.Group("/public")
	public.Get("/forms/:shareUrl", formHandler.GetPublicForm)
	public.Post("/forms/:shareUrl/responses", formHandler.SubmitPublicFormResponse)

	forms.Post("/:id/responses", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Submit form response"})
	})
}
