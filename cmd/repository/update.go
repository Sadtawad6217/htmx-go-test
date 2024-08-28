package repository

import (
	"fmt"

	"mywebsite.tv/name/cmd/model"
)

func (r *repository) UpdatePost(id string, updateData model.Posts) (model.Posts, error) {
	query := `UPDATE posts SET title = :title, content = :content, published = :published, 
              updated_at = :updated_at WHERE id = :id RETURNING *`
	stmt, err := r.db.PrepareNamed(query)
	if err != nil {
		return model.Posts{}, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	updateData.ID = id
	var updatedPost model.Posts
	err = stmt.Get(&updatedPost, updateData)
	if err != nil {
		return model.Posts{}, fmt.Errorf("failed to execute query: %w", err)
	}
	return updatedPost, nil
}

func (r *repository) IncrementViewCount(id string) error {
	query := `UPDATE posts SET view_count = view_count + 1 WHERE id = $1`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}
