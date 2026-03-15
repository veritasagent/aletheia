"use client";

import Link from "next/link";
import { useEffect, useMemo, useState } from "react";
import { Clock3, Plug, PlugZap } from "lucide-react";
import { useConnectionContext } from "@/lib/context";

function formatClock(date: Date) {
  return new Intl.DateTimeFormat("en-US", {
    hour: "2-digit",
    minute: "2-digit",
    second: "2-digit",
    hour12: false,
  }).format(date);
}

const navItems = [
  { href: "/monitor", label: "Verify" },
  { href: "#live-feed", label: "Live Feed" },
  { href: "#sources", label: "Sources" },
];

export default function AppHeader() {
  const { connected } = useConnectionContext();
  const [now, setNow] = useState<Date>(new Date());

  useEffect(() => {
    const timer = window.setInterval(() => {
      setNow(new Date());
    }, 1000);
    return () => window.clearInterval(timer);
  }, []);

  const statusStyles = useMemo(() => {
    if (connected) {
      return "text-[#059669]";
    }
    return "text-[#DC2626]";
  }, [connected]);

  return (
    <header className="sticky top-0 z-50 border-b border-[#E5E7EB] bg-white/80 backdrop-blur-md">
      <div className="mx-auto flex h-[74px] max-w-7xl items-center justify-between px-4 md:px-6">
        <div className="flex items-center gap-6">
          <Link
            href="/"
            className="font-heading text-xl font-bold tracking-tight text-[#111827] transition hover:opacity-80 md:text-2xl"
          >
            ALETHEIA
          </Link>
          <div className="hidden h-5 w-[1px] bg-[#E5E7EB] sm:block" />
          <span className="hidden items-center gap-1.5 text-[10px] font-bold uppercase tracking-[0.2em] text-[#4B5563] sm:inline-flex">
            <span className={`h-2 w-2 rounded-full ${connected ? "bg-[#059669]" : "bg-[#DC2626]"}`} />
            {connected ? "Live" : "Standby"}
          </span>
        </div>

        <nav className="hidden items-center gap-8 lg:flex">
          {navItems.map((item) => (
            <Link
              key={item.href}
              href={item.href}
              className="font-sans text-xs font-semibold text-[#6B7280] transition hover:text-[#111827]"
            >
              {item.label}
            </Link>
          ))}
        </nav>

        <div className="flex items-center gap-4">
          <span className="hidden items-center gap-2 font-mono text-[10px] font-medium text-[#9CA3AF] md:inline-flex">
            <Clock3 className="h-3.5 w-3.5" />
            {formatClock(now)}
          </span>
          <div className={`flex items-center gap-2 rounded-full bg-[#F3F4F6] px-3 py-1.5 transition-colors ${statusStyles}`}>
            {connected ? <PlugZap className="h-3.5 w-3.5" /> : <Plug className="h-3.5 w-3.5" />}
            <span className="font-mono text-[10px] font-bold uppercase tracking-wider">
              {connected ? "Connected" : "Offline"}
            </span>
          </div>
        </div>
      </div>
    </header>
  );
}
