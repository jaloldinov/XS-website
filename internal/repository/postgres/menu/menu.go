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
	detail.Index = request.Index
	detail.CreatedBy = dataCtx.UserId

	timeNow := time.Now()
	detail.CreatedAt = &timeNow

	// check requesting data is whether child or parent
	if detail.ParentId != nil && *detail.Type == "EXTRA" {

		// Get the maximum index of existing menu items with the same parent
		rows, err := r.QueryContext(ctx, "SELECT COALESCE(MAX(index), 0) FROM menu WHERE parent_id = ? AND deleted_at IS NULL", detail.ParentId)
		if err != nil {
			return CreateMenuResponse{}, &pkg.Error{
				Err:    pkg.WrapError(err, "cannot query max index value"),
				Status: http.StatusInternalServerError,
			}
		}
		var maxPosition int
		if rows.Next() {
			if err := rows.Scan(&maxPosition); err != nil {
				return CreateMenuResponse{}, &pkg.Error{
					Err:    pkg.WrapError(err, "cannot scan max index value"),
					Status: http.StatusInternalServerError,
				}
			}
		}

		if detail.Index == 1 {
			// Update all existing indexes by incrementing them by 1
			if _, err = r.NewUpdate().Table("menu").
				Where("parent_id = ? AND index > 0 AND deleted_at IS NULL", detail.ParentId).
				Set("index = index + 1").
				Exec(ctx); err != nil {
				return CreateMenuResponse{}, &pkg.Error{
					Err:    pkg.WrapError(err, "cannot update index by incrementing by 1 [ detail.Index == 1 ]"),
					Status: http.StatusInternalServerError,
				}
			}

		} else if detail.Index <= maxPosition {
			// if the given index is lower than the max index change all the menu items from given index
			if _, err = r.NewUpdate().
				Table("menu").
				Where("(index >= ?) and (parent_id = ?) AND deleted_at IS NULL", detail.Index, *detail.ParentId).
				Set("index = index + 1").
				Exec(ctx); err != nil {
				return CreateMenuResponse{}, &pkg.Error{
					Err:    pkg.WrapError(err, "updating indexes of menu items with indexes [ detail.Index <= maxPosition ]"),
					Status: http.StatusInternalServerError,
				}
			}

		} else if detail.Index >= maxPosition {
			// If the provided index is greater than the maximum index
			// Update new menu index = maxPosition + 1
			detail.Index = maxPosition + 1
		}

	} else if detail.ParentId == nil && *detail.Type == "MAIN" {

		// Get the maximum index of the existing only parent menu items
		rows, err := r.QueryContext(ctx, "SELECT COALESCE(MAX(index), 0) FROM menu WHERE parent_id IS NULL AND deleted_at IS NULL")
		if err != nil {
			return CreateMenuResponse{}, &pkg.Error{
				Err:    pkg.WrapError(err, "cannot query max index value"),
				Status: http.StatusInternalServerError,
			}
		}
		var maxPosition int
		if rows.Next() {
			if err := rows.Scan(&maxPosition); err != nil {
				return CreateMenuResponse{}, &pkg.Error{
					Err:    pkg.WrapError(err, "cannot scan max index value"),
					Status: http.StatusInternalServerError,
				}
			}
		}

		if detail.Index == 1 {
			// Update all existing indexes by incrementing them by 1
			if _, err = r.NewUpdate().Table("menu").
				Where("parent_id IS NULL AND index > 0 AND deleted_at IS NULL").
				Set("index = index + 1").
				Exec(ctx); err != nil {
				return CreateMenuResponse{}, &pkg.Error{
					Err:    pkg.WrapError(err, "cannot update index by incrementing by 1 [ detail.Index == 1 ]"),
					Status: http.StatusInternalServerError,
				}
			}

		} else if detail.Index <= maxPosition {
			// if the given index is lower than the max index change all the menu items from given index
			if _, err = r.NewUpdate().
				Table("menu").
				Where("(index >= ?) and parent_id is null AND deleted_at IS NULL", detail.Index).
				Set("index = index + 1").
				Exec(ctx); err != nil {
				return CreateMenuResponse{}, &pkg.Error{
					Err:    pkg.WrapError(err, "updating indexes of menu items with indexes [ detail.Index <= maxPosition ]"),
					Status: http.StatusInternalServerError,
				}
			}

		} else if detail.Index >= maxPosition {
			// If the provided index is greater than the maximum index
			// Update new menu index = maxPosition + 1
			detail.Index = maxPosition + 1
		}

	}

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
			parent_id,
			index
		FROM menu WHERE deleted_at IS NULL `)
	where := ""

	query += where

	if filter.Lang != nil {
		where += fmt.Sprintf(" AND title->>'%s' is not null", *filter.Lang)
	}

	query += where
	query += fmt.Sprintf(` ORDER BY index`)
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
	// menuMap := []*GetMenuListResponse{}

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
			&detail.Index,
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

	// Find the root-level menu
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

	// Get the menu item to be deleted
	var data UpdateMenuIndex
	err := r.NewSelect().Model(&data).Where("id = ? AND deleted_at IS NULL", id).Scan(ctx)
	if err != nil {
		return &pkg.Error{
			Err:    pkg.WrapError(err, "selecting menu by id"),
			Status: http.StatusInternalServerError,
		}
	}

	// Store the current index of the menu item
	oldIndex := data.Index

	// Delete the menu item
	result, err := r.NewUpdate().
		Table("menu").
		Where("deleted_at IS NULL AND id = ?", id).
		Set("deleted_at = ?, deleted_by = ?", time.Now(), dataCtx.UserId).
		Exec(ctx)

	if err != nil {
		return &pkg.Error{
			Err:    errors.New("delete row error, updating"),
			Status: http.StatusInternalServerError,
		}
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return &pkg.Error{
			Err:    errors.New("no matching ID found"),
			Status: http.StatusNotFound,
		}
	}

	// Update the indexes of other menu items if the deleted item had a parent
	if data.ParentId != nil {
		_, err := r.NewUpdate().
			Table("menu").
			Where("(index > ?) AND deleted_at IS NULL AND parent_id = ?", oldIndex, data.ParentId).
			Set("index = index - 1").
			Exec(ctx)
		if err != nil {
			return &pkg.Error{
				Err:    pkg.WrapError(err, "updating menu indexes"),
				Status: http.StatusInternalServerError,
			}
		}
	} else if data.ParentId == nil {
		_, err := r.NewUpdate().
			Table("menu").
			Where("(index > ?) AND deleted_at IS NULL AND parent_id is null", oldIndex).
			Set("index = index - 1").Exec(ctx)
		if err != nil {
			return &pkg.Error{
				Err:    pkg.WrapError(err, "updating menu indexes"),
				Status: http.StatusInternalServerError,
			}
		}
	}

	return nil
}

func (r Repository) UpdateIndex(ctx context.Context, request UpdateMenuIndex) *pkg.Error {
	var data UpdateMenuIndex

	err := r.NewSelect().Model(&data).Where("id = ? AND deleted_at IS NULL", request.Id).Scan(ctx)
	if err != nil {
		return &pkg.Error{
			Err:    pkg.WrapError(err, "selecting menu by id"),
			Status: http.StatusInternalServerError,
		}
	}

	// check requesting data is whether child or parent
	oldIndex := data.Index

	// Update the index based on the provided request
	if oldIndex != request.Index {

		if data.ParentId != nil {
			if oldIndex < request.Index {
				// Decrease the index of menu items between oldIndex and request.Index
				_, err = r.NewUpdate().
					Table("menu").
					Where("(index > ?) AND (index <= ?) AND deleted_at IS NULL AND parent_id = ?", oldIndex, request.Index, data.ParentId).
					Set("index = index - 1").
					Exec(ctx)
			} else {
				// Increase the index of menu items between request.Index and oldIndex
				_, err = r.NewUpdate().
					Table("menu").
					Where("(index >= ?) AND (index < ?) AND deleted_at IS NULL AND parent_id  = ?", request.Index, oldIndex, data.ParentId).
					Set("index = index + 1").
					Exec(ctx)
			}

			if err != nil {
				return &pkg.Error{
					Err:    pkg.WrapError(err, "updating menu indexes"),
					Status: http.StatusInternalServerError,
				}
			}

			// Update the index of the target menu item
			data.Index = request.Index
			_, err = r.NewUpdate().Model(&data).Where("id = ?", request.Id).Column("index").Exec(ctx)
			if err != nil {
				return &pkg.Error{
					Err:    pkg.WrapError(err, "updating menu index"),
					Status: http.StatusInternalServerError,
				}
			}

		} else if data.ParentId == nil {
			if oldIndex < request.Index {
				// Decrease the index of menu items between oldIndex and request.Index
				_, err = r.NewUpdate().
					Table("menu").
					Where("(index > ?) AND (index <= ?) AND deleted_at IS NULL AND parent_id IS NULL", oldIndex, request.Index).
					Set("index = index - 1").
					Exec(ctx)
			} else {
				// Increase the index of menu items between request.Index and oldIndex
				_, err = r.NewUpdate().
					Table("menu").
					Where("(index >= ?) AND (index < ?) AND deleted_at IS NULL AND parent_id IS NULL", request.Index, oldIndex).
					Set("index = index + 1").
					Exec(ctx)
			}

			if err != nil {
				return &pkg.Error{
					Err:    pkg.WrapError(err, "updating menu indexes"),
					Status: http.StatusInternalServerError,
				}
			}

			// Update the index of the target menu item
			data.Index = request.Index
			_, err = r.NewUpdate().Model(&data).Where("id = ?", request.Id).Column("index").Exec(ctx)
			if err != nil {
				return &pkg.Error{
					Err:    pkg.WrapError(err, "updating menu index"),
					Status: http.StatusInternalServerError,
				}
			}
		}
	}
	return nil
}
