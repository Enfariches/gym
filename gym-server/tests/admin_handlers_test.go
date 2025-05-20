package tests

import (
	authpb "health/protogen/v1/auth"
	adminpb "health/protogen/v1/users"
	"health/tests/suite"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

func TestAdmin_Get(t *testing.T) {
	ctx, st := suite.New(t)

	email := gofakeit.Email()
	pass := randomFakePassword()

	respReg, err := st.AuthClient.Register(ctx, &authpb.AuthRequest{
		Email:    email,
		Password: pass,
		Source:   authpb.AppSource_ADMIN,
	})
	require.NoError(t, err)

	respVerifyToken, err := st.AuthClient.VerifyRegister(ctx, &authpb.VerifyRegisterRequest{
		AuthToken: respReg.AuthToken,
	})
	require.NoError(t, err)
	require.NotNil(t, respVerifyToken)

	respLogin, err := st.AuthClient.Login(ctx, &authpb.AuthRequest{
		Email:    email,
		Password: pass,
		Source:   authpb.AppSource_ADMIN,
	})
	require.NoError(t, err)

	// Добавляем "Bearer <token>" в заголовки
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+respLogin.Token)

	getAdmin, err := st.AdminClient.GetAdmin(ctx, &emptypb.Empty{})
	require.NotNil(t, getAdmin)
	require.NoError(t, err)

	require.Equal(t, getAdmin.Email, email)
	require.Equal(t, getAdmin.Name, "")
	require.Equal(t, getAdmin.Surname, "")
	require.Equal(t, getAdmin.Department, "")

}

func TestAdmin_Update(t *testing.T) {
	name := "Testik"
	surname := "Testov"
	department := "TestLaba"

	ctx, st := suite.New(t)

	email := gofakeit.Email()
	pass := randomFakePassword()

	respReg, err := st.AuthClient.Register(ctx, &authpb.AuthRequest{
		Email:    email,
		Password: pass,
		Source:   authpb.AppSource_ADMIN,
	})
	require.NoError(t, err)

	respVerifyToken, err := st.AuthClient.VerifyRegister(ctx, &authpb.VerifyRegisterRequest{
		AuthToken: respReg.AuthToken,
	})
	require.NoError(t, err)
	require.NotNil(t, respVerifyToken)

	respLogin, err := st.AuthClient.Login(ctx, &authpb.AuthRequest{
		Email:    email,
		Password: pass,
		Source:   authpb.AppSource_ADMIN,
	})
	require.NoError(t, err)

	// Добавляем "Bearer <token>" в заголовки
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+respLogin.Token)

	updateAdmin, err := st.AdminClient.UpdateAdmin(ctx, &adminpb.UpdateAdminRequest{
		Admin: &adminpb.Admin{
			Name:       name,
			Surname:    surname,
			Department: "TestLaba",
		},
		FieldMask: &fieldmaskpb.FieldMask{
			Paths: []string{"name", "surname", "department"},
		},
	})
	require.NotNil(t, updateAdmin)
	require.NoError(t, err)

	require.Equal(t, updateAdmin.Email, email)
	require.Equal(t, updateAdmin.Name, name)
	require.Equal(t, updateAdmin.Surname, surname)
	require.Equal(t, updateAdmin.Department, department)

	// Проверяем админа еще через GetAdmin
	getAdmin, err := st.AdminClient.GetAdmin(ctx, &emptypb.Empty{})
	require.NotNil(t, getAdmin)
	require.NoError(t, err)

	require.Equal(t, updateAdmin.Email, email)
	require.Equal(t, updateAdmin.Name, name)
	require.Equal(t, updateAdmin.Surname, surname)
	require.Equal(t, updateAdmin.Department, department)
}

func TestAdmin_ListEmployees(t *testing.T) {
	ctx, st := suite.New(t)

	email := gofakeit.Email()
	pass := randomFakePassword()

	respReg, err := st.AuthClient.Register(ctx, &authpb.AuthRequest{
		Email:    email,
		Password: pass,
		Source:   authpb.AppSource_ADMIN,
	})
	require.NoError(t, err)

	respVerifyToken, err := st.AuthClient.VerifyRegister(ctx, &authpb.VerifyRegisterRequest{
		AuthToken: respReg.AuthToken,
	})
	require.NoError(t, err)
	require.NotNil(t, respVerifyToken)

	respLogin, err := st.AuthClient.Login(ctx, &authpb.AuthRequest{
		Email:    email,
		Password: pass,
		Source:   authpb.AppSource_ADMIN,
	})
	require.NoError(t, err)

	// Добавляем "Bearer <token>" в заголовки
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+respLogin.Token)

	listEmployees, err := st.AdminClient.ListAdminEmployees(ctx, &emptypb.Empty{})

	require.NoError(t, err)
	require.Len(t, listEmployees.Employees, 0)
	require.NotNil(t, listEmployees)
}
