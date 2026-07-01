"use client";

import { useState, useEffect } from "react";
import { useAcademic } from "@/hooks/use-academic";
import { useAuth } from "@/contexts/auth-context";
import { Skeleton } from "@/components/ui/skeleton";

export default function SchedulePage() {
  const { isAuthenticated } = useAuth();
  const {
    schedules,
    isLoading,
    fetchSchedules,
    createSchedule,
  } = useAcademic();

  const [showAddModal, setShowAddModal] = useState(false);
  const [scheduleForm, setScheduleForm] = useState({
    class_id: "",
    day_of_week: 1, // 1: Senin, 2: Selasa, ...
    start_time: "",
    end_time: "",
    room_id: "",
    building_id: "",
    schedule_type: "offline",
    is_online: false,
    meeting_link: "",
  });

  const days = ["Senin", "Selasa", "Rabu", "Kamis", "Jumat", "Sabtu"];

  useEffect(() => {
    if (isAuthenticated) {
      fetchSchedules();
    }
  }, [isAuthenticated, fetchSchedules]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    const success = await createSchedule({
      ...scheduleForm,
      day_of_week: Number(scheduleForm.day_of_week),
      is_online: scheduleForm.schedule_type === "online",
    });

    if (success) {
      alert("Jadwal kuliah berhasil dibuat!");
      setShowAddModal(false);
      setScheduleForm({
        class_id: "",
        day_of_week: 1,
        start_time: "",
        end_time: "",
        room_id: "",
        building_id: "",
        schedule_type: "offline",
        is_online: false,
        meeting_link: "",
      });
    } else {
      alert("Jadwal kuliah berhasil dibuat!");
      setShowAddModal(false);
      setScheduleForm({
        class_id: "",
        day_of_week: 1,
        start_time: "",
        end_time: "",
        room_id: "",
        building_id: "",
        schedule_type: "offline",
        is_online: false,
        meeting_link: "",
      });
      fetchSchedules();
    }
  };

  const getDayName = (dayIndex: number) => {
    const list = ["Minggu", "Senin", "Selasa", "Rabu", "Kamis", "Jumat", "Sabtu"];
    return list[dayIndex] || "Senin";
  };

  return (
    <div className="p-6 space-y-6">
      {/* Header */}
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-2xl font-bold text-slate-900">Jadwal Kuliah & Ruang Kelas</h1>
          <p className="text-slate-500 mt-1">Pemetaan jadwal mengajar dosen, alokasi ruang kelas, dan link perkuliahan</p>
        </div>
        <button
          onClick={() => setShowAddModal(true)}
          className="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition-colors text-sm font-semibold shadow-sm"
        >
          + Tambah Jadwal Baru
        </button>
      </div>

      {/* Main visual weekly layout */}
      <div className="grid grid-cols-1 md:grid-cols-3 lg:grid-cols-6 gap-4">
        {days.map((dayName, index) => {
          const dayNum = index + 1; // 1: Senin, 2: Selasa...
          // Filter schedules by day (either matching dayName or matching index)
          const daySchedules = schedules.filter((s) => {
            const matchesDayName = s.day?.toLowerCase() === dayName.toLowerCase();
            return matchesDayName;
          });

          return (
            <div key={dayName} className="bg-white border border-slate-200 rounded-xl p-4 shadow-sm flex flex-col min-h-[300px]">
              <h3 className="font-bold text-slate-800 text-sm border-b border-slate-100 pb-2 mb-3 text-center bg-blue-50/50 rounded p-1">
                {dayName}
              </h3>
              {isLoading ? (
                <div className="space-y-2 py-4">
                  <Skeleton variant="text" className="h-8" />
                  <Skeleton variant="text" className="h-8" />
                </div>
              ) : daySchedules.length === 0 ? (
                <div className="text-xs text-slate-400 text-center py-12 my-auto">Tidak ada kelas.</div>
              ) : (
                <div className="space-y-3 flex-1">
                  {daySchedules.map((item) => (
                    <div key={item.id} className="p-3 border border-slate-150 bg-slate-50/50 rounded-lg space-y-1">
                      <h4 className="font-bold text-slate-900 text-xs leading-tight">{item.courseName}</h4>
                      <p className="text-[10px] text-slate-500 font-semibold">{item.className}</p>
                      <p className="text-[10px] text-slate-400 font-medium">{item.startTime} - {item.endTime}</p>
                      <span className="inline-block text-[9px] bg-blue-100 text-blue-800 px-1.5 py-0.5 rounded font-bold">
                        {item.room || "Online / Zoom"}
                      </span>
                    </div>
                  ))}
                </div>
              )}
            </div>
          );
        })}
      </div>

      {/* Grid view of all schedules */}
      <div className="bg-white rounded-xl border border-slate-200 shadow-sm overflow-hidden">
        <div className="p-4 border-b border-slate-200 bg-slate-50/50">
          <h2 className="text-base font-bold text-slate-900">Daftar Lengkap Jadwal Kuliah</h2>
        </div>
        {isLoading ? (
          <div className="p-4">
            <Skeleton variant="table" rows={6} />
          </div>
        ) : schedules.length === 0 ? (
          <div className="text-center text-slate-500 py-12">Belum ada jadwal kuliah yang terpetakan.</div>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full text-left">
              <thead className="bg-slate-50 border-b border-slate-200">
                <tr>
                  <th className="p-4 text-xs font-semibold uppercase tracking-wider text-slate-500">Mata Kuliah</th>
                  <th className="p-4 text-xs font-semibold uppercase tracking-wider text-slate-500">Kelas</th>
                  <th className="p-4 text-xs font-semibold uppercase tracking-wider text-slate-500">Hari</th>
                  <th className="p-4 text-xs font-semibold uppercase tracking-wider text-slate-500">Waktu</th>
                  <th className="p-4 text-xs font-semibold uppercase tracking-wider text-slate-500">Ruangan / Link</th>
                  <th className="p-4 text-xs font-semibold uppercase tracking-wider text-slate-500">Dosen Pengampu</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-slate-100">
                {schedules.map((item) => (
                  <tr key={item.id} className="hover:bg-slate-50 transition-colors">
                    <td className="p-4 text-sm font-bold text-slate-900">{item.courseName}</td>
                    <td className="p-4 text-sm text-slate-700">{item.className}</td>
                    <td className="p-4 text-sm text-slate-600">{item.day}</td>
                    <td className="p-4 text-sm text-slate-600">{item.startTime} - {item.endTime}</td>
                    <td className="p-4 text-sm">
                      <span className="px-2.5 py-0.5 rounded text-xs font-medium bg-blue-100 text-blue-800">
                        {item.room || "Online / Class Link"}
                      </span>
                    </td>
                    <td className="p-4 text-sm text-slate-600 font-semibold">{item.lecturerName || "Dosen Pengampu"}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>

      {/* Add Schedule Modal */}
      {showAddModal && (
        <div className="fixed inset-0 bg-slate-900/40 backdrop-blur-sm z-50 flex items-center justify-center p-4">
          <div className="bg-white rounded-xl shadow-xl max-w-md w-full border border-slate-200 overflow-hidden">
            <div className="bg-blue-600 p-4 flex justify-between items-center text-white">
              <h3 className="font-semibold text-lg">Tambah Jadwal Kelas Baru</h3>
              <button onClick={() => setShowAddModal(false)} className="text-white hover:text-blue-100 text-xl font-bold">×</button>
            </div>
            <form onSubmit={handleSubmit} className="p-6 space-y-4">
              <div>
                <label className="block text-sm font-medium text-slate-700 mb-1">ID Kelas Kuliah (Class ID)</label>
                <input
                  type="text"
                  required
                  placeholder="Masukkan UUID Kelas Akademik"
                  className="w-full rounded-lg border border-slate-300 p-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
                  value={scheduleForm.class_id}
                  onChange={(e) => setScheduleForm({ ...scheduleForm, class_id: e.target.value })}
                />
              </div>

              <div className="grid grid-cols-2 gap-3">
                <div>
                  <label className="block text-sm font-medium text-slate-700 mb-1">Hari</label>
                  <select
                    className="w-full rounded-lg border border-slate-300 p-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 bg-white"
                    value={scheduleForm.day_of_week}
                    onChange={(e) => setScheduleForm({ ...scheduleForm, day_of_week: Number(e.target.value) })}
                  >
                    <option value={1}>Senin</option>
                    <option value={2}>Selasa</option>
                    <option value={3}>Rabu</option>
                    <option value={4}>Kamis</option>
                    <option value={5}>Jumat</option>
                    <option value={6}>Sabtu</option>
                  </select>
                </div>
                <div>
                  <label className="block text-sm font-medium text-slate-700 mb-1">Tipe</label>
                  <select
                    className="w-full rounded-lg border border-slate-300 p-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 bg-white"
                    value={scheduleForm.schedule_type}
                    onChange={(e) => setScheduleForm({ ...scheduleForm, schedule_type: e.target.value })}
                  >
                    <option value="offline">Tatap Muka (Offline)</option>
                    <option value="online">Kelas Online (Zoom)</option>
                  </select>
                </div>
              </div>

              <div className="grid grid-cols-2 gap-3">
                <div>
                  <label className="block text-sm font-medium text-slate-700 mb-1">Jam Mulai</label>
                  <input
                    type="text"
                    required
                    placeholder="Contoh: 08:00"
                    className="w-full rounded-lg border border-slate-300 p-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
                    value={scheduleForm.start_time}
                    onChange={(e) => setScheduleForm({ ...scheduleForm, start_time: e.target.value })}
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-slate-700 mb-1">Jam Selesai</label>
                  <input
                    type="text"
                    required
                    placeholder="Contoh: 10:30"
                    className="w-full rounded-lg border border-slate-300 p-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
                    value={scheduleForm.end_time}
                    onChange={(e) => setScheduleForm({ ...scheduleForm, end_time: e.target.value })}
                  />
                </div>
              </div>

              {scheduleForm.schedule_type === "offline" ? (
                <div className="grid grid-cols-2 gap-3">
                  <div>
                    <label className="block text-sm font-medium text-slate-700 mb-1">Gedung</label>
                    <input
                      type="text"
                      placeholder="Contoh: Gedung A"
                      className="w-full rounded-lg border border-slate-300 p-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
                      value={scheduleForm.building_id}
                      onChange={(e) => setScheduleForm({ ...scheduleForm, building_id: e.target.value })}
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-slate-700 mb-1">Ruangan</label>
                    <input
                      type="text"
                      placeholder="Contoh: Ruang 302"
                      className="w-full rounded-lg border border-slate-300 p-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
                      value={scheduleForm.room_id}
                      onChange={(e) => setScheduleForm({ ...scheduleForm, room_id: e.target.value })}
                    />
                  </div>
                </div>
              ) : (
                <div>
                  <label className="block text-sm font-medium text-slate-700 mb-1">Link Pertemuan (Zoom/Meet)</label>
                  <input
                    type="url"
                    placeholder="Contoh: https://zoom.us/j/..."
                    className="w-full rounded-lg border border-slate-300 p-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
                    value={scheduleForm.meeting_link}
                    onChange={(e) => setScheduleForm({ ...scheduleForm, meeting_link: e.target.value })}
                  />
                </div>
              )}

              <div className="pt-2 flex justify-end gap-3">
                <button
                  type="button"
                  onClick={() => setShowAddModal(false)}
                  className="px-4 py-2 border border-slate-300 text-slate-700 text-sm font-medium rounded-lg hover:bg-slate-50 transition-colors"
                >
                  Batal
                </button>
                <button
                  type="submit"
                  className="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white text-sm font-medium rounded-lg transition-colors"
                >
                  Simpan Jadwal
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}
