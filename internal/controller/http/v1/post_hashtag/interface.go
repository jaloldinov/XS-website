package post_hashtag

import (
	"context"
	"xs/internal/pkg"
	"xs/internal/repository/postgres/post_hashtag"
)

type PostHashtag interface {
	PostHashtagCreate(ctx context.Context, data post_hashtag.CreatePostHashtagRequest) (post_hashtag.CreatePostHashtagResponse, *pkg.Error)
	PostHashtagGetAll(ctx context.Context, post_id string) ([]post_hashtag.GetPostHashtagListResponse, int, *pkg.Error)
	PostHashtagDelete(ctx context.Context, id string) *pkg.Error
}
