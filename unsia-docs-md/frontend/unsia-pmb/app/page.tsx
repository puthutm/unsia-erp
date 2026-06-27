"use client";

import { useState, useEffect } from "react";
import { usePmb } from "../hooks/use-pmb";
import { useReference } from "@/contexts/reference-context";
import { useAuth } from "@/contexts/auth-context";

export default function PmbDashboardPage() {
  const { isAuthenticated } = useAuth();
  const { studyPrograms, pmbWaves, isLoading: refLoading } = useReference();
  const {
    applicants,
    waves,
    stats,
    isLoading: pmbLoading,
    fetchApplicants,
    fetchWaves,
    fetchStats,
    updateApplicantStatus,
  } = usePmb();

  const [activePanel, setActivePanel] = useState<
    "dashboard" | "monitoring" | "pendaftar" | "verifikasi" | "pembayaran" | "komunikasi" | "gelombang" | "pengaturan"
  >("dashboard");

  // Settings sub-panel state
  const [activeSettingSubTab, setActiveSettingSubTab] = useState<
    "admin" | "roles" | "prodi" | "biaya" | "surat" | "pddikti" | "audit"
  >("admin");

  const [searchQuery, setSearchQuery] = useState("");
  const [selectedWave, setSelectedWave] = useState("");
  const [selectedStatus, setSelectedStatus] = useState("");

  useEffect(() => {
    if (isAuthenticated) {
      fetchStats();
      fetchWaves();
      fetchApplicants();
    }
  }, [isAuthenticated]);

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    fetchApplicants({ waveId: selectedWave, status: selectedStatus, search: searchQuery });
  };

  const getStatusBadge = (status: string) => {
    const styles: Record<string, string> = {
      draft: "bg-slate-100 text-slate-700 border border-slate-200",
      submitted: "bg-blue-100 text-blue-800 border border-blue-200",
      verified: "bg-purple-100 text-purple-800 border border-purple-200",
      rejected: "bg-rose-100 text-rose-800 border border-rose-200",
      accepted: "bg-emerald-100 text-emerald-800 border border-emerald-200",
    };
    return styles[status.toLowerCase()] || "bg-gray-100 text-gray-800";
  };

  const getPaymentStatusBadge = (status: string) => {
    return status.toLowerCase() === "paid" || status.toLowerCase() === "lunas"
      ? "bg-green-100 text-green-800 border border-green-200"
      : "bg-amber-100 text-amber-800 border border-amber-200";
  };

  return (
    <div className="flex min-h-screen bg-slate-50">
      {/* Sidebar - Primary Brand Gradient */}
      <aside className="w-72 bg-gradient-to-b from-[#0f487b] to-[#0a3052] text-white flex flex-col shrink-0 shadow-xl">
        <div className="h-20 flex items-center px-6 border-b border-white/10 gap-3">
          <span className="text-xl font-bold font-display tracking-tight">UNSIA ERP</span>
          <span className="px-2 py-0.5 bg-[#FED524] text-[#0f487b] text-[9px] font-black uppercase tracking-wider rounded">PMB</span>
        </div>

        <nav className="flex-1 py-6 px-3 space-y-1.5 overflow-y-auto">
          <p className="px-3 text-[9px] font-bold text-white/50 uppercase tracking-widest mb-2">Operasional</p>
          
          <button
            onClick={() => setActivePanel("dashboard")}
            className={`w-full flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm font-semibold transition-all text-left ${
              activePanel === "dashboard" ? "bg-white/15 text-[#FED524]" : "text-white/70 hover:bg-white/10 hover:text-white"
            }`}
          >
            📊 Beranda PMB
          </button>

          <button
            onClick={() => setActivePanel("monitoring")}
            className={`w-full flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm font-semibold transition-all text-left ${
              activePanel === "monitoring" ? "bg-white/15 text-[#FED524]" : "text-white/70 hover:bg-white/10 hover:text-white"
            }`}
          >
            🔍 Monitoring & Funnel
          </button>

          <button
            onClick={() => setActivePanel("pendaftar")}
            className={`w-full flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm font-semibold transition-all text-left ${
              activePanel === "pendaftar" ? "bg-white/15 text-[#FED524]" : "text-white/70 hover:bg-white/10 hover:text-white"
            }`}
          >
            👥 Data Pendaftar
          </button>

          <button
            onClick={() => setActivePanel("verifikasi")}
            className={`w-full flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm font-semibold transition-all text-left ${
              activePanel === "verifikasi" ? "bg-white/15 text-[#FED524]" : "text-white/70 hover:bg-white/10 hover:text-white"
            }`}
          >
            ✅ Verifikasi Berkas
          </button>

          <button
            onClick={() => setActivePanel("pembayaran")}
            className={`w-full flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm font-semibold transition-all text-left ${
              activePanel === "pembayaran" ? "bg-white/15 text-[#FED524]" : "text-white/70 hover:bg-white/10 hover:text-white"
            }`}
          >
            💳 Pembayaran UKT
          </button>

          <button
            onClick={() => setActivePanel("komunikasi")}
            className={`w-full flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm font-semibold transition-all text-left ${
              activePanel === "komunikasi" ? "bg-white/15 text-[#FED524]" : "text-white/70 hover:bg-white/10 hover:text-white"
            }`}
          >
            ✉️ Komunikasi & WA
          </button>

          <p className="px-3 text-[9px] font-bold text-white/50 uppercase tracking-widest mb-2 mt-6 pt-4 border-t border-white/10">Konfigurasi</p>

          <button
            onClick={() => setActivePanel("gelombang")}
            className={`w-full flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm font-semibold transition-all text-left ${
              activePanel === "gelombang" ? "bg-white/15 text-[#FED524]" : "text-white/70 hover:bg-white/10 hover:text-white"
            }`}
          >
            📅 Gelombang & Kuota
          </button>

          <button
            onClick={() => setActivePanel("pengaturan")}
            className={`w-full flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm font-semibold transition-all text-left ${
              activePanel === "pengaturan" ? "bg-white/15 text-[#FED524]" : "text-white/70 hover:bg-white/10 hover:text-white"
            }`}
          >
            ⚙️ Referensi & Pengaturan
          </button>
        </nav>
      </aside>

      {/* Main Area */}
      <main className="flex-1 flex flex-col min-w-0 h-screen overflow-y-auto p-6 space-y-6">
        
        {/* Header Section */}
        <div className="flex justify-between items-center border-b border-slate-200 pb-5">
          <div>
            <h1 className="text-3xl font-display font-bold text-slate-800 tracking-tight">
              {activePanel === "dashboard" && "Dashboard PMB"}
              {activePanel === "monitoring" && "Monitoring PMB"}
              {activePanel === "pendaftar" && "Direktori Data Pendaftar"}
              {activePanel === "verifikasi" && "Verifikasi Dokumen Pendaftar"}
              {activePanel === "pembayaran" && "Konfirmasi Pembayaran"}
              {activePanel === "komunikasi" && "Pusat Komunikasi"}
              {activePanel === "gelombang" && "Konfigurasi Gelombang"}
              {activePanel === "pengaturan" && "Referensi Menu & Pengaturan Sistem"}
            </h1>
            <p className="text-sm text-slate-500 mt-1">Penerimaan Mahasiswa Baru Universitas Siber Asia</p>
          </div>
          <div className="text-sm text-slate-400 font-medium">Port Base: 8003</div>
        </div>

        {/* Dashboard Panel */}
        {activePanel === "dashboard" && (
          <div className="space-y-6">
            {/* Stats Cards */}
            <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
              <div className="bg-white rounded-2xl p-6 border border-slate-200 shadow-sm">
                <span className="text-xs font-bold text-slate-400 uppercase tracking-widest block">Total Pendaftar</span>
                <span className="text-3xl font-bold text-slate-800 block mt-2 font-mono">{stats?.totalApplicants || applicants.length}</span>
              </div>
              <div className="bg-white rounded-2xl p-6 border border-slate-200 shadow-sm">
                <span className="text-xs font-bold text-slate-400 uppercase tracking-widest block">Gelombang Aktif</span>
                <span className="text-3xl font-bold text-slate-800 block mt-2 font-mono">{stats?.activeWave || waves.filter(w => w.isActive).length}</span>
              </div>
              <div className="bg-white rounded-2xl p-6 border border-slate-200 shadow-sm">
                <span className="text-xs font-bold text-slate-400 uppercase tracking-widest block">Menunggu Pembayaran</span>
                <span className="text-3xl font-bold text-amber-600 block mt-2 font-mono">{stats?.pendingPayment || 12}</span>
              </div>
              <div className="bg-white rounded-2xl p-6 border border-slate-200 shadow-sm">
                <span className="text-xs font-bold text-slate-400 uppercase tracking-widest block">Diterima</span>
                <span className="text-3xl font-bold text-emerald-600 block mt-2 font-mono">{stats?.admitted || 48}</span>
              </div>
            </div>

            {/* Quick Chart Simulation */}
            <div className="bg-white rounded-2xl p-6 border border-slate-200 shadow-sm">
              <h3 className="font-bold text-slate-800 mb-4">Tren Pendaftaran Mingguan</h3>
              <div className="h-48 flex items-end gap-3 pt-6 border-b border-slate-100">
                <div className="flex-1 bg-slate-100 h-1/4 rounded-t relative group"><span className="absolute -top-6 left-1/2 -translate-x-1/2 text-xs font-mono font-bold hidden group-hover:block">150</span></div>
                <div className="flex-1 bg-slate-100 h-1/3 rounded-t relative group"><span className="absolute -top-6 left-1/2 -translate-x-1/2 text-xs font-mono font-bold hidden group-hover:block">220</span></div>
                <div className="flex-1 bg-slate-100 h-1/2 rounded-t relative group"><span className="absolute -top-6 left-1/2 -translate-x-1/2 text-xs font-mono font-bold hidden group-hover:block">380</span></div>
                <div className="flex-1 bg-blue-500 h-3/4 rounded-t relative group"><span className="absolute -top-6 left-1/2 -translate-x-1/2 text-xs font-mono font-bold text-blue-600">620</span></div>
                <div className="flex-1 bg-[#FED524] h-full rounded-t relative group"><span className="absolute -top-6 left-1/2 -translate-x-1/2 text-xs font-mono font-bold text-amber-600">950</span></div>
              </div>
              <div className="flex justify-between text-xs text-slate-400 font-bold uppercase mt-2">
                <span>Minggu 1</span>
                <span>Minggu 2</span>
                <span>Minggu 3</span>
                <span>Minggu 4</span>
                <span>Sekarang</span>
              </div>
            </div>
          </div>
        )}

        {/* Monitoring / Funnel Panel */}
        {activePanel === "monitoring" && (
          <div className="space-y-6">
            <h3 className="text-lg font-bold text-slate-800">Conversion Funnel PMB</h3>
            <div className="bg-white rounded-2xl border border-slate-200 shadow-sm p-6 space-y-5">
              <div>
                <div className="flex justify-between text-sm font-semibold mb-1">
                  <span>Isi Formulir (Pendaftar Baru)</span>
                  <span>100% (3,150)</span>
                </div>
                <div className="h-6 bg-blue-100 rounded-lg overflow-hidden">
                  <div className="h-full bg-blue-600" style={{ width: "100%" }}></div>
                </div>
              </div>
              <div>
                <div className="flex justify-between text-sm font-semibold mb-1">
                  <span>Upload Berkas Lengkap</span>
                  <span>82% (2,580)</span>
                </div>
                <div className="h-6 bg-purple-100 rounded-lg overflow-hidden">
                  <div className="h-full bg-purple-600" style={{ width: "82%" }}></div>
                </div>
              </div>
              <div>
                <div className="flex justify-between text-sm font-semibold mb-1">
                  <span>Pembayaran Tagihan Pendaftaran</span>
                  <span>65% (2,040)</span>
                </div>
                <div className="h-6 bg-amber-100 rounded-lg overflow-hidden">
                  <div className="h-full bg-amber-500" style={{ width: "65%" }}></div>
                </div>
              </div>
              <div>
                <div className="flex justify-between text-sm font-semibold mb-1">
                  <span>Lulus Seleksi & Generate NIM</span>
                  <span>45% (1,418)</span>
                </div>
                <div className="h-6 bg-emerald-100 rounded-lg overflow-hidden">
                  <div className="h-full bg-emerald-600" style={{ width: "45%" }}></div>
                </div>
              </div>
            </div>
          </div>
        )}

        {/* Data Pendaftar Panel */}
        {activePanel === "pendaftar" && (
          <div className="space-y-4">
            {/* Search and Filters */}
            <form onSubmit={handleSearch} className="grid grid-cols-1 md:grid-cols-3 gap-3 bg-slate-100 p-4 rounded-xl border border-slate-200">
              <input
                type="text"
                placeholder="Masukkan Nomor Pendaftaran atau Nama..."
                className="rounded-lg border border-slate-350 p-2.5 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 bg-white"
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
              />
              <select
                className="rounded-lg border border-slate-350 p-2.5 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 bg-white"
                value={selectedWave}
                onChange={(e) => setSelectedWave(e.target.value)}
              >
                <option value="">Semua Gelombang</option>
                {pmbWaves.map((w) => (
                  <option key={w.id} value={w.id}>{w.name}</option>
                ))}
              </select>
              <select
                className="rounded-lg border border-slate-350 p-2.5 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 bg-white"
                value={selectedStatus}
                onChange={(e) => setSelectedStatus(e.target.value)}
              >
                <option value="">Semua Status</option>
                <option value="draft">Draft</option>
                <option value="submitted">Submitted</option>
                <option value="verified">Verified</option>
                <option value="accepted">Accepted</option>
                <option value="rejected">Rejected</option>
              </select>
            </form>

            {/* Applicant Table */}
            <div className="bg-white rounded-2xl border border-slate-200 shadow-sm overflow-hidden">
              {pmbLoading ? (
                <div className="text-center py-12 text-slate-500">Memuat data pendaftar...</div>
              ) : applicants.length === 0 ? (
                <div className="text-center py-12 text-slate-500">Tidak ada pendaftar terdaftar.</div>
              ) : (
                <div className="overflow-x-auto">
                  <table className="w-full text-left">
                    <thead className="bg-slate-50 border-b border-slate-200">
                      <tr>
                        <th className="p-4 text-xs font-semibold uppercase tracking-wider text-slate-500">No. PMB</th>
                        <th className="p-4 text-xs font-semibold uppercase tracking-wider text-slate-500">Nama</th>
                        <th className="p-4 text-xs font-semibold uppercase tracking-wider text-slate-500">Prodi</th>
                        <th className="p-4 text-xs font-semibold uppercase tracking-wider text-slate-500">Status Seleksi</th>
                        <th className="p-4 text-xs font-semibold uppercase tracking-wider text-slate-500">Status Pembayaran</th>
                      </tr>
                    </thead>
                    <tbody className="divide-y divide-slate-100">
                      {applicants.map((app) => (
                        <tr key={app.id} className="hover:bg-slate-50 transition-colors">
                          <td className="p-4 text-sm font-bold text-slate-900 font-mono">{app.applicantNumber}</td>
                          <td className="p-4 text-sm text-slate-700">{app.name}</td>
                          <td className="p-4 text-sm text-slate-600">{app.studyProgramName}</td>
                          <td className="p-4 text-sm">
                            <span className={`px-2.5 py-1 rounded-full text-xs font-semibold ${getStatusBadge(app.status)}`}>
                              {app.status}
                            </span>
                          </td>
                          <td className="p-4 text-sm">
                            <span className={`px-2.5 py-1 rounded-full text-xs font-semibold ${getPaymentStatusBadge(app.paymentStatus)}`}>
                              {app.paymentStatus}
                            </span>
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              )}
            </div>
          </div>
        )}

        {/* Verifikasi Berkas Panel */}
        {activePanel === "verifikasi" && (
          <div className="space-y-4">
            <h3 className="font-bold text-slate-800 text-lg">Berkas Perlu Verifikasi</h3>
            <div className="bg-white rounded-2xl border border-slate-200 shadow-sm p-6 text-center text-slate-500 py-12">
              Tidak ada antrian berkas pendaftaran yang belum diverifikasi saat ini.
            </div>
          </div>
        )}

        {/* Pembayaran Panel */}
        {activePanel === "pembayaran" && (
          <div className="space-y-4">
            <h3 className="font-bold text-slate-800 text-lg">Konfirmasi Pembayaran Manual</h3>
            <div className="bg-white rounded-2xl border border-slate-200 shadow-sm p-6 text-center text-slate-500 py-12">
              Tidak ada pembayaran manual yang memerlukan verifikasi.
            </div>
          </div>
        )}

        {/* Komunikasi Panel */}
        {activePanel === "komunikasi" && (
          <div className="space-y-4">
            <h3 className="font-bold text-slate-800 text-lg">WhatsApp / Email Broadcast</h3>
            <div className="bg-white rounded-2xl border border-slate-200 shadow-sm p-6 space-y-4 max-w-xl">
              <div>
                <label className="block text-sm font-semibold text-slate-700 mb-1.5">Pilih Target Broadcast</label>
                <select className="w-full rounded-lg border border-slate-350 p-2.5 text-sm bg-white focus:ring-2 focus:ring-blue-500 focus:outline-none">
                  <option>Seluruh Calon Mahasiswa (Status: Draft)</option>
                  <option>Calon Mahasiswa (Status: Menunggu Pembayaran)</option>
                  <option>Calon Mahasiswa Diterima (Undangan Daftar Ulang)</option>
                </select>
              </div>
              <div>
                <label className="block text-sm font-semibold text-slate-700 mb-1.5">Template Pesan</label>
                <textarea rows={4} className="w-full rounded-lg border border-slate-350 p-2.5 text-sm focus:ring-2 focus:ring-blue-500 focus:outline-none" placeholder="Masukkan isi pesan WhatsApp..."></textarea>
              </div>
              <button className="px-5 py-2.5 bg-blue-600 text-white font-bold rounded-lg text-sm hover:bg-blue-700 transition-colors">
                Kirim Broadcast
              </button>
            </div>
          </div>
        )}

        {/* Gelombang Panel */}
        {activePanel === "gelombang" && (
          <div className="space-y-6">
            <div className="flex justify-between items-center">
              <h3 className="font-bold text-slate-800 text-lg">Gelombang Pendaftaran PMB</h3>
              <button className="px-4 py-2 bg-blue-600 text-white font-semibold rounded-lg text-xs hover:bg-blue-700">
                + Buat Gelombang
              </button>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              {pmbWaves.map((wave) => (
                <div key={wave.id} className="bg-white rounded-2xl border border-slate-200 shadow-sm p-6 space-y-3">
                  <div className="flex justify-between items-start">
                    <div>
                      <h4 className="font-bold text-slate-900 text-base">{wave.name}</h4>
                      <p className="text-xs text-slate-500 font-medium mt-0.5">Kode: {wave.code}</p>
                    </div>
                    <span className={`px-2.5 py-0.5 rounded-full text-xs font-bold ${
                      wave.isActive ? "bg-green-100 text-green-800" : "bg-slate-100 text-slate-800"
                    }`}>
                      {wave.isActive ? "Aktif" : "Non-Aktif"}
                    </span>
                  </div>
                  <div className="text-sm text-slate-600 flex justify-between">
                    <span>Tgl Pendaftaran:</span>
                    <span className="font-bold text-slate-700">
                      {wave.registrationStartAt ? new Date(wave.registrationStartAt).toLocaleDateString("id-ID") : "Belum diatur"} s/d {wave.registrationEndAt ? new Date(wave.registrationEndAt).toLocaleDateString("id-ID") : "Belum diatur"}
                    </span>
                  </div>
                </div>
              ))}
              {pmbWaves.length === 0 && (
                <div className="col-span-2 bg-white rounded-2xl border border-slate-200 shadow-sm p-6 text-center text-slate-500 py-12">
                  Belum ada gelombang PMB aktif terdaftar.
                </div>
              )}
            </div>
          </div>
        )}

        {/* Configuration & Referensi Menu Panel */}
        {activePanel === "pengaturan" && (
          <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
            
            {/* Referensi submenus sidebar (matching lines 1476-1502 from UI template) */}
            <div className="lg:col-span-1 bg-white rounded-2xl border border-slate-200 shadow-sm p-2 h-fit space-y-0.5">
              <span className="block px-3 text-[10px] font-bold text-slate-400 uppercase tracking-widest mb-2 mt-2">
                Submenu Referensi
              </span>

              <button
                onClick={() => setActiveSettingSubTab("admin")}
                className={`w-full text-left px-4 py-3 rounded-xl font-bold text-sm flex items-center gap-3 transition-colors ${
                  activeSettingSubTab === "admin" ? "bg-blue-50 text-blue-700" : "text-slate-700 hover:bg-slate-50"
                }`}
              >
                👥 Pengguna Admin
              </button>

              <button
                onClick={() => setActiveSettingSubTab("roles")}
                className={`w-full text-left px-4 py-3 rounded-xl font-bold text-sm flex items-center gap-3 transition-colors ${
                  activeSettingSubTab === "roles" ? "bg-blue-50 text-blue-700" : "text-slate-700 hover:bg-slate-50"
                }`}
              >
                🛡️ Peran & Akses
              </button>

              <button
                onClick={() => setActiveSettingSubTab("prodi")}
                className={`w-full text-left px-4 py-3 rounded-xl font-bold text-sm flex items-center gap-3 transition-colors ${
                  activeSettingSubTab === "prodi" ? "bg-blue-50 text-blue-700" : "text-slate-700 hover:bg-slate-50"
                }`}
              >
                🎓 Program Studi (SIAKAD)
              </button>

              <button
                onClick={() => setActiveSettingSubTab("biaya")}
                className={`w-full text-left px-4 py-3 rounded-xl font-bold text-sm flex items-center gap-3 transition-colors ${
                  activeSettingSubTab === "biaya" ? "bg-blue-50 text-blue-700" : "text-slate-700 hover:bg-slate-50"
                }`}
              >
                💵 Biaya & UKT
              </button>

              <button
                onClick={() => setActiveSettingSubTab("surat")}
                className={`w-full text-left px-4 py-3 rounded-xl font-bold text-sm flex items-center gap-3 transition-colors ${
                  activeSettingSubTab === "surat" ? "bg-blue-50 text-blue-700" : "text-slate-700 hover:bg-slate-50"
                }`}
              >
                📝 Template Surat
              </button>

              <button
                onClick={() => setActiveSettingSubTab("pddikti")}
                className={`w-full text-left px-4 py-3 rounded-xl font-bold text-sm flex items-center gap-3 transition-colors ${
                  activeSettingSubTab === "pddikti" ? "bg-blue-50 text-blue-700" : "text-slate-700 hover:bg-slate-50"
                }`}
              >
                🔌 Integrasi PDDikti
              </button>

              <button
                onClick={() => setActiveSettingSubTab("audit")}
                className={`w-full text-left px-4 py-3 rounded-xl font-bold text-sm flex items-center gap-3 transition-colors ${
                  activeSettingSubTab === "audit" ? "bg-blue-50 text-blue-700" : "text-slate-700 hover:bg-slate-50"
                }`}
              >
                📋 Audit Log Sistem
              </button>
            </div>

            {/* Referensi Sub-panel Content */}
            <div className="lg:col-span-2 bg-white rounded-2xl border border-slate-200 shadow-sm overflow-hidden p-6 space-y-4">
              
              {activeSettingSubTab === "admin" && (
                <div className="space-y-4">
                  <div className="flex justify-between items-center border-b border-slate-100 pb-3">
                    <div>
                      <h4 className="font-bold text-slate-800">Pengguna Admin PMB</h4>
                      <p className="text-xs text-slate-400 mt-0.5">Daftar staf administrator penerimaan mahasiswa baru</p>
                    </div>
                    <button className="px-3 py-1.5 bg-blue-600 text-white rounded-lg text-xs font-semibold hover:bg-blue-700">
                      + Tambah Admin
                    </button>
                  </div>
                  
                  <div className="divide-y divide-slate-100">
                    <div className="py-3 flex items-center gap-3">
                      <img src="https://ui-avatars.com/api/?name=Aris+Wijaya&background=FED524&color=0f487b&rounded=true&bold=true" class="w-10 h-10 rounded-full" />
                      <div>
                        <p className="font-bold text-sm text-slate-800">Aris Wijaya</p>
                        <p className="text-xs text-slate-400">Super Admin (aris.wijaya@unsia.ac.id)</p>
                      </div>
                      <span className="ml-auto text-xs text-emerald-600 font-bold bg-emerald-50 px-2 py-0.5 rounded">Online</span>
                    </div>
                    <div className="py-3 flex items-center gap-3">
                      <img src="https://ui-avatars.com/api/?name=Sari+Wulan&background=f0f4f8&color=0f487b&rounded=true&bold=true" class="w-10 h-10 rounded-full" />
                      <div>
                        <p className="font-bold text-sm text-slate-800">Sari Wulan</p>
                        <p className="text-xs text-slate-400">Verifikator (sari.wulan@unsia.ac.id)</p>
                      </div>
                      <span className="ml-auto text-xs text-slate-500 font-bold bg-slate-50 px-2 py-0.5 rounded">Offline</span>
                    </div>
                  </div>
                </div>
              )}

              {activeSettingSubTab === "roles" && (
                <div className="space-y-4">
                  <h4 className="font-bold text-slate-800 border-b border-slate-100 pb-3">Peran & Hak Akses Kontrol</h4>
                  <div className="space-y-3">
                    <div className="p-3 border border-slate-150 rounded-xl space-y-1">
                      <p className="font-bold text-sm text-slate-800">Super Admin PMB</p>
                      <p className="text-xs text-slate-500">Akses penuh ke konfigurasi gelombang, biaya, verifikasi, dan broadcast.</p>
                    </div>
                    <div className="p-3 border border-slate-150 rounded-xl space-y-1">
                      <p className="font-bold text-sm text-slate-800">Verifikator</p>
                      <p className="text-xs text-slate-500">Hanya memiliki akses verifikasi berkas pendaftaran dan approval kelulusan seleksi.</p>
                    </div>
                    <div className="p-3 border border-slate-150 rounded-xl space-y-1">
                      <p className="font-bold text-sm text-slate-800">Finance Staf</p>
                      <p className="text-xs text-slate-500">Hanya memiliki akses untuk memverifikasi pembayaran manual pendaftar.</p>
                    </div>
                  </div>
                </div>
              )}

              {activeSettingSubTab === "prodi" && (
                <div className="space-y-4">
                  <h4 className="font-bold text-slate-800 border-b border-slate-100 pb-3">Program Studi Terintegrasi SIAKAD</h4>
                  <div className="overflow-x-auto">
                    <table className="w-full text-left">
                      <thead className="bg-slate-50 border-b border-slate-200">
                        <tr>
                          <th className="p-3 text-xs font-semibold uppercase text-slate-500">Kode</th>
                          <th className="p-3 text-xs font-semibold uppercase text-slate-500">Program Studi</th>
                          <th className="p-3 text-xs font-semibold uppercase text-slate-500">Jenjang</th>
                        </tr>
                      </thead>
                      <tbody className="divide-y divide-slate-100">
                        {studyPrograms.map((p) => (
                          <tr key={p.id}>
                            <td className="p-3 text-sm font-bold text-slate-800 font-mono">{p.code}</td>
                            <td className="p-3 text-sm text-slate-700">{p.name}</td>
                            <td className="p-3 text-sm text-slate-500 font-semibold">{p.degree || "S1"}</td>
                          </tr>
                        ))}
                        {studyPrograms.length === 0 && (
                          <tr>
                            <td colSpan={3} className="p-3 text-sm text-slate-400 text-center py-6">Belum ada program studi aktif.</td>
                          </tr>
                        )}
                      </tbody>
                    </table>
                  </div>
                </div>
              )}

              {activeSettingSubTab === "biaya" && (
                <div className="space-y-4">
                  <h4 className="font-bold text-slate-800 border-b border-slate-100 pb-3">Konfigurasi UKT & Biaya Kuliah</h4>
                  <div className="divide-y divide-slate-150">
                    <div className="py-2.5 flex justify-between items-center text-sm">
                      <span className="font-semibold text-slate-700">Biaya Formulir Pendaftaran</span>
                      <span className="font-mono font-bold text-slate-900">Rp 250.000</span>
                    </div>
                    <div className="py-2.5 flex justify-between items-center text-sm">
                      <span className="font-semibold text-slate-700">UKT Informatika / S1</span>
                      <span className="font-mono font-bold text-slate-900">Rp 3.500.000</span>
                    </div>
                    <div className="py-2.5 flex justify-between items-center text-sm">
                      <span className="font-semibold text-slate-700">UKT Sistem Informasi / S1</span>
                      <span className="font-mono font-bold text-slate-900">Rp 3.300.000</span>
                    </div>
                    <div className="py-2.5 flex justify-between items-center text-sm">
                      <span className="font-semibold text-slate-700">UKT Manajemen / S1</span>
                      <span className="font-mono font-bold text-slate-900">Rp 3.000.000</span>
                    </div>
                  </div>
                </div>
              )}

              {activeSettingSubTab === "surat" && (
                <div className="space-y-4">
                  <h4 className="font-bold text-slate-800 border-b border-slate-100 pb-3">Template Surat & Email Korespondensi</h4>
                  <div className="space-y-3 text-sm">
                    <div className="p-3 border rounded-xl flex justify-between items-center">
                      <div>
                        <p className="font-bold text-slate-800">Email Verifikasi Pendaftaran</p>
                        <p className="text-xs text-slate-400">Dikirim saat calon mahasiswa mengunggah seluruh berkas</p>
                      </div>
                      <button className="text-xs text-blue-600 font-bold hover:underline">Edit</button>
                    </div>
                    <div className="p-3 border rounded-xl flex justify-between items-center">
                      <div>
                        <p className="font-bold text-slate-800">Surat Pengumuman Kelulusan (PDF)</p>
                        <p className="text-xs text-slate-400">Surat keterangan resmi lulus seleksi PMB</p>
                      </div>
                      <button className="text-xs text-blue-600 font-bold hover:underline">Edit</button>
                    </div>
                  </div>
                </div>
              )}

              {activeSettingSubTab === "pddikti" && (
                <div className="space-y-4">
                  <h4 className="font-bold text-slate-800 border-b border-slate-100 pb-3">Integrasi Sinkronisasi PDDikti Feeder</h4>
                  <div className="p-4 bg-blue-50 border border-blue-200 rounded-xl text-sm text-blue-900 flex justify-between items-center">
                    <div>
                      <strong>Status Koneksi WS Feeder:</strong>
                      <span className="block text-xs mt-0.5 text-blue-700">Tersambung (http://feeder.unsia.ac.id:8082)</span>
                    </div>
                    <button className="px-3 py-1.5 bg-blue-600 hover:bg-blue-700 text-white rounded-lg text-xs font-semibold">
                      Sinkronkan Sekarang
                    </button>
                  </div>
                </div>
              )}

              {activeSettingSubTab === "audit" && (
                <div className="space-y-4">
                  <h4 className="font-bold text-slate-800 border-b border-slate-100 pb-3">Log Histori Audit Sistem</h4>
                  <div className="space-y-2 text-xs font-mono text-slate-500">
                    <p>[2026-06-25 03:22] ADMIN 'aris.wijaya' membuat Gelombang 1 Ganjil aktif</p>
                    <p>[2026-06-25 02:10] SYSTEM verifikasi berkas pendaftar #PMB20260901 otomatis sukses</p>
                    <p>[2026-06-24 23:45] FINANCE verifikasi pembayaran invoice manual #INV991202 disetujui</p>
                  </div>
                </div>
              )}

            </div>
          </div>
        )}

      </main>
    </div>
  );
}
