"use client";

import { useState, useEffect } from "react";

// Grades Page - Next.js
// Matches: UI/AKADEMIK/SA/

interface GradeItem {
  id: string;
  courseName: string;
  courseCode: string;
  sks: number;
  numericGrade: number;
  letterGrade: string;
  gradePoint: number;
  semester: number;
  academicYear: string;
}

export default function GradePage() {
  const [grades, setGrades] = useState<GradeItem[]>([]);
  const [loading, setLoading] = useState(true);
  const [ipk, setIpk] = useState(0);
  const [totalSks, setTotalSks] = useState(0);

  useEffect(() => {
    fetchGrades();
  }, []);

  const fetchGrades = async () => {
    setLoading(true);
    try {
      const token = localStorage.getItem("accessToken");
      const response = await fetch("/api/v1/academic/grades", {
        headers: { Authorization: `Bearer ${token}` },
      });
      if (response.ok) {
        const data = await response.json();
        setGrades(data.data || []);
        calculateIpk(data.data || []);
      }
    } catch (error) {
      console.error("Error fetching grades:", error);
    } finally {
      setLoading(false);
    }
  };

  const calculateIpk = (gradeList: GradeItem[]) => {
    let totalPoints = 0;
    let totalSksCount = 0;
    gradeList.forEach((g) => {
      totalPoints += g.gradePoint * g.sks;
      totalSksCount += g.sks;
    });
    setTotalSks(totalSksCount);
    setIpk(totalSksCount > 0 ? totalPoints / totalSksCount : 0);
  };

const getGradeColor = (letter: string) => {
    if (letter === "A") return "bg-green-100 text-green-800";
    if (letter === "A-") return "bg-green-50 text-green-700";
    if (letter === "B+") return "bg-blue-100 text-blue-800";
    if (letter === "B") return "bg-blue-50 text-blue-700";
    if (letter === "B-") return "bg-cyan-100 text-cyan-800";
    if (letter === "C") return "bg-yellow-100 text-yellow-800";
    if (letter === "D") return "bg-orange-100 text-orange-800";
    if (letter === "E") return "bg-red-100 text-red-800";
    return "bg-gray-100 text-gray-800";
  };

  return (
    <div className="p-6 space-y-6">
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-2xl font-bold text-slate-900">Nilai Akademik</h1>
          <p className="text-slate-500 mt-1">Riwayat nilai mahasiswa</p>
        </div>
      </div>

      {/* IPK Stats */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <div className="bg-white rounded-xl p-6 border border-slate-200">
          <h3 className="text-sm font-medium text-slate-500">IPK</h3>
          <p className="text-3xl font-bold text-slate-900 mt-2">{ipk.toFixed(2)}</p>
        </div>
        <div className="bg-white rounded-xl p-6 border border-slate-200">
          <h3 className="text-sm font-medium text-slate-500">Total SKS</h3>
          <p className="text-3xl font-bold text-slate-900 mt-2">{totalSks}</p>
        </div>
        <div className="bg-white rounded-xl p-6 border border-slate-200">
          <h3 className="text-sm font-medium text-slate-500">Mata Kuliah</h3>
          <p className="text-3xl font-bold text-slate-900 mt-2">{grades.length}</p>
        </div>
      </div>

      {/* Grades Table */}
      <div className="bg-white rounded-xl border border-slate-200 overflow-hidden">
        <table className="w-full">
          <thead className="bg-slate-50">
            <tr>
              <th className="text-left p-4 text-sm font-medium text-slate-500">Kode</th>
              <th className="text-left p-4 text-sm font-medium text-slate-500">Mata Kuliah</th>
              <th className="text-left p-4 text-sm font-medium text-slate-500">SKS</th>
              <th className="text-left p-4 text-sm font-medium text-slate-500">Nilai Angka</th>
              <th className="text-left p-4 text-sm font-medium text-slate-500">Nilai Huruf</th>
              <th className="text-left p-4 text-sm font-medium text-slate-500">Bobot</th>
              <th className="text-left p-4 text-sm font-medium text-slate-500">Semester</th>
            </tr>
          </thead>
          <tbody>
            {loading ? (
              <tr>
                <td colSpan={7} className="p-8 text-center text-slate-500">
                  Memuat nilai...
                </td>
              </tr>
            ) : grades.length === 0 ? (
              <tr>
                <td colSpan={7} className="p-8 text-center text-slate-500">
                  Belum ada nilai
                </td>
              </tr>
            ) : (
              grades.map((grade) => (
                <tr key={grade.id} className="border-t border-slate-200">
                  <td className="p-4 text-slate-900 font-mono">{grade.courseCode}</td>
                  <td className="p-4 text-slate-900">{grade.courseName}</td>
                  <td className="p-4 text-slate-600">{grade.sks}</td>
                  <td className="p-4 text-slate-600">{grade.numericGrade}</td>
                  <td className="p-4">
                    <span className={`px-2 py-1 rounded-full text-sm font-medium ${getGradeColor(grade.letterGrade)}`}>
                      {grade.letterGrade}
                    </span>
                  </td>
                  <td className="p-4 text-slate-600">{grade.gradePoint}</td>
                  <td className="p-4 text-slate-600">Semester {grade.semester}</td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
}
