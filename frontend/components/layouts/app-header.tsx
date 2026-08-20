'use client'

import Link from 'next/link'
import { usePathname, useRouter } from 'next/navigation'
import { useState } from 'react'
import { LayoutGrid, ListChecks, LogOut, Menu, User, X } from 'lucide-react'
import { cn } from '@/lib/utils'
import { Button } from '@/components/ui/button'

const NAV = [
  { href: '/home', label: 'ホーム', icon: LayoutGrid },
  { href: '/tasks', label: 'タスク', icon: ListChecks },
  { href: '/profile', label: 'プロフィール', icon: User },
]

export function AppHeader() {
  const pathname = usePathname()
  const router = useRouter()
  const [open, setOpen] = useState(false)

  return (
    <header className="sticky top-0 z-40 border-b border-border bg-background/80 backdrop-blur">
      <div className="mx-auto flex h-16 max-w-6xl items-center justify-between gap-4 px-4 sm:px-6">
        <Link href="/home" className="flex items-center gap-2.5">
          <span className="grid h-8 w-8 place-items-center bg-primary text-primary-foreground">
            <span className="h-3 w-3 rotate-45 border-2 border-primary-foreground" />
          </span>
          <span className="font-heading text-lg font-bold tracking-tight">
            SyncTask
          </span>
        </Link>

        <nav className="hidden items-center gap-1 md:flex">
          {NAV.map((item) => {
            const active = pathname.startsWith(item.href)
            return (
              <Link
                key={item.href}
                href={item.href}
                className={cn(
                  'flex items-center gap-2 px-3 py-2 text-sm font-medium transition-colors',
                  active
                    ? 'bg-secondary text-foreground'
                    : 'text-muted-foreground hover:text-foreground',
                )}
              >
                <item.icon className="h-4 w-4" />
                {item.label}
              </Link>
            )
          })}
          <Button
            variant="ghost"
            size="sm"
            className="ml-2 gap-2 text-muted-foreground hover:text-foreground"
            onClick={() => router.push('/login')}
          >
            <LogOut className="h-4 w-4" />
            ログアウト
          </Button>
        </nav>

        <button
          className="grid h-9 w-9 place-items-center border border-border md:hidden"
          onClick={() => setOpen((v) => !v)}
          aria-label="メニューを開く"
        >
          {open ? <X className="h-5 w-5" /> : <Menu className="h-5 w-5" />}
        </button>
      </div>

      {open && (
        <div className="border-t border-border bg-background md:hidden">
          <nav className="mx-auto flex max-w-6xl flex-col px-4 py-2">
            {NAV.map((item) => {
              const active = pathname.startsWith(item.href)
              return (
                <Link
                  key={item.href}
                  href={item.href}
                  onClick={() => setOpen(false)}
                  className={cn(
                    'flex items-center gap-3 px-2 py-3 text-sm font-medium',
                    active ? 'text-primary' : 'text-muted-foreground',
                  )}
                >
                  <item.icon className="h-4 w-4" />
                  {item.label}
                </Link>
              )
            })}
            <button
              onClick={() => router.push('/login')}
              className="flex items-center gap-3 px-2 py-3 text-left text-sm font-medium text-muted-foreground"
            >
              <LogOut className="h-4 w-4" />
              ログアウト
            </button>
          </nav>
        </div>
      )}
    </header>
  )
}
