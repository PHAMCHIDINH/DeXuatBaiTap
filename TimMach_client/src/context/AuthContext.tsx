import { createContext, ReactNode, useCallback, useEffect, useState } from 'react';
import * as usersApi from '../api/users';
import { LoginRequest, RegisterRequest, User } from '../types/api';
import { clearToken, getToken as readToken, setToken as storeToken } from '../utils/storage';
import { setAuthToken } from '../api/client';

interface AuthContextValue {
  user: User | null;
  token: string | null;
  loading: boolean;
  login: (payload: LoginRequest) => Promise<void>;
  register: (payload: RegisterRequest) => Promise<void>;
  logout: () => void;
  refreshProfile: () => Promise<void>;
}

export const AuthContext = createContext<AuthContextValue | undefined>(undefined);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [token, setToken] = useState<string | null>(() => readToken());
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState<boolean>(!!readToken());

  const logout = useCallback(() => {
    clearToken();
    setAuthToken(null);
    setToken(null);
    setUser(null);
    setLoading(false);
  }, []);

  const refreshProfile = useCallback(async () => {
    if (!token) return;
    try {
      const me = await usersApi.me();
      setUser(me);
    } catch {
      logout();
    } finally {
      setLoading(false);
    }
  }, [logout, token]);

  useEffect(() => {
    setAuthToken(token);
    if (token) {
      refreshProfile();
    } else {
      setLoading(false);
    }
  }, [logout, refreshProfile, token]);

  const login = useCallback(async (payload: LoginRequest) => {
    const res = await usersApi.login(payload);
    storeToken(res.access_token);
    setAuthToken(res.access_token);
    setToken(res.access_token);
    setUser(res.user);
  }, []);

  const register = useCallback(async (payload: RegisterRequest) => {
    const res = await usersApi.register(payload);
    storeToken(res.token);
    setAuthToken(res.token);
    setToken(res.token);
    setUser(res.user);
  }, []);

  const value: AuthContextValue = {
    user,
    token,
    loading,
    login,
    register,
    logout,
    refreshProfile,
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}
