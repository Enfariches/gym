package tests

import (
	"health/lib/jwt"
	pb "health/protogen/v1/auth"
	"health/tests/suite"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegister_VerifyRegister_Login_ChangePassword_Admin(t *testing.T) {
	ctx, st := suite.New(t)

	email := gofakeit.Email()
	pass := randomFakePassword()
	source := "admins"

	respReg, err := st.AuthClient.Register(ctx, &pb.AuthRequest{
		Email:    email,
		Password: pass,
		Source:   pb.AppSource_ADMIN,
	})

	require.NoError(t, err)
	assert.NotEmpty(t, respReg.GetAuthToken())

	u, err := jwt.ParseAuthToken(respReg.AuthToken)
	require.NoError(t, err)

	assert.Equal(t, email, u.Email)
	assert.Equal(t, source, u.Source)

	respVerifyToken, err := st.AuthClient.VerifyRegister(ctx, &pb.VerifyRegisterRequest{
		AuthToken: respReg.AuthToken,
	})

	require.NoError(t, err)
	require.NotNil(t, respVerifyToken)

	respLogin, err := st.AuthClient.Login(ctx, &pb.AuthRequest{
		Email:    email,
		Password: pass,
		Source:   pb.AppSource_ADMIN,
	})

	require.NotNil(t, respLogin)
	require.NoError(t, err)

	_, _, err = jwt.ParseToken(respLogin.Token)
	require.NoError(t, err)

	resetPass, err := st.AuthClient.ChangePassword(ctx, &pb.ChangePasswordRequest{
		Email:  email,
		Source: pb.AppSource_ADMIN,
	})

	require.NotNil(t, resetPass)
	require.NoError(t, err)

	u, err = jwt.ParseResetToken(resetPass.ResetToken)
	require.NoError(t, err)

	assert.Equal(t, email, u.Email)
	assert.Equal(t, source, u.Source)

	verifyResetPass, err := st.AuthClient.VerifyChangePassword(ctx, &pb.VerifyChangePasswordRequest{
		ResetToken:  resetPass.ResetToken,
		NewPassword: randomFakePassword(),
	})

	require.NoError(t, err)
	require.NotNil(t, verifyResetPass)
}

func TestRegister_VerifyRegister_Login_ChangePassword_Employee(t *testing.T) {
	ctx, st := suite.New(t)

	email := gofakeit.Email()
	pass := randomFakePassword()
	source := "employees"

	respReg, err := st.AuthClient.Register(ctx, &pb.AuthRequest{
		Email:    email,
		Password: pass,
		Source:   pb.AppSource_EMPLOYEE,
	})

	require.NoError(t, err)
	assert.NotEmpty(t, respReg.GetAuthToken())

	u, err := jwt.ParseAuthToken(respReg.AuthToken)
	require.NoError(t, err)

	assert.Equal(t, email, u.Email)
	assert.Equal(t, source, u.Source)

	respVerifyToken, err := st.AuthClient.VerifyRegister(ctx, &pb.VerifyRegisterRequest{
		AuthToken: respReg.AuthToken,
	})

	require.NoError(t, err)
	require.NotNil(t, respVerifyToken)

	respLogin, err := st.AuthClient.Login(ctx, &pb.AuthRequest{
		Email:    email,
		Password: pass,
		Source:   pb.AppSource_EMPLOYEE,
	})

	require.NotNil(t, respLogin)
	require.NoError(t, err)

	_, _, err = jwt.ParseToken(respLogin.Token)
	require.NoError(t, err)

	resetPass, err := st.AuthClient.ChangePassword(ctx, &pb.ChangePasswordRequest{
		Email:  email,
		Source: pb.AppSource_EMPLOYEE,
	})

	require.NotNil(t, resetPass)
	require.NoError(t, err)

	u, err = jwt.ParseResetToken(resetPass.ResetToken)
	require.NoError(t, err)

	assert.Equal(t, email, u.Email)
	assert.Equal(t, source, u.Source)

	verifyResetPass, err := st.AuthClient.VerifyChangePassword(ctx, &pb.VerifyChangePasswordRequest{
		ResetToken:  resetPass.ResetToken,
		NewPassword: randomFakePassword(),
	})

	require.NoError(t, err)
	require.NotNil(t, verifyResetPass)
}

func randomFakePassword() string {
	return gofakeit.Password(true, true, true, true, false, 10)
}
