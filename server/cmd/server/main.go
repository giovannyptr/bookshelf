package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/giovannyptr/bookshelf/internal/books"
	"github.com/giovannyptr/bookshelf/internal/platform"
)

func env(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func main() {
	_ = godotenv.Load()

	db, err := platform.OpenGorm()
	if err != nil {
		log.Fatalf("DB connect failed: %v", err)
	}
	log.Println("✅ DB connected")

	// repositories & migrations
	bookRepo := books.NewRepository(db)
	if err := bookRepo.Migrate(); err != nil {
		log.Fatalf("migrate failed: %v", err)
	}

	// ensure uploads dir exists
	uploadDir := "uploads"
	_ = os.MkdirAll(uploadDir, 0o755)

	// gin
	r := gin.Default()
	r.Use(platform.CORS())

	// health
	r.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"ok": true}) })

	// books module
	books.NewHandler(bookRepo, uploadDir).RegisterRoutes(r)

	port := env("APP_PORT", "8080")
	log.Printf("🚀 Server on http://localhost:%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
