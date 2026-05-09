import appConfig from '../../appConfig.json';

export const API_URL: string = appConfig.server_url;

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

  const text = await res.text();
  const json = text ? (JSON.parse(text) as unknown) : undefined;

  if (!res.ok) {
    throw new ApiResponseError(res.status, (json as ApiError)?.error ?? `HTTP ${res.status}`);
  }

  return json as Res;
}
