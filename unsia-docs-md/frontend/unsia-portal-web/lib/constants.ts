// API Configuration for UNSIA Backend Services

// Service URLs - these should match the backend service ports
// Auth: 8001 (unsia-core-service)
// Reference: 8007 (unsia-reference-service)  
// PMB: 8003 (unsia-pmb-service)
// Finance: 8005 (unsia-finance-service)
// Academic: 8004 (unsia-academic-service)
// LMS: 8006 (unsia-lms-service)

export const API_BASE_URLS = {
  auth: process.env.NEXT_PUBLIC_AUTH_API || "http://localhost:8001",
  reference: process.env.NEXT_PUBLIC_REFERENCE_API || "http://localhost:8002",
  pmb: process.env.NEXT_PUBLIC_PMB_API || "http://localhost:8003",
  finance: process.env.NEXT_PUBLIC_FINANCE_API || "http://localhost:8005",
  academic: process.env.NEXT_PUBLIC_ACADEMIC_API || "http://localhost:8004",
  lms: process.env.NEXT_PUBLIC_LMS_API || "http://localhost:8006",
  assessment: process.env.NEXT_PUBLIC_ASSESSMENT_API || "http://localhost:8007",
  hris: process.env.NEXT_PUBLIC_HRIS_API || "http://localhost:8008",
  crm: process.env.NEXT_PUBLIC_CRM_API || "http://localhost:8009",
  portal: process.env.NEXT_PUBLIC_PORTAL_API || "http://localhost:8010",
} as const;

export const FRONTEND_URLS = {
  portal: process.env.NEXT_PUBLIC_PORTAL_URL || "http://localhost:3000",
  pmb: process.env.NEXT_PUBLIC_PMB_URL || "http://localhost:3001",
  academic: process.env.NEXT_PUBLIC_ACADEMIC_URL || "http://localhost:3002",
  finance: process.env.NEXT_PUBLIC_FINANCE_URL || "http://localhost:3003",
  lms: process.env.NEXT_PUBLIC_LMS_URL || "http://localhost:3004",
  hris: process.env.NEXT_PUBLIC_HRIS_URL || "http://localhost:3005",
  assessment: process.env.NEXT_PUBLIC_ASSESSMENT_URL || "http://localhost:3006",
  crm: process.env.NEXT_PUBLIC_CRM_URL || "http://localhost:3007",
  reference: process.env.NEXT_PUBLIC_REFERENCE_URL || "http://localhost:3008",
  core: process.env.NEXT_PUBLIC_CORE_URL || "http://localhost:3009",
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
  paymentComponents: "/api/v1/reference/payment-components",
  paymentMethods: "/api/v1/reference/payment-methods",
  documentTypes: "/api/v1/reference/document-types",
  pmbWaves: "/api/v1/reference/pmb-waves",
  religions: "/api/v1/reference/religions",
  countries: "/api/v1/reference/countries",
  admissionPaths: "/api/v1/reference/admission-paths",
  provinces: "/api/v1/reference/provinces",
  cities: "/api/v1/reference/cities",
  districts: "/api/v1/reference/districts",
  villages: "/api/v1/reference/villages",
} as const;

// PMB API endpoints (from unsia-pmb-service)
export const PMB_ENDPOINTS = {
  applicants: "/api/v1/pmb/applicants",
  waves: "/api/v1/pmb/waves",
  selection: "/api/v1/pmb/selection",
  dashboard: "/api/v1/pmb/dashboard",
} as const;

// Finance API endpoints (from unsia-finance-service)
export const FINANCE_ENDPOINTS = {
  invoices: "/api/v1/finance/invoices",
  payments: "/api/v1/finance/payments",
  scholarships: "/api/v1/finance/scholarships",
  clearance: "/api/v1/finance/clearance",
  journals: "/api/v1/finance/journals",
  budgets: "/api/v1/finance/budgets",
  vendors: "/api/v1/finance/vendors",
  purchaseOrders: "/api/v1/finance/purchase-orders",
  reports: "/api/v1/finance/reports",
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

// LMS API endpoints (from unsia-lms-service)
export const LMS_ENDPOINTS = {
  courses: "/api/v1/lms/courses",
  classes: "/api/v1/lms/classes",
  enrollments: "/api/v1/lms/enrollments",
  sessions: "/api/v1/lms/sessions",
  materials: "/api/v1/lms/materials",
  assignments: "/api/v1/lms/assignments",
  attendance: "/api/v1/lms/attendance",
} as const;

// CRM API endpoints (from unsia-crm-service)
export const CRM_ENDPOINTS = {
  leads: "/api/v1/crm/leads",
  campaigns: "/api/v1/crm/campaigns",
  agents: "/api/v1/crm/agents",
  referrals: "/api/v1/crm/referrals",
  contacts: "/api/v1/crm/contacts",
  opportunities: "/api/v1/crm/opportunities",
} as const;

// HRIS API endpoints (from unsia-hris-service)
export const HRIS_ENDPOINTS = {
  employees: "/api/v1/hris/employees",
  lecturers: "/api/v1/hris/lecturers",
  attendances: "/api/v1/hris/attendances",
  leaveRequests: "/api/v1/hris/leave-requests",
  bkdRecords: "/api/v1/hris/bkd-records",
} as const;

// Assessment API endpoints (from unsia-assessment-service)
export const ASSESSMENT_ENDPOINTS = {
  sessions: "/api/v1/assessment/sessions",
  participants: "/api/v1/assessment/participants",
  attempts: "/api/v1/assessment/attempts",
  resultsPublish: "/api/v1/assessment/results/publish",
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
