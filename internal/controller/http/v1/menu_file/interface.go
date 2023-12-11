package menu_file

import (
	"context"
	"xs/internal/pkg"
	"xs/internal/repository/postgres/menu_file"
)

type MenuFile interface {
	MenuFileCreate(ctx context.Context, data menu_file.CreateMenuFileRequest) (menu_file.CreateMenuFileResponse, *pkg.Error)
	MenuFileGetById(ctx context.Context, id string) (menu_file.GetMenuFileResponse, *pkg.Error)
	MenuFileGetAll(ctx context.Context, filter menu_file.Filter) ([]menu_file.GetMenuFileListResponse, int, *pkg.Error)
	MenuFileUpdate(ctx context.Context, data menu_file.UpdateMenuFileRequest) *pkg.Error
	MenuFileDelete(ctx context.Context, id string) *pkg.Error
}
