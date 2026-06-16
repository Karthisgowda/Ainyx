package repository

import (
	"context"

	"github.com/Karthisgowda/Ainyx/db/sqlc"
)

type UserRepository interface {
	Create(ctx context.Context, params sqlc.CreateUserParams) (sqlc.User, error)
	GetByID(ctx context.Context, id int32) (sqlc.User, error)
	List(ctx context.Context, params sqlc.ListUsersParams) ([]sqlc.User, error)
	Update(ctx context.Context, params sqlc.UpdateUserParams) (sqlc.User, error)
	Delete(ctx context.Context, id int32) (int64, error)
}

type userRepository struct {
	queries *sqlc.Queries
}

func NewUserRepository(queries *sqlc.Queries) UserRepository {
	return &userRepository{queries: queries}
}

func (r *userRepository) Create(ctx context.Context, params sqlc.CreateUserParams) (sqlc.User, error) {
	return r.queries.CreateUser(ctx, params)
}

func (r *userRepository) GetByID(ctx context.Context, id int32) (sqlc.User, error) {
	return r.queries.GetUser(ctx, id)
}

func (r *userRepository) List(ctx context.Context, params sqlc.ListUsersParams) ([]sqlc.User, error) {
	return r.queries.ListUsers(ctx, params)
}

func (r *userRepository) Update(ctx context.Context, params sqlc.UpdateUserParams) (sqlc.User, error) {
	return r.queries.UpdateUser(ctx, params)
}

func (r *userRepository) Delete(ctx context.Context, id int32) (int64, error) {
	return r.queries.DeleteUser(ctx, id)
}
