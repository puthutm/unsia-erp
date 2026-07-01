"use client";

import { useState, useEffect } from "react";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { useAuth } from "@/contexts/auth-context";
import { useReference } from "@/contexts/reference-context";
import { FRONTEND_URLS } from "@/lib/constants";

// Navigation items for each module
const navigation = {
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
  crm: [
    { name: "Beranda", href: "/crm", icon: "home" },
    { name: "Leads", href: "/crm/leads", icon: "users" },
  ],
  hris: [
    { name: "Beranda", href: "/hris", icon: "home" },
    { name: "Absensi", href: "/hris/attendance", icon: "check-circle" },
    { name: "Cuti", href: "/hris/leave", icon: "calendar" },
  ],
  assessment: [
    { name: "Ujian CBT", href: "/assessment", icon: "file-text" },
  ],
  reference: [
    { name: "Data Master", href: "/reference", icon: "database" },
  ],
};

const moduleColors = {
  pmb: { bg: "bg-purple-600", light: "bg-purple-50", text: "text-purple-600", border: "border-purple-200" },
  finance: { bg: "bg-emerald-600", light: "bg-emerald-50", text: "text-emerald-600", border: "border-emerald-200" },
  academic: { bg: "bg-blue-600", light: "bg-blue-50", text: "text-blue-600", border: "border-blue-200" },
  lms: { bg: "bg-orange-600", light: "bg-orange-50", text: "text-orange-600", border: "border-orange-200" },
  crm: { bg: "bg-violet-600", light: "bg-violet-50", text: "text-violet-600", border: "border-violet-200" },
  hris: { bg: "bg-rose-600", light: "bg-rose-50", text: "text-rose-600", border: "border-rose-200" },
  assessment: { bg: "bg-amber-600", light: "bg-amber-50", text: "text-amber-600", border: "border-amber-200" },
  reference: { bg: "bg-slate-600", light: "bg-slate-50", text: "text-slate-600", border: "border-slate-200" },
};

const modules = [
  { id: "pmb", name: "PMB", href: "/dashboard", icon: "graduation-cap" },
  { id: "finance", name: "Keuangan", href: "/finance", icon: "wallet" },
  { id: "academic", name: "Akademik", href: "/academic", icon: "book" },
  { id: "lms", name: "LMS", href: "/lms", icon: "monitor" },
  { id: "crm", name: "CRM", href: "/crm", icon: "users" },
  { id: "hris", name: "HRIS", href: "/hris", icon: "briefcase" },
  { id: "assessment", name: "CBT", href: "/assessment", icon: "file-text" },
  { id: "reference", name: "Referensi", href: "/reference", icon: "database" },
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
  "graduation-cap": ({ className }) => (
    <svg className={className} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
      <path d="M22 10v6M2 10l10-5 10 5-10 5z" />
      <path d="M6 12v5c3 3 9 3 12 0v-5" />
    </svg>
  ),
  wallet: ({ className }) => (
    <svg className={className} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
      <rect x="1" y="4" width="22" height="16" rx="2" ry="2" />
      <line x1="1" y1="10" x2="23" y2="10" />
      <path d="M16 14a2 2 0 100 4 2 2 0 000-4z" />
    </svg>
  ),
  monitor: ({ className }) => (
    <svg className={className} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
      <rect x="2" y="3" width="20" height="14" rx="2" ry="2" />
      <line x1="8" y1="21" x2="16" y2="21" />
      <line x1="12" y1="17" x2="12" y2="21" />
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
  bell: ({ className }) => (
    <svg className={className} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
      <path d="M18 8A6 6 0 006 8c0 7-3 9-3 9h18s-3-2-3-9M13.73 21a2 2 0 01-3.46 0" />
    </svg>
  ),
  settings: ({ className }) => (
    <svg className={className} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
      <circle cx="12" cy="12" r="3" />
      <path d="M19.4 15a1.65 1.65 0 00.33 1.82l.06.06a2 2 0 010 2.83 2 2 0 01-2.83 0l-.06-.06a1.65 1.65 0 00-1.82-.33 1.65 1.65 0 00-1 1.51V21a2 2 0 01-2 2 2 2 0 01-2-2v-.09A1.65 1.65 0 009 19.4a1.65 1.65 0 00-1.82.33l-.06.06a2 2 0 01-2.83 0 2 2 0 010-2.83l.06-.06a1.65 1.65 0 00.33-1.82 1.65 1.65 0 00-1.51-1H3a2 2 0 01-2-2 2 2 0 012-2h.09A1.65 1.65 0 004.6 9a1.65 1.65 0 00-.33-1.82l-.06-.06a2 2 0 010-2.83 2 2 0 012.83 0l.06.06a1.65 1.65 0 001.82.33H9a1.65 1.65 0 001-1.51V3a2 2 0 012-2 2 2 0 012 2v.09a1.65 1.65 0 001 1.51 1.65 1.65 0 001.82-.33l.06-.06a2 2 0 012.83 0 2 2 0 010 2.83l-.06.06a1.65 1.65 0 00-.33 1.82V9a1.65 1.65 0 001.51 1H21a2 2 0 012 2 2 2 0 01-2 2h-.09a1.65 1.65 0 00-1.51 1z" />
    </svg>
  ),
  search: ({ className }) => (
    <svg className={className} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
      <circle cx="11" cy="11" r="8" />
      <line x1="21" y1="21" x2="16.65" y2="16.65" />
    </svg>
  ),
  plus: ({ className }) => (
    <svg className={className} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
      <line x1="12" y1="5" x2="12" y2="19" />
      <line x1="5" y1="12" x2="19" y2="12" />
    </svg>
  ),
  briefcase: ({ className }) => (
    <svg className={className} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
      <rect x="2" y="7" width="20" height="14" rx="2" ry="2" />
      <path d="M16 21V5a2 2 0 00-2-2h-4a2 2 0 00-2 2v16" />
    </svg>
  ),
  database: ({ className }) => (
    <svg className={className} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
      <ellipse cx="12" cy="5" rx="9" ry="3" />
      <path d="M3 5v14c0 1.66 4 3 9 3s9-1.34 9-3V5" />
      <path d="M3 12c0 1.66 4 3 9 3s9-1.34 9-3" />
    </svg>
  ),
};

// Get current module based on pathname
const getCurrentModule = (pathname: string) => {
  if (pathname.startsWith("/dashboard") || pathname.startsWith("/pmb")) return "pmb";
  if (pathname.startsWith("/finance")) return "finance";
  if (pathname.startsWith("/academic")) return "academic";
  if (pathname.startsWith("/lms")) return "lms";
  if (pathname.startsWith("/crm")) return "crm";
  if (pathname.startsWith("/hris")) return "hris";
  if (pathname.startsWith("/assessment")) return "assessment";
  if (pathname.startsWith("/reference")) return "reference";
  return "pmb";
};

export default function PortalLayout({ children }: { children: React.ReactNode }) {
  const pathname = usePathname();
  const { user, logout } = useAuth();
  const [sidebarOpen, setSidebarOpen] = useState(true);
  const currentModule = getCurrentModule(pathname);
  const moduleNav = navigation[currentModule as keyof typeof navigation] || navigation.pmb;

  return (
    <div className="min-h-screen bg-slate-50 text-slate-800 font-sans">
      {/* Top Header - Glassmorphic */}
      <header className="fixed top-0 left-0 right-0 z-50 bg-white/95 backdrop-blur-md border-b border-slate-100 shadow-sm transition-all duration-300">
        <div className="flex items-center justify-between h-16 px-6">
          <div className="flex items-center gap-4">
            <button
              onClick={() => setSidebarOpen(!sidebarOpen)}
              className="p-2 rounded-lg hover:bg-slate-100 text-slate-600 transition-colors"
            >
              <icons.menu className="w-5 h-5" />
            </button>
            <Link href="/dashboard" className="flex items-center gap-3">
              <div className="w-8 h-8 bg-[#0f487b] rounded-lg flex items-center justify-center shadow relative overflow-hidden">
                <div className="absolute inset-0 bg-[#FED524]/20"></div>
                <span className="text-white font-bold text-sm relative z-10">U</span>
              </div>
              <div className="hidden sm:block">
                <span className="font-bold text-slate-800 tracking-wide text-sm">ERP UNSIA<span className="text-[#FED524]">.</span></span>
              </div>
            </Link>
          </div>
          
          <div className="flex items-center gap-6">
            {/* Module Selector - Pills */}
            <div className="flex items-center gap-1.5 bg-slate-100 rounded-xl p-1 overflow-x-auto max-w-[200px] sm:max-w-[400px] md:max-w-none whitespace-nowrap scrollbar-none shadow-inner">
              {modules.map((mod) => {
                const targetUrl = FRONTEND_URLS[mod.id as keyof typeof FRONTEND_URLS] || mod.href;
                const isActive = pathname.startsWith(mod.href);
                return (
                  <a
                    key={mod.id}
                    href={targetUrl}
                    className={`px-3 py-1.5 text-xs font-bold rounded-lg transition-all ${
                      isActive
                        ? "bg-white text-[#0f487b] shadow-sm scale-102"
                        : "text-slate-500 hover:text-slate-800 hover:bg-white/50"
                    }`}
                  >
                    {mod.name}
                  </a>
                );
              })}
            </div>
            
            {/* User Menu */}
            <div className="flex items-center gap-4 border-l border-slate-100 pl-6">
              <div className="text-right hidden md:block">
                <p className="text-sm font-semibold text-slate-800 leading-none">{user?.name || "Admin"}</p>
                <p className="text-[10px] font-bold text-slate-400 uppercase tracking-wider mt-1">{user?.role || "Administrator"}</p>
              </div>
              <button
                onClick={logout}
                className="p-2 rounded-lg hover:bg-rose-50 text-slate-400 hover:text-rose-600 transition-colors"
                title="Logout"
              >
                <icons.logout className="w-5 h-5" />
              </button>
            </div>
          </div>
        </div>
      </header>

      {/* Sidebar - Gradient Theme */}
      <aside
        className={`fixed left-0 top-16 bottom-0 bg-gradient-to-b from-[#0f487b] to-[#00719f] transition-all duration-300 z-40 border-r border-white/5 flex flex-col shadow-xl ${
          sidebarOpen ? "w-64" : "w-0 overflow-hidden"
        }`}
      >
        <div className="flex-1 overflow-y-auto no-scrollbar py-6 px-4 space-y-6">
          <div>
            <p className="px-3 text-[10px] font-bold uppercase tracking-widest text-white/50 mb-2">
              Menu Navigasi
            </p>
            <nav className="space-y-1.5">
              {moduleNav.map((item) => {
                const Icon = icons[item.icon as keyof typeof icons];
                const isActive = pathname === item.href;
                return (
                  <Link
                    key={item.name}
                    href={item.href}
                    className={`group flex items-center gap-3 px-3 py-2.5 rounded-lg transition-all text-sm border-l-4 ${
                      isActive
                        ? "bg-white/15 text-white font-bold border-[#FED524] shadow-md"
                        : "text-white/70 hover:bg-white/10 hover:text-white border-transparent"
                    }`}
                  >
                    {Icon && (
                      <Icon
                        className={`w-5 h-5 transition-colors ${
                          isActive ? "text-[#FED524]" : "text-white/60 group-hover:text-[#FED524]"
                        }`}
                      />
                    )}
                    <span className="font-semibold">{item.name}</span>
                  </Link>
                );
              })}
            </nav>
          </div>
        </div>
        
        {/* Footer Area inside Sidebar */}
        <div className="p-4 border-t border-white/10 text-center shrink-0">
          <p className="text-[10px] text-white/40 font-semibold tracking-wider uppercase">Universitas Siber Asia</p>
        </div>
      </aside>

      {/* Main Content */}
      <main className={`pt-16 min-h-screen transition-all duration-350 ${sidebarOpen ? "ml-64" : "ml-0"}`}>
        <div className="p-8 max-w-[1600px] mx-auto animate-fade-in">
          {children}
        </div>
      </main>
    </div>
  );
}
