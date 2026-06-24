// Base API client for connecting to backend microservices

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

export interface ApiResponse<T> {
  data: T;
  success: boolean;
  message?: string;
}

export interface PaginatedResponse<T> {
  data: T[];
  total: number;
  page: number;
  limit: number;
  totalPages: number;
}

export interface ApiError {
  code: string;
  message: string;
  details?: Record<string, string>;
}

// Base fetch wrapper with error handling
async function fetchApi<T>(
  endpoint: string,
  options: RequestInit = {}
): Promise<T> {
  const url = `${API_BASE_URL}${endpoint}`;
  
  const response = await fetch(url, {
    ...options,
    headers: {
      "Content-Type": "application/json",
      ...options.headers,
    },
  });

  if (!response.ok) {
    const error: ApiError = await response.json();
    throw new Error(error.message || "An error occurred");
  }

  return response.json();
}

// Auth API
export const authApi = {
  login: (email: string, password: string) =>
    fetchApi<{ token: string; user: any }>("/api/v1/auth/login", {
      method: "POST",
      body: JSON.stringify({ email, password }),
    }),
  
  logout: () =>
    fetchApi<void>("/api/v1/auth/logout", {
      method: "POST",
    }),
  
  me: () =>
    fetchApi<any>("/api/v1/auth/me"),
};

// Reference API
export const referenceApi = {
  getProvinces: () =>
    fetchApi<any[]>("/api/v1/reference/provinces"),
  
  getCities: (provinceId: string) =>
    fetchApi<any[]>(`/api/v1/reference/cities?province_id=${provinceId}`),
  
  getDistricts: (cityId: string) =>
    fetchApi<any[]>(`/api/v1/reference/districts?city_id=${cityId}`),
  
  getVillages: (districtId: string) =>
    fetchApi<any[]>(`/api/v1/reference/villages?district_id=${districtId}`),
  
  getStudyPrograms: () =>
    fetchApi<any[]>("/api/v1/reference/study-programs"),
  
  getFaculties: () =>
    fetchApi<any[]>("/api/v1/reference/faculties"),
  
  getDegrees: () =>
    fetchApi<any[]>("/api/v1/reference/degrees"),
  
  getAcademicYears: () =>
    fetchApi<any[]>("/api/v1/reference/academic-years"),
  
  getSemesters: () =>
    fetchApi<any[]>("/api/v1/reference/semesters"),
};

// PMB API
export const pmbApi = {
  // Applicants
  getApplicants: (params?: { page?: number; limit?: number; status?: string }) => {
    const query = new URLSearchParams(params as any).toString();
    return fetchApi<PaginatedResponse<any>>(`/api/v1/applicants?${query}`);
  },
  
  getApplicant: (id: string) =>
    fetchApi<any>(`/api/v1/applicants/${id}`),
  
  createApplicant: (data: any) =>
    fetchApi<any>("/api/v1/applicants", {
      method: "POST",
      body: JSON.stringify(data),
    }),
  
  updateApplicant: (id: string, data: any) =>
    fetchApi<any>(`/api/v1/applicants/${id}`, {
      method: "PUT",
      body: JSON.stringify(data),
    }),
  
  // Waves
  getWaves: () =>
    fetchApi<any[]>("/api/v1/waves"),
  
  createWave: (data: any) =>
    fetchApi<any>("/api/v1/waves", {
      method: "POST",
      body: JSON.stringify(data),
    }),
  
  // Documents
  getDocuments: (applicantId: string) =>
    fetchApi<any[]>(`/api/v1/applicants/${applicantId}/documents`),
  
  uploadDocument: (applicantId: string, data: any) =>
    fetchApi<any>(`/api/v1/applicants/${applicantId}/documents`, {
      method: "POST",
      body: JSON.stringify(data),
    }),
  
  verifyDocument: (documentId: string, status: string) =>
    fetchApi<any>(`/api/v1/documents/${documentId}/verify`, {
      method: "POST",
      body: JSON.stringify({ status }),
    }),
  
  // Payments
  getPayments: (params?: { page?: number; limit?: number }) => {
    const query = new URLSearchParams(params as any).toString();
    return fetchApi<PaginatedResponse<any>>(`/api/v1/payments?${query}`);
  },
  
  verifyPayment: (paymentId: string) =>
    fetchApi<any>(`/api/v1/payments/${paymentId}/verify`, {
      method: "POST",
    }),
  
  // Selection
  getSelectionResults: (waveId: string) =>
    fetchApi<any[]>(`/api/v1/selection/results?wave_id=${waveId}`),
  
  publishSelection: (waveId: string) =>
    fetchApi<any>(`/api/v1/selection/publish`, {
      method: "POST",
      body: JSON.stringify({ wave_id: waveId }),
    }),
};

// Finance API
export const financeApi = {
  getInvoices: (params?: { page?: number; limit?: number; status?: string }) => {
    const query = new URLSearchParams(params as any).toString();
    return fetchApi<PaginatedResponse<any>>(`/api/v1/invoices?${query}`);
  },
  
  createInvoice: (data: any) =>
    fetchApi<any>("/api/v1/invoices", {
      method: "POST",
      body: JSON.stringify(data),
    }),
  
  getPayments: (params?: { page?: number; limit?: number }) => {
    const query = new URLSearchParams(params as any).toString();
    return fetchApi<PaginatedResponse<any>>(`/api/v1/finance/payments?${query}`);
  },
  
  getDashboard: () =>
    fetchApi<any>("/api/v1/finance/dashboard"),
  
  getReports: (type: string, params?: any) => {
    const query = new URLSearchParams(params as any).toString();
    return fetchApi<any>(`/api/v1/finance/reports/${type}?${query}`);
  },
};

// Academic API
export const academicApi = {
  getStudents: (params?: { page?: number; limit?: number }) => {
    const query = new URLSearchParams(params as any).toString();
    return fetchApi<PaginatedResponse<any>>(`/api/v1/students?${query}`);
  },
  
  getStudent: (id: string) =>
    fetchApi<any>(`/api/v1/students/${id}`),
  
  createStudent: (data: any) =>
    fetchApi<any>("/api/v1/students", {
      method: "POST",
      body: JSON.stringify(data),
    }),
  
  getKRS: (studentId: string, academicYearId: string) =>
    fetchApi<any>(`/api/v1/krs?student_id=${studentId}&academic_year_id=${academicYearId}`),
  
  getSchedules: (params?: { day?: string; roomId?: string }) => {
    const query = new URLSearchParams(params as any).toString();
    return fetchApi<any[]>(`/api/v1/schedules?${query}`);
  },
  
  getGrades: (studentId: string) =>
    fetchApi<any[]>(`/api/v1/grades?student_id=${studentId}`),
  
  getTranscripts: (studentId: string) =>
    fetchApi<any>(`/api/v1/transcripts?student_id=${studentId}`),
};

// LMS API
export const lmsApi = {
  getCourses: (params?: { page?: number; limit?: number }) => {
    const query = new URLSearchParams(params as any).toString();
    return fetchApi<PaginatedResponse<any>>(`/api/v1/lms/courses?${query}`);
  },
  
  getCourse: (id: string) =>
    fetchApi<any>(`/api/v1/lms/courses/${id}`),
  
  getSessions: (courseId: string) =>
    fetchApi<any[]>(`/api/v1/lms/sessions?course_id=${courseId}`),
  
  getSession: (id: string) =>
    fetchApi<any>(`/api/v1/lms/sessions/${id}`),
  
  getAssignments: (courseId: string) =>
    fetchApi<any[]>(`/api/v1/lms/assignments?course_id=${courseId}`),
  
  getEnrollments: (courseId: string) =>
    fetchApi<any[]>(`/api/v1/lms/enrollments?course_id=${courseId}`),
  
  enrollStudent: (courseId: string, studentId: string) =>
    fetchApi<any>(`/api/v1/lms/enrollments`, {
      method: "POST",
      body: JSON.stringify({ course_id: courseId, student_id: studentId }),
    }),
};

export default fetchApi;
