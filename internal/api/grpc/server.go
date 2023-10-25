package grpc

import (
	"context"
	"fmt"
	userv1 "go-test-grpc-http/internal/api/grpc/gen/servertemplate/user/v1"
	"go-test-grpc-http/internal/api/grpc/presenter"
	"go-test-grpc-http/internal/db"
	"go-test-grpc-http/internal/repository"
	"go-test-grpc-http/internal/usecase"
	"net"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"

	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"goa.design/goa/v3/grpc/middleware"

	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
)

type server struct {
	addr   string
	server *grpc.Server
	db     *sqlx.DB
	logger *zap.Logger
}

func NewServer(addr string, db *sqlx.DB, logger *zap.Logger) *server {
	grpcServer := &server{
		addr:   addr,
		db:     db,
		logger: logger,
	}

	recoveryHandler := func(p interface{}) (err error) {
		return fmt.Errorf("grpc recovery from panic: %s", p)
	}
	recoveryOpts := []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandler(recoveryHandler),
	}

	interceptor := NewInterceptor()

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpc_recovery.UnaryServerInterceptor(recoveryOpts...),
			middleware.UnaryRequestID(
				middleware.UseXRequestIDMetadataOption(true),
				middleware.XRequestMetadataLimitOption(128),
			),
			grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpc_zap.UnaryServerInterceptor(logger, grpc_zap.WithLevels(grpcServer.grpcCodeToZapLevel)),
			interceptor.Unary(),
		),
		grpc.ChainStreamInterceptor(
			grpc_recovery.StreamServerInterceptor(recoveryOpts...),
			middleware.StreamRequestID(
				middleware.UseXRequestIDMetadataOption(true),
				middleware.XRequestMetadataLimitOption(128),
			),
			grpc_ctxtags.StreamServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpc_zap.StreamServerInterceptor(logger, grpc_zap.WithLevels(grpcServer.grpcCodeToZapLevel)),
			interceptor.Stream(),
		),
	)
	grpcServer.server = s

	return grpcServer
}

func (s *server) Run(ctx context.Context) error {
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("grpc listen error: %w", err)
	}

	s.registerServers()
	go func() {
		<-ctx.Done()
		err := s.Shutdown(ctx)
		if err != nil {
			s.logger.Fatal("can't shutdown grpc server", zap.Error(err))
			return
		}
	}()

	return s.server.Serve(lis)
}

func (s *server) Shutdown(_ context.Context) error {
	s.server.GracefulStop()
	return nil
}

func (s *server) grpcCodeToZapLevel(code codes.Code) zapcore.Level {
	switch code {
	case codes.OK:
		return zapcore.DebugLevel
	default:
		return grpc_zap.DefaultCodeToLevel(code)
	}
}

func (s *server) registerServers() {
	pgSource := db.NewSource(s.db)

	userRepository := repository.NewUserRepository(pgSource)
	userInteractor := usecase.NewUserInteractor(userRepository)
	userPresenter := presenter.NewUserPresenter()
	userv1.RegisterUserAPIServer(s.server, NewUserServer(userInteractor, userPresenter))

	// Серверная рефлексия
	reflection.Register(s.server)
}
