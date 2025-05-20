package tests

import (
	"health/lib/jwt"
	authpb "health/protogen/v1/auth"
	statspb "health/protogen/v1/statistics"
	"health/tests/suite"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
)

func TestStatistics_Create(t *testing.T) {
	ctx, st := suite.New(t)

	// Добавляем медиа заранее в БД, тк если указанного медиа нет, то расписание не добавится
	media_id, err := st.Storage.SaveMediaPostgres(1, 1, "mediaTitle1")
	require.NoError(t, err)

	email := gofakeit.Email()
	pass := randomFakePassword()

	respReg, err := st.AuthClient.Register(ctx, &authpb.AuthRequest{
		Email:    email,
		Password: pass,
		Source:   authpb.AppSource_EMPLOYEE,
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
		Source:   authpb.AppSource_EMPLOYEE,
	})
	require.NoError(t, err)
	// Добавляем "Bearer <token>" в заголовки
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+respLogin.Token)

	stat1, err := st.StatisticsClient.CreateStatistics(ctx, &statspb.CreateStatisticsRequest{
		MediaId:        media_id,
		Progress:       statspb.MediaProgress_COMPLETED,
		PercentageView: 100,
	})

	require.NoError(t, err)
	require.NotNil(t, stat1)

	stat2, err := st.StatisticsClient.CreateStatistics(ctx, &statspb.CreateStatisticsRequest{
		MediaId:        media_id,
		Progress:       statspb.MediaProgress_INCOMPLETE,
		PercentageView: 20,
	})

	require.NoError(t, err)
	require.NotNil(t, stat2)

	stat3, err := st.StatisticsClient.CreateStatistics(ctx, &statspb.CreateStatisticsRequest{
		MediaId:        media_id,
		Progress:       statspb.MediaProgress_SKIPPED,
		PercentageView: 0,
	})

	require.NoError(t, err)
	require.NotNil(t, stat3)
}

func TestStatistics_Get(t *testing.T) {
	ctx, st := suite.New(t)

	// Добавляем медиа заранее в БД, тк если указанного медиа нет, то расписание не добавится
	media_id, err := st.Storage.SaveMediaPostgres(1, 1, "mediaTitle1")
	require.NoError(t, err)

	email := gofakeit.Email()
	pass := randomFakePassword()

	respReg, err := st.AuthClient.Register(ctx, &authpb.AuthRequest{
		Email:    email,
		Password: pass,
		Source:   authpb.AppSource_EMPLOYEE,
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
		Source:   authpb.AppSource_EMPLOYEE,
	})
	require.NoError(t, err)

	// Добавляем "Bearer <token>" в заголовки
	ctxE := metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+respLogin.Token)

	stat1, err := st.StatisticsClient.CreateStatistics(ctxE, &statspb.CreateStatisticsRequest{
		MediaId:        media_id,
		Progress:       statspb.MediaProgress_COMPLETED,
		PercentageView: 100,
	})

	require.NoError(t, err)
	require.NotNil(t, stat1)

	// Добавляем "Bearer <token>" в заголовки
	ctxE = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+respLogin.Token)
	userId, _, err := jwt.ParseToken(respLogin.Token)
	require.NoError(t, err)

	getStat, err := st.StatisticsClient.GetEmployeeStatistics(ctxE, &statspb.GetEmployeeStatisticsRequest{
		EmployeeId: userId,
		MediaId:    media_id,
	})

	require.NoError(t, err)
	require.NotNil(t, getStat)

	require.Equal(t, getStat.EmployeeName, "")
	require.Equal(t, getStat.MediaTitle, "mediaTitle1")

	require.Equal(t, getStat.PercentageView, int64(100))
	require.Equal(t, getStat.Progress, statspb.MediaProgress_COMPLETED)

}

func TestStatistics_ListMediaStatistics(t *testing.T) {
	ctx, st := suite.New(t)

	email := gofakeit.Email()
	pass := randomFakePassword()

	_, err := st.AuthClient.Register(ctx, &authpb.AuthRequest{
		Email:    email,
		Password: pass,
		Source:   authpb.AppSource_ADMIN,
	})
	require.NoError(t, err)

}

func TestStatistics_ListEmployeeStatistics(t *testing.T) {
	ctx, st := suite.New(t)

	email := gofakeit.Email()
	pass := randomFakePassword()

	_, err := st.AuthClient.Register(ctx, &authpb.AuthRequest{
		Email:    email,
		Password: pass,
		Source:   authpb.AppSource_EMPLOYEE,
	})
	require.NoError(t, err)

}

func TestStatistics_ListDepartmentStatistics(t *testing.T) {
	ctx, st := suite.New(t)

	email := gofakeit.Email()
	pass := randomFakePassword()

	_, err := st.AuthClient.Register(ctx, &authpb.AuthRequest{
		Email:    email,
		Password: pass,
		Source:   authpb.AppSource_EMPLOYEE,
	})
	require.NoError(t, err)
}
