import { authorizedFetch } from '@/src/lib/apiClient';
import { palette } from '@/src/theme';

export interface Bean {
  id: string;
  user_id: string;
  name: string;
  roaster?: string | null;
  origin?: string | null;
  process?: string | null;
  roast_level?: string | null;
  tasting_notes?: string | null;
  notes?: string | null;
  created_at: string;
  updated_at: string;
}

export interface BeanBody {
  name: string;
  roaster?: string;
  origin?: string;
  process?: string;
  roast_level?: string;
  tasting_notes?: string;
  notes?: string;
}

export const PROCESSES = [
  { id: 'washed',    label: 'Washed'    },
  { id: 'natural',   label: 'Natural'   },
  { id: 'honey',     label: 'Honey'     },
  { id: 'anaerobic', label: 'Anaerobic' },
  { id: 'other',     label: 'Other'     },
] as const;

export const ROAST_LEVELS = [
  { id: 'light',       label: 'Light',        color: palette.caramel400  },
  { id: 'medium_light',label: 'Medium light',  color: palette.caramel600  },
  { id: 'medium',      label: 'Medium',        color: palette.espresso500 },
  { id: 'medium_dark', label: 'Medium dark',   color: palette.espresso600 },
  { id: 'dark',        label: 'Dark',          color: palette.espresso800 },
] as const;

export function processLabel(id: string | null | undefined): string {
  return PROCESSES.find(p => p.id === id)?.label ?? '';
}

export function roastLabel(id: string | null | undefined): string {
  return ROAST_LEVELS.find(r => r.id === id)?.label ?? '';
}

export function roastColor(id: string | null | undefined): string {
  return ROAST_LEVELS.find(r => r.id === id)?.color ?? palette.cream500;
}

export const beansApi = {
  list:   ()                          => authorizedFetch<Bean[]>('/beans'),
  create: (body: BeanBody)            => authorizedFetch<Bean>('/beans',         { method: 'POST', body }),
  update: (id: string, body: BeanBody)=> authorizedFetch<Bean>(`/beans/${id}`,   { method: 'PUT',  body }),
  delete: (id: string)                => authorizedFetch<void>(`/beans/${id}`,   { method: 'DELETE' }),
};
