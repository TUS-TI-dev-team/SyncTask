import type { ReactNode } from 'react'

export function AuthShell({
  title,
  subtitle,
  step,
  children,
}: {
  title: string
  subtitle?: string
  step?: string
  children: ReactNode
}) {
  return (
    <div className="grid min-h-svh lg:grid-cols-2">
      {/* Left geometric panel */}
      <div className="relative hidden overflow-hidden bg-primary lg:block">
        <div className="bg-grid absolute inset-0 opacity-20" />
        <div className="absolute -right-24 -top-24 h-96 w-96 rotate-45 border-[3px] border-primary-foreground/20" />
        <div className="absolute right-16 top-40 h-40 w-40 rotate-12 border-[3px] border-primary-foreground/30" />
        <div className="absolute -bottom-20 -left-16 h-80 w-80 rounded-full border-[3px] border-primary-foreground/20" />
        <div className="absolute bottom-32 left-32 h-4 w-4 bg-primary-foreground" />
        <div className="absolute left-1/2 top-1/3 h-2 w-2 bg-primary-foreground" />

        <div className="relative flex h-full flex-col justify-between p-12 text-primary-foreground">
          <div className="flex items-center gap-3">
            <span className="grid h-10 w-10 place-items-center border-2 border-primary-foreground">
              <span className="h-3.5 w-3.5 rotate-45 bg-primary-foreground" />
            </span>
            <span className="font-heading text-2xl font-bold tracking-tight">
              SyncTask
            </span>
          </div>
          <div className="max-w-md">
            <h2 className="font-heading text-4xl font-bold leading-tight text-balance">
              計画は、構造から始まる。
            </h2>
            <p className="mt-4 text-primary-foreground/80 leading-relaxed">
              優先度・締切・ピン留めの3つの軸でタスクを俯瞰。すべての作業を一つの幾何学的なワークスペースに。
            </p>
          </div>
          <div className="flex items-center gap-2 font-mono text-xs uppercase tracking-[0.3em] text-primary-foreground/70">
            <span className="h-2 w-2 rotate-45 bg-primary-foreground" />
            TASK MANAGEMENT SYSTEM
          </div>
        </div>
      </div>

      {/* Right form panel */}
      <div className="relative flex items-center justify-center px-4 py-12 sm:px-8">
        <div className="bg-grid bg-grid-fade absolute inset-0 opacity-[0.07]" />
        <div className="relative w-full max-w-sm">
          <div className="mb-8">
            {step && (
              <p className="mb-3 inline-flex items-center gap-2 font-mono text-xs uppercase tracking-widest text-primary">
                <span className="h-1.5 w-1.5 bg-primary" />
                {step}
              </p>
            )}
            <h1 className="font-heading text-2xl font-bold tracking-tight text-balance">
              {title}
            </h1>
            {subtitle && (
              <p className="mt-2 text-sm text-muted-foreground leading-relaxed">
                {subtitle}
              </p>
            )}
          </div>
          {children}
        </div>
      </div>
    </div>
  )
}
