"use client";

import { useState, useEffect } from "react";
import { use } from "react";
import { useLMS, Session, Material, Assignment, Discussion } from "../../../hooks/use-lms";

interface SessionPageProps {
  params: Promise<{ id: string }>;
}

export default function SessionPage({ params }: SessionPageProps) {
  const resolvedParams = use(params);
  const { fetchSessions, fetchMaterials, fetchAssignments, fetchDiscussions, createDiscussion, replyToDiscussion } = useLMS();
  
  const [session, setSession] = useState<Session | null>(null);
  const [materials, setMaterials] = useState<Material[]>([]);
  const [assignments, setAssignments] = useState<Assignment[]>([]);
  const [discussions, setDiscussions] = useState<Discussion[]>([]);
  const [loading, setLoading] = useState(true);
  const [activeTab, setActiveTab] = useState<"materials" | "assignments" | "discussions">("materials");
  const [newDiscussion, setNewDiscussion] = useState("");

  useEffect(() => {
    loadSessionData();
  }, [resolvedParams.id]);

  const loadSessionData = async () => {
    setLoading(true);
    try {
      const sessions = await fetchSessions();
      const foundSession = sessions.find((s: Session) => s.id === resolvedParams.id);
      setSession(foundSession || null);

      if (foundSession) {
        const mats = await fetchMaterials(resolvedParams.id);
        setMaterials(mats);
        
        const assigns = await fetchAssignments(resolvedParams.id);
        setAssignments(assigns);
        
        const discs = await fetchDiscussions(resolvedParams.id);
        setDiscussions(discs);
      }
    } catch (error) {
      console.error("Error loading session data:", error);
    } finally {
      setLoading(false);
    }
  };

  const handleCreateDiscussion = async () => {
    if (!newDiscussion.trim()) return;
    const success = await createDiscussion(resolvedParams.id, newDiscussion);
    if (success) {
      setNewDiscussion("");
      loadSessionData();
    }
  };

  const getStatusBadge = (status: string) => {
    const styles: Record<string, string> = {
      upcoming: "bg-blue-100 text-blue-800",
      ongoing: "bg-green-100 text-green-800",
      completed: "bg-gray-100 text-gray-800",
    };
    return styles[status] || "bg-gray-100 text-gray-800";
  };

  if (loading) {
    return (
      <div className="p-6 flex items-center justify-center min-h-[400px]">
        <div className="text-slate-500">Memuat data...</div>
      </div>
    );
  }

  if (!session) {
    return (
      <div className="p-6">
        <div className="text-center text-slate-500 py-8">Sesi tidak ditemukan</div>
      </div>
    );
  }

  return (
    <div className="p-6 space-y-6">
      {/* Header */}
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-2xl font-bold text-slate-900">{session.title}</h1>
          <p className="text-slate-500 mt-1">{session.description}</p>
        </div>
        <span className={`px-3 py-1 rounded-full text-sm font-medium ${getStatusBadge(session.status)}`}>
          {session.status}
        </span>
      </div>

      {/* Session Info */}
      <div className="bg-white rounded-xl p-6 border border-slate-200">
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
          <div>
            <h3 className="text-sm font-medium text-slate-500">Jadwal</h3>
            <p className="text-slate-900 mt-1">{new Date(session.scheduledAt).toLocaleString()}</p>
          </div>
          <div>
            <h3 className="text-sm font-medium text-slate-500">Durasi</h3>
            <p className="text-slate-900 mt-1">{session.duration} menit</p>
          </div>
          <div>
            <h3 className="text-sm font-medium text-slate-500">Materi</h3>
            <p className="text-slate-900 mt-1">{materials.length} file</p>
          </div>
          <div>
            <h3 className="text-sm font-medium text-slate-500">Tugas</h3>
            <p className="text-slate-900 mt-1">{assignments.length} tugas</p>
          </div>
        </div>
      </div>

      {/* Tabs */}
      <div className="bg-white rounded-xl border border-slate-200">
        <div className="flex border-b border-slate-200">
          <button
            onClick={() => setActiveTab("materials")}
            className={`px-6 py-3 text-sm font-medium transition-colors ${
              activeTab === "materials"
                ? "text-blue-600 border-b-2 border-blue-600"
                : "text-slate-500 hover:text-slate-700"
            }`}
          >
            Materi ({materials.length})
          </button>
          <button
            onClick={() => setActiveTab("assignments")}
            className={`px-6 py-3 text-sm font-medium transition-colors ${
              activeTab === "assignments"
                ? "text-blue-600 border-b-2 border-blue-600"
                : "text-slate-500 hover:text-slate-700"
            }`}
          >
            Tugas ({assignments.length})
          </button>
          <button
            onClick={() => setActiveTab("discussions")}
            className={`px-6 py-3 text-sm font-medium transition-colors ${
              activeTab === "discussions"
                ? "text-blue-600 border-b-2 border-blue-600"
                : "text-slate-500 hover:text-slate-700"
            }`}
          >
            Diskusi ({discussions.length})
          </button>
        </div>

        {/* Tab Content */}
        <div className="p-6">
          {activeTab === "materials" && (
            <div className="space-y-3">
              {materials.length === 0 ? (
                <div className="text-center text-slate-500 py-8">Tidak ada materi</div>
              ) : (
                materials.map((material) => (
                  <div key={material.id} className="p-4 border border-slate-200 rounded-lg flex justify-between items-center">
                    <div>
                      <h4 className="font-medium text-slate-900">{material.title}</h4>
                      <p className="text-sm text-slate-500">{material.type} - {material.description}</p>
                    </div>
                    <a
                      href={material.url}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="px-4 py-2 bg-blue-600 text-white rounded-lg text-sm hover:bg-blue-700"
                    >
                      Buka
                    </a>
                  </div>
                ))
              )}
            </div>
          )}

          {activeTab === "assignments" && (
            <div className="space-y-3">
              {assignments.length === 0 ? (
                <div className="text-center text-slate-500 py-8">Tidak ada tugas</div>
              ) : (
                assignments.map((assignment) => (
                  <div key={assignment.id} className="p-4 border border-slate-200 rounded-lg">
                    <h4 className="font-medium text-slate-900">{assignment.title}</h4>
                    <p className="text-sm text-slate-500 mt-1">{assignment.description}</p>
                    <div className="flex gap-4 mt-3 text-sm text-slate-500">
                      <span>Batas: {new Date(assignment.dueDate).toLocaleString()}</span>
                      <span>Skor max: {assignment.maxScore}</span>
                      <span>Dikumpulkan: {assignment.submissions}</span>
                    </div>
                  </div>
                ))
              )}
            </div>
          )}

          {activeTab === "discussions" && (
            <div className="space-y-4">
              {/* New Discussion Form */}
              <div className="flex gap-2">
                <input
                  type="text"
                  value={newDiscussion}
                  onChange={(e) => setNewDiscussion(e.target.value)}
                  placeholder="Tulis pertanyaan atau komentar..."
                  className="flex-1 px-4 py-2 border border-slate-200 rounded-lg"
                />
                <button
                  onClick={handleCreateDiscussion}
                  className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700"
                >
                  Kirim
                </button>
              </div>

              {/* Discussion List */}
              {discussions.length === 0 ? (
                <div className="text-center text-slate-500 py-8">Belum ada diskusi</div>
              ) : (
                discussions.map((discussion) => (
                  <div key={discussion.id} className="p-4 border border-slate-200 rounded-lg">
                    <div className="flex justify-between items-start">
                      <div>
                        <h4 className="font-medium text-slate-900">{discussion.userName}</h4>
                        <p className="text-slate-600 mt-1">{discussion.content}</p>
                        <p className="text-xs text-slate-400 mt-2">
                          {new Date(discussion.createdAt).toLocaleString()}
                        </p>
                      </div>
                    </div>
                    {/* Replies */}
                    {discussion.replies && discussion.replies.length > 0 && (
                      <div className="ml-4 mt-3 pl-4 border-l-2 border-slate-200 space-y-2">
                        {discussion.replies.map((reply) => (
                          <div key={reply.id} className="text-sm">
                            <span className="font-medium text-slate-700">{reply.userName}: </span>
                            <span className="text-slate-600">{reply.content}</span>
                          </div>
                        ))}
                      </div>
                    )}
                  </div>
                ))
              )}
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
