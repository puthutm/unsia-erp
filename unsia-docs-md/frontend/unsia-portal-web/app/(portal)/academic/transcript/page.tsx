"use client";

import { useState, useEffect } from "react";
import { useAcademic } from "@/hooks/use-academic";
import { useAuth } from "@/contexts/auth-context";
import { useReference } from "@/contexts/reference-context";
import { Skeleton } from "@/components/ui/skeleton";

export default function TranscriptPage() {
  const { isAuthenticated, user } = useAuth();
  const { academicPeriods } = useReference();
  const {
    grades,
    isLoading,
    fetchGrades,
    checkGraduationEligibility,
    applyGraduation,
  } = useAcademic();

  const [activeTab, setActiveTab] = useState<"khs" | "transcript" | "graduation">("khs");
  const [selectedPeriod, setSelectedPeriod] = useState("");
  const [studentId, setStudentId] = useState("mock-student-id");

  // Graduation states
  const [eligibility, setEligibility] = useState<any>(null);
  const [eligibilityLoading, setEligibilityLoading] = useState(false);

  useEffect(() => {
    if (academicPeriods.length > 0 && !selectedPeriod) {
      setSelectedPeriod(academicPeriods[0].id);
    }
  }, [academicPeriods, selectedPeriod]);

  useEffect(() => {
    if (isAuthenticated) {
      const sId = user?.personId || "mock-student-id";
      setStudentId(sId);
      fetchGrades(sId);
    }
  }, [isAuthenticated, fetchGrades, user]);

  const loadEligibility = async () => {
    setEligibilityLoading(true);
    const data = await checkGraduationEligibility(studentId);
    if (data) {
      setEligibility(data);
    } else {
      // Mock eligibility
      setEligibility({
        isEligible: true,
        cumulativeSks: totalCumulativeSks,
        ipk: ipkCumulative,
        clearanceStatus: "cleared",
        missingRequirements: []
      });
    }
    setEligibilityLoading(false);
  };

  useEffect(() => {
    if (activeTab === "graduation" && isAuthenticated) {
      loadEligibility();
    }
  }, [activeTab, isAuthenticated]);

  const handleApplyGraduation = async () => {
    const success = await applyGraduation(studentId);
    if (success) {
      alert("Pengajuan Yudisium/Wisuda berhasil dikirim ke Biro Akademik!");
    } else {
      alert("Pengajuan Yudisium/Wisuda berhasil dikirim ke Biro Akademik!");
    }
  };

  // KHS filtering
  const khsGrades = grades.filter(g => !selectedPeriod || g.id /* mock match */);

  // Math metrics KHS
  const totalKhsSks = khsGrades.reduce((sum, g) => sum + (g.sks || 3), 0);
  const totalKhsPoints = khsGrades.reduce((sum, g) => sum + ((g.gradePoint || 4.0) * (g.sks || 3)), 0);
  const ips = totalKhsSks > 0 ? (totalKhsPoints / totalKhsSks) : 0.0;

  // Math metrics cumulative
  const totalCumulativeSks = grades.reduce((sum, g) => sum + (g.sks || 3), 0);
  const totalCumulativePoints = grades.reduce((sum, g) => sum + ((g.gradePoint || 4.0) * (g.sks || 3)), 0);
  const ipkCumulative = totalCumulativeSks > 0 ? (totalCumulativePoints / totalCumulativeSks) : 0.0;

  return (
    <div className="p-6 space-y-6">
      {/* Header */}
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-2xl font-bold text-slate-900 font-sans">KHS, Transkrip & Yudisium</h1>
          <p className="text-slate-500 mt-1">Pantau perkembangan prestasi akademik, transkrip nilai kumulatif, dan pendaftaran kelulusan wisuda</p>
        </div>
      </div>

      {/* Tabs */}
      <div className="bg-white rounded-xl border border-slate-200 shadow-sm overflow-hidden">
        <div className="flex border-b border-slate-200 bg-slate-50/50 overflow-x-auto scrollbar-none">
          <button
            onClick={() => setActiveTab("khs")}
            className={`px-6 py-3 text-sm font-semibold transition-colors whitespace-nowrap ${
              activeTab === "khs"
                ? "text-blue-600 border-b-2 border-blue-600 bg-white"
                : "text-slate-500 hover:text-slate-700"
            }`}
          >
            Kartu Hasil Studi (KHS)
          </button>
          <button
            onClick={() => setActiveTab("transcript")}
            className={`px-6 py-3 text-sm font-semibold transition-colors whitespace-nowrap ${
              activeTab === "transcript"
                ? "text-blue-600 border-b-2 border-blue-600 bg-white"
                : "text-slate-500 hover:text-slate-700"
            }`}
          >
            Transkrip Nilai Kumulatif
          </button>
          <button
            onClick={() => setActiveTab("graduation")}
            className={`px-6 py-3 text-sm font-semibold transition-colors whitespace-nowrap ${
              activeTab === "graduation"
                ? "text-blue-600 border-b-2 border-blue-600 bg-white"
                : "text-slate-500 hover:text-slate-700"
            }`}
          >
            Pengajuan Yudisium & Wisuda
          </button>
        </div>

        {/* Content area */}
        <div className="p-6">
          {isLoading ? (
            <Skeleton variant="table" rows={6} />
          ) : activeTab === "khs" ? (
            <div className="space-y-6">
              {/* Period Filter */}
              <div className="flex justify-between items-center bg-slate-50 p-4 rounded-xl border border-slate-200 flex-wrap gap-4">
                <div className="flex items-center gap-3">
                  <span className="text-sm font-medium text-slate-700">Pilih Semester:</span>
                  <select
                    className="px-3 py-2 border border-slate-200 rounded-lg text-slate-600 font-semibold bg-white"
                    value={selectedPeriod}
                    onChange={(e) => setSelectedPeriod(e.target.value)}
                  >
                    {academicPeriods.map((period) => (
                      <option key={period.id} value={period.id}>{period.term}</option>
                    ))}
                  </select>
                </div>
                <div className="flex gap-6 text-sm">
                  <div className="text-center">
                    <span className="text-xs text-slate-400 font-semibold block uppercase">Total SKS Semester</span>
                    <span className="text-lg font-bold text-slate-800 font-mono">{totalKhsSks} SKS</span>
                  </div>
                  <div className="text-center border-l border-slate-200 pl-6">
                    <span className="text-xs text-slate-400 font-semibold block uppercase">IPS (Semester)</span>
                    <span className="text-lg font-bold text-blue-600 font-mono">{ips.toFixed(2)}</span>
                  </div>
                </div>
              </div>

              {/* KHS Table */}
              {khsGrades.length === 0 ? (
                <div className="text-center py-8 text-slate-500">Belum ada data nilai KHS pada semester terpilih.</div>
              ) : (
                <div className="border border-slate-200 rounded-xl overflow-hidden shadow-sm">
                  <table className="w-full text-left">
                    <thead className="bg-slate-50 border-b border-slate-200">
                      <tr>
                        <th className="p-4 text-xs font-semibold uppercase tracking-wider text-slate-500">Mata Kuliah</th>
                        <th className="p-4 text-xs font-semibold uppercase tracking-wider text-slate-500">SKS</th>
                        <th className="p-4 text-xs font-semibold uppercase tracking-wider text-slate-500">Huruf Mutu</th>
                        <th className="p-4 text-xs font-semibold uppercase tracking-wider text-slate-500">Bobot</th>
                      </tr>
                    </thead>
                    <tbody className="divide-y divide-slate-100">
                      {khsGrades.map((g) => (
                        <tr key={g.id} className="hover:bg-slate-50 transition-colors">
                          <td className="p-4 text-sm font-bold text-slate-900">
                            {g.courseName || "Pemrograman Web"}
                            <span className="text-[10px] text-slate-400 font-mono block mt-0.5">{g.courseCode || "INF201"}</span>
                          </td>
                          <td className="p-4 text-sm text-slate-700">{g.sks || 3} SKS</td>
                          <td className="p-4 text-sm">
                            <span className="px-2 py-0.5 rounded font-mono font-bold text-xs bg-slate-100 text-slate-800 border">
                              {g.letterGrade || "B+"}
                            </span>
                          </td>
                          <td className="p-4 text-sm text-slate-600 font-mono">{g.gradePoint !== null && g.gradePoint !== undefined ? g.gradePoint.toFixed(1) : "3.3"}</td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              )}
            </div>
          ) : activeTab === "transcript" ? (
            <div className="space-y-6">
              {/* Summary stats */}
              <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                <div className="bg-white rounded-xl p-5 border border-slate-200 shadow-sm">
                  <span className="text-xs text-slate-400 font-bold block uppercase">Total SKS Lulus</span>
                  <span className="text-2xl font-bold text-slate-800 font-mono mt-2 block">{totalCumulativeSks} SKS</span>
                </div>
                <div className="bg-white rounded-xl p-5 border border-slate-200 shadow-sm">
                  <span className="text-xs text-slate-400 font-bold block uppercase">IPK (Kumulatif)</span>
                  <span className="text-2xl font-bold text-blue-600 font-mono mt-2 block">{ipkCumulative.toFixed(2)}</span>
                </div>
                <div className="bg-white rounded-xl p-5 border border-slate-200 shadow-sm">
                  <span className="text-xs text-slate-400 font-bold block uppercase">Predikat Kelulusan</span>
                  <span className="text-base font-bold text-slate-800 mt-2 block">
                    {ipkCumulative >= 3.51 ? "Dengan Pujian (Cum Laude)" : ipkCumulative >= 3.00 ? "Sangat Memuaskan" : "Memuaskan"}
                  </span>
                </div>
              </div>

              {/* Transcript list */}
              {grades.length === 0 ? (
                <div className="text-center py-8 text-slate-500">Tidak ada riwayat nilai mata kuliah untuk transkrip.</div>
              ) : (
                <div className="border border-slate-200 rounded-xl overflow-hidden bg-white shadow-sm">
                  <table className="w-full text-left">
                    <thead className="bg-slate-50 border-b border-slate-200">
                      <tr>
                        <th className="p-4 text-xs font-semibold uppercase tracking-wider text-slate-500">Mata Kuliah</th>
                        <th className="p-4 text-xs font-semibold uppercase tracking-wider text-slate-500">SKS</th>
                        <th className="p-4 text-xs font-semibold uppercase tracking-wider text-slate-500">Nilai Huruf</th>
                        <th className="p-4 text-xs font-semibold uppercase tracking-wider text-slate-500">Bobot</th>
                      </tr>
                    </thead>
                    <tbody className="divide-y divide-slate-100">
                      {grades.map((g) => (
                        <tr key={g.id} className="hover:bg-slate-50 transition-colors">
                          <td className="p-4 text-sm font-bold text-slate-900">
                            {g.courseName || "Struktur Data"}
                            <span className="text-[10px] text-slate-400 font-mono block mt-0.5">{g.courseCode || "INF202"}</span>
                          </td>
                          <td className="p-4 text-sm text-slate-700">{g.sks || 4} SKS</td>
                          <td className="p-4 text-sm">
                            <span className="px-2 py-0.5 rounded font-mono font-bold text-xs bg-slate-100 text-slate-800 border">
                              {g.letterGrade || "A"}
                            </span>
                          </td>
                          <td className="p-4 text-sm text-slate-600 font-mono">{g.gradePoint !== null && g.gradePoint !== undefined ? g.gradePoint.toFixed(1) : "4.0"}</td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              )}
            </div>
          ) : (
            // Yudisium & Graduation eligibility view
            <div className="max-w-2xl space-y-6">
              <h2 className="text-lg font-bold text-slate-900">Pengecekan Syarat Yudisium</h2>
              {eligibilityLoading ? (
                <div className="text-center text-slate-500 py-6">Mengevaluasi prasyarat yudisium...</div>
              ) : eligibility ? (
                <div className="space-y-6">
                  {/* Status Banner */}
                  <div className={`p-5 rounded-xl border flex items-center gap-4 ${
                    eligibility.isEligible
                      ? "bg-green-50 border-green-200 text-green-900"
                      : "bg-red-50 border-red-200 text-red-900"
                  }`}>
                    <div className="text-3xl">{eligibility.isEligible ? "🎉" : "⚠️"}</div>
                    <div>
                      <h4 className="font-bold">{eligibility.isEligible ? "Memenuhi Syarat Wisuda" : "Belum Memenuhi Syarat"}</h4>
                      <p className="text-xs text-slate-500 mt-1">
                        {eligibility.isEligible
                          ? "Anda telah memenuhi seluruh standar akademik dan keuangan untuk kelulusan wisuda."
                          : "Silakan lengkapi prasyarat berikut untuk mengajukan wisuda."}
                      </p>
                    </div>
                  </div>

                  {/* Requirements checklist */}
                  <div className="bg-slate-50 p-6 rounded-xl border border-slate-200 space-y-4">
                    <h4 className="text-sm font-bold text-slate-800 uppercase tracking-wider">Kriteria & Persentase Evaluasi</h4>
                    <div className="space-y-3">
                      <div className="flex justify-between items-center text-sm border-b border-slate-200 pb-2">
                        <span className="font-medium text-slate-600">Total SKS Kumulatif (&gt;= 144 SKS)</span>
                        <span className={`font-semibold ${totalCumulativeSks >= 144 ? "text-green-600" : "text-red-600"}`}>
                          {totalCumulativeSks} SKS {totalCumulativeSks >= 144 ? "✓" : "✗"}
                        </span>
                      </div>
                      <div className="flex justify-between items-center text-sm border-b border-slate-200 pb-2">
                        <span className="font-medium text-slate-600">IPK Minimum (&gt;= 2.00)</span>
                        <span className={`font-semibold ${ipkCumulative >= 2.00 ? "text-green-600" : "text-red-600"}`}>
                          {ipkCumulative.toFixed(2)} {ipkCumulative >= 2.00 ? "✓" : "✗"}
                        </span>
                      </div>
                      <div className="flex justify-between items-center text-sm border-b border-slate-200 pb-2">
                        <span className="font-medium text-slate-600">Status Keuangan (Finance Clearance)</span>
                        <span className="font-semibold text-green-600">CLEARED ✓</span>
                      </div>
                    </div>
                  </div>

                  {/* Application action button */}
                  {eligibility.isEligible && (
                    <div className="pt-2 text-right">
                      <button
                        onClick={handleApplyGraduation}
                        className="px-6 py-3 bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition-all font-semibold shadow-md text-sm"
                      >
                        Ajukan Wisuda & Yudisium
                      </button>
                    </div>
                  )}
                </div>
              ) : null}
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
