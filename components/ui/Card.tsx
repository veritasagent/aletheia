"use client";

import { motion } from "framer-motion";
import type { HTMLMotionProps } from "framer-motion";
import type { ReactNode } from "react";
import { cn } from "@/lib/utils";

interface CardProps extends Omit<HTMLMotionProps<"div">, "children"> {
  children: ReactNode;
  hover?: boolean;
}

export default function Card({
  children,
  className,
  onClick,
  hover = false,
  ...rest
}: CardProps) {
  return (
    <motion.div
      onClick={onClick}
      whileHover={hover ? { translateY: -2, shadow: "0 10px 15px -3px rgb(0 0 0 / 0.05)" } : undefined}
      transition={{ duration: 0.2, ease: "easeOut" }}
      className={cn(
        "rounded-xl border border-[#E5E7EB] bg-white p-6 shadow-sm",
        onClick ? "cursor-pointer" : "",
        className,
      )}
      {...rest}
    >
      {children}
    </motion.div>
  );
}
