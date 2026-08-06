import axios, { AxiosError, type InternalAxiosRequestConfig } from "axios";

const TOKEN_KEY = "smsx.access_token";
const REFRESH_KEY = "smsx.refresh_token";

export const tokenStore = {
  get access() { return localStorage.getItem(TOKEN_KEY); },
  get refresh() { return localStorage.getItem(REFRESH_KEY); },
  set(access: string, refresh: string) {
    localStorage.setItem(TOKEN_KEY, access);
    localStorage.setItem(REFRESH_KEY, refresh);
  },
  clear() { localStorage.removeItem(TOKEN_KEY); localStorage.removeItem(REFRESH_KEY); },
};

const api = axios.create({ baseURL: import.meta.env.VITE_API_URL ?? "http://localhost", withCredentials: true });

api.interceptors.request.use((config: InternalAxiosRequestConfig) => {
  const token = tokenStore.access;
  if (token) config.headers.Authorization = `Bearer ${token}`;
  return config;
});

let refreshPromise: Promise<string> | null = null;
api.interceptors.response.use(undefined, async (error: AxiosError) => {
  const request = error.config as (InternalAxiosRequestConfig & { _retried?: boolean }) | undefined;
  if (!request || error.response?.status !== 401 || request._retried || request.url?.includes("/auth/")) throw error;
  const refreshToken = tokenStore.refresh;
  if (!refreshToken) throw error;
  request._retried = true;
  try {
    refreshPromise ??= axios.post<{ access_token: string; refresh_token: string }>(`${api.defaults.baseURL}/auth/refresh`, { refresh_token: refreshToken }, { withCredentials: true })
      .then(({ data }) => { tokenStore.set(data.access_token, data.refresh_token); return data.access_token; })
      .finally(() => { refreshPromise = null; });
    const access = await refreshPromise;
    request.headers.Authorization = `Bearer ${access}`;
    return api(request);
  } catch (refreshError) {
    tokenStore.clear();
    window.dispatchEvent(new Event("smsx:logout"));
    throw refreshError;
  }
});

export default api;
