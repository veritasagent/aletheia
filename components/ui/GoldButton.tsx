"use client";

import { motion } from "framer-motion";
import type { HTMLMotionProps } from "framer-motion";
import { Loader2 } from "lucide-react";
import type { ReactNode } from "react";
import { cn } from "@/lib/utils";

interface GoldButtonProps extends Omit<HTMLMotionProps<"button">, "children"> {
  loading?: boolean;
  children?: ReactNode;
}

export default function GoldButton({
  loading = false,
  className,
  children,
  disabled,
  ...rest
}: GoldButtonProps) {
  const isDisabled = disabled || loading;

  return (
    <motion.button
      type="button"
      whileHover={isDisabled ? undefined : { filter: "brightness(1.12)" }}
      whileTap={isDisabled ? undefined : { scale: 0.98 }}
      disabled={isDisabled}
      className={cn(
        "inline-flex items-center justify-center gap-2 rounded-[8px] bg-gold px-4 py-2 font-mono text-[10px] font-semibold uppercase tracking-[0.18em] text-black transition",
        isDisabled ? "cursor-not-allowed opacity-70" : "",
        className,
      )}
      {...rest}
    >
      {loading ? (
        <>
          <Loader2 className="h-3.5 w-3.5 animate-spin" />
          ANALYZING...
        </>
      ) : (
        children
      )}
    </motion.button>
  );
}
