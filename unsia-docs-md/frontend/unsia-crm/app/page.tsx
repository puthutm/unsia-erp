"use client";

import { useState, useEffect } from "react";
import { useCRM } from "@/hooks/use-crm";
import { useAuth } from "@/contexts/auth-context";

export default function CRMDashboardPage() {
  const { isAuthenticated } = useAuth();
  const {
    campaigns,
    agents,
    leads,
    isLoading,
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
    <div className="space-y-6">
      {/* Header */}
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-2xl font-bold text-slate-900">CRM Dashboard</h1>
          <p className="text-slate-500 mt-1">Manage marketing campaigns, lead conversions, and Mitra referrals</p>
        </div>
        <div className="flex gap-3">
          <button
            onClick={() => setShowCampaignModal(true)}
            className="px-4 py-2 bg-violet-600 hover:bg-violet-700 text-white rounded-lg transition-colors text-sm font-medium"
          >
            + New Campaign
          </button>
          <button
            onClick={() => setShowAgentModal(true)}
            className="px-4 py-2 bg-slate-800 hover:bg-slate-900 text-white rounded-lg transition-colors text-sm font-medium"
          >
            + Register Agent
          </button>
        </div>
      </div>

      {/* CRM Overview Stats */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <div className="bg-white rounded-xl p-6 border border-slate-200 shadow-sm">
          <h3 className="text-sm font-medium text-slate-500">Total Leads</h3>
          <p className="text-2xl font-bold text-slate-900 mt-2">{leads.length}</p>
        </div>
        <div className="bg-white rounded-xl p-6 border border-slate-200 shadow-sm">
          <h3 className="text-sm font-medium text-slate-500">Active Campaigns</h3>
          <p className="text-2xl font-bold text-slate-900 mt-2">
            {campaigns.filter(c => c.status.toLowerCase() === "active").length}
          </p>
        </div>
        <div className="bg-white rounded-xl p-6 border border-slate-200 shadow-sm">
          <h3 className="text-sm font-medium text-slate-500">Mitra Agents</h3>
          <p className="text-2xl font-bold text-slate-900 mt-2">{agents.length}</p>
        </div>
        <div className="bg-white rounded-xl p-6 border border-slate-200 shadow-sm">
          <h3 className="text-sm font-medium text-slate-500">Conversion Rate</h3>
          <p className="text-2xl font-bold text-slate-900 mt-2">
            {leads.length > 0
              ? `${Math.round((leads.filter(l => l.status === "converted").length / leads.length) * 100)}%`
              : "0%"}
          </p>
        </div>
      </div>

      {/* Tabs list */}
      <div className="bg-white rounded-xl border border-slate-200 shadow-sm">
        <div className="flex border-b border-slate-200 bg-slate-50/50">
          <button
            onClick={() => setActiveTab("campaigns")}
            className={`px-6 py-3 text-sm font-medium transition-colors ${
              activeTab === "campaigns"
                ? "text-violet-600 border-b-2 border-violet-600 bg-white"
                : "text-slate-500 hover:text-slate-700"
            }`}
          >
            Campaigns
          </button>
          <button
            onClick={() => setActiveTab("agents")}
            className={`px-6 py-3 text-sm font-medium transition-colors ${
              activeTab === "agents"
                ? "text-violet-600 border-b-2 border-violet-600 bg-white"
                : "text-slate-500 hover:text-slate-700"
            }`}
          >
            Mitra Agents
          </button>
        </div>

        {/* Tab Content */}
        <div className="p-6">
          {isLoading ? (
            <div className="text-center text-slate-500 py-8">Loading data...</div>
          ) : activeTab === "campaigns" ? (
            campaigns.length === 0 ? (
              <div className="text-center text-slate-500 py-8">No marketing campaigns registered.</div>
            ) : (
              <div className="overflow-x-auto">
                <table className="w-full text-left">
                  <thead className="bg-slate-50 border-b border-slate-200">
                    <tr>
                      <th className="p-4 text-xs font-semibold uppercase tracking-wider text-slate-500">Code</th>
                      <th className="p-4 text-xs font-semibold uppercase tracking-wider text-slate-500">Name</th>
                      <th className="p-4 text-xs font-semibold uppercase tracking-wider text-slate-500">Channel</th>
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
            <div className="text-center text-slate-500 py-8">No agents registered.</div>
          ) : (
            <div className="overflow-x-auto">
              <table className="w-full text-left">
                <thead className="bg-slate-50 border-b border-slate-200">
                  <tr>
                    <th className="p-4 text-xs font-semibold uppercase tracking-wider text-slate-500">Agent Code</th>
                    <th className="p-4 text-xs font-semibold uppercase tracking-wider text-slate-500">Organization</th>
                    <th className="p-4 text-xs font-semibold uppercase tracking-wider text-slate-500">Status</th>
                    <th className="p-4 text-xs font-semibold uppercase tracking-wider text-slate-500">Approval</th>
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
              <h3 className="font-semibold text-lg">Create Campaign</h3>
              <button onClick={() => setShowCampaignModal(false)} className="text-white hover:text-violet-100 text-xl font-bold">×</button>
            </div>
            <form onSubmit={handleCreateCampaign} className="p-6 space-y-4">
              <div>
                <label className="block text-sm font-medium text-slate-700 mb-1">Code</label>
                <input
                  type="text"
                  required
                  placeholder="e.g. PMB-FB-ADS"
                  className="w-full rounded-lg border border-slate-300 p-2 text-sm focus:outline-none focus:ring-2 focus:ring-violet-500"
                  value={newCampaign.code}
                  onChange={(e) => setNewCampaign({ ...newCampaign, code: e.target.value })}
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-slate-700 mb-1">Name</label>
                <input
                  type="text"
                  required
                  placeholder="e.g. Facebook Ads Wave 1"
                  className="w-full rounded-lg border border-slate-300 p-2 text-sm focus:outline-none focus:ring-2 focus:ring-violet-500"
                  value={newCampaign.name}
                  onChange={(e) => setNewCampaign({ ...newCampaign, name: e.target.value })}
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-slate-700 mb-1">Channel</label>
                <select
                  className="w-full rounded-lg border border-slate-300 p-2 text-sm focus:outline-none focus:ring-2 focus:ring-violet-500"
                  value={newCampaign.channel}
                  onChange={(e) => setNewCampaign({ ...newCampaign, channel: e.target.value })}
                >
                  <option value="Social Media">Social Media</option>
                  <option value="Email Marketing">Email Marketing</option>
                  <option value="SEO/Website">SEO / Website</option>
                  <option value="Partnership">Partnership / Agent</option>
                  <option value="Event/Exhibition">Exhibition / Event</option>
                </select>
              </div>
              <div className="pt-2 flex justify-end gap-3">
                <button
                  type="button"
                  onClick={() => setShowCampaignModal(false)}
                  className="px-4 py-2 border border-slate-300 text-slate-700 text-sm font-medium rounded-lg hover:bg-slate-50 transition-colors"
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  className="px-4 py-2 bg-violet-600 hover:bg-violet-700 text-white text-sm font-medium rounded-lg transition-colors"
                >
                  Save Campaign
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
              <h3 className="font-semibold text-lg">Register Agent</h3>
              <button onClick={() => setShowAgentModal(false)} className="text-white hover:text-slate-200 text-xl font-bold">×</button>
            </div>
            <form onSubmit={handleRegisterAgent} className="p-6 space-y-4">
              <div>
                <label className="block text-sm font-medium text-slate-700 mb-1">Agent Code</label>
                <input
                  type="text"
                  required
                  placeholder="e.g. AGENT-SMAS-01"
                  className="w-full rounded-lg border border-slate-300 p-2 text-sm focus:outline-none focus:ring-2 focus:ring-slate-500"
                  value={newAgent.agentCode}
                  onChange={(e) => setNewAgent({ ...newAgent, agentCode: e.target.value })}
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-slate-700 mb-1">Organization Name</label>
                <input
                  type="text"
                  required
                  placeholder="e.g. SMA Negeri 1 Jakarta"
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
                  Cancel
                </button>
                <button
                  type="submit"
                  className="px-4 py-2 bg-slate-800 hover:bg-slate-900 text-white text-sm font-medium rounded-lg transition-colors"
                >
                  Register Agent
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}
