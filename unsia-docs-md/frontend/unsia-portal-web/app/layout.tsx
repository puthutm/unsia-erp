import './globals.css';
import React from 'react';
import { AuthProvider } from '@/contexts/auth-context';
import { ReferenceProvider } from '@/contexts/reference-context';

export const metadata = {
  title: 'UNSIA ERP Portal',
  description: 'Integrated Identity, Academic, PMB, and Finance Portal for Universitas Siber Asia',
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="id">
      <body>
        <AuthProvider>
          <ReferenceProvider>
            <div className="bg-gradient-mesh"></div>
            {children}
          </ReferenceProvider>
        </AuthProvider>
      </body>
    </html>
  );
}
