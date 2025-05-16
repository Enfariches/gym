import { ScheduleServiceClient } from '../../protogen/v1/schedule/schedule.client';
import type {
  CreateSchedulesRequest,
  DeleteScheduleRequest,
  GetScheduleRequest,
  Schedule,
  UpdateScheduleRequest
} from '../../protogen/v1/schedule/schedule';
import { GrpcWebFetchTransport } from '@protobuf-ts/grpcweb-transport';
import type { RpcOptions } from '@protobuf-ts/runtime-rpc';
import type { Empty } from '../../protogen/google/protobuf/empty';
import { FieldMask } from '../../protogen/google/protobuf/field_mask';

// Получаем базовый URL из переменных окружения или используем значение по умолчанию
const API_URL = process.env.QUASAR_API_URL || 'http://localhost:8085';

// Функция для создания транспорта с JWT токеном
const createTransport = () => {
  const token = localStorage.getItem('auth_token');

  return new GrpcWebFetchTransport({
    baseUrl: API_URL,
    headers: token ? {
      'Authorization': `Bearer ${token}`,
    } : undefined,
  });
};

// Функция для создания опций RPC с JWT токеном
const createAuthOptions = (): RpcOptions => {
  const token = localStorage.getItem('auth_token');

  if (!token) {
    throw new Error('Требуется авторизация для доступа к API расписания');
  }

  return {
    meta: {
      'Authorization': `Bearer ${token}`,
    }
  };
};

// Получить список расписаний
export const listSchedules = async (): Promise<Schedule[]> => {
  const transport = createTransport();
  const scheduleService = new ScheduleServiceClient(transport);

  const request: Empty = {};

  try {
    const options = createAuthOptions();
    const call = await scheduleService.listSchedule(request, options);
    return call.response.schedules;
  } catch (error) {
    console.error('Error fetching schedules:', error);
    throw new Error('Не удалось получить список расписаний');
  }
};

// Получить конкретное расписание по ID
export const getSchedule = async (scheduleId: string): Promise<Schedule> => {
  const transport = createTransport();
  const scheduleService = new ScheduleServiceClient(transport);

  const request: GetScheduleRequest = {
    scheduleId: BigInt(scheduleId)
  };

  try {
    const options = createAuthOptions();
    const call = await scheduleService.getSchedule(request, options);
    return call.response;
  } catch (error) {
    console.error('Error fetching schedule:', error);
    throw new Error('Не удалось получить расписание');
  }
};

// Создать новое расписание (сервис ожидает CreateSchedules)
export const createSchedules = async (schedules: Schedule[]): Promise<Schedule[]> => {
  const transport = createTransport();
  const scheduleService = new ScheduleServiceClient(transport);

  const request: CreateSchedulesRequest = {
    schedules: schedules,
  };

  try {
    const options = createAuthOptions();
    const call = await scheduleService.createSchedules(request, options);
    return call.response.schedules;
  } catch (error) {
    console.error('Error creating schedule(s): ', error);
    throw new Error('Не удалось создать расписание(я)');
  }
};

// Удалить расписание
export const deleteSchedule = async (scheduleId: string): Promise<void> => {
  const transport = createTransport();
  const scheduleService = new ScheduleServiceClient(transport);

  const request: DeleteScheduleRequest = {
    scheduleId: BigInt(scheduleId)
  };

  try {
    const options = createAuthOptions();
    await scheduleService.deleteSchedule(request, options);
  } catch (error) {
    console.error('Error deleting schedule:', error);
    throw new Error('Не удалось удалить расписание');
  }
};

// Обновить расписание
export const updateSchedule = async (schedule: Schedule, fieldPaths: string[]): Promise<Schedule> => {
  const transport = createTransport();
  const scheduleService = new ScheduleServiceClient(transport);

  const request: UpdateScheduleRequest = {
    schedule: schedule,
    fieldMask: FieldMask.create({ paths: fieldPaths })
  };

  try {
    const options = createAuthOptions();
    const call = await scheduleService.updateSchedule(request, options);
    return call.response;
  } catch (error) {
    console.error('Error updating schedule:', error);
    throw new Error('Не удалось обновить расписание');
  }
};
