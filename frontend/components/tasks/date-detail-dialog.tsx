'use client'

import { MessageSquare, Pencil, Pin, Plus } from 'lucide-react'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import {
  PRIORITY_META,
  STATUS_META,
  type Status,
  type Task,
} from '@/lib/store'
import { cn } from '@/lib/utils'

const STATUS_ORDER: Status[] = ['todo', 'in_progress', 'done']

function timeLabel(due: string | null) {
  if (!due) return '--:--'
  return new Date(due).toLocaleTimeString('ja-JP', {
    hour: '2-digit',
    minute: '2-digit',
  })
}

export function DateDetailDialog({
  date,
  tasks,
  onOpenChange,
  onSelectTask,
  onCreate,
  onSetStatus,
  onTogglePin,
}: {
  date: Date | null
  tasks: Task[]
  onOpenChange: (open: boolean) => void
  onSelectTask: (task: Task) => void
  onCreate: (date: Date) => void
  onSetStatus: (id: string, status: Status) => void
  onTogglePin: (id: string) => void
}) {
  const label = date
    ? date.toLocaleDateString('ja-JP', {
        year: 'numeric',
        month: 'long',
        day: 'numeric',
        weekday: 'short',
      })
    : ''

  return (
    <Dialog open={Boolean(date)} onOpenChange={onOpenChange}>
      <DialogContent className="max-h-[85vh] overflow-y-auto sm:max-w-md">
        <DialogHeader>
          <DialogTitle className="font-heading">{label}</DialogTitle>
        </DialogHeader>

        <div className="flex flex-col gap-3">
          {tasks.length === 0 ? (
            <p className="py-6 text-center text-sm text-muted-foreground">
              この日に締切のタスクはありません。
            </p>
          ) : (
            tasks.map((t) => {
              const prio = PRIORITY_META[t.priority]
              return (
                <div
                  key={t.id}
                  className="relative flex flex-col gap-2 border border-border bg-card p-3"
                >
                  <span className={cn('absolute left-0 top-0 h-full w-1', prio.dot)} />
                  <div className="flex items-start justify-between gap-2 pl-2">
                    <div className="min-w-0">
                      <p
                        className={cn(
                          'font-medium leading-snug',
                          t.status === 'done' &&
                            'text-muted-foreground line-through',
                        )}
                      >
                        {t.title}
                      </p>
                      <p className="mt-0.5 flex items-center gap-2 text-xs text-muted-foreground">
                        <span className="font-mono">{timeLabel(t.due)}</span>
                        <span className={prio.text}>優先度 {prio.label}</span>
                      </p>
                    </div>
                    <div className="flex shrink-0 items-center gap-1">
                      <button
                        onClick={() => onTogglePin(t.id)}
                        className={cn(
                          'grid h-7 w-7 place-items-center transition-colors',
                          t.pinned
                            ? 'text-primary'
                            : 'text-muted-foreground/50 hover:text-foreground',
                        )}
                        aria-label={t.pinned ? 'ピン留めを解除' : 'ピン留め'}
                        title={t.pinned ? 'ピン留めを解除' : 'ピン留め'}
                      >
                        <Pin className={cn('h-3.5 w-3.5', t.pinned && 'fill-current')} />
                      </button>
                      <button
                        onClick={() => onSelectTask(t)}
                        className="grid h-7 w-7 shrink-0 place-items-center text-muted-foreground hover:bg-secondary hover:text-foreground"
                        aria-label="編集"
                      >
                        <Pencil className="h-3.5 w-3.5" />
                      </button>
                    </div>
                  </div>
                  {t.comment && (
                    <div className="flex items-start gap-1.5 pl-2 text-xs text-muted-foreground">
                      <MessageSquare className="mt-0.5 h-3 w-3 shrink-0" />
                      <p className="line-clamp-2">{t.comment}</p>
                    </div>
                  )}
                  {/* status selector */}
                  <div className="flex gap-1 pl-2">
                    {STATUS_ORDER.map((s) => (
                      <button
                        key={s}
                        onClick={() => onSetStatus(t.id, s)}
                        className={cn(
                          'flex-1 border px-2 py-1 text-xs font-medium transition-colors',
                          t.status === s
                            ? 'border-primary bg-primary/15 text-primary'
                            : 'border-border text-muted-foreground hover:border-primary/40',
                        )}
                      >
                        {STATUS_META[s].label}
                      </button>
                    ))}
                  </div>
                </div>
              )
            })
          )}

          {date && (
            <Button
              variant="outline"
              className="mt-1 w-full gap-2"
              onClick={() => onCreate(date)}
            >
              <Plus className="h-4 w-4" />
              この日に新規タスクを作成
            </Button>
          )}
        </div>
      </DialogContent>
    </Dialog>
  )
}
