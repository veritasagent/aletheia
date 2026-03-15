"use client";

import Link from "next/link";
import { motion } from "framer-motion";
import { ArrowRight, ShieldCheck, Sparkles, Workflow } from "lucide-react";

const features = [
  {
    title: "Real-Time Verification",
    description:
      "Ingests claim streams and evaluates credibility against trusted institutional evidence.",
    icon: ShieldCheck,
  },
  {
    title: "Agentic Reasoning",
    description:
      "Shows transparent reasoning chains for every verdict to support human oversight.",
    icon: Workflow,
  },
  {
    title: "Cryptographic Evidence",
    description:
      "Each decision is chained with SHA256 records for tamper-evident auditability.",
    icon: Sparkles,
  },
];

export default function Intro() {
  return (
    <main className="relative min-h-screen overflow-hidden bg-bg">
      <div className="absolute inset-0 bg-[radial-gradient(1100px_circle_at_15%_10%,rgba(77,123,255,0.2),transparent_45%),radial-gradient(900px_circle_at_80%_0%,rgba(200,168,74,0.2),transparent_40%)]" />
      <section className="relative z-10 mx-auto flex min-h-screen max-w-6xl flex-col justify-center px-6 py-20">
        <motion.p
          initial={{ opacity: 0, y: 10 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.5 }}
          className="mb-5 inline-flex w-fit items-center gap-2 rounded-full border border-s3 bg-s2/70 px-4 py-2 text-xs uppercase tracking-[0.2em] text-gold"
        >
          ALETHEIA UI · Sandbox Monitoring Interface
        </motion.p>

        <motion.h1
          initial={{ opacity: 0, y: 14 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.65 }}
          className="font-heading max-w-4xl text-4xl font-extrabold leading-tight text-t1 md:text-6xl"
        >
          Real-Time Misinformation Defense With Transparent AI Reasoning
        </motion.h1>

        <motion.p
          initial={{ opacity: 0, y: 18 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.1, duration: 0.65 }}
          className="mt-6 max-w-3xl text-base leading-relaxed text-t2 md:text-lg"
        >
          ALETHEIA monitors high-velocity digital narratives, verifies claims against trusted sources,
          and emits cryptographically linked evidence reports for institutional-grade auditability.
        </motion.p>

        <motion.div
          initial={{ opacity: 0, y: 16 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.2, duration: 0.6 }}
          className="mt-10"
        >
          <Link
            href="/monitor"
            className="inline-flex items-center gap-2 rounded-xl border border-blu/50 bg-blu/20 px-6 py-3 text-sm font-semibold text-t1 transition hover:bg-blu/30"
          >
            Launch Monitor
            <ArrowRight className="h-4 w-4" />
          </Link>
        </motion.div>

        <div className="mt-14 grid gap-4 md:grid-cols-3">
          {features.map((item, index) => {
            const Icon = item.icon;
            return (
              <motion.article
                key={item.title}
                initial={{ opacity: 0, y: 16 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: 0.28 + index * 0.08, duration: 0.55 }}
                className="rounded-2xl border border-s3 bg-s1/80 p-5"
              >
                <Icon className="h-5 w-5 text-gold" />
                <h2 className="font-heading mt-4 text-xl font-bold text-t1">{item.title}</h2>
                <p className="mt-2 text-sm leading-relaxed text-t2">{item.description}</p>
              </motion.article>
            );
          })}
        </div>
      </section>
    </main>
  );
}
