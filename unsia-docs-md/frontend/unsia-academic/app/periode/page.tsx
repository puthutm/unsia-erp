"use client";

import { useState, useEffect } from "react";

// Periode Akademik Page - Matches: UI/AKADEMIK/ADMIN/panel-periode
// Setting komprehensif per semester (Ganjil/Genap)

interface AcademicPeriod {
  id: string;
  academicYear: string;
  semester: "Ganjil" | "Genap";
  startDate: string;
  endDate: string;
  courseStartDate: string;
  courseEndDate: string;
  sessionsCount: number;
  utsStartDate: string;
  utsEndDate: string;
  uasStartDate: string;
  uasEndDate: string;
  quietWeek: string;
  graduationDate: string;
  totalStudents: number;
  status: "Aktif" | "Selesai" | "Arsip";
  nationalHolidays: { date: string; name: string }[];
  activities: { name: string; date: string }[];
}

export default function AcademicPeriodPage() {
  const [periods, setPeriods] = useState<AcademicPeriod[]>([]);
  const [loading, setLoading] = useState(true);
  const [showModal, setShowModal] = useState(false);

  useEffect(() => {
    fetchPeriods();
  }, []);

  const fetchPeriods = async () => {
    setLoading(true);
    try {
      const token = localStorage.getItem("accessToken");
      const response = await fetch("/api/v1/academic/academic-periods", {
        headers: { Authorization: `Bearer ${token}` },
      });
      if (response.ok) {
        const data = await response.json();
        setPeriods(data.data || getDefaultPeriods());
      } else {
        setPeriods(getDefaultPeriods());
      }
    } catch (error) {
      console.error("Error fetching periods:", error);
      setPeriods(getDefaultPeriods());
    } finally {
      setLoading(false);
    }
  };

  const getDefaultPeriods = (): AcademicPeriod[] => [
    {
      id: "P-2026-GANJIL",
      academicYear: "2026/2027",
      semester: "Ganjil",
      startDate: "01 Sep 2026",
      endDate: "15 Feb 2027",
      courseStartDate: "01 Sep 2026",
      courseEndDate: "20 Des 2026",
      sessionsCount: 16,
      utsStartDate: "27 Okt 2026",
      utsEndDate: "07 Nov 2026",
      uasStartDate: "12 Jan 2027",
      uasEndDate: "23 Jan 2027",
      quietWeek: "08-11 Jan 2027",
      graduationDate: "10 Feb 2027",
      totalStudents: 3719,
      status: "Aktif",
      nationalHolidays: [
        { date: "17 Aug 2026", name: "Hari Kemerdekaan RI" },
        { date: "10 Sep 2026", name: "Maulid Nabi Muhammad SAW" },
        { date: "25 Des 2026", name: "Hari Raya Natal" },
        { date: "01 Jan 2027", name: "Tahun Baru Masehi" },
      ],
      activities: [
        { name: "Pendaftaran Maba Gelombang 1", date: "01-30 Jun 2026" },
        { name: "Daftar Ulang Mahasiswa", date: "15-29 Aug 2026" },
        { name: "Yudisium Semester Ganjil", date: "10 Feb 2027" },
      ],
    },
    {
      id: "P-2025-GENAP",
      academicYear: "2025/2026",
      semester: "Genap",
      startDate: "16 Feb 2026",
      endDate: "31 Aug 2026",
      courseStartDate: "16 Feb 2026",
      courseEndDate: "05 Jun 2026",
      sessionsCount: 16,
      utsStartDate: "13 Apr 2026",
      utsEndDate: "24 Apr 2026",
      uasStartDate: "01 Jun 2026",
      uasEndDate: "12 Jun 2026",
      quietWeek: "29 Mei - 31 Mei 2026",
      graduationDate: "25 Aug 2026",
      totalStudents: 3580,
      status: "Selesai",
      nationalHolidays: [
        { date: "11 Mar 2026", name: "Hari Suci Nyepi" },
        { date: "29 Mar 2026", name: "Wafat Isa Al Masih" },
      ],
      activities: [],
    },
  ];

  const getStatusBadge = (status: string) => {
    const styles: Record<string, string> = {
      Aktif: "bg-emerald-100 text-emerald-800",
      Selesai: "bg-blue-100 text-blue-800",
      Arsip: "bg-slate-100 text-slate-800",
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
              <i className="ph-fill ph-calendar-check text-brand-accent"></i>
              <span className="text-[10px] uppercase tracking-widest font-bold text-brand-accent">Master · Periode Akademik</span>
            </div>
            <h2 className="font-display font-black text-2xl">Periode Akademik (Semester)</h2>
            <p className="text-brand-50 text-sm mt-1.5">
              Setiap Tahun Ajaran terdiri dari periode Ganjil + Genap. Setting periode ini otomatis mempengaruhi pembuatan kelas, jadwal sesi, dan kalender akademik mahasiswa.
            </p>
          </div>
          <button
            onClick={() => setShowModal(true)}
            className="px-3 py-2 bg-brand-accent text-brand-900 hover:bg-yellow-400 rounded-lg text-xs font-bold flex items-center gap-1.5 shrink-0"
          >
            <i className="ph-bold ph-plus-circle"></i> Tambah Periode
          </button>
        </div>
      </div>

      {/* Periods List */}
      <div className="space-y-4">
        {loading ? (
          <div className="text-center text-slate-500 py-8">Memuat...</div>
        ) : (
          periods.map((period) => (
            <div
              key={period.id}
              className={`bg-white border rounded-2xl shadow-soft overflow-hidden ${
                period.status === "Aktif" ? "border-brand-300 ring-2 ring-brand-100" : "border-slate-200"
              }`}
            >
              {/* Card Header */}
              <div className="px-5 py-4 bg-slate-50 border-b border-slate-200 flex items-center justify-between">
                <div className="flex items-center gap-3">
                  <div className={`w-11 h-11 rounded-xl flex items-center justify-center ${
                    period.status === "Aktif" ? "bg-brand-600 text-white" : "bg-slate-200 text-slate-600"
                  }`}>
                    <i className="ph-fill ph-calendar-check text-xl"></i>
                  </div>
                  <div>
                    <p className="text-[9px] uppercase tracking-wider font-bold text-slate-500">{period.id}</p>
                    <h3 className="font-display font-bold text-base text-slate-800">
                      {period.academicYear} · Semester {period.semester}
                    </h3>
                    <p className="text-[11px] text-slate-500">
                      {period.startDate} — {period.endDate} · {period.totalStudents.toLocaleString("id-ID")} mahasiswa
                    </p>
                  </div>
                </div>
                <div className="flex items-center gap-2">
                  <span className={`px-3 py-1 rounded-full text-xs font-bold ${getStatusBadge(period.status)}`}>
                    {period.status}
                  </span>
                  <button className="px-3 py-1.5 bg-brand-600 hover:bg-brand-700 text-white text-xs font-bold rounded-lg flex items-center gap-1.5">
                    <i className="ph-bold ph-eye"></i> Detail Setting
                  </button>
                </div>
              </div>

              {/* Card Body - Summary Cards */}
              <div className="p-5 grid grid-cols-2 md:grid-cols-4 gap-3">
                <div className="p-3 bg-emerald-50 border border-emerald-200 rounded-lg">
                  <p className="text-[9px] uppercase tracking-wider font-bold text-emerald-700">Jumlah Sesi/Pertemuan</p>
                  <p className="font-display font-black text-xl text-emerald-800 mt-1">{period.sessionsCount}</p>
                  <p className="text-[10px] text-emerald-600 mt-0.5">pertemuan per MK</p>
                </div>
                <div className="p-3 bg-amber-50 border border-amber-200 rounded-lg">
                  <p className="text-[9px] uppercase tracking-wider font-bold text-amber-700">Jadwal UTS</p>
                  <p className="text-xs font-bold text-amber-800 mt-1">{period.utsStartDate}</p>
                  <p className="text-[10px] text-amber-600">s/d {period.utsEndDate}</p>
                </div>
                <div className="p-3 bg-rose-50 border border-rose-200 rounded-lg">
                  <p className="text-[9px] uppercase tracking-wider font-bold text-rose-700">Jadwal UAS</p>
                  <p className="text-xs font-bold text-rose-800 mt-1">{period.uasStartDate}</p>
                  <p className="text-[10px] text-rose-600">s/d {period.uasEndDate}</p>
                </div>
                <div className="p-3 bg-violet-50 border border-violet-200 rounded-lg">
                  <p className="text-[9px] uppercase tracking-wider font-bold text-violet-700">Libur Nasional</p>
                  <p className="font-display font-black text-xl text-violet-800 mt-1">{period.nationalHolidays.length}</p>
                  <p className="text-[10px] text-violet-600">hari terdaftar</p>
                </div>
              </div>

              {/* National Holidays */}
              {period.nationalHolidays.length > 0 && (
                <div className="px-5 pb-5">
                  <p className="text-[10px] uppercase tracking-widest font-bold text-slate-400 mb-2">Libur Nasional ({period.nationalHolidays.length})</p>
                  <div className="flex flex-wrap gap-2">
                    {period.nationalHolidays.map((holiday, idx) => (
                      <div key={idx} className="flex items-center gap-2 px-3 py-1.5 bg-rose-50 border border-rose-200 rounded-lg">
                        <span className="text-[10px] font-bold text-rose-900">{holiday.name}</span>
                        <span className="text-[10px] text-rose-700 font-mono">{holiday.date}</span>
                      </div>
                    ))}
                  </div>
                </div>
              )}

              {/* Active Period Banner */}
              {period.status === "Aktif" && (
                <div className="px-5 py-3 bg-emerald-50 border-t border-emerald-200 flex items-center gap-2">
                  <i className="ph-fill ph-check-circle text-emerald-700"></i>
                  <p className="text-[11px] text-emerald-900">
                    <strong>Periode aktif</strong> — pembuatan kelas baru otomatis mengikuti periode ini. 
                    Tanggal kuliah dimulai: <strong>{period.courseStartDate}</strong>.
                  </p>
                </div>
              )}
            </div>
          ))
        )}
      </div>

      {/* Add Modal */}
      {showModal && (
        <div className="fixed inset-0 bg-slate-900/70 z-50 flex items-center justify-center p-4">
          <div className="bg-white rounded-2xl shadow-2xl max-w-xl w-full p-6">
            <h3 className="font-bold text-lg text-slate-800 mb-4">Buat Periode Akademik Baru</h3>
            <form className="space-y-4">
              <div className="grid grid-cols-2 gap-2">
                <div>
                  <label className="block text-xs font-bold text-slate-700 mb-1.5">Tahun Ajaran</label>
                  <select className="w-full px-3 py-2 bg-slate-50 border border-slate-200 rounded-lg text-sm focus:outline-none focus:border-brand-600">
                    <option>2026/2027</option>
                    <option>2025/2026</option>
                  </select>
                </div>
                <div>
                  <label className="block text-xs font-bold text-slate-700 mb-1.5">Semester</label>
                  <select className="w-full px-3 py-2 bg-slate-50 border border-slate-200 rounded-lg text-sm focus:outline-none focus:border-brand-600">
                    <option>Ganjil</option>
                    <option>Genap</option>
                  </select>
                </div>
              </div>
              <div className="bg-emerald-50 border border-emerald-200 rounded-lg p-3">
                <p className="text-xs font-bold text-emerald-900 mb-2">Setting Komprehensif Periode</p>
                <div className="grid grid-cols-2 gap-2 text-xs">
                  <div>
                    <label className="block text-[10px] font-bold text-emerald-700 mb-1">Tgl Mulai Kuliah</label>
                    <input type="date" className="w-full px-2 py-1.5 bg-white border border-emerald-200 rounded text-xs focus:outline-none" />
                  </div>
                  <div>
                    <label className="block text-[10px] font-bold text-emerald-700 mb-1">Tgl Akhir Kuliah</label>
                    <input type="date" className="w-full px-2 py-1.5 bg-white border border-emerald-200 rounded text-xs focus:outline-none" />
                  </div>
                  <div>
                    <label className="block text-[10px] font-bold text-emerald-700 mb-1">Jumlah Sesi per MK</label>
                    <input type="number" value="16" min="12" max="20" className="w-full px-2 py-1.5 bg-white border border-emerald-200 rounded text-xs focus:outline-none font-mono font-bold text-center" />
                  </div>
                  <div>
                    <label className="block text-[10px] font-bold text-emerald-700 mb-1">Tgl Mulai UTS</label>
                    <input type="date" className="w-full px-2 py-1.5 bg-white border border-emerald-200 rounded text-xs focus:outline-none" />
                  </div>
                </div>
              </div>
              <div className="flex gap-2 pt-2">
                <button type="button" onClick={() => setShowModal(false)} className="flex-1 px-4 py-2 text-slate-600 hover:bg-slate-100 rounded-lg text-sm font-bold">
                  Batal
                </button>
                <button type="submit" className="flex-1 px-4 py-2 bg-brand-600 hover:bg-brand-700 text-white rounded-lg text-sm font-bold flex items-center justify-center gap-1.5">
                  <i className="ph-bold ph-check"></i> Buat Periode
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}
