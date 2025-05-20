package tests

import (
	authpb "health/protogen/v1/auth"
	employeepb "health/protogen/v1/users"
	"health/tests/suite"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

func TestEmployee_Get(t *testing.T) {
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

	email_1 := gofakeit.Email()
	pass_1 := randomFakePassword()

	respRegEmployee, err := st.AuthClient.Register(ctx, &authpb.AuthRequest{
		Email:    email_1,
		Password: pass_1,
		Source:   authpb.AppSource_EMPLOYEE,
	})

	require.NoError(t, err)
	require.NotNil(t, respRegEmployee)

	respVerifyTokenEmployee, err := st.AuthClient.VerifyRegister(ctx, &authpb.VerifyRegisterRequest{
		AuthToken: respRegEmployee.AuthToken,
	})
	require.NoError(t, err)
	require.NotNil(t, respVerifyTokenEmployee)

	respLogin, err := st.AuthClient.Login(ctx, &authpb.AuthRequest{
		Email:    email_1,
		Password: pass_1,
		Source:   authpb.AppSource_EMPLOYEE,
	})

	require.NotNil(t, respLogin)
	require.NoError(t, err)

	// Добавляем "Bearer <token>" в заголовки
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+respLogin.Token)

	getEmployee, err := st.EmployeeClient.GetEmployee(ctx, &emptypb.Empty{})
	require.NotNil(t, getEmployee)
	require.NoError(t, err)

	require.Equal(t, getEmployee.Email, email_1)
	require.Equal(t, getEmployee.Name, "")
	require.Equal(t, getEmployee.SecondName, "")
	require.Equal(t, getEmployee.Sex, false)
	require.Equal(t, getEmployee.Age, int64(18))
	require.Equal(t, getEmployee.Post, "")
	require.Equal(t, getEmployee.Surname, "")
}

func TestEmployeeUpdate_And_Delete(t *testing.T) {
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

	email_1 := gofakeit.Email()
	pass_1 := randomFakePassword()

	respRegEmployee, err := st.AuthClient.Register(ctx, &authpb.AuthRequest{
		Email:    email_1,
		Password: pass_1,
		Source:   authpb.AppSource_EMPLOYEE,
	})

	require.NoError(t, err)
	require.NotNil(t, respRegEmployee)

	respVerifyTokenEmployee, err := st.AuthClient.VerifyRegister(ctx, &authpb.VerifyRegisterRequest{
		AuthToken: respRegEmployee.AuthToken,
	})
	require.NoError(t, err)
	require.NotNil(t, respVerifyTokenEmployee)

	respLogin, err := st.AuthClient.Login(ctx, &authpb.AuthRequest{
		Email:    email_1,
		Password: pass_1,
		Source:   authpb.AppSource_EMPLOYEE,
	})

	require.NotNil(t, respLogin)
	require.NoError(t, err)

	// Добавляем "Bearer <token>" в заголовки
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+respLogin.Token)

	updateEmployee, err := st.EmployeeClient.UpdateEmployee(ctx, &employeepb.UpdateEmployeeRequest{
		Employee: &employeepb.Employee{
			Name:       "Sonya",
			SecondName: "Olegovna",
			Age:        30,
			Sex:        true,
			Post:       "Leader",
			Surname:    "Limonova",
			Department: "TestLaba",
		},
		FieldMask: &fieldmaskpb.FieldMask{
			Paths: []string{"name", "second_name", "age", "sex", "post", "surname", "department"},
		},
	})
	require.NotNil(t, updateEmployee)
	require.NoError(t, err)

	require.Equal(t, updateEmployee.Email, email_1)
	require.Equal(t, updateEmployee.Name, "Sonya")
	require.Equal(t, updateEmployee.Surname, "Limonova")
	require.Equal(t, updateEmployee.Post, "Leader")
	require.Equal(t, updateEmployee.Age, int64(30))
	require.Equal(t, updateEmployee.Sex, true)
	require.Equal(t, updateEmployee.Department, "TestLaba")

	// respDelete, err := st.EmployeeClient.DeleteEmployee(ctx, &emptypb.Empty{})
	// require.NotNil(t, respDelete)
	// require.NoError(t, err)
}

func TestListDepartments (t *testing.T) {
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

	email_1 := gofakeit.Email()
	pass_1 := randomFakePassword()

	respRegEmployee, err := st.AuthClient.Register(ctx, &authpb.AuthRequest{
		Email:    email_1,
		Password: pass_1,
		Source:   authpb.AppSource_EMPLOYEE,
	})

	require.NoError(t, err)
	require.NotNil(t, respRegEmployee)

	respVerifyTokenEmployee, err := st.AuthClient.VerifyRegister(ctx, &authpb.VerifyRegisterRequest{
		AuthToken: respRegEmployee.AuthToken,
	})
	require.NoError(t, err)
	require.NotNil(t, respVerifyTokenEmployee)

	respLogin, err := st.AuthClient.Login(ctx, &authpb.AuthRequest{
		Email:    email_1,
		Password: pass_1,
		Source:   authpb.AppSource_EMPLOYEE,
	})

	require.NotNil(t, respLogin)
	require.NoError(t, err)

	// Добавляем "Bearer <token>" в заголовки
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+respLogin.Token)

	listDeps, err := st.EmployeeClient.ListDepartments(ctx, &emptypb.Empty{})
	require.NotNil(t, listDeps)
	require.NoError(t, err)
}
