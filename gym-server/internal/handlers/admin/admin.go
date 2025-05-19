package admin

import (
	"context"
	"health/internal/domain/models"
	"health/lib/ctxkey"
	pb "health/protogen/v1/users"

	"github.com/mennanov/fmutils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

type AdminService interface {
	GetAdmin(ctx context.Context, admin_id int64) (*models.Admin, error)
	UpdateAdmin(ctx context.Context, fieldMask map[string]any) (*models.Admin, error)
	ListAdminEmployees(ctx context.Context, department_id int64) ([]*models.Employee, error)
}

type AdminServerManagmentApi struct {
	pb.UnimplementedAdminServiceServer
	admin AdminService
}

func RegisterGRPCServer(gRPC *grpc.Server, admin AdminService) {
	pb.RegisterAdminServiceServer(gRPC, &AdminServerManagmentApi{admin: admin})
}

func (s *AdminServerManagmentApi) GetAdmin(ctx context.Context, _ *emptypb.Empty) (*pb.Admin, error) {
	admin_id := ctx.Value(ctxkey.UserKey).(int64)

	admin, err := s.admin.GetAdmin(ctx, admin_id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get admin: %v", err)
	}

	return &pb.Admin{
		Id:         admin.ID,
		Name:       admin.Name,
		Surname:    admin.Surname,
		Email:      admin.Email,
		Department: admin.Department,
	}, nil
}

func (s *AdminServerManagmentApi) UpdateAdmin(ctx context.Context, r *pb.UpdateAdminRequest) (*pb.Admin, error) {
	r.FieldMask.Normalize()

	if !r.FieldMask.IsValid(r.Admin) {
		return nil, status.Errorf(codes.InvalidArgument, "invalid field mask")
	}

	fmutils.Filter(r.GetAdmin(), r.FieldMask.GetPaths())

	updateFields, err := applyFieldMask(r.Admin, r.FieldMask)
	if err != nil {
		return nil, err
	}

	updatedAdmin, err := s.admin.UpdateAdmin(ctx, updateFields)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update admin: %v", err)
	}

	return &pb.Admin{
		Id:         updatedAdmin.ID,
		Name:       updatedAdmin.Name,
		Surname:    updatedAdmin.Surname,
		Email:      updatedAdmin.Email,
		Department: updatedAdmin.Department,
	}, nil
}

func (s *AdminServerManagmentApi) ListAdminEmployees(ctx context.Context, _ *emptypb.Empty) (*pb.ListAdminEmployeesResponse, error) {
	department_id := ctx.Value(ctxkey.DepartmentKey).(int64)

	employees, err := s.admin.ListAdminEmployees(ctx, department_id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list admin employees: %v", err)
	}

	employeesPb := makePbEmployees(employees)

	return &pb.ListAdminEmployeesResponse{
		Employees: employeesPb,
	}, nil
}

func applyFieldMask(req *pb.Admin, mask *fieldmaskpb.FieldMask) (map[string]any, error) {
	updateMap := make(map[string]any)

	for _, path := range mask.GetPaths() {
		switch path {
		case "name":
			updateMap[path] = req.Name
		case "surname":
			updateMap[path] = req.Surname
		case "department":
			updateMap[path] = req.Department
		default:
			return nil, status.Errorf(codes.InvalidArgument, "incorrect fields: %s", path)
		}
	}

	for _, value := range updateMap {
		if value == "" {
			return nil, status.Error(codes.InvalidArgument, "field value is required")
		}
	}
	return updateMap, nil
}

func makePbEmployees(savedEmployees []*models.Employee) []*pb.Employee {
	responseEmployees := make([]*pb.Employee, 0, len(savedEmployees))

	for _, e := range savedEmployees {
		responseEmployees = append(responseEmployees, &pb.Employee{
			Id:         e.ID,
			Name:       e.Name,
			SecondName: e.SecondName,
			Surname:    e.Surname,
			Email:      e.Email,
			Age:        e.Age,
			Sex:        e.Sex,
			Phone:      e.Phone,
			Department: e.Department,
			Post:       e.Post,
		})
	}
	return responseEmployees
}
