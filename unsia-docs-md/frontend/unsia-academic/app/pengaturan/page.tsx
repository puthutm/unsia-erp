"use client";

import { useState } from "react";

// Pengaturan Page - Matches: UI/AKADEMIK/ADMIN/panel-pengaturan
// Pengaturan & Setting Akademik

export default function PengaturanPage() {
  const [activeTab, setActiveTab] = useState("nim");

  // NIM Format (Drag & Drop - Simplified)
  const [nimFormat, setNimFormat] = useState([
    { id: "t1", label: "Tahun (2 digit)", example: "26", code: "YY" },
    { id: "t2", label: "Kode Prodi", example: "09", code: "PP" },
    { id: "t3", label: "Nomor Urut", example: "001", code: "NNN" },
  ]);

  // Bobot Nilai (Default)
  const [bobotNilai, setBobotNilai] = useState({
    tugas: 20,
    kuis: 10,
    uts: 30,
    uas: 40,
  });

  // Aturan Kehadiran
  const [kehadiran, setKehadiran] = useState({
    minKehadiran: 75,
    toleransiSakit: 3,
    toleransiIzin: 2,
  });

  // Hitung total bobot
  const totalBobot = Object.values(bobotNilai).reduce((a, b) => a + b, 0);

  const tabs = [
    { id: "nim", label: "Generate NIM", icon: "#" },
    { id: "bobot", label: "Bobot Nilai", icon: "%" },
    { id: "hadir", label: "Aturan Kehadiran", icon: "✓" },
    { id: "lulus", label: "Syarat Kelulusan", icon: "🎓" },
  ];

  const handleSave = () => {
    alert("Pengaturan berhasil disimpan!");
  };

  return (
    <div className="p-6 lg:p-8 space-y-6">
      {/* Header */}
      <div className="bg-gradient-to-br from-brand-700 via-brand-600 to-brand-500 rounded-2xl p-5 text-white">
        <p className="text-[10px] uppercase tracking-widest font-bold text-brand-accent">Pengaturan</p>
        <h2 className="font-display font-black text-2xl mt-1">Pengaturan & Setting Akademik</h2>
        <p className="text-brand-50 text-sm mt-1.5">
          Konfigurasi format NIM, bobot komponen nilai, aturan kehadiran, syarat kelulusan.
        </p>
      </div>

      {/* Sub-section Navigation */}
      <div className="flex gap-1 mb-5 overflow-x-auto pb-1">
        {tabs.map((tab) => (
          <button
            key={tab.id}
            onClick={() => setActiveTab(tab.id)}
            className={`px-4 py-2 rounded-lg text-xs font-bold shrink-0 transition-colors ${
              activeTab === tab.id
                ? "bg-brand-600 text-white"
                : "bg-white border border-slate-200 text-slate-700 hover:bg-slate-50"
            }`}
          >
            <span className="mr-1">{tab.icon}</span>
            {tab.label}
          </button>
        ))}
      </div>

      {/* Tab Content: Generate NIM */}
      {activeTab === "nim" && (
        <div className="bg-white rounded-2xl border border-slate-200 shadow-soft p-5">
          <h3 className="font-bold text-slate-800 text-sm mb-2 flex items-center gap-2">
            <span className="w-8 h-8 rounded-lg bg-brand-100 text-brand-600 flex items-center justify-center">#</span>
            Generate NIM — Format Custom
          </h3>
          <p className="text-[11px] text-slate-500 mb-4">
            Susun komponen NIM dengan drag & drop. Klik komponen untuk reorder.
          </p>

          {/* Current Format */}
          <div className="bg-slate-50 rounded-xl p-4 mb-4">
            <p className="text-[10px] uppercase tracking-wider font-bold text-slate-500 mb-3">Format NIM Saat Ini</p>
            <div className="flex flex-wrap gap-2">
              {nimFormat.map((item, idx) => (
                <div
                  key={item.id}
                  className="px-3 py-2 bg-white border border-brand-300 rounded-lg flex items-center gap-2 cursor-move"
                >
                  <span className="font-mono font-bold text-brand-700">{item.code}</span>
                  <span className="text-[10px] text-slate-500">: {item.label}</span>
                </div>
              ))}
            </div>
            <div className="mt-4 p-3 bg-brand-50 border border-brand-200 rounded-lg">
              <p className="text-[10px] uppercase tracking-wider font-bold text-brand-700">Preview NIM</p>
              <p className="font-mono font-black text-2xl text-brand-700 mt-1">26-09-001</p>
            </div>
          </div>

          {/* Available Components */}
          <div className="bg-white border border-slate-200 rounded-xl p-4">
            <p className="text-[10px] uppercase tracking-wider font-bold text-slate-500 mb-3">Komponen Tersedia</p>
            <div className="grid grid-cols-2 md:grid-cols-4 gap-2">
              {[
                { label: "Tahun (2 digit)", code: "YY", example: "26" },
                { label: "Tahun (4 digit)", code: "YYYY", example: "2026" },
                { label: "Kode Prodi", code: "PP", example: "09" },
                { label: "Kode Fakultas", code: "F", example: "1" },
                { label: "Gelombang", code: "G", example: "0" },
                { label: "Nomor Urut", code: "NNN", example: "001" },
              ].map((item) => (
                <button
                  key={item.code}
                  className="px-3 py-2 bg-slate-50 hover:bg-brand-50 border border-slate-200 hover:border-brand-300 rounded-lg text-left text-xs"
                >
                  <span className="font-mono font-bold text-slate-700">{item.code}</span>
                  <span className="block text-[10px] text-slate-500">{item.label}</span>
                </button>
              ))}
            </div>
          </div>

          <button
            onClick={handleSave}
            className="mt-4 w-full px-3 py-2.5 bg-brand-600 hover:bg-brand-700 text-white text-xs font-bold rounded-lg flex items-center justify-center gap-2"
          >
            <span>💾</span> Simpan Format NIM
          </button>
        </div>
      )}

      {/* Tab Content: Bobot Nilai */}
      {activeTab === "bobot" && (
        <div className="bg-white rounded-2xl border border-slate-200 shadow-soft p-5">
          <h3 className="font-bold text-slate-800 text-sm mb-2 flex items-center gap-2">
            <span className="w-8 h-8 rounded-lg bg-brand-100 text-brand-600 flex items-center justify-center">%</span>
            Bobot Komponen Nilai
          </h3>
          <p className="text-[11px] text-slate-500 mb-4">
            Bobot komponen nilai (Tugas, Kuis, UTS, UAS) berlaku global untuk semua kelas. Total harus = 100%.
          </p>

          <div className="space-y-3">
            {[
              { key: "tugas", label: "Tugas" },
              { key: "kuis", label: "Kuis" },
              { key: "uts", label: "UTS" },
              { key: "uas", label: "UAS" },
            ].map((item) => (
              <div key={item.key} className="flex items-center justify-between p-3 bg-slate-50 rounded-lg">
                <span className="text-xs font-bold text-slate-700">{item.label}</span>
                <div className="flex items-center gap-2">
                  <input
                    type="number"
                    value={bobotNilai[item.key as keyof typeof bobotNilai]}
                    onChange={(e) =>
                      setBobotNilai({
                        ...bobotNilai,
                        [item.key]: parseInt(e.target.value) || 0,
                      })
                    }
                    className="w-20 px-3 py-1.5 bg-white border border-slate-200 rounded-lg text-xs font-mono font-bold text-center"
                  />
                  <span className="text-sm font-bold text-slate-600">%</span>
                </div>
              </div>
            ))}
          </div>

          <div className="mt-4 p-3 bg-slate-50 border border-slate-200 rounded-lg flex items-center justify-between">
            <span className="text-xs font-bold text-slate-700">Total Bobot</span>
            <span className={`font-mono font-black text-xl ${totalBobot === 100 ? "text-emerald-600" : "text-rose-600"}`}>
              {totalBobot}%
            </span>
          </div>

          <button
            onClick={handleSave}
            className="mt-4 w-full px-3 py-2.5 bg-brand-600 hover:bg-brand-700 text-white text-xs font-bold rounded-lg flex items-center justify-center gap-2"
          >
            <span>✓</span> Set Bobot Nilai
          </button>
        </div>
      )}

      {/* Tab Content: Aturan Kehadiran */}
      {activeTab === "hadir" && (
        <div className="bg-white rounded-2xl border border-slate-200 shadow-soft p-5">
          <h3 className="font-bold text-slate-800 text-sm mb-2 flex items-center gap-2">
            <span className="w-8 h-8 rounded-lg bg-brand-100 text-brand-600 flex items-center justify-center">✓</span>
            Aturan Kehadiran
          </h3>
          <p className="text-[11px] text-slate-500 mb-4">
            Aturan kehadiran berlaku untuk seluruh mahasiswa aktif di LMS UNSIA.
          </p>

          <div className="space-y-3">
            <div className="flex items-center justify-between p-3 bg-slate-50 rounded-lg">
              <div>
                <p className="text-xs font-bold text-slate-700">Minimum Kehadiran</p>
                <p className="text-[10px] text-slate-500">Persentase minimal untuk lulus MK</p>
              </div>
              <div className="flex items-center gap-2">
                <input
                  type="number"
                  value={kehadiran.minKehadiran}
                  onChange={(e) => setKehadiran({ ...kehadiran, minKehadiran: parseInt(e.target.value) || 0 })}
                  className="w-20 px-3 py-1.5 bg-white border border-slate-200 rounded-lg text-xs font-mono font-bold text-center"
                />
                <span className="text-sm font-bold text-slate-600">%</span>
              </div>
            </div>

            <div className="flex items-center justify-between p-3 bg-slate-50 rounded-lg">
              <div>
                <p className="text-xs font-bold text-slate-700">Toleransi Sakit (dengan surat)</p>
                <p className="text-[10px] text-slate-500">Maks sesi yang dianggap hadir</p>
              </div>
              <input
                type="number"
                value={kehadiran.toleransiSakit}
                onChange={(e) => setKehadiran({ ...kehadiran, toleransiSakit: parseInt(e.target.value) || 0 })}
                className="w-20 px-3 py-1.5 bg-white border border-slate-200 rounded-lg text-xs font-mono font-bold text-center"
              />
            </div>

            <div className="flex items-center justify-between p-3 bg-slate-50 rounded-lg">
              <div>
                <p className="text-xs font-bold text-slate-700">Toleransi Izin (resmi)</p>
                <p className="text-[10px] text-slate-500">Maks sesi yang dianggap hadir</p>
              </div>
              <input
                type="number"
                value={kehadiran.toleransiIzin}
                onChange={(e) => setKehadiran({ ...kehadiran, toleransiIzin: parseInt(e.target.value) || 0 })}
                className="w-20 px-3 py-1.5 bg-white border border-slate-200 rounded-lg text-xs font-mono font-bold text-center"
              />
            </div>

            <div className="p-3 bg-rose-50 border border-rose-200 rounded-lg">
              <p className="text-xs font-bold text-rose-900">Sanksi Kehadiran &lt; 75%</p>
              <p className="text-[11px] text-rose-700 mt-1">
                Mahasiswa tidak boleh ikut UAS · nilai akhir = DO MK · wajib mengulang semester depan
              </p>
            </div>
          </div>

          <button
            onClick={handleSave}
            className="mt-4 w-full px-3 py-2.5 bg-brand-600 hover:bg-brand-700 text-white text-xs font-bold rounded-lg flex items-center justify-center gap-2"
          >
            <span>💾</span> Simpan Aturan Kehadiran
          </button>
        </div>
      )}

      {/* Tab Content: Syarat Kelulusan */}
      {activeTab === "lulus" && (
        <div className="bg-white rounded-2xl border border-slate-200 shadow-soft p-5">
          <h3 className="font-bold text-slate-800 text-sm mb-2 flex items-center gap-2">
            <span className="w-8 h-8 rounded-lg bg-brand-100 text-brand-600 flex items-center justify-center">🎓</span>
            Syarat Kelulusan
          </h3>
          <p className="text-[11px] text-slate-500 mb-4">
            Konfigurasi syarat minimum kelulusan per jenjang (S1, S2).
          </p>

          {/* Jenjang Tabs */}
          <div className="flex gap-2 mb-4">
            <button className="px-4 py-2 rounded-lg text-xs font-bold bg-brand-600 text-white">Sarjana (S1)</button>
            <button className="px-4 py-2 rounded-lg text-xs font-bold bg-white border border-slate-200 text-slate-700">Magister (S2)</button>
          </div>

          <div className="space-y-3">
            <div className="flex items-center justify-between p-3 bg-slate-50 rounded-lg">
              <span className="text-xs font-bold text-slate-700">Minimum SKS</span>
              <span className="text-xs font-mono font-bold text-slate-800">144</span>
            </div>
            <div className="flex items-center justify-between p-3 bg-slate-50 rounded-lg">
              <span className="text-xs font-bold text-slate-700">Minimum IPK</span>
              <span className="text-xs font-mono font-bold text-slate-800">2.75</span>
            </div>
            <div className="p-3 bg-emerald-50 border border-emerald-200 rounded-lg">
              <p className="text-[10px] uppercase tracking-wider font-bold text-emerald-700 mb-2">Syarat Tambahan</p>
              <ul className="text-[11px] text-emerald-800 space-y-1">
                <li>✓ Lulus semua MK Wajib Universitas (MKWU) - 18 SKS</li>
                <li>✓ Lulus Tugas Akhir / Skripsi (6 SKS)</li>
                <li>✓ Lulus seminar proposal & sidang skripsi</li>
                <li>✓ Bebas pinjaman perpustakaan</li>
                <li>✓ Bebas tanggungan keuangan</li>
                <li>✓ TOEFL min 450 / TOEIC min 600 / TOEP min 35</li>
              </ul>
            </div>
          </div>

          <button
            onClick={handleSave}
            className="mt-4 w-full px-3 py-2.5 bg-brand-600 hover:bg-brand-700 text-white text-xs font-bold rounded-lg flex items-center justify-center gap-2"
          >
            <span>💾</span> Simpan Syarat Kelulusan
          </button>
        </div>
      )}
    </div>
  );
}
