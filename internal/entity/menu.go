package entity

import (
	"time"

	"github.com/uptrace/bun"
)

type Menu struct {
	bun.BaseModel `bun:"table:menu"`

	Id        string             `json:"id" bun:"id,pk"`
	Title     *map[string]string `json:"title" bun:"title"`
	Content   map[string]string  `json:"content" bun:"content"`
	ParentId  *string            `json:"parent_id" bun:"parent_id"`
	IsStatic  *bool              `json:"is_static" bun:"is_static"`
	Status    *bool              `json:"status" bun:"status"`
	Slug      *string            `json:"slug" bun:"slug"`
	Type      *string            `json:"type" bun:"type"`
	CreatedAt time.Time          `json:"-" bun:"created_at"`
	CreatedBy string             `json:"-" bun:"created_by"`
	UpdatedAt *time.Time         `json:"-" bun:"updated_at"`
	UpdatedBy *string            `json:"-" bun:"updated_by"`
	DeletedAt *time.Time         `json:"-" bun:"deleted_at"`
	DeletedBy *string            `json:"-" bun:"deleted_by"`
}
