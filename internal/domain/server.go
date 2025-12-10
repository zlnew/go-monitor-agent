package domain

import (
	"context"
	"time"
)

type Server struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	IPAddress string    `json:"ip_address"`
	APIToken  string    `json:"-"`
	IsOnline  bool      `json:"is_online"`
	OSInfo    *OSInfo   `json:"os_info,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ServerCreateRequest struct {
	Name      string `json:"name" validate:"required"`
	IPAddress string `json:"ip_address" validate:"required"`
}

type ServerRepository interface {
	GetByToken(ctx context.Context, token string) (*Server, error)
	Create(ctx context.Context, s *Server) error
	List(ctx context.Context) ([]Server, error)
	UpdateStatus(ctx context.Context, id int64, isOnline bool) error
}

type ServerService interface {
	Register(ctx context.Context, name, ip string) (*Server, string, error)
	List(ctx context.Context) ([]Server, error)
}
