import { ScheduleServiceClient } from '../../protogen/v1/schedule/schedule.client';
import type { Schedule, CreateSchedulesRequest, UpdateScheduleRequest, DeleteScheduleRequest } from '../../protogen/v1/schedule/schedule';
import { Empty } from '../../protogen/google/protobuf/empty';
import { GrpcWebFetchTransport } from '@protobuf-ts/grpcweb-transport';
import type { RpcOptions } from '@protobuf-ts/runtime-rpc';

const GRPC_URL = 'http://localhost:8085'; // Замените на ваш адрес

const createTransport = () => {
  const token = localStorage.getItem('auth_token');
  return new GrpcWebFetchTransport({
    baseUrl: GRPC_URL,
    headers: token ? { 'Authorization': `Bearer ${token}` } : undefined,
  });
};

const createAuthOptions = (): RpcOptions => {
  const token = localStorage.getItem('auth_token');
  if (!token) throw new Error('Требуется авторизация для доступа к API');
  return { meta: { 'Authorization': `Bearer ${token}` } };
};

export const listSchedulesGrpc = async (): Promise<Schedule[]> => {
  const transport = createTransport();
  const client = new ScheduleServiceClient(transport);
  try {
    const options = createAuthOptions();
    const call = await client.listSchedule(Empty.create(), options);
    return call.response.schedules ?? [];
  } catch (error) {
    console.error('Ошибка при получении расписаний (gRPC):', error);
    throw new Error('Не удалось получить расписания');
  }
};

export const createSchedulesGrpc = async (schedules: Schedule[]): Promise<Schedule[]> => {
  const transport = createTransport();
  const client = new ScheduleServiceClient(transport);
  const request: CreateSchedulesRequest = { schedules };
  try {
    const options = createAuthOptions();
    const call = await client.createSchedules(request, options);
    return call.response.schedules ?? [];
  } catch (error) {
    console.error('Ошибка при создании расписания (gRPC):', error);
    throw new Error('Не удалось создать расписание');
  }
};

export const updateScheduleGrpc = async (schedule: Schedule, fieldMask?: string[]): Promise<Schedule> => {
  const transport = createTransport();
  const client = new ScheduleServiceClient(transport);
  const request: UpdateScheduleRequest = fieldMask
    ? { schedule, fieldMask: { paths: fieldMask } }
    : { schedule };
  try {
    const options = createAuthOptions();
    const call = await client.updateSchedule(request, options);
    return call.response;
  } catch (error) {
    console.error('Ошибка при обновлении расписания (gRPC):', error);
    throw new Error('Не удалось обновить расписание');
  }
};

export const deleteScheduleGrpc = async (scheduleId: bigint): Promise<void> => {
  const transport = createTransport();
  const client = new ScheduleServiceClient(transport);
  const request: DeleteScheduleRequest = { scheduleId };
  try {
    const options = createAuthOptions();
    await client.deleteSchedule(request, options);
  } catch (error) {
    console.error('Ошибка при удалении расписания (gRPC):', error);
    throw new Error('Не удалось удалить расписание');
  }
};
