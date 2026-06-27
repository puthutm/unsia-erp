"use client";

import { useState, useEffect } from "react";

// Transcript Page - Next.js
// Matches: UI/AKADEMIK/SA/Panduan Pengambilan KRS

interface TranscriptCourse {
  id: string;
  courseCode: string;
  courseName: string;
  semester: number;
  sks: number;
  numericGrade: number;
  letterGrade: string;
  gradePoint: number;
}

export default function TranscriptPage() {
  const [courses, setCourses] = useState<TranscriptCourse[]>([]);
  const [loading, setLoading] = useState(true);
  const [studentInfo, setStudentInfo] = useState({
    nim: "",
    name: "",
    studyProgram: "",
    entryYear: "",
    gpa: 0,
    totalSks: 0,
  });

  useEffect(() => {
    fetchTranscript();
  }, []);

  const fetchTranscript = async () => {
    setLoading(true);
    try {
      const token = localStorage.getItem("accessToken");
      const response = await fetch("/api/v1/academic/transcript", {
        headers: { Authorization: `Bearer ${token}` },
      });
      if (response.ok) {
        const data = await response.json();
        setCourses(data.courses || []);
        setStudentInfo(data.student || getDefaultStudent());
      } else {
        setCourses(getDefaultCourses());
        setStudentInfo(getDefaultStudent());
      }
    } catch (error) {
      console.error("Error fetching transcript:", error);
      setCourses(getDefaultCourses());
      setStudentInfo(getDefaultStudent());
    } finally {
      setLoading(false);
    }
  };

  const getDefaultStudent = () => ({
    nim: "2024100101",
    name: "Ahmad Faisal",
    studyProgram: "Teknik Informatika",
    entryYear: "2024",
    gpa: 3.75,
    totalSks: 110,
  });

  const getDefaultCourses = () => [
    { id: "1", courseCode: "TI101", courseName: "Dasar Pemrograman", semester: 1, sks: 4, numericGrade: 85, letterGrade: "A", gradePoint: 4.0 },
    { id: "2", courseCode: "TI102", courseName: "Kalkulus I", semester: 1, sks: 3, numericGrade: 78, letterGrade: "B+", gradePoint: 3.5 },
    { id: "3", courseCode: "TI103", courseName: "Pengantar TI", semester: 1, sks: 2, numericGrade: 82, letterGrade: "A-", gradePoint: 3.7 },
    { id: "4", courseCode: "TI201", courseName: "Struktur Data", semester: 2, sks: 4, numericGrade: 88, letterGrade: "A", gradePoint: 4.0 },
    { id: "5", courseCode: "TI202", courseName: "Kalkulus II", semester: 2, sks: 3, numericGrade: 75, letterGrade: "B", gradePoint: 3.0 },
    { id: "6", courseCode: "TI203", courseName: "Basis Data", semester: 2, sks: 4, numericGrade: 90, letterGrade: "A", gradePoint: 4.0 },
    { id: "7", courseCode: "TI301", courseName: "Algoritma & Kompleksitas", semester: 3, sks: 3, numericGrade: 80, letterGrade: "A-", gradePoint: 3.7 },
    { id: "8", courseCode: "TI302", courseName: "Sistem Operasi", semester: 3, sks: 3, numericGrade: 85, letterGrade: "A", gradePoint: 4.0 },
  ];

  const getLetterGradeColor = (letter: string) => {
    if (letter.startsWith("A")) return "bg-green-100 text-green-800";
    if (letter.startsWith("B")) return "bg-blue-100 text-blue-800";
    if (letter.startsWith("C")) return "bg-yellow-100 text-yellow-800";
    if (letter.startsWith("D")) return "bg-orange-100 text-orange-800";
    return "bg-red-100 text-red-800";
  };

  const groupedCourses = courses.reduce((acc, course) => {
    const sem = course.semester;
    if (!acc[sem]) acc[sem] = [];
    acc[sem].push(course);
    return acc;
  }, {} as Record<number, TranscriptCourse[]>);

  const calculateSemesterGPA = (semCourses: TranscriptCourse[]) => {
    const totalPoints = semCourses.reduce((sum, c) => sum + c.gradePoint * c.sks, 0);
    const totalSks = semCourses.reduce((sum, c) => sum + c.sks, 0);
    return totalSks > 0 ? (totalPoints / totalSks).toFixed(2) : "0.00";
  };

  return (
    <div className="p-6 space-y-6">
      {/* Header */}
      <div className="bg-white rounded-xl border border-slate-200 p-6">
        <h1 className="text-2xl font-bold text-slate-900">Transkip Nilai</h1>
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mt-4">
          <div>
            <p className="text-sm text-slate-500">NIM</p>
            <p className="font-medium text-slate-900">{studentInfo.nim}</p>
          </div>
          <div>
            <p className="text-sm text-slate-500">Nama</p>
            <p className="font-medium text-slate-900">{studentInfo.name}</p>
          </div>
          <div>
            <p className="text-sm text-slate-500">Program Studi</p>
            <p className="font-medium text-slate-900">{studentInfo.studyProgram}</p>
          </div>
          <div>
            <p className="text-sm text-slate-500">Tahun Masuk</p>
            <p className="font-medium text-slate-900">{studentInfo.entryYear}</p>
          </div>
        </div>
      </div>

      {/* Summary */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <div className="bg-white rounded-xl p-6 border border-slate-200">
          <h3 className="text-sm font-medium text-slate-500">IPK</h3>
          <p className="text-3xl font-bold text-blue-600">{studentInfo.gpa}</p>
        </div>
        <div className="bg-white rounded-xl p-6 border border-slate-200">
          <h3 className="text-sm font-medium text-slate-500">Total SKS</h3>
          <p className="text-3xl font-bold text-green-600">{studentInfo.totalSks}</p>
        </div>
        <div className="bg-white rounded-xl p-6 border border-slate-200">
          <h3 className="text-sm font-medium text-slate-500">Mata Kuliah</h3>
          <p className="text-3xl font-bold text-purple-600">{courses.length}</p>
        </div>
      </div>

      {/* Transcript by Semester */}
      {loading ? (
        <div className="text-center text-slate-500 py-8">Memuat transkip...</div>
      ) : (
        Object.entries(groupedCourses).map(([semester, semCourses]) => (
          <div key={semester} className="bg-white rounded-xl border border-slate-200 overflow-hidden">
            <div className="p-4 bg-slate-50 border-b border-slate-200 flex justify-between items-center">
              <h2 className="font-semibold text-slate-900">Semester {semester}</h2>
              <span className="text-sm text-slate-600">IPS: {calculateSemesterGPA(semCourses || [])}</span>
            </div>
            <table className="w-full">
              <thead className="bg-slate-50">
                <tr>
                  <th className="text-left p-3 text-sm font-medium text-slate-500">Kode</th>
                  <th className="text-left p-3 text-sm font-medium text-slate-500">Mata Kuliah</th>
                  <th className="text-center p-3 text-sm font-medium text-slate-500">SKS</th>
                  <th className="text-center p-3 text-sm font-medium text-slate-500">Nilai Angka</th>
                  <th className="text-center p-3 text-sm font-medium text-slate-500">Nilai Huruf</th>
                  <th className="text-center p-3 text-sm font-medium text-slate-500">Bobot</th>
                </tr>
              </thead>
              <tbody>
                {(semCourses || []).map((course) => (
                  <tr key={course.id} className="border-t border-slate-200">
                    <td className="p-3 text-slate-900 font-mono">{course.courseCode}</td>
                    <td className="p-3 text-slate-900">{course.courseName}</td>
                    <td className="p-3 text-center text-slate-600">{course.sks}</td>
                    <td className="p-3 text-center text-slate-600">{course.numericGrade}</td>
                    <td className="p-3 text-center">
                      <span className={`px-2 py-1 rounded-full text-xs font-medium ${getLetterGradeColor(course.letterGrade)}`}>
                        {course.letterGrade}
                      </span>
                    </td>
                    <td className="p-3 text-center text-slate-600">{course.gradePoint}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        ))
      )}

      {/* Print Button */}
      <div className="flex justify-end">
        <button className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700">
          Cetak Transkip
        </button>
      </div>
    </div>
  );
}
