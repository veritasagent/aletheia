"use client";

import type { ReactNode } from "react";

interface SectionLabelProps {
  children: ReactNode;
  className?: string;
}

export default function SectionLabel({ children, className }: SectionLabelProps) {
  return (
    <div className={`flex items-center gap-3 ${className ?? ""}`}>
      <span className="font-sans text-[10px] font-bold uppercase tracking-widest text-[#9CA3AF]">{children}</span>
      <div className="h-px flex-1 bg-[#E5E7EB]" />
    </div>
  );
}
