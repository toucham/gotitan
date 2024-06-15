package server

import "context"

type Server struct {
	ctx context.Context
}

func InitServer(ctx context.Context) *Server {
	return &Server{
		ctx,
	}
}
