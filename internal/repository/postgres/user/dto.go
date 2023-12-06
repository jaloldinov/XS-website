package user

import (
	"mime/multipart"
	"time"

	"github.com/uptrace/bun"
)

type Filter struct {
	Limit    *int
	Offset   *int
	Username *string
}

type CreateUserRequest struct {
	Avatar     *multipart.FileHeader `json:"-" form:"avatar"`
	AvatarLink string                `json:"-" form:"-"`
	Username   *string               `json:"username" form:"username"`
	Password   *string               `json:"password" form:"password"`
	FullName   string                `json:"full_name" form:"full_name"`
	Role       *string               `json:"role" form:"role"`
	BirthDate  string                `json:"birth_date" form:"birth_date"`
	Gender     string                `json:"gender" form:"gender"`
	Phone      string                `json:"phone" form:"phone"`
	CreatedBy  *string               `json:"-"`
}

type CreateUserResponse struct {
	bun.BaseModel `bun:"table:users"`

	Id        string    `json:"id" bun:"id"`
	Avatar    string    `json:"avatar" bun:"avatar"`
	Username  string    `json:"username" bun:"username"`
	Password  *string   `json:"password" bun:"password"`
	FullName  string    `json:"full_name" bun:"full_name"`
	Role      string    `json:"role" bun:"role"`
	BirthDate string    `json:"birth_date" bun:"birth_date"`
	Gender    string    `json:"gender" bun:"gender"`
	Phone     string    `json:"phone" bun:"phone"`
	Status    bool      `json:"status" bun:"status"`
	CreatedBy *string   `json:"created_by"`
	CreatedAt time.Time `json:"-" bun:"created_at"`
}

type GetUserResponse struct {
	bun.BaseModel `bun:"table:users"`

	Id        string `json:"id" bun:"id"`
	Avatar    string `json:"avatar" bun:"avatar"`
	Username  string `json:"username" bun:"username"`
	FullName  string `json:"full_name" bun:"full_name"`
	Status    bool   `json:"status" bun:"status"`
	Role      string `json:"role" bun:"role"`
	BirthDate string `json:"birth_date" bun:"birth_date"`
	Gender    string `json:"gender" bun:"gender"`
	Phone     string `json:"phone" bun:"phone"`
}

type GetUserListResponse struct {
	bun.BaseModel `bun:"table:users"`

	Id       string `json:"id" bun:"id"`
	Avatar   string `json:"avatar" bun:"avatar"`
	Username string `json:"username" bun:"username"`
	FullName string `json:"full_name" bun:"full_name"`
	Status   bool   `json:"status" bun:"status"`
	Role     string `json:"role" bun:"role"`
	Phone    string `json:"phone" bun:"phone"`
}

type UpdateUserRequest struct {
	Id         string                `json:"id" bun:"id"`
	Avatar     *multipart.FileHeader `json:"-" form:"avatar"`
	AvatarLink *string               `json:"-" form:"-"`
	Username   *string               `json:"username" form:"username"`
	FullName   *string               `json:"full_name" form:"full_name"`
	Status     *bool                 `json:"status" form:"status"`
	Role       *string               `json:"role" form:"role"`
	BirthDate  *string               `json:"birth_date" form:"birth_date"`
	Gender     *string               `json:"gender" form:"gender"`
	Phone      *string               `json:"phone" form:"phone"`
	UpdatedBy  *string               `json:"-"`
}

type DeleteUserRequest struct {
	Id        string `json:"id" bun:"id"`
	DeletedBy string `json:"deleted_by" bun:"deleted_by"`
}

type UpdatePasswordRequest struct {
	Id          *string `json:"id" form:"id"`
	NewPassword *string `json:"new_password" form:"new_password"`
	UpdatedBy   *string `json:"-"`
}

type DetailUserResponse struct {
	bun.BaseModel `bun:"table:users"`

	Id       string  `json:"id" bun:"id"`
	Username string  `json:"username" bun:"username"`
	Password *string `json:"-" bun:"password"`
	Role     string  `json:"-" bun:"role"`
	Status   bool    `json:"status" bun:"status"`
}
