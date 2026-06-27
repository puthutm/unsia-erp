"use client";

import { useState, useEffect } from "react";

// Laporan Page - Matches: UI/AKADEMIK/ADMIN/panel-laporan
// Laporan Akademik

export default function LaporanPage() {
  const [loading, setLoading] = useState(false);

  const handleGenerate = (reportType: string) => {
    setLoading(true);
    setTimeout(() => {
      setLoading(false);
    }, 1500);
  };

  return (
    <div className="p-6 lg:p-8 space-y-6">
      {/* Header */}
      <div className="bg-gradient-to-br from-brand-700 via-brand-600 to-brand-500 rounded-2xl p-6 text-white">
        <p className="text-[10px] uppercase tracking-widest font-bold text-brand-accent">Analitik</p>
        <h2 className="font-display font-black text-2xl mt-1">Laporan Akademik</h2>
        <p className="text-brand-50 text-sm mt-1.5">
          Laporan periodik dan ad-hoc untuk feeder dikti, akreditasi BAN-PT, dan internal management.
        </p>
      </div>

      {/* Report Buttons Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        <button
          onClick={() => handleGenerate("forlap")}
          className="bg-white border border-slate-200 hover:shadow-card rounded-2xl p-5 text-left group transition-all"
        >
          <div className="w-12 h-12 rounded-xl bg-brand-100 text-brand-600 flex items-center justify-center mb-4 group-hover:bg-brand-600 group-hover:text-white transition-all">
            <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 7v10c0 2.21 3.582 4 8 4s8-1.79 8-4V7M4 7c0 2.21 3.582 4 8 4s8-1.79 8-4M4 7c0-2.21 3.582-4 8-4s8 1.79 8 4m0 5c0 2.21-3.582 4-8 4s-8-1.79-8-4" />
            </svg>
          </div>
          <h3 className="font-bold text-slate-800 text-sm mt-3">Laporan Forlap Dikti</h3>
          <p className="text-[11px] text-slate-500 mt-1">Sync ke PDDikti per semester</p>
        </button>

        <button
          onClick={() => handleGenerate("khs")}
          className="bg-white border border-slate-200 hover:shadow-card rounded-2xl p-5 text-left group transition-all"
        >
          <div className="w-12 h-12 rounded-xl bg-emerald-100 text-emerald-600 flex items-center justify-center mb-4 group-hover:bg-emerald-600 group-hover:text-white transition-all">
            <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 17v-2m3 2v-4m3 4v-6m2 10H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
            </svg>
          </div>
          <h3 className="font-bold text-slate-800 text-sm mt-3">KHS & Transkrip Massal</h3>
          <p className="text-[11px] text-slate-500 mt-1">Generate per angkatan/prodi</p>
        </button>

        <button
          onClick={() => handleGenerate("beban")}
          className="bg-white border border-slate-200 hover:shadow-card rounded-2xl p-5 text-left group transition-all"
        >
          <div className="w-12 h-12 rounded-xl bg-violet-100 text-violet-600 flex items-center justify-center mb-4 group-hover:bg-violet-600 group-hover:text-white transition-all">
            <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
            </svg>
          </div>
          <h3 className="font-bold text-slate-800 text-sm mt-3">Beban Mengajar Dosen</h3>
          <p className="text-[11px] text-slate-500 mt-1">Distribusi SKS & jam/minggu</p>
        </button>

        <button
          onClick={() => handleGenerate("ekd")}
          className="bg-white border border-slate-200 hover:shadow-card rounded-2xl p-5 text-left group transition-all"
        >
          <div className="w-12 h-12 rounded-xl bg-amber-100 text-amber-600 flex items-center justify-center mb-4 group-hover:bg-amber-600 group-hover:text-white transition-all">
            <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 3v4M3 5h4M6 17v4m-2-2h4m5-16l2.286 6.857L21 12l-5.714 2.143L13 21l-2.286-6.857L5 12l5.714-2.143L13 3z" />
            </svg>
          </div>
          <h3 className="font-bold text-slate-800 text-sm mt-3">EKD Dosen</h3>
          <p className="text-[11px] text-slate-500 mt-1">Evaluasi Kinerja Dosen per semester</p>
        </button>

        <button
          onClick={() => handleGenerate("banpt")}
          className="bg-white border border-slate-200 hover:shadow-card rounded-2xl p-5 text-left group transition-all"
        >
          <div className="w-12 h-12 rounded-xl bg-blue-100 text-blue-600 flex items-center justify-center mb-4 group-hover:bg-blue-600 group-hover:text-white transition-all">
            <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4M7.835 4.697a3.42 3.42 0 001.241-.227 3.42 3.42 0 014.438 0 3.42 3.42 0 001.241.227 3.42 3.42 0 012.812 2.812c0 .696-.252 1.365-.697 1.873m0 3.745a3.42 3.42 0 001.241.227 3.42 3.42 0 014.438 0 3.42 3.42 0 001.241.227 3.42 3.42 0 012.812 2.812 3.42 3.42 0 01-.697 1.873m-7.939 2.873A3.42 3.42 0 0112 21c-.88 0-1.685-.224-2.395-.599m7.939-2.873A3.42 3.42 0 0112 18c.88 0 1.685.224 2.395.599m-2.395 0a3.42 3.42 0 01-2.395-.599M12 12a3.42 3.42 0 00-2.812 2.812c0 .696.252 1.365.697 1.873" />
            </svg>
          </div>
          <h3 className="font-bold text-slate-800 text-sm mt-3">Akreditasi BAN-PT</h3>
          <p className="text-[11px] text-slate-500 mt-1">Data 9 kriteria APS</p>
        </button>

        <button
          onClick={() => handleGenerate("yudisium")}
          className="bg-white border border-slate-200 hover:shadow-card rounded-2xl p-5 text-left group transition-all"
        >
          <div className="w-12 h-12 rounded-xl bg-rose-100 text-rose-600 flex items-center justify-center mb-4 group-hover:bg-rose-600 group-hover:text-white transition-all">
            <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-6 9l2 2 4-4" />
            </svg>
          </div>
          <h3 className="font-bold text-slate-800 text-sm mt-3">Daftar Yudisium</h3>
          <p className="text-[11px] text-slate-500 mt-1">Mhs siap wisuda & transkrip</p>
        </button>
      </div>

      {/* Loading Overlay */}
      {loading && (
        <div className="fixed inset-0 bg-slate-900/50 flex items-center justify-center z-50">
          <div className="bg-white rounded-2xl p-8 text-center">
            <div className="w-12 h-12 border-4 border-brand-600 border-t-transparent rounded-full animate-spin mx-auto mb-4"></div>
            <p className="text-slate-800 font-bold">Membuat Laporan...</p>
          </div>
        </div>
      )}
    </div>
  );
}
