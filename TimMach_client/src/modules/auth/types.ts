import { User } from '../users/types';

export interface PasswordLoginRequest {
  username: string;
  password: string;
}

export interface PasswordLoginResponse {
  access_token: string;
  refresh_token?: string;
  token_type?: string;
  expires_in?: number;
  refresh_expires_in?: number;
  scope?: string;
  user?: User; // Keycloak không trả user; giữ optional cho tương thích cũ
}
