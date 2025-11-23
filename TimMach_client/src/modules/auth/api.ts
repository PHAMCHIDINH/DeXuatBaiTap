import client from '../../api/client';
import { LoginRequest, LoginResponse, RegisterRequest, RegisterResponse } from '../../types/api';

export async function login(payload: LoginRequest): Promise<LoginResponse> {
  const { data } = await client.post<LoginResponse>('/users/login', payload);
  return data;
}

export async function register(payload: RegisterRequest): Promise<RegisterResponse> {
  const { data } = await client.post<RegisterResponse>('/users/register', payload);
  return data;
}
