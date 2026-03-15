export default function Loading() {
  return (
    <section className="mx-auto w-full max-w-7xl px-4 py-5 md:px-6 md:py-6">
      <div className="flex flex-row items-center gap-3 font-mono text-[9px] uppercase tracking-wider mb-6 pb-4 border-b border-white/[0.05]">
        <div className="h-4 w-20 bg-[var(--s2)] rounded animate-pulse" />
        <span className="text-[var(--t3)]">›</span>
        <div className="h-4 w-20 bg-[var(--s2)] rounded animate-pulse" />
      </div>

      <div className="grid gap-6 lg:grid-cols-[minmax(0,1fr)_280px]">
        <div className="min-w-0 space-y-4">
          <div className="h-6 w-32 bg-[var(--s2)] rounded animate-pulse mb-6" />

          <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4">
            {[1, 2, 3, 4].map((i) => (
              <div 
                key={i} 
                className="flex aspect-video w-full flex-col justify-between rounded-xl bg-[var(--s2)] p-4 shadow-sm border border-white/[0.04] animate-pulse"
              >
                <div className="flex justify-between">
                  <div className="h-6 w-6 rounded-full bg-[var(--s4)]" />
                  <div className="h-6 w-8 rounded bg-[var(--s4)]" />
                </div>
                <div className="space-y-2 mt-4">
                  <div className="h-4 w-3/4 rounded bg-[var(--s4)]" />
                  <div className="h-3 w-1/2 rounded bg-[var(--s4)] opacity-50" />
                </div>
              </div>
            ))}
          </div>
        </div>
        
        <div className="hidden lg:block">
          <div className="h-[600px] w-full rounded-xl bg-[var(--s2)] animate-pulse border border-white/[0.04]" />
        </div>
      </div>
    </section>
  );
}
