import { createContext, useCallback, useContext, useEffect, useRef, useState } from 'react';
import { gearApi, GearItem, Station } from '@/src/api/gear';

interface GearContextValue {
  gear: GearItem[];
  stations: Station[];
  loading: boolean;
  refresh: () => Promise<void>;
  addGear: (item: GearItem) => void;
  updateGear: (item: GearItem) => void;
  removeGear: (id: string) => void;
  addStation: (station: Station) => void;
  updateStation: (station: Station) => void;
  removeStation: (id: string) => void;
}

const GearContext = createContext<GearContextValue | null>(null);

export function GearProvider({ children }: { children: React.ReactNode }) {
  const [gear, setGear] = useState<GearItem[]>([]);
  const [stations, setStations] = useState<Station[]>([]);
  const [loading, setLoading] = useState(true);
  const loaded = useRef(false);

  const refresh = useCallback(async () => {
    try {
      const [g, s] = await Promise.all([gearApi.listGear(), gearApi.listStations()]);
      setGear(g);
      setStations(s);
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

  const addGear = useCallback((item: GearItem) => {
    setGear(prev => [...prev, item]);
  }, []);

  const updateGear = useCallback((item: GearItem) => {
    setGear(prev => prev.map(g => (g.id === item.id ? item : g)));
    setStations(prev =>
      prev.map(s => ({
        ...s,
        gear: s.gear.map(g => (g.id === item.id ? item : g)),
      })),
    );
  }, []);

  const removeGear = useCallback((id: string) => {
    setGear(prev => prev.filter(g => g.id !== id));
    setStations(prev =>
      prev.map(s => ({ ...s, gear: s.gear.filter(g => g.id !== id) })),
    );
  }, []);

  const addStation = useCallback((station: Station) => {
    setStations(prev => [...prev, station]);
  }, []);

  const updateStation = useCallback((station: Station) => {
    setStations(prev => prev.map(s => (s.id === station.id ? station : s)));
  }, []);

  const removeStation = useCallback((id: string) => {
    setStations(prev => prev.filter(s => s.id !== id));
  }, []);

  return (
    <GearContext.Provider
      value={{
        gear,
        stations,
        loading,
        refresh,
        addGear,
        updateGear,
        removeGear,
        addStation,
        updateStation,
        removeStation,
      }}
    >
      {children}
    </GearContext.Provider>
  );
}

export function useGear() {
  const ctx = useContext(GearContext);
  if (!ctx) throw new Error('useGear must be used inside GearProvider');
  return ctx;
}
