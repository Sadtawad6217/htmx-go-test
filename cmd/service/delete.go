package service

func (s *Service) DeletePost(id string) error {
	s.repository.DeletePost(id)
	return nil
}
