package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	v1 "github.com/AldiandyaIrsyad/author-notes/api/v1/auth"
)

// MockAuthRepository is a mock implementation of AuthRepository
type MockAuthRepository struct {
	mock.Mock
}

func (m *MockAuthRepository) CreateUser(ctx context.Context, user *User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockAuthRepository) FindByUsername(ctx context.Context, username string) (*User, error) {
	args := m.Called(ctx, username)
	if user := args.Get(0); user != nil {
		return user.(*User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockAuthRepository) FindByEmailOrUsername(ctx context.Context, email, username string) (*User, error) {
	args := m.Called(ctx, email, username)
	if user := args.Get(0); user != nil {
		return user.(*User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockAuthRepository) FindByID(ctx context.Context, id string) (*User, error) {
	args := m.Called(ctx, id)
	if user := args.Get(0); user != nil {
		return user.(*User), args.Error(1)
	}
	return nil, args.Error(1)
}

// Helper to create a hashed password for tests
func hashPasswordForTest(t *testing.T, password string) string {
	t.Helper()
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost) // Use MinCost for tests
	require.NoError(t, err)
	return string(hashed)
}

func TestAuthService_Register(t *testing.T) {
	mockRepo := new(MockAuthRepository)
	// Set JWT_SECRET for testing, ideally use a test-specific config
	t.Setenv("JWT_SECRET", "test_secret")
	service := NewAuthService(mockRepo)
	ctx := context.Background()

	registerReq := v1.RegisterRequest{
		Username: "testuser",
		Password: "password123",
		Email:    "test@example.com",
	}

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("FindByEmailOrUsername", ctx, registerReq.Email, registerReq.Username).Return(nil, ErrUserNotFound).Once()
		// We expect CreateUser to be called, but we don't need to inspect the user details deeply here,
		// just ensure it's called with a User object and returns no error.
		mockRepo.On("CreateUser", ctx, mock.AnythingOfType("*auth.User")).Return(nil).Once()

		user, err := service.Register(ctx, registerReq)

		require.NoError(t, err)
		require.NotNil(t, user)
		assert.Equal(t, registerReq.Username, user.Username)
		assert.Equal(t, registerReq.Email, user.Email)
		assert.Empty(t, user.Password, "Password should be empty in the response") // Important security check
		assert.NotZero(t, user.CreatedAt)
		assert.NotZero(t, user.UpdatedAt)
		mockRepo.AssertExpectations(t) // Verify that the expected methods were called
	})

	t.Run("Validation Failed - Short Username", func(t *testing.T) {
		invalidReq := v1.RegisterRequest{Username: "us", Password: "password123", Email: "test@example.com"}
		user, err := service.Register(ctx, invalidReq)
		require.Error(t, err)
		assert.Nil(t, user)
		assert.ErrorIs(t, err, ErrValidationFailed)
	})

	t.Run("Validation Failed - Short Password", func(t *testing.T) {
		invalidReq := v1.RegisterRequest{Username: "testuser", Password: "pass", Email: "test@example.com"}
		user, err := service.Register(ctx, invalidReq)
		require.Error(t, err)
		assert.Nil(t, user)
		assert.ErrorIs(t, err, ErrValidationFailed)
	})

	t.Run("Validation Failed - Invalid Email", func(t *testing.T) {
		invalidReq := v1.RegisterRequest{Username: "testuser", Password: "password123", Email: "invalid-email"}
		user, err := service.Register(ctx, invalidReq)
		require.Error(t, err)
		assert.Nil(t, user)
		assert.ErrorIs(t, err, ErrValidationFailed)
	})

	t.Run("User Already Exists", func(t *testing.T) {
		existingUser := &User{ID: "1", Username: "testuser", Email: "test@example.com"}
		mockRepo.On("FindByEmailOrUsername", ctx, registerReq.Email, registerReq.Username).Return(existingUser, nil).Once()

		user, err := service.Register(ctx, registerReq)

		require.Error(t, err)
		assert.Nil(t, user)
		assert.ErrorIs(t, err, ErrUserAlreadyExists)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Repository Find Error", func(t *testing.T) {
		dbError := errors.New("database connection error")
		mockRepo.On("FindByEmailOrUsername", ctx, registerReq.Email, registerReq.Username).Return(nil, dbError).Once()

		user, err := service.Register(ctx, registerReq)

		require.Error(t, err)
		assert.Nil(t, user)
		assert.NotErrorIs(t, err, ErrUserAlreadyExists) // Ensure it's not the specific domain error
		assert.NotErrorIs(t, err, ErrValidationFailed)
		assert.ErrorIs(t, err, dbError) // Or check for a wrapped internal error if you implement that
		mockRepo.AssertExpectations(t)
	})

	t.Run("Repository Create Error", func(t *testing.T) {
		dbError := errors.New("failed to insert user")
		mockRepo.On("FindByEmailOrUsername", ctx, registerReq.Email, registerReq.Username).Return(nil, ErrUserNotFound).Once()
		mockRepo.On("CreateUser", ctx, mock.AnythingOfType("*auth.User")).Return(dbError).Once()

		user, err := service.Register(ctx, registerReq)

		require.Error(t, err)
		assert.Nil(t, user)
		// Check for the generic error returned by the service in this case
		assert.EqualError(t, err, "failed to save user")
		mockRepo.AssertExpectations(t)
	})
}

func TestAuthService_Login(t *testing.T) {
	mockRepo := new(MockAuthRepository)
	t.Setenv("JWT_SECRET", "test_secret_for_login")
	service := NewAuthService(mockRepo) // Recreate service to pick up env var
	ctx := context.Background()

	loginReq := v1.LoginRequest{
		Username: "testuser",
		Password: "password123",
	}
	hashedPassword := hashPasswordForTest(t, loginReq.Password)
	existingUser := &User{
		ID:        "user123",
		Username:  loginReq.Username,
		Password:  hashedPassword,
		Email:     "test@example.com",
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("FindByUsername", ctx, loginReq.Username).Return(existingUser, nil).Once()

		token, err := service.Login(ctx, loginReq)

		require.NoError(t, err)
		require.NotEmpty(t, token)
		// Basic check: Add more robust JWT validation if needed (e.g., parse and check claims)
		assert.Greater(t, len(token), 20)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Validation Failed - Missing Username", func(t *testing.T) {
		invalidReq := v1.LoginRequest{Password: "password123"}
		token, err := service.Login(ctx, invalidReq)
		require.Error(t, err)
		assert.Empty(t, token)
		assert.ErrorIs(t, err, ErrValidationFailed)
	})

	t.Run("Validation Failed - Missing Password", func(t *testing.T) {
		invalidReq := v1.LoginRequest{Username: "testuser"}
		token, err := service.Login(ctx, invalidReq)
		require.Error(t, err)
		assert.Empty(t, token)
		assert.ErrorIs(t, err, ErrValidationFailed)
	})

	t.Run("User Not Found", func(t *testing.T) {
		mockRepo.On("FindByUsername", ctx, loginReq.Username).Return(nil, ErrUserNotFound).Once()

		token, err := service.Login(ctx, loginReq)

		require.Error(t, err)
		assert.Empty(t, token)
		assert.ErrorIs(t, err, ErrInvalidCredentials) // Should return invalid credentials, not user not found
		mockRepo.AssertExpectations(t)
	})

	t.Run("Repository Find Error", func(t *testing.T) {
		dbError := errors.New("database connection error")
		mockRepo.On("FindByUsername", ctx, loginReq.Username).Return(nil, dbError).Once()

		token, err := service.Login(ctx, loginReq)

		require.Error(t, err)
		assert.Empty(t, token)
		assert.NotErrorIs(t, err, ErrInvalidCredentials)
		assert.NotErrorIs(t, err, ErrValidationFailed)
		// Check for the generic error returned by the service
		assert.EqualError(t, err, "login failed")
		mockRepo.AssertExpectations(t)
	})

	t.Run("Incorrect Password", func(t *testing.T) {
		mockRepo.On("FindByUsername", ctx, loginReq.Username).Return(existingUser, nil).Once()
		incorrectLoginReq := v1.LoginRequest{
			Username: loginReq.Username,
			Password: "wrongpassword",
		}

		token, err := service.Login(ctx, incorrectLoginReq)

		require.Error(t, err)
		assert.Empty(t, token)
		assert.ErrorIs(t, err, ErrInvalidCredentials)
		mockRepo.AssertExpectations(t)
	})
}
