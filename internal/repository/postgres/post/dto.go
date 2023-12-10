package post

import (
	"time"

	"github.com/uptrace/bun"
)

type Filter struct {
	Limit         *int
	Offset        *int
	Lang          *string
	From          *string
	To            *string
	Title         *string
	Content       *string
	Status        *bool
	Order         *string
	PublishedAt   *string
	PublishedFrom *string
	PublishedTo   *string
}

type CreatePostRequest struct {
	bun.BaseModel `bun:"table:posts"`

	Id        string            `json:"id" bun:"id,pk"`
	Title     map[string]string `json:"title" bun:"title"`
	Content   map[string]string `json:"content" bun:"content"`
	Status    bool              `json:"status" bun:"status"`
	PubDate   *string           `json:"pub_date" bun:"pub_date"`
	Slug      *string           `json:"slug" bun:"slug"`
	AuthorId  *string           `json:"author_id" bun:"author_id"`
	MenuId    *string           `json:"menu_id" bun:"menu_id"`
	CreatedAt *time.Time        `json:"-" bun:"created_at"`
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
	Slug      *string           `json:"slug" bun:"slug"`
	AuthorId  *string           `json:"author_id" bun:"author_id"`
	MenuId    *string           `json:"menu_id" bun:"menu_id"`
	CreatedBy string            `json:"-" bun:"created_by"`
	CreatedAt *time.Time        `json:"-" bun:"created_at"`
}

type GetPostResponse struct {
	bun.BaseModel `bun:"table:posts"`

	Id           string            `json:"id" bun:"id,pk"`
	Title        map[string]string `json:"title" bun:"title"`
	Content      map[string]string `json:"content" bun:"content"`
	Status       bool              `json:"status" bun:"status"`
	PubDate      *string           `json:"pub_date" bun:"pub_date"`
	Slug         *string           `json:"slug" bun:"slug"`
	MenuId       *string           `json:"menu_id" bun:"menu_id"`
	AuthorName   string            `json:"author_name" bun:"full_name"`
	AuthorAvatar string            `json:"author_avatar" bun:"avatar"`
}

type GetPostListResponse struct {
	bun.BaseModel `bun:"table:posts"`

	Id               string   `json:"id" bun:"id,pk"`
	Title            string   `json:"title" bun:"title"`
	TitleLanguages   []string `json:"title_languages"`
	ContentLanguages []string `json:"content_languages"`
	Status           bool     `json:"status" bun:"status"`
	PubDate          *string  `json:"pub_date" bun:"pub_date"`
	AuthorName       string   `json:"author_name" bun:"full_name"`
	AuthorAvatar     string   `json:"author_avatar" bun:"avatar"`
}

type UpdatePostRequest struct {
	Id       string            `json:"id" bun:"id,pk"`
	Title    map[string]string `json:"title" bun:"title"`
	Content  map[string]string `json:"content" bun:"content"`
	Status   *bool             `json:"status" bun:"status"`
	PubDate  *string           `json:"pub_date" bun:"pub_date"`
	Slug     *string           `json:"slug" bun:"slug"`
	AuthorId *string           `json:"author_id" bun:"author_id"`
	MenuId   *string           `json:"menu_id" bun:"menu_id"`
}
