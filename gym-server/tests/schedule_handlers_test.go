package tests

import (
	"health/lib/jwt"
	authpb "health/protogen/v1/auth"
	schedulepb "health/protogen/v1/schedule"
	"health/tests/suite"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

func TestSchedule_Create(t *testing.T) {
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
	media_id, err := st.Storage.SaveMediaPostgres(userId, departmentId, "mediaTitle")
	require.NoError(t, err)

	schedulesCreated, err := st.ScheduleClient.CreateSchedules(ctx, &schedulepb.CreateSchedulesRequest{
		Schedules: []*schedulepb.Schedule{
			{
				CronExpression: "30 * * * *",
				IsActive:       true,
				MediaId:        media_id,
			},
			{
				CronExpression: "0 12 * * 1", // Понедельник в 12:00
				IsActive:       false,
				MediaId:        media_id,
			},
			{
				CronExpression: "1 12 * * 1", // Понедельник в 12:01
				IsActive:       true,
				MediaId:        media_id,
			},
		},
	})
	require.NoError(t, err)
	require.Len(t, schedulesCreated.Schedules, 3)

	s := schedulesCreated.Schedules[0]
	require.NotZero(t, s.Id)
	require.NotZero(t, s.AdminId)
	require.NotEmpty(t, s.CreatedAt)
	require.Equal(t, "30 * * * *", s.CronExpression)
}

func TestSchedule_Get_And_Update(t *testing.T) {
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
	media_id, err := st.Storage.SaveMediaPostgres(userId, departmentId, "mediaTitle")
	require.NoError(t, err)

	schedulesCreated, err := st.ScheduleClient.CreateSchedules(ctx, &schedulepb.CreateSchedulesRequest{
		Schedules: []*schedulepb.Schedule{
			{
				CronExpression: "* * * * *",
				IsActive:       true,
				MediaId:        media_id,
			},
		},
	})
	require.NoError(t, err)
	require.Len(t, schedulesCreated.Schedules, 1)
	schedule_id := schedulesCreated.Schedules[0].Id

	scheduleGetm, err := st.ScheduleClient.GetSchedule(ctx, &schedulepb.GetScheduleRequest{
		ScheduleId: schedule_id,
	})
	require.NoError(t, err)
	require.NotZero(t, scheduleGetm.Id)
	require.NotZero(t, scheduleGetm.AdminId)
	require.NotEmpty(t, scheduleGetm.CreatedAt)
	require.Equal(t, "* * * * *", scheduleGetm.CronExpression)

	updatedSchedule, err := st.ScheduleClient.UpdateSchedule(ctx, &schedulepb.UpdateScheduleRequest{
		Schedule: &schedulepb.Schedule{
			Id:     scheduleGetm.Id,
			CronExpression: "1 1 1 1 1",
			IsActive:       false,
		},
		FieldMask: &fieldmaskpb.FieldMask{
			Paths: []string{"cron_expression", "is_active"},
		},
	})

	require.NotNil(t, updatedSchedule)
	require.NoError(t, err)

	require.Equal(t, updatedSchedule.CronExpression, "1 1 1 1 1")
	require.Equal(t, updatedSchedule.IsActive, false)

}

func TestSchedule_List_And_Delete(t *testing.T) {
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
	media_id, err := st.Storage.SaveMediaPostgres(userId, departmentId, "mediaTitle")
	require.NoError(t, err)

	schedulesCreated, err := st.ScheduleClient.CreateSchedules(ctx, &schedulepb.CreateSchedulesRequest{
		Schedules: []*schedulepb.Schedule{
			{
				CronExpression: "* * * * *",
				IsActive:       true,
				MediaId:        media_id,
			},
			{
				CronExpression: "1 * * * *",
				IsActive:       false,
				MediaId:        media_id,
			},
			{
				CronExpression: "* * 3 * *",
				IsActive:       true,
				MediaId:        media_id,
			},
			{
				CronExpression: "* 2 * * *",
				IsActive:       true,
				MediaId:        media_id,
			},
		},
	})
	require.NoError(t, err)
	require.Len(t, schedulesCreated.Schedules, 4)

	listSchedule, err := st.ScheduleClient.ListSchedule(ctx, &emptypb.Empty{})

	require.Len(t, listSchedule.Schedules, 4)
	require.NoError(t, err)
}
