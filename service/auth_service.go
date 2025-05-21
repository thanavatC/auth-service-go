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
