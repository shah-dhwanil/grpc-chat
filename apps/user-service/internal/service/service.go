package service

import "github.com/shah-dhwanil/grpc-chat/apps/user-service/internal/repository"

type Service struct{
	UserService *UserService
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		UserService: NewUserService(repo),
	}
}