// API Configuration for UNSIA Backend Services

// Service URLs - these should match the backend service ports
export const API_BASE_URLS = {
  auth: process.env.NEXT_PUBLIC_API_CORE_URL || "http://localhost:8001",
  reference: process.env.NEXT_PUBLIC_API_REFERENCE_URL || "http://localhost:8002",
  finance: process.env.NEXT_PUBLIC_API_FINANCE_URL || "http://localhost:8005",
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
