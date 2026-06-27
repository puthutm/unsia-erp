"use client";

import { useAuth } from "@/contexts/auth-context";
import Link from "next/link";

// Module cards for the portal
const moduleCards = [
  { id: "pmb", name: "PMB", description: "Penerimaan Mahasiswa Baru", href: "/dashboard", icon: "🎓", bg: "bg-gradient-to-br from-purple-600 to-purple-700", hover: "hover:from-purple-700 hover:to-purple-800" },
  { id: "finance", name: "Keuangan", description: "Manajemen Keuangan Kampus", href: "/finance", icon: "💰", bg: "bg-gradient-to-br from-emerald-600 to-emerald-700", hover: "hover:from-emerald-700 hover:to-emerald-800" },
  { id: "academic", name: "Akademik", description: "Kelola Kegiatan Akademik", href: "/academic", icon: "📚", bg: "bg-gradient-to-br from-blue-600 to-blue-700", hover: "hover:from-blue-700 hover:to-blue-800" },
  { id: "lms", name: "LMS", description: "Learning Management System", href: "/lms", icon: "💻", bg: "bg-gradient-to-br from-orange-600 to-orange-700", hover: "hover:from-orange-700 hover:to-orange-800" },
];

// Stats card data
const stats = [
  { label: "Total Pendaftar", value: "1,234", change: "+12%", color: "bg-purple-100 text-purple-600", icon: "📋" },
  { label: "Pembayaran Lunas", value: "890", change: "+8%", color: "bg-green-100 text-green-600", icon: "✅" },
  { label: "Menunggu Verifikasi", value: "156", change: "-3%", color: "bg-yellow-100 text-yellow-600", icon: "⏳" },
  { label: "Calon Mahasiswa Baru", value: "567", change: "+15%", color: "bg-blue-100 text-blue-600", icon: "🎉" },
];

// Recent applicants
const recentApplicants = [
  { id: 1, name: "Ahmad Fauzi", prodi: "Teknik Informatika", status: "verified", date: "2024-01-15" },
  { id: 2, name: "Siti Rahayu", prodi: "Manajemen", status: "pending", date: "2024-01-15" },
  { id: 3, name: "Budi Santoso", prodi: "Akuntansi", status: "verified", date: "2024-01-14" },
  { id: 4, name: "Diana Putri", prodi: "Hukum", status: "rejected", date: "2024-01-14" },
  { id: 5, name: "Eko Prasetyo", prodi: "Teknik Elektro", status: "pending", date: "2024-01-13" },
];

const statusColors = {
  verified: "bg-green-100 text-green-800",
  pending: "bg-yellow-100 text-yellow-800",
  rejected: "bg-red-100 text-red-800",
};

export default function DashboardPage() {
  const { user } = useAuth();

  return (
    <div className="space-y-6">
      {/* Welcome Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">
            Dashboard PMB
          </h1>
          <p className="text-gray-500 mt-1">
            Selamat datang, {user?.name || "Admin"}. Berikut ringkasan aktivitas hari ini.
          </p>
        </div>
        <div className="flex items-center gap-2">
          <span className="px-3 py-1 bg-blue-100 text-blue-800 rounded-full text-sm font-medium">
            Gelombang 1 - 2024
          </span>
        </div>
      </div>

      {/* Module Cards Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        {moduleCards.map((mod) => (
          <Link
            key={mod.id}
            href={mod.href}
            className={`${mod.bg} ${mod.hover} rounded-xl p-5 text-white transition-all transform hover:scale-105 hover:shadow-lg`}
          >
            <div className="text-4xl mb-3">{mod.icon}</div>
            <h3 className="text-lg font-semibold">{mod.name}</h3>
            <p className="text-sm opacity-90 mt-1">{mod.description}</p>
          </Link>
        ))}
      </div>

{/* Stats Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        {stats.map((stat, index) => (
          <div key={index} className="bg-white rounded-xl shadow-sm border border-gray-200 p-5 hover:shadow-md transition-shadow">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-gray-500">{stat.label}</p>
                <p className="text-2xl font-bold text-gray-900 mt-1">{stat.value}</p>
              </div>
              <div className={`w-12 h-12 ${stat.color} rounded-xl flex items-center justify-center text-2xl`}>
                {stat.icon}
              </div>
            </div>
            <div className="mt-3 flex items-center text-sm">
              <span className={`font-medium ${stat.change.startsWith("+") ? "text-green-600" : "text-red-600"}`}>
                {stat.change}
              </span>
              <span className="text-gray-500 ml-2">dari bulan lalu</span>
            </div>
          </div>
        ))}
      </div>

{/* Quick Actions */}
      <div className="bg-white rounded-xl shadow-sm border border-gray-200 p-5">
        <h2 className="text-lg font-semibold text-gray-900 mb-4">Aksi Cepat</h2>
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
<button className="flex flex-col items-center gap-2 p-4 rounded-xl border border-gray-200 hover:border-blue-300 hover:bg-blue-50 transition-colors">
            <div className="w-10 h-10 bg-blue-100 rounded-full flex items-center justify-center">
              <svg className="w-5 h-5 text-blue-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M18 9v3m0 0v3m0-3h3m-3 0h-3m-2-5a4 4 0 11-8 0 4 4 0 018 0zM3 20a6 6 0 0112 0v1H8v-1a4 4 0 00-5-3.32M12 14v3m-3-3h3" />
              </svg>
            </div>
            <span className="text-sm font-medium text-gray-700">Tambah Pendaftar</span>
          </button>
          <button className="flex flex-col items-center gap-2 p-4 rounded-xl border border-gray-200 hover:border-green-300 hover:bg-green-50 transition-colors">
            <div className="w-10 h-10 bg-green-100 rounded-full flex items-center justify-center">
              <svg className="w-5 h-5 text-green-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </div>
            <span className="text-sm font-medium text-gray-700">Verifikasi Berkas</span>
          </button>
          <button className="flex flex-col items-center gap-2 p-4 rounded-xl border border-gray-200 hover:border-purple-300 hover:bg-purple-50 transition-colors">
            <div className="w-10 h-10 bg-purple-100 rounded-full flex items-center justify-center">
              <svg className="w-5 h-5 text-purple-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
              </svg>
            </div>
            <span className="text-sm font-medium text-gray-700">Kelola Gelombang</span>
          </button>
          <button className="flex flex-col items-center gap-2 p-4 rounded-xl border border-gray-200 hover:border-orange-300 hover:bg-orange-50 transition-colors">
            <div className="w-10 h-10 bg-orange-100 rounded-full flex items-center justify-center">
              <svg className="w-5 h-5 text-orange-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 17v-2m3 2v-4m3 4v-6m2 10H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
              </svg>
            </div>
            <span className="text-sm font-medium text-gray-700">Cetak Laporan</span>
          </button>
        </div>
      </div>

{/* Recent Applicants Table */}
      <div className="bg-white rounded-xl shadow-sm border border-gray-200 overflow-hidden">
        <div className="p-5 border-b border-gray-200">
          <h2 className="text-lg font-semibold text-gray-900">Pendaftar Terbaru</h2>
        </div>
        <table className="w-full">
          <thead className="bg-gray-50">
            <tr>
              <th className="px-5 py-3 text-left text-xs font-medium text-gray-500 uppercase">Nama</th>
              <th className="px-5 py-3 text-left text-xs font-medium text-gray-500 uppercase">Prodi</th>
              <th className="px-5 py-3 text-left text-xs font-medium text-gray-500 uppercase">Status</th>
              <th className="px-5 py-3 text-left text-xs font-medium text-gray-500 uppercase">Tanggal</th>
              <th className="px-5 py-3 text-right text-xs font-medium text-gray-500 uppercase">Aksi</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-200">
            {recentApplicants.map((applicant) => (
              <tr key={applicant.id} className="hover:bg-gray-50">
                <td className="px-5 py-4">
                  <div className="flex items-center gap-3">
                    <div className="w-8 h-8 bg-gray-200 rounded-full flex items-center justify-center">
                      <span className="text-sm font-medium text-gray-600">
                        {applicant.name.charAt(0)}
                      </span>
                    </div>
                    <span className="font-medium text-gray-900">{applicant.name}</span>
                  </div>
                </td>
                <td className="px-5 py-4 text-gray-600">{applicant.prodi}</td>
                <td className="px-5 py-4">
                  <span className={`px-2 py-1 rounded-full text-xs font-medium ${statusColors[applicant.status as keyof typeof statusColors]}`}>
                    {applicant.status === "verified" && "Terverifikasi"}
                    {applicant.status === "pending" && "Menunggu"}
                    {applicant.status === "rejected" && "Ditolak"}
                  </span>
                </td>
                <td className="px-5 py-4 text-gray-600">{applicant.date}</td>
                <td className="px-5 py-4 text-right">
                  <button className="text-blue-600 hover:text-blue-800 font-medium text-sm">
                    Lihat
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
        <div className="p-4 border-t border-gray-200">
          <button className="text-blue-600 hover:text-blue-800 font-medium text-sm">
            Lihat Semua Pendaftar →
          </button>
        </div>
      </div>
    </div>
  );
}
