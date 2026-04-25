import { authorizedFetch } from '@/src/lib/apiClient';

export interface GearItem {
  id: string;
  type_id: string;
  name: string;
  brand?: string | null;
  model?: string | null;
  year?: string | null;
  notes?: string | null;
  created_at: string;
  updated_at: string;
}

export interface Station {
  id: string;
  name: string;
  gear: GearItem[];
  created_at: string;
  updated_at: string;
}

export interface CreateGearBody {
  type_id: string;
  name: string;
  brand?: string;
  model?: string;
  year?: string;
  notes?: string;
}

export interface CreateStationBody {
  name: string;
  gear_ids: string[];
}

export const gearApi = {
  listGear: () => authorizedFetch<GearItem[]>('/gear'),
  createGear: (body: CreateGearBody) =>
    authorizedFetch<GearItem>('/gear', { method: 'POST', body }),
  updateGear: (id: string, body: CreateGearBody) =>
    authorizedFetch<GearItem>(`/gear/${id}`, { method: 'PUT', body }),
  deleteGear: (id: string) =>
    authorizedFetch<void>(`/gear/${id}`, { method: 'DELETE' }),

  listStations: () => authorizedFetch<Station[]>('/stations'),
  createStation: (body: CreateStationBody) =>
    authorizedFetch<Station>('/stations', { method: 'POST', body }),
  updateStation: (id: string, body: CreateStationBody) =>
    authorizedFetch<Station>(`/stations/${id}`, { method: 'PUT', body }),
  deleteStation: (id: string) =>
    authorizedFetch<void>(`/stations/${id}`, { method: 'DELETE' }),
};
