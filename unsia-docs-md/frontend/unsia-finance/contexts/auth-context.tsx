"use client";

import { createContext, useContext, useState, useEffect, ReactNode } from "react";
import { useRouter } from "next/navigation";
import { API_BASE_URLS, AUTH_ENDPOINTS, STORAGE_KEYS, type TokenResponse, type UserInfo } from "@/lib/constants";

interface User {
  id: string;
  personId: string;
  email: string;
  name: string;
  role: string;
  permissions: string[];
  scope: string;
  avatar?: string;
}

interface AuthContextType {
  user: User | null;
  isLoading: boolean;
  login: (email: string, password: string) => Promise<void>;
  logout: () => void;
  refreshToken: () => Promise<void>;
  isAuthenticated: boolean;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: ReactNode }) {
  const router = useRouter();
  const [user, setUser] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [, setError] = useState("");

  // Check for existing session on mount
  useEffect(() => {
    const initAuth = async () => {
      let accessToken = localStorage.getItem(STORAGE_KEYS.accessToken);
      let refreshTokenVal = localStorage.getItem(STORAGE_KEYS.refreshToken);

      // Check query params for SSO token injection
      if (typeof window !== "undefined") {
        const searchParams = new URLSearchParams(window.location.search);
        const urlToken = searchParams.get("token");
        const urlRefreshToken = searchParams.get("refresh_token");

        if (urlToken && urlRefreshToken) {
          accessToken = urlToken;
          refreshTokenVal = urlRefreshToken;
          localStorage.setItem(STORAGE_KEYS.accessToken, urlToken);
          localStorage.setItem(STORAGE_KEYS.refreshToken, urlRefreshToken);
          document.cookie = `${STORAGE_KEYS.accessToken}=${urlToken}; path=/; max-age=604800; SameSite=Lax`;
          document.cookie = `${STORAGE_KEYS.refreshToken}=${urlRefreshToken}; path=/; max-age=604800; SameSite=Lax`;
          
          // Clear query params from URL
          const cleanUrl = window.location.protocol + "//" + window.location.host + window.location.pathname;
          window.history.replaceState({ path: cleanUrl }, "", cleanUrl);
        }
      }

      if (accessToken) {
        try {
          const userInfo = await fetchUserInfo(accessToken);
          if (userInfo) {
            setUser(userInfo);
            localStorage.setItem(STORAGE_KEYS.user, JSON.stringify(userInfo));
          } else {
            if (refreshTokenVal) {
              const success = await tryRefreshToken(refreshTokenVal);
              if (!success) {
                clearAuthStorage();
                redirectToLogin();
              }
            } else {
              clearAuthStorage();
              redirectToLogin();
            }
          }
        } catch (error) {
          console.error("Auth init error:", error);
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
        window.location.href = `http://localhost:3000/login?redirect=${encodeURIComponent(currentUrl)}`;
      }
    };

    initAuth();
  }, []);

  const clearAuthStorage = () => {
    localStorage.removeItem(STORAGE_KEYS.accessToken);
    localStorage.removeItem(STORAGE_KEYS.refreshToken);
    localStorage.removeItem(STORAGE_KEYS.user);
  };

  const fetchUserInfo = async (token: string): Promise<User | null> => {
    try {
      const response = await fetch(`${API_BASE_URLS.auth}${AUTH_ENDPOINTS.me}`, {
        headers: {
          Authorization: `Bearer ${token}`,
          "Content-Type": "application/json",
        },
      });

      if (!response.ok) {
        return null;
      }

      const result = await response.json();
      const userInfo: UserInfo = result.data;
      
      return {
        id: userInfo.user_id,
        personId: userInfo.person_id,
        email: userInfo.email,
        name: userInfo.name,
        role: userInfo.active_role,
        permissions: userInfo.permissions,
        scope: userInfo.scope,
      };
    } catch (error) {
      console.error("Error fetching user info:", error);
      return null;
    }
  };

  const tryRefreshToken = async (refreshTokenVal: string): Promise<boolean> => {
    try {
      const response = await fetch(`${API_BASE_URLS.auth}${AUTH_ENDPOINTS.refresh}`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ refresh_token: refreshTokenVal }),
      });

      if (!response.ok) {
        return false;
      }

      const result = await response.json();
      const tokenData: TokenResponse = result.data;

      localStorage.setItem(STORAGE_KEYS.accessToken, tokenData.access_token);
      localStorage.setItem(STORAGE_KEYS.refreshToken, tokenData.refresh_token);

      const userInfo = await fetchUserInfo(tokenData.access_token);
      if (userInfo) {
        setUser(userInfo);
        localStorage.setItem(STORAGE_KEYS.user, JSON.stringify(userInfo));
        return true;
      }

      return false;
    } catch (error) {
      console.error("Error refreshing token:", error);
      return false;
    }
  };

  const login = async (email: string, password: string) => {
    setIsLoading(true);
    setError("");

    try {
      const response = await fetch(`${API_BASE_URLS.auth}${AUTH_ENDPOINTS.login}`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ username: email, password }),
      });

      const result = await response.json();

      if (!response.ok) {
        throw new Error(result.message || "Login failed");
      }

      const tokenData: TokenResponse = result.data;
      const userInfo = await fetchUserInfo(tokenData.access_token);
      
      if (!userInfo) {
        throw new Error("Failed to get user info");
      }

      localStorage.setItem(STORAGE_KEYS.accessToken, tokenData.access_token);
      localStorage.setItem(STORAGE_KEYS.refreshToken, tokenData.refresh_token);
      localStorage.setItem(STORAGE_KEYS.user, JSON.stringify(userInfo));
      
      setUser(userInfo);
      router.push("/");
    } catch (error) {
      const message = error instanceof Error ? error.message : "Login failed";
      setError(message);
      throw error;
    } finally {
      setIsLoading(false);
    }
  };

  const logout = async () => {
    const accessToken = localStorage.getItem(STORAGE_KEYS.accessToken);
    
    if (accessToken) {
      try {
        await fetch(`${API_BASE_URLS.auth}${AUTH_ENDPOINTS.login}`, {
          method: "POST",
          headers: {
            Authorization: `Bearer ${accessToken}`,
            "Content-Type": "application/json",
          },
        });
      } catch (error) {
        // Ignore API errors on logout
      }
    }

    clearAuthStorage();
    setUser(null);
    router.push("/login");
  };

  const refreshToken = async () => {
    const refreshTokenValue = localStorage.getItem(STORAGE_KEYS.refreshToken);
    if (refreshTokenValue) {
      await tryRefreshToken(refreshTokenValue);
    }
  };

  return (
    <AuthContext.Provider
      value={{
        user,
        isLoading,
        login,
        logout,
        refreshToken,
        isAuthenticated: !!user,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
}
