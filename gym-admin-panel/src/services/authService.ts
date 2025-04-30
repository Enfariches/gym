import axios from 'axios';
import type {
  AuthRequest,
  LoginResponse,
  RegisterResponse,
  ResetPasswordRequest,
  ResetPasswordResponse,
  ChangePasswordRequest,
  ChangePasswordResponse,
  VerifyRegisterRequest,
  VerifyRegisterResponse} from '../types/auth';
import {
  AppSource
} from '../types/auth';

const ENVOY_URL = process.env.VITE_ENVOY_URL || 'http://localhost:8085';

class AuthService {
  private client = axios.create({
    baseURL: ENVOY_URL,
    headers: {
      'Content-Type': 'application/json',
      'Accept': 'application/json'
    }
  });

  async register(email: string, password: string): Promise<RegisterResponse> {
    const request: AuthRequest = {
      email,
      password,
      source: AppSource.ADMIN
    };

    const response = await this.client.post('/auth.AuthService/Register', request);
    return response.data;
  }

  async login(email: string, password: string): Promise<LoginResponse> {
    const request: AuthRequest = {
      email,
      password,
      source: AppSource.ADMIN
    };

    const response = await this.client.post('/auth.AuthService/Login', request);
    return response.data;
  }

  async resetPassword(email: string): Promise<ResetPasswordResponse> {
    const request: ResetPasswordRequest = {
      email,
      source: AppSource.ADMIN
    };

    const response = await this.client.post('/auth.AuthService/ResetPassword', request);
    return response.data;
  }

  async changePassword(resetToken: string, newPassword: string): Promise<ChangePasswordResponse> {
    const request: ChangePasswordRequest = {
      reset_token: resetToken,
      new_password: newPassword
    };

    const response = await this.client.post('/auth.AuthService/ChangePassword', request);
    return response.data;
  }

  async verifyRegister(authToken: string): Promise<VerifyRegisterResponse> {
    const request: VerifyRegisterRequest = {
      auth_token: authToken
    };

    const response = await this.client.post('/auth.AuthService/VerifyRegister', request);
    return response.data;
  }
}

export const authService = new AuthService();
