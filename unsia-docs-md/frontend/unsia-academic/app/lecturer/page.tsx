"use client";

import { useState, useEffect } from "react";
import { useReference } from "@/contexts/reference-context";
import { useAuth } from "@/contexts/auth-context";
import { API_BASE_URLS, ACADEMIC_ENDPOINTS, STORAGE_KEYS } from "@/lib/constants";

interface Lecturer {
  id: string;
  nip: string;
  name: string;
  email: string;
  phone: string;
  department: string;
  studyProgramName: string;
  position: string;
  status: string;
}

interface Schedule {
  id: string;
  courseName: string;
  className: string;
  day: string;
  startTime: string;
  endTime: string;
  room: string;
}

export default function LecturerPage() {
  const { user, isAuthenticated } = useAuth();
  const { studyPrograms } = useReference();
  const [lecturers, setLecturers] = useState<Lecturer[]>([]);
  const [mySchedules, setMySchedules] = useState<Schedule[]>([]);
  const [loading, setLoading] = useState(true);
  const [activeTab, setActiveTab] = useState<"list" | "schedule">("list");

  useEffect(() => {
    if (isAuthenticated) {
      fetchLecturers();
      fetchMySchedules();
    }
  }, [isAuthenticated]);

  const fetchLecturers = async () => {
    const token = localStorage.getItem(STORAGE_KEYS.accessToken);
    if (!token) return;

    setLoading(true);
    try {
      const response = await fetch(`${API_BASE_URLS.hris}${HRIS_ENDPOINTS.lecturers}`, {
        headers: {
          Authorization: `Bearer ${token}`,
          "Content-Type": "application/json",
        },
      });
      if (response.ok) {
        const data = await response.json();
        setLecturers(data.data || []);
      }
    } catch (error) {
      console.error("Error fetching lecturers:", error);
    } finally {
      setLoading(false);
    }
  };

  const fetchMySchedules = async () => {
    const token = localStorage.getItem(STORAGE_KEYS.accessToken);
    if (!token) return;

    try {
      const response = await fetch(`${API_BASE_URLS.academic}${ACADEMIC_ENDPOINTS.schedules}/my`, {
        headers: {
          Authorization: `Bearer ${token}`,
          "Content-Type": "application/json",
        },
      });
      if (response.ok) {
        const data = await response.json();
        setMySchedules(data.data || []);
      }
    } catch (error) {
      console.error("Error fetching schedules:", error);
    }
  };

  const getStatusBadge = (status: string) => {
    const styles: Record<string, string> = {
      active: "bg-green-100 text-green-800",
      inactive: "bg-gray-100 text-gray-800",
      on_leave: "bg-yellow-100 text-yellow-800",
    };
    return styles[status] || "bg-gray-100 text-gray-800";
  };

  return (
    <div className="p-6 space-y-6">
      {/* Header */}
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-2xl font-bold text-slate-900">Dosen Akademik</h1>
          <p className="text-slate-500 mt-1">Kelola data dosen dan jadwal mengajar</p>
        </div>
      </div>

      {/* Tabs */}
      <div className="bg-white rounded-xl border border-slate-200">
        <div className="flex border-b border-slate-200">
          <button
            onClick={() => setActiveTab("list")}
            className={`px-6 py-3 text-sm font-medium transition-colors ${
              activeTab === "list"
                ? "text-blue-600 border-b-2 border-blue-600"
                : "text-slate-500 hover:text-slate-700"
            }`}
          >
            Daftar Dosen
          </button>
          <button
            onClick={() => setActiveTab("schedule")}
            className={`px-6 py-3 text-sm font-medium transition-colors ${
              activeTab === "schedule"
                ? "text-blue-600 border-b-2 border-blue-600"
                : "text-slate-500 hover:text-slate-700"
            }`}
          >
            Jadwal Mengajar
          </button>
        </div>

        {/* Tab Content */}
        <div className="p-6">
          {loading ? (
            <div className="text-center text-slate-500 py-8">Memuat data...</div>
          ) : activeTab === "list" ? (
            lecturers.length === 0 ? (
              <div className="text-center text-slate-500 py-8">Belum ada dosen</div>
            ) : (
              <table className="w-full">
                <thead className="bg-slate-50">
                  <tr>
                    <th className="text-left p-4 text-sm font-medium text-slate-500">NIP</th>
                    <th className="text-left p-4 text-sm font-medium text-slate-500">Nama</th>
                    <th className="text-left p-4 text-sm font-medium text-slate-500">Email</th>
                    <th className="text-left p-4 text-sm font-medium text-slate-500">Program Studi</th>
                    <th className="text-left p-4 text-sm font-medium text-slate-500">Jabatan</th>
                    <th className="text-left p-4 text-sm font-medium text-slate-500">Status</th>
                  </tr>
                </thead>
                <tbody>
                  {lecturers.slice(0, 10).map((lecturer) => (
                    <tr key={lecturer.id} className="border-t border-slate-200">
                      <td className="p-4 text-slate-900 font-mono">{lecturer.nip}</td>
                      <td className="p-4 text-slate-900">{lecturer.name}</td>
                      <td className="p-4 text-slate-600">{lecturer.email}</td>
                      <td className="p-4 text-slate-600">{lecturer.studyProgramName}</td>
                      <td className="p-4 text-slate-600">{lecturer.position}</td>
                      <td className="p-4">
                        <span className={`px-2 py-1 rounded-full text-xs font-medium ${getStatusBadge(lecturer.status)}`}>
                          {lecturer.status}
                        </span>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            )
          ) : activeTab === "schedule" ? (
            mySchedules.length === 0 ? (
              <div className="text-center text-slate-500 py-8">Belum ada jadwal</div>
            ) : (
              <table className="w-full">
                <thead className="bg-slate-50">
                  <tr>
                    <th className="text-left p-4 text-sm font-medium text-slate-500">Mata Kuliah</th>
                    <th className="text-left p-4 text-sm font-medium text-slate-500">Kelas</th>
                    <th className="text-left p-4 text-sm font-medium text-slate-500">Hari</th>
                    <th className="text-left p-4 text-sm font-medium text-slate-500">Jam</th>
                    <th className="text-left p-4 text-sm font-medium text-slate-500">Ruang</th>
                  </tr>
                </thead>
                <tbody>
                  {mySchedules.map((schedule) => (
                    <tr key={schedule.id} className="border-t border-slate-200">
                      <td className="p-4 text-slate-900">{schedule.courseName}</td>
                      <td className="p-4 text-slate-600">{schedule.className}</td>
                      <td className="p-4 text-slate-600">{schedule.day}</td>
                      <td className="p-4 text-slate-600">{schedule.startTime} - {schedule.endTime}</td>
                      <td className="p-4 text-slate-600">{schedule.room}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            )
          ) : null}
        </div>
      </div>
    </div>
  );
}

// HRIS Endpoints - Tambahkan di constants.ts jika belum ada
const HRIS_ENDPOINTS = {
  lecturers: "/api/v1/hris/lecturers",
  employees: "/api/v1/hris/employees",
};
