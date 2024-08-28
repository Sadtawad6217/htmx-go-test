package repository

import "mywebsite.tv/name/cmd/model"

func (r *repository) GetPostAll(limit, offset int, searchTitle string, published bool) ([]model.Posts, error) {
	var posts []model.Posts
	query := `SELECT * 
					FROM posts 
					WHERE title ILIKE '%' || $1 || '%' AND published = $2 
					ORDER BY title ASC 
					LIMIT $3 OFFSET $4;
					`
	err := r.db.Select(&posts, query, searchTitle, published, limit, offset)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (r *repository) GetTotalPostCount(searchTitle string, published bool) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM posts WHERE title ILIKE '%' || $1 || '%' AND published = $2`
	err := r.db.Get(&count, query, searchTitle, published)
	if err != nil {
		return 0, err
	}
	return count, nil
}
