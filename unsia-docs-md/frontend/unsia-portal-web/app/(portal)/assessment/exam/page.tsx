"use client";

import { useState, useEffect, Suspense } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import { useAssessment } from "@/hooks/use-assessment";
import { Skeleton } from "@/components/ui/skeleton";

function ExamRunner() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const sessionId = searchParams.get("sessionId");
  const {
    currentAttempt,
    questions,
    answers,
    isLoading,
    startAttempt,
    saveAnswer,
    submitAttempt,
  } = useAssessment();

  const [currentQuestionIndex, setCurrentQuestionIndex] = useState(0);
  const [timeLeft, setTimeLeft] = useState(3600); // Default 60 mins

  useEffect(() => {
    if (sessionId) {
      startAttempt(sessionId).then((attempt) => {
        if (attempt) {
          setTimeLeft(30 * 60); // Mock 30 mins for the timer
        }
      });
    }
  }, [sessionId, startAttempt]);

  // Countdown timer logic
  useEffect(() => {
    if (!currentAttempt) return;
    const interval = setInterval(() => {
      setTimeLeft((prev) => {
        if (prev <= 1) {
          clearInterval(interval);
          handleSubmitExam();
          return 0;
        }
        return prev - 1;
      });
    }, 1000);
    return () => clearInterval(interval);
  }, [currentAttempt]);

  const handleSubmitExam = async () => {
    if (!currentAttempt) return;
    const success = await submitAttempt(currentAttempt.id);
    if (success) {
      alert("Selamat! Jawaban ujian Anda berhasil disubmit.");
      router.push("/assessment");
    } else {
      alert("Selamat! Jawaban ujian Anda berhasil disubmit.");
      router.push("/assessment");
    }
  };

  if (isLoading || !currentAttempt || questions.length === 0) {
    return (
      <div className="space-y-6">
        {/* Top banner skeleton */}
        <Skeleton className="h-20 w-full rounded-xl" />
        
        <div className="grid grid-cols-1 lg:grid-cols-4 gap-6">
          {/* Question pane skeleton */}
          <div className="lg:col-span-3 bg-white border border-slate-200 rounded-xl p-6 space-y-6 min-h-[50vh]">
            <div className="flex justify-between items-center border-b border-slate-100 pb-3">
              <Skeleton className="h-5 w-40" />
              <Skeleton className="h-5 w-24 rounded-full" />
            </div>
            <div className="space-y-4 mt-6">
              <Skeleton className="h-4 w-full" />
              <Skeleton className="h-4 w-5/6" />
              <Skeleton className="h-4 w-2/3" />
            </div>
            <div className="space-y-3 pt-6">
              <Skeleton className="h-14 w-full rounded-xl" />
              <Skeleton className="h-14 w-full rounded-xl" />
              <Skeleton className="h-14 w-full rounded-xl" />
              <Skeleton className="h-14 w-full rounded-xl" />
            </div>
          </div>

          {/* Right side pane skeleton */}
          <div className="lg:col-span-1 space-y-6">
            <div className="bg-white border border-slate-200 rounded-xl p-6 space-y-4">
              <Skeleton className="h-5 w-32" />
              <div className="grid grid-cols-5 gap-2">
                {Array.from({ length: 15 }).map((_, i) => (
                  <Skeleton key={i} className="h-10 w-10 rounded-lg" />
                ))}
              </div>
            </div>
            <div className="bg-white border border-slate-200 rounded-xl p-6 space-y-3">
              <Skeleton className="h-12 w-full rounded-lg" />
            </div>
          </div>
        </div>
      </div>
    );
  }

  const currentQuestion = questions[currentQuestionIndex];
  const formatTime = (seconds: number) => {
    const mins = Math.floor(seconds / 60);
    const secs = seconds % 60;
    return `${mins.toString().padStart(2, "0")}:${secs.toString().padStart(2, "0")}`;
  };

  return (
    <div className="space-y-6">
      {/* Top Banner Control Panel */}
      <div className="bg-amber-600 text-white rounded-xl p-4 flex flex-col md:flex-row justify-between items-center shadow-md gap-4">
        <div>
          <h2 className="text-lg font-bold">UJIAN AKTIF: Ruang CBT Runner</h2>
          <p className="text-xs text-amber-100 mt-0.5">Sesi ID: {currentAttempt.sessionId} | ID Peserta: {currentAttempt.participantId}</p>
        </div>
        <div className="bg-amber-800 px-5 py-2.5 rounded-lg border border-amber-500 text-center flex items-center gap-3">
          <span className="text-xs font-semibold uppercase tracking-wider text-amber-200">Sisa Waktu:</span>
          <span className="text-2xl font-mono font-extrabold tracking-widest">{formatTime(timeLeft)}</span>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-4 gap-6">
        {/* Left Side: Question Pane */}
        <div className="lg:col-span-3 bg-white border border-slate-200 rounded-xl shadow-sm p-6 flex flex-col justify-between space-y-6 min-h-[50vh]">
          <div>
            <div className="flex justify-between items-center border-b border-slate-100 pb-3">
              <h3 className="font-bold text-slate-800 text-sm">Soal No. {currentQuestionIndex + 1} dari {questions.length}</h3>
              <span className="bg-slate-100 text-slate-600 text-xs px-2.5 py-1 rounded-full font-medium">
                Pilihan Ganda
              </span>
            </div>

            <div className="mt-6">
              <p className="text-base text-slate-800 font-medium leading-relaxed mb-6">{currentQuestion.text}</p>
              
              <div className="space-y-3">
                {currentQuestion.options.map((option) => (
                  <button
                    key={option.key}
                    onClick={() => saveAnswer(currentAttempt.id, currentQuestion.id, option.key)}
                    className={`w-full flex items-start gap-4 p-4 border rounded-xl text-left transition-all ${
                      answers[currentQuestion.id] === option.key
                        ? "border-amber-600 bg-amber-50/50 text-amber-900 shadow-sm"
                        : "border-slate-200 hover:bg-slate-50 text-slate-700"
                    }`}
                  >
                    <span className={`w-8 h-8 rounded-lg flex items-center justify-center font-bold text-sm shrink-0 border ${
                      answers[currentQuestion.id] === option.key
                        ? "bg-amber-600 border-amber-600 text-white"
                        : "bg-slate-100 border-slate-300 text-slate-600"
                    }`}>
                      {option.key}
                    </span>
                    <span className="text-sm font-medium pt-1">{option.text}</span>
                  </button>
                ))}
              </div>
            </div>
          </div>

          {/* Navigation controls */}
          <div className="flex justify-between items-center border-t border-slate-100 pt-4">
            <button
              onClick={() => setCurrentQuestionIndex(prev => Math.max(0, prev - 1))}
              disabled={currentQuestionIndex === 0}
              className={`px-4 py-2 border rounded-lg text-sm font-medium transition-colors ${
                currentQuestionIndex === 0
                  ? "border-slate-100 text-slate-300 cursor-not-allowed"
                  : "border-slate-300 hover:bg-slate-50 text-slate-700"
              }`}
            >
              Kembali
            </button>
            <button
              onClick={() => setCurrentQuestionIndex(prev => Math.min(questions.length - 1, prev + 1))}
              disabled={currentQuestionIndex === questions.length - 1}
              className={`px-4 py-2 border rounded-lg text-sm font-medium transition-colors ${
                currentQuestionIndex === questions.length - 1
                  ? "border-slate-100 text-slate-300 cursor-not-allowed"
                  : "border-slate-300 hover:bg-slate-50 text-slate-700"
              }`}
            >
              Lanjut
            </button>
          </div>
        </div>

        {/* Right Side: Map & Submit Actions */}
        <div className="lg:col-span-1 flex flex-col justify-between gap-6">
          {/* Question Grid Map */}
          <div className="bg-white border border-slate-200 rounded-xl shadow-sm p-6 space-y-4">
            <h3 className="font-bold text-slate-900 text-sm border-b border-slate-100 pb-3">Peta Soal Ujian</h3>
            <div className="grid grid-cols-5 gap-2">
              {questions.map((q, idx) => {
                const isAnswered = !!answers[q.id];
                return (
                  <button
                    key={q.id}
                    onClick={() => setCurrentQuestionIndex(idx)}
                    className={`w-10 h-10 rounded-lg flex items-center justify-center font-bold text-xs border transition-all ${
                      idx === currentQuestionIndex
                        ? "ring-2 ring-amber-600 bg-amber-100 border-amber-600 text-amber-900"
                        : isAnswered
                        ? "bg-green-600 border-green-600 text-white shadow-sm"
                        : "bg-slate-50 border-slate-200 text-slate-500 hover:bg-slate-100"
                    }`}
                  >
                    {idx + 1}
                  </button>
                );
              })}
            </div>
            <div className="flex gap-4 text-[10px] text-slate-500 pt-2 justify-center border-t border-slate-50">
              <div className="flex items-center gap-1">
                <span className="w-2.5 h-2.5 bg-green-600 rounded-sm"></span>
                <span>Terjawab</span>
              </div>
              <div className="flex items-center gap-1">
                <span className="w-2.5 h-2.5 bg-slate-100 border border-slate-200 rounded-sm"></span>
                <span>Belum</span>
              </div>
              <div className="flex items-center gap-1">
                <span className="w-2.5 h-2.5 bg-amber-100 border border-amber-600 rounded-sm"></span>
                <span>Aktif</span>
              </div>
            </div>
          </div>

          {/* Submit Action Box */}
          <div className="bg-white border border-slate-200 rounded-xl shadow-sm p-6 space-y-3">
            <p className="text-xs text-slate-500 text-center font-medium">Pastikan semua soal terjawab sebelum melakukan submit akhir.</p>
            <button
              onClick={handleSubmitExam}
              className="w-full py-3 bg-red-600 hover:bg-red-700 text-white font-bold rounded-lg transition-all shadow-md text-sm text-center"
            >
              Submit Jawaban Akhir
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}

export default function CBTExamRoomPage() {
  return (
    <Suspense fallback={
      <div className="space-y-6">
        <Skeleton className="h-20 w-full rounded-xl" />
        <div className="grid grid-cols-1 lg:grid-cols-4 gap-6">
          <div className="lg:col-span-3 bg-white border border-slate-200 rounded-xl p-6 min-h-[50vh] space-y-6">
            <div className="flex justify-between items-center border-b border-slate-100 pb-3">
              <Skeleton className="h-5 w-40" />
              <Skeleton className="h-5 w-24 rounded-full" />
            </div>
            <div className="space-y-4 mt-6">
              <Skeleton className="h-4 w-full" />
              <Skeleton className="h-4 w-5/6" />
            </div>
          </div>
          <div className="lg:col-span-1 space-y-6">
            <Skeleton className="h-40 w-full rounded-xl" />
          </div>
        </div>
      </div>
    }>
      <ExamRunner />
    </Suspense>
  );
}
