package hashtag

import (
	"time"

	"github.com/uptrace/bun"
)

type Filter struct {
	Limit  *int
	Offset *int
}

type CreateHashtagRequest struct {
	bun.BaseModel `bun:"table:hashtags"`

	Id        string     `json:"id" bun:"id,pk"`
	Name      *string    `json:"name" bun:"name"`
	Status    *bool      `json:"status" bun:"status"`
	CreatedAt time.Time  `json:"-" bun:"created_at"`
	CreatedBy string     `json:"-" bun:"created_by"`
	UpdatedAt *time.Time `json:"-" bun:"updated_at"`
	UpdatedBy *string    `json:"-" bun:"updated_by"`
	DeletedAt *time.Time `json:"-" bun:"deleted_at"`
	DeletedBy *string    `json:"-" bun:"deleted_by"`
}

type CreateHashtagResponse struct {
	bun.BaseModel `bun:"table:hashtags"`

	Id        string    `json:"id" bun:"id,pk"`
	Name      *string   `json:"name" bun:"name"`
	Status    *bool     `json:"status" bun:"status"`
	CreatedAt time.Time `json:"created_at" bun:"created_at"`
	CreatedBy string    `json:"-" bun:"created_by"`
}

type GetHashtagResponse struct {
	bun.BaseModel `bun:"table:hashtags"`

	Id     string  `json:"id" bun:"id,pk"`
	Name   *string `json:"name" bun:"name"`
	Status *bool   `json:"status" bun:"status"`
}

type GetHashtagListResponse struct {
	bun.BaseModel `bun:"table:hashtags"`

	Id     string  `json:"id" bun:"id,pk"`
	Name   *string `json:"name" bun:"name"`
	Status *bool   `json:"status" bun:"status"`
}

type UpdateHashtagRequest struct {
	bun.BaseModel `bun:"table:hashtags"`

	Id     string  `json:"id" bun:"id,pk"`
	Name   *string `json:"name" bun:"name"`
	Status *bool   `json:"status" bun:"status"`
}
