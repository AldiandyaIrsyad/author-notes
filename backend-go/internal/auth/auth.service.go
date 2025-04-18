package auth

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	v1 "github.com/AldiandyaIrsyad/author-notes/api/v1/auth"
)

// AuthService defines the interface for authentication related business logic.
type AuthService interface {
	Register(ctx context.Context, req v1.RegisterRequest) (*User, error)
	Login(ctx context.Context, req v1.LoginRequest) (string, error)
}

type authService struct {
	repo      AuthRepository
	validator *validator.Validate
	jwtSecret []byte
}

// NewAuthService creates a new instance of AuthService.
func NewAuthService(repo AuthRepository) AuthService {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		// Provide a default secret for development, but log a warning.
		// In production, this should cause a fatal error or use a secure default.
		jwtSecret = "default_dev_secret_key_please_change"
		// Consider adding logging here: log.Println("Warning: JWT_SECRET environment variable not set. Using default.")
	}
	return &authService{
		repo:      repo,
		validator: validator.New(),
		jwtSecret: []byte(jwtSecret),
	}
}

// Register handles user registration.
func (s *authService) Register(ctx context.Context, req v1.RegisterRequest) (*User, error) {
	// 1. Validate input
	if err := s.validator.Struct(req); err != nil {
		// Consider wrapping or logging the specific validation errors
		return nil, ErrValidationFailed
	}

	// 2. Check if user already exists (by email or username)
	existingUser, err := s.repo.FindByEmailOrUsername(ctx, req.Email, req.Username)
	if err != nil && !errors.Is(err, ErrUserNotFound) {
		// Handle potential database errors
		return nil, err // Or wrap in a more generic internal error
	}
	if existingUser != nil {
		return nil, ErrUserAlreadyExists
	}

	// 3. Hash password
	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		// Log error: log.Printf("Error hashing password: %v", err)
		return nil, errors.New("failed to process registration") // Generic error
	}

	// 4. Create user model
	now := time.Now().Unix()
	user := &User{
		Email:     req.Email,
		Username:  req.Username,
		Password:  hashedPassword,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// 5. Save user to repository
	if err := s.repo.CreateUser(ctx, user); err != nil {
		// Handle potential database errors (e.g., duplicate key if check failed due to race condition)
		// Log error: log.Printf("Error creating user: %v", err)
		return nil, errors.New("failed to save user") // Generic error
	}

	// Important: Clear password before returning
	user.Password = ""
	return user, nil
}

// Login handles user login.
func (s *authService) Login(ctx context.Context, req v1.LoginRequest) (string, error) {
	// 1. Validate input
	if err := s.validator.Struct(req); err != nil {
		return "", ErrValidationFailed
	}

	// 2. Find user by username
	user, err := s.repo.FindByUsername(ctx, req.Username)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return "", ErrInvalidCredentials // User not found treated as invalid credentials
		}
		// Log error: log.Printf("Error finding user by username: %v", err)
		return "", errors.New("login failed") // Generic internal error
	}

	// 3. Compare password
	if !checkPasswordHash(req.Password, user.Password) {
		return "", ErrInvalidCredentials
	}

	// 4. Generate JWT token
	token, err := s.generateJWT(user)
	if err != nil {
		// Log error: log.Printf("Error generating JWT token: %v", err)
		return "", errors.New("login failed") // Generic internal error
	}

	return token, nil
}

// hashPassword generates a bcrypt hash of the password.
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// checkPasswordHash compares a plain text password with a bcrypt hash.
func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// generateJWT creates a new JWT token for the given user.
func (s *authService) generateJWT(user *User) (string, error) {
	// Create the claims
	claims := jwt.MapClaims{
		"sub": user.ID,                               // Subject (user ID)
		"usr": user.Username,                         // Username
		"exp": time.Now().Add(time.Hour * 72).Unix(), // Expiration time (e.g., 72 hours)
		"iat": time.Now().Unix(),                     // Issued at
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token with secret
	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
