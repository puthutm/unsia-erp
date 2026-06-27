"use client";

import { useState, useCallback } from "react";
import { API_BASE_URLS, PMB_ENDPOINTS, STORAGE_KEYS } from "@/lib/constants";

export interface Applicant {
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

export interface PmbStats {
  totalApplicants: number;
  activeWave: number;
  pendingPayment: number;
  admitted: number;
}

export interface Wave {
  id: string;
  name: string;
  code: string;
  startDate: string;
  endDate: string;
  registrationStartAt: string;
  registrationEndAt: string;
  isActive: boolean;
  status: string;
}

export function usePmb() {
  const [applicants, setApplicants] = useState<Applicant[]>([]);
  const [waves, setWaves] = useState<Wave[]>([]);
  const [stats, setStats] = useState<PmbStats | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const getToken = () => localStorage.getItem(STORAGE_KEYS.accessToken);

  const fetchApplicants = useCallback(async (params?: { waveId?: string; status?: string; search?: string }) => {
    setIsLoading(true);
    setError(null);
    try {
      const token = getToken();
      if (!token) throw new Error("Not authenticated");

      const queryParams = new URLSearchParams();
      if (params?.waveId) queryParams.set("wave_id", params.waveId);
      if (params?.status) queryParams.set("status", params.status);
      if (params?.search) queryParams.set("search", params.search);

      const url = `${API_BASE_URLS.pmb}${PMB_ENDPOINTS.applicants}${queryParams.toString() ? `?${queryParams}` : ""}`;
      const response = await fetch(url, {
        headers: { Authorization: `Bearer ${token}`, "Content-Type": "application/json" },
      });

      if (!response.ok) throw new Error("Failed to fetch applicants");
      const data = await response.json();
      setApplicants(data.data || []);
      return data.data || [];
    } catch (err) {
      setError(err instanceof Error ? err.message : "Error");
      return [];
    } finally {
      setIsLoading(false);
    }
  }, []);

  const fetchWaves = useCallback(async () => {
    setIsLoading(true);
    setError(null);
    try {
      const token = getToken();
      if (!token) throw new Error("Not authenticated");

      const response = await fetch(`${API_BASE_URLS.pmb}${PMB_ENDPOINTS.waves}`, {
        headers: { Authorization: `Bearer ${token}`, "Content-Type": "application/json" },
      });

      if (!response.ok) throw new Error("Failed to fetch waves");
      const data = await response.json();
      setWaves(data.data || []);
      return data.data || [];
    } catch (err) {
      setError(err instanceof Error ? err.message : "Error");
      return [];
    } finally {
      setIsLoading(false);
    }
  }, []);

  const fetchStats = useCallback(async () => {
    try {
      const token = getToken();
      if (!token) return null;

      const response = await fetch(`${API_BASE_URLS.pmb}${PMB_ENDPOINTS.dashboard}`, {
        headers: { Authorization: `Bearer ${token}`, "Content-Type": "application/json" },
      });

      if (!response.ok) return null;
      const data = await response.json();
      setStats(data.data);
      return data.data;
    } catch {
      return null;
    }
  }, []);

  const updateApplicantStatus = useCallback(async (id: string, status: string) => {
    setIsLoading(true);
    try {
      const token = getToken();
      if (!token) throw new Error("Not authenticated");

      const response = await fetch(`${API_BASE_URLS.pmb}${PMB_ENDPOINTS.applicants}/${id}`, {
        method: "PUT",
        headers: { Authorization: `Bearer ${token}`, "Content-Type": "application/json" },
        body: JSON.stringify({ status }),
      });

      if (!response.ok) throw new Error("Failed to update applicant");
      await fetchApplicants();
      return true;
    } catch (err) {
      setError(err instanceof Error ? err.message : "Error");
      return false;
    } finally {
      setIsLoading(false);
    }
  }, [fetchApplicants]);

  return {
    applicants,
    waves,
    stats,
    isLoading,
    error,
    fetchApplicants,
    fetchWaves,
    fetchStats,
    updateApplicantStatus,
  };
}
