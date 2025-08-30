package server

import (
	"context"
	"fmt"

	interceptor "github.com/kkiling/goplatform/interseptors"
	"github.com/kkiling/goplatform/log"
	"github.com/kkiling/goplatform/server"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// CustomHandlerService хендлер касромного сервер
type CustomHandlerService interface {
	server.HandlerService
}

// CustomServer касромный сервер
type CustomServer struct {
	server   *server.Server
	logger   log.Logger
	services []CustomHandlerService
}

// NewCustomServer новый сервер
func NewCustomServer(
	logger log.Logger,
	serverConfig server.Config,
	services ...CustomHandlerService,
) *CustomServer {
	return &CustomServer{
		server:   server.NewServer(logger, serverConfig),
		logger:   logger.Named("media_delivery_server"),
		services: services,
	}
}

// NewPanicRecoverInterceptor интерсептор паники
func NewPanicRecoverInterceptor(logger log.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				logger.Error("panic recovered",
					zap.Any("panic_value", r),
					zap.Stack("stack"),
				)
				err = server.ErrInternal(fmt.Errorf("panic recovered: %v", r))
			}
		}()

		resp, err = handler(ctx, req)
		return resp, err
	}
}

func (p *CustomServer) unaryServerInterceptors() ([]grpc.UnaryServerInterceptor, error) {
	return []grpc.UnaryServerInterceptor{
		NewPanicRecoverInterceptor(p.logger),
		interceptor.NewLoggerInterceptor(p.logger),
	}, nil
}

// Start запустить сервер
func (p *CustomServer) Start(ctx context.Context, swaggerName string) error {
	unaryServerInterceptors, err := p.unaryServerInterceptors()
	if err != nil {
		return fmt.Errorf("p.unaryServerInterceptors(): %w", err)
	}

	p.server.WitUnaryServerInterceptor(unaryServerInterceptors...)

	var impl []server.HandlerService
	for _, service := range p.services {
		impl = append(impl, service)
	}

	if err = p.server.Start(ctx, swaggerName, impl...); err != nil {
		return fmt.Errorf("server.Start: %w", err)
	}

	return nil
}

// Stop остановить сервер
func (p *CustomServer) Stop() {
	p.server.Stop()
}
