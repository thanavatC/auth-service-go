package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/thanavatC/auth-service-go/model"
	"github.com/thanavatC/auth-service-go/repository"
	"golang.org/x/crypto/bcrypt"
)

type IAuthService interface {
	Register(req model.RegisterRequest) (model.AuthResponse, error)
	Login(req model.LoginRequest) (model.AuthResponse, error)
	InvalidateToken(token string) error
}

type AuthService struct {
	userRepo repository.IUserRepository
	jwtKey   []byte
}

func NewAuthService(userRepo repository.IUserRepository, jwtKey string) IAuthService {
	return &AuthService{
		userRepo: userRepo,
		jwtKey:   []byte(jwtKey),
	}
}

func (s *AuthService) Register(req model.RegisterRequest) (model.AuthResponse, error) {
	// Check if user already exists
	_, err := s.userRepo.FindByEmail(req.Email)
	if err == nil {
		return model.AuthResponse{}, errors.New("user already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return model.AuthResponse{}, err
	}

	// Create user
	user := model.User{
		ID:        uuid.New().String(),
		Email:     req.Email,
		Password:  string(hashedPassword),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Role:      "user",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.userRepo.Create(user); err != nil {
		return model.AuthResponse{}, err
	}

	// Generate token
	token, err := s.generateToken(user)
	if err != nil {
		return model.AuthResponse{}, err
	}

	return model.AuthResponse{
		Token: token,
		User:  user,
	}, nil
}

func (s *AuthService) Login(req model.LoginRequest) (model.AuthResponse, error) {
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return model.AuthResponse{}, errors.New("invalid credentials")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return model.AuthResponse{}, errors.New("invalid credentials")
	}

	// Generate token
	token, err := s.generateToken(user)
	if err != nil {
		return model.AuthResponse{}, err
	}

	return model.AuthResponse{
		Token: token,
		User:  user,
	}, nil
}

func (s *AuthService) generateToken(user model.User) (string, error) {
	claims := jwt.MapClaims{
		"sub":   user.ID,
		"email": user.Email,
		"role":  user.Role,
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtKey)
}

func (s *AuthService) InvalidateToken(token string) error {
	// Remove "Bearer " prefix if present
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	// Parse and validate the token
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.jwtKey, nil
	})

	if err != nil {
		return errors.New("invalid token")
	}

	// Add the token to a blacklist in Redis/database
	// Set a shorter expiration time
	// Or implement token revocation logic

	return nil
}
