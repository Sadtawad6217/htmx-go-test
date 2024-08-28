package repository

import "mywebsite.tv/name/cmd/model"

func (r *repository) CreatePosts(post model.Posts) (model.Posts, error) {
	query := `INSERT INTO posts (title, content, published, view_count, created_at, updated_at, deleted_at) 
              VALUES (:title, :content, :published, :view_count, :created_at, :updated_at, :deleted_at)
              RETURNING id`

	stmt, err := r.db.PrepareNamed(query)
	if err != nil {
		return model.Posts{}, err
	}
	defer stmt.Close()

	var id string
	err = stmt.QueryRow(post).Scan(&id)
	if err != nil {
		return model.Posts{}, err
	}

	post.ID = id
	return post, nil
}
