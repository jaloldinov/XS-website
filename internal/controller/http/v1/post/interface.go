package post

import (
	"context"
	"xs/internal/pkg"
	"xs/internal/repository/postgres/post"
)

type Post interface {
	PostCreate(ctx context.Context, data post.CreatePostRequest) (post.CreatePostResponse, *pkg.Error)
	PostGetById(ctx context.Context, id string) (post.GetPostResponse, *pkg.Error)
	PostGetAll(ctx context.Context, filter post.Filter) ([]post.GetPostListResponse, int, *pkg.Error)
	PostUpdate(ctx context.Context, data post.UpdatePostRequest) *pkg.Error
	PostDelete(ctx context.Context, id string) *pkg.Error
}
