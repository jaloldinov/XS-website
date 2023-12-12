package post_file

import (
	"mime/multipart"
	"time"

	"github.com/uptrace/bun"
)

type Filter struct {
	Limit  *int
	Offset *int
	PostId *string
	Type   *string
}

type CreatePostFileRequest struct {
	Id         *string               `json:"id" form:"id"`
	File       *multipart.FileHeader `json:"-" form:"file"`
	FileLink   *string               `json:"-" form:"-"`
	Type       string                `json:"type" form:"type"`
	MarkedLink *string               `json:"marked_link" form:"marked_link"`
	Grouping   *string               `json:"grouping" form:"grouping"`
	Carusel    *bool                 `json:"carusel" form:"carusel"`
	PostId     *string               `json:"post_id" form:"post_id"`
	AuthorId   *string               `json:"author_id" form:"author_id"`
	CreatedAt  *time.Time            `json:"-" bun:"created_at"`
	CreatedBy  string                `json:"-" bun:"created_by"`
	UpdatedAt  *time.Time            `json:"-" bun:"updated_at"`
	UpdatedBy  *string               `json:"-" bun:"updated_by"`
	DeletedAt  *time.Time            `json:"-" bun:"deleted_at"`
	DeletedBy  *string               `json:"-" bun:"deleted_by"`
}

type CreatePostFileResponse struct {
	bun.BaseModel `bun:"table:post_file"`

	Id         string     `json:"id" bun:"id,pk"`
	Link       string     `json:"link" bun:"link"`
	Type       string     `json:"type" bun:"type"`
	MarkedLink string     `json:"marked_link" bun:"marked_link"`
	Grouping   string     `json:"grouping" bun:"grouping"`
	Carusel    bool       `json:"carusel" bun:"carusel"`
	PostId     string     `json:"post_id" bun:"post_id"`
	AuthorId   string     `json:"author_id" bun:"author_id"`
	CreatedBy  string     `json:"-" bun:"created_by"`
	CreatedAt  *time.Time `json:"-" bun:"created_at"`
}

type GetPostFileResponse struct {
	bun.BaseModel `bun:"table:post_file"`

	Id         string    `json:"id" bun:"id,pk"`
	Link       string    `json:"link" bun:"link"`
	Type       string    `json:"type" bun:"type"`
	MarkedLink string    `json:"marked_link" bun:"marked_link"`
	Grouping   string    `json:"grouping" bun:"grouping"`
	Carusel    bool      `json:"carusel" bun:"carusel"`
	PostId     string    `json:"post_id" bun:"post_id"`
	AuthorId   string    `json:"author_id" bun:"author_id"`
	CreatedAt  time.Time `json:"created_at" bun:"created_at"`
	CreatedBy  *string   `json:"created_by" bun:"created_by"`
}

type GetPostFileListResponse struct {
	bun.BaseModel `bun:"table:post_file"`

	Id     *string `json:"id" bun:"id"`
	Link   *string `json:"link" bun:"link"`
	Type   *string `json:"type" bun:"type"`
	PostId *string `json:"post_id" bun:"post_id"`
	// MarkedLink string  `json:"marked_link" bun:"marked_link"`
	// Grouping   string  `json:"grouping" bun:"grouping"`
	Carusel  bool   `json:"carusel" bun:"carusel"`
	AuthorId string `json:"author_id" bun:"author_id"`
}

type UpdatePostFileRequest struct {
	Id         string                `json:"id" form:"id"`
	File       *multipart.FileHeader `json:"-" form:"file"`
	FileLink   *string               `json:"link" form:"link"`
	Type       string                `json:"type" form:"type"`
	MarkedLink *string               `json:"marked_link" form:"marked_link"`
	Grouping   *string               `json:"grouping" form:"grouping"`
	Carusel    *bool                 `json:"carusel" form:"carusel"`
	PostId     *string               `json:"post_id" form:"post_id"`
	AuthorId   *string               `json:"author_id" form:"author_id"`
}
