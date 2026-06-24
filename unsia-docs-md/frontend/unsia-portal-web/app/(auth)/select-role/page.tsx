'use client';

import React, { useEffect, useState } from 'react';

export default function SelectRolePage() {
  const [roles, setRoles] = useState<string[]>([]);
  const [selectedRole, setSelectedRole] = useState('');

  useEffect(() => {
    // Get mock user role or fetch available roles
    const mockRole = localStorage.getItem('mock_user_role') || 'admin';
    if (mockRole === 'admin') {
      setRoles(['super_admin', 'admin_pmb', 'admin_keuangan']);
    } else if (mockRole === 'mahasiswa') {
      setRoles(['mahasiswa']);
    } else if (mockRole === 'pendaftar') {
      setRoles(['pendaftar']);
    } else {
      setRoles(['public']);
    }
  }, []);

  const handleSelect = (role: string) => {
    setSelectedRole(role);
    localStorage.setItem('active_role', role);
    window.location.href = '/dashboard';
  };

  return (
    <div className="auth-container">
      <div className="card-glass auth-card">
        <h2 style={{ fontSize: '24px', fontWeight: 600, marginBottom: '8px', textAlign: 'center' }}>
          Pilih Peran Aktif
        </h2>
        <p style={{ color: 'var(--text-secondary)', marginBottom: '24px', fontSize: '14px', textAlign: 'center' }}>
          Akun Anda memiliki beberapa peran. Pilih peran untuk melanjutkan ke dashboard.
        </p>

        <div style={{ display: 'flex', flexDirection: 'column', gap: '12px' }}>
          {roles.map((role) => (
            <button
              key={role}
              onClick={() => handleSelect(role)}
              className="btn-secondary"
              style={{
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'space-between',
                padding: '16px 20px',
                textAlign: 'left',
                border: selectedRole === role ? '1px solid var(--accent-primary)' : '1px solid var(--border-glass)',
                backgroundColor: selectedRole === role ? 'rgba(99, 102, 241, 0.1)' : 'var(--bg-tertiary)',
              }}
            >
              <span style={{ fontWeight: 600, textTransform: 'capitalize' }}>
                {role.replace('_', ' ')}
              </span>
              <span style={{ fontSize: '12px', color: 'var(--text-secondary)' }}>
                Klik untuk Memilih &rarr;
              </span>
            </button>
          ))}
        </div>
      </div>
    </div>
  );
}
