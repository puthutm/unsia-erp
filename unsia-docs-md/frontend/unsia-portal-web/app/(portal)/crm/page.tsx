"use client";

import { useState, useEffect } from "react";
import { useCRM } from "@/hooks/use-crm";
import { useAuth } from "@/contexts/auth-context";
import { Skeleton } from "@/components/ui/skeleton";

export default function CRMDashboardPage() {
  const { isAuthenticated } = useAuth();
  const {
    campaigns,
    agents,
    leads,
    isLoading,
    error,
    fetchCampaigns,
    fetchAgents,
    fetchLeads,
    createCampaign,
    registerAgent
  } = useCRM();

  const [activeTab, setActiveTab] = useState<"campaigns" | "agents">("campaigns");
  const [showCampaignModal, setShowCampaignModal] = useState(false);
  const [showAgentModal, setShowAgentModal] = useState(false);

  // Form states
  const [newCampaign, setNewCampaign] = useState({ code: "", name: "", channel: "Social Media", status: "active" });
  const [newAgent, setNewAgent] = useState({ agentCode: "", organizationName: "", status: "active", approvalStatus: "pending" });

  useEffect(() => {
    if (isAuthenticated) {
      fetchCampaigns();
      fetchAgents();
      fetchLeads();
    }
  }, [isAuthenticated, fetchCampaigns, fetchAgents, fetchLeads]);

  const handleCreateCampaign = async (e: React.FormEvent) => {
    e.preventDefault();
    const success = await createCampaign(newCampaign);
    if (success) {
      setShowCampaignModal(false);
      setNewCampaign({ code: "", name: "", channel: "Social Media", status: "active" });
    }
  };

  const handleRegisterAgent = async (e: React.FormEvent) => {
    e.preventDefault();
    const success = await registerAgent({ ...newAgent, personId: "mock-person-id" });
    if (success) {
      setShowAgentModal(false);
      setNewAgent({ agentCode: "", organizationName: "", status: "active", approvalStatus: "pending" });
    }
  };

  const getStatusBadge = (status: string) => {
    const styles: Record<string, string> = {
      active: "bg-green-100 text-green-800",
      inactive: "bg-gray-100 text-gray-800",
      pending: "bg-yellow-100 text-yellow-800",
      approved: "bg-blue-100 text-blue-800",
      rejected: "bg-red-100 text-red-800",
    };
    return styles[status.toLowerCase()] || "bg-gray-100 text-gray-800";
  };

  return (
    <div className="p-6 space-y-6">
      {/* Header */}
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-2xl font-bold text-slate-900">Customer Relationship Management (CRM)</h1>
          <p className="text-slate-500 mt-1">Kelola relasi pendaftar, referral, dan kampanye pemasaran</p>
        </div>
        <div className="flex gap-3">
          <button
            onClick={() => setShowCampaignModal(true)}
            className="px-4 py-2 bg-violet-600 hover:bg-violet-700 text-white rounded-lg transition-colors text-sm font-medium"
          >
            + Buat Kampanye
          </button>
          <button
            onClick={() => setShowAgentModal(true)}
            className="px-4 py-2 bg-slate-800 hover:bg-slate-900 text-white rounded-lg transition-colors text-sm font-medium"
          >
            + Register Agen
          </button>
        </div>
      </div>

      {/* CRM Overview Stats */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <div className="bg-white rounded-xl p-6 border border-slate-200 shadow-sm">
          <h3 className="text-sm font-medium text-slate-500">Total Leads</h3>
          {isLoading ? (
            <Skeleton className="h-8 w-20 mt-2" />
          ) : (
            <p className="text-2xl font-bold text-slate-900 mt-2">{leads.length}</p>
          )}
          <div className="text-xs text-slate-400 mt-1">Calon mahasiswa terdaftar</div>
        </div>
        <div className="bg-white rounded-xl p-6 border border-slate-200 shadow-sm">
          <h3 className="text-sm font-medium text-slate-500">Kampanye Aktif</h3>
          {isLoading ? (
            <Skeleton className="h-8 w-20 mt-2" />
          ) : (
            <p className="text-2xl font-bold text-slate-900 mt-2">
              {campaigns.filter(c => c.status.toLowerCase() === "active").length}
            </p>
          )}
          <div className="text-xs text-slate-400 mt-1">Promosi sedang berjalan</div>
        </div>
        <div className="bg-white rounded-xl p-6 border border-slate-200 shadow-sm">
          <h3 className="text-sm font-medium text-slate-500">Mitra Agen</h3>
          {isLoading ? (
            <Skeleton className="h-8 w-20 mt-2" />
          ) : (
            <p className="text-2xl font-bold text-slate-900 mt-2">{agents.length}</p>
          )}
          <div className="text-xs text-slate-400 mt-1">Total agen pendaftaran</div>
        </div>
        <div className="bg-white rounded-xl p-6 border border-slate-200 shadow-sm">
          <h3 className="text-sm font-medium text-slate-500">Konversi Leads</h3>
          {isLoading ? (
            <Skeleton className="h-8 w-20 mt-2" />
          ) : (
            <p className="text-2xl font-bold text-slate-900 mt-2">
              {leads.length > 0
                ? `${Math.round((leads.filter(l => l.status === "converted").length / leads.length) * 100)}%`
                : "0%"}
            </p>
          )}
          <div className="text-xs text-slate-400 mt-1">Rasio leads menjadi mahasiswa</div>
        </div>
      </div>

      {/* Tabs list */}
      <div className="bg-white rounded-xl border border-slate-200 shadow-sm">
        <div className="flex border-b border-slate-200">
          <button
            onClick={() => setActiveTab("campaigns")}
            className={`px-6 py-3 text-sm font-medium transition-colors ${
              activeTab === "campaigns"
                ? "text-violet-600 border-b-2 border-violet-600"
                : "text-slate-500 hover:text-slate-700"
            }`}
          >
            Kampanye Pemasaran
          </button>
          <button
            onClick={() => setActiveTab("agents")}
            className={`px-6 py-3 text-sm font-medium transition-colors ${
              activeTab === "agents"
                ? "text-violet-600 border-b-2 border-violet-600"
                : "text-slate-500 hover:text-slate-700"
            }`}
          >
            Mitra & Agen Pendaftaran
          </button>
        </div>

        {/* Tab Content */}
        <div className="p-6">
          {isLoading ? (
            <Skeleton variant="table" rows={5} />
          ) : activeTab === "campaigns" ? (
            campaigns.length === 0 ? (
              <div className="text-center text-slate-500 py-8">Belum ada kampanye terdaftar.</div>
            ) : (
              <div className="overflow-x-auto">
                <table className="w-full text-left">
                  <thead className="bg-slate-50 border-b border-slate-200">
                    <tr>
                      <th className="p-4 text-xs font-semibold uppercase tracking-wider text-slate-500">Kode</th>
                      <th className="p-4 text-xs font-semibold uppercase tracking-wider text-slate-500">Nama Kampanye</th>
                      <th className="p-4 text-xs font-semibold uppercase tracking-wider text-slate-500">Saluran</th>
                      <th className="p-4 text-xs font-semibold uppercase tracking-wider text-slate-500">Status</th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-slate-100">
                    {campaigns.map((camp) => (
                      <tr key={camp.id} className="hover:bg-slate-50 transition-colors">
                        <td className="p-4 text-sm font-semibold text-slate-900">{camp.code}</td>
                        <td className="p-4 text-sm text-slate-700">{camp.name}</td>
                        <td className="p-4 text-sm text-slate-600">{camp.channel}</td>
                        <td className="p-4 text-sm">
                          <span className={`px-2.5 py-1 rounded-full text-xs font-medium ${getStatusBadge(camp.status)}`}>
                            {camp.status}
                          </span>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            )
          ) : agents.length === 0 ? (
            <div className="text-center text-slate-500 py-8">Belum ada agen kemitraan terdaftar.</div>
          ) : (
            <div className="overflow-x-auto">
              <table className="w-full text-left">
                <thead className="bg-slate-50 border-b border-slate-200">
                  <tr>
                    <th className="p-4 text-xs font-semibold uppercase tracking-wider text-slate-500">Kode Agen</th>
                    <th className="p-4 text-xs font-semibold uppercase tracking-wider text-slate-500">Nama Organisasi</th>
                    <th className="p-4 text-xs font-semibold uppercase tracking-wider text-slate-500">Status</th>
                    <th className="p-4 text-xs font-semibold uppercase tracking-wider text-slate-500">Verifikasi</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-slate-100">
                  {agents.map((agent) => (
                    <tr key={agent.id} className="hover:bg-slate-50 transition-colors">
                      <td className="p-4 text-sm font-semibold text-slate-900">{agent.agentCode}</td>
                      <td className="p-4 text-sm text-slate-700">{agent.organizationName}</td>
                      <td className="p-4 text-sm">
                        <span className={`px-2.5 py-1 rounded-full text-xs font-medium ${getStatusBadge(agent.status)}`}>
                          {agent.status}
                        </span>
                      </td>
                      <td className="p-4 text-sm">
                        <span className={`px-2.5 py-1 rounded-full text-xs font-medium ${getStatusBadge(agent.approvalStatus)}`}>
                          {agent.approvalStatus}
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

      {/* Campaign Modal */}
      {showCampaignModal && (
        <div className="fixed inset-0 bg-slate-900/40 backdrop-blur-sm z-50 flex items-center justify-center p-4">
          <div className="bg-white rounded-xl shadow-xl max-w-md w-full border border-slate-200 overflow-hidden">
            <div className="bg-violet-600 p-4 flex justify-between items-center text-white">
              <h3 className="font-semibold text-lg">Buat Kampanye Baru</h3>
              <button onClick={() => setShowCampaignModal(false)} className="text-white hover:text-violet-100 text-xl font-bold">×</button>
            </div>
            <form onSubmit={handleCreateCampaign} className="p-6 space-y-4">
              <div>
                <label className="block text-sm font-medium text-slate-700 mb-1">Kode Kampanye</label>
                <input
                  type="text"
                  required
                  placeholder="Contoh: PMB-FB-ADS"
                  className="w-full rounded-lg border border-slate-300 p-2 text-sm focus:outline-none focus:ring-2 focus:ring-violet-500"
                  value={newCampaign.code}
                  onChange={(e) => setNewCampaign({ ...newCampaign, code: e.target.value })}
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-slate-700 mb-1">Nama Kampanye</label>
                <input
                  type="text"
                  required
                  placeholder="Contoh: Facebook Ads Gelombang 1"
                  className="w-full rounded-lg border border-slate-300 p-2 text-sm focus:outline-none focus:ring-2 focus:ring-violet-500"
                  value={newCampaign.name}
                  onChange={(e) => setNewCampaign({ ...newCampaign, name: e.target.value })}
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-slate-700 mb-1">Saluran Pemasaran</label>
                <select
                  className="w-full rounded-lg border border-slate-300 p-2 text-sm focus:outline-none focus:ring-2 focus:ring-violet-500"
                  value={newCampaign.channel}
                  onChange={(e) => setNewCampaign({ ...newCampaign, channel: e.target.value })}
                >
                  <option value="Social Media">Media Sosial</option>
                  <option value="Email Marketing">Email Marketing</option>
                  <option value="SEO/Website">SEO / Website</option>
                  <option value="Partnership">Kemitraan/Agen</option>
                  <option value="Event/Exhibition">Pameran/Event</option>
                </select>
              </div>
              <div className="pt-2 flex justify-end gap-3">
                <button
                  type="button"
                  onClick={() => setShowCampaignModal(false)}
                  className="px-4 py-2 border border-slate-300 text-slate-700 text-sm font-medium rounded-lg hover:bg-slate-50 transition-colors"
                >
                  Batal
                </button>
                <button
                  type="submit"
                  className="px-4 py-2 bg-violet-600 hover:bg-violet-700 text-white text-sm font-medium rounded-lg transition-colors"
                >
                  Simpan Kampanye
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* Agent Modal */}
      {showAgentModal && (
        <div className="fixed inset-0 bg-slate-900/40 backdrop-blur-sm z-50 flex items-center justify-center p-4">
          <div className="bg-white rounded-xl shadow-xl max-w-md w-full border border-slate-200 overflow-hidden">
            <div className="bg-slate-800 p-4 flex justify-between items-center text-white">
              <h3 className="font-semibold text-lg">Registrasi Agen Baru</h3>
              <button onClick={() => setShowAgentModal(false)} className="text-white hover:text-slate-200 text-xl font-bold">×</button>
            </div>
            <form onSubmit={handleRegisterAgent} className="p-6 space-y-4">
              <div>
                <label className="block text-sm font-medium text-slate-700 mb-1">Kode Agen</label>
                <input
                  type="text"
                  required
                  placeholder="Contoh: AGENT-SMAS-01"
                  className="w-full rounded-lg border border-slate-300 p-2 text-sm focus:outline-none focus:ring-2 focus:ring-slate-500"
                  value={newAgent.agentCode}
                  onChange={(e) => setNewAgent({ ...newAgent, agentCode: e.target.value })}
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-slate-700 mb-1">Nama Organisasi / Institusi</label>
                <input
                  type="text"
                  required
                  placeholder="Contoh: SMA Negeri 1 Jakarta"
                  className="w-full rounded-lg border border-slate-300 p-2 text-sm focus:outline-none focus:ring-2 focus:ring-slate-500"
                  value={newAgent.organizationName}
                  onChange={(e) => setNewAgent({ ...newAgent, organizationName: e.target.value })}
                />
              </div>
              <div className="pt-2 flex justify-end gap-3">
                <button
                  type="button"
                  onClick={() => setShowAgentModal(false)}
                  className="px-4 py-2 border border-slate-300 text-slate-700 text-sm font-medium rounded-lg hover:bg-slate-50 transition-colors"
                >
                  Batal
                </button>
                <button
                  type="submit"
                  className="px-4 py-2 bg-slate-800 hover:bg-slate-900 text-white text-sm font-medium rounded-lg transition-colors"
                >
                  Daftarkan Agen
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}
