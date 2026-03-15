"use client";

import Link from "next/link";
import { motion } from "framer-motion";
import { CheckCircle2, ArrowRight } from "lucide-react";

const container = {
  hidden: { opacity: 0 },
  visible: {
    opacity: 1,
    transition: {
      staggerChildren: 0.1,
      delayChildren: 0.1,
    },
  },
};

const item = {
  hidden: { opacity: 0, y: 10 },
  visible: { opacity: 1, y: 0, transition: { duration: 0.4, ease: [0.25, 1, 0.5, 1] } },
};

const features = [
  "Monitors viral narratives across multiple segments",
  "Verifies claims using multi-source evidence",
  "Generates SHA-256 cryptographic verification reports",
];

export default function HomePage() {
  return (
    <section className="relative flex min-h-[calc(100vh-74px)] items-center justify-center bg-[#FAFAFB] px-5 py-20">
      <motion.div
        variants={container}
        initial="hidden"
        animate="visible"
        className="relative z-10 mx-auto flex max-w-4xl flex-col items-center text-center"
      >
        <motion.p
          variants={item}
          className="font-mono text-[10px] sm:text-xs uppercase tracking-[0.25em] font-bold text-[#D4A74F] mb-2"
        >
          Intelligence Engine
        </motion.p>

        <motion.h1
          variants={item}
          className="font-heading text-6xl font-extrabold tracking-tight text-[#111827] sm:text-7xl md:text-8xl"
        >
          ALETHEIA
        </motion.h1>

        <motion.p variants={item} className="mt-6 text-xl text-[#4B5563] font-medium">
          Real-Time Misinformation Intelligence
        </motion.p>

        <motion.div variants={item} className="mt-12 max-w-xl mx-auto">
          <ul className="space-y-4 text-left">
            {features.map((line) => (
              <li key={line} className="flex items-start gap-3 text-base text-[#374151]">
                <CheckCircle2 className="mt-1 h-5 w-5 flex-none text-[#111827]" />
                <span>{line}</span>
              </li>
            ))}
          </ul>
        </motion.div>

        <motion.div variants={item} className="mt-12">
          <Link
            href="/monitor"
            className="group inline-flex items-center gap-2 rounded-xl bg-[#111827] px-8 py-4 font-semibold text-white shadow-lg shadow-black/5 transition hover:bg-[#1f2937] hover:translate-y-[-1px] active:translate-y-0"
          >
            Start Monitoring <ArrowRight className="h-4 w-4 transition-transform group-hover:translate-x-1" />
          </Link>
        </motion.div>
      </motion.div>

      <div className="absolute bottom-8 left-1/2 -translate-x-1/2 text-center text-xs text-[#9CA3AF] pointer-events-none">
        ἀλήθεια — the state of not being hidden · unconcealedness · truth
      </div>
    </section>
  );
}
