"use client";

import { useState, useCallback } from "react";
import { API_BASE_URLS, FINANCE_ENDPOINTS, STORAGE_KEYS } from "@/lib/constants";

export interface Invoice {
  id: string;
  invoiceNumber: string;
  personName: string;
  personType: string;
  amount: number;
  paidAmount: number;
  status: string;
  dueDate: string;
  createdAt: string;
}

export interface Payment {
  id: string;
  paymentNumber: string;
  invoiceNumber: string;
  personName: string;
  amount: number;
  method: string;
  status: string;
  paidAt: string;
}

export interface FinanceStats {
  totalReceivable: number;
  totalPayable: number;
  collected: number;
  pendingPayment: number;
}

export function useFinance() {
  const [invoices, setInvoices] = useState<Invoice[]>([]);
  const [payments, setPayments] = useState<Payment[]>([]);
  const [stats, setStats] = useState<FinanceStats | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const getToken = () => localStorage.getItem(STORAGE_KEYS.accessToken);

  const fetchInvoices = useCallback(async (params?: { status?: string; personId?: string }) => {
    setIsLoading(true);
    setError(null);
    try {
      const token = getToken();
      if (!token) throw new Error("Not authenticated");

      const queryParams = new URLSearchParams();
      if (params?.status) queryParams.set("status", params.status);
      if (params?.personId) queryParams.set("person_id", params.personId);

      const url = `${API_BASE_URLS.finance}${FINANCE_ENDPOINTS.invoices}${queryParams.toString() ? `?${queryParams}` : ""}`;
      const response = await fetch(url, {
        headers: { Authorization: `Bearer ${token}`, "Content-Type": "application/json" },
      });

      if (!response.ok) throw new Error("Failed to fetch invoices");
      const data = await response.json();
      setInvoices(data.data || []);
      return data.data || [];
    } catch (err) {
      setError(err instanceof Error ? err.message : "Error");
      return [];
    } finally {
      setIsLoading(false);
    }
  }, []);

  const fetchPayments = useCallback(async (params?: { status?: string; invoiceId?: string }) => {
    setIsLoading(true);
    setError(null);
    try {
      const token = getToken();
      if (!token) throw new Error("Not authenticated");

      const queryParams = new URLSearchParams();
      if (params?.status) queryParams.set("status", params.status);
      if (params?.invoiceId) queryParams.set("invoice_id", params.invoiceId);

      const url = `${API_BASE_URLS.finance}${FINANCE_ENDPOINTS.payments}${queryParams.toString() ? `?${queryParams}` : ""}`;
      const response = await fetch(url, {
        headers: { Authorization: `Bearer ${token}`, "Content-Type": "application/json" },
      });

      if (!response.ok) throw new Error("Failed to fetch payments");
      const data = await response.json();
      setPayments(data.data || []);
      return data.data || [];
    } catch (err) {
      setError(err instanceof Error ? err.message : "Error");
      return [];
    } finally {
      setIsLoading(false);
    }
  }, []);

  const createInvoice = useCallback(async (invoiceData: Partial<Invoice>) => {
    setIsLoading(true);
    try {
      const token = getToken();
      if (!token) throw new Error("Not authenticated");

      const response = await fetch(`${API_BASE_URLS.finance}${FINANCE_ENDPOINTS.invoices}`, {
        method: "POST",
        headers: { Authorization: `Bearer ${token}`, "Content-Type": "application/json" },
        body: JSON.stringify(invoiceData),
      });

      if (!response.ok) throw new Error("Failed to create invoice");
      await fetchInvoices();
      return true;
    } catch (err) {
      setError(err instanceof Error ? err.message : "Error");
      return false;
    } finally {
      setIsLoading(false);
    }
  }, [fetchInvoices]);

  const verifyPayment = useCallback(async (paymentId: string) => {
    setIsLoading(true);
    try {
      const token = getToken();
      if (!token) throw new Error("Not authenticated");

      const response = await fetch(`${API_BASE_URLS.finance}${FINANCE_ENDPOINTS.payments}/${paymentId}/verify`, {
        method: "POST",
        headers: { Authorization: `Bearer ${token}`, "Content-Type": "application/json" },
      });

      if (!response.ok) throw new Error("Failed to verify payment");
      await fetchPayments();
      return true;
    } catch (err) {
      setError(err instanceof Error ? err.message : "Error");
      return false;
    }
  }, [fetchPayments]);

  // Format utilities
  const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat("id-ID", {
      style: "currency",
      currency: "IDR",
      minimumFractionDigits: 0,
    }).format(amount);
  };

  return {
    invoices,
    payments,
    stats,
    isLoading,
    error,
    fetchInvoices,
    fetchPayments,
    createInvoice,
    verifyPayment,
    formatCurrency,
  };
}
