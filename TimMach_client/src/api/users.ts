import client from './client';
import { LoginRequest, LoginResponse, RegisterRequest, RegisterResponse, User } from '../types/api';

export async function login(payload: LoginRequest): Promise<LoginResponse> {
  const { data } = await client.post<LoginResponse>('/users/login', payload);
  return data;
}

export async function register(payload: RegisterRequest): Promise<RegisterResponse> {
  const { data } = await client.post<RegisterResponse>('/users/register', payload);
  return data;
}

export async function me(): Promise<User> {
  const { data } = await client.get<User>('/users/me');
  return data;
}
