package auth

import (
	"context"
	"fmt"
	pb "health/protogen/v1/auth"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthService interface {
	RegisterNewUser(ctx context.Context, email, password, source string) (authToken string, err error)
	Login(ctx context.Context, email, password, source string) (token string, err error)
	VerifyRegister(ctx context.Context, authToken string) error
	ChangePassword(ctx context.Context, email string, source string) (string, error)
	VerifyChangePassword(ctx context.Context, resetToken, source string) error
}

type AuthServerManagmentApi struct {
	pb.UnimplementedAuthServiceServer
	auth AuthService
}

func RegisterGRPCServer(gRPC *grpc.Server, auth AuthService) {
	pb.RegisterAuthServiceServer(gRPC, &AuthServerManagmentApi{auth: auth})
}

func (s *AuthServerManagmentApi) Register(ctx context.Context, r *pb.AuthRequest) (*pb.RegisterResponse, error) {
	if err := validateAuth(r); err != nil {
		return nil, err
	}

	source, err := sourceFromProto(r.Source)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("unknown app source: %v", err))
	}

	authToken, err := s.auth.RegisterNewUser(context.Background(), r.Email, r.Password, source)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to register user: %v", err))
	}

	return &pb.RegisterResponse{AuthToken: authToken}, nil
}

func (s *AuthServerManagmentApi) Login(ctx context.Context, r *pb.AuthRequest) (*pb.LoginResponse, error) {
	if err := validateAuth(r); err != nil {
		return nil, err
	}

	source, err := sourceFromProto(r.Source)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("unknown app source: %v", err))
	}

	token, err := s.auth.Login(context.Background(), r.Email, r.Password, source)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to login user: %v", err))
	}

	return &pb.LoginResponse{Token: token}, nil
}

func (s *AuthServerManagmentApi) VerifyRegister(ctx context.Context, r *pb.VerifyRegisterRequest) (*pb.VerifyRegisterResponse, error) {
	if r.AuthToken == "" {
		return nil, status.Error(codes.InvalidArgument, "token is required")
	}

	if err := s.auth.VerifyRegister(context.Background(), r.AuthToken); err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to verify register: %v", err))
	}

	return &pb.VerifyRegisterResponse{}, nil
}

func (s *AuthServerManagmentApi) ChangePassword(ctx context.Context, r *pb.ChangePasswordRequest) (*pb.ChangePasswordResponse, error) {
	if r.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	source, err := sourceFromProto(r.Source)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("unknown app source: %v", err))
	}

	resetToken, err := s.auth.ChangePassword(context.Background(), r.Email, source)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to reset password: %v", err))
	}

	return &pb.ChangePasswordResponse{ResetToken: resetToken}, nil
}

func (s *AuthServerManagmentApi) VerifyChangePassword(ctx context.Context, r *pb.VerifyChangePasswordRequest) (*pb.VerifyChangePasswordResponse, error) {
	if r.ResetToken == "" {
		return nil, status.Error(codes.InvalidArgument, "reset token is required")
	}

	if r.NewPassword == "" {
		return nil, status.Error(codes.InvalidArgument, "new password is required")
	}

	if err := s.auth.VerifyChangePassword(context.Background(), r.ResetToken, r.NewPassword); err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to change password: %v", err))
	}

	return &pb.VerifyChangePasswordResponse{}, nil
}

func validateAuth(r *pb.AuthRequest) error {
	if r.Password == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}

	if r.Email == "" {
		return status.Error(codes.InvalidArgument, "email is required")
	}

	return nil
}

func sourceFromProto(pbSource pb.AppSource) (string, error) {
	switch pbSource {
	case pb.AppSource_EMPLOYEE:
		return "employees", nil
	case pb.AppSource_ADMIN:
		return "admins", nil
	default:
		return "", fmt.Errorf("unknown app source: %v", pbSource)
	}
}
