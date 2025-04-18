package auth

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	v1 "github.com/AldiandyaIrsyad/author-notes/api/v1/auth"
	app_auth "github.com/AldiandyaIrsyad/author-notes/internal/auth"
)

type AuthHTTPHandler struct {
	service app_auth.AuthService
}

func NewAuthHTTPHandler(service app_auth.AuthService) *AuthHTTPHandler {
	return &AuthHTTPHandler{service: service}
}

func (h *AuthHTTPHandler) RegisterRoutes(rg *gin.RouterGroup) {
	authGroup := rg.Group("/auth")
	authGroup.POST("/register", h.Register)
	authGroup.POST("/login", h.Login)
}

// Register handles the user registration request.
// @Summary Register a new user
// @Description Creates a new user account.
// @Tags auth
// @Accept json
// @Produce json
// @Param user body v1.RegisterRequest true "User registration details"
// @Success 201 {object} auth.User "User created successfully (Password field will be empty)"
// @Failure 400 {object} map[string]string "Validation error or bad request"
// @Failure 409 {object} map[string]string "User already exists"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /v1/auth/register [post]
func (h *AuthHTTPHandler) Register(c *gin.Context) {
	var req v1.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	user, err := h.service.Register(c.Request.Context(), req)
	if err != nil {

		switch {
		case errors.Is(err, app_auth.ErrValidationFailed): //
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, app_auth.ErrUserAlreadyExists):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			// Log the internal error (consider adding logging)
			// log.Printf("Internal server error during registration: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		}
		return
	}

	user.Password = "" // Service should handle this
	c.JSON(http.StatusCreated, user)
}

// Login handles the user login request.
// @Summary Log in a user
// @Description Authenticates a user and returns a JWT token.
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body v1.LoginRequest true "User login credentials"
// @Success 200 {object} v1.LoginResponse "Login successful, JWT token returned"
// @Failure 400 {object} map[string]string "Validation error or bad request"
// @Failure 401 {object} map[string]string "Invalid credentials"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /v1/auth/login [post]
func (h *AuthHTTPHandler) Login(c *gin.Context) {
	var req v1.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	token, err := h.service.Login(c.Request.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, app_auth.ErrValidationFailed):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, app_auth.ErrInvalidCredentials):
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		default:
			// Log the internal error (consider adding logging)
			// log.Printf("Internal server error during login: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to login"})
		}
		return
	}

	c.JSON(http.StatusOK, v1.LoginResponse{Token: token})
}
