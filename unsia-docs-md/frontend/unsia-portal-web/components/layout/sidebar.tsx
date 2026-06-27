"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { cn } from "@/lib/utils";

// Navigation configuration for each module
const moduleNavigation = {
  pmb: [
    { name: "Beranda", href: "/dashboard", icon: "home" },
    { name: "Pendaftar", href: "/dashboard/applicant", icon: "users" },
    { name: "Gelombang", href: "/dashboard/wave", icon: "calendar" },
    { name: "Dokumen", href: "/dashboard/document", icon: "file" },
    { name: "Pembayaran", href: "/dashboard/payment", icon: "credit-card" },
    { name: "Seleksi", href: "/dashboard/selection", icon: "check-circle" },
  ],
  finance: [
    { name: "Beranda", href: "/finance", icon: "home" },
    { name: "Invoice", href: "/finance/invoice", icon: "file-text" },
    { name: "Pembayaran", href: "/finance/payment", icon: "credit-card" },
    { name: "Komponen", href: "/finance/component", icon: "package" },
    { name: "Laporan", href: "/finance/report", icon: "bar-chart" },
  ],
  academic: [
    { name: "Beranda", href: "/academic", icon: "home" },
    { name: "Mahasiswa", href: "/academic/student", icon: "users" },
    { name: "KRS", href: "/academic/krs", icon: "book-open" },
    { name: "Jadwal", href: "/academic/schedule", icon: "calendar" },
    { name: "Nilai", href: "/academic/grade", icon: "award" },
    { name: "Transkrip", href: "/academic/transcript", icon: "file-text" },
  ],
  lms: [
    { name: "Beranda", href: "/lms", icon: "home" },
    { name: "Kursus", href: "/lms/course", icon: "book" },
    { name: "Kelas", href: "/lms/class", icon: "users" },
    { name: "Sesi", href: "/lms/session", icon: "video" },
    { name: "Tugas", href: "/lms/assignment", icon: "file-text" },
    { name: "Kehadiran", href: "/lms/attendance", icon: "check-circle" },
  ],
};

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
  file: ({ className }) => (
    <svg className={className} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
      <path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8z" />
      <polyline points="14 2 14 8 20 8" />
    </svg>
  ),
  "credit-card": ({ className }) => (
    <svg className={className} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
      <rect x="1" y="4" width="22" height="16" rx="2" ry="2" />
      <line x1="1" y1="10" x2="23" y2="10" />
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
  "bar-chart": ({ className }) => (
    <svg className={className} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
      <line x1="12" y1="20" x2="12" y2="10" />
      <line x1="18" y1="20" x2="18" y2="4" />
      <line x1="6" y1="20" x2="6" y2="16" />
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
  book: ({ className }) => (
    <svg className={className} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
      <path d="M4 19.5A2.5 2.5 0 016.5 17H20" />
      <path d="M6.5 2H20v20H6.5A2.5 2.5 0 014 19.5v-15A2.5 2.5 0 016.5 2z" />
    </svg>
  ),
  video: ({ className }) => (
    <svg className={className} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
      <polygon points="23 7 16 12 23 17 23 7" />
      <rect x="1" y="5" width="15" height="14" rx="2" ry="2" />
    </svg>
  ),
  "check-circle": ({ className }) => (
    <svg className={className} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
      <path d="M22 11.08V12a10 10 0 11-5.93-9.14" />
      <polyline points="22 4 12 14.01 9 11.01" />
    </svg>
  ),
  package: ({ className }) => (
    <svg className={className} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
      <path d="M16.5 9.4l-9-5.19M21 16V8a2 2 0 00-1-1.73l-7-4a2 2 0 00-2 0l-7 4A2 2 0 003 8v8a2 2 0 001 1.73l7 4a2 2 0 002 0l7-4A2 2 0 0021 16z" />
      <polyline points="3.27 6.96 12 12.01 20.73 6.96" />
      <line x1="12" y1="22.08" x2="12" y2="12" />
    </svg>
  ),
};

interface SidebarProps {
  module: "pmb" | "finance" | "academic" | "lms";
  isOpen?: boolean;
}

export function Sidebar({ module, isOpen = true }: SidebarProps) {
  const pathname = usePathname();
  const navItems = moduleNavigation[module] || moduleNavigation.pmb;

  return (
    <aside
      className={cn(
        "fixed left-0 top-16 bottom-0 bg-white border-r border-gray-200 transition-all duration-300 z-40",
        isOpen ? "w-64" : "w-0 overflow-hidden"
      )}
    >
      <nav className="p-4 space-y-1">
        {navItems.map((item) => {
          const Icon = icons[item.icon];
          const isActive = pathname === item.href || pathname.startsWith(item.href + "/");
          
          return (
            <Link
              key={item.href}
              href={item.href}
              className={cn(
                "flex items-center gap-3 px-3 py-2 rounded-lg transition-colors",
                isActive
                  ? "bg-blue-50 text-blue-600 font-medium"
                  : "text-gray-600 hover:bg-gray-100"
              )}
            >
              {Icon && <Icon className="w-5 h-5" />}
              <span>{item.name}</span>
            </Link>
          );
        })}
      </nav>
    </aside>
  );
}
