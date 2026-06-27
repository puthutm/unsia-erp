"use client";

import { useState, useCallback } from "react";
import { API_BASE_URLS, HRIS_ENDPOINTS, STORAGE_KEYS } from "@/lib/constants";

export interface Employee {
  id: string;
  employeeNumber: string;
  personId: string;
  name: string;
  email: string;
  role: string;
  status: string;
  joinDate: string;
}

export interface Attendance {
  id: string;
  employeeId: string;
  date: string;
  clockInTime: string;
  clockOutTime?: string;
  status: string;
  notes?: string;
}

export interface LeaveRequest {
  id: string;
  employeeId: string;
  leaveType: string;
  startDate: string;
  endDate: string;
  reason: string;
  status: string;
  createdAt: string;
}

export function useHRIS() {
  const [employees, setEmployees] = useState<Employee[]>([]);
  const [attendances, setAttendances] = useState<Attendance[]>([]);
  const [leaveRequests, setLeaveRequests] = useState<LeaveRequest[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const getToken = () => localStorage.getItem(STORAGE_KEYS.accessToken);

  const fetchEmployees = useCallback(async (): Promise<Employee[]> => {
    setIsLoading(true);
    setError(null);
    try {
      const token = getToken();
      if (!token) throw new Error("Not authenticated");

      const url = `${API_BASE_URLS.hris}${HRIS_ENDPOINTS.employees}`;
      const response = await fetch(url, {
        headers: { Authorization: `Bearer ${token}`, "Content-Type": "application/json" },
      });

      if (!response.ok) throw new Error("Failed to fetch employees");
      const data = await response.json();
      const records: Employee[] = data.data || [];
      setEmployees(records);
      return records;
    } catch (err) {
      setError(err instanceof Error ? err.message : "Error");
      return [];
    } finally {
      setIsLoading(false);
    }
  }, []);

  const fetchAttendances = useCallback(async (): Promise<Attendance[]> => {
    setIsLoading(true);
    setError(null);
    try {
      const token = getToken();
      if (!token) throw new Error("Not authenticated");

      const url = `${API_BASE_URLS.hris}${HRIS_ENDPOINTS.attendances}`;
      const response = await fetch(url, {
        headers: { Authorization: `Bearer ${token}`, "Content-Type": "application/json" },
      });

      if (!response.ok) throw new Error("Failed to fetch attendances");
      const data = await response.json();
      const records: Attendance[] = data.data || [];
      setAttendances(records);
      return records;
    } catch (err) {
      setError(err instanceof Error ? err.message : "Error");
      return [];
    } finally {
      setIsLoading(false);
    }
  }, []);

  const clockIn = useCallback(async (notes?: string) => {
    setIsLoading(true);
    setError(null);
    try {
      const token = getToken();
      if (!token) throw new Error("Not authenticated");

      const response = await fetch(`${API_BASE_URLS.hris}${HRIS_ENDPOINTS.attendances}/clock-in`, {
        method: "POST",
        headers: { Authorization: `Bearer ${token}`, "Content-Type": "application/json" },
        body: JSON.stringify({ notes }),
      });

      if (!response.ok) throw new Error("Failed to clock in");
      await fetchAttendances();
      return true;
    } catch (err) {
      setError(err instanceof Error ? err.message : "Error");
      return false;
    } finally {
      setIsLoading(false);
    }
  }, [fetchAttendances]);

  const clockOut = useCallback(async (notes?: string) => {
    setIsLoading(true);
    setError(null);
    try {
      const token = getToken();
      if (!token) throw new Error("Not authenticated");

      const response = await fetch(`${API_BASE_URLS.hris}${HRIS_ENDPOINTS.attendances}/clock-out`, {
        method: "POST",
        headers: { Authorization: `Bearer ${token}`, "Content-Type": "application/json" },
        body: JSON.stringify({ notes }),
      });

      if (!response.ok) throw new Error("Failed to clock out");
      await fetchAttendances();
      return true;
    } catch (err) {
      setError(err instanceof Error ? err.message : "Error");
      return false;
    } finally {
      setIsLoading(false);
    }
  }, [fetchAttendances]);

  const fetchLeaveRequests = useCallback(async (): Promise<LeaveRequest[]> => {
    setIsLoading(true);
    setError(null);
    try {
      const token = getToken();
      if (!token) throw new Error("Not authenticated");

      const url = `${API_BASE_URLS.hris}${HRIS_ENDPOINTS.leaveRequests}`;
      const response = await fetch(url, {
        headers: { Authorization: `Bearer ${token}`, "Content-Type": "application/json" },
      });

      if (!response.ok) throw new Error("Failed to fetch leave requests");
      const data = await response.json();
      const records: LeaveRequest[] = data.data || [];
      setLeaveRequests(records);
      return records;
    } catch (err) {
      setError(err instanceof Error ? err.message : "Error");
      return [];
    } finally {
      setIsLoading(false);
    }
  }, []);

  const createLeaveRequest = useCallback(async (leaveData: Partial<LeaveRequest>) => {
    setIsLoading(true);
    setError(null);
    try {
      const token = getToken();
      if (!token) throw new Error("Not authenticated");

      const response = await fetch(`${API_BASE_URLS.hris}${HRIS_ENDPOINTS.leaveRequests}`, {
        method: "POST",
        headers: { Authorization: `Bearer ${token}`, "Content-Type": "application/json" },
        body: JSON.stringify(leaveData),
      });

      if (!response.ok) throw new Error("Failed to create leave request");
      await fetchLeaveRequests();
      return true;
    } catch (err) {
      setError(err instanceof Error ? err.message : "Error");
      return false;
    } finally {
      setIsLoading(false);
    }
  }, [fetchLeaveRequests]);

  return {
    employees,
    attendances,
    leaveRequests,
    isLoading,
    error,
    fetchEmployees,
    fetchAttendances,
    clockIn,
    clockOut,
    fetchLeaveRequests,
    createLeaveRequest,
  };
}
