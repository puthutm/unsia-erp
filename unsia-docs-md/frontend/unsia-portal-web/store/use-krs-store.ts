"use client";

import { create } from "zustand";

interface SelectedClass {
  id: string;
  courseCode: string;
  courseName: string;
  className: string;
  sks: number;
  lecturerName: string;
  scheduleTime: string;
}

interface KrsStore {
  selectedClasses: SelectedClass[];
  addClass: (cls: SelectedClass) => void;
  removeClass: (classId: string) => void;
  clearStore: () => void;
  totalSks: () => number;
}

export const useKrsStore = create<KrsStore>((set, get) => ({
  selectedClasses: [],
  addClass: (cls) => {
    const exists = get().selectedClasses.some((c) => c.id === cls.id);
    if (!exists) {
      set({ selectedClasses: [...get().selectedClasses, cls] });
    }
  },
  removeClass: (classId) => {
    set({
      selectedClasses: get().selectedClasses.filter((c) => c.id !== classId),
    });
  },
  clearStore: () => set({ selectedClasses: [] }),
  totalSks: () => {
    return get().selectedClasses.reduce((sum, c) => sum + c.sks, 0);
  },
}));
