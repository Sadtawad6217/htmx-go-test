package service

import (
	"fmt"

	"mywebsite.tv/name/cmd/model"
)

func (s *Service) UpdatePost(id string, updateData model.Posts) (model.Posts, error) {
	existingPost, err := s.repository.GetPostByID(id)
	if err != nil {
		return model.Posts{}, fmt.Errorf("error getting post by ID: %w", err)
	}
	if existingPost == (model.Posts{}) { // Check if post is empty
		return model.Posts{}, fmt.Errorf("post not found")
	}

	title := updateData.Title
	content := updateData.Content
	published := updateData.Published
	viewCount := updateData.ViewCount
	existingPost.Update(
		&title,
		&content,
		published,
		&viewCount,
	)

	updatedPost, err := s.repository.UpdatePost(id, existingPost)
	if err != nil {
		return model.Posts{}, fmt.Errorf("error updating post: %w", err)
	}

	return updatedPost, nil
}

func (s *Service) IncrementViewCount(id string) error {
	post, err := s.repository.GetPostByID(id)
	if err != nil {
		return fmt.Errorf("error getting post by ID: %w", err)
	}
	if post == (model.Posts{}) {
		return fmt.Errorf("post not found")
	}

	post.ViewCount++
	_, err = s.repository.UpdatePost(id, post)
	return err
}
