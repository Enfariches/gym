package employees

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

type EmployeeService interface {
	GetEmployee(ctx context.Context, employee_id int64) (*models.Employee, error)
	UpdateEmployee(ctx context.Context, fieldMask map[string]any) (*models.Employee, error)
	DeleteEmployee(ctx context.Context, employee_id int64) error

	ListDepartments(ctx context.Context) ([]*models.Department, error)
}

type EmployeeServerManagmentApi struct {
	pb.UnimplementedEmployeeServiceServer
	employee EmployeeService
}

func RegisterGRPCServer(gRPC *grpc.Server, employee EmployeeService) {
	pb.RegisterEmployeeServiceServer(gRPC, &EmployeeServerManagmentApi{employee: employee})
}

func (s *EmployeeServerManagmentApi) GetEmployee(ctx context.Context, _ *emptypb.Empty) (*pb.Employee, error) {
	employee_id := ctx.Value(ctxkey.UserKey).(int64)

	employee, err := s.employee.GetEmployee(ctx, employee_id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get employee: %v", err)
	}

	return &pb.Employee{
		Id:         employee.ID,
		Name:       employee.Name,
		SecondName: employee.SecondName,
		Surname:    employee.Surname,
		Email:      employee.Email,
		Age:        employee.Age,
		Sex:        employee.Sex,
		Phone:      employee.Phone,
		Department: employee.Department,
		Post:       employee.Post,
	}, nil
}

func (s *EmployeeServerManagmentApi) UpdateEmployee(ctx context.Context, r *pb.UpdateEmployeeRequest) (*pb.Employee, error) {
	r.FieldMask.Normalize()

	if !r.FieldMask.IsValid(r.Employee) {
		return nil, status.Errorf(codes.InvalidArgument, "invalid field mask")
	}

	fmutils.Filter(r.GetEmployee(), r.FieldMask.GetPaths())

	updateFields, err := applyFieldMask(r.Employee, r.FieldMask)
	if err != nil {
		return nil, err
	}

	updatedEmployee, err := s.employee.UpdateEmployee(ctx, updateFields)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update employee: %v", err)
	}
	return &pb.Employee{
		Id:         updatedEmployee.ID,
		Name:       updatedEmployee.Name,
		SecondName: updatedEmployee.SecondName,
		Surname:    updatedEmployee.Surname,
		Email:      updatedEmployee.Email,
		Age:        updatedEmployee.Age,
		Sex:        updatedEmployee.Sex,
		Phone:      updatedEmployee.Phone,
		Department: updatedEmployee.Department,
		Post:       updatedEmployee.Post,
	}, nil
}

func (s *EmployeeServerManagmentApi) DeleteEmployee(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	employee_id := ctx.Value(ctxkey.UserKey).(int64)

	err := s.employee.DeleteEmployee(ctx, employee_id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete employee: %v", err)
	}

	return &emptypb.Empty{}, nil
}

func (s *EmployeeServerManagmentApi) ListDepartments(ctx context.Context, _ *emptypb.Empty) (*pb.ListDepartmentsResponse, error) {
	departments, err := s.employee.ListDepartments(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list departments: %v", err)
	}

	responseDepartments := makePbDepartments(departments)

	return &pb.ListDepartmentsResponse{
		Departments: responseDepartments,
	}, nil
}

func applyFieldMask(req *pb.Employee, mask *fieldmaskpb.FieldMask) (map[string]any, error) {
	updateMap := make(map[string]any)

	for _, path := range mask.GetPaths() {
		switch path {
		case "name":
			updateMap[path] = req.Name
		case "second_name":
			updateMap[path] = req.SecondName
		case "surname":
			updateMap[path] = req.Surname
		case "age":
			updateMap[path] = req.Age
		case "sex":
			updateMap[path] = req.Sex
		case "phone":
			updateMap[path] = req.Phone
		case "department":
			updateMap[path] = req.Department
		case "post":
			updateMap[path] = req.Post
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

func makePbDepartments(savedDepartment []*models.Department) []*pb.Department {
	responseDepartments := make([]*pb.Department, 0, len(savedDepartment))

	for _, d := range savedDepartment {
		responseDepartments = append(responseDepartments, &pb.Department{
			Id:   d.ID,
			Name: d.Name,
		})
	}

	return responseDepartments
}
