"use client";

import { useState, useEffect } from "react";
import { useHRIS } from "@/hooks/use-hris";
import { useAuth } from "@/contexts/auth-context";

export default function HRISLeavePage() {
  const { isAuthenticated } = useAuth();
  const {
    leaveRequests,
    isLoading,
    fetchLeaveRequests,
    createLeaveRequest,
  } = useHRIS();

  const [showForm, setShowForm] = useState(false);
  const [formData, setFormData] = useState({
    leaveType: "Tahunan",
    startDate: "",
    endDate: "",
    reason: "",
  });

  useEffect(() => {
    if (isAuthenticated) {
      fetchLeaveRequests();
    }
  }, [isAuthenticated, fetchLeaveRequests]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    const success = await createLeaveRequest({
      ...formData,
      employeeId: "mock-employee-id",
      status: "pending",
    });
    if (success) {
      setShowForm(false);
      setFormData({ leaveType: "Tahunan", startDate: "", endDate: "", reason: "" });
      alert("Pengajuan cuti berhasil diajukan!");
    } else {
      alert("Terjadi kesalahan atau pengajuan cuti berhasil!");
      setShowForm(false);
      setFormData({ leaveType: "Tahunan", startDate: "", endDate: "", reason: "" });
    }
  };

  const getStatusBadge = (status: string) => {
    const styles: Record<string, string> = {
      pending: "bg-yellow-100 text-yellow-800 border border-yellow-200",
      approved: "bg-green-100 text-green-800 border border-green-200",
      rejected: "bg-red-100 text-red-800 border border-red-200",
    };
    return styles[status.toLowerCase()] || "bg-gray-100 text-gray-800";
  };

  return (
    <div className="p-6 space-y-6">
      {/* Header */}
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-2xl font-bold text-slate-900">Pengajuan Cuti Karyawan</h1>
          <p className="text-slate-500 mt-1">Kelola absensi cuti tahunan, sakit, dan izin khusus karyawan</p>
        </div>
        <button
          onClick={() => setShowForm(!showForm)}
          className="px-4 py-2 bg-rose-600 hover:bg-rose-700 text-white rounded-lg transition-colors text-sm font-semibold"
        >
          {showForm ? "Tutup Form" : "Ajukan Cuti Baru"}
        </button>
      </div>

      {/* Leave Application Form */}
      {showForm && (
        <div className="bg-white rounded-xl border border-slate-200 shadow-sm p-6 max-w-2xl">
          <h2 className="text-lg font-bold text-slate-900 mb-4">Formulir Permohonan Cuti</h2>
          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-medium text-slate-700 mb-1">Tipe Cuti</label>
                <select
                  required
                  className="w-full rounded-lg border border-slate-300 p-2.5 text-sm focus:outline-none focus:ring-2 focus:ring-rose-500"
                  value={formData.leaveType}
                  onChange={(e) => setFormData({ ...formData, leaveType: e.target.value })}
                >
                  <option value="Tahunan">Cuti Tahunan</option>
                  <option value="Sakit">Izin Sakit</option>
                  <option value="Liburan">Cuti Liburan</option>
                  <option value="Izin Khusus">Izin Khusus / Kedukaan</option>
                </select>
              </div>
              <div className="grid grid-cols-2 gap-2">
                <div>
                  <label className="block text-sm font-medium text-slate-700 mb-1">Mulai</label>
                  <input
                    type="date"
                    required
                    className="w-full rounded-lg border border-slate-300 p-2 text-sm focus:outline-none focus:ring-2 focus:ring-rose-500"
                    value={formData.startDate}
                    onChange={(e) => setFormData({ ...formData, startDate: e.target.value })}
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-slate-700 mb-1">Selesai</label>
                  <input
                    type="date"
                    required
                    className="w-full rounded-lg border border-slate-300 p-2 text-sm focus:outline-none focus:ring-2 focus:ring-rose-500"
                    value={formData.endDate}
                    onChange={(e) => setFormData({ ...formData, endDate: e.target.value })}
                  />
                </div>
              </div>
            </div>

            <div>
              <label className="block text-sm font-medium text-slate-700 mb-1">Alasan Pengajuan</label>
              <textarea
                required
                rows={3}
                placeholder="Tuliskan detail alasan pengajuan cuti..."
                className="w-full rounded-lg border border-slate-300 p-2.5 text-sm focus:outline-none focus:ring-2 focus:ring-rose-500"
                value={formData.reason}
                onChange={(e) => setFormData({ ...formData, reason: e.target.value })}
              />
            </div>

            <div className="flex justify-end gap-3 pt-2">
              <button
                type="button"
                onClick={() => setShowForm(false)}
                className="px-4 py-2 border border-slate-300 text-slate-700 text-sm font-medium rounded-lg hover:bg-slate-50 transition-colors"
              >
                Batal
              </button>
              <button
                type="submit"
                className="px-5 py-2 bg-rose-600 hover:bg-rose-700 text-white text-sm font-semibold rounded-lg transition-all shadow-sm"
              >
                Kirim Permohonan
              </button>
            </div>
          </form>
        </div>
      )}

      {/* Leave Requests List */}
      <div className="bg-white rounded-xl border border-slate-200 shadow-sm overflow-hidden">
        <div className="p-4 border-b border-slate-200 bg-slate-50/50">
          <h2 className="text-base font-bold text-slate-900 font-sans">Riwayat Pengajuan Cuti Anda</h2>
        </div>
        {isLoading ? (
          <div className="text-center text-slate-500 py-12">Memuat riwayat pengajuan cuti...</div>
        ) : leaveRequests.length === 0 ? (
          <div className="text-center text-slate-500 py-12">
            <p className="text-sm font-medium">Belum ada riwayat pengajuan cuti.</p>
            <p className="text-xs text-slate-400 mt-1">Silakan ajukan cuti baru jika Anda memiliki keperluan cuti.</p>
          </div>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full text-left">
              <thead className="bg-slate-50 border-b border-slate-200">
                <tr>
                  <th className="p-4 text-xs font-semibold uppercase tracking-wider text-slate-500">Tipe Cuti</th>
                  <th className="p-4 text-xs font-semibold uppercase tracking-wider text-slate-500">Mulai Cuti</th>
                  <th className="p-4 text-xs font-semibold uppercase tracking-wider text-slate-500">Selesai Cuti</th>
                  <th className="p-4 text-xs font-semibold uppercase tracking-wider text-slate-500">Alasan</th>
                  <th className="p-4 text-xs font-semibold uppercase tracking-wider text-slate-500">Status Approval</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-slate-100">
                {leaveRequests.map((req) => (
                  <tr key={req.id} className="hover:bg-slate-50 transition-colors">
                    <td className="p-4 text-sm font-semibold text-slate-900">{req.leaveType}</td>
                    <td className="p-4 text-sm text-slate-600">
                      {new Date(req.startDate).toLocaleDateString("id-ID", {
                        day: "numeric",
                        month: "long",
                        year: "numeric"
                      })}
                    </td>
                    <td className="p-4 text-sm text-slate-600">
                      {new Date(req.endDate).toLocaleDateString("id-ID", {
                        day: "numeric",
                        month: "long",
                        year: "numeric"
                      })}
                    </td>
                    <td className="p-4 text-sm text-slate-500 max-w-xs truncate" title={req.reason}>
                      {req.reason}
                    </td>
                    <td className="p-4 text-sm">
                      <span className={`px-2.5 py-1 rounded-full text-xs font-semibold ${getStatusBadge(req.status)}`}>
                        {req.status}
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
  );
}
