// API Configuration for UNSIA Backend Services

// Service URLs - these should match the backend service ports
export const API_BASE_URLS = {
  auth: process.env.NEXT_PUBLIC_API_CORE_URL || "http://localhost:8001",
  reference: process.env.NEXT_PUBLIC_API_REFERENCE_URL || "http://localhost:8002",
  academic: process.env.NEXT_PUBLIC_API_ACADEMIC_URL || "http://localhost:8004",
  finance: process.env.NEXT_PUBLIC_API_FINANCE_URL || "http://localhost:8005",
  hris: process.env.NEXT_PUBLIC_HRIS_API || "http://localhost:8008",
} as const;

// Auth API endpoints (from unsia-core-service)
export const AUTH_ENDPOINTS = {
  login: "/api/v1/auth/login",
  refresh: "/api/v1/auth/refresh",
  me: "/api/v1/auth/me",
  switchRole: "/api/v1/auth/switch-role",
  applications: "/api/v1/auth/applications",
} as const;

// Reference API endpoints (from unsia-reference-service)
export const REFERENCE_ENDPOINTS = {
  studyPrograms: "/api/v1/reference/study-programs",
  academicYears: "/api/v1/reference/academic-years",
  academicPeriods: "/api/v1/reference/academic-periods",
  statusCodes: "/api/v1/reference/status-codes",
} as const;

// Academic API endpoints (from unsia-academic-service)
export const ACADEMIC_ENDPOINTS = {
  students: "/api/v1/academic/students",
  courses: "/api/v1/academic/courses",
  krs: "/api/v1/academic/krs",
  grades: "/api/v1/academic/grades",
  schedules: "/api/v1/academic/schedules",
  transcripts: "/api/v1/academic/transcripts",
  graduation: "/api/v1/academic/graduation",
  clearance: "/api/v1/academic/clearance",
} as const;

// HRIS API endpoints (from unsia-hris-service)
export const HRIS_ENDPOINTS = {
  employees: "/api/v1/hris/employees",
  lecturers: "/api/v1/hris/lecturers",
  attendances: "/api/v1/hris/attendances",
} as const;

// Token storage keys
export const STORAGE_KEYS = {
  accessToken: "unsia_access_token",
  refreshToken: "unsia_refresh_token",
  user: "unsia_user",
} as const;

// API Response types
export interface ApiResponse<T> {
  data: T;
  success: boolean;
  message?: string;
}

export interface TokenResponse {
  access_token: string;
  refresh_token: string;
  expires_in: number;
  token_type: string;
}

export interface UserInfo {
  user_id: string;
  person_id: string;
  name: string;
  email: string;
  active_role: string;
  permissions: string[];
  scope: string;
}
