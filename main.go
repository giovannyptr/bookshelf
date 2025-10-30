package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/giovannyptr/bookshelf/models" // ← adjust if your module path differs
)

// env reads an environment variable with a fallback default.
func env(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func main() {
	// Load environment variables from .env (if present)
	_ = godotenv.Load()

	// ---- Build DSN (URL form) ----
	host := env("DB_HOST", "127.0.0.1")
	port := env("DB_PORT", "5433") // using 5433 on host → 5432 in container
	user := env("DB_USER", "bookshelf")
	pass := env("DB_PASSWORD", "") // empty if using trust auth
	name := env("DB_NAME", "bookshelf")
	ssl := env("DB_SSLMODE", "disable")

	auth := user
	if pass != "" {
		auth = fmt.Sprintf("%s:%s", user, pass)
	}
	dsn := fmt.Sprintf("postgres://%s@%s:%s/%s?sslmode=%s", auth, host, port, name, ssl)

	// ---- Connect to PostgreSQL via GORM ----
	gdb, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// Ping once on startup with a short timeout
	sqlDB, err := gdb.DB()
	if err != nil {
		log.Fatalf("failed to get sql DB: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := sqlDB.PingContext(ctx); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}
	log.Println("✅ Connected to PostgreSQL")

	// ---- Auto-migrate tables ----
	if err := gdb.AutoMigrate(&models.Book{}); err != nil {
		log.Fatalf("auto-migrate failed: %v", err)
	}
	log.Println("📘 Auto-migration for Book completed")

	// Ensure uploads directory exists
	if err := os.MkdirAll("uploads", 0o755); err != nil {
		log.Fatalf("failed to create uploads dir: %v", err)
	}

	// ---- Gin HTTP server ----
	r := gin.Default()

	// limit upload size (8 MB)
	r.MaxMultipartMemory = 8 << 20
	// serve generated files at http://localhost:8080/uploads/<filename>
	r.Static("/uploads", "./uploads")

	// health check
	r.GET("/health", func(c *gin.Context) {
		if err := sqlDB.PingContext(c.Request.Context()); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"ok": false, "database": "down", "error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"ok": true, "database": "up"})
	})

	// GET /books — list all books
	r.GET("/books", func(c *gin.Context) {
		var books []models.Book
		if err := gdb.Find(&books).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"count": len(books),
			"data":  books,
		})
	})

	// GET /books/:id — get a single book by ID
	r.GET("/books/:id", func(c *gin.Context) {
		id := c.Param("id")

		var book models.Book
		if err := gdb.First(&book, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
			return
		}

		c.JSON(http.StatusOK, book)
	})

	// PUT /books/:id — update book info (supports re-uploading cover image)
	r.PUT("/books/:id", func(c *gin.Context) {
		id := c.Param("id")

		var book models.Book
		if err := gdb.First(&book, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
			return
		}

		// Read new form values (multipart form)
		title := c.PostForm("title")
		author := c.PostForm("author")
		category := c.PostForm("category")
		priceStr := c.PostForm("price")
		stockStr := c.PostForm("stock")

		if title != "" {
			book.Title = title
		}
		if author != "" {
			book.Author = author
		}
		if category != "" {
			book.Category = category
		}

		if priceStr != "" {
			if price, err := strconv.ParseFloat(priceStr, 64); err == nil {
				book.Price = price
			}
		}
		if stockStr != "" {
			if stock, err := strconv.Atoi(stockStr); err == nil {
				book.Stock = stock
			}
		}

		// Optional: new cover image upload
		file, err := c.FormFile("cover")
		if err == nil && file != nil {
			ext := strings.ToLower(filepath.Ext(file.Filename))
			filename := uuid.New().String() + ext
			dst := filepath.Join("uploads", filename)

			// delete old cover (if exists)
			if book.CoverURL != "" {
				oldFile := "." + book.CoverURL
				_ = os.Remove(oldFile)
			}

			// save new file
			if err := c.SaveUploadedFile(file, dst); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save new cover"})
				return
			}
			book.CoverURL = "/uploads/" + filename
		}

		// Save updates
		if err := gdb.Save(&book).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "book updated successfully",
			"data":    book,
		})
	})

	// DELETE /books/:id — remove a book (and delete its cover image)
	r.DELETE("/books/:id", func(c *gin.Context) {
		id := c.Param("id")

		var book models.Book
		if err := gdb.First(&book, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
			return
		}

		// Delete cover image file if exists
		if book.CoverURL != "" {
			filePath := "." + book.CoverURL // because stored as /uploads/filename
			if err := os.Remove(filePath); err != nil {
				log.Printf("⚠️ failed to delete image: %v", err)
			}
		}

		// Delete from DB
		if err := gdb.Delete(&book).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("book with id %s deleted successfully", id),
		})
	})

	// POST /books — create a new book (supports multipart form with 'cover' image)
	r.POST("/books", func(c *gin.Context) {
		// 1) read form fields
		title := c.PostForm("title")
		if title == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "title is required"})
			return
		}
		author := c.PostForm("author")
		category := c.PostForm("category")

		var price float64
		var stock int

		if p := c.PostForm("price"); p != "" {
			v, err := strconv.ParseFloat(p, 64)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "price must be a number"})
				return
			}
			price = v
		}
		if s := c.PostForm("stock"); s != "" {
			v, err := strconv.Atoi(s)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "stock must be an integer"})
				return
			}
			stock = v
		}

		// 2) optional file upload
		var coverURL string
		if file, err := c.FormFile("cover"); err == nil && file != nil {
			ext := strings.ToLower(filepath.Ext(file.Filename))
			if ext == "" {
				ext = ".jpg"
			}
			switch ext {
			case ".jpg", ".jpeg", ".png", ".webp":
			default:
				c.JSON(http.StatusBadRequest, gin.H{"error": "cover must be .jpg, .jpeg, .png, or .webp"})
				return
			}

			filename := uuid.New().String() + ext
			dst := filepath.Join("uploads", filename)

			if err := c.SaveUploadedFile(file, dst); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save cover: " + err.Error()})
				return
			}
			coverURL = "/uploads/" + filename
		}

		// 3) create record
		book := models.Book{
			Title:    title,
			Author:   author,
			Category: category,
			Price:    price,
			Stock:    stock,
			CoverURL: coverURL,
		}
		if err := gdb.Create(&book).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "created",
			"data":    book,
		})
	})

	// start server
	port = env("APP_PORT", "8080")
	log.Printf("🚀 Server running on http://localhost:%s", port)
	if err = r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
