package auth

import (
	util "go-fiber-api/internal/util/password"

	"github.com/google/uuid"
)

type AuthService interface {
	RegisterUser(req *RegisterUserRequest) error
	LoginUser(req *LoginRequest) (*AuthResponse, error)
	GetUser(id uuid.UUID) (*AuthResponse, error)
}

type authService struct {
	authRepo UserRepository
}

func NewAuthService(repo UserRepository) AuthService {
	return &authService{
		authRepo: repo,
	}
}

func (s *authService) RegisterUser(req *RegisterUserRequest) error {
	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return err
	}

	user := &User{
		Email:    req.Email,
		Password: hashedPassword,
	}

	return s.authRepo.RegisterUser(user)
}

func (s *authService) LoginUser(req *LoginRequest) (*AuthResponse, error) {
	user, err := s.authRepo.FindByEmail(req.Email)

	if err != nil {
		return nil, err
	}

	err = util.CheckPassword(user.Password, req.Password)

	if err != nil {
		return nil, err
	}
	result := &AuthResponse{
		ID:        user.ID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Time.String(),
		UpdatedAt: user.UpdatedAt.Time.String(),
	}

	return result, nil
}

func (s *authService) GetUser(id uuid.UUID) (*AuthResponse, error) {
	user, err := s.authRepo.FindByID(id)

	result := &AuthResponse{
		ID:        user.ID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Time.String(),
		UpdatedAt: user.UpdatedAt.Time.String(),
	}

	if err != nil {
		return nil, err
	}

	return result, nil
}
