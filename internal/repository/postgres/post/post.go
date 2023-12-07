package post

import (
	"context"
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

	if request.AuthorId != nil {
		detail.AuthorId = *request.AuthorId
	} else {
		detail.AuthorId = dataCtx.UserId
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
	detail.CreatedAt = time.Now()

	_, err := r.NewInsert().Model(&detail).Exec(ctx)

	if err != nil {
		return CreatePostResponse{}, &pkg.Error{
			Err:    pkg.WrapError(err, "creating post"),
			Status: http.StatusInternalServerError,
		}
	}

	return detail, nil
}

/*
func (r Repository) PostsGetById(ctx context.Context, id string) (GetPostsResponse, *pkg.Error) {
	var post GetPostsResponse

	err := r.NewSelect().Model(&post).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return GetPostsResponse{}, &pkg.Error{
			Err:    pkg.WrapError(err, "selecting post get by id"),
			Status: http.StatusInternalServerError,
		}
	}

	return post, nil
}

func (r Repository) PostsGetAll(ctx context.Context, filter Filter) ([]GetPostsListResponse, int, *pkg.Error) {
	var list []GetPostsListResponse

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

	if filter.Postsname != nil {
		q.WhereGroup(" and ", func(query *bun.SelectQuery) *bun.SelectQuery {
			query.Where("lower(postname) like lower(?)", "%"+*filter.Postsname+"%")
			return query
		})
	}

	q.Order("created_at desc")

	count, err := q.ScanAndCount(ctx)
	if err != nil {
		return nil, 0, &pkg.Error{
			Err:    pkg.WrapError(err, "selecting post list"),
			Status: http.StatusInternalServerError,
		}
	}
	return list, count, nil
}

func (r Repository) PostsUpdate(ctx context.Context, request UpdatePostsRequest) *pkg.Error {
	var detail entity.Posts

	err := r.NewSelect().Model(&detail).Where("id = ?", &request.Id).Scan(ctx)
	if err != nil {
		return &pkg.Error{
			Err:    pkg.WrapError(err, "updating post get by id"),
			Status: http.StatusInternalServerError,
		}
	}

	if request.Postsname != nil {
		detail.Postsname = *request.Postsname
	}
	if request.AvatarLink != nil {
		detail.Avatar = *request.AvatarLink
	}
	if request.FullName != nil {
		detail.FullName = *request.FullName
	}
	if request.Gender != nil {
		detail.Gender = *request.Gender
	}

	if request.Status != nil {
		detail.Status = *request.Status
	}

	if request.Role != nil {
		detail.Role = *request.Role
	}
	if request.BirthDate != nil {
		detail.BirthDate = *request.BirthDate
	}
	if request.Phone != nil {
		detail.Phone = *request.Phone
	}

	date := time.Now()
	detail.UpdatedAt = &date
	detail.UpdatedBy = request.UpdatedBy

	_, err = r.NewUpdate().Model(&detail).Where("id = ?", detail.Id).Exec(ctx)

	if err != nil {
		return &pkg.Error{
			Err:    pkg.WrapError(err, "updating post"),
			Status: http.StatusInternalServerError,
		}
	}
	return nil
}

func (r Repository) PostsDelete(ctx context.Context, req DeletePostsRequest) *pkg.Error {

	_, err := r.NewUpdate().
		Table("posts").
		Where("deleted_at is null AND id = ?", req.Id).
		Set("deleted_at = ?, deleted_by = ?", time.Now(), req.DeletedBy).
		Exec(ctx)

	if err != nil {
		return &pkg.Error{
			Err:    errors.New("delete row error, updating"),
			Status: http.StatusInternalServerError,
		}
	}

	return nil
}

func (r Repository) PostsUpdatePassword(ctx context.Context, req UpdatePasswordRequest) *pkg.Error {

	if req.NewPassword != nil {
		password, err := hash.HashPassword(*req.NewPassword)
		if err != nil {
			return &pkg.Error{
				Err:    pkg.WrapError(err, "creating post hash password"),
				Status: http.StatusInternalServerError,
			}
		}
		req.NewPassword = &password
	}

	_, err := r.NewUpdate().
		Table("posts").
		Where("deleted_at is null AND id = ?", req.Id).
		Set("password = ?, updated_at = ?, updated_by = ?", req.NewPassword, time.Now(), req.UpdatedBy).
		Exec(ctx)

	if err != nil {
		return &pkg.Error{
			Err:    errors.New("reset password row error, updating"),
			Status: http.StatusInternalServerError,
		}
	}

	return nil
}

func (r Repository) GetPostsByPostsname(ctx context.Context, postname string) (DetailPostsResponse, *pkg.Error) {
	var detail DetailPostsResponse

	err := r.NewSelect().Model(&detail).Where("postname = ?", postname).Scan(ctx)
	if err != nil {
		return DetailPostsResponse{}, &pkg.Error{
			Err:    pkg.WrapError(err, "selecting post get by postname"),
			Status: http.StatusInternalServerError,
		}
	}
	return detail, nil
}
*/
