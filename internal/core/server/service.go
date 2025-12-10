// Package server
package server

import (
	"context"
	"errors"

	"horizonx-server/internal/domain"
	"horizonx-server/pkg"
)

type Service struct {
	repo domain.ServerRepository
}

func NewService(repo domain.ServerRepository) domain.ServerService {
	return &Service{repo: repo}
}

func (s *Service) Register(ctx context.Context, name, ip string) (*domain.Server, string, error) {
	if name == "" {
		return nil, "", errors.New("server name cannot be empty")
	}

	token, err := pkg.GenerateToken()
	if err != nil {
		return nil, "", err
	}

	srv := &domain.Server{
		Name:      name,
		IPAddress: ip,
		APIToken:  token,
		IsOnline:  false,
	}

	if err := s.repo.Create(ctx, srv); err != nil {
		return nil, "", err
	}

	return srv, token, nil
}

func (s *Service) List(ctx context.Context) ([]domain.Server, error) {
	return s.repo.List(ctx)
}
