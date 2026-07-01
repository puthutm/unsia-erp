import './globals.css';
import React from 'react';
import { AuthProvider } from '@/contexts/auth-context';
import { ReferenceProvider } from '@/contexts/reference-context';
import QueryProvider from '@/components/providers/query-provider';

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
          <QueryProvider>
            <ReferenceProvider>
              <div className="bg-gradient-mesh"></div>
              {children}
            </ReferenceProvider>
          </QueryProvider>
        </AuthProvider>
      </body>
    </html>
  );
}
