export class ApiError extends Error {
  code: string;
  status: number;

  constructor(code: string, message: string, status: number) {
    super(message);
    this.code = code;
    this.status = status;
  }
}

interface Envelope<T> {
  data?: T;
  error?: { code: string; message: string };
}

function getCsrfToken(): string {
  const match = document.cookie.match(/(?:^|;\s*)toefl_csrf=([^;]+)/);
  return match ? decodeURIComponent(match[1]) : '';
}

const BASE = '/api/v1';

async function request<T>(path: string, init: RequestInit = {}): Promise<T> {
  const headers = new Headers(init.headers);
  headers.set('Accept', 'application/json');
  const isMutation = init.method && init.method !== 'GET' && init.method !== 'HEAD';
  if (isMutation) {
    headers.set('Content-Type', 'application/json');
    const csrf = getCsrfToken();
    if (csrf) headers.set('X-CSRF-Token', csrf);
  }

  const res = await fetch(`${BASE}${path}`, {
    ...init,
    headers,
    credentials: 'include',
  });

  if (res.status === 204) return undefined as T;

  let body: Envelope<T> | null = null;
  const text = await res.text();
  if (text) {
    try {
      body = JSON.parse(text) as Envelope<T>;
    } catch {
      body = null;
    }
  }

  if (!res.ok) {
    const err = body?.error;
    throw new ApiError(
      err?.code ?? 'internal',
      err?.message ?? `Request failed (${res.status})`,
      res.status,
    );
  }

  return body?.data as T;
}

export const api = {
  get: <T>(path: string) => request<T>(path),
  post: <T>(path: string, body?: unknown) =>
    request<T>(path, {
      method: 'POST',
      body: body !== undefined ? JSON.stringify(body) : undefined,
    }),
  put: <T>(path: string, body?: unknown) =>
    request<T>(path, {
      method: 'PUT',
      body: body !== undefined ? JSON.stringify(body) : undefined,
    }),
  del: <T>(path: string) => request<T>(path, { method: 'DELETE' }),
};