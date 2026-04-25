import { createContext, useCallback, useContext, useEffect, useState } from 'react';
import { router } from 'expo-router';
import { clearTokens } from '@/src/lib/auth';
import { authorizedFetch, setUnauthorizedHandler } from '@/src/lib/apiClient';

export interface User {
  id: string;
  first_name: string;
  last_name: string;
  email: string;
  created_at: string;
}

interface UserContextValue {
  user: User | null;
  loading: boolean;
  setUser: (user: User) => void;
  logout: () => Promise<void>;
  refreshUser: () => Promise<void>;
}

const UserContext = createContext<UserContextValue | null>(null);

export function UserProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);

  const logout = useCallback(async () => {
    await clearTokens();
    setUser(null);
    router.replace('/(auth)/login');
  }, []);

  const refreshUser = useCallback(async () => {
    try {
      const u = await authorizedFetch<User>('/user/me');
      setUser(u);
    } catch {
      // 401s are handled by the unauthorized handler above
    }
  }, []);

  useEffect(() => {
    setUnauthorizedHandler(logout);
  }, [logout]);

  useEffect(() => {
    refreshUser().finally(() => setLoading(false));
  }, []);

  return (
    <UserContext.Provider value={{ user, loading, setUser, logout, refreshUser }}>
      {children}
    </UserContext.Provider>
  );
}

export function useUser(): UserContextValue {
  const ctx = useContext(UserContext);
  if (!ctx) throw new Error('useUser must be used inside UserProvider');
  return ctx;
}
