import axios from 'axios';
import { PasswordLoginRequest, PasswordLoginResponse } from './types';

const KC_BASE = (import.meta.env.VITE_KEYCLOAK_URL as string | undefined) ?? 'http://localhost:8081';
const KC_REALM = (import.meta.env.VITE_KEYCLOAK_REALM as string | undefined) ?? 'timmach';
const KC_CLIENT_ID = (import.meta.env.VITE_KEYCLOAK_CLIENT_ID as string | undefined) ?? 'timmach-webapp';

// Gọi trực tiếp token endpoint của Keycloak bằng password grant.
export async function passwordLogin(payload: PasswordLoginRequest): Promise<PasswordLoginResponse> {
  try {
    const url = `${KC_BASE.replace(/\/$/, '')}/realms/${KC_REALM}/protocol/openid-connect/token`;
    const body = new URLSearchParams({
      client_id: KC_CLIENT_ID,
      grant_type: 'password',
      username: payload.username,
      password: payload.password,
    });
    const { data } = await axios.post<PasswordLoginResponse>(url, body, {
      headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
    });
    console.log('Token response:', data);
    return data;
  } catch (error) {
    if (axios.isAxiosError(error)) {
      console.error('Lỗi khi gọi Keycloak:', error.response?.data);
    } else {
      console.error('Lỗi không xác định:', error);
    }
    throw error;
  }
}
