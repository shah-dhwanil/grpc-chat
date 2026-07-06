package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/shah-dhwanil/grpc-chat/apps/user-service/internal/dto"
	"github.com/shah-dhwanil/grpc-chat/packages/database/postgres"
	"github.com/shah-dhwanil/grpc-chat/packages/pkgerror"
	errs "github.com/shah-dhwanil/grpc-chat/apps/user-service/internal/pkgerror"
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
		return nil, pkgerror.NewInternalError(err,"UUID_GEN_ERROR","Error while generating uuid v7",map[string]any{
			"operation":"user.create_user",
		})
	}
	args,err := postgres.StructToNamedArgs(user)
	if 	err != nil {
		return nil, postgres.NewStructToPayloadConversionError(err,"user.create_user")
	}
	args["id"] = id
	rows, _ := r.db.Query(ctx, createUserQuery, args)
	record,err := pgx.CollectOneRow(rows, pgx.RowToStructByName[dto.User])
	if err != nil {
		return nil, mapErrorToRepositoryError(err)
	}
	return &record, nil
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
	rows, _ := r.db.Query(ctx, getUserByIDQuery, args)
	record,err:= pgx.CollectOneRow(rows, pgx.RowToStructByName[dto.User])
	if err != nil {
		return nil, mapErrorToRepositoryError(err)
	}
	return &record, nil
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

	rows, _ := r.db.Query(ctx, getUsersQuery, args)
	records,err:= pgx.CollectRows(rows, pgx.RowToStructByName[dto.User])
	if err != nil {
		return nil, mapErrorToRepositoryError(err)
	}
	return records, nil
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
	rows, _ := r.db.Query(ctx, query, args)
	record,err:= pgx.CollectOneRow(rows, pgx.RowToStructByName[dto.User])
	if err != nil {
		return nil, mapErrorToRepositoryError(err)
	}
	return &record, nil
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
	_,err := r.db.Exec(ctx, deleteUserQuery, args)
	if err != nil {
		return mapErrorToRepositoryError(err)
	}
	return nil
}

func mapErrorToRepositoryError(err error) error {
	dbErr, ok := postgres.ConvertPgError(err)
	if !ok {
		return pkgerror.NewUnknownError(err, "DATABASE_ERROR", "Unknown Error while fetching record from postgres", nil)
	}
	pgErr, ok := dbErr.(*postgres.DatabaseError)
	if !ok {
		return pkgerror.NewUnknownError(err, "DATABASE_ERROR", "Unknown Error while fetching record from postgres", nil)
	}
	switch pgErr.Code {
	case postgres.NoRecordsFound:
		return errs.NewUserNotFoundError(err)
	case postgres.UniqueViolation:
		switch pgErr.ConstraintName {
			case "uq_users_primary_email":
				return errs.NewUserAlreadyExistsError(pgErr)
		}
	default:
		return pkgerror.NewUnknownError(pgErr, "DATABASE_ERROR", "Unknown Error while fetching record from postgres", nil)
	}
	return pkgerror.NewUnknownError(pgErr, "DATABASE_ERROR", "Unknown Error while fetching record from postgres", nil)
}