package repository

func (r *repository) DeletePost(id string) error {
	query := `DELETE FROM posts WHERE id = $1`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}
