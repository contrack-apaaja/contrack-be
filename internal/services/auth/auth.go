package auth

import (
	"database/sql"
	"fmt"
	"log"

	"contrack-be/internal/database"
	"contrack-be/internal/models"
	jwtService "contrack-be/internal/services/jwt"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	jwtService *jwtService.Service
}

// NewAuthService creates a new authentication service
func NewAuthService(jwtSvc *jwtService.Service) *Service {
	return &Service{
		jwtService: jwtSvc,
	}
}

// Register creates a new user account
func (s *Service) Register(req *models.UserRegistrationRequest) (*models.UserResponse, string, error) {
	// Check if user already exists
	existingUser, err := s.getUserByEmail(req.Email)
	if err != nil && err != sql.ErrNoRows {
		return nil, "", fmt.Errorf("failed to check existing user: %w", err)
	}
	if existingUser != nil {
		return nil, "", fmt.Errorf("user with email %s already exists", req.Email)
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user with default REGULAR role
	query := `
		INSERT INTO users (email, password, role)
		VALUES ($1, $2, $3)
		RETURNING id, email, role, created_at, updated_at
	`

	var user models.User
	err = database.DB.QueryRow(query, req.Email, string(hashedPassword), models.GetDefaultRole()).Scan(
		&user.ID, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create user: %w", err)
	}

	// Generate JWT token
	token, err := s.jwtService.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token: %w", err)
	}

	log.Printf("User registered successfully: %s", user.Email)
	return user.ToResponse(), token, nil
}

// Login authenticates a user and returns a JWT token
func (s *Service) Login(req *models.UserLoginRequest) (*models.UserResponse, string, error) {
	// Get user by email
	user, err := s.getUserByEmail(req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, "", fmt.Errorf("invalid credentials")
		}
		return nil, "", fmt.Errorf("failed to get user: %w", err)
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, "", fmt.Errorf("invalid credentials")
	}

	// Generate JWT token
	token, err := s.jwtService.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token: %w", err)
	}

	log.Printf("User logged in successfully: %s", user.Email)
	return user.ToResponse(), token, nil
}

// GetUserByID retrieves a user by their ID
func (s *Service) GetUserByID(userID string) (*models.UserResponse, error) {
	query := `
		SELECT id, email, role, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var user models.User
	err := database.DB.QueryRow(query, userID).Scan(
		&user.ID, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user.ToResponse(), nil
}

// RefreshToken generates a new token for the user
func (s *Service) RefreshToken(tokenString string) (string, error) {
	return s.jwtService.RefreshToken(tokenString)
}

// getUserByEmail is a helper method to get user by email (includes password)
func (s *Service) getUserByEmail(email string) (*models.User, error) {
	query := `
		SELECT id, email, password, role, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	var user models.User
	err := database.DB.QueryRow(query, email).Scan(
		&user.ID, &user.Email, &user.Password, &user.Role, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
