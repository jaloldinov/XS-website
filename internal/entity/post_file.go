package entity

import (
	"time"

	"github.com/uptrace/bun"
)

type PostFile struct {
	bun.BaseModel `bun:"table:post_file"`

	Id         string     `json:"id" bun:"id,pk"`
	Link       string     `json:"link" bun:"link"`
	Type       string     `json:"type" bun:"type"`
	MarkedLink string     `json:"marked_link" bun:"marked_link"`
	Grouping   string     `json:"grouping" bun:"grouping"`
	Carusel    bool       `json:"carusel" bun:"carusel"`
	PostId     string     `json:"post_id" bun:"post_id"`
	AuthorId   string     `json:"author_id" bun:"author_id"`
	CreatedAt  time.Time  `json:"created_at" bun:"created_at"`
	CreatedBy  *string    `json:"created_by" bun:"created_by"`
	UpdatedAt  *time.Time `json:"updated_at" bun:"updated_at"`
	UpdatedBy  *string    `json:"updated_by" bun:"updated_by"`
	DeletedAt  *time.Time `json:"deleted_at" bun:"deleted_at"`
	DeletedBy  *string    `json:"deleted_by" bun:"deleted_by"`
}
