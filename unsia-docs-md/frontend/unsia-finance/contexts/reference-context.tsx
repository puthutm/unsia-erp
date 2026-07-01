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

export interface PaymentComponent {
  id: string;
  code: string;
  name: string;
  defaultAmount: number;
  isActive: boolean;
}

export interface PaymentMethod {
  id: string;
  code: string;
  name: string;
  isActive: boolean;
}

interface ReferenceContextType {
  studyPrograms: StudyProgram[];
  academicYears: AcademicYear[];
  academicPeriods: AcademicPeriod[];
  paymentComponents: PaymentComponent[];
  paymentMethods: PaymentMethod[];
  isLoading: boolean;
  error: string | null;
  fetchStudyPrograms: () => Promise<void>;
  fetchAcademicYears: () => Promise<void>;
  fetchAcademicPeriods: () => Promise<void>;
  fetchPaymentComponents: () => Promise<void>;
  fetchPaymentMethods: () => Promise<void>;
  fetchAll: () => Promise<void>;
}

const ReferenceContext = createContext<ReferenceContextType | undefined>(undefined);

export function ReferenceProvider({ children }: { children: ReactNode }) {
  const { isAuthenticated } = useAuth();
  
  const [studyPrograms, setStudyPrograms] = useState<StudyProgram[]>([]);
  const [academicYears, setAcademicYears] = useState<AcademicYear[]>([]);
  const [academicPeriods, setAcademicPeriods] = useState<AcademicPeriod[]>([]);
  const [paymentComponents, setPaymentComponents] = useState<PaymentComponent[]>([]);
  const [paymentMethods, setPaymentMethods] = useState<PaymentMethod[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const getAccessToken = () => {
    return localStorage.getItem(STORAGE_KEYS.accessToken);
  };

  const fetchStudyPrograms = async () => {
    const token = getAccessToken();
    if (!token) return;
    try {
      const response = await fetch(`${API_BASE_URLS.reference}${REFERENCE_ENDPOINTS.studyPrograms}`, {
        headers: { Authorization: `Bearer ${token}` },
      });
      if (response.ok) {
        const result = await response.json();
        setStudyPrograms(result.data || []);
      }
    } catch (err) {
      console.error("Failed to fetch study programs", err);
    }
  };

  const fetchAcademicYears = async () => {
    const token = getAccessToken();
    if (!token) return;
    try {
      const response = await fetch(`${API_BASE_URLS.reference}${REFERENCE_ENDPOINTS.academicYears}`, {
        headers: { Authorization: `Bearer ${token}` },
      });
      if (response.ok) {
        const result = await response.json();
        setAcademicYears(result.data || []);
      }
    } catch (err) {
      console.error("Failed to fetch academic years", err);
    }
  };

  const fetchAcademicPeriods = async () => {
    const token = getAccessToken();
    if (!token) return;
    try {
      const response = await fetch(`${API_BASE_URLS.reference}${REFERENCE_ENDPOINTS.academicPeriods}`, {
        headers: { Authorization: `Bearer ${token}` },
      });
      if (response.ok) {
        const result = await response.json();
        setAcademicPeriods(result.data || []);
      }
    } catch (err) {
      console.error("Failed to fetch academic periods", err);
    }
  };

  const fetchPaymentComponents = async () => {
    const token = getAccessToken();
    if (!token) return;
    try {
      const response = await fetch(`${API_BASE_URLS.reference}${REFERENCE_ENDPOINTS.paymentComponents}`, {
        headers: { Authorization: `Bearer ${token}` },
      });
      if (response.ok) {
        const result = await response.json();
        setPaymentComponents(result.data || []);
      }
    } catch (err) {
      console.error("Failed to fetch payment components", err);
    }
  };

  const fetchPaymentMethods = async () => {
    const token = getAccessToken();
    if (!token) return;
    try {
      const response = await fetch(`${API_BASE_URLS.reference}${REFERENCE_ENDPOINTS.paymentMethods}`, {
        headers: { Authorization: `Bearer ${token}` },
      });
      if (response.ok) {
        const result = await response.json();
        setPaymentMethods(result.data || []);
      }
    } catch (err) {
      console.error("Failed to fetch payment methods", err);
    }
  };

  const fetchAll = async () => {
    setIsLoading(true);
    setError(null);
    try {
      await Promise.all([
        fetchStudyPrograms(),
        fetchAcademicYears(),
        fetchAcademicPeriods(),
        fetchPaymentComponents(),
        fetchPaymentMethods(),
      ]);
    } catch (err) {
      setError(err instanceof Error ? err.message : "An error occurred");
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    if (isAuthenticated) {
      fetchAll();
    }
  }, [isAuthenticated]);

  return (
    <ReferenceContext.Provider
      value={{
        studyPrograms,
        academicYears,
        academicPeriods,
        paymentComponents,
        paymentMethods,
        isLoading,
        error,
        fetchStudyPrograms,
        fetchAcademicYears,
        fetchAcademicPeriods,
        fetchPaymentComponents,
        fetchPaymentMethods,
        fetchAll,
      }}
    >
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
