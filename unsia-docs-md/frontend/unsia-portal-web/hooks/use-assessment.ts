"use client";

import { useState, useCallback } from "react";
import { API_BASE_URLS, ASSESSMENT_ENDPOINTS, STORAGE_KEYS } from "@/lib/constants";

export interface CBTSession {
  id: string;
  code: string;
  title: string;
  description?: string;
  startTime: string;
  endTime: string;
  durationMinutes: number;
  status: string;
}

export interface CBTQuestion {
  id: string;
  text: string;
  options: {
    key: string; // "A", "B", "C", "D", "E"
    text: string;
  }[];
}

export interface CBTAttempt {
  id: string;
  sessionId: string;
  participantId: string;
  startTime: string;
  endTime?: string;
  status: string; // "STARTED", "COMPLETED"
}

export function useAssessment() {
  const [sessions, setSessions] = useState<CBTSession[]>([]);
  const [currentAttempt, setCurrentAttempt] = useState<CBTAttempt | null>(null);
  const [questions, setQuestions] = useState<CBTQuestion[]>([]);
  const [answers, setAnswers] = useState<Record<string, string>>({}); // questionId -> selectedOptionKey
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const getToken = () => localStorage.getItem(STORAGE_KEYS.accessToken);

  const fetchSessions = useCallback(async () => {
    setIsLoading(true);
    setError(null);
    try {
      const token = getToken();
      if (!token) throw new Error("Not authenticated");

      const url = `${API_BASE_URLS.assessment}${ASSESSMENT_ENDPOINTS.sessions}`;
      const response = await fetch(url, {
        headers: { Authorization: `Bearer ${token}`, "Content-Type": "application/json" },
      });

      if (!response.ok) throw new Error("Failed to fetch sessions");
      const data = await response.json();
      setSessions(data.data || []);
      return data.data || [];
    } catch (err) {
      setError(err instanceof Error ? err.message : "Error");
      return [];
    } finally {
      setIsLoading(false);
    }
  }, []);

  const registerParticipant = useCallback(async (sessionId: string, personId: string) => {
    setIsLoading(true);
    setError(null);
    try {
      const token = getToken();
      if (!token) throw new Error("Not authenticated");

      const response = await fetch(`${API_BASE_URLS.assessment}${ASSESSMENT_ENDPOINTS.participants}`, {
        method: "POST",
        headers: { Authorization: `Bearer ${token}`, "Content-Type": "application/json" },
        body: JSON.stringify({ sessionId, personId }),
      });

      return response.ok;
    } catch (err) {
      setError(err instanceof Error ? err.message : "Error");
      return false;
    } finally {
      setIsLoading(false);
    }
  }, []);

  const startAttempt = useCallback(async (sessionId: string) => {
    setIsLoading(true);
    setError(null);
    try {
      const token = getToken();
      if (!token) throw new Error("Not authenticated");

      const response = await fetch(`${API_BASE_URLS.assessment}${ASSESSMENT_ENDPOINTS.attempts}`, {
        method: "POST",
        headers: { Authorization: `Bearer ${token}`, "Content-Type": "application/json" },
        body: JSON.stringify({ sessionId }),
      });

      if (!response.ok) throw new Error("Failed to start CBT attempt");
      const resData = await response.json();
      setCurrentAttempt(resData.data.attempt || resData.data);
      
      // Mock questions for the UI CBT exam room simulation
      const mockQuestions: CBTQuestion[] = [
        {
          id: "q1",
          text: "Manakah yang merupakan pilar dari Pemrograman Berorientasi Objek (OOP)?",
          options: [
            { key: "A", text: "Encapsulation, Inheritance, Polymorphism, Abstraction" },
            { key: "B", text: "Compilation, Interpretation, Execution, Debugging" },
            { key: "C", text: "HTML, CSS, JavaScript, SQL" },
            { key: "D", text: "GET, POST, PUT, DELETE" }
          ]
        },
        {
          id: "q2",
          text: "Apakah singkatan dari SQL dalam pengelolaan database?",
          options: [
            { key: "A", text: "Simple Query Language" },
            { key: "B", text: "Structured Query Language" },
            { key: "C", text: "System Query List" },
            { key: "D", text: "Sequential Query Layout" }
          ]
        },
        {
          id: "q3",
          text: "Protokol apa yang digunakan untuk transfer dokumen web secara aman (encrypted)?",
          options: [
            { key: "A", text: "HTTP" },
            { key: "B", text: "FTP" },
            { key: "C", text: "HTTPS" },
            { key: "D", text: "SMTP" }
          ]
        },
        {
          id: "q4",
          text: "Di bawah ini, manakah yang merupakan database NoSQL?",
          options: [
            { key: "A", text: "MySQL" },
            { key: "B", text: "PostgreSQL" },
            { key: "C", text: "MongoDB" },
            { key: "D", text: "Oracle" }
          ]
        },
        {
          id: "q5",
          text: "Fungsi utama dari DNS (Domain Name System) adalah untuk...",
          options: [
            { key: "A", text: "Mengamankan jaringan dari virus" },
            { key: "B", text: "Menerjemahkan nama domain menjadi IP address" },
            { key: "C", text: "Menyimpan file website" },
            { key: "D", text: "Mempercepat koneksi internet" }
          ]
        }
      ];
      setQuestions(mockQuestions);
      setAnswers({});
      return resData.data.attempt || resData.data;
    } catch (err) {
      setError(err instanceof Error ? err.message : "Error");
      return null;
    } finally {
      setIsLoading(false);
    }
  }, []);

  const saveAnswer = useCallback(async (attemptId: string, questionId: string, answerKey: string) => {
    setAnswers(prev => ({
      ...prev,
      [questionId]: answerKey
    }));

    try {
      const token = getToken();
      if (!token) return false;

      const response = await fetch(`${API_BASE_URLS.assessment}${ASSESSMENT_ENDPOINTS.attempts}/${attemptId}/answers`, {
        method: "POST",
        headers: { Authorization: `Bearer ${token}`, "Content-Type": "application/json" },
        body: JSON.stringify({ questionId, answerKey }),
      });
      return response.ok;
    } catch (err) {
      console.error("Failed to save answer to backend, saved locally instead", err);
      return true;
    }
  }, []);

  const submitAttempt = useCallback(async (attemptId: string) => {
    setIsLoading(true);
    setError(null);
    try {
      const token = getToken();
      if (!token) throw new Error("Not authenticated");

      const response = await fetch(`${API_BASE_URLS.assessment}${ASSESSMENT_ENDPOINTS.attempts}/${attemptId}/submit`, {
        method: "POST",
        headers: { Authorization: `Bearer ${token}`, "Content-Type": "application/json" },
        body: JSON.stringify({ answers }),
      });

      if (!response.ok) throw new Error("Failed to submit exam attempt");
      setCurrentAttempt(null);
      setQuestions([]);
      setAnswers({});
      return true;
    } catch (err) {
      setError(err instanceof Error ? err.message : "Error");
      return false;
    } finally {
      setIsLoading(false);
    }
  }, [answers]);

  return {
    sessions,
    currentAttempt,
    questions,
    answers,
    isLoading,
    error,
    fetchSessions,
    registerParticipant,
    startAttempt,
    saveAnswer,
    submitAttempt,
  };
}
