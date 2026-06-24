"use client";

import { useState, useEffect } from "react";
import { useAuth } from "@/contexts/auth-context";
import { API_BASE_URLS, LMS_ENDPOINTS, STORAGE_KEYS } from "@/lib/constants";

interface LmsCourse {
  id: string;
  code: string;
  name: string;
  studyProgramName: string;
  semester: number;
  lecturerName: string;
  isActive: boolean;
}

interface LmsClass {
  id: string;
  className: string;
  courseId: string;
  courseName: string;
  lecturerName: string;
  schedule: string;
  room: string;
  maxStudents: number;
  enrolledCount: number;
}

interface Enrollment {
  id: string;
  studentNim: string;
  studentName: string;
  classId: string;
  className: string;
  courseName: string;
  enrolledAt: string;
  status: string;
}

interface Session {
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

export default function LmsPage() {
  const { user, isAuthenticated } = useAuth();
  const [courses, setCourses] = useState<LmsCourse[]>([]);
  const [classes, setClasses] = useState<LmsClass[]>([]);
  const [enrollments, setEnrollments] = useState<Enrollment[]>([]);
  const [sessions, setSessions] = useState<Session[]>([]);
  const [loading, setLoading] = useState(true);
  const [activeTab, setActiveTab] = useState<"courses" | "classes" | "enrollments" | "sessions">("courses");

  useEffect(() => {
    if (isAuthenticated) {
      fetchLmsData();
    }
  }, [isAuthenticated]);

  const fetchLmsData = async () => {
    const token = localStorage.getItem(STORAGE_KEYS.accessToken);
    if (!token) return;

    setLoading(true);
    try {
      // Fetch courses
      const coursesRes = await fetch(`${API_BASE_URLS.lms}${LMS_ENDPOINTS.courses}`, {
        headers: {
          Authorization: `Bearer ${token}`,
          "Content-Type": "application/json",
        },
      });
      if (coursesRes.ok) {
        const data = await coursesRes.json();
        setCourses(data.data || []);
      }

      // Fetch classes
      const classesRes = await fetch(`${API_BASE_URLS.lms}${LMS_ENDPOINTS.classes}`, {
        headers: {
          Authorization: `Bearer ${token}`,
          "Content-Type": "application/json",
        },
      });
      if (classesRes.ok) {
        const data = await classesRes.json();
        setClasses(data.data || []);
      }

      // Fetch enrollments
      const enrollmentsRes = await fetch(`${API_BASE_URLS.lms}${LMS_ENDPOINTS.enrollments}`, {
        headers: {
          Authorization: `Bearer ${token}`,
          "Content-Type": "application/json",
        },
      });
      if (enrollmentsRes.ok) {
        const data = await enrollmentsRes.json();
        setEnrollments(data.data || []);
      }

      // Fetch sessions
      const sessionsRes = await fetch(`${API_BASE_URLS.lms}${LMS_ENDPOINTS.sessions}`, {
        headers: {
          Authorization: `Bearer ${token}`,
          "Content-Type": "application/json",
        },
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
      active: "bg-green-100 text-green-800",
      inactive: "bg-gray-100 text-gray-800",
      ongoing: "bg-blue-100 text-blue-800",
      completed: "bg-purple-100 text-purple-800",
      scheduled: "bg-yellow-100 text-yellow-800",
      enrolled: "bg-green-100 text-green-800",
      dropped: "bg-red-100 text-red-800",
    };
    return styles[status] || "bg-gray-100 text-gray-800";
  };

  return (
    <div className="p-6 space-y-6">
      {/* Header */}
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-2xl font-bold text-dark-900">LMS - Learning Management System</h1>
          <p className="text-dark-500 mt-1">Kelola pembelajaran online</p>
        </div>
        <button className="px-4 py-2 bg-brand-600 text-white rounded-lg hover:bg-brand-700 transition-colors">
          + Buat Kelas
        </button>
      </div>

      {/* Stats Cards */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <div className="bg-white rounded-xl p-6 border border-surface-border">
          <h3 className="text-sm font-medium text-dark-500">Mata Kuliah</h3>
          <p className="text-3xl font-bold text-dark-900 mt-2">{courses.length}</p>
        </div>
        <div className="bg-white rounded-xl p-6 border border-surface-border">
          <h3 className="text-sm font-medium text-dark-500">Kelas Aktif</h3>
          <p className="text-3xl font-bold text-dark-900 mt-2">
            {classes.filter(c => c.enrolledCount > 0 && c.enrolledCount < c.maxStudents).length}
          </p>
        </div>
        <div className="bg-white rounded-xl p-6 border border-surface-border">
          <h3 className="text-sm font-medium text-dark-500">Total Enrollment</h3>
          <p className="text-3xl font-bold text-dark-900 mt-2">{enrollments.length}</p>
        </div>
        <div className="bg-white rounded-xl p-6 border border-surface-border">
          <h3 className="text-sm font-medium text-dark-500">Sesi Aktif</h3>
          <p className="text-3xl font-bold text-dark-900 mt-2">
            {sessions.filter(s => s.status === "ongoing").length}
          </p>
        </div>
      </div>

      {/* Tabs */}
      <div className="bg-white rounded-xl border border-surface-border">
        <div className="flex border-b border-surface-border">
          <button
            onClick={() => setActiveTab("courses")}
            className={`px-6 py-3 text-sm font-medium transition-colors ${
              activeTab === "courses"
                ? "text-brand-600 border-b-2 border-brand-600"
                : "text-dark-500 hover:text-dark-700"
            }`}
          >
            Mata Kuliah
          </button>
          <button
            onClick={() => setActiveTab("classes")}
            className={`px-6 py-3 text-sm font-medium transition-colors ${
              activeTab === "classes"
                ? "text-brand-600 border-b-2 border-brand-600"
                : "text-dark-500 hover:text-dark-700"
            }`}
          >
            Kelas
          </button>
          <button
            onClick={() => setActiveTab("enrollments")}
            className={`px-6 py-3 text-sm font-medium transition-colors ${
              activeTab === "enrollments"
                ? "text-brand-600 border-b-2 border-brand-600"
                : "text-dark-500 hover:text-dark-700"
            }`}
          >
            Enrollment
          </button>
          <button
            onClick={() => setActiveTab("sessions")}
            className={`px-6 py-3 text-sm font-medium transition-colors ${
              activeTab === "sessions"
                ? "text-brand-600 border-b-2 border-brand-600"
                : "text-dark-500 hover:text-dark-700"
            }`}
          >
            Sesi Kuliah
          </button>
        </div>

        {/* Tab Content */}
        <div className="p-6">
          {loading ? (
            <div className="text-center text-dark-500 py-8">Memuat data...</div>
          ) : activeTab === "courses" && courses.length === 0 ? (
            <div className="text-center text-dark-500 py-8">Tidak ada mata kuliah</div>
          ) : activeTab === "courses" ? (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              {courses.map((course) => (
                <div key={course.id} className="p-4 border border-surface-border rounded-lg">
                  <div className="flex justify-between items-start">
                    <div>
                      <h4 className="font-medium text-dark-900">{course.name}</h4>
                      <p className="text-sm text-dark-500">{course.code}</p>
                    </div>
                    <span className={`px-2 py-1 rounded-full text-xs ${course.isActive ? "bg-green-100 text-green-800" : "bg-gray-100 text-gray-800"}`}>
                      {course.isActive ? "Aktif" : "Nonaktif"}
                    </span>
                  </div>
                  <div className="mt-3 space-y-1 text-sm text-dark-500">
                    <p>{course.studyProgramName}</p>
                    <p>Semester {course.semester}</p>
                    <p className="text-dark-400">{course.lecturerName}</p>
                  </div>
                </div>
              ))}
            </div>
          ) : activeTab === "classes" && classes.length === 0 ? (
            <div className="text-center text-dark-500 py-8">Tidak ada kelas</div>
          ) : activeTab === "classes" ? (
            <table className="w-full">
              <thead className="bg-surface-subtle">
                <tr>
                  <th className="text-left p-4 text-sm font-medium text-dark-500">Kelas</th>
                  <th className="text-left p-4 text-sm font-medium text-dark-500">Mata Kuliah</th>
                  <th className="text-left p-4 text-sm font-medium text-dark-500">Dosen</th>
                  <th className="text-left p-4 text-sm font-medium text-dark-500">Jadwal</th>
                  <th className="text-left p-4 text-sm font-medium text-dark-500">Kapasitas</th>
                </tr>
              </thead>
              <tbody>
                {classes.slice(0, 10).map((cls) => (
                  <tr key={cls.id} className="border-t border-surface-border">
                    <td className="p-4 text-dark-900">{cls.className}</td>
                    <td className="p-4 text-dark-600">{cls.courseName}</td>
                    <td className="p-4 text-dark-600">{cls.lecturerName}</td>
                    <td className="p-4 text-dark-600">
                      <div>{cls.schedule}</div>
                      <div className="text-xs text-dark-400">{cls.room}</div>
                    </td>
                    <td className="p-4 text-dark-600">
                      {cls.enrolledCount}/{cls.maxStudents}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          ) : activeTab === "enrollments" && enrollments.length === 0 ? (
            <div className="text-center text-dark-500 py-8">Tidak ada enrollment</div>
          ) : activeTab === "enrollments" ? (
            <table className="w-full">
              <thead className="bg-surface-subtle">
                <tr>
                  <th className="text-left p-4 text-sm font-medium text-dark-500">NIM</th>
                  <th className="text-left p-4 text-sm font-medium text-dark-500">Nama</th>
                  <th className="text-left p-4 text-sm font-medium text-dark-500">Kelas</th>
                  <th className="text-left p-4 text-sm font-medium text-dark-500">Mata Kuliah</th>
                  <th className="text-left p-4 text-sm font-medium text-dark-500">Status</th>
                </tr>
              </thead>
              <tbody>
                {enrollments.slice(0, 10).map((enrollment) => (
                  <tr key={enrollment.id} className="border-t border-surface-border">
                    <td className="p-4 text-dark-900 font-mono">{enrollment.studentNim}</td>
                    <td className="p-4 text-dark-900">{enrollment.studentName}</td>
                    <td className="p-4 text-dark-600">{enrollment.className}</td>
                    <td className="p-4 text-dark-600">{enrollment.courseName}</td>
                    <td className="p-4">
                      <span className={`px-2 py-1 rounded-full text-xs font-medium ${getStatusBadge(enrollment.status)}`}>
                        {enrollment.status}
                      </span>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          ) : activeTab === "sessions" && sessions.length === 0 ? (
            <div className="text-center text-dark-500 py-8">Tidak ada sesi</div>
          ) : activeTab === "sessions" ? (
            <table className="w-full">
              <thead className="bg-surface-subtle">
                <tr>
                  <th className="text-left p-4 text-sm font-medium text-dark-500">Topik</th>
                  <th className="text-left p-4 text-sm font-medium text-dark-500">Kelas</th>
                  <th className="text-left p-4 text-sm font-medium text-dark-500">Mata Kuliah</th>
                  <th className="text-left p-4 text-sm font-medium text-dark-500">Jam</th>
                  <th className="text-left p-4 text-sm font-medium text-dark-500">Status</th>
                  <th className="text-left p-4 text-sm font-medium text-dark-500">Hadir</th>
                </tr>
              </thead>
              <tbody>
                {sessions.slice(0, 10).map((session) => (
                  <tr key={session.id} className="border-t border-surface-border">
                    <td className="p-4 text-dark-900">{session.topic}</td>
                    <td className="p-4 text-dark-600">{session.className}</td>
                    <td className="p-4 text-dark-600">{session.courseName}</td>
                    <td className="p-4 text-dark-600">
                      {session.startTime} - {session.endTime}
                    </td>
                    <td className="p-4">
                      <span className={`px-2 py-1 rounded-full text-xs font-medium ${getStatusBadge(session.status)}`}>
                        {session.status}
                      </span>
                    </td>
                    <td className="p-4 text-dark-600">{session.attendanceCount}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          ) : null}
        </div>
      </div>
    </div>
  );
}
