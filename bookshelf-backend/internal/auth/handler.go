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
	Email    string `json:"email" example:"user@mail.com"`
	Password string `json:"password" example:"12345678"`
	Name     string `json:"name" example:"User Name"`
}

// register godoc
// @Summary Register new user
// @Tags    auth
// @Accept  json
// @Produce json
// @Param   payload body registerDTO true "Register payload"
// @Success 201 {object} map[string]any
// @Failure 400 {object} api.ErrorResponse
// @Failure 409 {object} api.ErrorResponse
// @Router  /auth/register [post]
func (h *Handler) register(c *gin.Context) {
	var in registerDTO
	if err := c.ShouldBindJSON(&in); err != nil {
		api.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	in.Email = strings.ToLower(strings.TrimSpace(in.Email))
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
	c.JSON(http.StatusCreated, gin.H{"token": token, "user": gin.H{"id": u.ID, "email": u.Email, "name": u.Name, "role": u.Role}})
}

type loginDTO struct {
	Email    string `json:"email" example:"admin@mail.com"`
	Password string `json:"password" example:"adminbookshelf"`
}

// login godoc
// @Summary Login
// @Tags    auth
// @Accept  json
// @Produce json
// @Param   payload body loginDTO true "Login payload"
// @Success 200 {object} map[string]any
// @Failure 400 {object} api.ErrorResponse
// @Failure 401 {object} api.ErrorResponse
// @Router  /auth/login [post]
func (h *Handler) login(c *gin.Context) {
	var in loginDTO
	if err := c.ShouldBindJSON(&in); err != nil {
		api.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	in.Email = strings.ToLower(strings.TrimSpace(in.Email))
	u, err := h.users.ByEmail(in.Email)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(in.Password)) != nil {
		api.Fail(c, http.StatusUnauthorized, "invalid credentials")
		return
	}
	token, _ := GenerateToken(u.ID, u.Email, u.Role)
	c.JSON(http.StatusOK, gin.H{"token": token, "user": gin.H{"id": u.ID, "email": u.Email, "name": u.Name, "role": u.Role}})
}

// me godoc
// @Summary Current user
// @Tags    auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]any
// @Router  /auth/me [get]
func (h *Handler) me(c *gin.Context) {
	uid, _ := GetUserID(c)
	role, _ := GetUserRole(c)
	c.JSON(http.StatusOK, gin.H{"id": uid, "role": role})
}
