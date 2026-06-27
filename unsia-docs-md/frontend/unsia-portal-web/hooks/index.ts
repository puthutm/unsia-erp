// Re-export hooks from contexts
export { useAuth } from "@/contexts/auth-context";
export { useReference } from "@/contexts/reference-context";

// Export module hooks
export { usePmb } from "./use-pmb";
export { useFinance } from "./use-finance";
export { useAcademic } from "./use-academic";
export { useLms } from "./use-lms";
export { useCRM } from "./use-crm";
export { useHRIS } from "./use-hris";
export { useAssessment } from "./use-assessment";

// Export type interfaces
export type { LmsCourse, LmsClass, Enrollment, Session, Material, Assignment } from "./use-lms";
export type { Student, Course, Schedule, KrsEntry, StudentGrade } from "./use-academic";
export type { Lead, Campaign, Agent } from "./use-crm";
export type { Employee, Attendance, LeaveRequest } from "./use-hris";
export type { CBTSession, CBTQuestion, CBTAttempt } from "./use-assessment";
