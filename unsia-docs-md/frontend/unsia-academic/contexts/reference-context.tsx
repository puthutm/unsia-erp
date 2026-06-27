"use client";

import { createContext, useContext, useState, useEffect, ReactNode } from "react";
import { REFERENCE_ENDPOINTS, API_BASE_URLS, STORAGE_KEYS } from "@/lib/constants";
import { useAuth } from "./auth-context";

// Types for reference data
export interface StudyProgram {
  id: string;
  code: string;
  name: string;
  degree: string;
  status: string;
}

export interface AcademicYear {
  id: string;
  code: string;
  name: string;
  status: string;
  startDate: string;
  endDate: string;
}

export interface AcademicPeriod {
  id: string;
  academicYearId: string;
  code: string;
  term: string;
  status: string;
  startDate: string;
  endDate: string;
}

interface ReferenceContextType {
  studyPrograms: StudyProgram[];
  academicYears: AcademicYear[];
  academicPeriods: AcademicPeriod[];
  isLoading: boolean;
  error: string | null;
  fetchStudyPrograms: () => Promise<void>;
  fetchAcademicYears: () => Promise<void>;
  fetchAcademicPeriods: () => Promise<void>;
  fetchAll: () => Promise<void>;
}

const ReferenceContext = createContext<ReferenceContextType | undefined>(undefined);

export function ReferenceProvider({ children }: { children: ReactNode }) {
  const { isAuthenticated } = useAuth();
  
  const [studyPrograms, setStudyPrograms] = useState<StudyProgram[]>([]);
  const [academicYears, setAcademicYears] = useState<AcademicYear[]>([]);
  const [academicPeriods, setAcademicPeriods] = useState<AcademicPeriod[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const getAccessToken = () => {
    return localStorage.getItem(STORAGE_KEYS.accessToken);
  };

  const fetchData = async (endpoint: string) => {
    const token = getAccessToken();
    const response = await fetch(`${API_BASE_URLS.reference}${endpoint}`, {
      headers: {
        Authorization: token ? `Bearer ${token}` : "",
        "Content-Type": "application/json",
      },
    });
    
    if (!response.ok) {
      throw new Error(`Failed to fetch ${endpoint}`);
    }
    
    const result = await response.json();
    return result.data || [];
  };

  const fetchStudyPrograms = async () => {
    try {
      const data = await fetchData(REFERENCE_ENDPOINTS.studyPrograms);
      setStudyPrograms(data);
    } catch (err) {
      console.error("Error fetching study programs:", err);
    }
  };

  const fetchAcademicYears = async () => {
    try {
      const data = await fetchData(REFERENCE_ENDPOINTS.academicYears);
      setAcademicYears(data);
    } catch (err) {
      console.error("Error fetching academic years:", err);
    }
  };

  const fetchAcademicPeriods = async () => {
    try {
      const data = await fetchData(REFERENCE_ENDPOINTS.academicPeriods);
      setAcademicPeriods(data);
    } catch (err) {
      console.error("Error fetching academic periods:", err);
    }
  };

  const fetchAll = async () => {
    if (!isAuthenticated) return;
    
    setIsLoading(true);
    setError(null);
    
    try {
      await Promise.all([
        fetchStudyPrograms(),
        fetchAcademicYears(),
        fetchAcademicPeriods(),
      ]);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to fetch reference data");
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    if (isAuthenticated) {
      fetchAll();
    }
  }, [isAuthenticated]);

  const value: ReferenceContextType = {
    studyPrograms,
    academicYears,
    academicPeriods,
    isLoading,
    error,
    fetchStudyPrograms,
    fetchAcademicYears,
    fetchAcademicPeriods,
    fetchAll,
  };

  return (
    <ReferenceContext.Provider value={value}>
      {children}
    </ReferenceContext.Provider>
  );
}

export function useReference() {
  const context = useContext(ReferenceContext);
  if (context === undefined) {
    throw new Error("useReference must be used within a ReferenceProvider");
  }
  return context;
}
