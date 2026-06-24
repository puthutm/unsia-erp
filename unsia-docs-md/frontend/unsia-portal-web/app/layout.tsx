import './globals.css';
import React from 'react';

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
        <div className="bg-gradient-mesh"></div>
        {children}
      </body>
    </html>
  );
}
