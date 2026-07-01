"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { Button } from "@/components/ui/button";
import { API_BASE_URLS, AUTH_ENDPOINTS, STORAGE_KEYS } from "@/lib/constants";

// Simple SVG Icons
function EyeIcon({ className }: { className?: string }) {
  return (
    <svg className={className} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
      <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z" />
      <circle cx="12" cy="12" r="3" />
    </svg>
  );
}

function EyeSlashIcon({ className }: { className?: string }) {
  return (
    <svg className={className} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
      <path d="M17.94 17.94A10.07 10.07 0 0112 20c-7 0-11-8-11-8a18.45 18.45 0 015.06-5.94M9.9 4.24A9.12 9.12 0 0112 4c7 0 11 8 11 8a18.5 18.5 0 01-2.16 3.19m-6.72-1.07a3 3 0 11-4.24-4.24" />
      <line x1="1" y1="1" x2="23" y2="23" />
    </svg>
  );
}

function EnvelopeIcon({ className }: { className?: string }) {
  return (
    <svg className={className} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
      <path d="M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z" />
      <polyline points="22,6 12,13 2,6" />
    </svg>
  );
}

function LockIcon({ className }: { className?: string }) {
  return (
    <svg className={className} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
      <rect x="3" y="11" width="18" height="11" rx="2" ry="2" />
      <path d="M7 11V7a5 5 0 0110 0v4" />
    </svg>
  );
}

// Login function that doesn't require AuthProvider
async function loginUser(email: string, password: string) {
  const response = await fetch(`${API_BASE_URLS.auth}${AUTH_ENDPOINTS.login}`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ username: email, password }),
  });

  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.message || "Login failed");
  }

  const result = await response.json();
  const { access_token, refresh_token } = result.data;

  // Store tokens in localStorage
  localStorage.setItem(STORAGE_KEYS.accessToken, access_token);
  localStorage.setItem(STORAGE_KEYS.refreshToken, refresh_token);

  // Store tokens in cookies for middleware
  document.cookie = `${STORAGE_KEYS.accessToken}=${access_token}; path=/; max-age=604800; SameSite=Lax`;
  document.cookie = `${STORAGE_KEYS.refreshToken}=${refresh_token}; path=/; max-age=604800; SameSite=Lax`;

  // Fetch user info
  const userResponse = await fetch(`${API_BASE_URLS.auth}${AUTH_ENDPOINTS.me}`, {
    headers: {
      Authorization: `Bearer ${access_token}`,
      "Content-Type": "application/json",
    },
  });

  if (userResponse.ok) {
    const userResult = await userResponse.json();
    localStorage.setItem(STORAGE_KEYS.user, JSON.stringify(userResult.data));
  }

  return result.data;
}

export default function LoginPage() {
  const router = useRouter();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [showPassword, setShowPassword] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState("");

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);
    setError("");

    try {
      const data = await loginUser(email, password);
      
      // Handle SSO redirect if present in query params
      const searchParams = new URLSearchParams(window.location.search);
      const redirectUrl = searchParams.get("redirect");
      if (redirectUrl) {
        window.location.href = `${redirectUrl}?token=${data.access_token}&refresh_token=${data.refresh_token}`;
      } else {
        router.push("/dashboard");
      }
    } catch (error) {
      const message = error instanceof Error ? error.message : "Login failed";
      setError(message);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-slate-900 p-4">
      {/* Background Pattern */}
      <div className="absolute inset-0 overflow-hidden">
        <div className="absolute -top-40 -right-40 w-96 h-96 bg-blue-500/20 rounded-full blur-3xl" />
        <div className="absolute -bottom-40 -left-40 w-96 h-96 bg-cyan-500/20 rounded-full blur-3xl" />
        <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[800px] h-[800px] bg-blue-500/10 rounded-full blur-3xl" />
      </div>

      {/* Login Card */}
      <div className="relative w-full max-w-md animate-in">
        <div className="bg-white/95 backdrop-blur-xl rounded-xl shadow-lg p-8 border border-white/20">
          {/* Logo */}
          <div className="text-center mb-8">
            <div className="inline-flex items-center justify-center w-20 h-20 bg-blue-600 rounded-xl mb-4">
              <svg
                className="w-12 h-12 text-white"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
                strokeWidth={2}
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  d="M12 14l9-5-9-5-9 5 9 5z"
                />
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  d="M12 14l6.16-3.422a12.083 12.083 0 01.665 6.479A11.952 11.952 0 0012 20.055a11.952 11.952 0 00-6.824-2.998 12.083 12.083 0 01.665-6.479L12 14z"
                />
              </svg>
            </div>
            <h1 className="text-2xl font-bold text-slate-800">
              UNSIA
            </h1>
            <p className="text-sm text-slate-500 mt-1">
              Universitas Siber Asia
            </p>
          </div>

          {/* Welcome Text */}
          <div className="text-center mb-6">
            <h2 className="text-xl font-semibold text-slate-800">
              Selamat Datang
            </h2>
            <p className="text-sm text-slate-500 mt-1">
              Masuk untuk mengakses sistem akademik
            </p>
          </div>

          {/* Error Message */}
          {error && (
            <div className="mb-4 p-3 bg-red-50 border border-red-200 rounded-lg text-red-600 text-sm">
              {error}
            </div>
          )}

          {/* Login Form */}
          <form onSubmit={handleSubmit} className="space-y-4">
            {/* Username Field */}
            <div className="space-y-2">
              <label
                htmlFor="username"
                className="text-sm font-medium text-slate-700"
              >
                Username / NIM / NIP
              </label>
              <div className="relative">
                <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                  <EnvelopeIcon className="h-5 w-5 text-slate-400" />
                </div>
                <input
                  type="text"
                  id="username"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  placeholder="Masukkan NIM, NIP, atau Email"
                  className="w-full pl-10 pr-4 py-3 bg-slate-50 border border-slate-300 rounded-lg text-slate-800 placeholder:text-slate-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all"
                  required
                />
              </div>
            </div>

            {/* Password Field */}
            <div className="space-y-2">
              <label
                htmlFor="password"
                className="text-sm font-medium text-slate-700"
              >
                Password
              </label>
              <div className="relative">
                <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                  <LockIcon className="h-5 w-5 text-slate-400" />
                </div>
                <input
                  type={showPassword ? "text" : "password"}
                  id="password"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  placeholder="••••••••"
                  className="w-full pl-10 pr-12 py-3 bg-slate-50 border border-slate-300 rounded-lg text-slate-800 placeholder:text-slate-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all"
                  required
                />
                <button
                  type="button"
                  onClick={() => setShowPassword(!showPassword)}
                  className="absolute inset-y-0 right-0 pr-3 flex items-center"
                >
                  {showPassword ? (
                    <EyeSlashIcon className="h-5 w-5 text-slate-400 hover:text-slate-600" />
                  ) : (
                    <EyeIcon className="h-5 w-5 text-slate-400 hover:text-slate-600" />
                  )}
                </button>
              </div>
            </div>

            {/* Forgot Password */}
            <div className="text-right">
              <a
                href="/forgot-password"
                className="text-sm text-blue-600 hover:text-blue-700"
              >
                Lupa password?
              </a>
            </div>

            {/* Submit Button */}
            <Button
              type="submit"
              className="w-full"
              size="lg"
              isLoading={isLoading}
            >
              Masuk
            </Button>
          </form>

          {/* Register Link */}
          <p className="text-center text-sm text-slate-500 mt-6">
            Belum punya akun?{" "}
            <a
              href="/register"
              className="text-blue-600 hover:text-blue-700 font-medium"
            >
              Daftar Sekarang
            </a>
          </p>
        </div>

        {/* Footer */}
        <p className="text-center text-xs text-white/60 mt-6">
          © 2024 Universitas Siber Asia. All rights reserved.
        </p>
      </div>
    </div>
  );
}
