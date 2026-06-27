"use client";

import { useState, useEffect } from "react";

// Admin Page - Next.js
// Matches: UI/AKADEMIK/ADMIN/

interface Stat {
  label: string;
  value: number;
  change: number;
  color: string;
}

export default function AdminPage() {
  const [stats, setStats] = useState<Stat[]>([]);
  const [loading, setLoading] = useState(true);
  const [activeModule, setActiveModule] = useState("dashboard");

  useEffect(() => {
    fetchStats();
  }, []);

  const fetchStats = async () => {
    setLoading(true);
    try {
      const token = localStorage.getItem("accessToken");
      const response = await fetch("/api/v1/academic/dashboard", {
        headers: { Authorization: `Bearer ${token}` },
      });
      if (response.ok) {
        const data = await response.json();
        setStats(data.stats || getDefaultStats());
      } else {
        setStats(getDefaultStats());
      }
    } catch (error) {
      console.error("Error fetching stats:", error);
      setStats(getDefaultStats());
    } finally {
      setLoading(false);
    }
  };

  const getDefaultStats = () => [
    { label: "Total Mahasiswa", value: 1250, change: 5.2, color: "blue" },
    { label: "Total Dosen", value: 85, change: 2.1, color: "green" },
    { label: "Mata Kuliah", value: 156, change: -1.3, color: "yellow" },
    { label: "Program Studi", value: 12, change: 0, color: "purple" },
  ];

  const modules = [
    { id: "dashboard", name: "Dashboard", icon: "📊" },
    { id: "mahasiswa", name: "Mahasiswa", icon: "🎓" },
    { id: "dosen", name: "Dosen", icon: "👨‍🏫" },
    { id: "matakuliah", name: "Mata Kuliah", icon: "📚" },
    { id: "krs", name: "KRS", icon: "📋" },
    { id: "nilai", name: "Nilai", icon: "📝" },
    { id: "jadwal", name: "Jadwal", icon: "🗓️" },
    { id: "absensi", name: "Absensi", icon: "✅" },
    { id: "lulus", name: " kelulusan", icon: "🎓" },
    { id: "laporan", name: "Laporan", icon: "📊" },
  ];

  const getChangeColor = (change: number) => {
    if (change > 0) return "text-green-600";
    if (change < 0) return "text-red-600";
    return "text-gray-600";
  };

  return (
    <div className="p-6 space-y-6">
      {/* Header */}
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-2xl font-bold text-slate-900">Admin Akademik</h1>
          <p className="text-slate-500 mt-1">Panel 管理 Akademik UNSIA</p>
        </div>
        <div className="flex items-center gap-4">
          <select className="px-3 py-2 border border-slate-200 rounded-lg">
            <option>Tahun Akademik 2024/2025</option>
            <option>Tahun Akademik 2023/2024</option>
          </select>
          <select className="px-3 py-2 border border-slate-200 rounded-lg">
            <option>Semester Ganjil</option>
            <option>Semester Genap</option>
          </select>
        </div>
      </div>

      {/* Module Tabs */}
      <div className="flex gap-2 flex-wrap">
        {modules.map((module) => (
          <button
            key={module.id}
            onClick={() => setActiveModule(module.id)}
            className={`px-4 py-2 rounded-lg font-medium flex items-center gap-2 ${
              activeModule === module.id
                ? "bg-blue-600 text-white"
                : "bg-slate-100 text-slate-600 hover:bg-slate-200"
            }`}
          >
            <span>{module.icon}</span>
            <span>{module.name}</span>
          </button>
        ))}
      </div>

      {/* Stats Grid */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        {loading ? (
          <div className="col-span-4 text-center text-slate-500 py-8">Memuat...</div>
        ) : (
          stats.map((stat, index) => (
            <div key={index} className="bg-white rounded-xl p-6 border border-slate-200">
              <h3 className="text-sm font-medium text-slate-500">{stat.label}</h3>
              <p className="text-3xl font-bold text-slate-900 mt-2">{stat.value}</p>
              <p className={`text-sm mt-2 ${getChangeColor(stat.change)}`}>
                {stat.change > 0 ? `+${stat.change}%` : stat.change}% dari semester lalu
              </p>
            </div>
          ))
        )}
      </div>

      {/* Quick Actions */}
      <div className="bg-white rounded-xl border border-slate-200 p-6">
        <h2 className="text-lg font-semibold text-slate-900 mb-4">Aksi Cepat</h2>
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
          <button className="p-4 bg-blue-50 rounded-lg text-blue-600 hover:bg-blue-100 text-left">
            <div className="text-2xl mb-2">🎓</div>
            <div className="font-medium">Generate NIM</div>
            <div className="text-sm">Dari Pendaftaran</div>
          </button>
          <button className="p-4 bg-green-50 rounded-lg text-green-600 hover:bg-green-100 text-left">
            <div className="text-2xl mb-2">📋</div>
            <div className="font-medium">Buat KRS</div>
            <div className="text-sm">Kartu Rencana Studi</div>
          </button>
          <button className="p-4 bg-yellow-50 rounded-lg text-yellow-600 hover:bg-yellow-100 text-left">
            <div className="text-2xl mb-2">📊</div>
            <div className="font-medium">Input Nilai</div>
            <div className="text-sm">Entri Nilai Mahasiswa</div>
          </button>
          <button className="p-4 bg-purple-50 rounded-lg text-purple-600 hover:bg-purple-100 text-left">
            <div className="text-2xl mb-2">🗓️</div>
            <div className="font-medium">Buat Jadwal</div>
            <div className="text-sm">Jadwal Kuliah</div>
          </button>
        </div>
      </div>

      {/* Recent Activity */}
      <div className="bg-white rounded-xl border border-slate-200 p-6">
        <h2 className="text-lg font-semibold text-slate-900 mb-4">Aktivitas Terkini</h2>
        <div className="space-y-3">
          <div className="flex items-center justify-between p-3 bg-slate-50 rounded-lg">
            <div className="flex items-center gap-3">
              <span className="text-2xl">✅</span>
              <div>
                <p className="font-medium text-slate-900">Absensi dikumpulkan</p>
                <p className="text-sm text-slate-500">Kelas TI-2024-01 - Struktur Data</p>
              </div>
            </div>
            <span className="text-sm text-slate-500">2 jam lalu</span>
          </div>
          <div className="flex items-center justify-between p-3 bg-slate-50 rounded-lg">
            <div className="flex items-center gap-3">
              <span className="text-2xl">📝</span>
              <div>
                <p className="font-medium text-slate-900">Nilai akhir diinput</p>
                <p className="text-sm text-slate-500">Kelas SI-2024-02 - Basis Data</p>
              </div>
            </div>
            <span className="text-sm text-slate-500">5 jam lalu</span>
          </div>
          <div className="flex items-center justify-between p-3 bg-slate-50 rounded-lg">
            <div className="flex items-center gap-3">
              <span className="text-2xl">🎓</span>
              <div>
                <p className="font-medium text-slate-900">Mahasiswa lulus</p>
                <p className="text-sm text-slate-500">12 mahasiswa teknik informatika</p>
              </div>
            </div>
            <span className="text-sm text-slate-500">1 hari lalu</span>
          </div>
        </div>
      </div>
    </div>
  );
}
