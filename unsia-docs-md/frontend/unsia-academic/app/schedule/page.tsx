"use client";

import { useState, useEffect } from "react";

// Schedule Page - Next.js
// Matches: UI/AKADEMIK/SA/activity diagram

interface ScheduleItem {
  id: string;
  courseName: string;
  courseCode: string;
  className: string;
  day: string;
  startTime: string;
  endTime: string;
  room: string;
  building: string;
  lecturerName: string;
  isOnline: boolean;
  meetingLink?: string;
}

export default function SchedulePage() {
  const [schedules, setSchedules] = useState<ScheduleItem[]>([]);
  const [loading, setLoading] = useState(true);
  const [selectedDay, setSelectedDay] = useState("all");

  useEffect(() => {
    fetchSchedules();
  }, []);

  const fetchSchedules = async () => {
    setLoading(true);
    try {
      const token = localStorage.getItem("accessToken");
      const response = await fetch("/api/v1/academic/schedules", {
        headers: { Authorization: `Bearer ${token}` },
      });
      if (response.ok) {
        const data = await response.json();
        setSchedules(data.data || []);
      }
    } catch (error) {
      console.error("Error fetching schedules:", error);
    } finally {
      setLoading(false);
    }
  };

  const days = ["Senin", "Selasa", "Rabu", "Kamis", "Jumat", "Sabtu"];
  const filteredSchedules = selectedDay === "all" 
    ? schedules 
    : schedules.filter(s => s.day === selectedDay);

  const getDayColor = (day: string) => {
    if (day === "Senin") return "bg-blue-100 text-blue-800";
    if (day === "Selasa") return "bg-green-100 text-green-800";
    if (day === "Rabu") return "bg-yellow-100 text-yellow-800";
    if (day === "Kamis") return "bg-orange-100 text-orange-800";
    if (day === "Jumat") return "bg-purple-100 text-purple-800";
    return "bg-gray-100 text-gray-800";
  };

  return (
    <div className="p-6 space-y-6">
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-2xl font-bold text-slate-900">Jadwal Kuliah</h1>
          <p className="text-slate-500 mt-1">Jadwal mata kuliah minggu ini</p>
        </div>
        <select 
          value={selectedDay}
          onChange={(e) => setSelectedDay(e.target.value)}
          className="px-3 py-2 border border-slate-200 rounded-lg"
        >
          <option value="all">Semua Hari</option>
          {days.map(day => (
            <option key={day} value={day}>{day}</option>
          ))}
        </select>
      </div>

      {/* Day Tabs */}
      <div className="flex gap-2 flex-wrap">
        <button
          onClick={() => setSelectedDay("all")}
          className={`px-4 py-2 rounded-lg font-medium ${
            selectedDay === "all" 
              ? "bg-blue-600 text-white" 
              : "bg-slate-100 text-slate-600"
          }`}
        >
          Semua
        </button>
        {days.map(day => (
          <button
            key={day}
            onClick={() => setSelectedDay(day)}
            className={`px-4 py-2 rounded-lg font-medium ${
              selectedDay === day 
                ? "bg-blue-600 text-white" 
                : "bg-slate-100 text-slate-600"
            }`}
          >
            {day}
          </button>
        ))}
      </div>

      {/* Schedule Grid */}
      {loading ? (
        <div className="text-center text-slate-500 py-8">Memuat jadwal...</div>
      ) : filteredSchedules.length === 0 ? (
        <div className="text-center text-slate-500 py-8">Tidak ada jadwal</div>
      ) : (
        <div className="grid gap-4">
          {filteredSchedules.map((schedule) => (
            <div 
              key={schedule.id} 
              className="bg-white rounded-xl p-4 border border-slate-200 flex items-center justify-between"
            >
              <div className="flex items-center gap-4">
                <div className={`px-3 py-1 rounded-full text-sm font-medium ${getDayColor(schedule.day)}`}>
                  {schedule.day}
                </div>
                <div>
                  <h3 className="font-medium text-slate-900">{schedule.courseName}</h3>
                  <p className="text-sm text-slate-500">{schedule.courseCode} - {schedule.className}</p>
                </div>
              </div>
              <div className="text-right">
                <p className="text-slate-900 font-medium">
                  {schedule.startTime} - {schedule.endTime}
                </p>
                <p className="text-sm text-slate-500">
                  {schedule.isOnline ? "Online" : `${schedule.room} - ${schedule.building}`}
                </p>
              </div>
            </div>
          ))}
        </div>
      )}

      {/* Weekly View */}
      <div className="bg-white rounded-xl border border-slate-200 overflow-hidden">
        <div className="p-4 border-b border-slate-200">
          <h2 className="font-semibold text-slate-900">Tampilan Mingguan</h2>
        </div>
        <div className="grid grid-cols-6 divide-x divide-slate-200">
          {days.map(day => (
            <div key={day} className="p-4">
              <h3 className="font-medium text-slate-900 text-sm mb-2">{day}</h3>
              <div className="space-y-2">
                {schedules
                  .filter(s => s.day === day)
                  .map(s => (
                    <div key={s.id} className="text-xs p-2 bg-slate-50 rounded">
                      <p className="font-medium">{s.startTime}</p>
                      <p className="text-slate-600 truncate">{s.courseName}</p>
                    </div>
                  ))}
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
