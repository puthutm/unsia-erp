"use client";

import { useState, useCallback } from "react";

export interface Lead {
  id: string;
  name: string;
  email: string;
  phone: string;
  company?: string;
  status: string;
  source: string;
  assignedTo?: string;
  createdAt: string;
  updatedAt: string;
  // Additional fields used in page
  leadNumber?: string;
  studyProgramId?: string;
  campaignId?: string;
}

export interface Campaign {
  id: string;
  code: string;
  name: string;
  description: string;
  startDate: string;
  endDate: string;
  budget: number;
  status: string;
  createdAt: string;
  channel: string;
}

export interface LeadActivity {
  id: string;
  leadId: string;
  type: string;
  note: string;
  createdAt: string;
}

export interface ConvertLeadData {
  opportunityName: string;
  value: number;
  admissionPathId?: string;
  // Support both camelCase and snake_case for API compatibility
  admission_path_id?: string;
}

export interface Pipeline {
  id: string;
  name: string;
  stages: PipelineStage[];
}

export interface PipelineStage {
  id: string;
  name: string;
  order: number;
  color: string;
}

export interface Agent {
  id: string;
  agentCode: string;
  organizationName: string;
  status: string;
  approvalStatus: string;
  personId?: string;
  createdAt: string;
  updatedAt: string;
}

const API_BASE_URL = process.env.NEXT_PUBLIC_CRM_API || "http://localhost:8083";
const STORAGE_KEY = "unsia_access_token";

export function useCRM() {
  const [leads, setLeads] = useState<Lead[]>([]);
  const [pipelines, setPipelines] = useState<Pipeline[]>([]);
  const [campaigns, setCampaigns] = useState<Campaign[]>([]);
  const [agents, setAgents] = useState<Agent[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const getToken = () => localStorage.getItem(STORAGE_KEY);

  const request = async (url: string, options: RequestInit = {}) => {
    const token = getToken();
    if (!token) throw new Error("Not authenticated");

    const response = await fetch(url, {
      ...options,
      headers: {
        Authorization: `Bearer ${token}`,
        "Content-Type": "application/json",
        ...(options.headers || {}),
      },
    });

    if (!response.ok) {
      throw new Error(`Request failed: ${response.status}`);
    }

    return response.json();
  };

  const fetchLeads = useCallback(async () => {
    setIsLoading(true);
    setError(null);
    try {
      const res = await request(`${API_BASE_URL}/api/v1/crm/leads`);
      const data = res.data || [];
      setLeads(data);
      return data;
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to fetch leads");
      return [];
    } finally {
      setIsLoading(false);
    }
  }, []);

  const createLead = useCallback(async (lead: Partial<Lead>) => {
    setIsLoading(true);
    setError(null);
    try {
      const res = await request(`${API_BASE_URL}/api/v1/crm/leads`, {
        method: "POST",
        body: JSON.stringify(lead),
      });
      return res.data;
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to create lead");
      return null;
    } finally {
      setIsLoading(false);
    }
  }, []);

  const updateLead = useCallback(async (id: string, lead: Partial<Lead>) => {
    setIsLoading(true);
    setError(null);
    try {
      const res = await request(`${API_BASE_URL}/api/v1/crm/leads/${id}`, {
        method: "PUT",
        body: JSON.stringify(lead),
      });
      return res.data;
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to update lead");
      return null;
    } finally {
      setIsLoading(false);
    }
  }, []);

  const deleteLead = useCallback(async (id: string) => {
    setIsLoading(true);
    setError(null);
    try {
      await request(`${API_BASE_URL}/api/v1/crm/leads/${id}`, {
        method: "DELETE",
      });
      return true;
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to delete lead");
      return false;
    } finally {
      setIsLoading(false);
    }
  }, []);

  const fetchPipelines = useCallback(async () => {
    setIsLoading(true);
    setError(null);
    try {
      const res = await request(`${API_BASE_URL}/api/v1/crm/pipelines`);
      const data = res.data || [];
      setPipelines(data);
      return data;
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to fetch pipelines");
      return [];
    } finally {
      setIsLoading(false);
    }
  }, []);

  const fetchCampaigns = useCallback(async () => {
    setIsLoading(true);
    setError(null);
    try {
      const res = await request(`${API_BASE_URL}/api/v1/crm/campaigns`);
      const data = res.data || [];
      setCampaigns(data);
      return data;
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to fetch campaigns");
      return [];
    } finally {
      setIsLoading(false);
    }
  }, []);

  const logLeadActivity = useCallback(async (leadId: string, activity: { activity_type: string; note: string }) => {
    try {
      await request(`${API_BASE_URL}/api/v1/crm/leads/${leadId}/activities`, {
        method: "POST",
        body: JSON.stringify(activity),
      });
      return true;
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to log activity");
      return false;
    }
  }, []);

const convertLead = useCallback(async (leadId: string, data: ConvertLeadData) => {
    setIsLoading(true);
    setError(null);
    try {
      await request(`${API_BASE_URL}/api/v1/crm/leads/${leadId}/convert`, {
        method: "POST",
        body: JSON.stringify(data),
      });
      return true;
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to convert lead");
      return false;
    } finally {
      setIsLoading(false);
    }
  }, []);

  const fetchAgents = useCallback(async () => {
    setIsLoading(true);
    setError(null);
    try {
      const res = await request(`${API_BASE_URL}/api/v1/crm/agents`);
      const data = res.data || [];
      setAgents(data);
      return data;
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to fetch agents");
      return [];
    } finally {
      setIsLoading(false);
    }
  }, []);

const registerAgent = useCallback(async (agent: { agentCode: string; organizationName: string; status: string; approvalStatus: string; personId?: string }) => {
    setIsLoading(true);
    setError(null);
    try {
      const res = await request(`${API_BASE_URL}/api/v1/crm/agents`, {
        method: "POST",
        body: JSON.stringify(agent),
      });
      return res.data;
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to register agent");
      return null;
    } finally {
      setIsLoading(false);
    }
  }, []);

  const createCampaign = useCallback(async (campaign: { code: string; name: string; channel: string; status: string }) => {
    setIsLoading(true);
    setError(null);
    try {
      await request(`${API_BASE_URL}/api/v1/crm/campaigns`, {
        method: "POST",
        body: JSON.stringify(campaign),
      });
      await fetchCampaigns();
      return true;
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to create campaign");
      return false;
    } finally {
      setIsLoading(false);
    }
  }, [fetchCampaigns]);

return {
    leads,
    pipelines,
    campaigns,
    agents,
    isLoading,
    error,
    fetchLeads,
    createLead,
    updateLead,
    deleteLead,
    fetchPipelines,
    fetchCampaigns,
    logLeadActivity,
    convertLead,
    fetchAgents,
    registerAgent,
    createCampaign,
  };
}
