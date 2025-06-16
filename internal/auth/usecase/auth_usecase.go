package usecase

import (
	"errors"
	"time"

	mongo "github.com/nightnice1st/testGridWhiz/internal/auth/repository"
	"github.com/nightnice1st/testGridWhiz/internal/pkg/jwt"
	"github.com/nightnice1st/testGridWhiz/internal/pkg/ratelimit"
	"github.com/nightnice1st/testGridWhiz/internal/pkg/validator"
	"github.com/nightnice1st/testGridWhiz/internal/users/domain"

	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase struct {
	userRepo    domain.UserRepository
	authRepo    *mongo.AuthRepository
	jwtSecret   string
	jwtExpiry   time.Duration
	rateLimiter *ratelimit.RateLimiter
}

func NewAuthUsecase(userRepo domain.UserRepository, authRepo *mongo.AuthRepository,
	jwtSecret string, jwtExpiry time.Duration, rateLimiter *ratelimit.RateLimiter) *AuthUsecase {
	return &AuthUsecase{
		userRepo:    userRepo,
		authRepo:    authRepo,
		jwtSecret:   jwtSecret,
		jwtExpiry:   jwtExpiry,
		rateLimiter: rateLimiter,
	}
}

func (u *AuthUsecase) Register(email, password, name string) (*domain.User, error) {
	// Validate email
	if err := validator.ValidateEmail(email); err != nil {
		return nil, err
	}

	// Validate password
	if err := validator.ValidatePassword(password); err != nil {
		return nil, err
	}

	// Check if user already exists
	existingUser, _ := u.userRepo.FindByEmail(email)
	if existingUser != nil {
		return nil, errors.New("user already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create new user
	user := &domain.User{
		Email:    email,
		Password: string(hashedPassword),
		Name:     name,
	}

	if err := u.userRepo.Create(user); err != nil {
		return nil, err
	}

	// Clear password before returning
	user.Password = ""
	return user, nil
}

func (u *AuthUsecase) Login(email, password string) (string, error) {
	// Check rate limit
	if !u.rateLimiter.Allow(email) {
		return "", errors.New("too many login attempts, please try again later")
	}

	// Record login attempt
	if err := u.authRepo.RecordLoginAttempt(email); err != nil {
		// Log error but don't fail login
	}

	// Find user
	user, err := u.userRepo.FindByEmail(email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	// Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	// Generate JWT token
	token, err := jwt.GenerateToken(user.ID, user.Email, u.jwtSecret, u.jwtExpiry)
	if err != nil {
		return "", err
	}

	// Reset login attempts on successful login
	u.authRepo.ResetLoginAttempts(email)

	return token, nil
}

func (u *AuthUsecase) Logout(token string) error {
	// Validate token first
	claims, err := jwt.ValidateToken(token, u.jwtSecret)
	if err != nil {
		return errors.New("invalid token")
	}

	// Revoke token
	return u.authRepo.RevokeToken(token, claims.UserID)
}

func (u *AuthUsecase) ValidateToken(token string) (*jwt.Claims, error) {
	// Check if token is revoked
	isRevoked, err := u.authRepo.IsTokenRevoked(token)
	if err != nil {
		return nil, err
	}

	if isRevoked {
		return nil, errors.New("token has been revoked")
	}

	// Validate token
	return jwt.ValidateToken(token, u.jwtSecret)
}
