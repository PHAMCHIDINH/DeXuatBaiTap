import { create } from 'zustand';
import { me } from '../users/api';
import { User } from '../users/types';
import { passwordLogin } from './api';
import { PasswordLoginRequest } from './types';
import { clearToken, getToken, setToken } from '../../utils/storage';
import { setAuthToken } from '../../api/client';

type AuthState = {
  token: string | null;
  user: User | null;
  loading: boolean;
  authenticated: boolean;
  profileRequested: boolean;
  login: (payload: PasswordLoginRequest) => Promise<void>;
  logout: () => void;
  refreshProfile: (force?: boolean) => Promise<boolean>;
  initAuth: () => Promise<void>;
};

const initialToken = getToken();
if (initialToken) {
  setAuthToken(initialToken);
}

export const useAuthStore = create<AuthState>((set, get) => ({
  token: initialToken,
  user: null,
  loading: Boolean(initialToken),
  authenticated: Boolean(initialToken),
  profileRequested: false,

  login: async (payload: PasswordLoginRequest) => {
    set({ loading: true });
    try {
      const res = await passwordLogin(payload);
      setToken(res.access_token);
      setAuthToken(res.access_token);
      set({ token: res.access_token, authenticated: true, profileRequested: false });
      const loaded = await get().refreshProfile(true);
      if (!loaded) {
        throw new Error('Không thể tải thông tin người dùng');
      }
    } finally {
      set({ loading: false });
    }
  },

  logout: () => {
    clearToken();
    setAuthToken(null);
    set({ token: null, user: null, authenticated: false, loading: false, profileRequested: false });
  },

  refreshProfile: async (force = false) => {
    const token = get().token;
    if (!token) {
      set({ loading: false, user: null, authenticated: false, profileRequested: false });
      return false;
    }
    if (get().profileRequested && !force) {
      return Boolean(get().user);
    }
    try {
      set({ loading: true, profileRequested: true });
      const profile = await me();
      set({ user: profile, authenticated: true });
      return true;
    } catch (err) {
      get().logout();
      return false;
    } finally {
      set({ loading: false, profileRequested: false });
    }
  },

  initAuth: async () => {
    await get().refreshProfile();
  },
}));
