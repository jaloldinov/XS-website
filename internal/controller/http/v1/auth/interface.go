package auth

import (
	"xs/internal/auth"
	"xs/internal/pkg"
	"xs/internal/repository/postgres/user"
	"context"
)

type Auth interface {
	GenerateToken(ctx context.Context, data auth.GenerateToken) (string, error)
}

type User interface {
	GetUserByUsername(ctx context.Context, username string) (user.DetailUserResponse, *pkg.Error)
}
