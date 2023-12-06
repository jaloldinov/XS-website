package user

import (
	"context"
	"xs/internal/pkg"
	"xs/internal/repository/postgres/user"
)

type User interface {
	UserCreate(ctx context.Context, data user.CreateUserRequest) (user.CreateUserResponse, *pkg.Error)
	UserGetById(ctx context.Context, id string) (user.GetUserResponse, *pkg.Error)
	UserGetAll(ctx context.Context, filter user.Filter) ([]user.GetUserListResponse, int, *pkg.Error)
	UserUpdate(ctx context.Context, data user.UpdateUserRequest) *pkg.Error
	UserDelete(ctx context.Context, req user.DeleteUserRequest) *pkg.Error
	UserUpdatePassword(ctx context.Context, req user.UpdatePasswordRequest) *pkg.Error
}
