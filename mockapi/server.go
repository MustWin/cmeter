package mockapi

import (
	"github.com/MustWin/cmeter/configuration"
	"github.com/MustWin/cmeter/context"
)

type Server struct {
}

func (server *Server) ListenAndServe() error {
	return nil
}

func NewServer(ctx context.Context, config *configuration.Config) (*Server, error) {
	return &Server{}, nil
}
