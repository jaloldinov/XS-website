package post_hashtag

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
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

func (r Repository) PostHashtagCreate(ctx context.Context, request CreatePostHashtagRequest) (CreatePostHashtagResponse, *pkg.Error) {
	var detail CreatePostHashtagResponse
	dataCtx, er := r.CheckCtx(ctx)
	if er != nil {
		return CreatePostHashtagResponse{}, er
	}

	detail.Id = uuid.NewString()
	detail.PostId = request.PostId
	detail.HashtagId = request.HashtagId

	timeNow := time.Now()
	detail.CreatedAt = timeNow
	detail.CreatedBy = dataCtx.UserId

	_, err := r.NewInsert().Model(&detail).Exec(ctx)

	if err != nil {
		return CreatePostHashtagResponse{}, &pkg.Error{
			Err:    pkg.WrapError(err, "creating hashtag"),
			Status: http.StatusInternalServerError,
		}
	}

	return detail, nil
}

func (r Repository) PostHashtagGetAll(ctx context.Context, post_id string) ([]GetPostHashtagListResponse, int, *pkg.Error) {

	var list []GetPostHashtagListResponse

	query := fmt.Sprintf(`
	SELECT
		post_hashtag.id,
		hashtags.name
	FROM
		hashtags
	JOIN
		post_hashtag ON hashtags.id = post_hashtag.hashtag_id
	WHERE
		post_hashtag.post_id = '%s'
		AND post_hashtag.deleted_at IS NULL
		AND hashtags.deleted_at IS NULL
`, post_id)

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, &pkg.Error{
			Err:    pkg.WrapError(err, "selecting population request list"),
			Status: http.StatusInternalServerError,
		}
	}

	for rows.Next() {
		var detail GetPostHashtagListResponse

		if err = rows.Scan(
			&detail.Id,
			&detail.Name,
		); err != nil {
			return nil, 0, &pkg.Error{
				Err:    pkg.WrapError(err, "selecting population request list"),
				Status: http.StatusInternalServerError,
			}
		}

		list = append(list, detail)
	}

	var count int

	countQuery := fmt.Sprintf(
		`SELECT
				count(hashtags.id)
			FROM hashtags 
			JOIN post_hashtag ON hashtags.id = post_hashtag.hashtag_id
			WHERE post_hashtag.post_id = '%s' AND post_hashtag.deleted_at IS NULL AND hashtags.deleted_at IS NULL`, post_id)

	err = r.QueryRowContext(ctx, countQuery).Scan(&count)
	if err != nil {
		return nil, 0, &pkg.Error{
			Err:    pkg.WrapError(err, "selecting population request list"),
			Status: http.StatusInternalServerError,
		}
	}
	return list, count, nil
}

func (r Repository) PostHashtagDelete(ctx context.Context, id string) *pkg.Error {
	dataCtx, er := r.CheckCtx(ctx)
	if er != nil {
		return er
	}

	result, err := r.NewUpdate().
		Table("post_hashtag").
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
