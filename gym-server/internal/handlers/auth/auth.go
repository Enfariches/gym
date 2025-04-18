package auth

import (
	"context"
	"fmt"
	pb "health/protogen/v1/auth"

	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var validate *validator.Validate

type Auth interface {
	RegisterNewUser(email, password string) (err error)
	Login(email, password string) (token string, err error)
}

type AuthServerManagmentApi struct {
	auth Auth
	pb.UnimplementedAuthServiceServer
}

func RegisterGRPCServer(gRPC *grpc.Server) {
	pb.RegisterAuthServiceServer(gRPC, &AuthServerManagmentApi{})
}

func (s *AuthServerManagmentApi) Register(ctx context.Context, r *pb.AuthRequest) (*pb.RegisterResponse, error) {

	if err := validateAuth(r); err != nil {
		return nil, err
	}

	if err := s.auth.RegisterNewUser(r.Email, r.Password); err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	return &pb.RegisterResponse{}, nil

}

func (s *AuthServerManagmentApi) Login(ctx context.Context, r *pb.AuthRequest) (*pb.LoginResponse, error) {

	if err := validateAuth(r); err != nil {
		return nil, err
	}

	token, err := s.auth.Login(r.Email, r.Password)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	return &pb.LoginResponse{Token: token}, nil
}

func validateAuth(r *pb.AuthRequest) error {
	if r.Password == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}

	if r.Email == "" {
		return status.Error(codes.InvalidArgument, "email is required")
	}

	if err := validate.Var(r.Email, "required,email"); err != nil {
		return status.Error(codes.InvalidArgument, "Email is not valid format")
	}

	return nil
}
