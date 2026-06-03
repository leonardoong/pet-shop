package dashboard

import (
	dashdto "petshop/internal/dto/dashboard"
	dashrepo "petshop/internal/repository/dashboard"
)

type Service interface {
	GetStats() (*dashdto.Stats, error)
}

type service struct {
	repo dashrepo.Repository
}

func NewService(repo dashrepo.Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetStats() (*dashdto.Stats, error) {
	return s.repo.GetStats()
}
