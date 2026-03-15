"use client";

import { motion } from "framer-motion";

interface ProgressBarProps {
  value: number;
  color?: string;
  animated?: boolean;
  className?: string;
}

export default function ProgressBar({
  value,
  color = "var(--gold)",
  animated = true,
  className,
}: ProgressBarProps) {
  const safeValue = Math.max(0, Math.min(100, value));

  return (
    <div className={`h-2 w-full overflow-hidden rounded-full bg-[#F3F4F6] ${className ?? ""}`}>
      <motion.div
        className="h-full rounded-full"
        style={{ backgroundColor: color }}
        initial={{ width: 0 }}
        animate={{ width: `${safeValue}%` }}
        transition={animated ? { type: "spring", stiffness: 90, damping: 20 } : { duration: 0 }}
      />
    </div>
  );
}
