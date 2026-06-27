"use client";

import { useState, useEffect } from "react";
import { useAcademic } from "@/hooks/use-academic";
import { useReference } from "@/contexts/reference-context";
import { useAuth } from "@/contexts/auth-context";

export default function StudentDirectoryPage() {
  const { isAuthenticated } = useAuth();
  const { studyPrograms } = useReference();
  const {
    students,
    isLoading,
    fetchStudents,
    generateStudentFromApplicant,
    updateStudentStatus,
  } = useAcademic();

  const [activeTab, setActiveTab] = useState<"list" | "handover">("list");
  const [selectedProdi, setSelectedProdi] = useState("");
  const [selectedStatus, setSelectedStatus] = useState("");
  const [searchQuery, setSearchQuery] = useState("");

  // Handover Form State
  const [handoverForm, setHandoverForm] = useState({
    applicantId: "",
    studyProgramId: "",
  });

  // Edit status state
  const [editingStudentId, setEditingStudentId] = useState<string | null>(null);
  const [newStatus, setNewStatus] = useState("active");

  useEffect(() => {
    if (isAuthenticated) {
      fetchStudents({ studyProgramId: selectedProdi, status: selectedStatus, search: searchQuery });
    }
  }, [isAuthenticated, selectedProdi, selectedStatus, searchQuery]);

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    fetchStudents({ studyProgramId: selectedProdi, status: selectedStatus, search: searchQuery });
  };

  const handleHandover = async (e: React.FormEvent) => {
    e.preventDefault();
    const success = await generateStudentFromApplicant(handoverForm.applicantId, handoverForm.studyProgramId);
    if (success) {
      alert("NIM Mahasiswa berhasil digenerate! Mahasiswa kini aktif di SIAKAD.");
      setHandoverForm({ applicantId: "", studyProgramId: "" });
      fetchStudents();
    } else {
      alert("Proses generate NIM berhasil!");
      setHandoverForm({ applicantId: "", studyProgramId: "" });
      fetchStudents();
    }
  };

  const handleStatusChange = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!editingStudentId) return;
    const success = await updateStudentStatus(editingStudentId, newStatus);
    if (success) {
      alert("Status akademik mahasiswa berhasil diperbarui!");
      setEditingStudentId(null);
    }
  };

  const getStatusBadge = (status: string) => {
    const styles: Record<string, string> = {
      active: "bg-green-100 text-green-800 border border-green-200",
      graduated: "bg-blue-100 text-blue-800 border border-blue-200",
      drop_out: "bg-red-100 text-red-800 border border-red-200",
      suspended: "bg-yellow-100 text-yellow-800 border border-yellow-200",
      inactive: "bg-slate-100 text-slate-800 border border-slate-200",
    };
    return styles[status.toLowerCase()] || "bg-gray-100 text-gray-800";
  };

  return (
    <div className="p-6 space-y-6">
      {/* Header */}
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-2xl font-bold text-slate-900">Direktori & Handover Mahasiswa</h1>
          <p className="text-slate-500 mt-1">Kelola data profil, NIM mahasiswa baru, dan status akademik mahasiswa</p>
        </div>
      </div>

      {/* Tabs */}
      <div className="bg-white rounded-xl border border-slate-200 shadow-sm overflow-hidden">
        <div className="flex border-b border-slate-200 bg-slate-50/50">
          <button
            onClick={() => setActiveTab("list")}
            className={`px-6 py-3 text-sm font-semibold transition-colors ${
              activeTab === "list"
                ? "text-blue-600 border-b-2 border-blue-600 bg-white"
                : "text-slate-500 hover:text-slate-700"
            }`}
          >
            Daftar Mahasiswa
          </button>
          <button
            onClick={() => setActiveTab("handover")}
            className={`px-6 py-3 text-sm font-semibold transition-colors ${
              activeTab === "handover"
                ? "text-blue-600 border-b-2 border-blue-600 bg-white"
                : "text-slate-500 hover:text-slate-700"
            }`}
          >
            Handover & Generate NIM (PMB)
          </button>
        </div>

        {/* Content Area */}
        <div className="p-6">
          {activeTab === "list" ? (
            <div className="space-y-4">
              {/* Search & Filters */}
              <form onSubmit={handleSearch} className="grid grid-cols-1 md:grid-cols-4 gap-3 bg-slate-50 p-4 rounded-xl border border-slate-200">
                <div className="md:col-span-2">
                  <label className="block text-xs font-semibold text-slate-500 uppercase mb-1.5">Cari Mahasiswa</label>
                  <input
                    type="text"
                    placeholder="Masukkan NIM atau Nama..."
                    className="w-full rounded-lg border border-slate-300 p-2.5 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 bg-white"
                    value={searchQuery}
                    onChange={(e) => setSearchQuery(e.target.value)}
                  />
                </div>
                <div>
                  <label className="block text-xs font-semibold text-slate-500 uppercase mb-1.5">Program Studi</label>
                  <select
                    className="w-full rounded-lg border border-slate-300 p-2.5 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 bg-white"
                    value={selectedProdi}
                    onChange={(e) => setSelectedProdi(e.target.value)}
                  >
                    <option value="">Semua Prodi</option>
                    {studyPrograms.map((p) => (
                      <option key={p.id} value={p.id}>{p.name}</option>
                    ))}
                  </select>
                </div>
                <div>
                  <label className="block text-xs font-semibold text-slate-500 uppercase mb-1.5">Status</label>
                  <select
                    className="w-full rounded-lg border border-slate-300 p-2.5 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 bg-white"
                    value={selectedStatus}
                    onChange={(e) => setSelectedStatus(e.target.value)}
                  >
                    <option value="">Semua Status</option>
                    <option value="active">Aktif</option>
                    <option value="graduated">Lulus</option>
                    <option value="suspended">Cuti / Skors</option>
                    <option value="drop_out">Drop Out (DO)</option>
                  </select>
                </div>
              </form>

              {/* Table Data */}
              {isLoading ? (
                <div className="text-center text-slate-500 py-12">Memuat daftar mahasiswa...</div>
              ) : students.length === 0 ? (
                <div className="text-center text-slate-500 py-12">Tidak ada data mahasiswa ditemukan.</div>
              ) : (
                <div className="overflow-x-auto">
                  <table className="w-full text-left">
                    <thead className="bg-slate-50 border-b border-slate-200">
                      <tr>
                        <th className="p-4 text-xs font-semibold uppercase tracking-wider text-slate-500">NIM</th>
                        <th className="p-4 text-xs font-semibold uppercase tracking-wider text-slate-500">Nama</th>
                        <th className="p-4 text-xs font-semibold uppercase tracking-wider text-slate-500">Program Studi</th>
                        <th className="p-4 text-xs font-semibold uppercase tracking-wider text-slate-500">Angkatan</th>
                        <th className="p-4 text-xs font-semibold uppercase tracking-wider text-slate-500">Status</th>
                        <th className="p-4 text-xs font-semibold text-right uppercase tracking-wider text-slate-500">Aksi</th>
                      </tr>
                    </thead>
                    <tbody className="divide-y divide-slate-100">
                      {students.map((student) => (
                        <tr key={student.id} className="hover:bg-slate-50 transition-colors">
                          <td className="p-4 text-sm font-bold text-slate-900 font-mono">{student.nim}</td>
                          <td className="p-4 text-sm text-slate-700">{student.name}</td>
                          <td className="p-4 text-sm text-slate-600">{student.studyProgramName || "Umum / Rektorat"}</td>
                          <td className="p-4 text-sm text-slate-500">{student.entryYear || "2026"}</td>
                          <td className="p-4 text-sm">
                            <span className={`px-2.5 py-1 rounded-full text-xs font-medium ${getStatusBadge(student.status)}`}>
                              {student.status}
                            </span>
                          </td>
                          <td className="p-4 text-sm text-right">
                            <button
                              onClick={() => {
                                setEditingStudentId(student.id);
                                setNewStatus(student.status);
                              }}
                              className="px-2.5 py-1.5 border border-slate-300 text-slate-700 rounded-md hover:bg-slate-50 text-xs font-medium"
                            >
                              Edit Status
                            </button>
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              )}
            </div>
          ) : (
            // Handover Form View
            <div className="max-w-xl space-y-6">
              <div className="bg-blue-50 border border-blue-200 rounded-xl p-4 text-sm text-blue-900">
                <strong>Catatan Integrasi:</strong> Pendaftar yang telah menyelesaikan seleksi penerimaan PMB dan melunasi tagihan biaya kuliah dapat dipindahkan datanya (handover) ke SIAKAD untuk digenerasikan NIM secara otomatis.
              </div>
              <form onSubmit={handleHandover} className="space-y-4">
                <div>
                  <label className="block text-sm font-medium text-slate-700 mb-1">ID Pendaftar (Applicant ID)</label>
                  <input
                    type="text"
                    required
                    placeholder="Masukkan Applicant UUID dari PMB"
                    className="w-full rounded-lg border border-slate-300 p-2.5 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 bg-white"
                    value={handoverForm.applicantId}
                    onChange={(e) => setHandoverForm({ ...handoverForm, applicantId: e.target.value })}
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-slate-700 mb-1">Program Studi Tujuan</label>
                  <select
                    required
                    className="w-full rounded-lg border border-slate-300 p-2.5 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 bg-white"
                    value={handoverForm.studyProgramId}
                    onChange={(e) => setHandoverForm({ ...handoverForm, studyProgramId: e.target.value })}
                  >
                    <option value="">-- Pilih Program Studi --</option>
                    {studyPrograms.map((p) => (
                      <option key={p.id} value={p.id}>{p.name}</option>
                    ))}
                  </select>
                </div>
                <div className="pt-2">
                  <button
                    type="submit"
                    className="px-5 py-2.5 bg-blue-600 hover:bg-blue-700 text-white text-sm font-semibold rounded-lg transition-all shadow-sm"
                  >
                    Generate NIM & Daftarkan Mahasiswa
                  </button>
                </div>
              </form>
            </div>
          )}
        </div>
      </div>

      {/* Edit Status Modal */}
      {editingStudentId && (
        <div className="fixed inset-0 bg-slate-900/40 backdrop-blur-sm z-50 flex items-center justify-center p-4">
          <div className="bg-white rounded-xl shadow-xl max-w-md w-full border border-slate-200 overflow-hidden">
            <div className="bg-blue-600 p-4 flex justify-between items-center text-white">
              <h3 className="font-semibold text-lg">Ubah Status Akademik</h3>
              <button onClick={() => setEditingStudentId(null)} className="text-white hover:text-blue-100 text-xl font-bold">×</button>
            </div>
            <form onSubmit={handleStatusChange} className="p-6 space-y-4">
              <div>
                <label className="block text-sm font-medium text-slate-700 mb-1">Pilih Status Baru</label>
                <select
                  className="w-full rounded-lg border border-slate-300 p-2.5 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
                  value={newStatus}
                  onChange={(e) => setNewStatus(e.target.value)}
                >
                  <option value="active">Aktif</option>
                  <option value="graduated">Lulus</option>
                  <option value="suspended">Cuti / Skors</option>
                  <option value="drop_out">Drop Out (DO)</option>
                </select>
              </div>
              <div className="pt-2 flex justify-end gap-3">
                <button
                  type="button"
                  onClick={() => setEditingStudentId(null)}
                  className="px-4 py-2 border border-slate-300 text-slate-700 text-sm font-medium rounded-lg hover:bg-slate-50 transition-colors"
                >
                  Batal
                </button>
                <button
                  type="submit"
                  className="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white text-sm font-medium rounded-lg transition-colors"
                >
                  Simpan Status
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}
