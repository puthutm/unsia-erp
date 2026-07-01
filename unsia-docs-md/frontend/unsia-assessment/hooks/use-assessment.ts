"use client";

import { useState, useCallback } from "react";

export interface AssessmentSession {
  id: string;
  code: string;
  title: string;
  description?: string;
  durationMinutes: number;
  status?: string;
}

export interface Participant {
  id: string;
  sessionId: string;
  personId: string;
  status: string;
}

const API_BASE_URL = process.env.NEXT_PUBLIC_ASSESSMENT_API || "http://localhost:8007";
const STORAGE_KEY = "unsia_access_token";

export function useAssessment() {
  const [sessions, setSessions] = useState<AssessmentSession[]>([]);
  const [participants, setParticipants] = useState<Participant[]>([]);
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

  const fetchSessions = useCallback(async () => {
    setIsLoading(true);
    setError(null);
    try {
      const res = await request(`${API_BASE_URL}/api/v1/assessment/sessions`);
      const data = res.data || [];
      setSessions(data);
      return data;
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to fetch sessions");
      return [];
    } finally {
      setIsLoading(false);
    }
  }, []);

  const registerParticipant = useCallback(async (sessionId: string, personId: string) => {
    setIsLoading(true);
    setError(null);
    try {
      await request(`${API_BASE_URL}/api/v1/assessment/sessions/${sessionId}/participants`, {
        method: "POST",
        body: JSON.stringify({ personId }),
      });
      return true;
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to register participant");
      return false;
    } finally {
      setIsLoading(false);
    }
  }, []);

  const fetchParticipants = useCallback(async (sessionId: string) => {
    setIsLoading(true);
    setError(null);
    try {
      const res = await request(`${API_BASE_URL}/api/v1/assessment/sessions/${sessionId}/participants`);
      const data = res.data || [];
      setParticipants(data);
      return data;
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to fetch participants");
      return [];
    } finally {
      setIsLoading(false);
    }
  }, []);

  return {
    sessions,
    participants,
    isLoading,
    error,
    fetchSessions,
    registerParticipant,
    fetchParticipants,
  };
}
