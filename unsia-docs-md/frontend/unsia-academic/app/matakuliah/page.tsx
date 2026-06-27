"use client";

import { useState, useEffect } from "react";

// Mata Kuliah Page - Next.js
// Matches: UI/AKADEMIK/ADMIN/panel-matakuliah

interface Course {
  id: string;
  code: string;
  name: string;
  sks: number;
  semester: number;
  prodi: string;
  type: string; // Wajib/Pilihan
  coordinator: string;
  activeClasses: number;
  totalStudents: number;
}

export default function MataKuliahPage() {
  const [courses, setCourses] = useState<Course[]>([]);
  const [loading, setLoading] = useState(true);
  const [searchQuery, setSearchQuery] = useState("");
  const [filterProdi, setFilterProdi] = useState("");
  const [filterSemester, setFilterSemester] = useState("");
  const [filterType, setFilterType] = useState("");

  useEffect(() => {
    fetchCourses();
  }, []);

  const fetchCourses = async () => {
    setLoading(true);
    try {
      const token = localStorage.getItem("accessToken");
      const response = await fetch("/api/v1/academic/courses", {
        headers: { Authorization: `Bearer ${token}` },
      });
      if (response.ok) {
        const data = await response.json();
        setCourses(data.data || getDefaultCourses());
      } else {
        setCourses(getDefaultCourses());
      }
    } catch (error) {
      console.error("Error fetching courses:", error);
      setCourses(getDefaultCourses());
    } finally {
      setLoading(false);
    }
  };

  const getDefaultCourses = (): Course[] => [
    { id: "MK-IF201", code: "IF201", name: "Algoritma & Struktur Data", sks: 3, semester: 2, prodi: "S1 Informatika", type: "Wajib", coordinator: "Dr. Aulia Rahman, M.Kom.", activeClasses: 3, totalStudents: 87 },
    { id: "MK-IF203", code: "IF203", name: "Pemrograman Berorientasi Objek", sks: 4, semester: 2, prodi: "S1 Informatika", type: "Wajib", coordinator: "Noviandri, S.Kom., MMSI.", activeClasses: 3, totalStudents: 87 },
    { id: "MK-IF205", code: "IF205", name: "Basis Data", sks: 3, semester: 2, prodi: "S1 Informatika", type: "Wajib", coordinator: "Dr. Bayu Setiawan, M.T.", activeClasses: 3, totalStudents: 87 },
    { id: "MK-IF207", code: "IF207", name: "Jaringan Komputer", sks: 3, semester: 2, prodi: "S1 Informatika", type: "Wajib", coordinator: "Prof. Dr. Hendro Wijaksono", activeClasses: 2, totalStudents: 58 },
    { id: "MK-MK101", code: "MK101", name: "Pendidikan Pancasila", sks: 2, semester: 1, prodi: "Universal", type: "Wajib", coordinator: "Bp. Surya Hartanto", activeClasses: 6, totalStudents: 320 },
    { id: "MK-MK103", code: "MK103", name: "Bahasa Inggris", sks: 2, semester: 1, prodi: "Universal", type: "Wajib", coordinator: "Ms. Diana Kartika", activeClasses: 6, totalStudents: 320 },
    { id: "MK-MK105", code: "MK105", name: "Kewirausahaan", sks: 2, semester: 3, prodi: "Universal", type: "Wajib", coordinator: "Dr. Rini Susilowati", activeClasses: 5, totalStudents: 280 },
    { id: "MK-IF209", code: "IF209", name: "Pemrograman Web", sks: 4, semester: 3, prodi: "S1 Informatika", type: "Wajib", coordinator: "Noviandri, S.Kom., MMSI.", activeClasses: 3, totalStudents: 85 },
    { id: "MK-IF301", code: "IF301", name: "Rekayasa Perangkat Lunak", sks: 3, semester: 4, prodi: "S1 Informatika", type: "Wajib", coordinator: "Dr. Aulia Rahman, M.Kom.", activeClasses: 2, totalStudents: 76 },
    { id: "MK-IF401", code: "IF401", name: "AI & Machine Learning", sks: 3, semester: 5, prodi: "S1 Informatika", type: "Pilihan", coordinator: "Prof. Dr. Hendro Wijaksono", activeClasses: 1, totalStudents: 32 },
  ];

  const filteredCourses = courses.filter(course => {
    const matchesSearch = searchQuery === "" || 
      course.code.toLowerCase().includes(searchQuery.toLowerCase()) ||
      course.name.toLowerCase().includes(searchQuery.toLowerCase());
    const matchesProdi = filterProdi === "" || course.prodi === filterProdi;
    const matchesSemester = filterSemester === "" || course.semester === parseInt(filterSemester);
    const matchesType = filterType === "" || course.type === filterType;
    return matchesSearch && matchesProdi && matchesSemester && matchesType;
  });

  const getTypeBadge = (type: string) => {
    return type === "Wajib" 
      ? "bg-emerald-100 text-emerald-800" 
      : "bg-violet-100 text-violet-800";
  };

  return (
    <div className="p-6 space-y-6">
      {/* Header */}
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-2xl font-bold text-slate-900">Mata Kuliah</h1>
          <p className="text-slate-500 mt-1">Master mata kuliah seluruh prodi</p>
        </div>
        <button className="px-4 py-2 bg-blue-600 text-white rounded-lg font-medium flex items-center gap-2 hover:bg-blue-700">
          <span className="text-lg">+</span>
          <span>Tambah Mata Kuliah</span>
        </button>
      </div>

      {/* Filters */}
      <div className="bg-white rounded-xl border border-slate-200 p-4">
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
          <input
            type="text"
            placeholder="🔍 Cari kode/nama MK..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="px-3 py-2 bg-slate-50 border border-slate-200 rounded-lg text-sm focus:outline-none focus:border-blue-600"
          />
          <select
            value={filterProdi}
            onChange={(e) => setFilterProdi(e.target.value)}
            className="px-3 py-2 bg-slate-50 border border-slate-200 rounded-lg text-sm focus:outline-none focus:border-blue-600"
          >
            <option value="">Semua Prodi</option>
            <option>S1 Informatika</option>
            <option>S1 Sistem Informasi</option>
            <option>S1 Manajemen</option>
            <option>S1 Akuntansi</option>
            <option>S1 Psikologi</option>
          </select>
          <select
            value={filterSemester}
            onChange={(e) => setFilterSemester(e.target.value)}
            className="px-3 py-2 bg-slate-50 border border-slate-200 rounded-lg text-sm focus:outline-none focus:border-blue-600"
          >
            <option value="">Semua Semester</option>
            <option value="1">Semester 1</option>
            <option value="2">Semester 2</option>
            <option value="3">Semester 3</option>
            <option value="4">Semester 4</option>
            <option value="5">Semester 5</option>
            <option value="6">Semester 6</option>
            <option value="7">Semester 7</option>
            <option value="8">Semester 8</option>
          </select>
          <select
            value={filterType}
            onChange={(e) => setFilterType(e.target.value)}
            className="px-3 py-2 bg-slate-50 border border-slate-200 rounded-lg text-sm focus:outline-none focus:border-blue-600"
          >
            <option value="">Semua Jenis</option>
            <option>Wajib</option>
            <option>Pilihan</option>
          </select>
        </div>
      </div>

      {/* Table */}
      <div className="bg-white rounded-xl border border-slate-200 overflow-hidden">
        <div className="px-5 py-3 border-b border-slate-200 flex items-center justify-between">
          <h3 className="font-semibold text-slate-900">Daftar Mata Kuliah</h3>
          <span className="text-sm text-slate-500">{filteredCourses.length} mata kuliah</span>
        </div>
        <div className="overflow-x-auto">
          <table className="w-full text-sm">
            <thead className="bg-slate-50 border-b border-slate-200">
              <tr>
                <th className="text-left px-4 py-3 text-xs font-semibold text-slate-500 uppercase">Kode</th>
                <th className="text-left px-4 py-3 text-xs font-semibold text-slate-500 uppercase">Nama Mata Kuliah</th>
                <th className="text-center px-4 py-3 text-xs font-semibold text-slate-500 uppercase">SKS</th>
                <th className="text-center px-4 py-3 text-xs font-semibold text-slate-500 uppercase">Smt</th>
                <th className="text-left px-4 py-3 text-xs font-semibold text-slate-500 uppercase">Prodi</th>
                <th className="text-left px-4 py-3 text-xs font-semibold text-slate-500 uppercase">Jenis</th>
                <th className="text-left px-4 py-3 text-xs font-semibold text-slate-500 uppercase">Dosen Koordinator</th>
                <th className="text-center px-4 py-3 text-xs font-semibold text-slate-500 uppercase">Kelas</th>
                <th className="text-center px-4 py-3 text-xs font-semibold text-slate-500 uppercase">Aksi</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-100">
              {loading ? (
                <tr>
                  <td colSpan={9} className="p-8 text-center text-slate-500">
                    Memuat...
                  </td>
                </tr>
              ) : filteredCourses.length === 0 ? (
                <tr>
                  <td colSpan={9} className="p-8 text-center text-slate-500">
                    Tidak ada mata kuliah
                  </td>
                </tr>
              ) : (
                filteredCourses.map((course) => (
                  <tr key={course.id} className="hover:bg-slate-50">
                    <td className="px-4 py-3 font-mono font-bold text-blue-600">{course.code}</td>
                    <td className="px-4 py-3">
                      <div className="font-medium text-slate-900">{course.name}</div>
                    </td>
                    <td className="px-4 py-3 text-center">
                      <span className="inline-flex items-center justify-center w-8 h-8 bg-slate-100 rounded font-bold text-slate-700">
                        {course.sks}
                      </span>
                    </td>
                    <td className="px-4 py-3 text-center font-bold text-slate-700">{course.semester}</td>
                    <td className="px-4 py-3 text-slate-600">{course.prodi}</td>
                    <td className="px-4 py-3">
                      <span className={`px-2 py-1 rounded-full text-xs font-medium ${getTypeBadge(course.type)}`}>
                        {course.type}
                      </span>
                    </td>
                    <td className="px-4 py-3 text-slate-600">{course.coordinator}</td>
                    <td className="px-4 py-3 text-center">
                      <span className="px-2 py-1 bg-blue-100 text-blue-800 rounded text-xs font-medium">
                        {course.activeClasses} kelas
                      </span>
                      <div className="text-xs text-slate-500 mt-1">{course.totalStudents} mhs</div>
                    </td>
                    <td className="px-4 py-3 text-center">
                      <button className="px-3 py-1.5 bg-blue-600 hover:bg-blue-700 text-white rounded text-xs font-medium">
                        Buka Kelas
                      </button>
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}
