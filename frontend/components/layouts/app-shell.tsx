import type { ReactNode } from 'react'
import { AppHeader } from '@/components/layouts/app-header'

export function AppShell({ children }: { children: ReactNode }) {
  return (
    <div className="flex min-h-svh flex-col">
      <AppHeader />
      <main className="mx-auto w-full max-w-6xl flex-1 px-4 py-8 sm:px-6">
        {children}
      </main>
      <footer className="border-t border-border">
        <div className="mx-auto flex max-w-6xl flex-col items-center justify-between gap-2 px-4 py-6 text-xs text-muted-foreground sm:flex-row sm:px-6">
          <div className="flex items-center gap-2 font-mono uppercase tracking-widest">
            <span className="h-2 w-2 rotate-45 bg-primary" />
            SyncTask / TASK SYSTEM
          </div>
          <p>© 2026 SyncTask. すべての作業は計画から。</p>
        </div>
      </footer>
    </div>
  )
}
