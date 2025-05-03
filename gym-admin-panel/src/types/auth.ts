export enum AppSource {
  APP_SOURCE_UNSPECIFIED = 0,
  EMPLOYEE = 1,
  ADMIN = 2
}

export interface AuthRequest {
  email: string;
  password: string;
  source: AppSource;
}

export interface ResetPasswordRequest {
  email: string;
  source: AppSource;
}

export interface ChangePasswordRequest {
  reset_token: string;
  new_password: string;
}

export interface VerifyRegisterRequest {
  auth_token: string;
}

export interface LoginResponse {
  token: string;
}

export interface ResetPasswordResponse {
  reset_token: string;
}

export interface ChangePasswordResponse {
  status: 'OK' | 'ERROR';
}

export interface VerifyRegisterResponse {
  status: 'OK' | 'ERROR';
}
