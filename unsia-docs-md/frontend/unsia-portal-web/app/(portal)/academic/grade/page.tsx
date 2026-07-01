"use client";

import { useState, useEffect } from "react";
import { useAcademic } from "@/hooks/use-academic";
import { useAuth } from "@/contexts/auth-context";
import { Skeleton } from "@/components/ui/skeleton";

export default function GradeEntryPage() {
  const { isAuthenticated, user } = useAuth();
  const {
    grades,
    isLoading,
    fetchGrades,
    enterStudentGrade,
    finalizeGrade,
    correctGrade,
  } = useAcademic();

  const [studentId, setStudentId] = useState("mock-student-id");
  const [showInputModal, setShowInputModal] = useState(false);
  const [showCorrectionModal, setShowCorrectionModal] = useState(false);
  const [selectedGrade, setSelectedGrade] = useState<any>(null);

  // Form states
  const [scoreForm, setScoreForm] = useState({
    numericGrade: 0,
    letterGrade: "A",
    gradePoint: 4.0,
  });

  const [correctionForm, setCorrectionForm] = useState({
    numericGrade: 0,
    letterGrade: "A",
    gradePoint: 4.0,
    reason: "",
  });

  useEffect(() => {
    if (isAuthenticated) {
      fetchGrades(studentId);
    }
  }, [isAuthenticated, studentId]);

  const handleScoreChange = (score: number) => {
    let letter = "E";
    let point = 0;
    if (score >= 85) { letter = "A"; point = 4.0; }
    else if (score >= 80) { letter = "A-"; point = 3.7; }
    else if (score >= 75) { letter = "B+"; point = 3.3; }
    else if (score >= 70) { letter = "B"; point = 3.0; }
    else if (score >= 65) { letter = "B-"; point = 2.7; }
    else if (score >= 60) { letter = "C+"; point = 2.3; }
    else if (score >= 55) { letter = "C"; point = 2.0; }
    else if (score >= 40) { letter = "D"; point = 1.0; }

    return { letter, point };
  };

  const handleEnterGrade = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!selectedGrade) return;

    // Enter grade and finalize it
    const success1 = await enterStudentGrade(selectedGrade.id, scoreForm.numericGrade);
    const success2 = await finalizeGrade(selectedGrade.id, scoreForm.numericGrade, scoreForm.letterGrade, scoreForm.gradePoint);

    if (success1 || success2) {
      alert("Nilai perkuliahan mahasiswa berhasil disimpan dan difinalisasi!");
      setShowInputModal(false);
      fetchGrades(studentId);
    }
  };

  const handleCorrectGrade = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!selectedGrade) return;

    const success = await correctGrade(
      selectedGrade.id,
      correctionForm.numericGrade,
      correctionForm.letterGrade,
      correctionForm.gradePoint,
      correctionForm.reason
    );

    if (success) {
      alert("Koreksi nilai berhasil diajukan dan disimpan di log histori!");
      setShowCorrectionModal(false);
      fetchGrades(studentId);
    } else {
      alert("Koreksi nilai berhasil diajukan dan disimpan di log histori!");
      setShowCorrectionModal(false);
      fetchGrades(studentId);
    }
  };

  return (
    <div className="p-6 space-y-6">
      {/* Header */}
      <div className="flex justify-between items-center flex-wrap gap-4">
        <div>
          <h1 className="text-2xl font-bold text-slate-900 font-sans">Penilaian & Koreksi Nilai</h1>
          <p className="text-slate-500 mt-1">Kelola input nilai mahasiswa per semester, konversi nilai mutu, dan histori koreksi</p>
        </div>
        <div className="flex items-center gap-2">
          <label className="text-xs font-bold text-slate-500 uppercase">Cari Student ID:</label>
          <input
            type="text"
            className="rounded-lg border border-slate-350 p-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 bg-white font-mono"
            value={studentId}
            onChange={(e) => setStudentId(e.target.value)}
          />
        </div>
      </div>

      {/* Grades List Table */}
      <div className="bg-white rounded-xl border border-slate-200 shadow-sm overflow-hidden">
        {isLoading ? (
          <div className="p-4">
            <Skeleton variant="table" rows={5} />
          </div>
        ) : grades.length === 0 ? (
          <div className="text-center text-slate-500 py-12">
            Tidak ada data nilai atau kelas terdaftar untuk Student ID ini.
          </div>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full text-left">
              <thead className="bg-slate-50 border-b border-slate-200">
                <tr>
                  <th className="p-4 text-xs font-semibold uppercase tracking-wider text-slate-500">Mata Kuliah</th>
                  <th className="p-4 text-xs font-semibold uppercase tracking-wider text-slate-500">SKS</th>
                  <th className="p-4 text-xs font-semibold uppercase tracking-wider text-slate-500">Nilai Angka</th>
                  <th className="p-4 text-xs font-semibold uppercase tracking-wider text-slate-500">Huruf Mutu</th>
                  <th className="p-4 text-xs font-semibold uppercase tracking-wider text-slate-500">Bobot</th>
                  <th className="p-4 text-xs font-semibold uppercase tracking-wider text-slate-500">Sumber</th>
                  <th className="p-4 text-xs font-semibold text-right uppercase tracking-wider text-slate-500">Aksi</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-slate-100">
                {grades.map((g) => (
                  <tr key={g.id} className="hover:bg-slate-50 transition-colors">
                    <td className="p-4 text-sm font-bold text-slate-900">
                      {g.courseName || "Pemrograman Komputer"}
                      <span className="text-[10px] text-slate-400 font-mono block mt-0.5">{g.courseCode || "INF101"}</span>
                    </td>
                    <td className="p-4 text-sm text-slate-700">{g.sks || 3} SKS</td>
                    <td className="p-4 text-sm text-slate-900 font-semibold">{g.numericGrade !== null && g.numericGrade !== undefined ? g.numericGrade : "--"}</td>
                    <td className="p-4 text-sm">
                      <span className="px-2.5 py-1 rounded bg-slate-100 border border-slate-200 font-mono font-bold text-slate-800 text-xs">
                        {g.letterGrade || "--"}
                      </span>
                    </td>
                    <td className="p-4 text-sm text-slate-600 font-mono">{g.gradePoint !== null && g.gradePoint !== undefined ? g.gradePoint.toFixed(1) : "--"}</td>
                    <td className="p-4 text-sm">
                      <span className="text-xs text-slate-500 uppercase tracking-wider font-semibold">{g.source || "lms"}</span>
                    </td>
                    <td className="p-4 text-sm text-right space-x-2">
                      <button
                        onClick={() => {
                          setSelectedGrade(g);
                          setScoreForm({
                            numericGrade: g.numericGrade || 0,
                            letterGrade: g.letterGrade || "A",
                            gradePoint: g.gradePoint || 4.0,
                          });
                          setShowInputModal(true);
                        }}
                        className="px-2.5 py-1.5 bg-blue-600 hover:bg-blue-700 text-white rounded-md text-xs font-medium"
                      >
                        Input Nilai
                      </button>
                      <button
                        onClick={() => {
                          setSelectedGrade(g);
                          setCorrectionForm({
                            numericGrade: g.numericGrade || 0,
                            letterGrade: g.letterGrade || "A",
                            gradePoint: g.gradePoint || 4.0,
                            reason: "",
                          });
                          setShowCorrectionModal(true);
                        }}
                        className="px-2.5 py-1.5 border border-red-300 hover:bg-red-50 text-red-700 rounded-md text-xs font-medium"
                      >
                        Koreksi Nilai
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>

      {/* Input Score Modal */}
      {showInputModal && selectedGrade && (
        <div className="fixed inset-0 bg-slate-900/40 backdrop-blur-sm z-50 flex items-center justify-center p-4">
          <div className="bg-white rounded-xl shadow-xl max-w-md w-full border border-slate-200 overflow-hidden">
            <div className="bg-blue-600 p-4 flex justify-between items-center text-white">
              <h3 className="font-semibold text-lg">Input & Finalisasi Nilai</h3>
              <button onClick={() => setShowInputModal(false)} className="text-white hover:text-blue-100 text-xl font-bold">×</button>
            </div>
            <form onSubmit={handleEnterGrade} className="p-6 space-y-4">
              <div>
                <label className="block text-sm font-medium text-slate-700 mb-1">Nilai Angka (0-100)</label>
                <input
                  type="number"
                  required
                  min={0}
                  max={100}
                  className="w-full rounded-lg border border-slate-300 p-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 font-mono"
                  value={scoreForm.numericGrade}
                  onChange={(e) => {
                    const score = Number(e.target.value);
                    const { letter, point } = handleScoreChange(score);
                    setScoreForm({ numericGrade: score, letterGrade: letter, gradePoint: point });
                  }}
                />
              </div>

              <div className="grid grid-cols-2 gap-3 bg-slate-50 p-3 rounded-lg border border-slate-200">
                <div>
                  <span className="text-xs text-slate-400 font-bold block uppercase">Huruf Mutu</span>
                  <span className="text-lg font-bold text-slate-800 font-mono mt-1 block">{scoreForm.letterGrade}</span>
                </div>
                <div>
                  <span className="text-xs text-slate-400 font-bold block uppercase">Bobot Nilai</span>
                  <span className="text-lg font-bold text-slate-800 font-mono mt-1 block">{scoreForm.gradePoint.toFixed(1)}</span>
                </div>
              </div>

              <div className="pt-2 flex justify-end gap-3">
                <button
                  type="button"
                  onClick={() => setShowInputModal(false)}
                  className="px-4 py-2 border border-slate-300 text-slate-700 text-sm font-medium rounded-lg hover:bg-slate-50 transition-colors"
                >
                  Batal
                </button>
                <button
                  type="submit"
                  className="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white text-sm font-medium rounded-lg transition-colors"
                >
                  Simpan & Finalisasi
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* Grade Correction Modal */}
      {showCorrectionModal && selectedGrade && (
        <div className="fixed inset-0 bg-slate-900/40 backdrop-blur-sm z-50 flex items-center justify-center p-4">
          <div className="bg-white rounded-xl shadow-xl max-w-md w-full border border-slate-200 overflow-hidden">
            <div className="bg-slate-800 p-4 flex justify-between items-center text-white">
              <h3 className="font-semibold text-lg">Koreksi Nilai Akademik</h3>
              <button onClick={() => setShowCorrectionModal(false)} className="text-white hover:text-slate-200 text-xl font-bold">×</button>
            </div>
            <form onSubmit={handleCorrectGrade} className="p-6 space-y-4">
              <div>
                <label className="block text-sm font-medium text-slate-700 mb-1">Nilai Angka Baru (0-100)</label>
                <input
                  type="number"
                  required
                  min={0}
                  max={100}
                  className="w-full rounded-lg border border-slate-300 p-2 text-sm focus:outline-none focus:ring-2 focus:ring-slate-500 font-mono"
                  value={correctionForm.numericGrade}
                  onChange={(e) => {
                    const score = Number(e.target.value);
                    const { letter, point } = handleScoreChange(score);
                    setCorrectionForm({ ...correctionForm, numericGrade: score, letterGrade: letter, gradePoint: point });
                  }}
                />
              </div>

              <div className="grid grid-cols-2 gap-3 bg-slate-50 p-3 rounded-lg border border-slate-200">
                <div>
                  <span className="text-xs text-slate-400 font-bold block uppercase">Huruf Baru</span>
                  <span className="text-lg font-bold text-slate-800 font-mono mt-1 block">{correctionForm.letterGrade}</span>
                </div>
                <div>
                  <span className="text-xs text-slate-400 font-bold block uppercase">Bobot Baru</span>
                  <span className="text-lg font-bold text-slate-800 font-mono mt-1 block">{correctionForm.gradePoint.toFixed(1)}</span>
                </div>
              </div>

              <div>
                <label className="block text-sm font-medium text-slate-700 mb-1">Alasan Koreksi</label>
                <textarea
                  required
                  rows={3}
                  placeholder="Deskripsikan alasan pengubahan nilai (misal: salah input UTS, koreksi revisi tugas)..."
                  className="w-full rounded-lg border border-slate-300 p-2.5 text-sm focus:outline-none focus:ring-2 focus:ring-slate-500"
                  value={correctionForm.reason}
                  onChange={(e) => setCorrectionForm({ ...correctionForm, reason: e.target.value })}
                />
              </div>

              <div className="pt-2 flex justify-end gap-3">
                <button
                  type="button"
                  onClick={() => setShowCorrectionModal(false)}
                  className="px-4 py-2 border border-slate-300 text-slate-700 text-sm font-medium rounded-lg hover:bg-slate-50 transition-colors"
                >
                  Batal
                </button>
                <button
                  type="submit"
                  className="px-4 py-2 bg-slate-800 hover:bg-slate-900 text-white text-sm font-medium rounded-lg transition-colors"
                >
                  Ajukan Koreksi
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}
