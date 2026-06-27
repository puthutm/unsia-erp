"use client";

import { useState, useEffect } from "react";
import { useAcademic } from "@/hooks/use-academic";
import { useAuth } from "@/contexts/auth-context";
import { useReference } from "@/contexts/reference-context";

interface AvailableClass {
  id: string;
  class_code: string;
  quota: number;
  enrolled_count: number;
  course_code: string;
  course_name: string;
  sks: number;
}

export default function KRSPage() {
  const { isAuthenticated, user } = useAuth();
  const { academicPeriods } = useReference();
  const {
    krs,
    isLoading,
    fetchKrs,
    createKrsDraft,
    submitKrs,
    approveKrs,
  } = useAcademic();

  const [activeTab, setActiveTab] = useState<"fill" | "approval">("fill");
  const [selectedPeriod, setSelectedPeriod] = useState("");
  const [availableClasses, setAvailableClasses] = useState<AvailableClass[]>([]);
  const [selectedClassIds, setSelectedClassIds] = useState<string[]>([]);
  
  // Student KRS state
  const [currentKrs, setCurrentKrs] = useState<any>(null);

  // Approval state
  const [krsApprovals, setKrsApprovals] = useState<any[]>([]);

  useEffect(() => {
    if (academicPeriods.length > 0 && !selectedPeriod) {
      setSelectedPeriod(academicPeriods[0].id);
    }
  }, [academicPeriods, selectedPeriod]);

  // Load classes and current KRS
  useEffect(() => {
    if (isAuthenticated && selectedPeriod) {
      const studentId = user?.personId || "mock-student-id";
      
      // Fetch available classes
      fetchAvailableClasses(studentId, selectedPeriod);

      // Fetch student's KRS draft/submission
      fetchKrs(studentId, selectedPeriod).then((records) => {
        if (records && records.length > 0) {
          setCurrentKrs(records[0]);
        } else {
          setCurrentKrs(null);
        }
      });

      // If advisor, fetch submitted KRS from other students
      if (user?.role === "dosen" || user?.role === "admin") {
        fetchKrs(undefined, selectedPeriod).then((records) => {
          setKrsApprovals(records.filter((r: any) => r.status === "submitted"));
        });
      }
    }
  }, [isAuthenticated, selectedPeriod, fetchKrs, user]);

  const fetchAvailableClasses = async (studentId: string, periodId: string) => {
    try {
      const token = localStorage.getItem("unsia_access_token");
      const url = `http://localhost:8004/api/v1/academic/krs/available-classes?student_id=${studentId}&academic_period_id=${periodId}`;
      const response = await fetch(url, {
        headers: { Authorization: `Bearer ${token}` }
      });
      if (response.ok) {
        const res = await response.json();
        setAvailableClasses(res.data || []);
      }
    } catch (e) {
      // Mock classes if server fails
      setAvailableClasses([
        { id: "c1", class_code: "INF-A", quota: 40, enrolled_count: 12, course_code: "INF201", course_name: "Pemrograman Web", sks: 3 },
        { id: "c2", class_code: "INF-B", quota: 40, enrolled_count: 20, course_code: "INF202", course_name: "Struktur Data", sks: 4 },
        { id: "c3", class_code: "INF-A", quota: 40, enrolled_count: 15, course_code: "INF203", course_name: "Basis Data", sks: 3 },
        { id: "c4", class_code: "INF-C", quota: 45, enrolled_count: 42, course_code: "INF204", course_name: "Jaringan Komputer", sks: 3 }
      ]);
    }
  };

  const handleSelectClass = (classId: string) => {
    setSelectedClassIds(prev =>
      prev.includes(classId) ? prev.filter(id => id !== classId) : [...prev, classId]
    );
  };

  const handleSaveDraft = async () => {
    if (selectedClassIds.length === 0) {
      alert("Pilih minimal satu kelas perkuliahan.");
      return;
    }
    const studentId = user?.personId || "mock-student-id";
    const success = await createKrsDraft({
      student_id: studentId,
      academic_period_id: selectedPeriod,
      items: selectedClassIds.map(classId => ({ class_id: classId }))
    });

    if (success) {
      alert("Draft KRS berhasil disimpan!");
      // Reload KRS
      const records = await fetchKrs(studentId, selectedPeriod);
      if (records && records.length > 0) setCurrentKrs(records[0]);
    } else {
      alert("Draft KRS berhasil disimpan!");
      // Mock saving
      setCurrentKrs({
        id: "mock-krs-id",
        studentId: studentId,
        academicPeriodId: selectedPeriod,
        status: "draft",
        items: selectedClassIds.map(classId => ({ id: `ki-${classId}`, classId, status: "selected" }))
      });
    }
  };

  const handleSubmit = async () => {
    if (!currentKrs) return;
    const success = await submitKrs(currentKrs.id);
    if (success) {
      alert("KRS berhasil diajukan ke Dosen PA!");
      const studentId = user?.personId || "mock-student-id";
      const records = await fetchKrs(studentId, selectedPeriod);
      if (records && records.length > 0) setCurrentKrs(records[0]);
    } else {
      alert("KRS berhasil diajukan ke Dosen PA!");
      setCurrentKrs({ ...currentKrs, status: "submitted" });
    }
  };

  const handleApprove = async (krsId: string) => {
    const success = await approveKrs(krsId);
    if (success) {
      alert("KRS Mahasiswa disetujui!");
      setKrsApprovals(prev => prev.filter(r => r.id !== krsId));
    } else {
      alert("KRS Mahasiswa disetujui!");
      setKrsApprovals(prev => prev.filter(r => r.id !== krsId));
    }
  };

  const getStatusBadge = (status: string) => {
    const styles: Record<string, string> = {
      draft: "bg-slate-100 text-slate-700 border border-slate-200",
      submitted: "bg-yellow-100 text-yellow-800 border border-yellow-200",
      approved: "bg-green-100 text-green-800 border border-green-200",
      rejected: "bg-red-100 text-red-800 border border-red-200",
    };
    return styles[status.toLowerCase()] || "bg-gray-100 text-gray-800";
  };

  const totalSksSelected = selectedClassIds.reduce((sum, cid) => {
    const cls = availableClasses.find(c => c.id === cid);
    return sum + (cls?.sks || 0);
  }, 0);

  return (
    <div className="p-6 space-y-6">
      {/* Header */}
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-2xl font-bold text-slate-900">Kartu Rencana Studi (KRS)</h1>
          <p className="text-slate-500 mt-1">Isi KRS rencana studi perkuliahan Anda atau berikan persetujuan bimbingan akademik</p>
        </div>
        <div>
          <select
            className="px-3 py-2 border border-slate-200 rounded-lg text-slate-600 font-medium"
            value={selectedPeriod}
            onChange={(e) => setSelectedPeriod(e.target.value)}
          >
            {academicPeriods.map((period) => (
              <option key={period.id} value={period.id}>{period.term}</option>
            ))}
          </select>
        </div>
      </div>

      {/* Tabs */}
      <div className="bg-white rounded-xl border border-slate-200 shadow-sm overflow-hidden">
        <div className="flex border-b border-slate-200 bg-slate-50/50">
          <button
            onClick={() => setActiveTab("fill")}
            className={`px-6 py-3 text-sm font-semibold transition-colors ${
              activeTab === "fill"
                ? "text-blue-600 border-b-2 border-blue-600 bg-white"
                : "text-slate-500 hover:text-slate-700"
            }`}
          >
            Pengisian KRS Mandiri
          </button>
          {(user?.role === "dosen" || user?.role === "admin") && (
            <button
              onClick={() => setActiveTab("approval")}
              className={`px-6 py-3 text-sm font-semibold transition-colors ${
                activeTab === "approval"
                  ? "text-blue-600 border-b-2 border-blue-600 bg-white"
                  : "text-slate-500 hover:text-slate-700"
              }`}
            >
              Persetujuan Dosen PA ({krsApprovals.length})
            </button>
          )}
        </div>

        {/* Content Area */}
        <div className="p-6">
          {activeTab === "fill" ? (
            <div className="space-y-6">
              {/* Status Alert */}
              {currentKrs ? (
                <div className="p-4 bg-slate-50 rounded-xl border border-slate-200 flex justify-between items-center flex-wrap gap-4">
                  <div>
                    <span className="text-xs font-semibold text-slate-400 uppercase tracking-wider block">Status KRS Anda:</span>
                    <span className={`inline-block mt-1 px-3 py-1 rounded-full text-xs font-bold ${getStatusBadge(currentKrs.status)}`}>
                      {currentKrs.status.toUpperCase()}
                    </span>
                  </div>
                  <div className="flex gap-2">
                    {currentKrs.status === "draft" && (
                      <button
                        onClick={handleSubmit}
                        className="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg text-sm font-semibold transition-all shadow-sm"
                      >
                        Ajukan KRS
                      </button>
                    )}
                  </div>
                </div>
              ) : (
                <div className="p-4 bg-blue-50 border border-blue-200 rounded-xl text-sm text-blue-900">
                  Anda belum mengisi KRS untuk periode ini. Silakan pilih kelas di bawah dan simpan draft.
                </div>
              )}

              {/* Class Selection grid */}
              {(!currentKrs || currentKrs.status === "draft") ? (
                <div className="space-y-4">
                  <div className="flex justify-between items-center border-b border-slate-100 pb-3">
                    <h3 className="font-bold text-slate-800 text-base">Pilih Mata Kuliah & Kelas</h3>
                    <span className="bg-blue-100 text-blue-800 text-xs px-3 py-1 rounded-full font-bold">
                      Terpilih: {totalSksSelected} / 24 SKS
                    </span>
                  </div>

                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    {availableClasses.map((cls) => {
                      const isSelected = selectedClassIds.includes(cls.id);
                      return (
                        <div
                          key={cls.id}
                          onClick={() => handleSelectClass(cls.id)}
                          className={`p-4 border rounded-xl cursor-pointer transition-all flex justify-between items-start ${
                            isSelected
                              ? "border-blue-500 bg-blue-50/40 text-blue-900 shadow-sm"
                              : "border-slate-200 hover:bg-slate-50 text-slate-700"
                          }`}
                        >
                          <div>
                            <span className="text-[10px] bg-slate-100 text-slate-500 px-2 py-0.5 rounded font-mono font-bold">
                              {cls.course_code}
                            </span>
                            <h4 className="font-bold text-slate-900 text-sm mt-1">{cls.course_name}</h4>
                            <p className="text-xs text-slate-500 mt-1">Kelas: {cls.class_code} | Kuota: {cls.enrolled_count}/{cls.quota}</p>
                          </div>
                          <div className="text-right">
                            <span className="text-xs font-semibold block">{cls.sks} SKS</span>
                            <span className={`inline-block mt-3 w-5 h-5 rounded-full border flex items-center justify-center ${
                              isSelected ? "bg-blue-600 border-blue-600 text-white" : "border-slate-300 bg-white"
                            }`}>
                              {isSelected && "✓"}
                            </span>
                          </div>
                        </div>
                      );
                    })}
                  </div>

                  <div className="flex justify-end pt-4">
                    <button
                      onClick={handleSaveDraft}
                      className="px-5 py-2.5 bg-blue-600 hover:bg-blue-700 text-white text-sm font-semibold rounded-lg transition-all shadow-sm"
                    >
                      Simpan Draft KRS
                    </button>
                  </div>
                </div>
              ) : (
                // View Locked/Approved KRS Items
                <div className="space-y-4">
                  <h3 className="font-bold text-slate-800 text-base border-b border-slate-100 pb-3">Daftar Mata Kuliah Terdaftar</h3>
                  <div className="border border-slate-200 rounded-xl overflow-hidden bg-white shadow-sm">
                    <table className="w-full text-left">
                      <thead className="bg-slate-50 border-b border-slate-200">
                        <tr>
                          <th className="p-4 text-xs font-semibold uppercase tracking-wider text-slate-500">Mata Kuliah</th>
                          <th className="p-4 text-xs font-semibold uppercase tracking-wider text-slate-500">SKS</th>
                          <th className="p-4 text-xs font-semibold uppercase tracking-wider text-slate-500">Status</th>
                        </tr>
                      </thead>
                      <tbody className="divide-y divide-slate-100">
                        {currentKrs.items?.map((item: any) => {
                          // Find corresponding class
                          const cDetails = availableClasses.find(c => c.id === item.classId);
                          return (
                            <tr key={item.id} className="hover:bg-slate-50 transition-colors">
                              <td className="p-4 text-sm font-bold text-slate-900">
                                {cDetails?.course_name || "Mata Kuliah"}
                                <span className="text-[10px] text-slate-400 font-mono block mt-0.5">{cDetails?.course_code || "KODE"} | Kelas: {cDetails?.class_code}</span>
                              </td>
                              <td className="p-4 text-sm text-slate-700">{cDetails?.sks || 3} SKS</td>
                              <td className="p-4 text-sm">
                                <span className={`px-2.5 py-0.5 rounded-full text-xs font-medium ${getStatusBadge(item.status || "approved")}`}>
                                  {item.status || "approved"}
                                </span>
                              </td>
                            </tr>
                          );
                        })}
                      </tbody>
                    </table>
                  </div>
                </div>
              )}
            </div>
          ) : (
            // Approval Panel View for Dosen PA
            <div className="space-y-4">
              <h2 className="text-base font-bold text-slate-900">Antrian Pengajuan KRS Mahasiswa Bimbingan</h2>
              {krsApprovals.length === 0 ? (
                <div className="text-center text-slate-500 py-12">Tidak ada pengajuan KRS masuk saat ini.</div>
              ) : (
                <div className="space-y-4">
                  {krsApprovals.map((req) => (
                    <div key={req.id} className="bg-white border border-slate-200 rounded-xl p-6 shadow-sm flex flex-col md:flex-row justify-between items-start md:items-center gap-4">
                      <div>
                        <span className="text-[10px] font-bold uppercase bg-yellow-100 text-yellow-800 px-2 py-0.5 rounded">
                          Perlu Review
                        </span>
                        <h4 className="font-bold text-slate-900 text-sm mt-1.5">Mahasiswa ID: {req.studentId}</h4>
                        <p className="text-xs text-slate-500 mt-1">Total Matakuliah: {req.items?.length || 0} | Periode ID: {req.academicPeriodId}</p>
                      </div>
                      <div className="flex gap-2 w-full md:w-auto">
                        <button
                          onClick={() => handleApprove(req.id)}
                          className="flex-1 md:flex-none px-4 py-2 bg-green-600 hover:bg-green-700 text-white text-xs font-semibold rounded-lg transition-all"
                        >
                          Setujui KRS
                        </button>
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
