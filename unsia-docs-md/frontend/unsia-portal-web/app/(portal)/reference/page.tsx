"use client";

import { useState, useEffect } from "react";
import { useReference } from "@/contexts/reference-context";
import { useAuth } from "@/contexts/auth-context";
import { Skeleton } from "@/components/ui/skeleton";

export default function ReferencePage() {
  const { isAuthenticated } = useAuth();
  const {
    studyPrograms,
    academicPeriods,
    pmbWaves,
    provinces,
    cities,
    districts,
    isLoading,
    fetchAll,
    fetchCities,
    fetchDistricts,
  } = useReference();

  const [activeTab, setActiveTab] = useState<"programs" | "periods" | "waves" | "regions">("programs");

  // Region drilldown states
  const [selectedProvince, setSelectedProvince] = useState("");
  const [selectedCity, setSelectedCity] = useState("");

  useEffect(() => {
    if (isAuthenticated) {
      fetchAll();
    }
  }, [isAuthenticated]);

  // Fetch cities when province selection changes
  const handleProvinceChange = (provinceId: string) => {
    setSelectedProvince(provinceId);
    setSelectedCity("");
    if (provinceId) {
      fetchCities(provinceId);
    }
  };

  // Fetch districts when city selection changes
  const handleCityChange = (cityId: string) => {
    setSelectedCity(cityId);
    if (cityId) {
      fetchDistricts(cityId);
    }
  };

  const getStatusBadge = (status: string | boolean) => {
    const isActive = typeof status === "boolean" ? status : status?.toLowerCase() === "active" || status?.toLowerCase() === "aktif";
    return isActive
      ? "bg-green-100 text-green-800 border border-green-200"
      : "bg-gray-100 text-gray-800 border border-gray-200";
  };

  return (
    <div className="p-6 space-y-6">
      {/* Header */}
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-2xl font-bold text-slate-900">Data Master & Referensi</h1>
          <p className="text-slate-500 mt-1">Pusat data referensi akademik, gelombang pendaftaran, dan data administratif wilayah</p>
        </div>
        <button
          onClick={fetchAll}
          className="px-4 py-2 border border-slate-300 hover:bg-slate-50 text-slate-700 rounded-lg transition-colors text-sm font-medium"
        >
          🔄 Refresh Data
        </button>
      </div>

      {/* Tabs */}
      <div className="bg-white rounded-xl border border-slate-200 shadow-sm overflow-hidden">
        <div className="flex border-b border-slate-200 bg-slate-50/50 overflow-x-auto scrollbar-none">
          <button
            onClick={() => setActiveTab("programs")}
            className={`px-6 py-3 text-sm font-semibold transition-colors whitespace-nowrap ${
              activeTab === "programs"
                ? "text-slate-800 border-b-2 border-slate-800 bg-white"
                : "text-slate-500 hover:text-slate-700"
            }`}
          >
            Program Studi
          </button>
          <button
            onClick={() => setActiveTab("periods")}
            className={`px-6 py-3 text-sm font-semibold transition-colors whitespace-nowrap ${
              activeTab === "periods"
                ? "text-slate-800 border-b-2 border-slate-800 bg-white"
                : "text-slate-500 hover:text-slate-700"
            }`}
          >
            Periode Akademik
          </button>
          <button
            onClick={() => setActiveTab("waves")}
            className={`px-6 py-3 text-sm font-semibold transition-colors whitespace-nowrap ${
              activeTab === "waves"
                ? "text-slate-800 border-b-2 border-slate-800 bg-white"
                : "text-slate-500 hover:text-slate-700"
            }`}
          >
            Gelombang PMB
          </button>
          <button
            onClick={() => setActiveTab("regions")}
            className={`px-6 py-3 text-sm font-semibold transition-colors whitespace-nowrap ${
              activeTab === "regions"
                ? "text-slate-800 border-b-2 border-slate-800 bg-white"
                : "text-slate-500 hover:text-slate-700"
            }`}
          >
            Wilayah Indonesia
          </button>
        </div>

        {/* Content Area */}
        <div className="p-6">
          {isLoading ? (
            <Skeleton variant="table" rows={5} />
          ) : activeTab === "programs" ? (
            studyPrograms.length === 0 ? (
              <div className="text-center text-slate-500 py-8">Tidak ada data program studi tersedia.</div>
            ) : (
              <div className="overflow-x-auto">
                <table className="w-full text-left">
                  <thead className="bg-slate-50 border-b border-slate-200">
                    <tr>
                      <th className="p-4 text-xs font-semibold uppercase tracking-wider text-slate-500">Kode</th>
                      <th className="p-4 text-xs font-semibold uppercase tracking-wider text-slate-500">Nama Program Studi</th>
                      <th className="p-4 text-xs font-semibold uppercase tracking-wider text-slate-500">Jenjang</th>
                      <th className="p-4 text-xs font-semibold uppercase tracking-wider text-slate-500">Status</th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-slate-100">
                    {studyPrograms.map((prog) => (
                      <tr key={prog.id} className="hover:bg-slate-50 transition-colors">
                        <td className="p-4 text-sm font-bold text-slate-900">{prog.code}</td>
                        <td className="p-4 text-sm text-slate-700">{prog.name}</td>
                        <td className="p-4 text-sm text-slate-600 font-semibold">{prog.degree || "S1"}</td>
                        <td className="p-4 text-sm">
                          <span className={`px-2.5 py-1 rounded-full text-xs font-medium ${getStatusBadge(prog.status)}`}>
                            {prog.status}
                          </span>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            )
          ) : activeTab === "periods" ? (
            academicPeriods.length === 0 ? (
              <div className="text-center text-slate-500 py-8">Tidak ada data periode akademik tersedia.</div>
            ) : (
              <div className="overflow-x-auto">
                <table className="w-full text-left">
                  <thead className="bg-slate-50 border-b border-slate-200">
                    <tr>
                      <th className="p-4 text-xs font-semibold uppercase tracking-wider text-slate-500">Kode Periode</th>
                      <th className="p-4 text-xs font-semibold uppercase tracking-wider text-slate-500">Term / Semester</th>
                      <th className="p-4 text-xs font-semibold uppercase tracking-wider text-slate-500">Tanggal Mulai</th>
                      <th className="p-4 text-xs font-semibold uppercase tracking-wider text-slate-500">Tanggal Selesai</th>
                      <th className="p-4 text-xs font-semibold uppercase tracking-wider text-slate-500">Status</th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-slate-100">
                    {academicPeriods.map((period) => (
                      <tr key={period.id} className="hover:bg-slate-50 transition-colors">
                        <td className="p-4 text-sm font-bold text-slate-900">{period.code}</td>
                        <td className="p-4 text-sm text-slate-700">{period.term}</td>
                        <td className="p-4 text-sm text-slate-600">
                          {new Date(period.startDate).toLocaleDateString("id-ID", { dateStyle: "medium" })}
                        </td>
                        <td className="p-4 text-sm text-slate-600">
                          {new Date(period.endDate).toLocaleDateString("id-ID", { dateStyle: "medium" })}
                        </td>
                        <td className="p-4 text-sm">
                          <span className={`px-2.5 py-1 rounded-full text-xs font-medium ${getStatusBadge(period.status)}`}>
                            {period.status}
                          </span>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            )
          ) : activeTab === "waves" ? (
            pmbWaves.length === 0 ? (
              <div className="text-center text-slate-500 py-8">Tidak ada data gelombang pendaftaran PMB tersedia.</div>
            ) : (
              <div className="overflow-x-auto">
                <table className="w-full text-left">
                  <thead className="bg-slate-50 border-b border-slate-200">
                    <tr>
                      <th className="p-4 text-xs font-semibold uppercase tracking-wider text-slate-500">Kode Gelombang</th>
                      <th className="p-4 text-xs font-semibold uppercase tracking-wider text-slate-500">Nama Gelombang</th>
                      <th className="p-4 text-xs font-semibold uppercase tracking-wider text-slate-500">Tgl Registrasi</th>
                      <th className="p-4 text-xs font-semibold uppercase tracking-wider text-slate-500">Status Gelombang</th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-slate-100">
                    {pmbWaves.map((wave) => (
                      <tr key={wave.id} className="hover:bg-slate-50 transition-colors">
                        <td className="p-4 text-sm font-bold text-slate-900">{wave.code}</td>
                        <td className="p-4 text-sm text-slate-700">{wave.name}</td>
                        <td className="p-4 text-sm text-slate-600">
                          {wave.registrationStartAt
                            ? `${new Date(wave.registrationStartAt).toLocaleDateString("id-ID", { dateStyle: "short" })} s/d ${new Date(wave.registrationEndAt || "").toLocaleDateString("id-ID", { dateStyle: "short" })}`
                            : "Belum Diatur"}
                        </td>
                        <td className="p-4 text-sm">
                          <span className={`px-2.5 py-1 rounded-full text-xs font-medium ${getStatusBadge(wave.isActive)}`}>
                            {wave.isActive ? "Aktif" : "Nonaktif"}
                          </span>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            )
          ) : (
            // Regions Drilldown View
            <div className="space-y-6">
              <div className="grid grid-cols-1 md:grid-cols-3 gap-4 bg-slate-50 p-4 rounded-xl border border-slate-200">
                {/* Province select */}
                <div>
                  <label className="block text-xs font-semibold text-slate-500 uppercase mb-1.5">Pilih Provinsi</label>
                  <select
                    className="w-full rounded-lg border border-slate-300 p-2.5 text-sm focus:outline-none focus:ring-2 focus:ring-slate-800 bg-white"
                    value={selectedProvince}
                    onChange={(e) => handleProvinceChange(e.target.value)}
                  >
                    <option value="">-- Pilih Provinsi --</option>
                    {provinces.map((prov) => (
                      <option key={prov.id} value={prov.id}>
                        {prov.name} ({prov.code})
                      </option>
                    ))}
                  </select>
                </div>

                {/* City select */}
                <div>
                  <label className="block text-xs font-semibold text-slate-500 uppercase mb-1.5">Pilih Kabupaten/Kota</label>
                  <select
                    className="w-full rounded-lg border border-slate-300 p-2.5 text-sm focus:outline-none focus:ring-2 focus:ring-slate-800 bg-white disabled:bg-slate-100 disabled:text-slate-400"
                    value={selectedCity}
                    disabled={!selectedProvince}
                    onChange={(e) => handleCityChange(e.target.value)}
                  >
                    <option value="">-- Pilih Kota --</option>
                    {cities.map((city) => (
                      <option key={city.id} value={city.id}>
                        {city.name} ({city.code})
                      </option>
                    ))}
                  </select>
                </div>

                {/* Info Text */}
                <div className="flex items-center text-xs text-slate-400 font-medium">
                  {selectedProvince && !selectedCity && <p>Silakan pilih kota untuk melihat daftar kecamatan.</p>}
                  {selectedCity && <p>Menampilkan kecamatan dari wilayah terpilih.</p>}
                  {!selectedProvince && <p>Mulai dengan memilih Provinsi terlebih dahulu.</p>}
                </div>
              </div>

              {/* Districts output */}
              {selectedCity && (
                <div className="border border-slate-200 rounded-xl overflow-hidden">
                  <div className="bg-slate-50 p-4 border-b border-slate-200">
                    <h4 className="text-sm font-bold text-slate-800">Daftar Kecamatan Terdaftar</h4>
                  </div>
                  {districts.length === 0 ? (
                    <div className="text-center py-6 text-slate-500 text-sm">Tidak ada data kecamatan.</div>
                  ) : (
                    <div className="grid grid-cols-2 md:grid-cols-4 gap-3 p-4">
                      {districts.map((dist) => (
                        <div key={dist.id} className="p-3 border border-slate-200 bg-white rounded-lg text-sm text-slate-700 hover:shadow-sm transition-all">
                          <p className="font-semibold text-slate-900">{dist.name}</p>
                          <p className="text-[10px] text-slate-400 mt-0.5">Kode: {dist.code}</p>
                        </div>
                      ))}
                    </div>
                  )}
                </div>
              )}
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
