import { AuthServiceClient } from '../../protogen/v1/auth/auth.client';
import type {
  AuthRequest,
  RegisterResponse,
  LoginResponse,
  ChangePasswordRequest,
  ChangePasswordResponse,
  VerifyRegisterRequest,
  VerifyChangePasswordRequest,
} from '../../protogen/v1/auth/auth';
import {
  AppSource
} from '../../protogen/v1/auth/auth';

// Конфигурация gRPC-Web клиента
import { GrpcWebFetchTransport } from '@protobuf-ts/grpcweb-transport';

const transport = new GrpcWebFetchTransport({
  baseUrl: 'http://localhost:8085' // Замените на адрес вашего gRPC-сервера
});

const authService = new AuthServiceClient(transport);

// Регистрация пользователя
export const register = async (email: string, password: string): Promise<RegisterResponse> => {
  const request: AuthRequest = {
    email,
    password,
    source: AppSource.ADMIN
  };
  return (await authService.register(request)).response;
};

// Аутентификация пользователя
export const login = async (email: string, password: string): Promise<LoginResponse> => {
  const request: AuthRequest = {
    email,
    password,
    source: AppSource.ADMIN
  };
  return (await authService.login(request)).response;
};

// Сброс пароля (отправка reset token)
export const resetPassword = async (email: string): Promise<ChangePasswordResponse> => {
  const request: ChangePasswordRequest = {
    email,
    source: AppSource.ADMIN
  };
  return (await authService.changePassword(request)).response;
};

// Подтверждение сброса пароля (смена пароля)
export const changePassword = async (resetToken: string, newPassword: string) => {
  const request: VerifyChangePasswordRequest = {
    resetToken,
    newPassword
  };
  return (await authService.verifyChangePassword(request)).response;
};

// Подтверждение регистрации
export const verifyRegister = async (authToken: string) => {
  const request: VerifyRegisterRequest = {
    authToken
  };
  return (await authService.verifyRegister(request)).response;
};
