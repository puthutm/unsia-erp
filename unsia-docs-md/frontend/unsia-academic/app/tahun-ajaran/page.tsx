"use client";

import { useState, useEffect } from "react";

// Tahun Ajaran Page - Matches: UI/AKADEMIK/ADMIN/panel-tahun-ajaran
// Format YYYY/YYYY tanpa ganjil/genap

interface AcademicYear {
  id: string;
  code: string;
  startDate: string;
  endDate: string;
  status: "Aktif" | "Selesai" | "Arsip";
  totalPeriods: number;
  activeStudents: number;
}

export default function AcademicYearPage() {
  const [years, setYears] = useState<AcademicYear[]>([]);
  const [loading, setLoading] = useState(true);
  const [showModal, setShowModal] = useState(false);
  const [formData, setFormData] = useState({ yearStart: "", yearEnd: "" });

  useEffect(() => {
    fetchAcademicYears();
  }, []);

  const fetchAcademicYears = async () => {
    setLoading(true);
    try {
      const token = localStorage.getItem("accessToken");
      const response = await fetch("/api/v1/academic/academic-years", {
        headers: { Authorization: `Bearer ${token}` },
      });
      if (response.ok) {
        const data = await response.json();
        setYears(data.data || getDefaultYears());
      } else {
        setYears(getDefaultYears());
      }
    } catch (error) {
      console.error("Error fetching academic years:", error);
      setYears(getDefaultYears());
    } finally {
      setLoading(false);
    }
  };

  const getDefaultYears = (): AcademicYear[] => [
    { id: "TA-2026", code: "2026/2027", startDate: "01 Sep 2026", endDate: "31 Aug 2027", status: "Aktif", totalPeriods: 2, activeStudents: 3719 },
    { id: "TA-2025", code: "2025/2026", startDate: "01 Sep 2025", endDate: "31 Aug 2026", status: "Selesai", totalPeriods: 2, activeStudents: 3580 },
    { id: "TA-2024", code: "2024/2025", startDate: "01 Sep 2024", endDate: "31 Aug 2025", status: "Arsip", totalPeriods: 2, activeStudents: 3420 },
    { id: "TA-2023", code: "2023/2024", startDate: "01 Sep 2023", endDate: "31 Aug 2024", status: "Arsip", totalPeriods: 2, activeStudents: 3210 },
  ];

  const getStatusBadge = (status: string) => {
    const styles: Record<string, string> = {
      Aktif: "bg-emerald-100 text-emerald-800",
      Selesai: "bg-blue-100 text-blue-800",
      Arsip: "bg-slate-100 text-slate-800",
    };
    return styles[status] || "bg-gray-100 text-gray-800";
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    // API call to create new academic year
    const newYear: AcademicYear = {
      id: `TA-${formData.yearStart}`,
      code: `${formData.yearStart}/${formData.yearEnd}`,
      startDate: `01 Sep ${formData.yearStart}`,
      endDate: `31 Aug ${formData.yearEnd}`,
      status: "Selesai",
      totalPeriods: 2,
      activeStudents: 0,
    };
    setYears([newYear, ...years]);
    setShowModal(false);
    setFormData({ yearStart: "", yearEnd: "" });
  };

  return (
    <div className="p-6 lg:p-8 space-y-6">
      {/* Header */}
      <div className="bg-gradient-to-br from-brand-700 via-brand-600 to-brand-500 rounded-2xl p-6 text-white relative overflow-hidden">
        <div className="absolute -right-12 -top-12 w-48 h-48 bg-white/5 rounded-full"></div>
        <div className="relative flex items-start justify-between flex-wrap gap-4">
          <div className="flex-1 min-w-0">
            <div className="flex items-center gap-2 mb-2">
              <i className="ph-fill ph-calendar text-brand-accent"></i>
              <span className="text-[10px] uppercase tracking-widest font-bold text-brand-accent">Master · Tahun Ajaran</span>
            </div>
            <h2 className="font-display font-black text-2xl">Tahun Ajaran UNSIA</h2>
            <p className="text-brand-50 text-sm mt-1.5">
              Master tahun ajaran berformat <strong>YYYY/YYYY</strong> saja. Setting Ganjil/Genap ada di menu <strong>Periode Akademik</strong>.
            </p>
          </div>
          <button
            onClick={() => setShowModal(true)}
            className="px-3 py-2 bg-brand-accent text-brand-900 hover:bg-yellow-400 rounded-lg text-xs font-bold flex items-center gap-1.5 shrink-0"
          >
            <i className="ph-bold ph-plus-circle"></i> Tambah Tahun Ajaran
          </button>
        </div>
      </div>

      {/* Info Banner */}
      <div className="bg-blue-50 border border-blue-200 rounded-xl p-3 flex items-start gap-2">
        <i className="ph-fill ph-info text-blue-600 text-lg shrink-0"></i>
        <p className="text-xs text-blue-900">
          <strong>Catatan:</strong> Tahun Ajaran hanya menyimpan tahun saja (mis. <strong>2026/2027</strong>). 
          Pemecahan menjadi semester Ganjil & Genap dilakukan di <strong>Periode Akademik</strong> yang berisi tanggal kuliah, UTS, UAS,-libur nasional, dan aktivitas akademik per semester.
        </p>
      </div>

      {/* Years Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        {loading ? (
          <div className="col-span-4 text-center text-slate-500 py-8">Memuat...</div>
        ) : (
          years.map((year) => (
            <div
              key={year.id}
              className={`border rounded-2xl p-5 hover:shadow-soft transition-shadow ${
                year.status === "Aktif" ? "border-brand-300 bg-gradient-to-br from-brand-50 to-white" : "border-slate-200 bg-white"
              }`}
            >
              <div className="flex items-start justify-between mb-3">
                <div>
                  <p className="text-[9px] uppercase tracking-wider font-bold text-slate-500">{year.id}</p>
                  <h3 className="font-display font-black text-xl text-slate-800 mt-0.5">{year.code}</h3>
                </div>
                <span className={`px-2 py-1 rounded-full text-[10px] font-bold ${getStatusBadge(year.status)}`}>
                  {year.status}
                </span>
              </div>
              <div className="space-y-2 text-xs text-slate-600">
                <div className="flex items-center gap-2">
                  <i className="ph-bold ph-calendar-plus text-slate-400"></i>
                  <span>Mulai: <strong>{year.startDate}</strong></span>
                </div>
                <div className="flex items-center gap-2">
                  <i className="ph-bold ph-calendar-x text-slate-400"></i>
                  <span>Selesai: <strong>{year.endDate}</strong></span>
                </div>
                <div className="flex items-center gap-2">
                  <i className="ph-bold ph-students text-slate-400"></i>
                  <span>Mhs aktif: <strong>{year.activeStudents.toLocaleString("id-ID")}</strong></span>
                </div>
                <div className="flex items-center gap-2">
                  <i className="ph-bold ph-calendar-check text-slate-400"></i>
                  <span>Periode: <strong>{year.totalPeriods} semester</strong></span>
                </div>
              </div>
              <div className="flex gap-2 mt-4">
                <button className="flex-1 px-2.5 py-1.5 bg-brand-50 hover:bg-brand-100 text-brand-700 text-[10px] font-bold rounded flex items-center justify-center gap-1">
                  <i className="ph-bold ph-calendar-check"></i> Lihat Periode
                </button>
                <button className="px-2.5 py-1.5 bg-white border border-slate-200 hover:bg-slate-50 text-slate-600 text-[10px] font-bold rounded">
                  <i className="ph-bold ph-pencil-simple"></i>
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
            <h3 className="font-bold text-lg text-slate-800 mb-4">Tambah Tahun Ajaran Baru</h3>
            <form onSubmit={handleSubmit} className="space-y-4">
              <div>
                <label className="block text-xs font-bold text-slate-700 mb-1.5">Tahun Ajaran <span className="text-rose-500">*</span></label>
                <div className="grid grid-cols-3 items-center gap-2">
                  <input
                    type="number"
                    min="2020"
                    max="2099"
                    value={formData.yearStart}
                    onChange={(e) => setFormData({ ...formData, yearStart: e.target.value })}
                    className="px-3 py-2 bg-slate-50 border border-slate-200 rounded-lg text-sm focus:outline-none focus:border-brand-600 text-center font-mono font-bold"
                    placeholder="2027"
                    required
                  />
                  <p className="text-center text-slate-400">/</p>
                  <input
                    type="number"
                    min="2020"
                    max="2099"
                    value={formData.yearEnd}
                    onChange={(e) => setFormData({ ...formData, yearEnd: e.target.value })}
                    className="px-3 py-2 bg-slate-50 border border-slate-200 rounded-lg text-sm focus:outline-none focus:border-brand-600 text-center font-mono font-bold"
                    placeholder="2028"
                    required
                  />
                </div>
              </div>
              <div className="grid grid-cols-2 gap-2">
                <div>
                  <label className="block text-xs font-bold text-slate-700 mb-1.5">Tanggal Mulai</label>
                  <input
                    type="date"
                    className="w-full px-3 py-2 bg-slate-50 border border-slate-200 rounded-lg text-sm focus:outline-none focus:border-brand-600"
                  />
                </div>
                <div>
                  <label className="block text-xs font-bold text-slate-700 mb-1.5">Tanggal Selesai</label>
                  <input
                    type="date"
                    className="w-full px-3 py-2 bg-slate-50 border border-slate-200 rounded-lg text-sm focus:outline-none focus:border-brand-600"
                  />
                </div>
              </div>
              <div className="flex gap-2 pt-2">
                <button
                  type="button"
                  onClick={() => setShowModal(false)}
                  className="flex-1 px-4 py-2 text-slate-600 hover:bg-slate-100 rounded-lg text-sm font-bold"
                >
                  Batal
                </button>
                <button
                  type="submit"
                  className="flex-1 px-4 py-2 bg-brand-600 hover:bg-brand-700 text-white rounded-lg text-sm font-bold flex items-center justify-center gap-1.5"
                >
                  <i className="ph-bold ph-check"></i> Simpan
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}
