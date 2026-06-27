"use client";

import { useState, useEffect } from "react";

// Student Page - Next.js
// Matches: UI/AKADEMIK/MAHASISWA/

interface Student {
  id: string;
  nim: string;
  name: string;
  email: string;
  phone: string;
  studyProgramName: string;
  facultyName: string;
  entryYear: string;
  status: string;
  currentSemester: number;
  advisorName: string;
}

export default function StudentPage() {
  const [students, setStudents] = useState<Student[]>([]);
  const [loading, setLoading] = useState(true);
  const [searchQuery, setSearchQuery] = useState("");
  const [selectedStatus, setSelectedStatus] = useState("all");

  useEffect(() => {
    fetchStudents();
  }, []);

  const fetchStudents = async () => {
    setLoading(true);
    try {
      const token = localStorage.getItem("accessToken");
      const response = await fetch("/api/v1/academic/students", {
        headers: { Authorization: `Bearer ${token}` },
      });
      if (response.ok) {
        const data = await response.json();
        setStudents(data.data || []);
      }
    } catch (error) {
      console.error("Error fetching students:", error);
    } finally {
      setLoading(false);
    }
  };

  const filteredStudents = students.filter(s => {
    const matchesSearch = searchQuery === "" || 
      s.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
      s.nim.includes(searchQuery);
    const matchesStatus = selectedStatus === "all" || s.status === selectedStatus;
    return matchesSearch && matchesStatus;
  });

  const getStatusBadge = (status: string) => {
    if (status === "active") return "bg-green-100 text-green-800";
    if (status === "inactive") return "bg-gray-100 text-gray-800";
    if (status === "graduated") return "bg-blue-100 text-blue-800";
    if (status === "drop_out") return "bg-red-100 text-red-800";
    return "bg-gray-100 text-gray-800";
  };

  return (
    <div className="p-6 space-y-6">
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-2xl font-bold text-slate-900">Data Mahasiswa</h1>
          <p className="text-slate-500 mt-1">Kelola data mahasiswa</p>
        </div>
        <button className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700">
          + Tambah Mahasiswa
        </button>
      </div>

      {/* Filters */}
      <div className="flex gap-4">
        <input
          type="text"
          placeholder="Cari NIM atau nama..."
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          className="flex-1 px-4 py-2 border border-slate-200 rounded-lg"
        />
        <select
          value={selectedStatus}
          onChange={(e) => setSelectedStatus(e.target.value)}
          className="px-4 py-2 border border-slate-200 rounded-lg"
        >
          <option value="all">Semua Status</option>
          <option value="active">Aktif</option>
          <option value="inactive">Tidak Aktif</option>
          <option value="graduated">Lulus</option>
          <option value="drop_out">DO</option>
        </select>
      </div>

      {/* Stats */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <div className="bg-white rounded-xl p-6 border border-slate-200">
          <h3 className="text-sm font-medium text-slate-500">Total</h3>
          <p className="text-2xl font-bold text-slate-900">{students.length}</p>
        </div>
        <div className="bg-white rounded-xl p-6 border border-slate-200">
          <h3 className="text-sm font-medium text-slate-500">Aktif</h3>
          <p className="text-2xl font-bold text-green-600">
            {students.filter(s => s.status === "active").length}
          </p>
        </div>
        <div className="bg-white rounded-xl p-6 border border-slate-200">
          <h3 className="text-sm font-medium text-slate-500">Lulus</h3>
          <p className="text-2xl font-bold text-blue-600">
            {students.filter(s => s.status === "graduated").length}
          </p>
        </div>
        <div className="bg-white rounded-xl p-6 border border-slate-200">
          <h3 className="text-sm font-medium text-slate-500">DO</h3>
          <p className="text-2xl font-bold text-red-600">
            {students.filter(s => s.status === "drop_out").length}
          </p>
        </div>
      </div>

      {/* Student Table */}
      <div className="bg-white rounded-xl border border-slate-200 overflow-hidden">
        <table className="w-full">
          <thead className="bg-slate-50">
            <tr>
              <th className="text-left p-4 text-sm font-medium text-slate-500">NIM</th>
              <th className="text-left p-4 text-sm font-medium text-slate-500">Nama</th>
              <th className="text-left p-4 text-sm font-medium text-slate-500">Program Studi</th>
              <th className="text-left p-4 text-sm font-medium text-slate-500">Tahun Masuk</th>
              <th className="text-left p-4 text-sm font-medium text-slate-500">Semester</th>
              <th className="text-left p-4 text-sm font-medium text-slate-500">Status</th>
              <th className="text-left p-4 text-sm font-medium text-slate-500">Aksi</th>
            </tr>
          </thead>
          <tbody>
            {loading ? (
              <tr>
                <td colSpan={7} className="p-8 text-center text-slate-500">
                  Memuat...
                </td>
              </tr>
            ) : filteredStudents.length === 0 ? (
              <tr>
                <td colSpan={7} className="p-8 text-center text-slate-500">
                  Tidak ada mahasiswa
                </td>
              </tr>
            ) : (
              filteredStudents.map((student) => (
                <tr key={student.id} className="border-t border-slate-200">
                  <td className="p-4 text-slate-900 font-mono">{student.nim}</td>
                  <td className="p-4 text-slate-900">{student.name}</td>
                  <td className="p-4 text-slate-600">{student.studyProgramName}</td>
                  <td className="p-4 text-slate-600">{student.entryYear}</td>
                  <td className="p-4 text-slate-600">{student.currentSemester}</td>
                  <td className="p-4">
                    <span className={`px-2 py-1 rounded-full text-xs font-medium ${getStatusBadge(student.status)}`}>
                      {student.status}
                    </span>
                  </td>
                  <td className="p-4">
                    <button className="text-blue-600 hover:underline">Detail</button>
                  </td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
}
