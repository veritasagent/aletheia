import type { Config } from "tailwindcss";

const config: Config = {
  darkMode: "class",
  content: [
    "./app/**/*.{js,ts,jsx,tsx,mdx}",
    "./components/**/*.{js,ts,jsx,tsx,mdx}",
    "./lib/**/*.{js,ts,jsx,tsx,mdx}",
  ],
  theme: {
    extend: {
      colors: {
        bg: "var(--bg)",
        s1: "var(--s1)",
        s2: "var(--s2)",
        s3: "var(--s3)",
        s4: "var(--s4)",
        gold: "var(--gold)",
        grn: "var(--grn)",
        red: "var(--red)",
        amb: "var(--amb)",
        blu: "var(--blu)",
        t1: "var(--t1)",
        t2: "var(--t2)",
        t3: "var(--t3)",
      },
      fontFamily: {
        sans: ["var(--font-inter)", "Inter", "sans-serif"],
        heading: ["Sora", "var(--font-inter)", "sans-serif"],
      },
    },
  },
  plugins: [],
};

export default config;
