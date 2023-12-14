package entity

import (
	"time"

	"github.com/uptrace/bun"
)

type PostHashtag struct {
	bun.BaseModel `bun:"table:post_hashtag"`

	Id        string     `json:"id" bun:"id,pk"`
	PostId    *string    `json:"post_id" bun:"post_id"`
	HashtagId *string    `json:"hashtag_id" bun:"hashtag_id"`
	CreatedAt time.Time  `json:"-" bun:"created_at"`
	CreatedBy string     `json:"-" bun:"created_by"`
	DeletedAt *time.Time `json:"-" bun:"deleted_at"`
	DeletedBy *string    `json:"-" bun:"deleted_by"`
}
