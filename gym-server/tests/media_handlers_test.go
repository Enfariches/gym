package tests

import (
	"health/lib/jwt"
	authpb "health/protogen/v1/auth"
	mediapb "health/protogen/v1/media"
	"health/tests/suite"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestMedia_Get(t *testing.T) {
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

	userId, departmentId, err := jwt.ParseToken(respLogin.Token)
	require.NoError(t, err)

	// Добавляем медиа заранее в БД, тк если указанного медиа нет, то расписание не добавится
	media_id_1, err := st.Storage.SaveMediaPostgres(userId, departmentId, "mediaTitle1")
	require.NoError(t, err)

	media_id_2, err := st.Storage.SaveMediaPostgres(userId, departmentId, "mediaTitle2")
	require.NoError(t, err)

	getMedia, err := st.MediaClient.GetMedia(ctx, &mediapb.GetMediaRequest{
		MediaId:      media_id_1,
		DepartmentId: departmentId,
		Expiry:       durationpb.New(5 * time.Minute),
	})
	require.NoError(t, err)
	require.NotNil(t, media_id_2)

	require.Equal(t, getMedia.Title, "mediaTitle1")
	require.NotNil(t, getMedia.PressignedUrl)
	require.Equal(t, getMedia.AdminId, userId)
}

func TestMedia_List_And_Delete(t *testing.T) {
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

	userId, departmentId, err := jwt.ParseToken(respLogin.Token)
	require.NoError(t, err)

	// Добавляем медиа заранее в БД, тк если указанного медиа нет, то расписание не добавится
	media_id_1, err := st.Storage.SaveMediaPostgres(userId, departmentId, "mediaTitle1")
	require.NoError(t, err)

	_, err = st.Storage.SaveMediaPostgres(userId, departmentId, "mediaTitle2")
	require.NoError(t, err)
	
	listMedia, err := st.MediaClient.ListMedia(ctx, &emptypb.Empty{})
	require.Len(t, listMedia.Medias, 2)
	require.NotNil(t, listMedia)
	require.NoError(t, err)

	media := listMedia.Medias[1]
	require.Equal(t, media.Title, "mediaTitle2")

	respDelete, err := st.MediaClient.DeleteMedia(ctx, &mediapb.DeleteMediaRequest{
		MediaId: media_id_1,
	})

	require.NotNil(t, respDelete)
	require.NoError(t, err)
}