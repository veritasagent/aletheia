"use client";

import { cn } from "@/lib/utils";
import type { ReactNode } from "react";

type BadgeVariant = "FALSE" | "MISLEADING" | "TRUE" | "UNVERIFIABLE" | "neutral";

interface BadgeProps {
  variant?: BadgeVariant;
  children: ReactNode;
  className?: string;
}

const stylesByVariant: Record<BadgeVariant, string> = {
  FALSE: "bg-[#FEE2E2] text-[#DC2626] border border-[#FECACA]",
  MISLEADING: "bg-[#FEF3C7] text-[#D97706] border border-[#FDE68A]",
  TRUE: "bg-[#D1FAE5] text-[#059669] border border-[#A7F3D0]",
  UNVERIFIABLE: "bg-[#DBEAFE] text-[#2563EB] border border-[#BFDBFE]",
  neutral: "bg-[#F3F4F6] text-[#6B7280] border border-[#E5E7EB]",
};

export default function Badge({ variant = "neutral", children, className }: BadgeProps) {
  return (
    <span
      className={cn(
        "inline-flex items-center rounded-full px-2 py-1 font-mono text-[8px] uppercase tracking-[0.22em]",
        stylesByVariant[variant],
        className,
      )}
    >
      {children}
    </span>
  );
}
