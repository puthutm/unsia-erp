"use client";

import { useEffect } from "react";
import { usePmb } from "@/hooks";

export default function CommandCenterPage() {
  const { waves, stats, isLoading, fetchWaves, fetchStats } = usePmb();

  useEffect(() => {
    fetchWaves();
    fetchStats();
  }, [fetchWaves, fetchStats]);

  return (
    <div className="p-6 space-y-6">
      {/* Header */}
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-2xl font-bold text-slate-900">Command Center PMB</h1>
<p className="text-slate-500 mt-1">Pengawasan process PMB secara real-time</p>
        </div>
        <div className="flex gap-2">
          <button className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700">
            Refresh Data
          </button>
        </div>
      </div>

{/* Stats Cards */}
      <div className="grid grid-cols-1 md:grid-cols-5 gap-4">
        <div className="bg-white rounded-xl p-6 border border-slate-200">
          <h3 className="text-sm font-medium text-slate-500">Total Pendaftar</h3>
          <p className="text-3xl font-bold text-slate-900 mt-2">{stats?.totalApplicants || 0}</p>
        </div>
        <div className="bg-white rounded-xl p-6 border border-slate-200">
          <h3 className="text-sm font-medium text-slate-500">Verifikasi</h3>
          <p className="text-3xl font-bold text-blue-600 mt-2">{stats?.verifikasi || 0}</p>
        </div>
        <div className="bg-white rounded-xl p-6 border border-slate-200">
          <h3 className="text-sm font-medium text-slate-500">Seleksi</h3>
          <p className="text-3xl font-bold text-yellow-600 mt-2">{stats?.seleksi || 0}</p>
        </div>
        <div className="bg-white rounded-xl p-6 border border-slate-200">
          <h3 className="text-sm font-medium text-slate-500">Lulus</h3>
          <p className="text-3xl font-bold text-green-600 mt-2">{stats?.lulus || 0}</p>
        </div>
        <div className="bg-white rounded-xl p-6 border border-slate-200">
          <h3 className="text-sm font-medium text-slate-500">Gelombang Aktif</h3>
          <p className="text-3xl font-bold text-purple-600 mt-2">{stats?.activeWave || "-"}</p>
        </div>
      </div>

      {/* Wave Management */}
      <div className="bg-white rounded-xl border border-slate-200">
        <div className="p-6 border-b border-slate-200">
          <h2 className="text-lg font-semibold text-slate-900">Pengelolaan Gelombang</h2>
        </div>
        <div className="p-6">
          {isLoading ? (
            <div className="text-center text-slate-500 py-8">Memuat data...</div>
          ) : waves.length === 0 ? (
            <div className="text-center text-slate-500 py-8">Belum ada gelombang</div>
          ) : (
            <div className="space-y-4">
              {waves.map((wave) => (
                <div key={wave.id} className="flex items-center justify-between p-4 border border-slate-200 rounded-lg">
                  <div>
                    <h4 className="font-medium text-slate-900">{wave.name}</h4>
                    <p className="text-sm text-slate-500">
                      {new Date(wave.startDate).toLocaleDateString("id-ID")} - {new Date(wave.endDate).toLocaleDateString("id-ID")}
                    </p>
                  </div>
                  <div className="flex gap-2">
                    <span className={`px-3 py-1 rounded-full text-sm ${
                      wave.isActive ? "bg-green-100 text-green-800" : "bg-gray-100 text-gray-800"
                    }`}>
                      {wave.isActive ? "Aktif" : "Tidak Aktif"}
                    </span>
                    <button className="px-3 py-1 text-blue-600 border border-blue-600 rounded-lg hover:bg-blue-50">
                      Edit
                    </button>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>

      {/* Quick Actions */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div className="bg-white rounded-xl border border-slate-200 p-6">
          <h3 className="font-semibold text-slate-900 mb-4">Aksi Cepat</h3>
          <div className="space-y-2">
            <button className="w-full px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700">
              Kirim Email Massal
            </button>
            <button className="w-full px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700">
              Export Data Pendaftar
            </button>
            <button className="w-full px-4 py-2 bg-purple-600 text-white rounded-lg hover:bg-purple-700">
              Generate Laporan
            </button>
          </div>
        </div>
        <div className="bg-white rounded-xl border border-slate-200 p-6">
          <h3 className="font-semibold text-slate-900 mb-4">Aktivitas Terkini</h3>
          <div className="space-y-2 text-sm text-slate-600">
            <p>• Pendaftar baru: 5 orang</p>
            <p>• Verifikasi selesai: 12 orang</p>
            <p>• Seleksi akademik: 8 orang</p>
            <p>• Daftar ulang: 3 orang</p>
          </div>
        </div>
      </div>
    </div>
  );
}
