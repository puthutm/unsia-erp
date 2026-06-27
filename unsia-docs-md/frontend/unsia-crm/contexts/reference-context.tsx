"use client";

import { createContext, useContext, useState, ReactNode } from "react";
import { API_BASE_URLS, REFERENCE_ENDPOINTS, STORAGE_KEYS } from "@/lib/constants";

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
  code: string;
  name: string;
  status: string;
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
    if (typeof window === "undefined") return null;
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

  const fetchStudyPrograms = async () => setStudyPrograms(await fetchData(REFERENCE_ENDPOINTS.studyPrograms));
  const fetchAcademicYears = async () => setAcademicYears(await fetchData(REFERENCE_ENDPOINTS.academicYears));
  const fetchAcademicPeriods = async () => setAcademicPeriods(await fetchData(REFERENCE_ENDPOINTS.academicPeriods));
  const fetchPaymentComponents = async () => setPaymentComponents(await fetchData(REFERENCE_ENDPOINTS.paymentComponents));
  const fetchPaymentMethods = async () => setPaymentMethods(await fetchData(REFERENCE_ENDPOINTS.paymentMethods));
  const fetchDocumentTypes = async () => setDocumentTypes(await fetchData(REFERENCE_ENDPOINTS.documentTypes));
  const fetchPmbWaves = async () => setPmbWaves(await fetchData(REFERENCE_ENDPOINTS.pmbWaves));
  const fetchProvinces = async () => setProvinces(await fetchData(REFERENCE_ENDPOINTS.provinces));
  const fetchReligions = async () => setReligions(await fetchData(REFERENCE_ENDPOINTS.religions));
  const fetchAdmissionPaths = async () => setAdmissionPaths(await fetchData(REFERENCE_ENDPOINTS.admissionPaths));

  const fetchCities = async (provinceId?: string) => {
    const endpoint = provinceId ? `${REFERENCE_ENDPOINTS.cities}?province_id=${provinceId}` : REFERENCE_ENDPOINTS.cities;
    setCities(await fetchData(endpoint));
  };

  const fetchDistricts = async (cityId?: string) => {
    const endpoint = cityId ? `${REFERENCE_ENDPOINTS.districts}?city_id=${cityId}` : REFERENCE_ENDPOINTS.districts;
    setDistricts(await fetchData(endpoint));
  };

  const fetchVillages = async (districtId?: string) => {
    const endpoint = districtId ? `${REFERENCE_ENDPOINTS.villages}?district_id=${districtId}` : REFERENCE_ENDPOINTS.villages;
    setVillages(await fetchData(endpoint));
  };

  const fetchAll = async () => {
    const token = getAccessToken();
    if (!token) return;

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

  return <ReferenceContext.Provider value={value}>{children}</ReferenceContext.Provider>;
}

export function useReference() {
  const context = useContext(ReferenceContext);
  if (context === undefined) {
    throw new Error("useReference must be used within a ReferenceProvider");
  }
  return context;
}
