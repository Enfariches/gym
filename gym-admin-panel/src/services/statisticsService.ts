import { StatisticsServiceClient } from '../../protogen/v1/statistics/statistics.client';
import type {
  CreateStatisticsRequest,
  GetEmployeeStatisticsRequest,
  ListEmployeeStatisticsRequest,
  ListMediaStatisticsRequest,
  Statistic,
  MediaProgress
} from '../../protogen/v1/statistics/statistics';
import { Empty } from '../../protogen/google/protobuf/empty';
import { GrpcWebFetchTransport } from '@protobuf-ts/grpcweb-transport';
import type { RpcOptions } from '@protobuf-ts/runtime-rpc';

// API URL configuration
const GRPC_URL = 'http://localhost:8085';
const HTTP_URL = 'http://localhost:3000';

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

/**
 * Создает запись о просмотре видео
 */
export const createStatistics = async (
  progress: MediaProgress,
  percentageView: number,
  mediaId: number
): Promise<void> => {
  const transport = createTransport();
  const client = new StatisticsServiceClient(transport);
  const request: CreateStatisticsRequest = {
    progress,
    percentageView: BigInt(percentageView),
    mediaId: BigInt(mediaId)
  };

  try {
    const options = createAuthOptions();
    await client.createStatistics(request, options);
  } catch (error) {
    console.error('Ошибка при создании статистики:', error);
    throw new Error('Не удалось сохранить данные о просмотре');
  }
};

/**
 * Получить статистику сотрудника по конкретному видео
 */
export const getEmployeeStatistics = async (
  employeeId: number,
  mediaId: number
): Promise<Statistic> => {
  const transport = createTransport();
  const client = new StatisticsServiceClient(transport);
  const request: GetEmployeeStatisticsRequest = {
    employeeId: BigInt(employeeId),
    mediaId: BigInt(mediaId)
  };

  try {
    const options = createAuthOptions();
    const { response } = await client.getEmployeeStatistics(request, options);
    return response;
  } catch (error) {
    console.error('Ошибка при получении статистики сотрудника:', error);
    throw new Error('Не удалось получить статистику сотрудника');
  }
};

/**
 * Получить статистику всех сотрудников по конкретному видео
 */
export const listMediaStatistics = async (
  mediaId: number
): Promise<Statistic[]> => {
  const transport = createTransport();
  const client = new StatisticsServiceClient(transport);
  const request: ListMediaStatisticsRequest = {
    mediaId: BigInt(mediaId)
  };

  try {
    const options = createAuthOptions();
    const { response } = await client.listMediaStatistics(request, options);
    return response.statistics || [];
  } catch (error) {
    console.error('Ошибка при получении статистики по видео:', error);
    throw new Error('Не удалось получить статистику по видео');
  }
};

/**
 * Получить статистику конкретного сотрудника по всем видео
 */
export const listEmployeeStatistics = async (
  employeeId: number
): Promise<Statistic[]> => {
  const transport = createTransport();
  const client = new StatisticsServiceClient(transport);
  const request: ListEmployeeStatisticsRequest = {
    employeeId: BigInt(employeeId)
  };

  try {
    const options = createAuthOptions();
    const { response } = await client.listEmployeeStatistics(request, options);
    return response.statistics || [];
  } catch (error) {
    console.error('Ошибка при получении статистики сотрудника:', error);
    throw new Error('Не удалось получить статистику сотрудника');
  }
};

/**
 * Получить статистику по отделу (по всем сотрудникам и всем видео)
 */
export const listDepartmentStatistics = async (): Promise<Statistic[]> => {
  const transport = createTransport();
  const client = new StatisticsServiceClient(transport);

  try {
    const options = createAuthOptions();
    const { response } = await client.listDepartmentStatistics(Empty.create(), options);
    return response.statistics || [];
  } catch (error) {
    console.error('Ошибка при получении статистики отдела:', error);
    throw new Error('Не удалось получить статистику отдела');
  }
};

/**
 * Экспорт статистики в PDF формат
 */
export const exportStatisticsToPDF = async (
  startDate?: Date | null,
  endDate?: Date | null,
  departmentId?: number | null
): Promise<Blob> => {
  try {
    const token = localStorage.getItem('auth_token');
    if (!token) throw new Error('Требуется авторизация для доступа к API');

    let url = `${HTTP_URL}/api/export`;
    const params = new URLSearchParams();

    if (startDate) {
      params.append('startDate', startDate.toISOString());
    }

    if (endDate) {
      params.append('endDate', endDate.toISOString());
    }

    if (departmentId) {
      params.append('departmentId', departmentId.toString());
    }

    if (params.toString()) {
      url += `?${params.toString()}`;
    }

    const response = await fetch(url, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${token}`
      }
    });

    if (!response.ok) {
      throw new Error(`HTTP error! Status: ${response.status}`);
    }

    return await response.blob();
  } catch (error) {
    console.error('Ошибка при экспорте статистики в PDF:', error);
    throw new Error('Не удалось экспортировать статистику в PDF');
  }
};
