import type { HttpMethod, FetchOptions, FetchParams, Api_Route } from './types';

export class HttpError extends Error {
  public status: number;
  public response: Response;
  constructor(message: string, status: number, response: Response) {
    super(message);
    this.status = status;
    this.response = response;
  }
}

function buildQuery(params?: FetchParams): string {
  if (!params) return '';
  const usp = new URLSearchParams();
  Object.entries(params).forEach(([k, v]) => {
    if (v !== undefined && v !== null) usp.append(k, String(v));
  });
  const s = usp.toString();
  return s ? '?' + s : '';
}

export const client = {
  async request<T = unknown>(
    url: Api_Route,
    method: HttpMethod,
    opts: FetchOptions = {}
  ): Promise<T> {
    const fullUrl = url + buildQuery(opts.params);

    const headers: Record<string, string> = { ...opts.headers };
    let body: BodyInit | undefined = undefined;

    if (opts.body && !(opts.body instanceof FormData)) {
      // If object, encode as JSON
      if (typeof opts.body === 'object' && !(opts.body instanceof Blob)) {
        body = JSON.stringify(opts.body);
        headers['Content-Type'] = 'application/json';
      } else if (typeof opts.body === 'string' || opts.body instanceof Blob) {
        body = opts.body as BodyInit;
        if (typeof opts.body === 'string') {
          headers['Content-Type'] = headers['Content-Type'] || 'text/plain';
        }
      }
    } else if (opts.body instanceof FormData) {
      body = opts.body;
      // Content-Type set automatically
    }

    const res = await fetch(fullUrl, {
      ...opts,
      method,
      headers,
      body: method === 'GET' ? undefined : body,
    });

    if (!res.ok) {
      const msg = await res.text();
      throw new HttpError(msg || res.statusText, res.status, res);
    }

    if (opts.parseJson === false) {
      // Return raw Response object
      return res as unknown as T;
    }

    // Default: parse as JSON (if not 204)
    return res.status === 204 ? (undefined as T) : await res.json();
  },

  get<T = unknown>(url: Api_Route, opts?: FetchOptions) {
    return client.request<T>(url, 'GET', opts);
  },
  post<T = unknown>(
    url: Api_Route,
    body?: BodyInit | object,
    opts?: FetchOptions
  ) {
    return client.request<T>(url, 'POST', { ...opts, body });
  },
  put<T = unknown>(
    url: Api_Route,
    body?: BodyInit | object,
    opts?: FetchOptions
  ) {
    return client.request<T>(url, 'PUT', { ...opts, body });
  },
  delete<T = unknown>(url: Api_Route, opts?: FetchOptions) {
    return client.request<T>(url, 'DELETE', opts);
  },
  patch<T = unknown>(
    url: Api_Route,
    body?: BodyInit | object,
    opts?: FetchOptions
  ) {
    return client.request<T>(url, 'PATCH', { ...opts, body });
  },
};
