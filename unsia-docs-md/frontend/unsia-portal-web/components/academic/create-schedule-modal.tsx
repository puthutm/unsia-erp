"use client";

import { useState } from "react";
import { Modal } from "@/components/ui/modal";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Select } from "@/components/ui/select";
import { useAcademic } from "@/hooks/use-academic";
import { useLms, LmsClass } from "@/hooks/use-lms";

interface CreateScheduleModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSuccess: () => void;
  classes: LmsClass[];
}

export function CreateScheduleModal({ isOpen, onClose, onSuccess, classes }: CreateScheduleModalProps) {
  const { createSchedule, isLoading } = useAcademic();
  const [errors, setErrors] = useState<Record<string, string>>({});

  const [formData, setFormData] = useState({
    classId: "",
    dayOfWeek: "",
    startTime: "",
    endTime: "",
    roomId: "",
    scheduleType: "regular",
    isOnline: false,
    meetingLink: "",
  });

  const resetForm = () => {
    setFormData({
      classId: "",
      dayOfWeek: "",
      startTime: "",
      endTime: "",
      roomId: "",
      scheduleType: "regular",
      isOnline: false,
      meetingLink: "",
    });
    setErrors({});
  };

  const handleClose = () => {
    resetForm();
    onClose();
  };

  const validateForm = () => {
    const newErrors: Record<string, string> = {};

    if (!formData.classId) {
      newErrors.classId = "Kelas wajib dipilih";
    }
    if (!formData.dayOfWeek) {
      newErrors.dayOfWeek = "Hari wajib dipilih";
    }
    if (!formData.startTime) {
      newErrors.startTime = "Jam mulai wajib diisi";
    }
    if (!formData.endTime) {
      newErrors.endTime = "Jam selesai wajib diisi";
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = async () => {
    if (!validateForm()) return;

    const dayMap: Record<string, number> = {
      "senin": 1,
      "selasa": 2,
      "rabu": 3,
      "kamis": 4,
      "jumat": 5,
      "sabtu": 6,
      "minggu": 7,
    };

    const success = await createSchedule({
      class_id: formData.classId,
      day_of_week: dayMap[formData.dayOfWeek] || 1,
      start_time: formData.startTime,
      end_time: formData.endTime,
      room_id: formData.roomId || undefined,
      schedule_type: formData.scheduleType || "regular",
      is_online: formData.isOnline,
      meeting_link: formData.isOnline ? formData.meetingLink || undefined : undefined,
    });

    if (success) {
      handleClose();
      onSuccess();
    }
  };

  const classOptions = [
    { value: "", label: "Pilih Kelas" },
    ...classes.map((cls) => ({
      value: cls.id,
      label: `${cls.className} - ${cls.courseName}`,
    })),
  ];

  const dayOptions = [
    { value: "", label: "Pilih Hari" },
    { value: "senin", label: "Senin" },
    { value: "selasa", label: "Selasa" },
    { value: "rabu", label: "Rabu" },
    { value: "kamis", label: "Kamis" },
    { value: "jumat", label: "Jumat" },
    { value: "sabtu", label: "Sabtu" },
  ];

  const scheduleTypeOptions = [
    { value: "regular", label: "Regular" },
    { value: "practice", label: "Praktikum" },
    { value: "field_practice", label: "Praktik Lapangan" },
    { value: "online", label: "Online" },
  ];

  return (
    <Modal isOpen={isOpen} onClose={handleClose} title="Buat Jadwal Kuliah" size="lg">
      <div className="space-y-4">
        {/* Kelas */}
        <Select
          label="Kelas"
          options={classOptions}
          value={formData.classId}
          onChange={(value) => setFormData({ ...formData, classId: value })}
          error={errors.classId}
        />

        {/* Hari */}
        <Select
          label="Hari"
          options={dayOptions}
          value={formData.dayOfWeek}
          onChange={(value) => setFormData({ ...formData, dayOfWeek: value })}
          error={errors.dayOfWeek}
        />

        <div className="grid grid-cols-2 gap-4">
          {/* Jam Mulai */}
          <Input
            label="Jam Mulai"
            type="time"
            value={formData.startTime}
            onChange={(e) => setFormData({ ...formData, startTime: e.target.value })}
            error={errors.startTime}
          />

          {/* Jam Selesai */}
          <Input
            label="Jam Selesai"
            type="time"
            value={formData.endTime}
            onChange={(e) => setFormData({ ...formData, endTime: e.target.value })}
            error={errors.endTime}
          />
        </div>

        {/* Ruang */}
        <Input
          label="Ruang (Opsional)"
          placeholder="Contoh: R.101"
          value={formData.roomId}
          onChange={(e) => setFormData({ ...formData, roomId: e.target.value })}
        />

        {/* Tipe Jadwal */}
        <Select
          label="Tipe Jadwal"
          options={scheduleTypeOptions}
          value={formData.scheduleType}
          onChange={(value) => setFormData({ ...formData, scheduleType: value })}
        />

        {/* Online Toggle */}
        <div className="flex items-center gap-3">
          <input
            type="checkbox"
            id="isOnline"
            checked={formData.isOnline}
            onChange={(e) => setFormData({ ...formData, isOnline: e.target.checked })}
            className="w-4 h-4 text-blue-600"
          />
          <label htmlFor="isOnline" className="text-sm text-slate-700">
            Jadwal Online
          </label>
        </div>

        {/* Meeting Link (if online) */}
        {formData.isOnline && (
          <Input
            label="Link Meeting (Opsional)"
            placeholder="https://meet.google.com/..."
            value={formData.meetingLink}
            onChange={(e) => setFormData({ ...formData, meetingLink: e.target.value })}
          />
        )}

        {/* Actions */}
        <div className="flex justify-end gap-3 pt-4 border-t border-slate-200">
          <Button variant="outline" onClick={handleClose} disabled={isLoading}>
            Batal
          </Button>
          <Button onClick={handleSubmit} isLoading={isLoading}>
            Simpan
          </Button>
        </div>
      </div>
    </Modal>
  );
}
