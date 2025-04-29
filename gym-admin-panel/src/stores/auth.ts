import { defineStore } from 'pinia';
import { ref } from 'vue';
import { authService } from '../services/authService';
import type {
  RegisterResponse,
  ResetPasswordResponse,
  ChangePasswordResponse,
  VerifyRegisterResponse
} from '../types/auth';

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(localStorage.getItem('token'));
  const loading = ref(false);
  const error = ref<string | null>(null);

  const isAuthenticated = ref(!!token.value);

  const setToken = (newToken: string | null) => {
    token.value = newToken;
    isAuthenticated.value = !!newToken;
    if (newToken) {
      localStorage.setItem('token', newToken);
    } else {
      localStorage.removeItem('token');
    }
  };

  const register = async (email: string, password: string): Promise<RegisterResponse> => {
    try {
      loading.value = true;
      error.value = null;
      const response = await authService.register(email, password);
      return response;
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Registration failed';
      throw err;
    } finally {
      loading.value = false;
    }
  };

  const login = async (email: string, password: string): Promise<void> => {
    try {
      loading.value = true;
      error.value = null;
      const response = await authService.login(email, password);
      setToken(response.token);
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Login failed';
      throw err;
    } finally {
      loading.value = false;
    }
  };

  const logout = () => {
    setToken(null);
  };

  const resetPassword = async (email: string): Promise<ResetPasswordResponse> => {
    try {
      loading.value = true;
      error.value = null;
      const response = await authService.resetPassword(email);
      return response;
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Password reset failed';
      throw err;
    } finally {
      loading.value = false;
    }
  };

  const changePassword = async (resetToken: string, newPassword: string): Promise<ChangePasswordResponse> => {
    try {
      loading.value = true;
      error.value = null;
      const response = await authService.changePassword(resetToken, newPassword);
      return response;
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Password change failed';
      throw err;
    } finally {
      loading.value = false;
    }
  };

  const verifyRegister = async (authToken: string): Promise<VerifyRegisterResponse> => {
    try {
      loading.value = true;
      error.value = null;
      const response = await authService.verifyRegister(authToken);
      return response;
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Verification failed';
      throw err;
    } finally {
      loading.value = false;
    }
  };

  return {
    token,
    loading,
    error,
    isAuthenticated,
    register,
    login,
    logout,
    resetPassword,
    changePassword,
    verifyRegister
  };
});
