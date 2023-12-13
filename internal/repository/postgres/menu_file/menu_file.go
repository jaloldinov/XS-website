package menu_file

import (
	"context"
	"errors"
	"net/http"
	"time"
	"xs/internal/entity"
	"xs/internal/pkg"
	"xs/internal/pkg/repository/postgres"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Repository struct {
	*postgres.Database
}

func NewRepository(postgresDB *postgres.Database) *Repository {
	return &Repository{postgresDB}
}

func (r Repository) MenuFileCreate(ctx context.Context, request CreateMenuFileRequest) (CreateMenuFileResponse, *pkg.Error) {
	var detail CreateMenuFileResponse

	dataCtx, er := r.CheckCtx(ctx)
	if er != nil {
		return CreateMenuFileResponse{}, er
	}

	detail.Id = uuid.NewString()
	detail.Link = *request.FileLink
	detail.Type = request.Type
	detail.MenuId = *request.MenuId

	if request.MarkedLink != nil {
		detail.MarkedLink = *request.MarkedLink
	}
	if request.Grouping != nil {
		detail.Grouping = *request.Grouping
	}
	if request.Carusel != nil {
		detail.Carusel = *request.Carusel
	}
	if request.AuthorId != nil {
		detail.AuthorId = *request.AuthorId
	} else {
		detail.AuthorId = dataCtx.UserId
	}

	detail.CreatedBy = dataCtx.UserId

	_, err := r.NewInsert().Model(&detail).Exec(ctx)

	if err != nil {
		return CreateMenuFileResponse{}, &pkg.Error{
			Err:    pkg.WrapError(err, "creating menu file"),
			Status: http.StatusInternalServerError,
		}
	}

	return detail, nil
}

func (r Repository) MenuFileGetById(ctx context.Context, id string) (GetMenuFileResponse, *pkg.Error) {
	var menu_file GetMenuFileResponse

	err := r.NewSelect().Model(&menu_file).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return GetMenuFileResponse{}, &pkg.Error{
			Err:    pkg.WrapError(err, "selecting menu_file get by id"),
			Status: http.StatusInternalServerError,
		}
	}

	return menu_file, nil
}

func (r Repository) MenuFileGetAll(ctx context.Context, filter Filter) ([]GetMenuFileListResponse, int, *pkg.Error) {
	var list []GetMenuFileListResponse

	q := r.NewSelect().Model(&list)
	q.WhereGroup(" and ", func(query *bun.SelectQuery) *bun.SelectQuery {
		query.Where("deleted_at is null")
		return query
	})
	if filter.Limit != nil {
		q.Limit(*filter.Limit)
	}

	if filter.Offset != nil {
		q.Offset(*filter.Offset)
	}

	if filter.MenuId != nil {
		q.WhereGroup(" and ", func(query *bun.SelectQuery) *bun.SelectQuery {
			query.Where("menu_id = ?", *filter.MenuId)
			return query
		})
	}
	if filter.Type != nil {
		q.WhereGroup(" and ", func(query *bun.SelectQuery) *bun.SelectQuery {
			query.Where("type = ?", *filter.Type)
			return query
		})
	}

	q.Order("created_at desc")

	count, err := q.ScanAndCount(ctx)
	if err != nil {
		return nil, 0, &pkg.Error{
			Err:    pkg.WrapError(err, "selecting menu file list"),
			Status: http.StatusInternalServerError,
		}
	}
	return list, count, nil
}

func (r Repository) MenuFileUpdate(ctx context.Context, request UpdateMenuFileRequest) *pkg.Error {
	var detail entity.MenuFile
	dataCtx, er := r.CheckCtx(ctx)
	if er != nil {
		return er
	}

	err := r.NewSelect().Model(&detail).Where("id = ?", &request.Id).Scan(ctx)
	if err != nil {
		return &pkg.Error{
			Err:    pkg.WrapError(err, "updating menu_file get by id"),
			Status: http.StatusInternalServerError,
		}
	}

	if request.FileLink != nil && *request.FileLink != "" {
		detail.Link = *request.FileLink
	}

	detail.Type = request.Type

	if request.MarkedLink != nil {
		detail.MarkedLink = *request.MarkedLink
	}
	if request.Grouping != nil {
		detail.Grouping = *request.Grouping
	}
	if request.Carusel != nil {
		detail.Carusel = *request.Carusel
	}
	if request.MenuId != nil {
		detail.MenuId = *request.MenuId
	}
	if request.AuthorId != nil {
		detail.AuthorId = *request.AuthorId
	}

	date := time.Now()
	detail.UpdatedAt = &date
	detail.UpdatedBy = &dataCtx.UserId

	_, err = r.NewUpdate().Model(&detail).Where("id = ?", detail.Id).Exec(ctx)

	if err != nil {
		return &pkg.Error{
			Err:    pkg.WrapError(err, "updating menu_file"),
			Status: http.StatusInternalServerError,
		}
	}
	return nil
}

func (r Repository) MenuFileDelete(ctx context.Context, id string) *pkg.Error {
	dataCtx, er := r.CheckCtx(ctx)
	if er != nil {
		return er
	}

	result, err := r.NewUpdate().
		Table("menu_file").
		Where("deleted_at is null AND id = ?", id).
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

	return nil
}
