"use client";

import { useState, useEffect } from "react";
import { usePmb } from "@/hooks";

interface ExamSession {
  id: string;
  name: string;
  examDate: string;
  startTime: string;
  endTime: string;
  room: string;
  capacity: number;
  registered: number;
  status: "upcoming" | "ongoing" | "completed";
}

export default function UjianMasukPage() {
  const { selectionResults, isLoading, fetchSelectionResults } = usePmb();
  const [sessions, setSessions] = useState<ExamSession[]>([]);
  const [activeSession, setActiveSession] = useState<ExamSession | null>(null);

  useEffect(() => {
    fetchSelectionResults();
    // Load mock exam sessions
    setSessions([
      {
        id: "1",
        name: "Ujian Masuk Gelombang 1 - Sesi Pagi",
        examDate: "2024-03-15",
        startTime: "08:00",
        endTime: "10:00",
        room: "Gedung A - Ruang 101",
        capacity: 50,
        registered: 45,
        status: "completed",
      },
      {
        id: "2",
        name: "Ujian Masuk Gelombang 1 - Sesi Siang",
        examDate: "2024-03-15",
        startTime: "13:00",
        endTime: "15:00",
        room: "Gedung A - Ruang 102",
        capacity: 50,
        registered: 38,
        status: "completed",
      },
      {
        id: "3",
        name: "Ujian Masuk Gelombang 2 - Sesi Pagi",
        examDate: "2024-04-20",
        startTime: "08:00",
        endTime: "10:00",
        room: "Gedung B - Ruang 201",
        capacity: 60,
        registered: 55,
        status: "upcoming",
      },
    ]);
  }, [fetchSelectionResults]);

  return (
    <div className="p-6 space-y-6">
      {/* Header */}
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-2xl font-bold text-slate-900">Ujian Masuk</h1>
          <p className="text-slate-500 mt-1">Pengelolaan ujian masuk mahasiswa baru</p>
        </div>
        <button className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700">
          + Buat Sesi Ujian
        </button>
      </div>

      {/* Stats */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <div className="bg-white rounded-xl p-6 border border-slate-200">
          <h3 className="text-sm font-medium text-slate-500">Total Sesi</h3>
          <p className="text-3xl font-bold text-slate-900 mt-2">{sessions.length}</p>
        </div>
        <div className="bg-white rounded-xl p-6 border border-slate-200">
          <h3 className="text-sm font-medium text-slate-500">Akan Datang</h3>
          <p className="text-3xl font-bold text-yellow-600 mt-2">
            {sessions.filter(s => s.status === "upcoming").length}
          </p>
        </div>
        <div className="bg-white rounded-xl p-6 border border-slate-200">
          <h3 className="text-sm font-medium text-slate-500">Sedang Berlangsung</h3>
          <p className="text-3xl font-bold text-green-600 mt-2">
            {sessions.filter(s => s.status === "ongoing").length}
          </p>
        </div>
        <div className="bg-white rounded-xl p-6 border border-slate-200">
          <h3 className="text-sm font-medium text-slate-500">Selesai</h3>
          <p className="text-3xl font-bold text-slate-600 mt-2">
            {sessions.filter(s => s.status === "completed").length}
          </p>
        </div>
      </div>

      {/* Exam Sessions Table */}
      <div className="bg-white rounded-xl border border-slate-200">
        <div className="p-6 border-b border-slate-200">
          <h2 className="text-lg font-semibold text-slate-900">Daftar Sesi Ujian</h2>
        </div>
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead className="bg-slate-50">
              <tr>
                <th className="text-left p-4 text-sm font-medium text-slate-500">Nama Sesi</th>
                <th className="text-left p-4 text-sm font-medium text-slate-500">Tanggal</th>
                <th className="text-left p-4 text-sm font-medium text-slate-500">Waktu</th>
                <th className="text-left p-4 text-sm font-medium text-slate-500">Ruang</th>
                <th className="text-left p-4 text-sm font-medium text-slate-500">Kuota</th>
                <th className="text-left p-4 text-sm font-medium text-slate-500">Status</th>
                <th className="text-left p-4 text-sm font-medium text-slate-500">Aksi</th>
              </tr>
            </thead>
            <tbody>
              {sessions.map((session) => (
                <tr key={session.id} className="border-t border-slate-200">
                  <td className="p-4 text-slate-900 font-medium">{session.name}</td>
                  <td className="p-4 text-slate-600">
                    {new Date(session.examDate).toLocaleDateString("id-ID")}
                  </td>
                  <td className="p-4 text-slate-600">
                    {session.startTime} - {session.endTime}
                  </td>
                  <td className="p-4 text-slate-600">{session.room}</td>
                  <td className="p-4 text-slate-600">
                    {session.registered}/{session.capacity}
                  </td>
                  <td className="p-4">
                    <span className={`px-2 py-1 rounded-full text-xs ${
                      session.status === "upcoming" ? "bg-yellow-100 text-yellow-800" :
                      session.status === "ongoing" ? "bg-green-100 text-green-800" :
                      "bg-gray-100 text-gray-800"
                    }`}>
                      {session.status === "upcoming" ? "Akan Datang" :
                       session.status === "ongoing" ? "Berlangsung" : "Selesai"}
                    </span>
                  </td>
                  <td className="p-4">
                    <div className="flex gap-2">
                      <button className="text-blue-600 hover:underline">
                        Edit
                      </button>
                      <button className="text-green-600 hover:underline">
                        Mulai
                      </button>
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>

      {/* Quick Actions */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <div className="bg-white rounded-xl border border-slate-200 p-6">
          <h3 className="font-semibold text-slate-900 mb-4">Cetak Daftar Hadir</h3>
          <button className="w-full px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700">
            Generate Daftar Hadir
          </button>
        </div>
        <div className="bg-white rounded-xl border border-slate-200 p-6">
          <h3 className="font-semibold text-slate-900 mb-4">Input Nilai</h3>
          <button className="w-full px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700">
            Input Nilai Ujian
          </button>
        </div>
        <div className="bg-white rounded-xl border border-slate-200 p-6">
          <h3 className="font-semibold text-slate-900 mb-4">Hasil Seleksi</h3>
          <button className="w-full px-4 py-2 bg-purple-600 text-white rounded-lg hover:bg-purple-700">
            Proses Hasil
          </button>
        </div>
      </div>
    </div>
  );
}
