/* eslint-disable @typescript-eslint/no-unused-vars */
import { defineStore } from 'pinia';
import {
  register,
  login,
  resetPassword,
  changePassword,
  verifyRegister
} from '../services/authService';
import type {
  RegisterResponse,
  LoginResponse,
  ChangePasswordResponse,
} from '../../protogen/v1/auth/auth';
import { VerifyRegisterResponse } from 'src/types/auth';

export const useAuthStore = defineStore('auth', {
  state: () => ({
    isAuthenticated: false,
    token: null as string | null,
  }),

  actions: {
    // Регистрация
    async register(email: string, password: string): Promise<RegisterResponse> {
      try {
        const response = await register(email, password);
        return response;
      } catch (error) {
        throw new Error('Registration failed');
      }
    },

    // Аутентификация
    async login(email: string, password: string): Promise<LoginResponse> {
      try {
        const response = await login(email, password);
        this.isAuthenticated = true;
        this.token = response.token;
        localStorage.setItem('auth_token', this.token!); // Сохраняем токен
        return response;
      } catch (error) {
        throw new Error('Login failed');
      }
    },

    // Сброс пароля
    async resetPassword(email: string): Promise<ChangePasswordResponse> {
      try {
        const response = await resetPassword(email);
        return response;
      } catch (error) {
        throw new Error('Reset password failed');
      }
    },

    // Подтверждение сброса пароля
    async changePassword(resetToken: string, newPassword: string) {
      try {
        const response = await changePassword(resetToken, newPassword);
      } catch (error) {
        throw new Error('Change password failed');
      }
    },

    // Подтверждение регистрации
    async verifyRegister(authToken: string) {
      try {
        const response = await verifyRegister(authToken);
        this.isAuthenticated = true;
      } catch (error) {
        throw new Error('Verification failed');
      }
    },

    // Выход из системы
    logout() {
      this.isAuthenticated = false;
      this.token = null;
      localStorage.removeItem('auth_token');
    },
  },
});
