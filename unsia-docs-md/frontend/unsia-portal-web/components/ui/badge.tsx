"use client";

import { forwardRef, type HTMLAttributes } from "react";
import { cn } from "@/lib/utils";

export interface BadgeProps extends HTMLAttributes<HTMLSpanElement> {
  variant?: "default" | "success" | "warning" | "danger" | "info" | "outline";
}

export const Badge = forwardRef<HTMLSpanElement, BadgeProps>(
  ({ className, variant = "default", ...props }, ref) => {
    const variants = {
      default: "bg-slate-100 text-slate-800",
      success: "bg-green-100 text-green-800",
      warning: "bg-yellow-100 text-yellow-800",
      danger: "bg-red-100 text-red-800",
      info: "bg-blue-100 text-blue-800",
      outline: "border border-slate-300 text-slate-700",
    };

    return (
      <span
        ref={ref}
        className={cn(
          "inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium",
          variants[variant],
          className
        )}
        {...props}
      />
    );
  }
);

Badge.displayName = "Badge";

export const StatusBadge = ({ status }: { status: string }) => {
  const statusMap: Record<string, { variant: BadgeProps["variant"]; label: string }> = {
    active: { variant: "success", label: "Aktif" },
    inactive: { variant: "default", label: "Nonaktif" },
    pending: { variant: "warning", label: "Menunggu" },
    paid: { variant: "success", label: "Lunas" },
    verified: { variant: "success", label: "Terverifikasi" },
    rejected: { variant: "danger", label: "Ditolak" },
    accepted: { variant: "success", label: "Diterima" },
    submitted: { variant: "info", label: "Terkirim" },
    draft: { variant: "default", label: "Draft" },
    published: { variant: "success", label: "Dipublikasikan" },
    archived: { variant: "default", label: "Diarsipkan" },
  };

  const config = statusMap[status] || { variant: "default", label: status };

  return (
    <Badge variant={config.variant}>
      {config.label}
    </Badge>
  );
};
