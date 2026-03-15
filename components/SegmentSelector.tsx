"use client";

import {
  Briefcase,
  Cpu,
  FlaskConical,
  Globe2,
  GraduationCap,
  HeartPulse,
  Landmark,
  type LucideIcon,
  ShieldCheck,
  Trophy,
  Vote,
} from "lucide-react";
import Card from "@/components/ui/Card";
import SectionLabel from "@/components/ui/SectionLabel";
import type { Segment } from "@/lib/types";
import { fallbackSegments } from "@/lib/sandbox";

interface SegmentSelectorProps {
  onSelect: (segment: Segment) => void;
  selectedSegmentId?: string;
  segments?: Segment[];
}

const iconMap: Record<string, LucideIcon> = {
  Landmark,
  Cpu,
  Briefcase,
  FlaskConical,
  Vote,
  Trophy,
  HeartPulse,
  Globe2,
  GraduationCap,
  ShieldCheck,
};

export default function SegmentSelector({
  onSelect,
  selectedSegmentId,
  segments = fallbackSegments,
}: SegmentSelectorProps) {
  return (
    <section id="segments" className="space-y-4">
      <SectionLabel>Select Monitoring Segment</SectionLabel>
      <div className="grid gap-3 md:grid-cols-2 xl:grid-cols-3">
        {segments.map((segment) => {
          const Icon = iconMap[segment.icon] ?? ShieldCheck;
          const selected = selectedSegmentId === segment.id;

          return (
            <Card
              key={segment.id}
              hover
              onClick={() => onSelect(segment)}
              className={`group cursor-pointer transition-all duration-300 border-[1px] ${
                selected 
                  ? "border-[#D4A74F] bg-white shadow-md shadow-[#D4A74F]/5 ring-1 ring-[#D4A74F]/5" 
                  : "border-[#E5E7EB] bg-white hover:border-[#D1D5DB] hover:shadow-sm"
              }`}
            >
              <div className="flex items-center justify-between">
                <span
                  className={`inline-flex h-10 w-10 items-center justify-center rounded-xl border transition-colors ${
                    selected
                      ? "border-[#D4A74F]/20 bg-[#D4A74F]/5 text-[#D4A74F]"
                      : "border-[#F3F4F6] bg-[#FAFAFB] text-[#9CA3AF] group-hover:text-[#6B7280]"
                  }`}
                >
                  <Icon className="h-5 w-5" />
                </span>
                <span className={`font-mono text-[9px] font-bold uppercase tracking-[0.2em] ${selected ? "text-[#D4A74F]" : "text-[#9CA3AF]"}`}>
                  ID: {segment.id}
                </span>
              </div>
              <h3 className={`font-heading mt-5 text-xl font-bold transition-colors ${selected ? "text-[#111827]" : "text-[#374151] group-hover:text-[#111827]"}`}>
                {segment.label}
              </h3>
              <p className="mt-1.5 text-xs text-[#6B7280] leading-relaxed">
                Analyze active narratives and verification flows in this segment.
              </p>
            </Card>
          );
        })}
      </div>
    </section>
  );
}
