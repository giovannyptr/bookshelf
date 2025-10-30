package books

import (
	"fmt"
	// "net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/giovannyptr/bookshelf/internal/api"
	"github.com/giovannyptr/bookshelf/models"
	"github.com/google/uuid"
)

type Handler struct {
	repo      *Repository
	uploadDir string
}

func NewHandler(repo *Repository, uploadDir string) *Handler {
	return &Handler{repo: repo, uploadDir: uploadDir}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	r.MaxMultipartMemory = 8 << 20
	r.Static("/uploads", h.uploadDir)

	g := r.Group("/books")
	g.GET("", h.list)
	g.GET("/:id", h.detail)
	g.POST("", h.create)
	g.PUT("/:id", h.update)
	g.DELETE("/:id", h.delete)
}

// Exported wrappers so cmd/server can wire routes per auth group.
func (h *Handler) List(c *gin.Context)   { h.list(c) }
func (h *Handler) Detail(c *gin.Context) { h.detail(c) }
func (h *Handler) Create(c *gin.Context) { h.create(c) }
func (h *Handler) Update(c *gin.Context) { h.update(c) }
func (h *Handler) Delete(c *gin.Context) { h.delete(c) }

// ------- handlers -------

func (h *Handler) list(c *gin.Context) {
	q := strings.TrimSpace(c.Query("q"))
	category := c.Query("category")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	sort := c.DefaultQuery("sort", "created_at")
	order := strings.ToUpper(c.DefaultQuery("order", "DESC"))

	items, total, err := h.repo.List(q, category, page, limit, sort, order)
	if err != nil {
		api.Fail(c, 500, err.Error())
		return
	}

	api.OK(c, gin.H{
		"items": items, "total": total, "page": page, "limit": limit,
	})
}

func (h *Handler) detail(c *gin.Context) {
	id := c.Param("id")
	b, err := h.repo.ByID(id)
	if err != nil {
		api.Fail(c, 404, "book not found")
		return
	}
	api.OK(c, b)
}

func (h *Handler) create(c *gin.Context) {
	title := c.PostForm("title")
	if title == "" {
		api.Fail(c, 400, "title is required")
		return
	}
	author := c.PostForm("author")
	category := c.PostForm("category")

	var price float64
	var stock int
	if s := c.PostForm("price"); s != "" {
		v, err := strconv.ParseFloat(s, 64)
		if err != nil {
			api.Fail(c, 400, "price must be a number")
			return
		}
		price = v
	}
	if s := c.PostForm("stock"); s != "" {
		v, err := strconv.Atoi(s)
		if err != nil {
			api.Fail(c, 400, "stock must be an integer")
			return
		}
		stock = v
	}

	coverURL := ""
	if file, err := c.FormFile("cover"); err == nil && file != nil {
		ext := strings.ToLower(filepath.Ext(file.Filename))
		if ext == "" {
			ext = ".jpg"
		}
		switch ext {
		case ".jpg", ".jpeg", ".png", ".webp":
		default:
			api.Fail(c, 400, "cover must be .jpg/.jpeg/.png/.webp")
			return
		}
		filename := uuid.New().String() + ext
		dst := filepath.Join(h.uploadDir, filename)
		if err := c.SaveUploadedFile(file, dst); err != nil {
			api.Fail(c, 500, "failed to save cover")
			return
		}
		coverURL = "/uploads/" + filename
	}

	b := models.Book{Title: title, Author: author, Category: category, Price: price, Stock: stock, CoverURL: coverURL}
	if err := h.repo.Create(&b); err != nil {
		api.Fail(c, 500, err.Error())
		return
	}
	api.Created(c, b)
}

func (h *Handler) update(c *gin.Context) {
	id := c.Param("id")
	b, err := h.repo.ByID(id)
	if err != nil {
		api.Fail(c, 404, "book not found")
		return
	}

	if v := c.PostForm("title"); v != "" {
		b.Title = v
	}
	if v := c.PostForm("author"); v != "" {
		b.Author = v
	}
	if v := c.PostForm("category"); v != "" {
		b.Category = v
	}
	if s := c.PostForm("price"); s != "" {
		if v, err := strconv.ParseFloat(s, 64); err == nil {
			b.Price = v
		}
	}
	if s := c.PostForm("stock"); s != "" {
		if v, err := strconv.Atoi(s); err == nil {
			b.Stock = v
		}
	}

	if file, err := c.FormFile("cover"); err == nil && file != nil {
		ext := strings.ToLower(filepath.Ext(file.Filename))
		if ext == "" {
			ext = ".jpg"
		}
		switch ext {
		case ".jpg", ".jpeg", ".png", ".webp":
		default:
			api.Fail(c, 400, "cover must be .jpg/.jpeg/.png/.webp")
			return
		}
		filename := uuid.New().String() + ext
		dst := filepath.Join(h.uploadDir, filename)
		if b.CoverURL != "" {
			_ = os.Remove("." + b.CoverURL)
		}
		if err := c.SaveUploadedFile(file, dst); err != nil {
			api.Fail(c, 500, "failed to save new cover")
			return
		}
		b.CoverURL = "/uploads/" + filename
	}

	if err := h.repo.Save(&b); err != nil {
		api.Fail(c, 500, err.Error())
		return
	}
	api.OK(c, b)
}

func (h *Handler) delete(c *gin.Context) {
	id := c.Param("id")
	b, err := h.repo.ByID(id)
	if err != nil {
		api.Fail(c, 404, "book not found")
		return
	}

	if b.CoverURL != "" {
		_ = os.Remove("." + b.CoverURL)
	}
	if err := h.repo.Delete(&b); err != nil {
		api.Fail(c, 500, err.Error())
		return
	}
	api.OK(c, gin.H{"message": fmt.Sprintf("book %s deleted", id)})
}
