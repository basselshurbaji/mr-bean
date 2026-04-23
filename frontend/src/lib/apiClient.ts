import { apiFetch, ApiResponseError } from '@/src/config/api';
import { getAccessToken, getRefreshToken, saveTokens, clearTokens } from '@/src/lib/auth';

export class UnauthorizedError extends Error {
  constructor() {
    super('Unauthorized');
    this.name = 'UnauthorizedError';
  }
}

type Method = 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE';

interface AuthorizedRequestOptions<B> {
  method?: Method;
  body?: B;
  signal?: AbortSignal;
}

let _unauthorizedHandler: (() => void) | null = null;
let _refreshPromise: Promise<string> | null = null;

export function setUnauthorizedHandler(handler: () => void) {
  _unauthorizedHandler = handler;
}

async function refreshAccessToken(): Promise<string> {
  if (_refreshPromise) return _refreshPromise;

  _refreshPromise = (async () => {
    const refreshToken = await getRefreshToken();
    if (!refreshToken) throw new UnauthorizedError();

    const res = await apiFetch<{ token: string; refresh_token: string }>('/auth/refresh', {
      method: 'POST',
      body: { refresh_token: refreshToken },
    });
    await saveTokens(res.token, res.refresh_token);
    return res.token;
  })()
    .catch(async (e) => {
      await clearTokens();
      _unauthorizedHandler?.();
      throw e instanceof UnauthorizedError ? e : new UnauthorizedError();
    })
    .finally(() => {
      _refreshPromise = null;
    });

  return _refreshPromise;
}

export async function authorizedFetch<Res, Body = unknown>(
  path: string,
  options: AuthorizedRequestOptions<Body> = {}
): Promise<Res> {
  let token = await getAccessToken();

  if (!token) {
    token = await refreshAccessToken();
  }

  try {
    return await apiFetch<Res, Body>(path, { ...options, token });
  } catch (err) {
    if (!(err instanceof ApiResponseError) || err.status !== 401) throw err;
    const newToken = await refreshAccessToken();
    return apiFetch<Res, Body>(path, { ...options, token: newToken });
  }
}
