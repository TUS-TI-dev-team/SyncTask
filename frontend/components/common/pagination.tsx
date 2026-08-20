'use client'

import { ChevronLeft, ChevronRight } from 'lucide-react'
import { cn } from '@/lib/utils'

export function Pagination({
  page,
  totalPages,
  onChange,
}: {
  page: number
  totalPages: number
  onChange: (page: number) => void
}) {
  if (totalPages <= 1) return null

  const getPageNumbers = (): (number | 'ellipsis')[] => {
    if (totalPages <= 10) {
      return Array.from({ length: totalPages }, (_, i) => i + 1)
    }
    // N > 10
    if (page === 1) {
      return [1, 2, 3, 'ellipsis', totalPages]
    }
    if (page === totalPages) {
      return [1, 2, 'ellipsis', totalPages - 2, totalPages - 1, totalPages]
    }
    const result: (number | 'ellipsis')[] = [1, 2]
    if (page - 1 > 3) {
      result.push('ellipsis')
    } else if (page - 1 === 3) {
      result.push(3)
    }
    const mid = [page - 1, page, page + 1].filter(
      (p) => p > 2 && p < totalPages,
    )
    for (const p of mid) {
      if (!result.includes(p)) result.push(p)
    }
    if (page + 1 < totalPages - 2) {
      result.push('ellipsis')
    } else if (page + 1 === totalPages - 2) {
      if (!result.includes(totalPages - 1)) result.push(totalPages - 1)
    }
    if (!result.includes(totalPages)) {
      result.push(totalPages)
    }
    return result
  }

  const pageNumbers = getPageNumbers()

  return (
    <nav
      className="flex items-center justify-center gap-1"
      aria-label="ページネーション"
    >
      <button
        onClick={() => onChange(Math.max(1, page - 1))}
        disabled={page === 1}
        className="grid h-9 w-9 place-items-center border border-border text-muted-foreground transition-colors hover:bg-secondary hover:text-foreground disabled:cursor-not-allowed disabled:opacity-40"
        aria-label="前のページ"
      >
        <ChevronLeft className="h-4 w-4" />
      </button>
      {pageNumbers.map((p, i) =>
        p === 'ellipsis' ? (
          <span
            key={`ellipsis-${i}`}
            className="grid h-9 min-w-9 place-items-center text-sm text-muted-foreground"
          >
            …
          </span>
        ) : (
          <button
            key={p}
            onClick={() => onChange(p)}
            aria-current={p === page ? 'page' : undefined}
            className={cn(
              'h-9 min-w-9 border px-2 text-sm font-medium transition-colors',
              p === page
                ? 'border-primary bg-primary text-primary-foreground'
                : 'border-border text-muted-foreground hover:bg-secondary hover:text-foreground',
            )}
          >
            {p}
          </button>
        ),
      )}
      <button
        onClick={() => onChange(Math.min(totalPages, page + 1))}
        disabled={page === totalPages}
        className="grid h-9 w-9 place-items-center border border-border text-muted-foreground transition-colors hover:bg-secondary hover:text-foreground disabled:cursor-not-allowed disabled:opacity-40"
        aria-label="次のページ"
      >
        <ChevronRight className="h-4 w-4" />
      </button>
    </nav>
  )
}

export function CompletedToggle({
  checked,
  onChange,
  disabled,
}: {
  checked: boolean
  onChange: (v: boolean) => void
  disabled?: boolean
}) {
  return (
    <button
      type="button"
      onClick={() => !disabled && onChange(!checked)}
      disabled={disabled}
      className={cn(
        'inline-flex items-center gap-2 border border-border px-3 py-2 text-xs font-medium transition-colors',
        disabled
          ? 'cursor-not-allowed opacity-40'
          : 'hover:border-primary/40',
      )}
      title={disabled ? '締切間近ビューでは完了タスクは常に除外されます' : undefined}
    >
      <span
        className={cn(
          'relative h-4 w-7 shrink-0 rounded-full transition-colors',
          checked && !disabled ? 'bg-primary' : 'bg-secondary',
        )}
      >
        <span
          className={cn(
            'absolute top-0.5 h-3 w-3 rounded-full bg-background transition-all',
            checked && !disabled ? 'left-3.5' : 'left-0.5',
          )}
        />
      </span>
      完了タスクを表示
    </button>
  )
}
