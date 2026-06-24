import React from 'react';

export default function Home() {
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
