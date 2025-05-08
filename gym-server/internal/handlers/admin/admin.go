package admin

import (
	"context"
	"health/internal/domain/models"
	pb "health/protogen/v1/users"

	"github.com/mennanov/fmutils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

type AdminService interface {
	GetAdmin(ctx context.Context, admin_id int64) (*models.Admin, error)
	UpdateAdmin(ctx context.Context, fieldMask map[string]interface{}) (*models.Admin, error)
}

type AdminServerManagmentApi struct {
	pb.UnimplementedAdminServiceServer
	admin AdminService
}

func RegisterGRPCServer(gRPC *grpc.Server, admin AdminService) {
	pb.RegisterAdminServiceServer(gRPC, &AdminServerManagmentApi{admin: admin})
}

func (s *AdminServerManagmentApi) GetAdmin(ctx context.Context, r *pb.GetAdminRequest) (*pb.Admin, error) {
	if r.AdminId == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "admin_id is required")
	}

	admin, err := s.admin.GetAdmin(ctx, r.AdminId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get admin: %v", err)
	}

	departament := ""
	if admin.Departament != nil {
		departament = *admin.Departament
	}

	name := ""
	if admin.Name != nil {
		name = *admin.Name
	}

	surname := ""
	if admin.Surname != nil {
		surname = *admin.Surname
	}

	return &pb.Admin{
		Id:          admin.ID,
		Name:        name,
		Surname:     surname,
		Email:       admin.Email,
		Departament: departament,
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
		return nil, status.Errorf(codes.Internal, "failed to update admin")
	}

	name := ""
	if updatedAdmin.Name != nil {
		name = *updatedAdmin.Name
	}

	surname := ""
	if updatedAdmin.Surname != nil {
		surname = *updatedAdmin.Surname
	}

	departament := ""
	if updatedAdmin.Departament != nil {
		departament = *updatedAdmin.Departament
	}

	return &pb.Admin{
		Id:          updatedAdmin.ID,
		Name:        name,
		Surname:     surname,
		Email:       updatedAdmin.Email,
		Departament: departament,
	}, nil
}

func applyFieldMask(req *pb.Admin, mask *fieldmaskpb.FieldMask) (map[string]interface{}, error) {
	updateMap := make(map[string]interface{})

	for _, path := range mask.GetPaths() {
		switch path {
		case "name":
			updateMap["name"] = req.Name
		case "surname":
			updateMap["surname"] = req.Surname
		case "departament":
			updateMap["departament"] = req.Departament
		default:
			return nil, status.Errorf(codes.InvalidArgument, "id and email immutable fields")
		}
	}
	return updateMap, nil
}
