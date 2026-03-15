import type { Metadata } from "next";
import { Inter } from "next/font/google";
import "@fontsource/sora/800.css";
import "./globals.css";
import AppHeader from "@/components/AppHeader";
import { ConnectionProvider } from "@/lib/context";

const inter = Inter({
  subsets: ["latin"],
  variable: "--font-inter",
  display: "swap",
});

export const metadata: Metadata = {
  title: "ALETHEIA UI",
  description: "Real-time misinformation defense interface for ALETHEIA",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body className={`${inter.variable} bg-bg text-t1 antialiased`}>
        <ConnectionProvider>
          <div className="relative z-10 min-h-screen bg-bg">
            <AppHeader />
            <main>{children}</main>
          </div>
        </ConnectionProvider>
      </body>
    </html>
  );
}
