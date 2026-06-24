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
  provider: string;
  isActive: boolean;
}

export interface DocumentType {
  id: string;
  code: string;
  name: string;
  isMandatory: boolean;
  isActive: boolean;
}

export interface PmbWave {
  id: string;
  academicYearId?: string;
  targetEntryPeriodId: string;
  admissionPathId?: string;
  code: string;
  name: string;
  status: string;
  startDate?: string;
  endDate?: string;
  registrationStartAt?: string;
  registrationEndAt?: string;
  isActive: boolean;
}

export interface Province {
  id: string;
  code: string;
  name: string;
}

export interface City {
  id: string;
  provinceId: string;
  code: string;
  name: string;
}

export interface District {
  id: string;
  cityId: string;
  code: string;
  name: string;
}

export interface Village {
  id: string;
  districtId: string;
  code: string;
  name: string;
}

export interface Religion {
  id: string;
  name: string;
}

export interface AdmissionPath {
  id: string;
  code: string;
  name: string;
  isActive: boolean;
}

interface ReferenceContextType {
  // States
  studyPrograms: StudyProgram[];
  academicYears: AcademicYear[];
  academicPeriods: AcademicPeriod[];
  paymentComponents: PaymentComponent[];
  paymentMethods: PaymentMethod[];
  documentTypes: DocumentType[];
  pmbWaves: PmbWave[];
  provinces: Province[];
  cities: City[];
  districts: District[];
  villages: Village[];
  religions: Religion[];
  admissionPaths: AdmissionPath[];
  isLoading: boolean;
  error: string | null;
  
  // Actions
  fetchStudyPrograms: () => Promise<void>;
  fetchAcademicYears: () => Promise<void>;
  fetchAcademicPeriods: () => Promise<void>;
  fetchPaymentComponents: () => Promise<void>;
  fetchPaymentMethods: () => Promise<void>;
  fetchDocumentTypes: () => Promise<void>;
  fetchPmbWaves: () => Promise<void>;
  fetchProvinces: () => Promise<void>;
  fetchCities: (provinceId?: string) => Promise<void>;
  fetchDistricts: (cityId?: string) => Promise<void>;
  fetchVillages: (districtId?: string) => Promise<void>;
  fetchReligions: () => Promise<void>;
  fetchAdmissionPaths: () => Promise<void>;
  fetchAll: () => Promise<void>;
}

const ReferenceContext = createContext<ReferenceContextType | undefined>(undefined);

export function ReferenceProvider({ children }: { children: ReactNode }) {
  const { isAuthenticated } = useAuth();
  
  // State
  const [studyPrograms, setStudyPrograms] = useState<StudyProgram[]>([]);
  const [academicYears, setAcademicYears] = useState<AcademicYear[]>([]);
  const [academicPeriods, setAcademicPeriods] = useState<AcademicPeriod[]>([]);
  const [paymentComponents, setPaymentComponents] = useState<PaymentComponent[]>([]);
  const [paymentMethods, setPaymentMethods] = useState<PaymentMethod[]>([]);
  const [documentTypes, setDocumentTypes] = useState<DocumentType[]>([]);
  const [pmbWaves, setPmbWaves] = useState<PmbWave[]>([]);
  const [provinces, setProvinces] = useState<Province[]>([]);
  const [cities, setCities] = useState<City[]>([]);
  const [districts, setDistricts] = useState<District[]>([]);
  const [villages, setVillages] = useState<Village[]>([]);
  const [religions, setReligions] = useState<Religion[]>([]);
  const [admissionPaths, setAdmissionPaths] = useState<AdmissionPath[]>([]);
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

  const fetchPaymentComponents = async () => {
    try {
      const data = await fetchData(REFERENCE_ENDPOINTS.paymentComponents);
      setPaymentComponents(data);
    } catch (err) {
      console.error("Error fetching payment components:", err);
    }
  };

  const fetchPaymentMethods = async () => {
    try {
      const data = await fetchData(REFERENCE_ENDPOINTS.paymentMethods);
      setPaymentMethods(data);
    } catch (err) {
      console.error("Error fetching payment methods:", err);
    }
  };

  const fetchDocumentTypes = async () => {
    try {
      const data = await fetchData(REFERENCE_ENDPOINTS.documentTypes);
      setDocumentTypes(data);
    } catch (err) {
      console.error("Error fetching document types:", err);
    }
  };

  const fetchPmbWaves = async () => {
    try {
      const data = await fetchData(REFERENCE_ENDPOINTS.pmbWaves);
      setPmbWaves(data);
    } catch (err) {
      console.error("Error fetching PMB waves:", err);
    }
  };

  const fetchProvinces = async () => {
    try {
      const data = await fetchData(REFERENCE_ENDPOINTS.provinces);
      setProvinces(data);
    } catch (err) {
      console.error("Error fetching provinces:", err);
    }
  };

  const fetchCities = async (provinceId?: string) => {
    try {
      const endpoint = provinceId 
        ? `${REFERENCE_ENDPOINTS.cities}?province_id=${provinceId}`
        : REFERENCE_ENDPOINTS.cities;
      const data = await fetchData(endpoint);
      setCities(data);
    } catch (err) {
      console.error("Error fetching cities:", err);
    }
  };

  const fetchDistricts = async (cityId?: string) => {
    try {
      const endpoint = cityId 
        ? `${REFERENCE_ENDPOINTS.districts}?city_id=${cityId}`
        : REFERENCE_ENDPOINTS.districts;
      const data = await fetchData(endpoint);
      setDistricts(data);
    } catch (err) {
      console.error("Error fetching districts:", err);
    }
  };

  const fetchVillages = async (districtId?: string) => {
    try {
      const endpoint = districtId 
        ? `${REFERENCE_ENDPOINTS.villages}?district_id=${districtId}`
        : REFERENCE_ENDPOINTS.villages;
      const data = await fetchData(endpoint);
      setVillages(data);
    } catch (err) {
      console.error("Error fetching villages:", err);
    }
  };

  const fetchReligions = async () => {
    try {
      const data = await fetchData(REFERENCE_ENDPOINTS.religions);
      setReligions(data);
    } catch (err) {
      console.error("Error fetching religions:", err);
    }
  };

  const fetchAdmissionPaths = async () => {
    try {
      const data = await fetchData(REFERENCE_ENDPOINTS.admissionPaths);
      setAdmissionPaths(data);
    } catch (err) {
      console.error("Error fetching admission paths:", err);
    }
  };

  // Fetch all reference data
  const fetchAll = async () => {
    if (!isAuthenticated) return;
    
    setIsLoading(true);
    setError(null);
    
    try {
      await Promise.all([
        fetchStudyPrograms(),
        fetchAcademicYears(),
        fetchAcademicPeriods(),
        fetchPaymentComponents(),
        fetchPaymentMethods(),
        fetchDocumentTypes(),
        fetchPmbWaves(),
        fetchProvinces(),
        fetchReligions(),
        fetchAdmissionPaths(),
      ]);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to fetch reference data");
    } finally {
      setIsLoading(false);
    }
  };

  // Fetch data when authenticated
  useEffect(() => {
    if (isAuthenticated) {
      fetchAll();
    }
  }, [isAuthenticated]);

  const value: ReferenceContextType = {
    studyPrograms,
    academicYears,
    academicPeriods,
    paymentComponents,
    paymentMethods,
    documentTypes,
    pmbWaves,
    provinces,
    cities,
    districts,
    villages,
    religions,
    admissionPaths,
    isLoading,
    error,
    fetchStudyPrograms,
    fetchAcademicYears,
    fetchAcademicPeriods,
    fetchPaymentComponents,
    fetchPaymentMethods,
    fetchDocumentTypes,
    fetchPmbWaves,
    fetchProvinces,
    fetchCities,
    fetchDistricts,
    fetchVillages,
    fetchReligions,
    fetchAdmissionPaths,
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
