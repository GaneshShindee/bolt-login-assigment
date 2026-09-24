package service

import (
	"context"

	"github.com/ganeshshinde/boltapp/backend/internal/biz"
)

type UserService struct {
	users *biz.UserBiz
}

func NewUserService(users *biz.UserBiz) *UserService {
	return &UserService{users: users}
}

func (s *UserService) Register(ctx context.Context, req RegisterRequest) (RegisterResponse, error) {
	u, code, err := s.users.Register(ctx, req.Email, req.FirstName, req.LastName)
	if err != nil {
		return RegisterResponse{}, err
	}
	return RegisterResponse{User: toUserDTO(u), Code: code}, nil
}

func (s *UserService) Recognize(ctx context.Context, req RecognizeRequest) (RecognizeResponse, error) {
	recognized, err := s.users.Recognize(ctx, req.Email)
	if err != nil {
		return RecognizeResponse{}, err
	}
	return RecognizeResponse{Recognized: recognized}, nil
}

func (s *UserService) Login(ctx context.Context, req LoginRequest) (LoginResponse, error) {
	u, token, err := s.users.Login(ctx, req.Email, req.Code)
	if err != nil {
		return LoginResponse{}, err
	}
	return LoginResponse{User: toUserDTO(u), Token: token}, nil
}

func (s *UserService) Me(ctx context.Context, token string) (MeResponse, error) {
	u, err := s.users.Authenticate(ctx, token)
	if err != nil {
		return MeResponse{}, err
	}
	return MeResponse{User: toUserDTO(u)}, nil
}
