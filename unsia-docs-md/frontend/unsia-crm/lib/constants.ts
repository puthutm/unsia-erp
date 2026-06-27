// API Configuration for UNSIA Backend Services

export const API_BASE_URLS = {
  auth: "http://localhost:8001",
  reference: "http://localhost:8002",
  pmb: "http://localhost:8003",
  finance: "http://localhost:8005",
  academic: "http://localhost:8004",
  lms: "http://localhost:8006",
  assessment: "http://localhost:8007",
  hris: "http://localhost:8008",
  crm: "http://localhost:8009",
  portal: "http://localhost:8010",
} as const;

export const CRM_ENDPOINTS = {
  leads: "/api/v1/crm/leads",
  campaigns: "/api/v1/crm/campaigns",
  agents: "/api/v1/crm/agents",
  referrals: "/api/v1/crm/referrals",
  contacts: "/api/v1/crm/contacts",
  opportunities: "/api/v1/crm/opportunities",
} as const;

export const REFERENCE_ENDPOINTS = {
  studyPrograms: "/api/v1/reference/study-programs",
  academicYears: "/api/v1/reference/academic-years",
  academicPeriods: "/api/v1/reference/academic-periods",
  paymentComponents: "/api/v1/reference/payment-components",
  paymentMethods: "/api/v1/reference/payment-methods",
  documentTypes: "/api/v1/reference/document-types",
  pmbWaves: "/api/v1/reference/pmb-waves",
  religions: "/api/v1/reference/religions",
  admissionPaths: "/api/v1/reference/admission-paths",
  provinces: "/api/v1/reference/provinces",
  cities: "/api/v1/reference/cities",
  districts: "/api/v1/reference/districts",
  villages: "/api/v1/reference/villages",
} as const;

export const STORAGE_KEYS = {
  accessToken: "unsia_access_token",
  refreshToken: "unsia_refresh_token",
  user: "unsia_user",
} as const;
