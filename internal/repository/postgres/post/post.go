package post

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

func (r Repository) PostCreate(ctx context.Context, request CreatePostRequest) (CreatePostResponse, *pkg.Error) {
	var detail CreatePostResponse
	dataCtx, er := r.CheckCtx(ctx)
	if er != nil {
		return CreatePostResponse{}, er
	}
	if request.PubDate == nil {
		now := time.Now().Format("2006-01-02")
		request.PubDate = &now
	}

	detail.Id = uuid.NewString()
	detail.Title = request.Title
	detail.Content = request.Content
	detail.Slug = request.Slug
	detail.AuthorId = request.AuthorId
	detail.MenuId = request.MenuId

	if request.AuthorId != nil {
		detail.AuthorId = request.AuthorId
	} else {
		detail.AuthorId = &dataCtx.UserId
	}

	if request.PubDate != nil {
		publishDate, err := time.Parse("2006-01-02", *request.PubDate)
		if err != nil {

			return CreatePostResponse{}, &pkg.Error{
				Err:    pkg.WrapError(err, "creating post publish date"),
				Status: http.StatusInternalServerError,
			}
		}
		detail.PubDate = &publishDate
	}
	detail.CreatedBy = dataCtx.UserId

	detail.Status = true
	timeNow := time.Now()
	detail.CreatedAt = &timeNow

	_, err := r.NewInsert().Model(&detail).Exec(ctx)

	if err != nil {
		return CreatePostResponse{}, &pkg.Error{
			Err:    pkg.WrapError(err, "creating post"),
			Status: http.StatusInternalServerError,
		}
	}

	return detail, nil
}

func (r Repository) PostGetById(ctx context.Context, id string) (GetPostResponse, *pkg.Error) {
	var post GetPostResponse
	query := fmt.Sprintf(
		`SELECT
				posts.id,
				posts.title,
				posts.content,
				posts.status,
				to_char(posts.pub_date, 'DD.MM.YYYY') as pub_date,
				users.full_name as author_name,
				users.avatar as author_avatar,
				posts.slug,
				posts.menu_id
			FROM posts
			JOIN users ON users.id = posts.author_id
			WHERE posts.deleted_at IS NULL AND posts.id = '%s' `, id)

	row := r.DB.QueryRowContext(ctx, query)
	var titleB, contentB []byte
	err := row.Scan(
		&post.Id,
		&titleB,
		&contentB,
		&post.Status,
		&post.PubDate,
		&post.AuthorName,
		&post.AuthorAvatar,
		&post.Slug,
		&post.MenuId,
	)
	if err != nil {
		return GetPostResponse{}, &pkg.Error{
			Err:    pkg.WrapError(err, "selecting post get by id"),
			Status: http.StatusInternalServerError,
		}
	}
	var title, content map[string]string

	err = json.Unmarshal(titleB, &title)
	if err != nil {
		return GetPostResponse{}, &pkg.Error{
			Err:    pkg.WrapError(err, "json.Unmarshal(titleB, &title)"),
			Status: http.StatusInternalServerError,
		}
	}

	err = json.Unmarshal(contentB, &content)
	if err != nil {
		return GetPostResponse{}, &pkg.Error{
			Err:    pkg.WrapError(err, "Unmarshal(contentB, &content)"),
			Status: http.StatusInternalServerError,
		}
	}

	post.Title = title
	post.Content = content

	return post, nil
}

func (r Repository) PostGetAll(ctx context.Context, filter Filter) ([]GetPostListResponse, int, *pkg.Error) {
	var list []GetPostListResponse
	dataCtx, er := r.CheckCtx(ctx)
	if er != nil {
		return nil, 0, er
	}
	filter.Lang = &dataCtx.Lang

	query := fmt.Sprintf(
		`SELECT
				posts.id,
				posts.title,
				posts.content,
				posts.status,
				to_char(posts.pub_date, 'DD.MM.YYYY') as pub_date,
				users.full_name as author_name,
				users.avatar as author_avatar
			FROM posts
			JOIN users ON users.id = posts.author_id
			WHERE posts.deleted_at IS NULL `)
	where := ""

	query += where

	if filter.Content != nil {
		where += fmt.Sprintf(" AND lower(content->>'%s') similar to lower('%s')", *filter.Lang, "%"+*filter.Content+"%")
	}

	if filter.Title != nil {
		where += fmt.Sprintf(" AND lower(title->>'%s') similar to lower('%s')", *filter.Lang, "%"+*filter.Title+"%")
	}

	if filter.Lang != nil {
		where += fmt.Sprintf(" AND title->>'%s' is not null", *filter.Lang)
	}
	/*
		// if filter.From != nil {
		// 	where += fmt.Sprintf(" AND (created_at > %s0) or (pub_date > %s0)", *filter.From, *filter.From)
		// }

		// if filter.To != nil {
		// 	where += fmt.Sprintf(" AND (created_at < %s0) or (pub_date < %s0)", *filter.To, *filter.To)

		// }

		// if filter.PublishedAt != nil {
		// 	publishedAt, err := time.Parse("02.01.2006", *filter.PublishedAt)
		// 	if err != nil {
		// 		return nil, 0, &pkg.Error{
		// 			Err:    pkg.WrapError(err, "selecting population request list"),
		// 			Status: http.StatusInternalServerError,
		// 		}
		// 	}
		// 	where += fmt.Sprintf(" AND pub_date = %v", publishedAt)
		// }

		// if filter.PublishedFrom != nil {
		// 	publishedFrom, err := time.Parse("02.01.2006", *filter.PublishedFrom)
		// 	if err != nil {
		// 		return nil, 0, &pkg.Error{
		// 			Err:    pkg.WrapError(err, "selecting population request list"),
		// 			Status: http.StatusInternalServerError,
		// 		}
		// 	}
		// 	where += fmt.Sprintf(" AND pub_date >= %v", publishedFrom)
		// }

		// if filter.PublishedTo != nil {
		// 	publishedTo, err := time.Parse("02.01.2006", *filter.PublishedTo)
		// 	if err != nil {
		// 		return nil, 0, &pkg.Error{
		// 			Err:    pkg.WrapError(err, "selecting population request list"),
		// 			Status: http.StatusInternalServerError,
		// 		}
		// 	}
		// 	where += fmt.Sprintf(" AND pub_date <= %v", publishedTo)
		// }
	*/
	query += where

	if filter.Order != nil {
		query += fmt.Sprintf(" ORDER BY posts.created_at %s", *filter.Order)
	} else {
		query += " ORDER BY posts.created_at asc"
	}

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

	for rows.Next() {
		var detail GetPostListResponse
		var titleB, contentB []byte
		if err = rows.Scan(
			&detail.Id,
			&titleB,
			&contentB,
			&detail.Status,
			&detail.PubDate,
			&detail.AuthorName,
			&detail.AuthorAvatar,
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

		err = json.Unmarshal(contentB, &content)
		if err != nil {
			return nil, 0, &pkg.Error{
				Err:    pkg.WrapError(err, "selecting population request list"),
				Status: http.StatusInternalServerError,
			}
		}

		for k := range content {
			if content[k] != "" {
				detail.ContentLanguages = append(detail.ContentLanguages, k)
			}
		}
		for k, v := range title {
			if title[k] != "" {
				detail.TitleLanguages = append(detail.TitleLanguages, k)
				detail.Title = v
			}
		}

		list = append(list, detail)
	}

	var count int

	countQuery := `
		SELECT
		    count(id)
		FROM posts WHERE deleted_at IS NULL`

	err = r.QueryRowContext(ctx, countQuery+where).Scan(&count)
	if err != nil {
		return nil, 0, &pkg.Error{
			Err:    pkg.WrapError(err, "selecting population request list"),
			Status: http.StatusInternalServerError,
		}
	}
	return list, count, nil
}

func (r Repository) PostUpdate(ctx context.Context, request UpdatePostRequest) *pkg.Error {
	var detail entity.Post
	dataCtx, er := r.CheckCtx(ctx)
	if er != nil {
		return er
	}

	err := r.NewSelect().Model(&detail).Where("id = ?", &request.Id).Scan(ctx)
	if err != nil {
		return &pkg.Error{
			Err:    pkg.WrapError(err, "updating post get by id"),
			Status: http.StatusInternalServerError,
		}
	}

	if request.Title != nil {
		detail.Title = request.Title
	}
	if request.Content != nil {
		detail.Content = request.Content
	}
	if request.PubDate != nil {
		publishDate, err := time.Parse("2006-01-02", *request.PubDate)
		if err != nil {

			return &pkg.Error{
				Err:    pkg.WrapError(err, "updating post publish date"),
				Status: http.StatusInternalServerError,
			}
		}
		detail.PubDate = &publishDate
	}
	if request.Status != nil {
		detail.Status = request.Status
	}
	if request.AuthorId != nil {
		detail.AuthorId = *request.AuthorId
	}
	detail.UpdatedBy = &dataCtx.UserId
	date := time.Now()
	detail.UpdatedAt = &date

	_, err = r.NewUpdate().Model(&detail).Where("id = ?", detail.Id).Exec(ctx)

	if err != nil {
		return &pkg.Error{
			Err:    pkg.WrapError(err, "updating post"),
			Status: http.StatusInternalServerError,
		}
	}
	return nil
}

func (r Repository) PostDelete(ctx context.Context, id string) *pkg.Error {
	dataCtx, er := r.CheckCtx(ctx)
	if er != nil {
		return er
	}

	result, err := r.NewUpdate().
		Table("posts").
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

func (r Repository) IsMenuStatic(ctx context.Context, id string) (bool, *pkg.Error) {
	var isStatic bool
	query := fmt.Sprintf(
		`SELECT
				is_static
			FROM menu
			WHERE deleted_at IS NULL AND menu.id = '%s' `, id)

	row := r.DB.QueryRowContext(ctx, query)

	err := row.Scan(
		&isStatic,
	)

	if err != nil {
		return false, &pkg.Error{
			Err:    errors.New("menu not found"),
			Status: http.StatusNotFound,
		}
	}

	return isStatic, nil
}
