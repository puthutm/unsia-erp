import type { Metadata } from "next";
import "./globals.css";
import { AuthProvider } from "../contexts/auth-context";
import { ReferenceProvider } from "../contexts/reference-context";

export const metadata: Metadata = {
  title: "PMB UNSIA - Penerimaan Mahasiswa Baru",
  description: "Sistem Pengelolaan Penerimaan Mahasiswa Baru Universitas Siber Asia",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="id">
      <body className="antialiased">
        <AuthProvider>
          <ReferenceProvider>
            {children}
          </ReferenceProvider>
        </AuthProvider>
      </body>
    </html>
  );
}
