"use client";

import { useState, useCallback } from "react";
import { API_BASE_URLS, CRM_ENDPOINTS, STORAGE_KEYS } from "@/lib/constants";

export interface Lead {
	id: string;
	leadNumber: string;
	personId: string;
	studyProgramId?: string;
	status: string;
	createdAt: string;
	convertedAt?: string;
	campaignId?: string;
}

export interface Campaign {
	id: string;
	code: string;
	name: string;
	channel: string;
	status: string;
}

export interface Agent {
	id: string;
	personId: string;
	agentCode: string;
	organizationName: string;
	status: string;
	approvalStatus: string;
}

export function useCRM() {
	const [leads, setLeads] = useState<Lead[]>([]);
	const [campaigns, setCampaigns] = useState<Campaign[]>([]);
	const [agents, setAgents] = useState<Agent[]>([]);
	const [isLoading, setIsLoading] = useState(false);
	const [error, setError] = useState<string | null>(null);

	const getToken = () => localStorage.getItem(STORAGE_KEYS.accessToken);

	const fetchLeads = useCallback(async () => {
		setIsLoading(true);
		setError(null);
		try {
			const token = getToken();
			if (!token) throw new Error("Not authenticated");

			const url = `${API_BASE_URLS.crm}${CRM_ENDPOINTS.leads}`;
			const response = await fetch(url, {
				headers: { Authorization: `Bearer ${token}`, "Content-Type": "application/json" },
			});

			if (!response.ok) throw new Error("Failed to fetch leads");
			const data = await response.json();
			setLeads(data.data || []);
			return data.data || [];
		} catch (err) {
			setError(err instanceof Error ? err.message : "Error");
			return [];
		} finally {
			setIsLoading(false);
		}
	}, []);

	const createLead = useCallback(async (leadData: Partial<Lead>) => {
		setIsLoading(true);
		setError(null);
		try {
			const token = getToken();
			if (!token) throw new Error("Not authenticated");

			const response = await fetch(`${API_BASE_URLS.crm}${CRM_ENDPOINTS.leads}`, {
				method: "POST",
				headers: { Authorization: `Bearer ${token}`, "Content-Type": "application/json" },
				body: JSON.stringify(leadData),
			});

			if (!response.ok) throw new Error("Failed to create lead");
			await fetchLeads();
			return true;
		} catch (err) {
			setError(err instanceof Error ? err.message : "Error");
			return false;
		} finally {
			setIsLoading(false);
		}
	}, [fetchLeads]);

	const logLeadActivity = useCallback(async (leadId: string, activity: { activity_type: string; note: string }) => {
		setIsLoading(true);
		try {
			const token = getToken();
			if (!token) throw new Error("Not authenticated");

			const response = await fetch(`${API_BASE_URLS.crm}${CRM_ENDPOINTS.leads}/${leadId}/activities`, {
				method: "POST",
				headers: { Authorization: `Bearer ${token}`, "Content-Type": "application/json" },
				body: JSON.stringify(activity),
			});

			return response.ok;
		} catch (err) {
			setError(err instanceof Error ? err.message : "Error");
			return false;
		} finally {
			setIsLoading(false);
		}
	}, []);

	const convertLead = useCallback(async (leadId: string, convertData: any) => {
		setIsLoading(true);
		try {
			const token = getToken();
			if (!token) throw new Error("Not authenticated");

			const response = await fetch(`${API_BASE_URLS.crm}${CRM_ENDPOINTS.leads}/${leadId}/convert-to-applicant`, {
				method: "POST",
				headers: { Authorization: `Bearer ${token}`, "Content-Type": "application/json" },
				body: JSON.stringify(convertData),
			});

			if (response.ok) {
				await fetchLeads();
				return true;
			}
			return false;
		} catch (err) {
			setError(err instanceof Error ? err.message : "Error");
			return false;
		} finally {
			setIsLoading(false);
		}
	}, [fetchLeads]);

	const fetchCampaigns = useCallback(async () => {
		setIsLoading(true);
		try {
			const token = getToken();
			if (!token) throw new Error("Not authenticated");

			const response = await fetch(`${API_BASE_URLS.crm}${CRM_ENDPOINTS.campaigns}`, {
				headers: { Authorization: `Bearer ${token}` },
			});
			if (!response.ok) throw new Error("Failed to fetch campaigns");
			const data = await response.json();
			setCampaigns(data.data || []);
		} catch (err) {
			setError(err instanceof Error ? err.message : "Error");
		} finally {
			setIsLoading(false);
		}
	}, []);

	const createCampaign = useCallback(async (campaignData: Partial<Campaign>) => {
		setIsLoading(true);
		try {
			const token = getToken();
			if (!token) throw new Error("Not authenticated");

			const response = await fetch(`${API_BASE_URLS.crm}${CRM_ENDPOINTS.campaigns}`, {
				method: "POST",
				headers: { Authorization: `Bearer ${token}`, "Content-Type": "application/json" },
				body: JSON.stringify(campaignData),
			});
			if (response.ok) {
				await fetchCampaigns();
				return true;
			}
			return false;
		} catch (err) {
			setError(err instanceof Error ? err.message : "Error");
			return false;
		} finally {
			setIsLoading(false);
		}
	}, [fetchCampaigns]);

	const fetchAgents = useCallback(async () => {
		setIsLoading(true);
		try {
			const token = getToken();
			if (!token) throw new Error("Not authenticated");

			const response = await fetch(`${API_BASE_URLS.crm}${CRM_ENDPOINTS.agents}`, {
				headers: { Authorization: `Bearer ${token}` },
			});
			if (!response.ok) throw new Error("Failed to fetch agents");
			const data = await response.json();
			setAgents(data.data || []);
		} catch (err) {
			setError(err instanceof Error ? err.message : "Error");
		} finally {
			setIsLoading(false);
		}
	}, []);

	const registerAgent = useCallback(async (agentData: Partial<Agent>) => {
		setIsLoading(true);
		try {
			const token = getToken();
			if (!token) throw new Error("Not authenticated");

			const response = await fetch(`${API_BASE_URLS.crm}${CRM_ENDPOINTS.agents}`, {
				method: "POST",
				headers: { Authorization: `Bearer ${token}`, "Content-Type": "application/json" },
				body: JSON.stringify(agentData),
			});
			if (response.ok) {
				await fetchAgents();
				return true;
			}
			return false;
		} catch (err) {
			setError(err instanceof Error ? err.message : "Error");
			return false;
		} finally {
			setIsLoading(false);
		}
	}, [fetchAgents]);

	return {
		leads,
		campaigns,
		agents,
		isLoading,
		error,
		fetchLeads,
		createLead,
		logLeadActivity,
		convertLead,
		fetchCampaigns,
		createCampaign,
		fetchAgents,
		registerAgent,
	};
}
