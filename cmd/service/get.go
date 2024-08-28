package service

import "mywebsite.tv/name/cmd/model"

func (s *Service) GetPostByID(id string) (model.Posts, error) {
	s.repository.IncrementViewCount(id)
	return s.repository.GetPostByID(id)
}

func (s *Service) GetTotalPostCount(searchTitle string, published bool) (int, error) {
	// Call repository method to get total count
	count, err := s.repository.GetTotalPostCount(searchTitle, published)
	if err != nil {
		return 0, err
	}

	return count, nil
}
