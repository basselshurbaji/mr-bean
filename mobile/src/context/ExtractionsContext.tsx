import { createContext, useCallback, useContext, useEffect, useRef, useState } from 'react';
import { extractionApi, Extraction } from '@/src/api/extractions';

interface ExtractionsContextValue {
  extractions: Extraction[];
  loading: boolean;
  refresh: () => Promise<void>;
  addExtraction: (e: Extraction) => void;
  removeExtraction: (id: string) => void;
}

const ExtractionsContext = createContext<ExtractionsContextValue | null>(null);

export function ExtractionsProvider({ children }: { children: React.ReactNode }) {
  const [extractions, setExtractions] = useState<Extraction[]>([]);
  const [loading, setLoading] = useState(true);
  const loaded = useRef(false);

  const refresh = useCallback(async () => {
    try {
      const result = await extractionApi.list({ limit: 20 });
      setExtractions(result);
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

  const addExtraction = useCallback(
    (e: Extraction) => setExtractions(prev => [e, ...prev]),
    [],
  );

  const removeExtraction = useCallback(
    (id: string) => setExtractions(prev => prev.filter(e => e.id !== id)),
    [],
  );

  return (
    <ExtractionsContext.Provider value={{ extractions, loading, refresh, addExtraction, removeExtraction }}>
      {children}
    </ExtractionsContext.Provider>
  );
}

export function useExtractions() {
  const ctx = useContext(ExtractionsContext);
  if (!ctx) throw new Error('useExtractions must be used inside ExtractionsProvider');
  return ctx;
}
