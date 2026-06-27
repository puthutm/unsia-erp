"use client";

import { useState, useEffect } from "react";
import { useReference } from "@/contexts/reference-context";
import { useAuth } from "@/contexts/auth-context";
import { API_BASE_URLS, ACADEMIC_ENDPOINTS, STORAGE_KEYS } from "@/lib/constants";

interface Student {
  id: string;
  nim: string;
  name: string;
  studyProgramName: string;
  entryYear: string;
  status: string;
}

interface Course {
  id: string;
  code: string;
  name: string;
  sks: number;
  semester: number;
  isActive: boolean;
}

interface Schedule {
  id: string;
  courseName: string;
  className: string;
  day: string;
  startTime: string;
  endTime: string;
  room: string;
  lecturerName: string;
}

export default function AcademicPage() {
  const { user, isAuthenticated } = useAuth();
  const { studyPrograms, academicYears, academicPeriods } = useReference();
  const [students, setStudents] = useState<Student[]>([]);
  const [courses, setCourses] = useState<Course[]>([]);
  const [schedules, setSchedules] = useState<Schedule[]>([]);
  const [loading, setLoading] = useState(true);
  const [activeTab, setActiveTab] = useState<"students" | "courses" | "schedules">("students");

  useEffect(() => {
    if (isAuthenticated) {
      fetchAcademicData();
    }
  }, [isAuthenticated]);

  const fetchAcademicData = async () => {
    const token = localStorage.getItem(STORAGE_KEYS.accessToken);
    if (!token) return;

    setLoading(true);
    try {
      // Fetch students
      const studentsRes = await fetch(`${API_BASE_URLS.academic}${ACADEMIC_ENDPOINTS.students}`, {
        headers: {
          Authorization: `Bearer ${token}`,
          "Content-Type": "application/json",
        },
      });
      if (studentsRes.ok) {
        const data = await studentsRes.json();
        setStudents(data.data || []);
      }

      // Fetch courses
      const coursesRes = await fetch(`${API_BASE_URLS.academic}${ACADEMIC_ENDPOINTS.courses}`, {
        headers: {
          Authorization: `Bearer ${token}`,
          "Content-Type": "application/json",
        },
      });
      if (coursesRes.ok) {
        const data = await coursesRes.json();
        setCourses(data.data || []);
      }

      // Fetch schedules
      const schedulesRes = await fetch(`${API_BASE_URLS.academic}${ACADEMIC_ENDPOINTS.schedules}`, {
        headers: {
          Authorization: `Bearer ${token}`,
          "Content-Type": "application/json",
        },
      });
      if (schedulesRes.ok) {
        const data = await schedulesRes.json();
        setSchedules(data.data || []);
      }
    } catch (error) {
      console.error("Error fetching academic data:", error);
    } finally {
      setLoading(false);
    }
  };

  const getStatusBadge = (status: string) => {
    const styles: Record<string, string> = {
      active: "bg-green-100 text-green-800",
      inactive: "bg-gray-100 text-gray-800",
      graduated: "bg-blue-100 text-blue-800",
      drop_out: "bg-red-100 text-red-800",
      suspended: "bg-yellow-100 text-yellow-800",
    };
    return styles[status] || "bg-gray-100 text-gray-800";
  };

  return (
    <div className="p-6 space-y-6">
      {/* Header */}
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-2xl font-bold text-slate-900">Akademik</h1>
          <p className="text-slate-500 mt-1">Kelola kegiatan akademik</p>
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
          <h3 className="text-sm font-medium text-slate-500">Total Mahasiswa</h3>
          <p className="text-3xl font-bold text-slate-900 mt-2">{students.length}</p>
        </div>
        <div className="bg-white rounded-xl p-6 border border-slate-200">
          <h3 className="text-sm font-medium text-slate-500">Aktif</h3>
          <p className="text-3xl font-bold text-slate-900 mt-2">
            {students.filter(s => s.status === "active").length}
          </p>
        </div>
        <div className="bg-white rounded-xl p-6 border border-slate-200">
          <h3 className="text-sm font-medium text-slate-500">Mata Kuliah</h3>
          <p className="text-3xl font-bold text-slate-900 mt-2">{courses.length}</p>
        </div>
        <div className="bg-white rounded-xl p-6 border border-slate-200">
          <h3 className="text-sm font-medium text-slate-500">Program Studi</h3>
          <p className="text-3xl font-bold text-slate-900 mt-2">{studyPrograms.length}</p>
        </div>
      </div>

      {/* Tabs */}
      <div className="bg-white rounded-xl border border-slate-200">
        <div className="flex border-b border-slate-200">
          <button
            onClick={() => setActiveTab("students")}
            className={`px-6 py-3 text-sm font-medium transition-colors ${
              activeTab === "students"
                ? "text-blue-600 border-b-2 border-blue-600"
                : "text-slate-500 hover:text-slate-700"
            }`}
          >
            Mahasiswa
          </button>
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
            onClick={() => setActiveTab("schedules")}
            className={`px-6 py-3 text-sm font-medium transition-colors ${
              activeTab === "schedules"
                ? "text-blue-600 border-b-2 border-blue-600"
                : "text-slate-500 hover:text-slate-700"
            }`}
          >
            Jadwal Kuliah
          </button>
        </div>

        {/* Tab Content */}
        <div className="p-6">
          {loading ? (
            <div className="text-center text-slate-500 py-8">Memuat data...</div>
          ) : activeTab === "students" && students.length === 0 ? (
            <div className="text-center text-slate-500 py-8">Tidak ada mahasiswa</div>
          ) : activeTab === "students" ? (
            <table className="w-full">
              <thead className="bg-slate-50">
                <tr>
                  <th className="text-left p-4 text-sm font-medium text-slate-500">NIM</th>
                  <th className="text-left p-4 text-sm font-medium text-slate-500">Nama</th>
                  <th className="text-left p-4 text-sm font-medium text-slate-500">Program Studi</th>
                  <th className="text-left p-4 text-sm font-medium text-slate-500">Tahun Masuk</th>
                  <th className="text-left p-4 text-sm font-medium text-slate-500">Status</th>
                </tr>
              </thead>
              <tbody>
                {students.slice(0, 10).map((student) => (
                  <tr key={student.id} className="border-t border-slate-200">
                    <td className="p-4 text-slate-900 font-mono">{student.nim}</td>
                    <td className="p-4 text-slate-900">{student.name}</td>
                    <td className="p-4 text-slate-600">{student.studyProgramName}</td>
                    <td className="p-4 text-slate-600">{student.entryYear}</td>
                    <td className="p-4">
                      <span className={`px-2 py-1 rounded-full text-xs font-medium ${getStatusBadge(student.status)}`}>
                        {student.status}
                      </span>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          ) : activeTab === "courses" && courses.length === 0 ? (
            <div className="text-center text-slate-500 py-8">Tidak ada mata kuliah</div>
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
                  <div className="flex gap-4 mt-3 text-sm text-slate-500">
                    <span>{course.sks} SKS</span>
                    <span>Semester {course.semester}</span>
                  </div>
                </div>
              ))}
            </div>
          ) : activeTab === "schedules" && schedules.length === 0 ? (
            <div className="text-center text-slate-500 py-8">Tidak ada jadwal</div>
          ) : activeTab === "schedules" ? (
            <table className="w-full">
              <thead className="bg-slate-50">
                <tr>
                  <th className="text-left p-4 text-sm font-medium text-slate-500">Mata Kuliah</th>
                  <th className="text-left p-4 text-sm font-medium text-slate-500">Kelas</th>
                  <th className="text-left p-4 text-sm font-medium text-slate-500">Hari</th>
                  <th className="text-left p-4 text-sm font-medium text-slate-500">Jam</th>
                  <th className="text-left p-4 text-sm font-medium text-slate-500">Ruang</th>
                  <th className="text-left p-4 text-sm font-medium text-slate-500">Dosen</th>
                </tr>
              </thead>
              <tbody>
                {schedules.slice(0, 10).map((schedule) => (
                  <tr key={schedule.id} className="border-t border-slate-200">
                    <td className="p-4 text-slate-900">{schedule.courseName}</td>
                    <td className="p-4 text-slate-600">{schedule.className}</td>
                    <td className="p-4 text-slate-600">{schedule.day}</td>
                    <td className="p-4 text-slate-600">{schedule.startTime} - {schedule.endTime}</td>
                    <td className="p-4 text-slate-600">{schedule.room}</td>
                    <td className="p-4 text-slate-600">{schedule.lecturerName}</td>
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
