import { createContext, useCallback, useContext, useEffect, useState } from 'react';
import { getAccessToken } from '@/src/lib/auth';

interface AuthContextValue {
  isAuthenticated: boolean;
  ready: boolean;
  setIsAuthenticated: (v: boolean) => void;
}

const AuthContext = createContext<AuthContextValue | null>(null);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [ready, setReady] = useState(false);

  const setAuth = useCallback((v: boolean) => setIsAuthenticated(v), []);

  useEffect(() => {
    getAccessToken().then(token => {
      setIsAuthenticated(!!token);
      setReady(true);
    });
  }, []);

  return (
    <AuthContext.Provider value={{ isAuthenticated, ready, setIsAuthenticated: setAuth }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth(): AuthContextValue {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error('useAuth must be inside AuthProvider');
  return ctx;
}