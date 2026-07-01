"use client";

import { useState, useEffect } from "react";
import { useAuth } from "@/contexts/auth-context";
import { useReference } from "@/contexts/reference-context";

interface Course {
  id: string;
  code: string;
  name: string;
  lecturerName: string;
  semester: number;
  isActive: boolean;
}

interface Session {
  id: string;
  courseId: string;
  courseName: string;
  title: string;
  description: string;
  scheduledAt: string;
  duration: number;
  status: "upcoming" | "ongoing" | "completed";
  materialCount: number;
  assignmentCount: number;
}

export default function LMSPage() {
  const { user, isAuthenticated } = useAuth();
  const { academicYears, academicPeriods } = useReference();
  const [courses, setCourses] = useState<Course[]>([]);
  const [sessions, setSessions] = useState<Session[]>([]);
  const [loading, setLoading] = useState(true);
  const [activeTab, setActiveTab] = useState<"courses" | "sessions" | "materials">("courses");

  useEffect(() => {
    if (isAuthenticated) {
      fetchLMSData();
    }
  }, [isAuthenticated]);

  const fetchLMSData = async () => {
    const token = localStorage.getItem("unsia_access_token");
    if (!token) return;

    setLoading(true);
    try {
      // Fetch courses
      const coursesRes = await fetch("http://localhost:8006/api/v1/lms/courses", {
        headers: { Authorization: `Bearer ${token}` },
      });
      if (coursesRes.ok) {
        const data = await coursesRes.json();
        setCourses(data.data || []);
      }

      // Fetch sessions
      const sessionsRes = await fetch("http://localhost:8006/api/v1/lms/sessions", {
        headers: { Authorization: `Bearer ${token}` },
      });
      if (sessionsRes.ok) {
        const data = await sessionsRes.json();
        setSessions(data.data || []);
      }
    } catch (error) {
      console.error("Error fetching LMS data:", error);
    } finally {
      setLoading(false);
    }
  };

  const getStatusBadge = (status: string) => {
    const styles: Record<string, string> = {
      upcoming: "bg-blue-100 text-blue-800",
      ongoing: "bg-green-100 text-green-800",
      completed: "bg-gray-100 text-gray-800",
    };
    return styles[status] || "bg-gray-100 text-gray-800";
  };

  return (
    <div className="p-6 space-y-6">
      {/* Header */}
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-2xl font-bold text-slate-900">LMS - Learning Management System</h1>
          <p className="text-slate-500 mt-1">Kelola pembelajaran online</p>
        </div>
        <div className="flex gap-2">
          <select className="px-3 py-2 border border-slate-200 rounded-lg text-slate-600">
            {academicYears.map((year) => (
              <option key={year.id} value={year.id}>{year.name}</option>
            ))}
          </select>
          <select className="px-3 py-2 border border-slate-200 rounded-lg text-slate-600">
            {academicPeriods.map((period) => (
              <option key={period.id} value={period.id}>{period.term}</option>
            ))}
          </select>
        </div>
      </div>

      {/* Stats Cards */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <div className="bg-white rounded-xl p-6 border border-slate-200">
          <h3 className="text-sm font-medium text-slate-500">Total Mata Kuliah</h3>
          <p className="text-3xl font-bold text-slate-900 mt-2">{courses.length}</p>
        </div>
        <div className="bg-white rounded-xl p-6 border border-slate-200">
          <h3 className="text-sm font-medium text-slate-500">Aktif</h3>
          <p className="text-3xl font-bold text-slate-900 mt-2">
            {courses.filter(c => c.isActive).length}
          </p>
        </div>
        <div className="bg-white rounded-xl p-6 border border-slate-200">
          <h3 className="text-sm font-medium text-slate-500">Sesi Berlangsung</h3>
          <p className="text-3xl font-bold text-slate-900 mt-2">
            {sessions.filter(s => s.status === "ongoing").length}
          </p>
        </div>
        <div className="bg-white rounded-xl p-6 border border-slate-200">
          <h3 className="text-sm font-medium text-slate-500">Sesi Selesai</h3>
          <p className="text-3xl font-bold text-slate-900 mt-2">
            {sessions.filter(s => s.status === "completed").length}
          </p>
        </div>
      </div>

      {/* Tabs */}
      <div className="bg-white rounded-xl border border-slate-200">
        <div className="flex border-b border-slate-200">
          <button
            onClick={() => setActiveTab("courses")}
            className={`px-6 py-3 text-sm font-medium transition-colors ${
              activeTab === "courses"
                ? "text-blue-600 border-b-2 border-blue-600"
                : "text-slate-500 hover:text-slate-700"
            }`}
          >
            Mata Kuliah
          </button>
          <button
            onClick={() => setActiveTab("sessions")}
            className={`px-6 py-3 text-sm font-medium transition-colors ${
              activeTab === "sessions"
                ? "text-blue-600 border-b-2 border-blue-600"
                : "text-slate-500 hover:text-slate-700"
            }`}
          >
            Sesi Kuliah
          </button>
          <button
            onClick={() => setActiveTab("materials")}
            className={`px-6 py-3 text-sm font-medium transition-colors ${
              activeTab === "materials"
                ? "text-blue-600 border-b-2 border-blue-600"
                : "text-slate-500 hover:text-slate-700"
            }`}
          >
            Materi
          </button>
        </div>

        {/* Tab Content */}
        <div className="p-6">
          {loading ? (
            <div className="text-center text-slate-500 py-8">Memuat data...</div>
          ) : activeTab === "courses" && courses.length === 0 ? (
            <div className="text-center text-slate-500 py-8">Tidak ada Mata Kuliah</div>
          ) : activeTab === "courses" ? (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              {courses.map((course) => (
                <div key={course.id} className="p-4 border border-slate-200 rounded-lg">
                  <div className="flex justify-between items-start">
                    <div>
                      <h4 className="font-medium text-slate-900">{course.name}</h4>
                      <p className="text-sm text-slate-500">{course.code}</p>
                    </div>
                    <span className={`px-2 py-1 rounded-full text-xs ${course.isActive ? "bg-green-100 text-green-800" : "bg-gray-100 text-gray-800"}`}>
                      {course.isActive ? "Aktif" : "Nonaktif"}
                    </span>
                  </div>
                  <div className="mt-3 text-sm text-slate-500">
                    <p>Dosen: {course.lecturerName}</p>
                    <p>Semester {course.semester}</p>
                  </div>
                </div>
              ))}
            </div>
          ) : activeTab === "sessions" && sessions.length === 0 ? (
            <div className="text-center text-slate-500 py-8">Tidak ada sesi</div>
          ) : activeTab === "sessions" ? (
            <table className="w-full">
              <thead className="bg-slate-50">
                <tr>
                  <th className="text-left p-4 text-sm font-medium text-slate-500">Mata Kuliah</th>
                  <th className="text-left p-4 text-sm font-medium text-slate-500">Topik</th>
                  <th className="text-left p-4 text-sm font-medium text-slate-500">Jadwal</th>
                  <th className="text-left p-4 text-sm font-medium text-slate-500">Durasi</th>
                  <th className="text-left p-4 text-sm font-medium text-slate-500">Status</th>
                </tr>
              </thead>
              <tbody>
                {sessions.slice(0, 10).map((session) => (
                  <tr key={session.id} className="border-t border-slate-200">
                    <td className="p-4 text-slate-900">{session.courseName}</td>
                    <td className="p-4 text-slate-600">{session.title}</td>
                    <td className="p-4 text-slate-600">{new Date(session.scheduledAt).toLocaleString()}</td>
                    <td className="p-4 text-slate-600">{session.duration} menit</td>
                    <td className="p-4">
                      <span className={`px-2 py-1 rounded-full text-xs font-medium ${getStatusBadge(session.status)}`}>
                        {session.status}
                      </span>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          ) : activeTab === "materials" ? (
            <div className="text-center text-slate-500 py-8">Materi - Dalam Pengembangan</div>
          ) : null}
        </div>
      </div>
    </div>
  );
}
