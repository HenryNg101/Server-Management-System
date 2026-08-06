import { useCallback, useEffect, useMemo, useState, type ReactNode } from "react";
import api, { tokenStore } from "../api/client";
import type { Role } from "../types/api";
import { AuthContext } from "./auth";

function roleFromToken(token: string | null): Role | null {
  if (!token) return null;
  try { return JSON.parse(atob(token.split(".")[1].replace(/-/g, "+").replace(/_/g, "/"))).role === "admin" ? "admin" : "user"; } catch { return null; }
}

export function AuthProvider({ children }: { children: ReactNode }) {
  const [role, setRole] = useState<Role | null>(() => roleFromToken(tokenStore.access));
  const login = useCallback(async (email: string, password: string) => {
    const { data } = await api.post<{ access_token: string; refresh_token: string }>("/auth/login", { email, password });
    tokenStore.set(data.access_token, data.refresh_token); setRole(roleFromToken(data.access_token));
  }, []);
  const logout = useCallback(async () => {
    try { if (tokenStore.refresh) await api.post("/auth/logout", { refresh_token: tokenStore.refresh }); } finally { tokenStore.clear(); setRole(null); }
  }, []);
  useEffect(() => { const handle = () => setRole(null); window.addEventListener("smsx:logout", handle); return () => window.removeEventListener("smsx:logout", handle); }, []);
  const value = useMemo(() => ({ role, isAuthenticated: Boolean(role), login, logout }), [role, login, logout]);
  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}
