"use client";

import { useState, useEffect } from "react";
import { useHRIS } from "@/hooks/use-hris";
import { useAuth } from "@/contexts/auth-context";
import { Skeleton } from "@/components/ui/skeleton";

export default function HRISDashboardPage() {
  const { isAuthenticated } = useAuth();
  const {
    employees,
    attendances,
    isLoading,
    error,
    fetchEmployees,
    fetchAttendances,
    clockIn,
    clockOut,
  } = useHRIS();

  const [notes, setNotes] = useState("");
  const [currentTime, setCurrentTime] = useState("");
  const [isClockedIn, setIsClockedIn] = useState(false);

  // Digital clock update
  useEffect(() => {
    const timer = setInterval(() => {
      const date = new Date();
      setCurrentTime(date.toLocaleTimeString("id-ID", { hour12: false }));
    }, 1000);
    return () => clearInterval(timer);
  }, []);

  useEffect(() => {
    if (isAuthenticated) {
      fetchEmployees();
      fetchAttendances().then((records) => {
        // Simple logic to check if already clocked-in today
        const todayStr = new Date().toISOString().split("T")[0];
        const todayRecord = records.find((r) => r.date === todayStr);
        if (todayRecord) {
          setIsClockedIn(!!todayRecord.clockInTime && !todayRecord.clockOutTime);
        }
      });
    }
  }, [isAuthenticated, fetchEmployees, fetchAttendances]);

  const handleClockIn = async () => {
    const success = await clockIn(notes || "Hadir kerja");
    if (success) {
      setIsClockedIn(true);
      setNotes("");
      alert("Clock-In Presensi Berhasil!");
    }
  };

  const handleClockOut = async () => {
    const success = await clockOut(notes || "Selesai kerja");
    if (success) {
      setIsClockedIn(false);
      setNotes("");
      alert("Clock-Out Presensi Berhasil!");
    }
  };

  const getStatusBadge = (status: string) => {
    const styles: Record<string, string> = {
      aktif: "bg-green-100 text-green-800",
      active: "bg-green-100 text-green-800",
      cuti: "bg-amber-100 text-amber-800",
      nonaktif: "bg-red-100 text-red-800",
      hadir: "bg-green-100 text-green-800",
      sakit: "bg-orange-100 text-orange-800",
      alpha: "bg-rose-100 text-rose-800",
    };
    return styles[status.toLowerCase()] || "bg-gray-100 text-gray-800";
  };

  return (
    <div className="p-6 space-y-6">
      {/* Header */}
      <div>
        <h1 className="text-2xl font-bold text-slate-900">Human Resource Information System (HRIS)</h1>
        <p className="text-slate-500 mt-1">Sistem informasi manajemen kehadiran dan database karyawan</p>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Attendance Widget */}
        <div className="bg-white rounded-xl border border-slate-200 shadow-sm p-6 space-y-4 lg:col-span-1">
          <h2 className="text-lg font-bold text-slate-900 border-b border-slate-100 pb-3">Presensi Kerja Mandiri</h2>
          <div className="text-center py-6 bg-rose-50/50 rounded-xl border border-rose-100">
            <p className="text-sm font-medium text-rose-600">Waktu Saat Ini (WIB)</p>
            <p className="text-4xl font-extrabold text-rose-700 tracking-wider mt-1">{currentTime || "00:00:00"}</p>
            <p className="text-xs text-slate-500 mt-2">
              {new Date().toLocaleDateString("id-ID", { weekday: "long", day: "numeric", month: "long", year: "numeric" })}
            </p>
          </div>

          <div className="space-y-2">
            <label className="block text-sm font-medium text-slate-700">Catatan Lokasi / Aktivitas</label>
            <input
              type="text"
              placeholder="Contoh: WFH (Home), WFO (Gedung Rektorat)"
              className="w-full rounded-lg border border-slate-300 p-2 text-sm focus:outline-none focus:ring-2 focus:ring-rose-500"
              value={notes}
              onChange={(e) => setNotes(e.target.value)}
            />
          </div>

          <div className="grid grid-cols-2 gap-3 pt-2">
            <button
              onClick={handleClockIn}
              disabled={isClockedIn || isLoading}
              className={`w-full py-2.5 rounded-lg text-sm font-semibold transition-all ${
                isClockedIn
                  ? "bg-slate-100 text-slate-400 cursor-not-allowed"
                  : "bg-rose-600 hover:bg-rose-700 text-white shadow-sm"
              }`}
            >
              Clock In
            </button>
            <button
              onClick={handleClockOut}
              disabled={!isClockedIn || isLoading}
              className={`w-full py-2.5 rounded-lg text-sm font-semibold transition-all ${
                !isClockedIn
                  ? "bg-slate-100 text-slate-400 cursor-not-allowed"
                  : "bg-slate-800 hover:bg-slate-900 text-white shadow-sm"
              }`}
            >
              Clock Out
            </button>
          </div>
        </div>

        {/* Directory/History Container */}
        <div className="lg:col-span-2 space-y-6">
          {/* Employee Directory */}
          <div className="bg-white rounded-xl border border-slate-200 shadow-sm overflow-hidden">
            <div className="p-4 border-b border-slate-200 bg-slate-50/50 flex justify-between items-center">
              <h2 className="text-base font-bold text-slate-900">Daftar Karyawan UNSIA</h2>
              <span className="bg-rose-100 text-rose-800 text-xs px-2.5 py-1 rounded-full font-semibold">
                Total: {employees.length}
              </span>
            </div>
            {isLoading ? (
              <Skeleton variant="table" rows={5} />
            ) : employees.length === 0 ? (
              <div className="text-center text-slate-500 py-8">Tidak ada data karyawan terdaftar.</div>
            ) : (
              <div className="overflow-x-auto">
                <table className="w-full text-left">
                  <thead className="bg-slate-50 border-b border-slate-200">
                    <tr>
                      <th className="p-3 text-xs font-semibold text-slate-500">NIP</th>
                      <th className="p-3 text-xs font-semibold text-slate-500">Nama</th>
                      <th className="p-3 text-xs font-semibold text-slate-500">Email</th>
                      <th className="p-3 text-xs font-semibold text-slate-500">Jabatan</th>
                      <th className="p-3 text-xs font-semibold text-slate-500">Status</th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-slate-100">
                    {employees.map((emp) => (
                      <tr key={emp.id} className="hover:bg-slate-50 transition-colors">
                        <td className="p-3 text-sm font-semibold text-slate-900">{emp.employeeNumber}</td>
                        <td className="p-3 text-sm text-slate-700">{emp.name || "Karyawan UNSIA"}</td>
                        <td className="p-3 text-sm text-slate-500">{emp.email || "karyawan@unsia.ac.id"}</td>
                        <td className="p-3 text-sm text-slate-600">{emp.role}</td>
                        <td className="p-3 text-sm">
                          <span className={`px-2 py-0.5 rounded-full text-xs font-medium ${getStatusBadge(emp.status)}`}>
                            {emp.status}
                          </span>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            )}
          </div>

          {/* Attendance History */}
          <div className="bg-white rounded-xl border border-slate-200 shadow-sm overflow-hidden">
            <div className="p-4 border-b border-slate-200 bg-slate-50/50">
              <h2 className="text-base font-bold text-slate-900">Riwayat Presensi Terbaru</h2>
            </div>
            {attendances.length === 0 ? (
              <div className="text-center text-slate-500 py-8">Belum ada riwayat kehadiran tercatat hari ini.</div>
            ) : (
              <div className="overflow-x-auto">
                <table className="w-full text-left">
                  <thead className="bg-slate-50 border-b border-slate-200">
                    <tr>
                      <th className="p-3 text-xs font-semibold text-slate-500">Tanggal</th>
                      <th className="p-3 text-xs font-semibold text-slate-500">Jam Masuk</th>
                      <th className="p-3 text-xs font-semibold text-slate-500">Jam Keluar</th>
                      <th className="p-3 text-xs font-semibold text-slate-500">Keterangan</th>
                      <th className="p-3 text-xs font-semibold text-slate-500">Status</th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-slate-100">
                    {attendances.slice(0, 5).map((att) => (
                      <tr key={att.id} className="hover:bg-slate-50 transition-colors">
                        <td className="p-3 text-sm text-slate-900">{att.date}</td>
                        <td className="p-3 text-sm text-green-700 font-semibold">{att.clockInTime}</td>
                        <td className="p-3 text-sm text-slate-500">{att.clockOutTime || "--:--:--"}</td>
                        <td className="p-3 text-sm text-slate-600">{att.notes || "-"}</td>
                        <td className="p-3 text-sm">
                          <span className={`px-2.5 py-1 rounded-full text-xs font-semibold ${getStatusBadge(att.status || "Hadir")}`}>
                            {att.status || "Hadir"}
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
      </div>
    </div>
  );
}
