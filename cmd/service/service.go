package service

import "mywebsite.tv/name/cmd/repository"

type Service struct {
	repository repository.PostRepository
}

func NewService(repo repository.PostRepository) *Service {
	return &Service{repository: repo}
}
