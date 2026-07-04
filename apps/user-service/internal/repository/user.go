package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/shah-dhwanil/grpc-chat/apps/user-service/internal/dto"
	"github.com/shah-dhwanil/grpc-chat/packages/database/postgres"
)

type UserRepository struct{
	db postgres.DBTX
}

func NewUserRepository(pgPool postgres.PgPool) *UserRepository {
	return &UserRepository{
		db: pgPool,
	}
}

func (r *UserRepository) WithTransaction(tx postgres.DBTX) *UserRepository {
	return &UserRepository{
		db: tx,
	}
}

const createUserQuery = `
INSERT INTO users.users (id, name, primary_email)
VALUES (@id, @name, @primary_email)
RETURNING id, name, primary_email, created_at, updated_at
`

func (r *UserRepository) CreateUser(ctx context.Context,user *dto.CreateUserRequest) (*dto.User,error) {
	id,err := uuid.NewV7()
	if err != nil {
		return nil, fmt.Errorf("failed to generate user ID: %w", err)
	}
	args,err := postgres.StructToNamedArgs(user)
	if 	err != nil {
		return nil, fmt.Errorf("failed to convert user struct to named args: %w", err)
	}
	args["id"] = id
	rows, err := postgres.QueryInTransaction(ctx,r.db,
		func(executor postgres.Tx) (dto.User, error) {
			rows, _ := executor.Query(ctx, createUserQuery, args)
			return pgx.CollectOneRow(rows, pgx.RowToStructByName[dto.User])
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return &rows, nil
}

const getUserByIDQuery = `
SELECT id, name, primary_email, created_at, updated_at
FROM users.users
WHERE id = @id and is_deleted = false
`

func (r *UserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*dto.User, error) {
	args := pgx.NamedArgs{
		"id": id,
	}
	rows, err := postgres.QueryInTransaction(ctx,r.db,
		func(executor postgres.Tx) (dto.User, error) {
			rows, _ := executor.Query(ctx, getUserByIDQuery, args)
			return pgx.CollectOneRow(rows, pgx.RowToStructByName[dto.User])
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}
	return &rows, nil
}

const getUsersQuery = `
SELECT id, name, primary_email, created_at, updated_at
FROM users.users
WHERE id = ANY(@ids) and is_deleted = false
`

func (r *UserRepository) GetUsers(ctx context.Context, ids []uuid.UUID) ([]dto.User, error) {
	args := pgx.NamedArgs{
		"ids": ids,
	}
	rows, err := postgres.QueryInTransaction(ctx,r.db,
		func(executor postgres.Tx) ([]dto.User, error) {
			rows, _ := executor.Query(ctx, getUsersQuery, args)
			return pgx.CollectRows(rows, pgx.RowToStructByName[dto.User])
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}
	return rows, nil
}

const updateUserQuery = `
UPDATE users.users
SET %s
WHERE id = @id and is_deleted = false
RETURNING id, name, primary_email, created_at, updated_at
`

func (r *UserRepository) UpdateUser(ctx context.Context, id uuid.UUID, user *dto.UpdateUserRequest) (*dto.User, error) {
	args, err := postgres.StructToNamedArgs(user)
	if err != nil {
		return nil, fmt.Errorf("failed to convert user struct to named args: %w", err)
	}
	args["id"] = id

	setClause := make([]string, 0, 2)
	if user.Name != nil {
		setClause = append(setClause, "name = @name")
	}
	if user.PrimaryEmail != nil {
		setClause = append(setClause, "primary_email = @primary_email")
	}
	if len(setClause) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	query := fmt.Sprintf(updateUserQuery, postgres.ConstructSetClause(setClause))
	row,err:= postgres.QueryInTransaction(ctx,r.db,
		func(executor postgres.Tx) (dto.User, error) {
			rows, _ := executor.Query(ctx, query, args)
			return pgx.CollectOneRow(rows, pgx.RowToStructByName[dto.User])
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}
	return &row, nil
}

const deleteUserQuery = `
UPDATE users.users
SET is_deleted = true
WHERE id = @id and is_deleted = false
`

func (r *UserRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	args := pgx.NamedArgs{
		"id": id,
	}
	_, err := postgres.ExecuteInTransaction(ctx,r.db,
		func(executor postgres.Tx) (pgconn.CommandTag, error) {
			result, _ := executor.Exec(ctx, deleteUserQuery, args)
			return result, nil
		},
	)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}