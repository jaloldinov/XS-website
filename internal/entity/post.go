package entity

import (
	"time"

	"github.com/uptrace/bun"
)

type Post struct {
	bun.BaseModel `bun:"table:posts"`

	Id        string            `json:"id" bun:"id,pk"`
	Title     map[string]string `json:"title" bun:"title"`
	Content   map[string]string `json:"content" bun:"content"`
	Status    *bool             `json:"status" bun:"status"`
	PubDate   *time.Time        `json:"pub_date" bun:"pub_date"`
	AuthorId  string            `json:"author_id" bun:"author_id"`
	CreatedAt time.Time         `json:"-" bun:"created_at"`
	CreatedBy string            `json:"-" bun:"created_by"`
	UpdatedAt *time.Time        `json:"-" bun:"updated_at"`
	UpdatedBy *string           `json:"-" bun:"updated_by"`
	DeletedAt *time.Time        `json:"-" bun:"deleted_at"`
	DeletedBy *string           `json:"-" bun:"deleted_by"`
}
