"use client";

import { cn } from "@/lib/utils";

interface SkeletonProps extends React.HTMLAttributes<HTMLDivElement> {
  variant?: "base" | "card" | "stats" | "table" | "text" | "circle";
  rows?: number;
}

export function Skeleton({
  className,
  variant = "base",
  rows = 5,
  ...props
}: SkeletonProps) {
  if (variant === "circle") {
    return (
      <div
        className={cn("animate-pulse rounded-full bg-slate-200 dark:bg-slate-700", className)}
        {...props}
      />
    );
  }

  if (variant === "card") {
    return (
      <div className={cn("bg-white rounded-xl shadow-sm border border-slate-200 p-6 space-y-4", className)} {...props}>
        <div className="animate-pulse h-4 w-1/3 bg-slate-200 rounded" />
        <div className="animate-pulse h-8 w-2/3 bg-slate-200 rounded" />
      </div>
    );
  }

  if (variant === "stats") {
    return (
      <div className={cn("grid grid-cols-1 md:grid-cols-4 gap-4", className)} {...props}>
        {Array.from({ length: 4 }).map((_, i) => (
          <div key={i} className="bg-white rounded-xl shadow-sm border border-slate-200 p-6 space-y-4">
            <div className="animate-pulse h-4 w-1/3 bg-slate-200 rounded" />
            <div className="animate-pulse h-8 w-2/3 bg-slate-200 rounded" />
          </div>
        ))}
      </div>
    );
  }

  if (variant === "table") {
    return (
      <div className={cn("overflow-x-auto", className)} {...props}>
        <table className="w-full">
          <thead className="bg-slate-50">
            <tr>
              <th className="p-4 text-left"><div className="animate-pulse h-4 w-20 bg-slate-200 rounded" /></th>
              <th className="p-4 text-left"><div className="animate-pulse h-4 w-28 bg-slate-200 rounded" /></th>
              <th className="p-4 text-left"><div className="animate-pulse h-4 w-24 bg-slate-200 rounded" /></th>
              <th className="p-4 text-left"><div className="animate-pulse h-4 w-16 bg-slate-200 rounded" /></th>
              <th className="p-4 text-left"><div className="animate-pulse h-4 w-16 bg-slate-200 rounded" /></th>
            </tr>
          </thead>
          <tbody>
            {Array.from({ length: rows }).map((_, i) => (
              <tr key={i} className="border-t border-slate-200">
                <td className="p-4"><div className="animate-pulse h-5 w-24 bg-slate-200 rounded" /></td>
                <td className="p-4"><div className="animate-pulse h-5 w-48 bg-slate-200 rounded" /></td>
                <td className="p-4"><div className="animate-pulse h-5 w-32 bg-slate-200 rounded" /></td>
                <td className="p-4"><div className="animate-pulse h-5 w-28 bg-slate-200 rounded" /></td>
                <td className="p-4"><div className="animate-pulse h-5 w-16 bg-slate-200 rounded" /></td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    );
  }

  // base / text
  return (
    <div
      className={cn(
        "animate-pulse rounded bg-slate-200 dark:bg-slate-700",
        variant === "text" ? "h-4 w-full" : "h-5 w-full",
        className
      )}
      {...props}
    />
  );
}
