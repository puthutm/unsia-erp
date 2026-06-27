"use client";

import { useState, useEffect } from "react";

// Persuratan Page - Matches: UI/AKADEMIK/ADMIN/panel-persuratan
// Layanan persuratan akademik

interface LetterRequest {
  id: string;
  type: string;
  applicant: string;
  purpose: string;
  requestDate: string;
  status: "Diproses" | "Disetujui" | "Menunggu TTD" | "Selesai" | "Ditolak";
  priority: "Normal" | "Tinggi";
}

export default function PersuratanPage() {
  const [letters, setLetters] = useState<LetterRequest[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchLetters();
  }, []);

  const fetchLetters = async () => {
    setLoading(true);
    try {
      const token = localStorage.getItem("accessToken");
      const response = await fetch("/api/v1/academic/letter-requests", {
        headers: { Authorization: `Bearer ${token}` },
      });
      if (response.ok) {
        const data = await response.json();
        setLetters(data.data || getDefaultLetters());
      } else {
        setLetters(getDefaultLetters());
      }
    } catch (error) {
      console.error("Error fetching letters:", error);
      setLetters(getDefaultLetters());
    } finally {
      setLoading(false);
    }
  };

  const getDefaultLetters = (): LetterRequest[] => [
    { id: "SRT-2026-0512", type: "Surat Aktif Kuliah", applicant: "Budi Santoso (26090182)", purpose: "Bank BCA — Pengajuan Beasiswa", requestDate: "20 Mei 2026", status: "Diproses", priority: "Normal" },
    { id: "SRT-2026-0511", type: "Surat Keterangan Aktif", applicant: "Reni Aprilia (26110045)", purpose: "Kantor Imigrasi — Visa Studi", requestDate: "19 Mei 2026", status: "Selesai", priority: "Tinggi" },
    { id: "SRT-2026-0510", type: "Surat Cuti Akademik", applicant: "Lina Maharani (24070033)", purpose: "BAAK — Pengajuan Cuti", requestDate: "18 Mei 2026", status: "Disetujui", priority: "Normal" },
    { id: "SRT-2026-0509", type: "Surat Rekomendasi", applicant: "Daffa Hidayat (26020891)", purpose: "Beasiswa LPDP", requestDate: "17 Mei 2026", status: "Menunggu TTD", priority: "Tinggi" },
    { id: "SRT-2026-0508", type: "Transkrip Nilai Sementara", applicant: "Sasmita Rahayu (25110078)", purpose: "Lampiran Magang", requestDate: "16 Mei 2026", status: "Diproses", priority: "Normal" },
    { id: "SRT-2026-0507", type: "Surat Keterangan Lulus", applicant: "Andika Pratama (22090111)", purpose: "PT Tokopedia — Lamar Kerja", requestDate: "15 Mei 2026", status: "Selesai", priority: "Normal" },
    { id: "SRT-2026-0506", type: "Surat Pengantar Penelitian", applicant: "Bagas Cahyono (24080052)", purpose: "Disdukcapil Jakarta Pusat", requestDate: "14 Mei 2026", status: "Ditolak", priority: "Normal" },
  ];

  const getStatusBadge = (status: string) => {
    const styles: Record<string, string> = {
      Selesai: "bg-emerald-100 text-emerald-800",
      Disetujui: "bg-emerald-100 text-emerald-800",
      Diproses: "bg-blue-100 text-blue-800",
      "Menunggu TTD": "bg-amber-100 text-amber-800",
      Ditolak: "bg-rose-100 text-rose-800",
    };
    return styles[status] || "bg-gray-100 text-gray-800";
  };

  const getPriorityBadge = (priority: string) => {
    return priority === "Tinggi" ? "bg-rose-100 text-rose-800" : "bg-slate-100 text-slate-800";
  };

  return (
    <div className="p-6 lg:p-8 space-y-6">
      {/* Header */}
      <div className="bg-gradient-to-br from-brand-700 via-brand-600 to-brand-500 rounded-2xl p-6 text-white relative overflow-hidden">
        <div className="absolute -right-12 -top-12 w-56 h-56 bg-white/5 rounded-full"></div>
        <div className="relative flex items-start justify-between gap-4 flex-wrap">
          <div className="flex-1 min-w-0">
            <div className="flex items-center gap-2 mb-2">
              <i className="ph-fill ph-files text-brand-accent"></i>
              <span className="text-[10px] uppercase tracking-widest font-bold text-brand-accent">Layanan</span>
            </div>
            <h2 className="font-display font-black text-2xl mt-1">Persuratan Akademik</h2>
            <p className="text-brand-50 text-sm mt-1.5">
              Request surat: aktif kuliah, cuti, rekomendasi, transkrip, pengantar penelitian. 7 surat menunggu proses.
            </p>
          </div>
          <button className="px-3 py-2 bg-brand-accent text-brand-900 hover:bg-yellow-400 rounded-lg text-xs font-bold flex items-center gap-1.5 shrink-0">
            <i className="ph-bold ph-plus-circle"></i> Buat Surat Baru
          </button>
        </div>
      </div>

      {/* KPI Stats */}
      <div className="grid grid-cols-2 lg:grid-cols-4 gap-3">
        <div className="kpi-tile rounded-2xl p-4">
          <p className="text-[10px] uppercase tracking-wider font-bold text-slate-500">Total Bulan Ini</p>
          <p className="font-display font-black text-2xl text-slate-800 mt-1">42</p>
          <p className="text-[10px] text-emerald-600 mt-1 font-bold">+8 vs bulan lalu</p>
        </div>
        <div className="kpi-tile rounded-2xl p-4">
          <p className="text-[10px] uppercase tracking-wider font-bold text-amber-700">Diproses</p>
          <p className="font-display font-black text-2xl text-amber-700 mt-1">7</p>
          <p className="text-[10px] text-amber-600 mt-1 font-bold">SLA avg 2.4 hari</p>
        </div>
        <div className="kpi-tile rounded-2xl p-4">
          <p className="text-[10px] uppercase tracking-wider font-bold text-blue-700">Menunggu TTD</p>
          <p className="font-display font-black text-2xl text-blue-700 mt-1">2</p>
          <p className="text-[10px] text-blue-600 mt-1 font-bold">Kabir BAAK</p>
        </div>
        <div className="kpi-tile rounded-2xl p-4">
          <p className="text-[10px] uppercase tracking-wider font-bold text-emerald-700">Selesai</p>
          <p className="font-display font-black text-2xl text-emerald-700 mt-1">33</p>
          <p className="text-[10px] text-emerald-600 mt-1 font-bold">100% terkirim email</p>
        </div>
      </div>

      {/* Templates Section */}
      <div className="bg-white rounded-2xl border border-slate-200 shadow-soft p-5">
        <h3 className="font-bold text-slate-800 text-sm mb-3 flex items-center gap-2">
          <i className="ph-fill ph-files text-violet-600"></i> Template Surat Tersedia
        </h3>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-2">
          <div className="p-3 bg-slate-50 rounded-lg border border-slate-200">
            <p className="text-xs font-bold text-slate-800">Surat Keterangan Aktif Kuliah</p>
            <p className="text-[10px] text-slate-500 mt-1">SLA: 1 hari · Auto · butuh TTD Kepala BAAK</p>
          </div>
          <div className="p-3 bg-slate-50 rounded-lg border border-slate-200">
            <p className="text-xs font-bold text-slate-800">Surat Cuti Akademik</p>
            <p className="text-[10px] text-slate-500 mt-1">SLA: 3 hari · Manual · butuh TTDWR I</p>
          </div>
          <div className="p-3 bg-slate-50 rounded-lg border border-slate-200">
            <p className="text-xs font-bold text-slate-800">Surat Rekomendasi (Beasiswa/Kerja)</p>
            <p className="text-[10px] text-slate-500 mt-1">SLA: 5 hari · Manual · butuh TTD Dekan</p>
          </div>
          <div className="p-3 bg-slate-50 rounded-lg border border-slate-200">
            <p className="text-xs font-bold text-slate-800">Surat Keterangan Lulus</p>
            <p className="text-[10px] text-slate-500 mt-1">SLA: 2 hari · Auto · butuh TTD Kepala BAAK</p>
          </div>
          <div className="p-3 bg-slate-50 rounded-lg border border-slate-200">
            <p className="text-xs font-bold text-slate-800">Surat Pengantar Penelitian</p>
            <p className="text-[10px] text-slate-500 mt-1">SLA: 3 hari · Manual · butuh TTD Kaprodi</p>
          </div>
          <div className="p-3 bg-slate-50 rounded-lg border border-slate-200">
            <p className="text-xs font-bold text-slate-800">Transkrip Nilai Sementara</p>
            <p className="text-[10px] text-slate-500 mt-1">SLA: 2 hari · Auto · butuh TTD Kepala BAAK</p>
          </div>
        </div>
      </div>

      {/* Request Table */}
      <div className="bg-white rounded-2xl border border-slate-200 shadow-soft overflow-hidden">
        <div className="px-5 py-3 border-b border-slate-200">
          <h3 className="font-bold text-slate-800 text-sm">Daftar Request Surat</h3>
        </div>
        <div className="overflow-x-auto">
          <table className="w-full text-sm">
            <thead className="bg-slate-50 border-b border-slate-200">
              <tr>
                <th className="text-left px-4 py-2.5 text-[10px] uppercase tracking-wider font-bold text-slate-500">ID</th>
                <th className="text-left px-4 py-2.5 text-[10px] uppercase tracking-wider font-bold text-slate-500">Jenis Surat</th>
                <th className="text-left px-4 py-2.5 text-[10px] uppercase tracking-wider font-bold text-slate-500">Pemohon</th>
                <th className="text-left px-4 py-2.5 text-[10px] uppercase tracking-wider font-bold text-slate-500">Tujuan</th>
                <th className="text-left px-4 py-2.5 text-[10px] uppercase tracking-wider font-bold text-slate-500">Tanggal</th>
                <th className="text-center px-4 py-2.5 text-[10px] uppercase tracking-wider font-bold text-slate-500">Prioritas</th>
                <th className="text-center px-4 py-2.5 text-[10px] uppercase tracking-wider font-bold text-slate-500">Status</th>
                <th className="text-center px-4 py-2.5 text-[10px] uppercase tracking-wider font-bold text-slate-500">Aksi</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-100">
              {loading ? (
                <tr>
                  <td colSpan={8} className="p-8 text-center text-slate-500">Memuat...</td>
                </tr>
              ) : letters.length === 0 ? (
                <tr>
                  <td colSpan={8} className="p-8 text-center text-slate-500">Belum ada request surat</td>
                </tr>
              ) : (
                letters.map((letter) => (
                  <tr key={letter.id} className="hover:bg-slate-50">
                    <td className="px-4 py-2.5 font-mono font-bold text-brand-600 text-[10px]">{letter.id}</td>
                    <td className="px-4 py-2.5 text-xs font-bold text-slate-800">{letter.type}</td>
                    <td className="px-4 py-2.5 text-xs text-slate-700">{letter.applicant}</td>
                    <td className="px-4 py-2.5 text-[11px] text-slate-600">{letter.purpose}</td>
                    <td className="px-4 py-2.5 text-xs text-slate-700">{letter.requestDate}</td>
                    <td className="px-4 py-2.5 text-center">
                      <span className={`px-2 py-1 rounded-full text-[10px] font-bold ${getPriorityBadge(letter.priority)}`}>
                        {letter.priority}
                      </span>
                    </td>
                    <td className="px-4 py-2.5 text-center">
                      <span className={`px-2 py-1 rounded-full text-[10px] font-bold ${getStatusBadge(letter.status)}`}>
                        {letter.status}
                      </span>
                    </td>
                    <td className="px-4 py-2.5 text-center">
                      <div className="flex justify-center gap-1">
                        {letter.status === "Menunggu TTD" && (
                          <button className="px-2 py-1 bg-emerald-600 hover:bg-emerald-700 text-white rounded text-[10px] font-bold">
                            <i className="ph-bold ph-signature"></i> TTD
                          </button>
                        )}
                        {letter.status === "Diproses" && (
                          <button className="px-2 py-1 bg-blue-600 hover:bg-blue-700 text-white rounded text-[10px] font-bold">
                            <i className="ph-bold ph-arrow-right"></i> Proses
                          </button>
                        )}
                        {(letter.status === "Selesai" || letter.status === "Disetujui") && (
                          <button className="px-2 py-1 bg-brand-600 hover:bg-brand-700 text-white rounded text-[10px] font-bold">
                            <i className="ph-bold ph-printer"></i> Cetak
                          </button>
                        )}
                        <button className="px-2 py-1 bg-white border border-slate-200 hover:bg-slate-50 text-slate-600 rounded text-[10px] font-bold">
                          <i className="ph-bold ph-paper-plane"></i>
                        </button>
                      </div>
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}
