"use client";

import { useState } from "react";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { useAuth } from "@/contexts/auth-context";
import { FRONTEND_URLS } from "@/lib/constants";

// Navigation items for academic module
const navigation = [
  { name: "Beranda", href: "/", icon: "home" },
  { name: "Mahasiswa", href: "/student", icon: "users" },
  { name: "KRS", href: "/krs", icon: "book-open" },
  { name: "Jadwal", href: "/schedule", icon: "calendar" },
  { name: "Nilai", href: "/grade", icon: "award" },
  { name: "Transkrip", href: "/transcript", icon: "file-text" },
];

const modules = [
  { id: "pmb", name: "PMB", href: FRONTEND_URLS.pmb, icon: "graduation-cap" },
  { id: "finance", name: "Keuangan", href: FRONTEND_URLS.finance, icon: "wallet" },
  { id: "academic", name: "Akademik", href: "/", icon: "book" },
  { id: "lms", name: "LMS", href: FRONTEND_URLS.lms, icon: "monitor" },
  { id: "crm", name: "CRM", href: FRONTEND_URLS.crm, icon: "users" },
  { id: "hris", name: "HRIS", href: FRONTEND_URLS.hris, icon: "briefcase" },
  { id: "assessment", name: "CBT", href: FRONTEND_URLS.assessment, icon: "file-text" },
  { id: "reference", name: "Referensi", href: FRONTEND_URLS.reference, icon: "database" },
];

// Icon components
const icons: Record<string, React.FC<{ className?: string }>> = {
  home: ({ className }) => (
    <svg className={className} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
      <path d="M3 9l9-7 9 7v11a2 2 0 01-2 2H5a2 2 0 01-2-2z" />
      <polyline points="9 22 9 12 15 12 15 22" />
    </svg>
  ),
  users: ({ className }) => (
    <svg className={className} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
      <path d="M17 21v-2a4 4 0 00-4-4H5a4 4 0 00-4 4v2" />
      <circle cx="9" cy="7" r="4" />
      <path d="M23 21v-2a4 4 0 00-3-3.87M16 3.13a4 4 0 010 7.75" />
    </svg>
  ),
  calendar: ({ className }) => (
    <svg className={className} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
      <rect x="3" y="4" width="18" height="18" rx="2" ry="2" />
      <line x1="16" y1="2" x2="16" y2="6" />
      <line x1="8" y1="2" x2="8" y2="6" />
      <line x1="3" y1="10" x2="21" y2="10" />
    </svg>
  ),
  "file-text": ({ className }) => (
    <svg className={className} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
      <path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8z" />
      <polyline points="14 2 14 8 20 8" />
      <line x1="16" y1="13" x2="8" y2="13" />
      <line x1="16" y1="17" x2="8" y2="17" />
    </svg>
  ),
  "book-open": ({ className }) => (
    <svg className={className} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
      <path d="M2 3h6a4 4 0 014 4v14a3 3 0 00-3-3H2z" />
      <path d="M22 3h-6a4 4 0 00-4 4v14a3 3 0 013-3h7z" />
    </svg>
  ),
  award: ({ className }) => (
    <svg className={className} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
      <circle cx="12" cy="8" r="7" />
      <polyline points="8.21 13.89 7 23 12 20 17 23 15.79 13.88" />
    </svg>
  ),
  logout: ({ className }) => (
    <svg className={className} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
      <path d="M9 21H5a2 2 0 01-2-2V5a2 2 0 012-2h4" />
      <polyline points="16 17 21 12 16 7" />
      <line x1="21" y1="12" x2="9" y2="12" />
    </svg>
  ),
  menu: ({ className }) => (
    <svg className={className} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
      <line x1="3" y1="12" x2="21" y2="12" />
      <line x1="3" y1="6" x2="21" y2="6" />
      <line x1="3" y1="18" x2="21" y2="18" />
    </svg>
  ),
};

export default function PortalLayout({ children }: { children: React.ReactNode }) {
  const pathname = usePathname();
  const { user, logout } = useAuth();
  const [sidebarOpen, setSidebarOpen] = useState(true);

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Top Header */}
      <header className="fixed top-0 left-0 right-0 z-50 bg-white border-b border-gray-200">
        <div className="flex items-center justify-between h-16 px-4">
          <div className="flex items-center gap-4">
            <button
              onClick={() => setSidebarOpen(!sidebarOpen)}
              className="p-2 rounded-lg hover:bg-gray-100"
            >
              <icons.menu className="w-5 h-5" />
            </button>
            <a href={FRONTEND_URLS.portal} className="flex items-center gap-2">
              <div className="w-8 h-8 bg-blue-600 rounded-lg flex items-center justify-center">
                <span className="text-white font-bold text-sm">U</span>
              </div>
              <span className="font-semibold text-gray-900">UNSIA</span>
            </a>
          </div>
          
          <div className="flex items-center gap-4">
            {/* Module Selector */}
            <div className="flex items-center gap-1 bg-gray-100 rounded-lg p-1 overflow-x-auto max-w-[200px] sm:max-w-[400px] md:max-w-none whitespace-nowrap scrollbar-none">
              {modules.map((mod) => (
                <a
                  key={mod.id}
                  href={mod.href}
                  className={`px-3 py-1.5 text-sm font-medium rounded-md transition-colors ${
                    mod.id === "academic"
                      ? "bg-white text-gray-900 shadow-sm"
                      : "text-gray-600 hover:text-gray-900"
                  }`}
                >
                  {mod.name}
                </a>
              ))}
            </div>
            
            {/* User Menu */}
            <div className="flex items-center gap-3">
              <div className="text-right">
                <p className="text-sm font-medium text-gray-900">{user?.name || "Admin"}</p>
                <p className="text-xs text-gray-500">{user?.role || "Administrator"}</p>
              </div>
              <button
                onClick={logout}
                className="p-2 rounded-lg hover:bg-gray-100 text-gray-500"
                title="Logout"
              >
                <icons.logout className="w-5 h-5" />
              </button>
            </div>
          </div>
        </div>
      </header>

      {/* Sidebar */}
      <aside
        className={`fixed left-0 top-16 bottom-0 bg-white border-r border-gray-200 transition-all duration-300 z-40 ${
          sidebarOpen ? "w-64" : "w-0 overflow-hidden"
        }`}
      >
        <nav className="p-4 space-y-1">
          {navigation.map((item) => {
            const Icon = icons[item.icon as keyof typeof icons];
            const isCurrent = pathname === item.href;
            return (
              <Link
                key={item.name}
                href={item.href}
                className={`flex items-center gap-3 px-3 py-2 rounded-lg transition-colors ${
                  isCurrent
                    ? "bg-blue-50 text-blue-600"
                    : "text-gray-600 hover:bg-gray-100"
                }`}
              >
                {Icon && <Icon className="w-5 h-5" />}
                <span className="font-medium">{item.name}</span>
              </Link>
            );
          })}
        </nav>
      </aside>

      {/* Main Content */}
      <main className={`pt-16 min-h-screen transition-all duration-300 ${sidebarOpen ? "ml-64" : "ml-0"}`}>
        <div className="p-6">
          {children}
        </div>
      </main>
    </div>
  );
}
