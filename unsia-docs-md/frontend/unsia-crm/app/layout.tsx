import "@/app/globals.css";
import { AuthProvider } from "@/contexts/auth-context";
import { ReferenceProvider } from "@/contexts/reference-context";

export const metadata = {
  title: "CRM - UNSIA ERP",
  description: "Standalone Module for unsia-crm",
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
            <div className="min-h-screen bg-gray-50 p-6">
              {children}
            </div>
          </ReferenceProvider>
        </AuthProvider>
      </body>
    </html>
  );
}
