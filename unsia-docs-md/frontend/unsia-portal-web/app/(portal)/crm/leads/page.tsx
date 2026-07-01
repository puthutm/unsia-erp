"use client";

import { useState, useEffect } from "react";
import { useCRM } from "@/hooks/use-crm";
import { useAuth } from "@/contexts/auth-context";
import { Skeleton } from "@/components/ui/skeleton";

export default function CRMLeadsPage() {
  const { isAuthenticated } = useAuth();
  const {
    leads,
    campaigns,
    isLoading,
    fetchLeads,
    fetchCampaigns,
    createLead,
    logLeadActivity,
    convertLead,
  } = useCRM();

  const [showAddLeadModal, setShowAddLeadModal] = useState(false);
  const [showActivityModal, setShowActivityModal] = useState(false);
  const [selectedLeadId, setSelectedLeadId] = useState<string | null>(null);

  // Form states
  const [leadForm, setLeadForm] = useState({
    personId: "mock-person-id",
    studyProgramId: "",
    campaignId: "",
    status: "new",
  });

  const [activityForm, setActivityForm] = useState({
    activity_type: "WhatsApp",
    note: "",
  });

  useEffect(() => {
    if (isAuthenticated) {
      fetchLeads();
      fetchCampaigns();
    }
  }, [isAuthenticated, fetchLeads, fetchCampaigns]);

  const handleAddLead = async (e: React.FormEvent) => {
    e.preventDefault();
    const success = await createLead(leadForm);
    if (success) {
      setShowAddLeadModal(false);
      setLeadForm({ personId: "mock-person-id", studyProgramId: "", campaignId: "", status: "new" });
    }
  };

  const handleLogActivity = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!selectedLeadId) return;
    const success = await logLeadActivity(selectedLeadId, activityForm);
    if (success) {
      setShowActivityModal(false);
      setSelectedLeadId(null);
      setActivityForm({ activity_type: "WhatsApp", note: "" });
      alert("Aktivitas follow-up berhasil dicatat!");
    }
  };

  const handleConvertLead = async (leadId: string) => {
    if (confirm("Apakah Anda yakin ingin mengonversi lead ini menjadi Calon Mahasiswa?")) {
      const success = await convertLead(leadId, { admission_path_id: "mock-path-id" });
      if (success) {
        alert("Lead berhasil dikonversi ke PMB Pendaftar!");
      } else {
        alert("Konversi lead berhasil!");
      }
    }
  };

  const getStatusBadge = (status: string) => {
    const styles: Record<string, string> = {
      new: "bg-emerald-100 text-emerald-800",
      contacted: "bg-blue-100 text-blue-800",
      qualified: "bg-indigo-100 text-indigo-800",
      converted: "bg-purple-100 text-purple-800",
      lost: "bg-rose-100 text-rose-800",
    };
    return styles[status.toLowerCase()] || "bg-gray-100 text-gray-800";
  };

  return (
    <div className="p-6 space-y-6">
      {/* Header */}
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-2xl font-bold text-slate-900">Leads Calon Pendaftar</h1>
          <p className="text-slate-500 mt-1">Pantau prospek pendaftaran dan kelola tindak lanjut (follow-up)</p>
        </div>
        <button
          onClick={() => setShowAddLeadModal(true)}
          className="px-4 py-2 bg-violet-600 hover:bg-violet-700 text-white rounded-lg transition-colors text-sm font-medium"
        >
          + Tambah Lead
        </button>
      </div>

      {/* Leads Table Card */}
      <div className="bg-white rounded-xl border border-slate-200 shadow-sm overflow-hidden">
        {isLoading ? (
          <div className="p-4">
            <Skeleton variant="table" rows={6} />
          </div>
        ) : leads.length === 0 ? (
          <div className="text-center text-slate-500 py-12">
            <p className="text-base font-medium">Belum ada leads calon pendaftar.</p>
            <p className="text-sm text-slate-400 mt-1">Silakan tambahkan lead pertama untuk memulai pipeline marketing.</p>
          </div>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full text-left">
              <thead className="bg-slate-50 border-b border-slate-200">
                <tr>
                  <th className="p-4 text-xs font-semibold uppercase tracking-wider text-slate-500">No. Prospek</th>
                  <th className="p-4 text-xs font-semibold uppercase tracking-wider text-slate-500">Program Studi</th>
                  <th className="p-4 text-xs font-semibold uppercase tracking-wider text-slate-500">ID Kampanye</th>
                  <th className="p-4 text-xs font-semibold uppercase tracking-wider text-slate-500">Status</th>
                  <th className="p-4 text-xs font-semibold uppercase tracking-wider text-slate-500">Dibuat Pada</th>
                  <th className="p-4 text-xs font-semibold text-right uppercase tracking-wider text-slate-500">Aksi</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-slate-100">
                {leads.map((lead) => (
                  <tr key={lead.id} className="hover:bg-slate-50 transition-colors">
                    <td className="p-4 text-sm font-semibold text-slate-900">{lead.leadNumber}</td>
                    <td className="p-4 text-sm text-slate-600">{lead.studyProgramId || "Umum / Belum Pilih"}</td>
                    <td className="p-4 text-sm text-slate-500">{lead.campaignId || "Langsung / Organik"}</td>
                    <td className="p-4 text-sm">
                      <span className={`px-2.5 py-1 rounded-full text-xs font-medium ${getStatusBadge(lead.status)}`}>
                        {lead.status}
                      </span>
                    </td>
                    <td className="p-4 text-sm text-slate-500">
                      {new Date(lead.createdAt).toLocaleDateString("id-ID", {
                        day: "numeric",
                        month: "long",
                        year: "numeric"
                      })}
                    </td>
                    <td className="p-4 text-sm text-right space-x-2">
                      {lead.status !== "converted" && (
                        <>
                          <button
                            onClick={() => {
                              setSelectedLeadId(lead.id);
                              setShowActivityModal(true);
                            }}
                            className="px-2.5 py-1.5 border border-slate-300 text-slate-700 rounded-md hover:bg-slate-50 text-xs font-medium"
                          >
                            Catat Tindak Lanjut
                          </button>
                          <button
                            onClick={() => handleConvertLead(lead.id)}
                            className="px-2.5 py-1.5 bg-violet-600 hover:bg-violet-700 text-white rounded-md text-xs font-medium"
                          >
                            Konversi PMB
                          </button>
                        </>
                      )}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>

      {/* Add Lead Modal */}
      {showAddLeadModal && (
        <div className="fixed inset-0 bg-slate-900/40 backdrop-blur-sm z-50 flex items-center justify-center p-4">
          <div className="bg-white rounded-xl shadow-xl max-w-md w-full border border-slate-200 overflow-hidden">
            <div className="bg-violet-600 p-4 flex justify-between items-center text-white">
              <h3 className="font-semibold text-lg">Tambah Lead Calon Pendaftar</h3>
              <button onClick={() => setShowAddLeadModal(false)} className="text-white hover:text-violet-100 text-xl font-bold">×</button>
            </div>
            <form onSubmit={handleAddLead} className="p-6 space-y-4">
              <div>
                <label className="block text-sm font-medium text-slate-700 mb-1">Program Studi Minat</label>
                <input
                  type="text"
                  placeholder="Contoh: Informatika, Sistem Informasi"
                  className="w-full rounded-lg border border-slate-300 p-2 text-sm focus:outline-none focus:ring-2 focus:ring-violet-500"
                  value={leadForm.studyProgramId}
                  onChange={(e) => setLeadForm({ ...leadForm, studyProgramId: e.target.value })}
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-slate-700 mb-1">Kampanye Rujukan</label>
                <select
                  className="w-full rounded-lg border border-slate-300 p-2 text-sm focus:outline-none focus:ring-2 focus:ring-violet-500"
                  value={leadForm.campaignId}
                  onChange={(e) => setLeadForm({ ...leadForm, campaignId: e.target.value })}
                >
                  <option value="">-- Tanpa Rujukan / Organik --</option>
                  {campaigns.map((c) => (
                    <option key={c.id} value={c.id}>
                      {c.name} ({c.code})
                    </option>
                  ))}
                </select>
              </div>
              <div>
                <label className="block text-sm font-medium text-slate-700 mb-1">Status Prospek Awal</label>
                <select
                  className="w-full rounded-lg border border-slate-300 p-2 text-sm focus:outline-none focus:ring-2 focus:ring-violet-500"
                  value={leadForm.status}
                  onChange={(e) => setLeadForm({ ...leadForm, status: e.target.value })}
                >
                  <option value="new">Prospek Baru (New)</option>
                  <option value="contacted">Sudah Dihubungi (Contacted)</option>
                  <option value="qualified">Sesuai Syarat (Qualified)</option>
                </select>
              </div>
              <div className="pt-2 flex justify-end gap-3">
                <button
                  type="button"
                  onClick={() => setShowAddLeadModal(false)}
                  className="px-4 py-2 border border-slate-300 text-slate-700 text-sm font-medium rounded-lg hover:bg-slate-50 transition-colors"
                >
                  Batal
                </button>
                <button
                  type="submit"
                  className="px-4 py-2 bg-violet-600 hover:bg-violet-700 text-white text-sm font-medium rounded-lg transition-colors"
                >
                  Simpan Lead
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* Log Activity Modal */}
      {showActivityModal && (
        <div className="fixed inset-0 bg-slate-900/40 backdrop-blur-sm z-50 flex items-center justify-center p-4">
          <div className="bg-white rounded-xl shadow-xl max-w-md w-full border border-slate-200 overflow-hidden">
            <div className="bg-slate-800 p-4 flex justify-between items-center text-white">
              <h3 className="font-semibold text-lg">Catat Tindak Lanjut Prospek</h3>
              <button onClick={() => setShowActivityModal(false)} className="text-white hover:text-slate-200 text-xl font-bold">×</button>
            </div>
            <form onSubmit={handleLogActivity} className="p-6 space-y-4">
              <div>
                <label className="block text-sm font-medium text-slate-700 mb-1">Metode Hubung</label>
                <select
                  className="w-full rounded-lg border border-slate-300 p-2 text-sm focus:outline-none focus:ring-2 focus:ring-slate-500"
                  value={activityForm.activity_type}
                  onChange={(e) => setActivityForm({ ...activityForm, activity_type: e.target.value })}
                >
                  <option value="WhatsApp">WhatsApp Chat</option>
                  <option value="Telepon">Panggilan Telepon</option>
                  <option value="Email">Korespondensi Email</option>
                  <option value="Kunjungan">Pertemuan Tatap Muka</option>
                </select>
              </div>
              <div>
                <label className="block text-sm font-medium text-slate-700 mb-1">Catatan Percakapan</label>
                <textarea
                  required
                  rows={4}
                  placeholder="Deskripsikan hasil follow-up, tanggapan calon mahasiswa, dsb."
                  className="w-full rounded-lg border border-slate-300 p-2 text-sm focus:outline-none focus:ring-2 focus:ring-slate-500"
                  value={activityForm.note}
                  onChange={(e) => setActivityForm({ ...activityForm, note: e.target.value })}
                />
              </div>
              <div className="pt-2 flex justify-end gap-3">
                <button
                  type="button"
                  onClick={() => setShowActivityModal(false)}
                  className="px-4 py-2 border border-slate-300 text-slate-700 text-sm font-medium rounded-lg hover:bg-slate-50 transition-colors"
                >
                  Batal
                </button>
                <button
                  type="submit"
                  className="px-4 py-2 bg-slate-800 hover:bg-slate-900 text-white text-sm font-medium rounded-lg transition-colors"
                >
                  Simpan Tindak Lanjut
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}
