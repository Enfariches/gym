import { AuthServiceClient } from '../../protogen/v1/auth/AuthServiceClientPb';
import type {
  RegisterResponse,
  LoginResponse,
  ResetPasswordResponse,
  ChangePasswordResponse,
  VerifyRegisterResponse} from '../../protogen/v1/auth/auth_pb';
import {
  AuthRequest,
  ResetPasswordRequest,
  ChangePasswordRequest,
  VerifyRegisterRequest
} from '../../protogen/v1/auth/auth_pb';

// Конфигурация gRPC-Web клиента
const authService = new AuthServiceClient(
  'http://localhost:8085', // Замените на адрес вашего gRPC-сервера
  null,
  null
);

// Регистрация пользователя
export const register = async (email: string, password: string): Promise<RegisterResponse> => {
  const request = new AuthRequest();
  request.setEmail(email);
  request.setPassword(password);
  request.setSource(2); // ADMIN (AppSource.ADMIN)

  return new Promise((resolve, reject) => {
    authService.register(request, {}, (err, response) => {
      if (err) {
        console.error('gRPC register error:', err);
        reject(err);
        return;
      }
      resolve(response);
    });
  });
};

// Аутентификация пользователя
export const login = async (email: string, password: string): Promise<LoginResponse> => {
  const request = new AuthRequest();
  request.setEmail(email);
  request.setPassword(password);
  request.setSource(2); // ADMIN (AppSource.ADMIN)

  return new Promise((resolve, reject) => {
    authService.login(request, {}, (err, response) => {
      if (err) {
        console.error('gRPC login error:', err);
        reject(err);
        return;
      }
      resolve(response);
    });
  });
};

// Сброс пароля
export const resetPassword = async (email: string): Promise<ResetPasswordResponse> => {
  const request = new ResetPasswordRequest();
  request.setEmail(email);
  request.setSource(2); // ADMIN (AppSource.ADMIN)

  return new Promise((resolve, reject) => {
    authService.resetPassword(request, {}, (err, response) => {
      if (err) {
        console.error('gRPC resetPassword error:', err);
        reject(err);
        return;
      }
      resolve(response);
    });
  });
};

// Подтверждение сброса пароля
export const changePassword = async (resetToken: string, newPassword: string): Promise<ChangePasswordResponse> => {
  const request = new ChangePasswordRequest();
  request.setResetToken(resetToken);
  request.setNewPassword(newPassword);

  return new Promise((resolve, reject) => {
    authService.changePassword(request, {}, (err, response) => {
      if (err) {
        console.error('gRPC changePassword error:', err);
        reject(err);
        return;
      }
      resolve(response);
    });
  });
};

// Подтверждение регистрации
export const verifyRegister = async (authToken: string): Promise<VerifyRegisterResponse> => {
  const request = new VerifyRegisterRequest();
  request.setAuthToken(authToken);

  return new Promise((resolve, reject) => {
    authService.verifyRegister(request, {}, (err, response) => {
      if (err) {
        console.error('gRPC verifyRegister error:', err);
        reject(err);
        return;
      }
      resolve(response);
    });
  });
};