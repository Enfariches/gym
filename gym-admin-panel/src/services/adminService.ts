import { AdminServiceClient } from '../../protogen/v1/users/admin.client';
import type {
  UpdateAdminRequest,
  Admin
} from '../../protogen/v1/users/admin';
import { Empty } from '../../protogen/google/protobuf/empty';

// Конфигурация gRPC-Web клиента
import { GrpcWebFetchTransport } from '@protobuf-ts/grpcweb-transport';
import type { RpcOptions } from '@protobuf-ts/runtime-rpc';

// Получаем базовый URL из переменных окружения или используем значение по умолчанию
const API_URL = 'http://localhost:8085';

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
    throw new Error('Требуется авторизация для доступа к API администратора');
  }

  return {
    meta: {
      'Authorization': `Bearer ${token}`,
    }
  };
};

// Получить информацию об администраторе
export const getAdmin = async (): Promise<Admin> => {
  const transport = createTransport();
  const adminService = new AdminServiceClient(transport);

  try {
    const options = createAuthOptions();
    const call = await adminService.getAdmin(Empty.create(), options);
    return call.response;
  } catch (error) {
    console.error('Error fetching admin data:', error);
    throw new Error('Не удалось получить данные администратора');
  }
};

// Обновить данные администратора
export const updateAdmin = async (
  admin: Admin,
  fieldsToUpdate: string[]
): Promise<Admin> => {
  const transport = createTransport();
  const adminService = new AdminServiceClient(transport);

  const request: UpdateAdminRequest = {
    admin,
    fieldMask: {
      paths: fieldsToUpdate
    }
  };

  try {
    const options = createAuthOptions();
    const call = await adminService.updateAdmin(request, options);
    return call.response;
  } catch (error) {
    console.error('Error updating admin data:', error);
    throw new Error('Не удалось обновить данные администратора');
  }
};
