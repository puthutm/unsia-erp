import "./globals.css";
import { AuthProvider } from "@/contexts/auth-context";
import { ReferenceProvider } from "@/contexts/reference-context";
import PortalLayout from "@/components/layout/PortalLayout";

export const metadata = {
  title: "UNSIA Finance Portal",
  description: "Finance Management System",
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
            <PortalLayout>
              {children}
            </PortalLayout>
          </ReferenceProvider>
        </AuthProvider>
      </body>
    </html>
  );
}
