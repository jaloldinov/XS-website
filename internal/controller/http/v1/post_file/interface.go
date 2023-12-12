package post_file

import (
	"context"
	"xs/internal/pkg"
	"xs/internal/repository/postgres/post_file"
)

type PostFile interface {
	PostFileCreate(ctx context.Context, data post_file.CreatePostFileRequest) (post_file.CreatePostFileResponse, *pkg.Error)
	PostFileGetById(ctx context.Context, id string) (post_file.GetPostFileResponse, *pkg.Error)
	PostFileGetAll(ctx context.Context, filter post_file.Filter) ([]post_file.GetPostFileListResponse, int, *pkg.Error)
	PostFileUpdate(ctx context.Context, data post_file.UpdatePostFileRequest) *pkg.Error
	PostFileDelete(ctx context.Context, id string) *pkg.Error
}
