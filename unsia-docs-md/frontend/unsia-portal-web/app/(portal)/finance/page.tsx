"use client";

import { useState, useEffect } from "react";
import { useReference } from "@/contexts/reference-context";
import { useAuth } from "@/contexts/auth-context";
import { API_BASE_URLS, FINANCE_ENDPOINTS, STORAGE_KEYS } from "@/lib/constants";

interface FinanceStats {
  totalReceivable: number;
  totalPayable: number;
  collected: number;
  pendingPayment: number;
}

interface Invoice {
  id: string;
  invoiceNumber: string;
  personName: string;
  personType: string;
  amount: number;
  status: string;
  dueDate: string;
  createdAt: string;
}

interface Payment {
  id: string;
  paymentNumber: string;
  invoiceNumber: string;
  personName: string;
  amount: number;
  method: string;
  status: string;
  paidAt: string;
}

export default function FinancePage() {
  const { user, isAuthenticated } = useAuth();
  const { paymentComponents, paymentMethods } = useReference();
  const [stats, setStats] = useState<FinanceStats | null>(null);
  const [invoices, setInvoices] = useState<Invoice[]>([]);
  const [payments, setPayments] = useState<Payment[]>([]);
  const [loading, setLoading] = useState(true);
  const [activeTab, setActiveTab] = useState<"invoices" | "payments" | "components">("invoices");

  useEffect(() => {
    if (isAuthenticated) {
      fetchFinanceData();
    }
  }, [isAuthenticated]);

  const fetchFinanceData = async () => {
    const token = localStorage.getItem(STORAGE_KEYS.accessToken);
    if (!token) return;

    setLoading(true);
    try {
      // Fetch invoices
      const invoicesRes = await fetch(`${API_BASE_URLS.finance}${FINANCE_ENDPOINTS.invoices}`, {
        headers: {
          Authorization: `Bearer ${token}`,
          "Content-Type": "application/json",
        },
      });
      if (invoicesRes.ok) {
        const data = await invoicesRes.json();
        setInvoices(data.data || []);
      }

      // Fetch payments
      const paymentsRes = await fetch(`${API_BASE_URLS.finance}${FINANCE_ENDPOINTS.payments}`, {
        headers: {
          Authorization: `Bearer ${token}`,
          "Content-Type": "application/json",
        },
      });
      if (paymentsRes.ok) {
        const data = await paymentsRes.json();
        setPayments(data.data || []);
      }
    } catch (error) {
      console.error("Error fetching finance data:", error);
    } finally {
      setLoading(false);
    }
  };

  const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat("id-ID", {
      style: "currency",
      currency: "IDR",
      minimumFractionDigits: 0,
    }).format(amount);
};

  const getStatusBadge = (status: string) => {
    const styles: Record<string, string> = {
      pending: "bg-yellow-100 text-yellow-800",
      paid: "bg-green-100 text-green-800",
      paid_: "bg-green-100 text-green-800",
      overdue: "bg-red-100 text-red-800",
      cancelled: "bg-gray-100 text-gray-800",
      partial: "bg-blue-100 text-blue-800",
      draft: "bg-gray-100 text-gray-800",
      issued: "bg-blue-100 text-blue-800",
      partially_paid: "bg-blue-100 text-blue-800",
      expired: "bg-red-100 text-red-800",
      verified: "bg-green-100 text-green-800",
      received: "bg-yellow-100 text-yellow-800",
      failed: "bg-red-100 text-red-800",
    };
    return styles[status] || "bg-gray-100 text-gray-800";
  };

return (
    <div className="p-6 space-y-6">
      {/* Header */}
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-2xl font-bold text-slate-900">Keuangan</h1>
          <p className="text-slate-500 mt-1">Kelola keuangan kampus</p>
        </div>
        <button className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors">
          + Buat Invoice
        </button>
      </div>

      {/* Stats Cards */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <div className="bg-white rounded-xl p-6 border border-slate-200">
          <h3 className="text-sm font-medium text-slate-500">Piutang</h3>
          <p className="text-2xl font-bold text-slate-900 mt-2">{formatCurrency(stats?.totalReceivable || 0)}</p>
        </div>
        <div className="bg-white rounded-xl p-6 border border-slate-200">
          <h3 className="text-sm font-medium text-slate-500">Hutang</h3>
          <p className="text-2xl font-bold text-slate-900 mt-2">{formatCurrency(stats?.totalPayable || 0)}</p>
        </div>
        <div className="bg-white rounded-xl p-6 border border-slate-200">
          <h3 className="text-sm font-medium text-slate-500">Terpilih</h3>
          <p className="text-2xl font-bold text-slate-900 mt-2">{formatCurrency(stats?.collected || 0)}</p>
        </div>
        <div className="bg-white rounded-xl p-6 border border-slate-200">
          <h3 className="text-sm font-medium text-slate-500">Menunggu Pembayaran</h3>
          <p className="text-2xl font-bold text-slate-900 mt-2">{stats?.pendingPayment || 0}</p>
        </div>
      </div>

      {/* Tabs */}
      <div className="bg-white rounded-xl border border-slate-200">
        <div className="flex border-b border-slate-200">
          <button
            onClick={() => setActiveTab("invoices")}
            className={`px-6 py-3 text-sm font-medium transition-colors ${
              activeTab === "invoices"
                ? "text-blue-600 border-b-2 border-blue-600"
                : "text-slate-500 hover:text-slate-700"
            }`}
          >
            Invoice
          </button>
          <button
            onClick={() => setActiveTab("payments")}
            className={`px-6 py-3 text-sm font-medium transition-colors ${
              activeTab === "payments"
                ? "text-blue-600 border-b-2 border-blue-600"
                : "text-slate-500 hover:text-slate-700"
            }`}
          >
            Pembayaran
          </button>
          <button
            onClick={() => setActiveTab("components")}
            className={`px-6 py-3 text-sm font-medium transition-colors ${
              activeTab === "components"
                ? "text-blue-600 border-b-2 border-blue-600"
                : "text-slate-500 hover:text-slate-700"
            }`}
          >
            Komponen Pembayaran
          </button>
        </div>

        {/* Tab Content */}
        <div className="p-6">
          {loading ? (
            <div className="text-center text-slate-500 py-8">Memuat data...</div>
          ) : activeTab === "invoices" && invoices.length === 0 ? (
            <div className="text-center text-slate-500 py-8">Tidak ada invoice</div>
          ) : activeTab === "invoices" ? (
            <table className="w-full">
              <thead className="bg-slate-50">
                <tr>
                  <th className="text-left p-4 text-sm font-medium text-slate-500">No. Invoice</th>
                  <th className="text-left p-4 text-sm font-medium text-slate-500">Pelanggan</th>
                  <th className="text-left p-4 text-sm font-medium text-slate-500">Jumlah</th>
                  <th className="text-left p-4 text-sm font-medium text-slate-500">Jatuh Tempo</th>
                  <th className="text-left p-4 text-sm font-medium text-slate-500">Status</th>
                </tr>
              </thead>
              <tbody>
                {invoices.slice(0, 10).map((invoice) => (
                  <tr key={invoice.id} className="border-t border-slate-200">
                    <td className="p-4 text-slate-900">{invoice.invoiceNumber}</td>
                    <td className="p-4 text-slate-600">
                      {invoice.personName}
                      <span className="text-xs text-slate-400 ml-2">({invoice.personType})</span>
                    </td>
                    <td className="p-4 text-slate-900">{formatCurrency(invoice.amount)}</td>
                    <td className="p-4 text-slate-600">{invoice.dueDate}</td>
                    <td className="p-4">
                      <span className={`px-2 py-1 rounded-full text-xs font-medium ${getStatusBadge(invoice.status)}`}>
                        {invoice.status}
                      </span>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          ) : activeTab === "payments" && payments.length === 0 ? (
            <div className="text-center text-slate-500 py-8">Tidak ada pembayaran</div>
          ) : activeTab === "payments" ? (
            <table className="w-full">
              <thead className="bg-slate-50">
                <tr>
                  <th className="text-left p-4 text-sm font-medium text-slate-500">No. Pembayaran</th>
                  <th className="text-left p-4 text-sm font-medium text-slate-500">Invoice</th>
                  <th className="text-left p-4 text-sm font-medium text-slate-500">Pelanggan</th>
                  <th className="text-left p-4 text-sm font-medium text-slate-500">Jumlah</th>
                  <th className="text-left p-4 text-sm font-medium text-slate-500">Status</th>
                </tr>
              </thead>
              <tbody>
                {payments.slice(0, 10).map((payment) => (
                  <tr key={payment.id} className="border-t border-slate-200">
                    <td className="p-4 text-slate-900">{payment.paymentNumber}</td>
                    <td className="p-4 text-slate-600">{payment.invoiceNumber}</td>
                    <td className="p-4 text-slate-600">{payment.personName}</td>
                    <td className="p-4 text-slate-900">{formatCurrency(payment.amount)}</td>
                    <td className="p-4">
                      <span className={`px-2 py-1 rounded-full text-xs font-medium ${getStatusBadge(payment.status)}`}>
                        {payment.status}
                      </span>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              {paymentComponents.map((component) => (
                <div key={component.id} className="p-4 border border-slate-200 rounded-lg">
                  <div className="flex justify-between items-start">
                    <div>
                      <h4 className="font-medium text-slate-900">{component.name}</h4>
                      <p className="text-sm text-slate-500">{component.code}</p>
                    </div>
                    <span className={`px-2 py-1 rounded-full text-xs ${component.isActive ? "bg-green-100 text-green-800" : "bg-gray-100 text-gray-800"}`}>
                      {component.isActive ? "Aktif" : "Nonaktif"}
                    </span>
                  </div>
                  <p className="text-lg font-semibold text-slate-900 mt-2">{formatCurrency(component.defaultAmount)}</p>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
