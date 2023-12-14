package hashtag

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

func (r Repository) HashtagCreate(ctx context.Context, request CreateHashtagRequest) (CreateHashtagResponse, *pkg.Error) {
	var detail CreateHashtagResponse
	dataCtx, er := r.CheckCtx(ctx)
	if er != nil {
		return CreateHashtagResponse{}, er
	}

	detail.Id = uuid.NewString()
	detail.Name = request.Name
	detail.Status = request.Status

	if request.Status != nil {
		detail.Status = request.Status
	}
	timeNow := time.Now()
	detail.CreatedAt = timeNow
	detail.CreatedBy = dataCtx.UserId

	_, err := r.NewInsert().Model(&detail).Exec(ctx)

	if err != nil {
		return CreateHashtagResponse{}, &pkg.Error{
			Err:    pkg.WrapError(err, "creating hashtag"),
			Status: http.StatusInternalServerError,
		}
	}

	return detail, nil
}

func (r Repository) HashtagGetById(ctx context.Context, id string) (GetHashtagResponse, *pkg.Error) {
	var hashtag GetHashtagResponse

	err := r.NewSelect().Model(&hashtag).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return GetHashtagResponse{}, &pkg.Error{
			Err:    pkg.WrapError(err, "selecting hashtag get by id"),
			Status: http.StatusInternalServerError,
		}
	}

	return hashtag, nil
}

func (r Repository) HashtagGetAll(ctx context.Context, filter Filter) ([]GetHashtagListResponse, int, *pkg.Error) {

	var list []GetHashtagListResponse

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

	q.Order("created_at desc")

	count, err := q.ScanAndCount(ctx)
	if err != nil {
		return nil, 0, &pkg.Error{
			Err:    pkg.WrapError(err, "selecting hashtag list"),
			Status: http.StatusInternalServerError,
		}
	}
	return list, count, nil
}

func (r Repository) HashtagUpdate(ctx context.Context, request UpdateHashtagRequest) *pkg.Error {
	var detail entity.Hashtag
	dataCtx, er := r.CheckCtx(ctx)
	if er != nil {
		return er
	}

	err := r.NewSelect().Model(&detail).Where("id = ?", &request.Id).Scan(ctx)
	if err != nil {
		return &pkg.Error{
			Err:    pkg.WrapError(err, "updating hashtag get by id"),
			Status: http.StatusInternalServerError,
		}
	}

	if request.Name != nil {
		detail.Name = request.Name
	}
	if request.Status != nil {
		detail.Status = request.Status
	}

	detail.UpdatedBy = &dataCtx.UserId
	date := time.Now()
	detail.UpdatedAt = &date

	_, err = r.NewUpdate().Model(&detail).Where("id = ?", detail.Id).Exec(ctx)

	if err != nil {
		return &pkg.Error{
			Err:    pkg.WrapError(err, "updating hashtag"),
			Status: http.StatusInternalServerError,
		}
	}
	return nil
}

func (r Repository) HashtagDelete(ctx context.Context, id string) *pkg.Error {
	dataCtx, er := r.CheckCtx(ctx)
	if er != nil {
		return er
	}

	result, err := r.NewUpdate().
		Table("hashtags").
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
