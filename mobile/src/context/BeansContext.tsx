import { createContext, useCallback, useContext, useEffect, useRef, useState } from 'react';
import { beansApi, Bean } from '@/src/api/beans';

interface BeansContextValue {
  beans: Bean[];
  loading: boolean;
  refresh: () => Promise<void>;
  addBean: (bean: Bean) => void;
  updateBean: (bean: Bean) => void;
  removeBean: (id: string) => void;
}

const BeansContext = createContext<BeansContextValue | null>(null);

export function BeansProvider({ children }: { children: React.ReactNode }) {
  const [beans, setBeans] = useState<Bean[]>([]);
  const [loading, setLoading] = useState(true);
  const loaded = useRef(false);

  const refresh = useCallback(async () => {
    try {
      const result = await beansApi.list();
      setBeans(result);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    if (!loaded.current) {
      loaded.current = true;
      refresh();
    }
  }, [refresh]);

  const addBean    = useCallback((bean: Bean) => setBeans(prev => [...prev, bean]), []);
  const updateBean = useCallback((bean: Bean) => setBeans(prev => prev.map(b => b.id === bean.id ? bean : b)), []);
  const removeBean = useCallback((id: string) => setBeans(prev => prev.filter(b => b.id !== id)), []);

  return (
    <BeansContext.Provider value={{ beans, loading, refresh, addBean, updateBean, removeBean }}>
      {children}
    </BeansContext.Provider>
  );
}

export function useBeans() {
  const ctx = useContext(BeansContext);
  if (!ctx) throw new Error('useBeans must be used inside BeansProvider');
  return ctx;
}
