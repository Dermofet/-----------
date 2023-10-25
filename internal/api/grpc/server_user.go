package grpc

import (
	"context"
	userv1 "go-test-grpc-http/internal/api/grpc/gen/servertemplate/user/v1"
	"go-test-grpc-http/internal/api/grpc/presenter"
	"go-test-grpc-http/internal/usecase"

	"google.golang.org/grpc/codes"
)

type userServer struct {
	interactor usecase.UserInteractor
	presenter  presenter.UserPresenter
	userv1.UnimplementedUserAPIServer
}

func NewUserServer(interactor usecase.UserInteractor, presenter presenter.UserPresenter) userv1.UserAPIServer {
	return &userServer{
		interactor: interactor,
		presenter:  presenter,
	}
}

func (s *userServer) GetById(ctx context.Context, request *userv1.GetByIdRequest) (*userv1.GetByIdResponse, error) {
	userId := s.presenter.ToUserID(request.GetId())
	if userId == nil {
		return nil, NewApiError(codes.Internal, "get user error: id is invalid")
	}
	user, err := s.interactor.GetById(ctx, userId)
	if err != nil {
		return nil, NewApiError(codes.Internal, "get user error: %v", err)
	}

	return &userv1.GetByIdResponse{
		User: s.presenter.FromUser(user),
	}, nil
}

func (s *userServer) GetByEmail(ctx context.Context, request *userv1.GetByEmailRequest) (*userv1.GetByEmailResponse, error) {
	user, err := s.interactor.GetByEmail(ctx, request.GetEmail())
	if err != nil {
		return nil, NewApiError(codes.Internal, "get user error: %v", err)
	}

	return &userv1.GetByEmailResponse{
		User: s.presenter.FromUser(user),
	}, nil
}

func (s *userServer) Update(ctx context.Context, request *userv1.UpdateRequest) (*userv1.UpdateResponse, error) {
	userId := s.presenter.ToUserID(request.GetId())
	if userId == nil {
		return nil, NewApiError(codes.Internal, "get user error: id is invalid")
	}

	user := s.presenter.ToUserCreate(request.User)

	userDB, err := s.interactor.Update(ctx, userId, user)
	if err != nil {
		return nil, NewApiError(codes.Internal, "get user error", err)
	}

	return &userv1.UpdateResponse{
		User: s.presenter.FromUser(userDB),
	}, nil
}

func (s *userServer) Delete(ctx context.Context, request *userv1.DeleteRequest) (*userv1.DeleteResponse, error) {
	userId := s.presenter.ToUserID(request.GetId())
	if userId == nil {
		return nil, NewApiError(codes.Internal, "get user error: id is invalid")
	}
	err := s.interactor.Delete(ctx, userId)
	if err != nil {
		return nil, NewApiError(codes.Internal, "get user error: %v", err)
	}

	return &userv1.DeleteResponse{}, nil
}
