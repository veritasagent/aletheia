"use client";

import { useState, useEffect } from "react";
import { motion, AnimatePresence } from "framer-motion";

type ToastMessage = { id: number; text: string };
let toastCounter = 0;
let toastSubscribers: ((toast: ToastMessage) => void)[] = [];

export function useToast() {
  return (text: string) => {
    const newToast = { id: ++toastCounter, text };
    toastSubscribers.forEach((sub) => sub(newToast));
  };
}

export function Toaster() {
  const [toasts, setToasts] = useState<ToastMessage[]>([]);

  useEffect(() => {
    const handleAdd = (t: ToastMessage) => {
      setToasts((prev) => [...prev, t]);
      setTimeout(() => {
        setToasts((prev) => prev.filter((item) => item.id !== t.id));
      }, 3000);
    };
    toastSubscribers.push(handleAdd);
    return () => {
      toastSubscribers = toastSubscribers.filter((s) => s !== handleAdd);
    };
  }, []);

  return (
    <div
      style={{
        position: "fixed",
        bottom: "24px",
        left: "50%",
        transform: "translateX(-50%)",
        display: "flex",
        flexDirection: "column",
        gap: "8px",
        zIndex: 9999,
        pointerEvents: "none",
        alignItems: "center",
      }}
    >
      <AnimatePresence>
        {toasts.map((t) => (
          <motion.div
            key={t.id}
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, scale: 0.95 }}
            transition={{ duration: 0.2 }}
            style={{
              background: "var(--s4)",
              color: "var(--t1)",
              border: "1px solid rgba(255,255,255,0.1)",
              padding: "10px 20px",
              borderRadius: "8px",
              fontSize: "13px",
              fontWeight: 500,
              boxShadow: "0 4px 12px rgba(0,0,0,0.5)",
            }}
          >
            {t.text}
          </motion.div>
        ))}
      </AnimatePresence>
    </div>
  );
}
