package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// HTTPConfig holds HTTP server configuration.
type HTTPConfig struct {
	Host         string
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// HTTPServer wraps the Gin HTTP server.
type HTTPServer struct {
	server *http.Server
	engine *gin.Engine
}

func NewHTTPServer(cfg HTTPConfig) *HTTPServer {
	engine := gin.Default()
	engine.Use(gin.Recovery())

	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Handler:      engine,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	return &HTTPServer{
		server: server,
		engine: engine,
	}
}

func (s *HTTPServer) Start() error {
	return s.server.ListenAndServe()
}

func (s *HTTPServer) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *HTTPServer) Engine() *gin.Engine {
	return s.engine
}

// GRPCConfig holds gRPC server configuration.
type GRPCConfig struct {
	Host string
	Port int
}

// GRPCServer wraps the gRPC server.
type GRPCServer struct {
	server *grpc.Server
}

func NewGRPCServer(cfg GRPCConfig) *GRPCServer {
	server := grpc.NewServer(
		grpc.UnaryInterceptor(chainUnaryInterceptors()),
		grpc.StreamInterceptor(chainStreamInterceptors()),
	)

	// Enable reflection for development
	reflection.Register(server)

	return &GRPCServer{
		server: server,
	}
}

func (s *GRPCServer) Serve(lis net.Listener) error {
	return s.server.Serve(lis)
}

func (s *GRPCServer) GracefulStop() {
	s.server.GracefulStop()
}

func (s *GRPCServer) Server() *grpc.Server {
	return s.server
}

func chainUnaryInterceptors() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		return handler(ctx, req)
	}
}

func chainStreamInterceptors() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		return handler(srv, ss)
	}
}
