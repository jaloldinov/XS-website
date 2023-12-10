package menu

import (
	"time"

	"github.com/uptrace/bun"
)

type Filter struct {
	Limit    *int
	Offset   *int
	Lang     *string
	ParentId *string
}

type CreateMenuRequest struct {
	bun.BaseModel `bun:"table:menu"`

	Id        string             `json:"id" bun:"id,pk"`
	Title     *map[string]string `json:"title" bun:"title"`
	Content   map[string]string  `json:"content" bun:"content"`
	ParentId  *string            `json:"parent_id" bun:"parent_id"`
	IsStatic  *bool              `json:"is_static" bun:"is_static"`
	Status    *bool              `json:"status" bun:"status"`
	Slug      *string            `json:"slug" bun:"slug"`
	Type      string             `json:"type" bun:"type"`
	CreatedAt time.Time          `json:"-" bun:"created_at"`
	CreatedBy string             `json:"-" bun:"created_by"`
	UpdatedAt *time.Time         `json:"-" bun:"updated_at"`
	UpdatedBy *string            `json:"-" bun:"updated_by"`
	DeletedAt *time.Time         `json:"-" bun:"deleted_at"`
	DeletedBy *string            `json:"-" bun:"deleted_by"`
}

type CreateMenuResponse struct {
	bun.BaseModel `bun:"table:menu"`

	Id        string             `json:"id" bun:"id,pk"`
	Title     *map[string]string `json:"title" bun:"title"`
	Content   map[string]string  `json:"content" bun:"content"`
	ParentId  *string            `json:"parent_id" bun:"parent_id"`
	IsStatic  *bool              `json:"is_static" bun:"is_static"`
	Status    *bool              `json:"status" bun:"status"`
	Slug      *string            `json:"slug" bun:"slug"`
	Type      *string            `json:"type" bun:"type"`
	CreatedBy string             `json:"-" bun:"created_by"`
	CreatedAt *time.Time         `json:"-" bun:"created_at"`
}

type GetMenuListResponse struct {
	bun.BaseModel `bun:"table:menu"`

	Id               string                 `json:"id"`
	Title            string                 `json:"title"`
	Content          string                 `json:"content"`
	TitleLanguages   []string               `json:"title_languages"`
	ContentLanguages []string               `json:"content_languages"`
	Status           bool                   `json:"status"`
	Slug             string                 `json:"slug"`
	ParentId         *string                `json:"parent_id"`
	Children         *[]GetMenuListResponse `json:"children"`
}

type UpdateMenuRequest struct {
	Id        string            `json:"id" bun:"id,pk"`
	Title     map[string]string `json:"title" bun:"title"`
	Content   map[string]string `json:"content" bun:"content"`
	Status    *bool             `json:"status" bun:"status"`
	PubDate   *string           `json:"pub_date" bun:"pub_date"`
	AuthorId  *string           `json:"author_id" bun:"author_id"`
	UpdatedBy *string           `json:"-"`
}
