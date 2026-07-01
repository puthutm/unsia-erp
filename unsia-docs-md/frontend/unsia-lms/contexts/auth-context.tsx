"use client";

import { createContext, useContext, useState, useEffect, ReactNode } from "react";

// Types
export interface UserInfo {
  user_id: string;
  person_id: string;
  name: string;
  email: string;
  active_role: string;
  permissions: string[];
  scope: string;
}

interface AuthContextType {
  user: UserInfo | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  login: (username: string, password: string) => Promise<boolean>;
  logout: () => void;
  switchRole: (role: string) => Promise<void>;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

// API Configuration
const API_BASE_URL = process.env.NEXT_PUBLIC_AUTH_API || "http://localhost:8001";
const STORAGE_KEYS = {
  accessToken: "unsia_access_token",
  refreshToken: "unsia_refresh_token",
  user: "unsia_user",
};

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<UserInfo | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  const getToken = () => localStorage.getItem(STORAGE_KEYS.accessToken);

  const fetchUserInfo = async (token: string) => {
    try {
      const response = await fetch(`${API_BASE_URL}/api/v1/auth/me`, {
        headers: {
          Authorization: `Bearer ${token}`,
          "Content-Type": "application/json",
        },
      });

      if (response.ok) {
        const data = await response.json();
        return data.data;
      }
      return null;
    } catch {
      return null;
    }
  };

  useEffect(() => {
    const initAuth = async () => {
      let token = getToken();

      // Check query params for SSO token injection
      if (typeof window !== "undefined") {
        const searchParams = new URLSearchParams(window.location.search);
        const urlToken = searchParams.get("token");
        const urlRefreshToken = searchParams.get("refresh_token");

        if (urlToken && urlRefreshToken) {
          token = urlToken;
          localStorage.setItem(STORAGE_KEYS.accessToken, urlToken);
          localStorage.setItem(STORAGE_KEYS.refreshToken, urlRefreshToken);
          document.cookie = `${STORAGE_KEYS.accessToken}=${urlToken}; path=/; max-age=604800; SameSite=Lax`;
          document.cookie = `${STORAGE_KEYS.refreshToken}=${urlRefreshToken}; path=/; max-age=604800; SameSite=Lax`;
          
          // Clear query params from URL
          const cleanUrl = window.location.protocol + "//" + window.location.host + window.location.pathname;
          window.history.replaceState({ path: cleanUrl }, "", cleanUrl);
        }
      }

      if (token) {
        const userData = await fetchUserInfo(token);
        if (userData) {
          setUser(userData);
        } else {
          clearAuthStorage();
          redirectToLogin();
        }
      } else {
        redirectToLogin();
      }
      setIsLoading(false);
    };

    const redirectToLogin = () => {
      if (typeof window !== "undefined") {
        const currentUrl = window.location.protocol + "//" + window.location.host + window.location.pathname;
        const portalUrl = process.env.NEXT_PUBLIC_PORTAL_URL || "http://localhost:3010";
        window.location.href = `${portalUrl}/login?redirect=${encodeURIComponent(currentUrl)}`;
      }
    };

    const clearAuthStorage = () => {
      localStorage.removeItem(STORAGE_KEYS.accessToken);
      localStorage.removeItem(STORAGE_KEYS.refreshToken);
      localStorage.removeItem(STORAGE_KEYS.user);
    };

    initAuth();
  }, []);

  const login = async (username: string, password: string): Promise<boolean> => {
    setIsLoading(true);
    try {
      const response = await fetch(`${API_BASE_URL}/api/v1/auth/login`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ username, password }),
      });

      if (!response.ok) throw new Error("Login failed");

      const data = await response.json();

      localStorage.setItem(STORAGE_KEYS.accessToken, data.data.access_token);
      localStorage.setItem(STORAGE_KEYS.refreshToken, data.data.refresh_token);
      localStorage.setItem(STORAGE_KEYS.user, JSON.stringify(data.data.user));

      setUser(data.data.user);
      return true;
    } catch {
      return false;
    } finally {
      setIsLoading(false);
    }
  };

  const logout = () => {
    localStorage.removeItem(STORAGE_KEYS.accessToken);
    localStorage.removeItem(STORAGE_KEYS.refreshToken);
    localStorage.removeItem(STORAGE_KEYS.user);
    setUser(null);
  };

  const switchRole = async (role: string) => {
    const token = getToken();
    if (!token) return;

    try {
      const response = await fetch(`${API_BASE_URL}/api/v1/auth/switch-role`, {
        method: "POST",
        headers: {
          Authorization: `Bearer ${token}`,
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ role }),
      });

      if (response.ok) {
        const data = await response.json();
        setUser(data.data);
        localStorage.setItem(STORAGE_KEYS.user, JSON.stringify(data.data));
      }
    } catch {
      console.error("Failed to switch role");
    }
  };

  const value: AuthContextType = {
    user,
    isAuthenticated: !!user,
    isLoading,
    login,
    logout,
    switchRole,
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
}
