package repository

import "github.com/shah-dhwanil/grpc-chat/packages/database/postgres"


type Repository struct {
	UserRepository *UserRepository
}

func NewRepository(pool postgres.PgPool) *Repository {
	return &Repository{
		UserRepository: NewUserRepository(pool),
	}
}