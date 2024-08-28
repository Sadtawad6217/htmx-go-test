package repository

import (
	"mywebsite.tv/name/cmd/model"

	"github.com/jmoiron/sqlx"
)

type PostRepository interface {
	GetPostByID(id string) (model.Posts, error)
	GetPostAll(limit, offset int, searchTitle string, published bool) ([]model.Posts, error)
	CreatePosts(post model.Posts) (model.Posts, error)
	UpdatePost(id string, updateData model.Posts) (model.Posts, error)
	DeletePost(id string) error
	IncrementViewCount(id string) error
	GetTotalPostCount(searchTitle string, published bool) (int, error)
}

type repository struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) PostRepository {
	return &repository{db: db}
}
