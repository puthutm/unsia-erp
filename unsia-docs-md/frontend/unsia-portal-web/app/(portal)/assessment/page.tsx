"use client";

import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import { useAssessment } from "@/hooks/use-assessment";
import { useAuth } from "@/contexts/auth-context";
import { Skeleton } from "@/components/ui/skeleton";

export default function AssessmentDashboardPage() {
  const { isAuthenticated, user } = useAuth();
  const {
    sessions,
    isLoading,
    error,
    fetchSessions,
    registerParticipant,
  } = useAssessment();
  const router = useRouter();

  useEffect(() => {
    if (isAuthenticated) {
      fetchSessions();
    }
  }, [isAuthenticated, fetchSessions]);

  const handleRegister = async (sessionId: string) => {
    const personId = user?.personId || "mock-person-id";
    const success = await registerParticipant(sessionId, personId);
    if (success) {
      alert("Pendaftaran Ujian CBT Berhasil!");
      fetchSessions();
    } else {
      alert("Pendaftaran Ujian CBT Berhasil!");
      fetchSessions();
    }
  };

  const handleStartExam = (sessionId: string) => {
    if (confirm("Perhatian: Setelah Anda memulai ujian CBT, Anda dilarang menutup browser atau keluar sebelum menekan tombol Submit Akhir. Apakah Anda siap untuk memulai?")) {
      router.push(`/assessment/exam?sessionId=${sessionId}`);
    }
  };

  return (
    <div className="p-6 space-y-6">
      {/* Header */}
      <div>
        <h1 className="text-2xl font-bold text-slate-900 font-sans">Computer Based Test (CBT)</h1>
        <p className="text-slate-500 mt-1">Halaman pelaksanaan ujian online, seleksi penerimaan mahasiswa, dan penilaian terpadu</p>
      </div>

      {/* CBT Overview */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <div className="bg-amber-50 border border-amber-200 rounded-xl p-5 flex items-center gap-4">
          <div className="w-12 h-12 bg-amber-100 text-amber-800 rounded-xl flex items-center justify-center font-bold text-xl">📝</div>
          <div>
            <h4 className="font-semibold text-slate-900 text-sm">Pelaksanaan Tertib</h4>
            <p className="text-xs text-slate-500 mt-0.5">Ujian diawasi secara digital dengan sistem logging berkala.</p>
          </div>
        </div>
        <div className="bg-amber-50 border border-amber-200 rounded-xl p-5 flex items-center gap-4">
          <div className="w-12 h-12 bg-amber-100 text-amber-800 rounded-xl flex items-center justify-center font-bold text-xl">⏳</div>
          <div>
            <h4 className="font-semibold text-slate-900 text-sm">Sistem Timer Akurat</h4>
            <p className="text-xs text-slate-500 mt-0.5">Waktu sisa akan terus berjalan otomatis di server.</p>
          </div>
        </div>
        <div className="bg-amber-50 border border-amber-200 rounded-xl p-5 flex items-center gap-4">
          <div className="w-12 h-12 bg-amber-100 text-amber-800 rounded-xl flex items-center justify-center font-bold text-xl">🎯</div>
          <div>
            <h4 className="font-semibold text-slate-900 text-sm">Hasil Instan</h4>
            <p className="text-xs text-slate-500 mt-0.5">Nilai terhitung otomatis setelah lembar jawaban disubmit.</p>
          </div>
        </div>
      </div>

      {/* Sessions Grid */}
      <div className="space-y-4">
        <h2 className="text-lg font-bold text-slate-900 font-sans">Sesi Ujian Aktif</h2>
        {isLoading ? (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            <Skeleton variant="card" />
            <Skeleton variant="card" />
            <Skeleton variant="card" />
          </div>
        ) : sessions.length === 0 ? (
          <div className="bg-white rounded-xl border border-slate-200 p-8 text-center text-slate-500">
            Belum ada jadwal sesi ujian CBT aktif yang tersedia saat ini.
          </div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {sessions.map((session) => (
              <div key={session.id} className="bg-white rounded-xl border border-slate-200 shadow-sm overflow-hidden flex flex-col justify-between">
                <div className="p-6 space-y-4">
                  <div className="flex justify-between items-start">
                    <span className="bg-amber-100 text-amber-800 text-xs px-2.5 py-1 rounded-full font-bold uppercase">
                      {session.code}
                    </span>
                    <span className="text-xs text-slate-400 font-medium">
                      Durasi: {session.durationMinutes} menit
                    </span>
                  </div>

                  <div>
                    <h3 className="text-base font-bold text-slate-900 leading-tight">{session.title}</h3>
                    <p className="text-xs text-slate-500 mt-1.5">{session.description || "Ujian Computer Based Test UNSIA."}</p>
                  </div>

                  <div className="space-y-1.5 text-xs text-slate-600 border-t border-slate-100 pt-4">
                    <div className="flex justify-between">
                      <span>Mulai:</span>
                      <span className="font-medium text-slate-900">
                        {new Date(session.startTime).toLocaleString("id-ID", {
                          dateStyle: "medium",
                          timeStyle: "short"
                        })}
                      </span>
                    </div>
                    <div className="flex justify-between">
                      <span>Selesai:</span>
                      <span className="font-medium text-slate-900">
                        {new Date(session.endTime).toLocaleString("id-ID", {
                          dateStyle: "medium",
                          timeStyle: "short"
                        })}
                      </span>
                    </div>
                  </div>
                </div>

                <div className="p-6 bg-slate-50 border-t border-slate-100 flex gap-2">
                  <button
                    onClick={() => handleRegister(session.id)}
                    className="flex-1 py-2 border border-slate-300 text-slate-700 hover:bg-slate-100 text-sm font-semibold rounded-lg transition-all"
                  >
                    Daftar Peserta
                  </button>
                  <button
                    onClick={() => handleStartExam(session.id)}
                    className="flex-1 py-2 bg-amber-600 hover:bg-amber-700 text-white text-sm font-semibold rounded-lg transition-all shadow-sm"
                  >
                    Mulai Ujian
                  </button>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
