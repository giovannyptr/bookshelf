package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/giovannyptr/bookshelf/internal/auth"
	"github.com/giovannyptr/bookshelf/internal/books"
	"github.com/giovannyptr/bookshelf/internal/users"
	"github.com/giovannyptr/bookshelf/models"
)

/* ---------- tiny helper ---------- */
func env(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

/* -------------------- main -------------------- */
func main() {
	_ = godotenv.Load(".env", "cmd/server/.env")

	wd, _ := os.Getwd()
	log.Printf("DEBUG CWD=%s", wd)

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		// Fallback to DB_* pieces if DATABASE_URL not provided
		host := env("DB_HOST", "127.0.0.1")
		port := env("DB_PORT", "5433")
		user := env("DB_USER", "bookshelf")
		pass := env("DB_PASSWORD", "bookshelf")
		name := env("DB_NAME", "bookshelf")
		ssl := env("DB_SSLMODE", "disable")
		dsn = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", user, pass, host, port, name, ssl)
	}
	log.Printf("DEBUG DATABASE_URL=%q", dsn)

	// Connect DB
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	log.Println("Connected to PostgreSQL successfully")

	// Auto-migrate
	if err := db.AutoMigrate(&models.Book{}, &models.User{}); err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}
	log.Println(" Auto-migration completed")

	// Repos
	bookRepo := books.NewRepository(db)
	userRepo := users.NewRepository(db)

	// Seed admin user if missing
	ensureAdmin(userRepo)

	// Handlers
	bookHandler := books.NewHandler(bookRepo, "./uploads")
	authHandler := auth.NewHandler(userRepo)

	// Ensure uploads dir exists
	_ = os.MkdirAll("uploads", 0o755)

	// Gin engine
	r := gin.Default()

	// Silence proxy warning (no reverse proxy in dev)
	r.SetTrustedProxies(nil)

	// CORS so Vue (Vite) can call us
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:5173",
			"http://127.0.0.1:5173",
			"http://localhost:4173",
			"http://127.0.0.1:4173",
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Static for uploaded covers
	r.Static("/uploads", "./uploads")

	// Health
	r.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"ok": true}) })

	// Auth
	authHandler.RegisterRoutes(r) // /auth/register, /auth/login, /auth/me

	// Public book reads
	r.GET("/books", bookHandler.List)
	r.GET("/books/:id", bookHandler.Detail)

	// Protected writes
	secured := r.Group("/")
	secured.Use(auth.AuthRequired())
	secured.POST("/books", bookHandler.Create)
	secured.PUT("/books/:id", bookHandler.Update)
	secured.DELETE("/books/:id", bookHandler.Delete)

	// Run server
	port := env("APP_PORT", "8080")
	log.Printf("Server running on http://localhost:%s", port)
	log.Fatal(r.Run(":" + port))
}

/* ---------- admin seeding ---------- */
func ensureAdmin(userRepo *users.Repository) {
	email := env("ADMIN_EMAIL", "admin@mail.com")
	pass := env("ADMIN_PASSWORD", "adminbookshelf")
	name := env("ADMIN_NAME", "Admin")

	if u, _ := userRepo.ByEmail(email); u != nil {
		log.Printf(" admin already exists: %s", email)
		return
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	admin := models.User{
		Email:    email,
		Password: string(hash),
		Name:     name,
		Role:     "admin",
	}
	if err := userRepo.Create(&admin); err != nil {
		log.Printf(" failed to create admin: %v", err)
		return
	}
	log.Printf(" Admin user created: %s (password from ADMIN_PASSWORD)", email)
}
