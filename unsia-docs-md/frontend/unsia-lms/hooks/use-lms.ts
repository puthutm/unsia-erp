"use client";

import { useState, useCallback } from "react";

const LMS_API_BASE_URL = "http://localhost:8006/api/v1/lms";

export interface Course {
  id: string;
  code: string;
  name: string;
  description: string;
  lecturerId: string;
  lecturerName: string;
  semester: number;
  academicYearId: string;
  isActive: boolean;
}

export interface Session {
  id: string;
  courseId: string;
  title: string;
  description: string;
  scheduledAt: string;
  duration: number;
  status: "upcoming" | "ongoing" | "completed";
  meetingLink?: string;
  materials: Material[];
  assignments: Assignment[];
}

export interface Material {
  id: string;
  sessionId: string;
  title: string;
  type: "pdf" | "video" | "link" | "document";
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
  submissions: number;
}

export interface Discussion {
  id: string;
  sessionId: string;
  userId: string;
  userName: string;
  content: string;
  createdAt: string;
  replies: DiscussionReply[];
}

export interface DiscussionReply {
  id: string;
  discussionId: string;
  userId: string;
  userName: string;
  content: string;
  createdAt: string;
}

export function useLMS() {
  const [courses, setCourses] = useState<Course[]>([]);
  const [sessions, setSessions] = useState<Session[]>([]);
  const [materials, setMaterials] = useState<Material[]>([]);
  const [assignments, setAssignments] = useState<Assignment[]>([]);
  const [discussions, setDiscussions] = useState<Discussion[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const getToken = () => localStorage.getItem("unsia_access_token");

  // Course hooks
  const fetchCourses = useCallback(async () => {
    setIsLoading(true);
    setError(null);
    try {
      const token = getToken();
      if (!token) throw new Error("Not authenticated");

      const response = await fetch(`${LMS_API_BASE_URL}/courses`, {
        headers: { Authorization: `Bearer ${token}` },
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

  const createCourse = useCallback(async (courseData: Partial<Course>) => {
    setIsLoading(true);
    setError(null);
    try {
      const token = getToken();
      if (!token) throw new Error("Not authenticated");

      const response = await fetch(`${LMS_API_BASE_URL}/courses`, {
        method: "POST",
        headers: { Authorization: `Bearer ${token}`, "Content-Type": "application/json" },
        body: JSON.stringify(courseData),
      });
      if (!response.ok) throw new Error("Failed to create course");
      await fetchCourses();
      return true;
    } catch (err) {
      setError(err instanceof Error ? err.message : "Error");
      return false;
    } finally {
      setIsLoading(false);
    }
  }, [fetchCourses]);

  // Session hooks
  const fetchSessions = useCallback(async (courseId?: string) => {
    setIsLoading(true);
    setError(null);
    try {
      const token = getToken();
      if (!token) throw new Error("Not authenticated");

      const queryParams = courseId ? `?course_id=${courseId}` : "";
      const response = await fetch(`${LMS_API_BASE_URL}/sessions${queryParams}`, {
        headers: { Authorization: `Bearer ${token}` },
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

  const createSession = useCallback(async (sessionData: Partial<Session>) => {
    setIsLoading(true);
    setError(null);
    try {
      const token = getToken();
      if (!token) throw new Error("Not authenticated");

      const response = await fetch(`${LMS_API_BASE_URL}/sessions`, {
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

  // Material hooks
  const fetchMaterials = useCallback(async (sessionId: string) => {
    setIsLoading(true);
    setError(null);
    try {
      const token = getToken();
      if (!token) throw new Error("Not authenticated");

      const response = await fetch(`${LMS_API_BASE_URL}/sessions/${sessionId}/materials`, {
        headers: { Authorization: `Bearer ${token}` },
      });
      if (!response.ok) throw new Error("Failed to fetch materials");
      const data = await response.json();
      setMaterials(data.data || []);
      return data.data || [];
    } catch (err) {
      setError(err instanceof Error ? err.message : "Error");
      return [];
    } finally {
      setIsLoading(false);
    }
  }, []);

  const uploadMaterial = useCallback(async (sessionId: string, materialData: Partial<Material>) => {
    setIsLoading(true);
    setError(null);
    try {
      const token = getToken();
      if (!token) throw new Error("Not authenticated");

      const response = await fetch(`${LMS_API_BASE_URL}/sessions/${sessionId}/materials`, {
        method: "POST",
        headers: { Authorization: `Bearer ${token}`, "Content-Type": "application/json" },
        body: JSON.stringify(materialData),
      });
      if (!response.ok) throw new Error("Failed to upload material");
      await fetchMaterials(sessionId);
      return true;
    } catch (err) {
      setError(err instanceof Error ? err.message : "Error");
      return false;
    } finally {
      setIsLoading(false);
    }
  }, [fetchMaterials]);

  // Assignment hooks
  const fetchAssignments = useCallback(async (sessionId: string) => {
    setIsLoading(true);
    setError(null);
    try {
      const token = getToken();
      if (!token) throw new Error("Not authenticated");

      const response = await fetch(`${LMS_API_BASE_URL}/sessions/${sessionId}/assignments`, {
        headers: { Authorization: `Bearer ${token}` },
      });
      if (!response.ok) throw new Error("Failed to fetch assignments");
      const data = await response.json();
      setAssignments(data.data || []);
      return data.data || [];
    } catch (err) {
      setError(err instanceof Error ? err.message : "Error");
      return [];
    } finally {
      setIsLoading(false);
    }
  }, []);

  const createAssignment = useCallback(async (sessionId: string, assignmentData: Partial<Assignment>) => {
    setIsLoading(true);
    setError(null);
    try {
      const token = getToken();
      if (!token) throw new Error("Not authenticated");

      const response = await fetch(`${LMS_API_BASE_URL}/sessions/${sessionId}/assignments`, {
        method: "POST",
        headers: { Authorization: `Bearer ${token}`, "Content-Type": "application/json" },
        body: JSON.stringify(assignmentData),
      });
      if (!response.ok) throw new Error("Failed to create assignment");
      await fetchAssignments(sessionId);
      return true;
    } catch (err) {
      setError(err instanceof Error ? err.message : "Error");
      return false;
    } finally {
      setIsLoading(false);
    }
  }, [fetchAssignments]);

  // Discussion hooks
  const fetchDiscussions = useCallback(async (sessionId: string) => {
    setIsLoading(true);
    setError(null);
    try {
      const token = getToken();
      if (!token) throw new Error("Not authenticated");

      const response = await fetch(`${LMS_API_BASE_URL}/sessions/${sessionId}/discussions`, {
        headers: { Authorization: `Bearer ${token}` },
      });
      if (!response.ok) throw new Error("Failed to fetch discussions");
      const data = await response.json();
      setDiscussions(data.data || []);
      return data.data || [];
    } catch (err) {
      setError(err instanceof Error ? err.message : "Error");
      return [];
    } finally {
      setIsLoading(false);
    }
  }, []);

  const createDiscussion = useCallback(async (sessionId: string, content: string) => {
    setIsLoading(true);
    setError(null);
    try {
      const token = getToken();
      if (!token) throw new Error("Not authenticated");

      const response = await fetch(`${LMS_API_BASE_URL}/sessions/${sessionId}/discussions`, {
        method: "POST",
        headers: { Authorization: `Bearer ${token}`, "Content-Type": "application/json" },
        body: JSON.stringify({ content }),
      });
      if (!response.ok) throw new Error("Failed to create discussion");
      await fetchDiscussions(sessionId);
      return true;
    } catch (err) {
      setError(err instanceof Error ? err.message : "Error");
      return false;
    } finally {
      setIsLoading(false);
    }
  }, [fetchDiscussions]);

  const replyToDiscussion = useCallback(async (discussionId: string, content: string) => {
    setIsLoading(true);
    setError(null);
    try {
      const token = getToken();
      if (!token) throw new Error("Not authenticated");

      const response = await fetch(`${LMS_API_BASE_URL}/discussions/${discussionId}/replies`, {
        method: "POST",
        headers: { Authorization: `Bearer ${token}`, "Content-Type": "application/json" },
        body: JSON.stringify({ content }),
      });
      if (!response.ok) throw new Error("Failed to reply to discussion");
      return true;
    } catch (err) {
      setError(err instanceof Error ? err.message : "Error");
      return false;
    } finally {
      setIsLoading(false);
    }
  }, []);

  return {
    courses,
    sessions,
    materials,
    assignments,
    discussions,
    isLoading,
    error,
    fetchCourses,
    createCourse,
    fetchSessions,
    createSession,
    fetchMaterials,
    uploadMaterial,
    fetchAssignments,
    createAssignment,
    fetchDiscussions,
    createDiscussion,
    replyToDiscussion,
  };
}
