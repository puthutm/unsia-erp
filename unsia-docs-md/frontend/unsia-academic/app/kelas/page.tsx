"use client";

import { useState, useEffect } from "react";

// Kelas Kuliah Page - Next.js
// Matches: UI/AKADEMIK/ADMIN/panel-kelas

interface Class {
  id: string;
  courseCode: string;
  courseName: string;
  classCode: string;
  lecturer: string;
  day: string;
  time: string;
  room: string;
  quota: number;
  enrolled: number;
  period: string;
  status: string;
  isParallel: boolean;
  parallelFrom?: string;
}

export default function KelasPage() {
  const [classes, setClasses] = useState<Class[]>([]);
  const [loading, setLoading] = useState(true);
  const [selectedPeriod, setSelectedPeriod] = useState("2026/2027 Ganjil");

  useEffect(() => {
    fetchClasses();
  }, []);

  const fetchClasses = async () => {
    setLoading(true);
    try {
      const token = localStorage.getItem("accessToken");
      const response = await fetch("/api/v1/academic/classes", {
        headers: { Authorization: `Bearer ${token}` },
      });
      if (response.ok) {
        const data = await response.json();
        setClasses(data.data || getDefaultClasses());
      } else {
        setClasses(getDefaultClasses());
      }
    } catch (error) {
      console.error("Error fetching classes:", error);
      setClasses(getDefaultClasses());
    } finally {
      setLoading(false);
    }
  };

  const getDefaultClasses = (): Class[] => [
    { id: "KLS-IF201-A", courseCode: "IF201", courseName: "Algoritma & Struktur Data", classCode: "A", lecturer: "Dr. Aulia Rahman, M.Kom.", day: "Senin", time: "08:00-11:00", room: "Lab Komputer 1 / Online", quota: 35, enrolled: 32, period: "2026/2027 Ganjil", status: "Aktif", isParallel: false },
    { id: "KLS-IF201-B", courseCode: "IF201", courseName: "Algoritma & Struktur Data", classCode: "B", lecturer: "Bp. Yusuf Andi, S.Kom., M.T.", day: "Selasa", time: "13:00-16:00", room: "Lab Komputer 1 / Online", quota: 35, enrolled: 30, period: "2026/2027 Ganjil", status: "Aktif", isParallel: true, parallelFrom: "KLS-IF201-A" },
    { id: "KLS-IF201-C", courseCode: "IF201", courseName: "Algoritma & Struktur Data", classCode: "C", lecturer: "Noviandri, S.Kom., MMSI.", day: "Rabu", time: "19:00-22:00", room: "Online (Zoom)", quota: 35, enrolled: 25, period: "2026/2027 Ganjil", status: "Aktif", isParallel: true, parallelFrom: "KLS-IF201-A" },
    { id: "KLS-IF203-A", courseCode: "IF203", courseName: "Pemrograman Berorientasi Objek", classCode: "A", lecturer: "Noviandri, S.Kom., MMSI.", day: "Senin", time: "13:00-17:00", room: "Lab Komputer 2", quota: 35, enrolled: 33, period: "2026/2027 Ganjil", status: "Aktif", isParallel: false },
    { id: "KLS-IF205-A", courseCode: "IF205", courseName: "Basis Data", classCode: "A", lecturer: "Dr. Bayu Setiawan, M.T.", day: "Selasa", time: "08:00-11:00", room: "Lab Basis Data", quota: 35, enrolled: 31, period: "2026/2027 Ganjil", status: "Aktif", isParallel: false },
    { id: "KLS-IF207-A", courseCode: "IF207", courseName: "Jaringan Komputer", classCode: "A", lecturer: "Prof. Dr. Hendro Wijaksono", day: "Kamis", time: "13:00-16:00", room: "Lab Jaringan", quota: 30, enrolled: 28, period: "2026/2027 Ganjil", status: "Aktif", isParallel: false },
    { id: "KLS-MK101-A", courseCode: "MK101", courseName: "Pendidikan Pancasila", classCode: "A", lecturer: "Bp. Surya Hartanto", day: "Jumat", time: "08:00-10:00", room: "R201", quota: 50, enrolled: 48, period: "2026/2027 Ganjil", status: "Aktif", isParallel: false },
    { id: "KLS-MK103-A", courseCode: "MK103", courseName: "Bahasa Inggris", classCode: "A", lecturer: "Ms. Diana Kartika", day: "Sabtu", time: "08:00-10:00", room: "Online", quota: 50, enrolled: 47, period: "2026/2027 Ganjil", status: "Aktif", isParallel: false },
  ];

  const totalClasses = classes.length;
  const activeClasses = classes.filter(c => c.status === "Aktif").length;
  const parallelClasses = classes.filter(c => c.isParallel).length;
  const onlineClasses = classes.filter(c => c.room.toLowerCase().includes("online") || c.room.toLowerCase().includes("zoom")).length;

  const getQuotaPercentage = (enrolled: number, quota: number) => {
    return Math.round((enrolled / quota) * 100);
  };

  return (
    <div className="p-6 space-y-6">
      {/* Header */}
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-2xl font-bold text-slate-900">Kelas Kuliah</h1>
          <p className="text-slate-500 mt-1">Periode {selectedPeriod}</p>
        </div>
        <div className="flex items-center gap-4">
          <select 
            value={selectedPeriod}
            onChange={(e) => setSelectedPeriod(e.target.value)}
            className="px-3 py-2 border border-slate-200 rounded-lg text-sm"
          >
            <option>2026/2027 Ganjil</option>
            <option>2025/2026 Genap</option>
            <option>2025/2026 Ganjil</option>
          </select>
          <button className="px-4 py-2 bg-blue-600 text-white rounded-lg font-medium flex items-center gap-2 hover:bg-blue-700">
            <span className="text-lg">+</span>
            <span>Buka Kelas Paralel</span>
          </button>
        </div>
      </div>

      {/* KPI Cards */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <div className="bg-white rounded-xl p-5 border border-slate-200">
          <h3 className="text-sm font-medium text-slate-500">Total Kelas</h3>
          <p className="text-3xl font-bold text-slate-900 mt-2">{totalClasses}</p>
          <p className="text-sm text-slate-500 mt-1">22 prodi</p>
        </div>
        <div className="bg-white rounded-xl p-5 border border-emerald-200 bg-emerald-50">
          <h3 className="text-sm font-medium text-emerald-700">Kelas Aktif</h3>
          <p className="text-3xl font-bold text-emerald-700 mt-2">{activeClasses}</p>
          <p className="text-sm text-emerald-600 mt-1">Quota terisi {'>'}80%</p>
        </div>
        <div className="bg-white rounded-xl p-5 border border-violet-200 bg-violet-50">
          <h3 className="text-sm font-medium text-violet-700">Kelas Paralel</h3>
          <p className="text-3xl font-bold text-violet-700 mt-2">{parallelClasses}</p>
          <p className="text-sm text-violet-600 mt-1">3 MK dengan ≥2 kelas</p>
        </div>
        <div className="bg-white rounded-xl p-5 border border-amber-200 bg-amber-50">
          <h3 className="text-sm font-medium text-amber-700">Kelas Online</h3>
          <p className="text-3xl font-bold text-amber-700 mt-2">{onlineClasses}</p>
          <p className="text-sm text-amber-600 mt-1">Zoom + LMS</p>
        </div>
      </div>

      {/* Table */}
      <div className="bg-white rounded-xl border border-slate-200 overflow-hidden">
        <div className="px-5 py-3 border-b border-slate-200 flex items-center justify-between">
          <h3 className="font-semibold text-slate-900">Daftar Kelas Berjalan</h3>
          <button className="px-3 py-1.5 bg-blue-600 hover:bg-blue-700 text-white text-sm font-medium rounded-lg flex items-center gap-2">
            <span>+</span> Buka Kelas Paralel
          </button>
        </div>
        <div className="overflow-x-auto">
          <table className="w-full text-sm">
            <thead className="bg-slate-50 border-b border-slate-200">
              <tr>
                <th className="text-left px-4 py-3 text-xs font-semibold text-slate-500 uppercase">ID Kelas</th>
                <th className="text-left px-4 py-3 text-xs font-semibold text-slate-500 uppercase">Mata Kuliah</th>
                <th className="text-left px-4 py-3 text-xs font-semibold text-slate-500 uppercase">Dosen Pengajar</th>
                <th className="text-left px-4 py-3 text-xs font-semibold text-slate-500 uppercase">Jadwal</th>
                <th className="text-left px-4 py-3 text-xs font-semibold text-slate-500 uppercase">Ruang</th>
                <th className="text-center px-4 py-3 text-xs font-semibold text-slate-500 uppercase">Kuota</th>
                <th className="text-center px-4 py-3 text-xs font-semibold text-slate-500 uppercase">Status</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-100">
              {loading ? (
                <tr>
                  <td colSpan={7} className="p-8 text-center text-slate-500">
                    Memuat...
                  </td>
                </tr>
              ) : classes.length === 0 ? (
                <tr>
                  <td colSpan={7} className="p-8 text-center text-slate-500">
                    Tidak ada kelas
                  </td>
                </tr>
              ) : (
                classes.map((cls) => (
                  <tr key={cls.id} className="hover:bg-slate-50">
                    <td className="px-4 py-3">
                      <div className="font-mono font-bold text-blue-600">{cls.courseCode}-{cls.classCode}</div>
                      {cls.isParallel && (
                        <div className="text-xs text-violet-600">Paralel dari {cls.parallelFrom}</div>
                      )}
                    </td>
                    <td className="px-4 py-3">
                      <div className="font-medium text-slate-900">{cls.courseName}</div>
                    </td>
                    <td className="px-4 py-3 text-slate-600">{cls.lecturer}</td>
                    <td className="px-4 py-3 text-slate-600">
                      <div>{cls.day}</div>
                      <div className="text-xs text-slate-500">{cls.time}</div>
                    </td>
                    <td className="px-4 py-3 text-slate-600">{cls.room}</td>
                    <td className="px-4 py-3 text-center">
                      <div className="flex items-center justify-center gap-2">
                        <span className="font-bold text-slate-700">{cls.enrolled}/{cls.quota}</span>
                        <div className="w-16 h-1.5 bg-slate-100 rounded-full overflow-hidden">
                          <div 
                            className="h-full bg-blue-500 rounded-full" 
                            style={{ width: `${getQuotaPercentage(cls.enrolled, cls.quota)}%` }}
                          />
                        </div>
                      </div>
                    </td>
                    <td className="px-4 py-3 text-center">
                      <span className="px-2 py-1 bg-emerald-100 text-emerald-800 rounded-full text-xs font-medium">
                        {cls.status}
                      </span>
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
