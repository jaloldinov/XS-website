package post

import (
	"time"

	"github.com/uptrace/bun"
)

type Filter struct {
	Limit    *int
	Offset   *int
	Username *string
}

type CreatePostRequest struct {
	bun.BaseModel `bun:"table:posts"`

	Id        string            `json:"id" bun:"id,pk"`
	Title     map[string]string `json:"title" bun:"title"`
	Content   map[string]string `json:"content" bun:"content"`
	Status    bool              `json:"status" bun:"status"`
	PubDate   *string           `json:"pub_date" bun:"pub_date"`
	AuthorId  *string           `json:"author_id" bun:"author_id"`
	CreatedAt time.Time         `json:"-" bun:"created_at"`
	CreatedBy string            `json:"-" bun:"created_by"`
	UpdatedAt *time.Time        `json:"-" bun:"updated_at"`
	UpdatedBy *string           `json:"-" bun:"updated_by"`
	DeletedAt *time.Time        `json:"-" bun:"deleted_at"`
	DeletedBy *string           `json:"-" bun:"deleted_by"`
}

type CreatePostResponse struct {
	bun.BaseModel `bun:"table:posts"`

	Id        string            `json:"id" bun:"id,pk"`
	Title     map[string]string `json:"title" bun:"title"`
	Content   map[string]string `json:"content" bun:"content"`
	Status    bool              `json:"status" bun:"status"`
	PubDate   *time.Time        `json:"pub_date" bun:"pub_date"`
	AuthorId  string            `json:"author_id" bun:"author_id"`
	CreatedBy string            `json:"-" bun:"created_by"`
	CreatedAt time.Time         `json:"-" bun:"created_at"`
}

/*
type GetPostResponse struct {
	bun.BaseModel `bun:"table:posts"`

	Id        string `json:"id" bun:"id"`
	Avatar    string `json:"avatar" bun:"avatar"`
	Postname  string `json:"username" bun:"username"`
	FullName  string `json:"full_name" bun:"full_name"`
	Status    bool   `json:"status" bun:"status"`
	Role      string `json:"role" bun:"role"`
	BirthDate string `json:"birth_date" bun:"birth_date"`
	Gender    string `json:"gender" bun:"gender"`
	Phone     string `json:"phone" bun:"phone"`
}

type GetPostListResponse struct {
	bun.BaseModel `bun:"table:posts"`

	Id       string `json:"id" bun:"id"`
	Avatar   string `json:"avatar" bun:"avatar"`
	Postname string `json:"username" bun:"username"`
	FullName string `json:"full_name" bun:"full_name"`
	Status   bool   `json:"status" bun:"status"`
	Role     string `json:"role" bun:"role"`
	Phone    string `json:"phone" bun:"phone"`
}

type UpdatePostRequest struct {
	Id         string                `json:"id" bun:"id"`
	Avatar     *multipart.FileHeader `json:"-" form:"avatar"`
	AvatarLink *string               `json:"-" form:"-"`
	Postname   *string               `json:"username" form:"username"`
	FullName   *string               `json:"full_name" form:"full_name"`
	Status     *bool                 `json:"status" form:"status"`
	Role       *string               `json:"role" form:"role"`
	BirthDate  *string               `json:"birth_date" form:"birth_date"`
	Gender     *string               `json:"gender" form:"gender"`
	Phone      *string               `json:"phone" form:"phone"`
	UpdatedBy  *string               `json:"-"`
}

type DeletePostRequest struct {
	Id        string `json:"id" bun:"id"`
	DeletedBy string `json:"deleted_by" bun:"deleted_by"`
}

type UpdatePasswordRequest struct {
	Id          *string `json:"id" form:"id"`
	NewPassword *string `json:"new_password" form:"new_password"`
	UpdatedBy   *string `json:"-"`
}

type DetailPostResponse struct {
	bun.BaseModel `bun:"table:posts"`

	Id       string  `json:"id" bun:"id"`
	Postname string  `json:"username" bun:"username"`
	Password *string `json:"-" bun:"password"`
	Role     string  `json:"-" bun:"role"`
	Status   bool    `json:"status" bun:"status"`
}
*/
