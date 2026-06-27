'use client';

import { useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useAuth } from '@/contexts/auth-context';
import { STORAGE_KEYS } from '@/lib/constants';

export default function Home() {
  const router = useRouter();
  const { isAuthenticated, isLoading } = useAuth();

  useEffect(() => {
    if (!isLoading) {
      if (isAuthenticated) {
        router.replace('/dashboard');
      } else {
        // Check if there's a stored session
        const accessToken = localStorage.getItem(STORAGE_KEYS.accessToken);
        if (accessToken) {
          router.replace('/dashboard');
        }
        // Otherwise stay on home page (show login prompt)
      }
    }
  }, [isAuthenticated, isLoading, router]);

  // Show loading while checking auth
  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-slate-900">
        <div className="text-center">
          <div className="w-16 h-16 border-4 border-blue-500 border-t-transparent rounded-full animate-spin mx-auto mb-4"></div>
          <p className="text-white/60">Memuat...</p>
        </div>
      </div>
    );
  }

  // If not authenticated, show landing page
  return (
    <div className="auth-container">
      <div className="card-glass auth-card" style={{ textAlign: 'center' }}>
        <h1 style={{ fontSize: '32px', marginBottom: '10px', fontWeight: 700 }}>
          UNSIA ERP
        </h1>
        <p style={{ color: 'var(--text-secondary)', marginBottom: '30px' }}>
          Portal ERP Terintegrasi Universitas Siber Asia. Silakan masuk untuk mengakses SIAKAD, tagihan keuangan, dan kelas online Anda.
        </p>
        <a href="/login" style={{ textDecoration: 'none' }}>
          <button className="btn-primary" style={{ width: '100%' }}>
            Masuk ke Portal
          </button>
        </a>
      </div>
    </div>
  );
}
