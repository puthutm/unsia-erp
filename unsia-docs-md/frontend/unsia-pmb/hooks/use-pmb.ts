"use client";

import { useState, useCallback } from "react";

// API Configuration (standalone module)
const API_BASE_URLS = {
  pmb: process.env.NEXT_PUBLIC_PMB_API || "http://localhost:8003",
  reference: process.env.NEXT_PUBLIC_REFERENCE_API || "http://localhost:8002",
};

const PMB_ENDPOINTS = {
  applicants: "/api/v1/pmb/applicants",
  waves: "/api/v1/pmb/waves",
  selection: "/api/v1/pmb/selection",
  dashboard: "/api/v1/pmb/dashboard",
  documents: "/api/v1/pmb/documents",
  payment: "/api/v1/pmb/payment",
};

const STORAGE_KEYS = {
  accessToken: "unsia_access_token",
};

// Types
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

export interface ApplicantDocument {
  id: string;
  applicantId: string;
  documentTypeCode: string;
  documentTypeName: string;
  fileUrl: string;
  verificationStatus: "pending" | "verified" | "rejected";
  verifiedBy?: string;
  verifiedAt?: string;
  rejectReason?: string;
  createdAt: string;
}

export interface SelectionResult {
  id: string;
  applicantId: string;
  testScore: number;
  result: "passed" | "failed" | "pending";
  notes?: string;
  createdAt: string;
}

export function usePmb() {
  const [applicants, setApplicants] = useState<Applicant[]>([]);
  const [waves, setWaves] = useState<Wave[]>([]);
  const [stats, setStats] = useState<PmbStats | null>(null);
  const [documents, setDocuments] = useState<ApplicantDocument[]>([]);
  const [selectionResults, setSelectionResults] = useState<SelectionResult[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const getToken = () => {
    if (typeof window === "undefined") return null;
    return localStorage.getItem(STORAGE_KEYS.accessToken);
  };

  // Fetch applicants list
  const fetchApplicants = useCallback(async (params?: { waveId?: string; status?: string; search?: string; page?: number; limit?: number }) => {
    setIsLoading(true);
    setError(null);
    try {
      const token = getToken();
      if (!token) throw new Error("Not authenticated");

      const queryParams = new URLSearchParams();
      if (params?.waveId) queryParams.set("wave_id", params.waveId);
      if (params?.status) queryParams.set("status", params.status);
      if (params?.search) queryParams.set("search", params.search);
      if (params?.page) queryParams.set("page", params.page.toString());
      if (params?.limit) queryParams.set("limit", params.limit.toString());

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

  // Fetch waves list
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

  // Fetch dashboard stats
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

  // Fetch applicant documents
  const fetchDocuments = useCallback(async (applicantId: string) => {
    setIsLoading(true);
    setError(null);
    try {
      const token = getToken();
      if (!token) throw new Error("Not authenticated");

      const response = await fetch(`${API_BASE_URLS.pmb}${PMB_ENDPOINTS.documents}?applicant_id=${applicantId}`, {
        headers: { Authorization: `Bearer ${token}`, "Content-Type": "application/json" },
      });

      if (!response.ok) throw new Error("Failed to fetch documents");
      const data = await response.json();
      setDocuments(data.data || []);
      return data.data || [];
    } catch (err) {
      setError(err instanceof Error ? err.message : "Error");
      return [];
    } finally {
      setIsLoading(false);
    }
  }, []);

  // Fetch selection results
  const fetchSelectionResults = useCallback(async (params?: { waveId?: string }) => {
    setIsLoading(true);
    setError(null);
    try {
      const token = getToken();
      if (!token) throw new Error("Not authenticated");

      const queryParams = new URLSearchParams();
      if (params?.waveId) queryParams.set("wave_id", params.waveId);

      const url = `${API_BASE_URLS.pmb}${PMB_ENDPOINTS.selection}${queryParams.toString() ? `?${queryParams}` : ""}`;
      const response = await fetch(url, {
        headers: { Authorization: `Bearer ${token}`, "Content-Type": "application/json" },
      });

      if (!response.ok) throw new Error("Failed to fetch selection results");
      const data = await response.json();
      setSelectionResults(data.data || []);
      return data.data || [];
    } catch (err) {
      setError(err instanceof Error ? err.message : "Error");
      return [];
    } finally {
      setIsLoading(false);
    }
  }, []);

  // Update applicant status
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
  }, []);

  // Verify document
  const verifyDocument = useCallback(async (documentId: string, status: "verified" | "rejected", reason?: string) => {
    setIsLoading(true);
    try {
      const token = getToken();
      if (!token) throw new Error("Not authenticated");

      const response = await fetch(`${API_BASE_URLS.pmb}${PMB_ENDPOINTS.documents}/${documentId}/verify`, {
        method: "PUT",
        headers: { Authorization: `Bearer ${token}`, "Content-Type": "application/json" },
        body: JSON.stringify({ status, rejectReason: reason }),
      });

      if (!response.ok) throw new Error("Failed to verify document");
      return true;
    } catch (err) {
      setError(err instanceof Error ? err.message : "Error");
      return false;
    } finally {
      setIsLoading(false);
    }
  }, []);

  // Input selection result
  const inputSelectionResult = useCallback(async (applicantId: string, testScore: number, result: "passed" | "failed", notes?: string) => {
    setIsLoading(true);
    try {
      const token = getToken();
      if (!token) throw new Error("Not authenticated");

      const response = await fetch(`${API_BASE_URLS.pmb}${PMB_ENDPOINTS.selection}`, {
        method: "POST",
        headers: { Authorization: `Bearer ${token}`, "Content-Type": "application/json" },
        body: JSON.stringify({ applicantId, testScore, result, notes }),
      });

      if (!response.ok) throw new Error("Failed to input selection result");
      return true;
    } catch (err) {
      setError(err instanceof Error ? err.message : "Error");
      return false;
    } finally {
      setIsLoading(false);
    }
  }, []);

  // Create wave
  const createWave = useCallback(async (waveData: Partial<Wave>) => {
    setIsLoading(true);
    try {
      const token = getToken();
      if (!token) throw new Error("Not authenticated");

      const response = await fetch(`${API_BASE_URLS.pmb}${PMB_ENDPOINTS.waves}`, {
        method: "POST",
        headers: { Authorization: `Bearer ${token}`, "Content-Type": "application/json" },
        body: JSON.stringify(waveData),
      });

      if (!response.ok) throw new Error("Failed to create wave");
      await fetchWaves();
      return true;
    } catch (err) {
      setError(err instanceof Error ? err.message : "Error");
      return false;
    } finally {
      setIsLoading(false);
    }
  }, []);

  // Update wave
  const updateWave = useCallback(async (waveId: string, waveData: Partial<Wave>) => {
    setIsLoading(true);
    try {
      const token = getToken();
      if (!token) throw new Error("Not authenticated");

      const response = await fetch(`${API_BASE_URLS.pmb}${PMB_ENDPOINTS.waves}/${waveId}`, {
        method: "PUT",
        headers: { Authorization: `Bearer ${token}`, "Content-Type": "application/json" },
        body: JSON.stringify(waveData),
      });

      if (!response.ok) throw new Error("Failed to update wave");
      await fetchWaves();
      return true;
    } catch (err) {
      setError(err instanceof Error ? err.message : "Error");
      return false;
    } finally {
      setIsLoading(false);
    }
  }, []);

  return {
    applicants,
    waves,
    stats,
    documents,
    selectionResults,
    isLoading,
    error,
    fetchApplicants,
    fetchWaves,
    fetchStats,
    fetchDocuments,
    fetchSelectionResults,
    updateApplicantStatus,
    verifyDocument,
    inputSelectionResult,
    createWave,
    updateWave,
  };
}
