"use client";

import { useState, useEffect } from "react";

// Kurikulum Page - Matches: UI/AKADEMIK/ADMIN/panel-kurikulum
// Master kurikulum program studi

interface Curriculum {
  id: string;
  code: string;
  studyProgram: string;
  yearIssued: number;
  totalSks: number;
  totalCourses: number;
  mandatoryCourses: number;
  electiveCourses: number;
  faculty: string;
  status: "Aktif" | "Phase-out";
}

export default function CurriculumPage() {
  const [curricula, setCurricula] = useState<Curriculum[]>([]);
  const [loading, setLoading] = useState(true);
  const [showModal, setShowModal] = useState(false);

  useEffect(() => {
    fetchCurricula();
  }, []);

  const fetchCurricula = async () => {
    setLoading(true);
    try {
      const token = localStorage.getItem("accessToken");
      const response = await fetch("/api/v1/academic/curricula", {
        headers: { Authorization: `Bearer ${token}` },
      });
      if (response.ok) {
        const data = await response.json();
        setCurricula(data.data || getDefaultCurricula());
      } else {
        setCurricula(getDefaultCurricula());
      }
    } catch (error) {
      console.error("Error fetching curricula:", error);
      setCurricula(getDefaultCurricula());
    } finally {
      setLoading(false);
    }
  };

  const getDefaultCurricula = (): Curriculum[] => [
    { id: "KUR-2024-IF", code: "KUR-2024", studyProgram: "S1 Informatika", yearIssued: 2024, totalSks: 144, totalCourses: 56, mandatoryCourses: 42, electiveCourses: 14, faculty: "FTI", status: "Aktif" },
    { id: "KUR-2024-SI", code: "KUR-2024", studyProgram: "S1 Sistem Informasi", yearIssued: 2024, totalSks: 144, totalCourses: 54, mandatoryCourses: 40, electiveCourses: 14, faculty: "FTI", status: "Aktif" },
    { id: "KUR-2024-MJ", code: "KUR-2024", studyProgram: "S1 Manajemen", yearIssued: 2024, totalSks: 144, totalCourses: 52, mandatoryCourses: 38, electiveCourses: 14, faculty: "FEB", status: "Aktif" },
    { id: "KUR-2024-AK", code: "KUR-2024", studyProgram: "S1 Akuntansi", yearIssued: 2024, totalSks: 144, totalCourses: 52, mandatoryCourses: 38, electiveCourses: 14, faculty: "FEB", status: "Aktif" },
    { id: "KUR-2024-PSI", code: "KUR-2024", studyProgram: "S1 Psikologi", yearIssued: 2024, totalSks: 144, totalCourses: 50, mandatoryCourses: 36, electiveCourses: 14, faculty: "Psi", status: "Aktif" },
    { id: "KUR-2024-MM", code: "KUR-2024", studyProgram: "S2 Magister Manajemen", yearIssued: 2024, totalSks: 42, totalCourses: 14, mandatoryCourses: 10, electiveCourses: 4, faculty: "FEB", status: "Aktif" },
    { id: "KUR-2020-IF", code: "KUR-2020", studyProgram: "S1 Informatika", yearIssued: 2020, totalSks: 146, totalCourses: 58, mandatoryCourses: 44, electiveCourses: 14, faculty: "FTI", status: "Phase-out" },
  ];

  const getStatusBadge = (status: string) => {
    const styles: Record<string, string> = {
      Aktif: "bg-emerald-100 text-emerald-800",
      "Phase-out": "bg-amber-100 text-amber-800",
    };
    return styles[status] || "bg-gray-100 text-gray-800";
  };

  return (
    <div className="p-6 lg:p-8 space-y-6">
      {/* Header */}
      <div className="bg-gradient-to-br from-brand-700 via-brand-600 to-brand-500 rounded-2xl p-6 text-white relative overflow-hidden">
        <div className="absolute -right-12 -top-12 w-48 h-48 bg-white/5 rounded-full"></div>
        <div className="relative flex items-start justify-between flex-wrap gap-4">
          <div className="flex-1 min-w-0">
            <div className="flex items-center gap-2 mb-2">
              <i className="ph-fill ph-tree-structure text-brand-accent"></i>
              <span className="text-[10px] uppercase tracking-widest font-bold text-brand-accent">Master · Kurikulum Prodi</span>
            </div>
            <h2 className="font-display font-black text-2xl">Kurikulum Program Studi</h2>
            <p className="text-brand-50 text-sm mt-1.5">
              Master kurikulum per program studi. Tiap kurikulum berisi daftar mata kuliah wajib & pilihan beserta SKS dan distribusi semester.
            </p>
          </div>
          <button
            onClick={() => setShowModal(true)}
            className="px-3 py-2 bg-brand-accent text-brand-900 hover:bg-yellow-400 rounded-lg text-xs font-bold flex items-center gap-1.5 shrink-0"
          >
            <i className="ph-bold ph-plus-circle"></i> Tambah Kurikulum
          </button>
        </div>
      </div>

      {/* Curricula Grid */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
        {loading ? (
          <div className="col-span-2 text-center text-slate-500 py-8">Memuat...</div>
        ) : (
          curricula.map((curriculum) => (
            <div key={curriculum.id} className="bg-white border border-slate-200 rounded-2xl shadow-soft p-5 hover:shadow-card transition-shadow">
              <div className="flex items-start justify-between mb-3">
                <div className="flex-1 min-w-0">
                  <div className="flex items-center gap-2 mb-1">
                    <span className="px-2 py-1 bg-brand-100 text-brand-700 rounded text-[10px] font-bold">{curriculum.code}</span>
                    <span className={`px-2 py-1 rounded-full text-[10px] font-bold ${getStatusBadge(curriculum.status)}`}>
                      {curriculum.status}
                    </span>
                    <span className="text-[10px] text-slate-500">{curriculum.faculty}</span>
                  </div>
                  <h3 className="font-display font-bold text-base text-slate-800">{curriculum.studyProgram}</h3>
                  <p className="text-[11px] text-slate-500 mt-0.5">Terbit {curriculum.yearIssued} · {curriculum.id}</p>
                </div>
              </div>
              <div className="grid grid-cols-3 gap-2 mt-3">
                <div className="p-2 bg-slate-50 rounded text-center">
                  <p className="font-display font-black text-lg text-slate-800">{curriculum.totalSks}</p>
                  <p className="text-[9px] text-slate-500 font-bold uppercase">Total SKS</p>
                </div>
                <div className="p-2 bg-emerald-50 rounded text-center">
                  <p className="font-display font-black text-lg text-emerald-700">{curriculum.mandatoryCourses}</p>
                  <p className="text-[9px] text-emerald-600 font-bold uppercase">MK Wajib</p>
                </div>
                <div className="p-2 bg-violet-50 rounded text-center">
                  <p className="font-display font-black text-lg text-violet-700">{curriculum.electiveCourses}</p>
                  <p className="text-[9px] text-violet-600 font-bold uppercase">MK Pilihan</p>
                </div>
              </div>
              <div className="flex gap-1.5 mt-3">
                <button className="flex-1 px-2.5 py-1.5 bg-brand-50 hover:bg-brand-100 text-brand-700 text-[10px] font-bold rounded flex items-center justify-center gap-1">
                  <i className="ph-bold ph-list-checks"></i> Daftar MK
                </button>
                <button className="flex-1 px-2.5 py-1.5 bg-emerald-50 hover:bg-emerald-100 text-emerald-700 text-[10px] font-bold rounded flex items-center justify-center gap-1">
                  <i className="ph-bold ph-plus-circle"></i> Tambah MK
                </button>
                <button className="px-2.5 py-1.5 bg-white border border-slate-200 hover:bg-slate-50 text-slate-600 text-[10px] font-bold rounded">
                  <i className="ph-bold ph-file-xls"></i>
                </button>
              </div>
            </div>
          ))
        )}
      </div>

      {/* Add Modal */}
      {showModal && (
        <div className="fixed inset-0 bg-slate-900/70 z-50 flex items-center justify-center p-4">
          <div className="bg-white rounded-2xl shadow-2xl max-w-md w-full p-6">
            <h3 className="font-bold text-lg text-slate-800 mb-4">Tambah Kurikulum Baru</h3>
            <form className="space-y-4">
              <div className="grid grid-cols-2 gap-2">
                <div>
                  <label className="block text-xs font-bold text-slate-700 mb-1.5">Kode Kurikulum</label>
                  <input type="text" defaultValue="KUR-2027" className="w-full px-3 py-2 bg-slate-50 border border-slate-200 rounded-lg text-sm font-mono font-bold focus:outline-none focus:border-brand-600" />
                </div>
                <div>
                  <label className="block text-xs font-bold text-slate-700 mb-1.5">Tahun Terbit</label>
                  <input type="number" defaultValue="2027" className="w-full px-3 py-2 bg-slate-50 border border-slate-200 rounded-lg text-sm focus:outline-none focus:border-brand-600" />
                </div>
              </div>
              <div>
                <label className="block text-xs font-bold text-slate-700 mb-1.5">Program Studi</label>
                <select className="w-full px-3 py-2 bg-slate-50 border border-slate-200 rounded-lg text-sm focus:outline-none focus:border-brand-600">
                  <option>S1 Informatika</option>
                  <option>S1 Sistem Informasi</option>
                  <option>S1 Manajemen</option>
                  <option>S1 Akuntansi</option>
                  <option>S1 Psikologi</option>
                  <option>S2 Magister Manajemen</option>
                </select>
              </div>
              <div className="grid grid-cols-3 gap-2">
                <div>
                  <label className="block text-xs font-bold text-slate-700 mb-1.5">Total SKS</label>
                  <input type="number" defaultValue="144" className="w-full px-3 py-2 bg-slate-50 border border-slate-200 rounded-lg text-sm font-mono text-center focus:outline-none focus:border-brand-600" />
                </div>
                <div>
                  <label className="block text-xs font-bold text-slate-700 mb-1.5">MK Wajib</label>
                  <input type="number" defaultValue="42" className="w-full px-3 py-2 bg-slate-50 border border-slate-200 rounded-lg text-sm font-mono text-center focus:outline-none focus:border-brand-600" />
                </div>
                <div>
                  <label className="block text-xs font-bold text-slate-700 mb-1.5">MK Pilihan</label>
                  <input type="number" defaultValue="14" className="w-full px-3 py-2 bg-slate-50 border border-slate-200 rounded-lg text-sm font-mono text-center focus:outline-none focus:border-brand-600" />
                </div>
              </div>
              <div className="flex gap-2 pt-2">
                <button type="button" onClick={() => setShowModal(false)} className="flex-1 px-4 py-2 text-slate-600 hover:bg-slate-100 rounded-lg text-sm font-bold">Batal</button>
                <button type="submit" className="flex-1 px-4 py-2 bg-brand-600 hover:bg-brand-700 text-white rounded-lg text-sm font-bold flex items-center justify-center gap-1.5">
                  <i className="ph-bold ph-check"></i> Buat Kurikulum
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}
