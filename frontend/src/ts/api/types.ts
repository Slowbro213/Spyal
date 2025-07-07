export type Api_Route = '/create/remote' | '/join';

export type HttpMethod = 'GET' | 'POST' | 'PUT' | 'DELETE' | 'PATCH';

export type FetchParams = Record<
  string,
  string | number | boolean | undefined | null
>;

export type FetchOptions = Omit<RequestInit, 'method' | 'body'> & {
  params?: FetchParams;
  headers?: Record<string, string>;
  body?: BodyInit | object;
  parseJson?: boolean;
};
