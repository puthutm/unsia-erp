"use client";

import { useState, useCallback } from "react";

export interface Employee {
  id: string;
  employeeCode: string;
  personId: string;
  employmentType: string;
  joinDate: string;
  status: string;
  organizationUnitId?: string;
}

export interface Attendance {
  id: string;
  employeeId: string;
  checkInAt?: string;
  checkOutAt?: string;
  status: string;
  workDate: string;
}

export interface LeaveRequest {
  id: string;
  employeeId: string;
  leaveType: string;
  startDate: string;
  endDate: string;
  reason: string;
  status: string;
}

const API_BASE_URL = process.env.NEXT_PUBLIC_HRIS_API || "http://localhost:8008";
const STORAGE_KEY = "unsia_access_token";

export function useHRIS() {
  const [employees, setEmployees] = useState<Employee[]>([]);
  const [attendance, setAttendance] = useState<Attendance[]>([]);
  const [leaveRequests, setLeaveRequests] = useState<LeaveRequest[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const getToken = () => localStorage.getItem(STORAGE_KEY);

  const request = async (url: string, options: RequestInit = {}) => {
    const token = getToken();
    if (!token) throw new Error("Not authenticated");

    const response = await fetch(url, {
      ...options,
      headers: {
        Authorization: `Bearer ${token}`,
        "Content-Type": "application/json",
        ...(options.headers || {}),
      },
    });

    if (!response.ok) {
      throw new Error(`Request failed: ${response.status}`);
    }

    return response.json();
  };

  const fetchEmployees = useCallback(async () => {
    setIsLoading(true);
    setError(null);
    try {
      const res = await request(`${API_BASE_URL}/api/v1/hris/employees`);
      const data = res.data || [];
      setEmployees(data);
      return data;
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to fetch employees");
      return [];
    } finally {
      setIsLoading(false);
    }
  }, []);

  const fetchAttendance = useCallback(async () => {
    setIsLoading(true);
    setError(null);
    try {
      const res = await request(`${API_BASE_URL}/api/v1/hris/attendance`);
      const data = res.data || [];
      setAttendance(data);
      return data;
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to fetch attendance");
      return [];
    } finally {
      setIsLoading(false);
    }
  }, []);

  const fetchLeaveRequests = useCallback(async () => {
    setIsLoading(true);
    setError(null);
    try {
      const res = await request(`${API_BASE_URL}/api/v1/hris/leave`);
      const data = res.data || [];
      setLeaveRequests(data);
      return data;
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to fetch leave requests");
      return [];
    } finally {
      setIsLoading(false);
    }
  }, []);

  const createLeaveRequest = useCallback(async (payload: Partial<LeaveRequest>) => {
    setIsLoading(true);
    setError(null);
    try {
      await request(`${API_BASE_URL}/api/v1/hris/leave`, {
        method: "POST",
        body: JSON.stringify(payload),
      });
      await fetchLeaveRequests();
      return true;
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to create leave request");
      return false;
    } finally {
      setIsLoading(false);
    }
  }, [fetchLeaveRequests]);

  return {
    employees,
    attendance,
    leaveRequests,
    isLoading,
    error,
    fetchEmployees,
    fetchAttendance,
    fetchLeaveRequests,
    createLeaveRequest,
  };
}
