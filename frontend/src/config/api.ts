import Constants from 'expo-constants';

export const API_URL: string =
  (Constants.expoConfig?.extra?.apiUrl as string | undefined) ??
  'http://localhost:8080';

type Method = 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE';

interface RequestOptions<B> {
  method?: Method;
  body?: B;
  token?: string;
  signal?: AbortSignal;
}

interface ApiError {
  error: string;
}

export class ApiResponseError extends Error {
  constructor(public readonly status: number, message: string) {
    super(message);
    this.name = 'ApiResponseError';
  }
}

export async function apiFetch<Res, Body = unknown>(
  path: string,
  options: RequestOptions<Body> = {}
): Promise<Res> {
  const { method = 'GET', body, token, signal } = options;

  const headers: HeadersInit = { 'Content-Type': 'application/json' };
  if (token) headers['Authorization'] = `Bearer ${token}`;

  const res = await fetch(`${API_URL}${path}`, {
    method,
    headers,
    body: body !== undefined ? JSON.stringify(body) : undefined,
    signal,
  });

  const json = await res.json();

  if (!res.ok) {
    throw new ApiResponseError(res.status, (json as ApiError).error ?? `HTTP ${res.status}`);
  }

  return json as Res;
}
