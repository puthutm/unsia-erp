"use client";

import { useState, useEffect } from "react";
import { useReference } from "@/contexts/reference-context";
import { useAuth } from "@/contexts/auth-context";
import { API_BASE_URLS, PMB_ENDPOINTS, STORAGE_KEYS } from "@/lib/constants";
import { Skeleton } from "@/components/ui/skeleton";

interface PmbStats {
  totalApplicants: number;
  activeWave: number;
  pendingPayment: number;
  admitted: number;
}

interface Applicant {
  id: string;
  applicantNumber: string;
  name: string;
  email: string;
  phone: string;
  studyProgramName: string;
  waveName: string;
  status: string;
  paymentStatus: string;
  createdAt: string;
}

export default function PmbPage() {
  const { user, isAuthenticated } = useAuth();
  const { pmbWaves, studyPrograms, isLoading: refLoading } = useReference();
  const [stats, setStats] = useState<PmbStats | null>(null);
  const [applicants, setApplicants] = useState<Applicant[]>([]);
  const [loading, setLoading] = useState(true);
  const [selectedWave, setSelectedWave] = useState<string>("");

  useEffect(() => {
    if (isAuthenticated && !refLoading) {
      fetchPmbData();
    }
  }, [isAuthenticated, refLoading]);

  const fetchPmbData = async () => {
    const token = localStorage.getItem(STORAGE_KEYS.accessToken);
    if (!token) return;

    setLoading(true);
    try {
      // Fetch stats
      const statsRes = await fetch(`${API_BASE_URLS.pmb}${PMB_ENDPOINTS.dashboard}`, {
        headers: {
          Authorization: `Bearer ${token}`,
          "Content-Type": "application/json",
        },
      });
      if (statsRes.ok) {
        const statsData = await statsRes.json();
        setStats(statsData.data);
      }

      // Fetch applicants
      const applicantsRes = await fetch(`${API_BASE_URLS.pmb}${PMB_ENDPOINTS.applicants}`, {
        headers: {
          Authorization: `Bearer ${token}`,
          "Content-Type": "application/json",
        },
      });
      if (applicantsRes.ok) {
        const applicantsData = await applicantsRes.json();
        setApplicants(applicantsData.data || []);
      }
    } catch (error) {
      console.error("Error fetching PMB data:", error);
    } finally {
      setLoading(false);
    }
  };

  const getStatusBadge = (status: string) => {
    const styles: Record<string, string> = {
      draft: "bg-gray-100 text-gray-800",
      submitted: "bg-blue-100 text-blue-800",
      verified: "bg-green-100 text-green-800",
      rejected: "bg-red-100 text-red-800",
      accepted: "bg-emerald-100 text-emerald-800",
    };
    return styles[status] || "bg-gray-100 text-gray-800";
  };

  return (
    <div className="p-6 space-y-6">
      {/* Header */}
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-2xl font-bold text-dark-900">PMB - Penerimaan Mahasiswa Baru</h1>
          <p className="text-dark-500 mt-1">Kelola pendaftaran mahasiswa baru</p>
        </div>
        <button className="px-4 py-2 bg-brand-600 text-white rounded-lg hover:bg-brand-700 transition-colors">
          + Gelombang Baru
        </button>
      </div>

      {/* Stats Cards */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <div className="bg-white rounded-xl p-6 border border-surface-border">
          <h3 className="text-sm font-medium text-dark-500">Total Pendaftaran</h3>
          {loading ? (
            <Skeleton className="h-9 w-20 mt-2" />
          ) : (
            <p className="text-3xl font-bold text-dark-900 mt-2">{stats?.totalApplicants || 0}</p>
          )}
        </div>
        <div className="bg-white rounded-xl p-6 border border-surface-border">
          <h3 className="text-sm font-medium text-dark-500">Gelombang Aktif</h3>
          {loading ? (
            <Skeleton className="h-9 w-20 mt-2" />
          ) : (
            <p className="text-3xl font-bold text-dark-900 mt-2">{stats?.activeWave || 0}</p>
          )}
        </div>
        <div className="bg-white rounded-xl p-6 border border-surface-border">
          <h3 className="text-sm font-medium text-dark-500">Menunggu Pembayaran</h3>
          {loading ? (
            <Skeleton className="h-9 w-20 mt-2" />
          ) : (
            <p className="text-3xl font-bold text-dark-900 mt-2">{stats?.pendingPayment || 0}</p>
          )}
        </div>
        <div className="bg-white rounded-xl p-6 border border-surface-border">
          <h3 className="text-sm font-medium text-dark-500">Diterima</h3>
          {loading ? (
            <Skeleton className="h-9 w-20 mt-2" />
          ) : (
            <p className="text-3xl font-bold text-dark-900 mt-2">{stats?.admitted || 0}</p>
          )}
        </div>
      </div>

      {/* Wave Selection */}
      <div className="bg-white rounded-xl p-6 border border-surface-border">
        <h2 className="text-lg font-semibold text-dark-900 mb-4">Gelombang PMB</h2>
        <div className="flex gap-4">
          {pmbWaves.filter(w => w.isActive).map((wave) => (
            <button
              key={wave.id}
              onClick={() => setSelectedWave(wave.id)}
              className={`px-4 py-2 rounded-lg border transition-colors ${
                selectedWave === wave.id
                  ? "bg-brand-50 border-brand-500 text-brand-700"
                  : "border-surface-border text-dark-600 hover:bg-gray-50"
              }`}
            >
              {wave.name}
            </button>
          ))}
          {pmbWaves.length === 0 && (
            <p className="text-dark-500">Belum ada gelombang PMB aktif</p>
          )}
        </div>
      </div>

      {/* Applicant Table */}
      <div className="bg-white rounded-xl border border-surface-border overflow-hidden">
        <div className="p-6 border-b border-surface-border">
          <h2 className="text-lg font-semibold text-dark-900">Daftar Pendaftar</h2>
        </div>
        {loading ? (
          <Skeleton variant="table" rows={5} />
        ) : applicants.length === 0 ? (
          <div className="p-8 text-center text-dark-500">Belum ada pendaftar</div>
        ) : (
          <table className="w-full">
            <thead className="bg-surface-subtle">
              <tr>
                <th className="text-left p-4 text-sm font-medium text-dark-500">No. Pendaftaran</th>
                <th className="text-left p-4 text-sm font-medium text-dark-500">Nama</th>
                <th className="text-left p-4 text-sm font-medium text-dark-500">Program Studi</th>
                <th className="text-left p-4 text-sm font-medium text-dark-500">Gelombang</th>
                <th className="text-left p-4 text-sm font-medium text-dark-500">Status</th>
                <th className="text-left p-4 text-sm font-medium text-dark-500">Pembayaran</th>
              </tr>
            </thead>
            <tbody>
              {applicants.slice(0, 10).map((applicant, index) => (
                <tr key={index} className="border-t border-surface-border">
                  <td className="p-4 text-dark-900">{applicant.applicantNumber}</td>
                  <td className="p-4 text-dark-900">{applicant.name}</td>
                  <td className="p-4 text-dark-600">{applicant.studyProgramName}</td>
                  <td className="p-4 text-dark-600">{applicant.waveName}</td>
                  <td className="p-4">
                    <span className={`px-2 py-1 rounded-full text-xs font-medium ${getStatusBadge(applicant.status)}`}>
                      {applicant.status}
                    </span>
                  </td>
                  <td className="p-4">
                    <span className={`px-2 py-1 rounded-full text-xs font-medium ${
                      applicant.paymentStatus === "paid" 
                        ? "bg-green-100 text-green-800" 
                        : "bg-yellow-100 text-yellow-800"
                    }`}>
                      {applicant.paymentStatus === "paid" ? "Lunas" : "Menunggu"}
                    </span>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>
    </div>
  );
}
