package model

import (
	"time"
)

type Posts struct {
	ID        string     `db:"id" json:"id"`
	Title     string     `db:"title" json:"title"`
	Content   string     `db:"content" json:"content"`
	Published bool       `db:"published" json:"published"`
	ViewCount int64      `db:"view_count" json:"view_count"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
}

func New(title, content string, published bool) *Posts {
	return &Posts{
		Title:     title,
		Content:   content,
		Published: published,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
func (p *Posts) Update(title, content *string, published bool, viewCount *int64) {
	if title != nil {
		p.Title = *title
	}
	if content != nil {
		p.Content = *content
	}
	if published {
		p.Published = true
	}
	if viewCount != nil {
		p.ViewCount = *viewCount
	}
	p.UpdatedAt = time.Now()
}
