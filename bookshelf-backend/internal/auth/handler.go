package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/giovannyptr/bookshelf/internal/api"
	"github.com/giovannyptr/bookshelf/internal/users"
	"github.com/giovannyptr/bookshelf/models"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	users *users.Repository
}

func NewHandler(ur *users.Repository) *Handler { return &Handler{users: ur} }

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	g := r.Group("/auth")
	g.POST("/register", h.register)
	g.POST("/login", h.login)
	g.GET("/me", AuthRequired(), h.me)
}

type registerDTO struct {
	Email    string `form:"email" json:"email" binding:"required,email"`
	Password string `form:"password" json:"password" binding:"required,min=6"`
	Name     string `form:"name" json:"name"`
}

func (h *Handler) register(c *gin.Context) {
	var in registerDTO
	if err := c.ShouldBind(&in); err != nil {
		api.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	in.Email = strings.ToLower(strings.TrimSpace(in.Email))

	// ensure unique
	if _, err := h.users.ByEmail(in.Email); err == nil {
		api.Fail(c, http.StatusConflict, "email already registered")
		return
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	u := models.User{Email: in.Email, Password: string(hash), Name: in.Name, Role: "user"}
	if err := h.users.Create(&u); err != nil {
		api.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	token, _ := GenerateToken(u.ID, u.Email, u.Role)
	api.Created(c, gin.H{"token": token, "user": gin.H{"id": u.ID, "email": u.Email, "name": u.Name, "role": u.Role}})
}

type loginDTO struct {
	Email    string `form:"email" json:"email" binding:"required,email"`
	Password string `form:"password" json:"password" binding:"required"`
}

func (h *Handler) login(c *gin.Context) {
	var in loginDTO
	if err := c.ShouldBind(&in); err != nil {
		api.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	in.Email = strings.ToLower(strings.TrimSpace(in.Email))
	u, err := h.users.ByEmail(in.Email)
	if err != nil {
		api.Fail(c, http.StatusUnauthorized, "invalid credentials")
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(in.Password)); err != nil {
		api.Fail(c, http.StatusUnauthorized, "invalid credentials")
		return
	}
	token, _ := GenerateToken(u.ID, u.Email, u.Role)
	api.OK(c, gin.H{"token": token, "user": gin.H{"id": u.ID, "email": u.Email, "name": u.Name, "role": u.Role}})
}

func (h *Handler) me(c *gin.Context) {
	uid, _ := GetUserID(c)
	role, _ := GetUserRole(c)
	api.OK(c, gin.H{"id": uid, "role": role})
}
