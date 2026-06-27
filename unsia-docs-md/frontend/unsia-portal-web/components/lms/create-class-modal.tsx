"use client";

import { useState } from "react";
import { Modal } from "@/components/ui/modal";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Select } from "@/components/ui/select";
import { useLms, LmsCourse } from "@/hooks/use-lms";

interface CreateClassModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSuccess: () => void;
  courses: LmsCourse[];
}

export function CreateClassModal({ isOpen, onClose, onSuccess, courses }: CreateClassModalProps) {
  const { createClass, isLoading } = useLms();
  const [errors, setErrors] = useState<Record<string, string>>({});

  const [formData, setFormData] = useState({
    courseId: "",
    academicClassId: "",
    classCode: "",
    semester: "",
    academicYear: "",
    maxStudents: 30,
    lecturerId: "",
  });

  const resetForm = () => {
    setFormData({
      courseId: "",
      academicClassId: "",
      classCode: "",
      semester: "",
      academicYear: "",
      maxStudents: 30,
      lecturerId: "",
    });
    setErrors({});
  };

  const handleClose = () => {
    resetForm();
    onClose();
  };

  const validateForm = () => {
    const newErrors: Record<string, string> = {};

    if (!formData.courseId) {
      newErrors.courseId = "Mata kuliah wajib dipilih";
    }
    if (!formData.academicClassId) {
      newErrors.academicClassId = "ID kelas akademik wajib diisi";
    }
    if (!formData.classCode) {
      newErrors.classCode = "Kode kelas wajib diisi";
    }
    if (!formData.semester) {
      newErrors.semester = "Semester wajib diisi";
    }
    if (!formData.academicYear) {
      newErrors.academicYear = "Tahun ajaran wajib diisi";
    }
    if (formData.maxStudents < 1) {
      newErrors.maxStudents = "Kapasitas minimal 1 mahasiswa";
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = async () => {
    if (!validateForm()) return;

    const success = await createClass({
      academic_class_id: formData.academicClassId,
      course_id: formData.courseId,
      class_code: formData.classCode,
      semester: formData.semester,
      academic_year: formData.academicYear,
      max_students: formData.maxStudents,
      lecturer_id: formData.lecturerId || undefined,
    });

    if (success) {
      handleClose();
      onSuccess();
    }
  };

  const courseOptions = [
    { value: "", label: "Pilih Mata Kuliah" },
    ...courses.map((course) => ({
      value: course.id,
      label: `${course.code} - ${course.name}`,
    })),
  ];

  const semesterOptions = [
    { value: "", label: "Pilih Semester" },
    ...Array.from({ length: 8 }, (_, i) => ({
      value: String(i + 1),
      label: `Semester ${i + 1}`,
    })),
  ];

  const currentYear = new Date().getFullYear();
  const academicYearOptions = [
    { value: "", label: "Pilih Tahun Ajaran" },
    { value: `${currentYear}/${currentYear + 1}`, label: `${currentYear}/${currentYear + 1}` },
    { value: `${currentYear + 1}/${currentYear + 2}`, label: `${currentYear + 1}/${currentYear + 2}` },
    { value: `${currentYear - 1}/${currentYear}`, label: `${currentYear - 1}/${currentYear}` },
  ];

  return (
    <Modal isOpen={isOpen} onClose={handleClose} title="Buat Kelas Baru" size="lg">
      <div className="space-y-4">
        <Select
          label="Mata Kuliah"
          options={courseOptions}
          value={formData.courseId}
          onChange={(value) => setFormData({ ...formData, courseId: value })}
          error={errors.courseId}
        />

        <Input
          label="ID Kelas Akademik"
          placeholder="Contoh: CLASS-2024-001"
          value={formData.academicClassId}
          onChange={(e) => setFormData({ ...formData, academicClassId: e.target.value })}
          error={errors.academicClassId}
        />

        <Input
          label="Kode Kelas"
          placeholder="Contoh: TI-2024-A"
          value={formData.classCode}
          onChange={(e) => setFormData({ ...formData, classCode: e.target.value })}
          error={errors.classCode}
        />

        <div className="grid grid-cols-2 gap-4">
          <Select
            label="Semester"
            options={semesterOptions}
            value={formData.semester}
            onChange={(value) => setFormData({ ...formData, semester: value })}
            error={errors.semester}
          />

          <Select
            label="Tahun Ajaran"
            options={academicYearOptions}
            value={formData.academicYear}
            onChange={(value) => setFormData({ ...formData, academicYear: value })}
            error={errors.academicYear}
          />
        </div>

        <Input
          label="Kapasitas Maksimal Mahasiswa"
          type="number"
          min={1}
          max={500}
          value={formData.maxStudents}
          onChange={(e) => setFormData({ ...formData, maxStudents: parseInt(e.target.value) || 0 })}
          error={errors.maxStudents}
        />

        <Input
          label="ID Dosen Pengampu (Opsional)"
          placeholder="Kosongkan jika belum ditentukan"
          value={formData.lecturerId}
          onChange={(e) => setFormData({ ...formData, lecturerId: e.target.value })}
        />

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
