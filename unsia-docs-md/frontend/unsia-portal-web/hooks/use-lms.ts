"use client";

import { useState, useCallback } from "react";
import { API_BASE_URLS, LMS_ENDPOINTS, STORAGE_KEYS } from "@/lib/constants";

export interface LmsCourse {
  id: string;
  code: string;
  name: string;
  studyProgramName: string;
  semester: number;
  lecturerName: string;
  isActive: boolean;
}

export interface LmsClass {
  id: string;
  className: string;
  courseId: string;
  courseName: string;
  lecturerName: string;
  schedule: string;
  room: string;
  maxStudents: number;
  enrolledCount: number;
  status?: string;
}

export interface Enrollment {
  id: string;
  studentNim: string;
  studentName: string;
  classId: string;
  className: string;
  courseName: string;
  enrolledAt: string;
  status: string;
}

export interface Session {
  id: string;
  classId: string;
  className: string;
  courseName: string;
  topic: string;
  startTime: string;
  endTime: string;
  status: string;
  attendanceCount: number;
}

export interface Material {
  id: string;
  sessionId: string;
  title: string;
  type: string;
  url: string;
  description: string;
  uploadedAt: string;
}

export interface Assignment {
  id: string;
  sessionId: string;
  title: string;
  description: string;
  dueDate: string;
  maxScore: number;
  isActive: boolean;
}

export function useLms() {
  const [courses, setCourses] = useState<LmsCourse[]>([]);
  const [classes, setClasses] = useState<LmsClass[]>([]);
  const [enrollments, setEnrollments] = useState<Enrollment[]>([]);
  const [sessions, setSessions] = useState<Session[]>([]);
  const [materials, setMaterials] = useState<Material[]>([]);
  const [assignments, setAssignments] = useState<Assignment[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const getToken = () => localStorage.getItem(STORAGE_KEYS.accessToken);

  const fetchCourses = useCallback(async (params?: { semester?: number; isActive?: boolean }) => {
    setIsLoading(true);
    setError(null);
    try {
      const token = getToken();
      if (!token) throw new Error("Not authenticated");

      const queryParams = new URLSearchParams();
      if (params?.semester) queryParams.set("semester", params.semester.toString());
      if (params?.isActive !== undefined) queryParams.set("is_active", params.isActive.toString());

      const url = `${API_BASE_URLS.lms}${LMS_ENDPOINTS.courses}${queryParams.toString() ? `?${queryParams}` : ""}`;
      const response = await fetch(url, {
        headers: { Authorization: `Bearer ${token}`, "Content-Type": "application/json" },
      });

      if (!response.ok) throw new Error("Failed to fetch courses");
      const data = await response.json();
      setCourses(data.data || []);
      return data.data || [];
    } catch (err) {
      setError(err instanceof Error ? err.message : "Error");
      return [];
    } finally {
      setIsLoading(false);
    }
  }, []);

  const fetchClasses = useCallback(async (params?: { courseId?: string }) => {
    setIsLoading(true);
    setError(null);
    try {
      const token = getToken();
      if (!token) throw new Error("Not authenticated");

      const queryParams = new URLSearchParams();
      if (params?.courseId) queryParams.set("course_id", params.courseId);

      const url = `${API_BASE_URLS.lms}${LMS_ENDPOINTS.classes}${queryParams.toString() ? `?${queryParams}` : ""}`;
      const response = await fetch(url, {
        headers: { Authorization: `Bearer ${token}`, "Content-Type": "application/json" },
      });

      if (!response.ok) throw new Error("Failed to fetch classes");
      const data = await response.json();
      setClasses(data.data || []);
      return data.data || [];
    } catch (err) {
      setError(err instanceof Error ? err.message : "Error");
      return [];
    } finally {
      setIsLoading(false);
    }
  }, []);

  const fetchEnrollments = useCallback(async (params?: { classId?: string; status?: string }) => {
    setIsLoading(true);
    setError(null);
    try {
      const token = getToken();
      if (!token) throw new Error("Not authenticated");

      const queryParams = new URLSearchParams();
      if (params?.classId) queryParams.set("class_id", params.classId);
      if (params?.status) queryParams.set("status", params.status);

      const url = `${API_BASE_URLS.lms}${LMS_ENDPOINTS.enrollments}${queryParams.toString() ? `?${queryParams}` : ""}`;
      const response = await fetch(url, {
        headers: { Authorization: `Bearer ${token}`, "Content-Type": "application/json" },
      });

      if (!response.ok) throw new Error("Failed to fetch enrollments");
      const data = await response.json();
      setEnrollments(data.data || []);
      return data.data || [];
    } catch (err) {
      setError(err instanceof Error ? err.message : "Error");
      return [];
    } finally {
      setIsLoading(false);
    }
  }, []);

  const fetchSessions = useCallback(async (params?: { classId?: string; status?: string }) => {
    setIsLoading(true);
    setError(null);
    try {
      const token = getToken();
      if (!token) throw new Error("Not authenticated");

      const queryParams = new URLSearchParams();
      if (params?.classId) queryParams.set("class_id", params.classId);
      if (params?.status) queryParams.set("status", params.status);

      const url = `${API_BASE_URLS.lms}${LMS_ENDPOINTS.sessions}${queryParams.toString() ? `?${queryParams}` : ""}`;
      const response = await fetch(url, {
        headers: { Authorization: `Bearer ${token}`, "Content-Type": "application/json" },
      });

      if (!response.ok) throw new Error("Failed to fetch sessions");
      const data = await response.json();
      setSessions(data.data || []);
      return data.data || [];
    } catch (err) {
      setError(err instanceof Error ? err.message : "Error");
      return [];
    } finally {
      setIsLoading(false);
    }
  }, []);

  const enrollStudent = useCallback(async (classId: string, studentId: string) => {
    setIsLoading(true);
    try {
      const token = getToken();
      if (!token) throw new Error("Not authenticated");

      const response = await fetch(`${API_BASE_URLS.lms}${LMS_ENDPOINTS.enrollments}`, {
        method: "POST",
        headers: { Authorization: `Bearer ${token}`, "Content-Type": "application/json" },
        body: JSON.stringify({ class_id: classId, student_id: studentId }),
      });

      if (!response.ok) throw new Error("Failed to enroll student");
      await fetchEnrollments();
      return true;
    } catch (err) {
      setError(err instanceof Error ? err.message : "Error");
      return false;
    } finally {
      setIsLoading(false);
    }
  }, [fetchEnrollments]);

const createSession = useCallback(async (sessionData: Partial<Session>) => {
    setIsLoading(true);
    try {
      const token = getToken();
      if (!token) throw new Error("Not authenticated");

      const response = await fetch(`${API_BASE_URLS.lms}${LMS_ENDPOINTS.sessions}`, {
        method: "POST",
        headers: { Authorization: `Bearer ${token}`, "Content-Type": "application/json" },
        body: JSON.stringify(sessionData),
      });

      if (!response.ok) throw new Error("Failed to create session");
      await fetchSessions();
      return true;
    } catch (err) {
      setError(err instanceof Error ? err.message : "Error");
      return false;
    } finally {
      setIsLoading(false);
    }
  }, [fetchSessions]);

  const createClass = useCallback(async (classData: {
    academic_class_id: string;
    course_id: string;
    lecturer_id?: string;
    class_code: string;
    semester: string;
    academic_year: string;
    max_students: number;
  }) => {
    setIsLoading(true);
    try {
      const token = getToken();
      if (!token) throw new Error("Not authenticated");

      const response = await fetch(`${API_BASE_URLS.lms}${LMS_ENDPOINTS.classes}`, {
        method: "POST",
        headers: { Authorization: `Bearer ${token}`, "Content-Type": "application/json" },
        body: JSON.stringify(classData),
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.message || "Failed to create class");
      }
      await fetchClasses();
      return true;
    } catch (err) {
      setError(err instanceof Error ? err.message : "Error");
      return false;
    } finally {
      setIsLoading(false);
    }
  }, [fetchClasses]);

  return {
    courses,
    classes,
    enrollments,
    sessions,
    materials,
    assignments,
    isLoading,
    error,
    fetchCourses,
    fetchClasses,
    fetchEnrollments,
    fetchSessions,
    enrollStudent,
    createSession,
    createClass,
  };
}
