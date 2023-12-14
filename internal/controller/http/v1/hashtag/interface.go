package hashtag

import (
	"context"
	"xs/internal/pkg"
	"xs/internal/repository/postgres/hashtag"
)

type Hashtag interface {
	HashtagCreate(ctx context.Context, data hashtag.CreateHashtagRequest) (hashtag.CreateHashtagResponse, *pkg.Error)
	HashtagGetById(ctx context.Context, id string) (hashtag.GetHashtagResponse, *pkg.Error)
	HashtagGetAll(ctx context.Context, filter hashtag.Filter) ([]hashtag.GetHashtagListResponse, int, *pkg.Error)
	HashtagUpdate(ctx context.Context, data hashtag.UpdateHashtagRequest) *pkg.Error
	HashtagDelete(ctx context.Context, id string) *pkg.Error
}
