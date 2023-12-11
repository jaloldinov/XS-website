package menu_file

import (
	"mime/multipart"
	"time"

	"github.com/uptrace/bun"
)

type Filter struct {
	Limit  *int
	Offset *int
	MenuId *string
	Type   *string
}

type CreateMenuFileRequest struct {
	Id         *string               `json:"id" form:"id"`
	File       *multipart.FileHeader `json:"-" form:"file"`
	FileLink   *string               `json:"-" form:"-"`
	Type       string                `json:"type" form:"type"`
	MarkedLink *string               `json:"marked_link" form:"marked_link"`
	Grouping   *string               `json:"grouping" form:"grouping"`
	Carusel    *bool                 `json:"carusel" form:"carusel"`
	MenuId     *string               `json:"menu_id" form:"menu_id"`
	AuthorId   *string               `json:"author_id" form:"author_id"`
	CreatedAt  *time.Time            `json:"-" bun:"created_at"`
	CreatedBy  string                `json:"-" bun:"created_by"`
	UpdatedAt  *time.Time            `json:"-" bun:"updated_at"`
	UpdatedBy  *string               `json:"-" bun:"updated_by"`
	DeletedAt  *time.Time            `json:"-" bun:"deleted_at"`
	DeletedBy  *string               `json:"-" bun:"deleted_by"`
}

type CreateMenuFileResponse struct {
	bun.BaseModel `bun:"table:menu_file"`

	Id         string     `json:"id" bun:"id,pk"`
	Link       string     `json:"link" bun:"link"`
	Type       string     `json:"type" bun:"type"`
	MarkedLink string     `json:"marked_link" bun:"marked_link"`
	Grouping   string     `json:"grouping" bun:"grouping"`
	Carusel    bool       `json:"carusel" bun:"carusel"`
	MenuId     string     `json:"menu_id" bun:"menu_id"`
	AuthorId   string     `json:"author_id" bun:"author_id"`
	CreatedBy  string     `json:"-" bun:"created_by"`
	CreatedAt  *time.Time `json:"-" bun:"created_at"`
}

type GetMenuFileResponse struct {
	bun.BaseModel `bun:"table:menu_file"`

	Id         string    `json:"id" bun:"id,pk"`
	Link       string    `json:"link" bun:"link"`
	Type       string    `json:"type" bun:"type"`
	MarkedLink string    `json:"marked_link" bun:"marked_link"`
	Grouping   string    `json:"grouping" bun:"grouping"`
	Carusel    bool      `json:"carusel" bun:"carusel"`
	MenuId     string    `json:"menu_id" bun:"menu_id"`
	AuthorId   string    `json:"author_id" bun:"author_id"`
	CreatedAt  time.Time `json:"created_at" bun:"created_at"`
	CreatedBy  *string   `json:"created_by" bun:"created_by"`
}

type GetMenuFileListResponse struct {
	bun.BaseModel `bun:"table:menu_file"`

	Id     *string `json:"id" bun:"id"`
	Link   *string `json:"link" bun:"link"`
	Type   *string `json:"type" bun:"type"`
	MenuId *string `json:"menu_id" bun:"menu_id"`
	// MarkedLink string  `json:"marked_link" bun:"marked_link"`
	// Grouping   string  `json:"grouping" bun:"grouping"`
	Carusel  bool   `json:"carusel" bun:"carusel"`
	AuthorId string `json:"author_id" bun:"author_id"`
}

type UpdateMenuFileRequest struct {
	Id         string                `json:"id" form:"id"`
	File       *multipart.FileHeader `json:"-" form:"file"`
	FileLink   *string               `json:"link" form:"link"`
	Type       string                `json:"type" form:"type"`
	MarkedLink *string               `json:"marked_link" form:"marked_link"`
	Grouping   *string               `json:"grouping" form:"grouping"`
	Carusel    *bool                 `json:"carusel" form:"carusel"`
	MenuId     *string               `json:"menu_id" form:"menu_id"`
	AuthorId   *string               `json:"author_id" form:"author_id"`
}
