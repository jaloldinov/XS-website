package menu

import (
	"context"
	"xs/internal/pkg"
	"xs/internal/repository/postgres/menu"
)

type Menu interface {
	MenuCreate(ctx context.Context, data menu.CreateMenuRequest) (menu.CreateMenuResponse, *pkg.Error)
	// MenuGetById(ctx context.Context, id string) (menu.GetMenuResponse, *pkg.Error)
	// MenuGetAll(ctx context.Context, filter menu.Filter) ([]menu.GetMenuListResponse, int, *pkg.Error)
	// MenuUpdate(ctx context.Context, data menu.UpdateMenuRequest) *pkg.Error
	// MenuDelete(ctx context.Context, id string) *pkg.Error
}
