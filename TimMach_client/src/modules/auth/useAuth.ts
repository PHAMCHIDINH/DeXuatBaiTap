import { useEffect } from 'react';
import { useAuthStore } from './store';

// Hook tiện dụng để sử dụng store Zustand và tự động refresh profile khi đã có token.
export function useAuth() {
  const { token, user, refreshProfile } = useAuthStore((state) => ({
    token: state.token,
    user: state.user,
    refreshProfile: state.refreshProfile,
  }));

  useEffect(() => {
    if (token && !user) {
      refreshProfile();
    }
  }, [refreshProfile, token, user]);

  return useAuthStore();
}
