package user

import (
	"context"
	"errors"
	"net/http"
	"time"
	"xs/internal/entity"
	"xs/internal/pkg"
	"xs/internal/pkg/repository/postgres"
	"xs/internal/service/hash"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Repository struct {
	*postgres.Database
}

func NewRepository(postgresDB *postgres.Database) *Repository {
	return &Repository{postgresDB}
}

func (r Repository) UserCreate(ctx context.Context, request CreateUserRequest) (CreateUserResponse, *pkg.Error) {
	var detail CreateUserResponse
	dataCtx, er := r.CheckCtx(ctx)
	if er != nil {
		return CreateUserResponse{}, er
	}

	detail.Id = uuid.NewString()
	detail.Avatar = request.AvatarLink
	detail.Username = *request.Username
	detail.FullName = request.FullName
	detail.Role = *request.Role
	detail.BirthDate = request.BirthDate
	detail.Gender = request.Gender
	detail.Phone = request.Phone
	detail.CreatedBy = &dataCtx.UserId

	if request.Password != nil {
		password, err := hash.HashPassword(*request.Password)
		if err != nil {
			return CreateUserResponse{}, &pkg.Error{
				Err:    pkg.WrapError(err, "creating user hash password"),
				Status: http.StatusInternalServerError,
			}
		}
		detail.Password = &password
	}
	detail.Status = true
	detail.CreatedAt = time.Now()

	_, err := r.NewInsert().Model(&detail).Exec(ctx)

	if err != nil {
		return CreateUserResponse{}, &pkg.Error{
			Err:    pkg.WrapError(err, "creating user"),
			Status: http.StatusInternalServerError,
		}
	}

	detail.Password = nil
	return detail, nil
}

func (r Repository) UserGetById(ctx context.Context, id string) (GetUserResponse, *pkg.Error) {
	var user GetUserResponse

	err := r.NewSelect().Model(&user).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return GetUserResponse{}, &pkg.Error{
			Err:    pkg.WrapError(err, "selecting user get by id"),
			Status: http.StatusInternalServerError,
		}
	}

	return user, nil
}

func (r Repository) UserGetAll(ctx context.Context, filter Filter) ([]GetUserListResponse, int, *pkg.Error) {
	var list []GetUserListResponse

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

	if filter.Username != nil {
		q.WhereGroup(" and ", func(query *bun.SelectQuery) *bun.SelectQuery {
			query.Where("lower(username) like lower(?)", "%"+*filter.Username+"%")
			return query
		})
	}

	q.Order("created_at desc")

	count, err := q.ScanAndCount(ctx)
	if err != nil {
		return nil, 0, &pkg.Error{
			Err:    pkg.WrapError(err, "selecting user list"),
			Status: http.StatusInternalServerError,
		}
	}
	return list, count, nil
}

func (r Repository) UserUpdate(ctx context.Context, request UpdateUserRequest) *pkg.Error {
	var detail entity.User
	dataCtx, er := r.CheckCtx(ctx)
	if er != nil {
		return er
	}

	err := r.NewSelect().Model(&detail).Where("id = ?", &request.Id).Scan(ctx)
	if err != nil {
		return &pkg.Error{
			Err:    pkg.WrapError(err, "updating user get by id"),
			Status: http.StatusInternalServerError,
		}
	}

	if request.Username != nil {
		detail.Username = *request.Username
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
	detail.UpdatedBy = &dataCtx.UserId

	_, err = r.NewUpdate().Model(&detail).Where("id = ?", detail.Id).Exec(ctx)

	if err != nil {
		return &pkg.Error{
			Err:    pkg.WrapError(err, "updating user"),
			Status: http.StatusInternalServerError,
		}
	}
	return nil
}

func (r Repository) UserDelete(ctx context.Context, id string) *pkg.Error {

	dataCtx, er := r.CheckCtx(ctx)
	if er != nil {
		return er
	}
	_, err := r.NewUpdate().
		Table("users").
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

func (r Repository) UserUpdatePassword(ctx context.Context, req UpdatePasswordRequest) *pkg.Error {

	dataCtx, er := r.CheckCtx(ctx)
	if er != nil {
		return er
	}
	if req.NewPassword != nil {
		password, err := hash.HashPassword(*req.NewPassword)
		if err != nil {
			return &pkg.Error{
				Err:    pkg.WrapError(err, "creating user hash password"),
				Status: http.StatusInternalServerError,
			}
		}
		req.NewPassword = &password
	}

	_, err := r.NewUpdate().
		Table("users").
		Where("deleted_at is null AND id = ?", req.Id).
		Set("password = ?, updated_at = ?, updated_by = ?", req.NewPassword, time.Now(), dataCtx.UserId).
		Exec(ctx)

	if err != nil {
		return &pkg.Error{
			Err:    errors.New("reset password row error, updating"),
			Status: http.StatusInternalServerError,
		}
	}

	return nil
}

func (r Repository) GetUserByUsername(ctx context.Context, username string) (DetailUserResponse, *pkg.Error) {
	var detail DetailUserResponse

	err := r.NewSelect().Model(&detail).Where("username = ?", username).Scan(ctx)
	if err != nil {
		return DetailUserResponse{}, &pkg.Error{
			Err:    pkg.WrapError(err, "selecting user get by username"),
			Status: http.StatusInternalServerError,
		}
	}
	return detail, nil
}
