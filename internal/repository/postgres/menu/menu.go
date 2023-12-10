package menu

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
	"xs/internal/entity"
	"xs/internal/pkg"
	"xs/internal/pkg/repository/postgres"

	"github.com/google/uuid"
)

type Repository struct {
	*postgres.Database
}

func NewRepository(postgresDB *postgres.Database) *Repository {
	return &Repository{postgresDB}
}

func (r Repository) MenuCreate(ctx context.Context, request CreateMenuRequest) (CreateMenuResponse, *pkg.Error) {
	var detail CreateMenuResponse
	dataCtx, er := r.CheckCtx(ctx)
	if er != nil {
		return CreateMenuResponse{}, er
	}

	detail.Id = uuid.NewString()
	detail.Title = request.Title
	detail.Content = request.Content
	if request.IsStatic != nil {
		detail.IsStatic = request.IsStatic
	}
	if request.Status != nil {
		detail.Status = request.Status
	}
	if request.ParentId != nil {
		detail.ParentId = request.ParentId
	}

	detail.Slug = request.Slug
	detail.Type = &request.Type
	detail.CreatedBy = dataCtx.UserId

	timeNow := time.Now()
	detail.CreatedAt = &timeNow

	_, err := r.NewInsert().Model(&detail).Exec(ctx)

	if err != nil {
		return CreateMenuResponse{}, &pkg.Error{
			Err:    pkg.WrapError(err, "creating menu"),
			Status: http.StatusInternalServerError,
		}
	}

	return detail, nil
}

func (r Repository) MenuGetAll(ctx context.Context, filter Filter) ([]GetMenuListResponse, int, *pkg.Error) {
	// var list []GetMenuListResponse
	dataCtx, er := r.CheckCtx(ctx)
	if er != nil {
		return nil, 0, er
	}
	filter.Lang = &dataCtx.Lang

	query := fmt.Sprintf(
		`SELECT
			id,
			title,
			content,
			status,
			slug,
			parent_id
		FROM menu WHERE deleted_at IS NULL`)
	where := ""

	query += where

	if filter.Lang != nil {
		where += fmt.Sprintf(" AND title->>'%s' is not null", *filter.Lang)
	}

	query += where

	if filter.Offset != nil {
		query += fmt.Sprintf(" OFFSET  %d", *filter.Offset)
	}

	if filter.Limit != nil {
		query += fmt.Sprintf(" LIMIT %d", *filter.Limit)
	}

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, &pkg.Error{
			Err:    pkg.WrapError(err, "selecting population request list"),
			Status: http.StatusInternalServerError,
		}
	}

	menuMap := make(map[string]*GetMenuListResponse)

	for rows.Next() {
		var detail GetMenuListResponse
		var titleB, contentB []byte
		if err = rows.Scan(
			&detail.Id,
			&titleB,
			&contentB,
			&detail.Status,
			&detail.Slug,
			&detail.ParentId,
		); err != nil {
			return nil, 0, &pkg.Error{
				Err:    pkg.WrapError(err, "selecting population request list"),
				Status: http.StatusInternalServerError,
			}
		}

		var title, content map[string]string

		err = json.Unmarshal(titleB, &title)
		if err != nil {
			return nil, 0, &pkg.Error{
				Err:    pkg.WrapError(err, "selecting population request list"),
				Status: http.StatusInternalServerError,
			}
		}

		detail.Title = title[*filter.Lang]
		detail.Content = content[*filter.Lang]
		err = json.Unmarshal(contentB, &content)
		if err != nil {
			return nil, 0, &pkg.Error{
				Err:    pkg.WrapError(err, "selecting population request list"),
				Status: http.StatusInternalServerError,
			}
		}

		for k, v := range content {
			if content[k] != "" {
				detail.ContentLanguages = append(detail.ContentLanguages, k)
				detail.Content = v
			}
		}
		for k, v := range title {
			if title[k] != "" {
				detail.TitleLanguages = append(detail.TitleLanguages, k)
				detail.Title = v
			}
		}

		menuMap[detail.Id] = &detail
	}

	// Build the children hierarchy
	for _, menu := range menuMap {
		parentID := menu.ParentId
		if parentID != nil {
			parent, exists := menuMap[*parentID]
			if exists {
				if parent.Children == nil {
					children := []GetMenuListResponse{*menu}
					parent.Children = &children
				} else {
					*parent.Children = append(*parent.Children, *menu)
				}
			}
		}
	}

	// Find the root-level menus
	var rootMenus []GetMenuListResponse
	for _, menu := range menuMap {
		if menu.ParentId == nil {
			rootMenus = append(rootMenus, *menu)
		}
	}

	var count int

	countQuery := `
		SELECT
		    count(id)
		FROM menu WHERE deleted_at IS NULL`

	err = r.QueryRowContext(ctx, countQuery+where).Scan(&count)
	if err != nil {
		return nil, 0, &pkg.Error{
			Err:    pkg.WrapError(err, "selecting population request list"),
			Status: http.StatusInternalServerError,
		}
	}

	return rootMenus, count, nil
}

func (r Repository) MenuUpdate(ctx context.Context, request UpdateMenuRequest) *pkg.Error {
	var detail entity.Menu
	dataCtx, er := r.CheckCtx(ctx)
	if er != nil {
		return er
	}

	err := r.NewSelect().Model(&detail).Where("id = ?", &request.Id).Scan(ctx)
	if err != nil {
		return &pkg.Error{
			Err:    pkg.WrapError(err, "updating menu get by id"),
			Status: http.StatusInternalServerError,
		}
	}

	if request.Title != nil {
		detail.Title = request.Title
	}
	if request.Content != nil {
		detail.Content = *request.Content
	}
	if request.ParentId != nil {
		detail.ParentId = request.ParentId
	}
	if request.IsStatic != nil {
		detail.IsStatic = request.IsStatic
	}
	if request.Status != nil {
		detail.Status = request.Status
	}
	if request.Slug != nil {
		detail.Slug = request.Slug
	}

	if request.Type != nil {
		detail.Type = request.Type
	}

	detail.UpdatedBy = &dataCtx.UserId
	date := time.Now()
	detail.UpdatedAt = &date

	_, err = r.NewUpdate().Model(&detail).Where("id = ?", detail.Id).Exec(ctx)

	if err != nil {
		return &pkg.Error{
			Err:    pkg.WrapError(err, "updating menu"),
			Status: http.StatusInternalServerError,
		}
	}
	return nil
}

func (r Repository) MenuDelete(ctx context.Context, id string) *pkg.Error {
	dataCtx, er := r.CheckCtx(ctx)
	if er != nil {
		return er
	}

	_, err := r.NewUpdate().
		Table("menu").
		Where("deleted_at is null AND id = ?", id).
		Set("deleted_at = ?, deleted_by = ?", time.Now(), dataCtx.UserId).
		Exec(ctx)

	if err != nil {
		return &pkg.Error{
			Err:    errors.New("delete row error, updating"),
			Status: http.StatusInternalServerError,
		}
	}

	return nil
}
