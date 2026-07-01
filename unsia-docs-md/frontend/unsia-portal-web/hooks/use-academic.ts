"use client";

import { useState, useCallback } from "react";
import { useQueryClient } from "@tanstack/react-query";
import { API_BASE_URLS, ACADEMIC_ENDPOINTS, STORAGE_KEYS } from "@/lib/constants";

export interface Student {
  id: string;
  nim: string;
  name: string;
  email: string;
  phone: string;
  studyProgramName: string;
  entryYear: string;
  status: string;
}

export interface Course {
  id: string;
  code: string;
  name: string;
  sks: number;
  semester: number;
  isActive: boolean;
}

export interface Schedule {
  id: string;
  courseName: string;
  className: string;
  day: string;
  startTime: string;
  endTime: string;
  room: string;
  lecturerName: string;
}

export interface KrsEntry {
  id: string;
  studentId: string;
  academicPeriodId: string;
  status: string;
  items: {
    id: string;
    classId: string;
    status: string;
  }[];
}

export interface StudentGrade {
  id: string;
  krsItemId: string;
  numericGrade: number;
  letterGrade: string;
  gradePoint: number;
  source: string;
  courseName?: string;
  courseCode?: string;
  sks?: number;
}

export function useAcademic() {
  const queryClient = useQueryClient();
  const [students, setStudents] = useState<Student[]>([]);
  const [courses, setCourses] = useState<Course[]>([]);
  const [schedules, setSchedules] = useState<Schedule[]>([]);
  const [krs, setKrs] = useState<KrsEntry[]>([]);
  const [grades, setGrades] = useState<StudentGrade[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const getToken = () => localStorage.getItem(STORAGE_KEYS.accessToken);

  const fetchStudents = useCallback(async (params?: { studyProgramId?: string; status?: string; search?: string }) => {
    setIsLoading(true);
    setError(null);
    try {
      const queryParams = new URLSearchParams();
      if (params?.studyProgramId) queryParams.set("study_program_id", params.studyProgramId);
      if (params?.status) queryParams.set("status", params.status);
      if (params?.search) queryParams.set("search", params.search);

      const records = await queryClient.fetchQuery({
        queryKey: ["academic", "students", params],
        queryFn: async () => {
          const token = getToken();
          if (!token) throw new Error("Not authenticated");

          const url = `${API_BASE_URLS.academic}${ACADEMIC_ENDPOINTS.students}${queryParams.toString() ? `?${queryParams}` : ""}`;
          const response = await fetch(url, {
            headers: { Authorization: `Bearer ${token}`, "Content-Type": "application/json" },
          });

          if (!response.ok) throw new Error("Failed to fetch students");
          const data = await response.json();
          return data.data || [];
        }
      });
      setStudents(records);
      return records;
    } catch (err) {
      setError(err instanceof Error ? err.message : "Error");
      return [];
    } finally {
      setIsLoading(false);
    }
  }, [queryClient]);

  const generateStudentFromApplicant = useCallback(async (applicantId: string, studyProgramId: string) => {
    setIsLoading(true);
    setError(null);
    try {
      const token = getToken();
      if (!token) throw new Error("Not authenticated");

      const response = await fetch(`${API_BASE_URLS.academic}${ACADEMIC_ENDPOINTS.students}/generate-from-applicant`, {
        method: "POST",
        headers: { Authorization: `Bearer ${token}`, "Content-Type": "application/json" },
        body: JSON.stringify({ applicant_id: applicantId, study_program_id: studyProgramId }),
      });

      if (!response.ok) throw new Error("Failed to generate student from applicant");
      
      await queryClient.invalidateQueries({ queryKey: ["academic", "students"] });
      return true;
    } catch (err) {
      setError(err instanceof Error ? err.message : "Error");
      return false;
    } finally {
      setIsLoading(false);
    }
  }, [queryClient]);

  const fetchCourses = useCallback(async (params?: { semester?: number; isActive?: boolean }) => {
    setIsLoading(true);
    setError(null);
    try {
      const queryParams = new URLSearchParams();
      if (params?.semester) queryParams.set("semester", params.semester.toString());
      if (params?.isActive !== undefined) queryParams.set("is_active", params.isActive.toString());

      const records = await queryClient.fetchQuery({
        queryKey: ["academic", "courses", params],
        queryFn: async () => {
          const token = getToken();
          if (!token) throw new Error("Not authenticated");

          const url = `${API_BASE_URLS.academic}${ACADEMIC_ENDPOINTS.courses}${queryParams.toString() ? `?${queryParams}` : ""}`;
          const response = await fetch(url, {
            headers: { Authorization: `Bearer ${token}`, "Content-Type": "application/json" },
          });

          if (!response.ok) throw new Error("Failed to fetch courses");
          const data = await response.json();
          return data.data || [];
        }
      });
      setCourses(records);
      return records;
    } catch (err) {
      setError(err instanceof Error ? err.message : "Error");
      return [];
    } finally {
      setIsLoading(false);
    }
  }, [queryClient]);

  const fetchSchedules = useCallback(async (params?: { day?: string; roomId?: string }) => {
    setIsLoading(true);
    setError(null);
    try {
      const queryParams = new URLSearchParams();
      if (params?.day) queryParams.set("day", params.day);
      if (params?.roomId) queryParams.set("room_id", params.roomId);

      const records = await queryClient.fetchQuery({
        queryKey: ["academic", "schedules", params],
        queryFn: async () => {
          const token = getToken();
          if (!token) throw new Error("Not authenticated");

          const url = `${API_BASE_URLS.academic}${ACADEMIC_ENDPOINTS.schedules}${queryParams.toString() ? `?${queryParams}` : ""}`;
          const response = await fetch(url, {
            headers: { Authorization: `Bearer ${token}`, "Content-Type": "application/json" },
          });

          if (!response.ok) throw new Error("Failed to fetch schedules");
          const data = await response.json();
          return data.data || [];
        }
      });
      setSchedules(records);
      return records;
    } catch (err) {
      setError(err instanceof Error ? err.message : "Error");
      return [];
    } finally {
      setIsLoading(false);
    }
  }, [queryClient]);

  const fetchKrs = useCallback(async (studentId?: string, periodId?: string) => {
    setIsLoading(true);
    setError(null);
    try {
      const queryParams = new URLSearchParams();
      if (studentId) queryParams.set("student_id", studentId);
      if (periodId) queryParams.set("academic_period_id", periodId);

      const records = await queryClient.fetchQuery({
        queryKey: ["academic", "krs", { studentId, periodId }],
        queryFn: async () => {
          const token = getToken();
          if (!token) throw new Error("Not authenticated");

          const url = `${API_BASE_URLS.academic}${ACADEMIC_ENDPOINTS.krs}${queryParams.toString() ? `?${queryParams}` : ""}`;
          const response = await fetch(url, {
            headers: { Authorization: `Bearer ${token}`, "Content-Type": "application/json" },
          });

          if (!response.ok) throw new Error("Failed to fetch KRS");
          const data = await response.json();
          return data.data || [];
        }
      });
      setKrs(records);
      return records;
    } catch (err) {
      setError(err instanceof Error ? err.message : "Error");
      return [];
    } finally {
      setIsLoading(false);
    }
  }, [queryClient]);

  const createKrsDraft = useCallback(async (krsData: { student_id: string; academic_period_id: string; items: { class_id: string }[] }) => {
    setIsLoading(true);
    setError(null);
    try {
      const token = getToken();
      if (!token) throw new Error("Not authenticated");

      const response = await fetch(`${API_BASE_URLS.academic}${ACADEMIC_ENDPOINTS.krs}`, {
        method: "POST",
        headers: { Authorization: `Bearer ${token}`, "Content-Type": "application/json" },
        body: JSON.stringify(krsData),
      });

      if (!response.ok) throw new Error("Failed to create KRS draft");
      
      await queryClient.invalidateQueries({ queryKey: ["academic", "krs"] });
      return true;
    } catch (err) {
      setError(err instanceof Error ? err.message : "Error");
      return false;
    } finally {
      setIsLoading(false);
    }
  }, [queryClient]);

  const submitKrs = useCallback(async (krsId: string) => {
    setIsLoading(true);
    setError(null);
    try {
      const token = getToken();
      if (!token) throw new Error("Not authenticated");

      const response = await fetch(`${API_BASE_URLS.academic}${ACADEMIC_ENDPOINTS.krs}/${krsId}/submit`, {
        method: "POST",
        headers: { Authorization: `Bearer ${token}`, "Content-Type": "application/json" },
      });

      if (!response.ok) throw new Error("Failed to submit KRS");
      
      await queryClient.invalidateQueries({ queryKey: ["academic", "krs"] });
      return true;
    } catch (err) {
      setError(err instanceof Error ? err.message : "Error");
      return false;
    } finally {
      setIsLoading(false);
    }
  }, [queryClient]);

  const approveKrs = useCallback(async (krsId: string) => {
    setIsLoading(true);
    setError(null);
    try {
      const token = getToken();
      if (!token) throw new Error("Not authenticated");

      const response = await fetch(`${API_BASE_URLS.academic}${ACADEMIC_ENDPOINTS.krs}/${krsId}/approve`, {
        method: "POST",
        headers: { Authorization: `Bearer ${token}`, "Content-Type": "application/json" },
      });

      if (!response.ok) throw new Error("Failed to approve KRS");
      
      await queryClient.invalidateQueries({ queryKey: ["academic", "krs"] });
      return true;
    } catch (err) {
      setError(err instanceof Error ? err.message : "Error");
      return false;
    } finally {
      setIsLoading(false);
    }
  }, [queryClient]);

  const fetchGrades = useCallback(async (studentId: string) => {
    setIsLoading(true);
    setError(null);
    try {
      const records = await queryClient.fetchQuery({
        queryKey: ["academic", "grades", studentId],
        queryFn: async () => {
          const token = getToken();
          if (!token) throw new Error("Not authenticated");

          const url = `${API_BASE_URLS.academic}${ACADEMIC_ENDPOINTS.grades}/student/${studentId}`;
          const response = await fetch(url, {
            headers: { Authorization: `Bearer ${token}`, "Content-Type": "application/json" },
          });

          if (!response.ok) throw new Error("Failed to fetch grades");
          const data = await response.json();
          return data.data || [];
        }
      });
      setGrades(records);
      return records;
    } catch (err) {
      setError(err instanceof Error ? err.message : "Error");
      return [];
    } finally {
      setIsLoading(false);
    }
  }, [queryClient]);

  const enterStudentGrade = useCallback(async (gradeId: string, score: number) => {
    setIsLoading(true);
    setError(null);
    try {
      const token = getToken();
      if (!token) throw new Error("Not authenticated");

      const response = await fetch(`${API_BASE_URLS.academic}${ACADEMIC_ENDPOINTS.grades}/${gradeId}/entries`, {
        method: "POST",
        headers: { Authorization: `Bearer ${token}`, "Content-Type": "application/json" },
        body: JSON.stringify({ numeric_grade: score }),
      });

      if (!response.ok) throw new Error("Failed to enter grade entry");
      
      await queryClient.invalidateQueries({ queryKey: ["academic", "grades"] });
      return true;
    } catch (err) {
      setError(err instanceof Error ? err.message : "Error");
      return false;
    } finally {
      setIsLoading(false);
    }
  }, [queryClient]);

  const finalizeGrade = useCallback(async (gradeId: string, numericGrade: number, letterGrade: string, gradePoint: number) => {
    setIsLoading(true);
    setError(null);
    try {
      const token = getToken();
      if (!token) throw new Error("Not authenticated");

      const response = await fetch(`${API_BASE_URLS.academic}${ACADEMIC_ENDPOINTS.grades}/${gradeId}/finalize`, {
        method: "POST",
        headers: { Authorization: `Bearer ${token}`, "Content-Type": "application/json" },
        body: JSON.stringify({ numeric_grade: numericGrade, letter_grade: letterGrade, grade_point: gradePoint }),
      });

      if (!response.ok) throw new Error("Failed to finalize grade");
      
      await queryClient.invalidateQueries({ queryKey: ["academic", "grades"] });
      return true;
    } catch (err) {
      setError(err instanceof Error ? err.message : "Error");
      return false;
    } finally {
      setIsLoading(false);
    }
  }, [queryClient]);

  const correctGrade = useCallback(async (gradeId: string, numericGrade: number, letterGrade: string, gradePoint: number, reason: string) => {
    setIsLoading(true);
    setError(null);
    try {
      const token = getToken();
      if (!token) throw new Error("Not authenticated");

      const response = await fetch(`${API_BASE_URLS.academic}${ACADEMIC_ENDPOINTS.grades}/${gradeId}/corrections`, {
        method: "POST",
        headers: { Authorization: `Bearer ${token}`, "Content-Type": "application/json" },
        body: JSON.stringify({ numeric_grade: numericGrade, letter_grade: letterGrade, grade_point: gradePoint, reason }),
      });

      if (!response.ok) throw new Error("Failed to correct grade");
      
      await queryClient.invalidateQueries({ queryKey: ["academic", "grades"] });
      return true;
    } catch (err) {
      setError(err instanceof Error ? err.message : "Error");
      return false;
    } finally {
      setIsLoading(false);
    }
  }, [queryClient]);

  const checkGraduationEligibility = useCallback(async (studentId: string) => {
    setIsLoading(true);
    setError(null);
    try {
      const data = await queryClient.fetchQuery({
        queryKey: ["academic", "graduation", "eligibility", studentId],
        queryFn: async () => {
          const token = getToken();
          if (!token) throw new Error("Not authenticated");

          const response = await fetch(`${API_BASE_URLS.academic}${ACADEMIC_ENDPOINTS.graduation}/eligibility/${studentId}`, {
            headers: { Authorization: `Bearer ${token}`, "Content-Type": "application/json" },
          });

          if (!response.ok) throw new Error("Failed to check graduation eligibility");
          const res = await response.json();
          return res.data;
        }
      });
      return data;
    } catch (err) {
      setError(err instanceof Error ? err.message : "Error");
      return null;
    } finally {
      setIsLoading(false);
    }
  }, [queryClient]);

  const applyGraduation = useCallback(async (studentId: string) => {
    setIsLoading(true);
    setError(null);
    try {
      const token = getToken();
      if (!token) throw new Error("Not authenticated");

      const response = await fetch(`${API_BASE_URLS.academic}${ACADEMIC_ENDPOINTS.graduation}/apply`, {
        method: "POST",
        headers: { Authorization: `Bearer ${token}`, "Content-Type": "application/json" },
        body: JSON.stringify({ student_id: studentId }),
      });

      if (!response.ok) throw new Error("Failed to apply for graduation");
      
      await queryClient.invalidateQueries({ queryKey: ["academic", "graduation"] });
      return true;
    } catch (err) {
      setError(err instanceof Error ? err.message : "Error");
      return false;
    } finally {
      setIsLoading(false);
    }
  }, [queryClient]);

  const updateStudentStatus = useCallback(async (studentId: string, status: string) => {
    setIsLoading(true);
    try {
      const token = getToken();
      if (!token) throw new Error("Not authenticated");

      const response = await fetch(`${API_BASE_URLS.academic}${ACADEMIC_ENDPOINTS.students}/${studentId}`, {
        method: "PUT",
        headers: { Authorization: `Bearer ${token}`, "Content-Type": "application/json" },
        body: JSON.stringify({ status }),
      });

      if (!response.ok) throw new Error("Failed to update student");
      
      await queryClient.invalidateQueries({ queryKey: ["academic", "students"] });
      await fetchStudents();
      return true;
    } catch (err) {
      setError(err instanceof Error ? err.message : "Error");
      return false;
    } finally {
      setIsLoading(false);
    }
  }, [queryClient, fetchStudents]);

  const createSchedule = useCallback(async (scheduleData: {
    class_id: string;
    day_of_week: number;
    start_time: string;
    end_time: string;
    room_id?: string;
    building_id?: string;
    schedule_type?: string;
    is_online?: boolean;
    meeting_link?: string;
  }) => {
    setIsLoading(true);
    try {
      const token = getToken();
      if (!token) throw new Error("Not authenticated");

      const response = await fetch(`${API_BASE_URLS.academic}${ACADEMIC_ENDPOINTS.schedules}`, {
        method: "POST",
        headers: { Authorization: `Bearer ${token}`, "Content-Type": "application/json" },
        body: JSON.stringify(scheduleData),
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.message || "Failed to create schedule");
      }
      
      await queryClient.invalidateQueries({ queryKey: ["academic", "schedules"] });
      await fetchSchedules();
      return true;
    } catch (err) {
      setError(err instanceof Error ? err.message : "Error");
      return false;
    } finally {
      setIsLoading(false);
    }
  }, [queryClient, fetchSchedules]);

  return {
    students,
    courses,
    schedules,
    krs,
    grades,
    isLoading,
    error,
    fetchStudents,
    generateStudentFromApplicant,
    fetchCourses,
    fetchSchedules,
    fetchKrs,
    createKrsDraft,
    submitKrs,
    approveKrs,
    fetchGrades,
    enterStudentGrade,
    finalizeGrade,
    correctGrade,
    checkGraduationEligibility,
    applyGraduation,
    updateStudentStatus,
    createSchedule,
  };
}
