import { authorizedFetch } from '@/src/lib/apiClient';

export interface BeanSummary {
  id: string;
  name: string;
  roaster: string | null;
  roast: string | null;
}

export interface GearSummary {
  id: string;
  type_id: string;
  name: string;
}

export interface Extraction {
  id: string;
  user_id: string;
  bean: BeanSummary;
  dose_in: number;
  yield_out: number;
  time: number;
  target_time: number;
  grind_size: number;
  pre_infusion: boolean;
  tasting_note: string | null;
  gear: GearSummary[];
  created_at: string;
  updated_at: string;
}

export interface CreateExtractionBody {
  bean_id: string;
  dose_in: number;
  yield_out: number;
  time: number;
  target_time: number;
  grind_size: number;
  gear_ids: string[];
  pre_infusion: boolean;
  tasting_note?: string | null;
}

export type ExtractionZone = 'under' | 'perfect' | 'over';

export function computeZone(time: number, targetTime: number): ExtractionZone {
  if (time < targetTime - 4) return 'under';
  if (time > targetTime + 4) return 'over';
  return 'perfect';
}

export const extractionApi = {
  list: (params?: { limit?: number; page?: number }) => {
    const q = new URLSearchParams();
    if (params?.limit != null) q.set('limit', String(params.limit));
    if (params?.page != null) q.set('page', String(params.page));
    const qs = q.toString();
    return authorizedFetch<Extraction[]>(`/extractions${qs ? `?${qs}` : ''}`);
  },
  get: (id: string) => authorizedFetch<Extraction>(`/extractions/${id}`),
  create: (body: CreateExtractionBody) =>
    authorizedFetch<Extraction>('/extractions', { method: 'POST', body }),
  update: (id: string, body: CreateExtractionBody) =>
    authorizedFetch<Extraction>(`/extractions/${id}`, { method: 'PUT', body }),
  delete: (id: string) =>
    authorizedFetch<void>(`/extractions/${id}`, { method: 'DELETE' }),
};
