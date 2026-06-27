"use client";

import { useState, useEffect } from "react";

// KRS (Course Enrollment) Page - Next.js
// Matches: UI/AKADEMIK/MAHASISWA/SIAKAD MAHASISWA/onboarding

interface KrsItem {
  id: string;
  courseName: string;
  courseCode: string;
  className: string;
  sks: number;
  schedule: string;
  lecturerName: string;
  status: string;
}

interface Krs {
  id: string;
  academicPeriodName: string;
  totalSks: number;
  items: KrsItem[];
  status: string;
}

export default function KrsPage() {
  const [krs, setKrs] = useState<Krs | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchKrs();
  }, []);

  const fetchKrs = async () => {
    setLoading(true);
    try {
      const token = localStorage.getItem("accessToken");
      const response = await fetch("/api/v1/academic/krs", {
        headers: { Authorization: `Bearer ${token}` },
      });
      if (response.ok) {
        const data = await response.json();
        setKrs(data.data || null);
      }
    } catch (error) {
      console.error("Error fetching KRS:", error);
    } finally {
      setLoading(false);
    }
  };

  const getStatusBadge = (status: string) => {
    const styles: Record<string, string> = {
      draft: "bg-yellow-100 text-yellow-800",
      submitted: "bg-blue-100 text-blue-800",
      approved: "bg-green-100 text-green-800",
      rejected: "bg-red-100 text-red-800",
    };
    return styles[status] || "bg-gray-100 text-gray-800";
  };

  return (
    <div className="p-6 space-y-6">
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-2xl font-bold text-slate-900">KRS (Kartu Rencana Studi)</h1>
          <p className="text-slate-500 mt-1">Pengelolaan KRS mahasiswa</p>
        </div>
      </div>

      {/* KRS Info */}
      {krs && (
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div className="bg-white rounded-xl p-6 border border-slate-200">
            <h3 className="text-sm font-medium text-slate-500">Periode Akademik</h3>
            <p className="text-xl font-bold text-slate-900 mt-2">{krs.academicPeriodName}</p>
          </div>
          <div className="bg-white rounded-xl p-6 border border-slate-200">
            <h3 className="text-sm font-medium text-slate-500">Total SKS</h3>
            <p className="text-xl font-bold text-slate-900 mt-2">{krs.totalSks}</p>
          </div>
          <div className="bg-white rounded-xl p-6 border border-slate-200">
            <h3 className="text-sm font-medium text-slate-500">Status</h3>
            <span className={`inline-block px-3 py-1 rounded-full text-sm font-medium mt-2 ${getStatusBadge(krs.status)}`}>
              {krs.status.toUpperCase()}
            </span>
          </div>
        </div>
      )}

      {/* KRS Items Table */}
      <div className="bg-white rounded-xl border border-slate-200 overflow-hidden">
        <table className="w-full">
          <thead className="bg-slate-50">
            <tr>
              <th className="text-left p-4 text-sm font-medium text-slate-500">Kode</th>
              <th className="text-left p-4 text-sm font-medium text-slate-500">Mata Kuliah</th>
              <th className="text-left p-4 text-sm font-medium text-slate-500">Kelas</th>
              <th className="text-left p-4 text-sm font-medium text-slate-500">SKS</th>
              <th className="text-left p-4 text-sm font-medium text-slate-500">Jadwal</th>
              <th className="text-left p-4 text-sm font-medium text-slate-500">Dosen</th>
              <th className="text-left p-4 text-sm font-medium text-slate-500">Status</th>
            </tr>
          </thead>
          <tbody>
            {loading ? (
              <tr>
                <td colSpan={7} className="p-8 text-center text-slate-500">
                  Memuat KRS...
                </td>
              </tr>
            ) : !krs || krs.items.length === 0 ? (
              <tr>
                <td colSpan={7} className="p-8 text-center text-slate-500">
                  KRS masih kosong
                </td>
              </tr>
            ) : (
              krs.items.map((item) => (
                <tr key={item.id} className="border-t border-slate-200">
                  <td className="p-4 text-slate-900 font-mono">{item.courseCode}</td>
                  <td className="p-4 text-slate-900">{item.courseName}</td>
                  <td className="p-4 text-slate-600">{item.className}</td>
                  <td className="p-4 text-slate-600">{item.sks}</td>
                  <td className="p-4 text-slate-600">{item.schedule}</td>
                  <td className="p-4 text-slate-600">{item.lecturerName}</td>
                  <td className="p-4">
                    <span className={`px-2 py-1 rounded-full text-xs font-medium ${getStatusBadge(item.status)}`}>
                      {item.status}
                    </span>
                  </td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
}
