package service

import "mywebsite.tv/name/cmd/model"

func (s *Service) GetPostAll(limit, offset int, searchTitle string, published bool) ([]model.Posts, error) {
	return s.repository.GetPostAll(limit, offset, searchTitle, published)
}
