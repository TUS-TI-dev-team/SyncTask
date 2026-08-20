'use client'

import {
  CircleCheckBig,
  CircleDashed,
  CircleDot,
  Clock,
  Pencil,
  Pin,
  Trash2,
} from 'lucide-react'
import {
  daysUntil,
  formatDue,
  PRIORITY_META,
  STATUS_META,
  type Status,
  type Task,
} from '@/lib/store'
import { cn } from '@/lib/utils'

const STATUS_ICON: Record<Status, typeof CircleDashed> = {
  todo: CircleDashed,
  in_progress: CircleDot,
  done: CircleCheckBig,
}

const STATUS_COLOR: Record<Status, string> = {
  todo: 'text-muted-foreground',
  in_progress: 'text-chart-3',
  done: 'text-primary',
}

function dueLabel(task: Task) {
  if (task.status === 'done' || !task.due) return null
  const d = daysUntil(task.due)
  if (d < 0) return { text: `${Math.abs(d)}日超過`, cls: 'text-destructive' }
  if (d === 0) return { text: '今日締切', cls: 'text-chart-3' }
  if (d === 1) return { text: '明日締切', cls: 'text-chart-3' }
  return null
}

export function TaskCard({
  task,
  onCycleStatus,
  onTogglePin,
  onEdit,
  onDelete,
}: {
  task: Task
  onCycleStatus: () => void
  onTogglePin: () => void
  onEdit: () => void
  onDelete: () => void
}) {
  const prio = PRIORITY_META[task.priority]
  const StatusIcon = STATUS_ICON[task.status]
  const due = dueLabel(task)
  const done = task.status === 'done'

  return (
    <div className="group relative flex flex-col gap-3 border border-border bg-card p-4 transition-colors hover:border-primary/50">
      {/* priority accent bar */}
      <span className={cn('absolute left-0 top-0 h-full w-1', prio.dot)} />

      <div className="flex items-start justify-between gap-3 pl-2">
        <div className="flex min-w-0 items-start gap-3">
          <button
            onClick={onCycleStatus}
            className={cn(
              'mt-0.5 shrink-0 transition-colors hover:opacity-80',
              STATUS_COLOR[task.status],
            )}
            aria-label="状態を変更"
            title={`状態: ${STATUS_META[task.status].label}（クリックで変更）`}
          >
            <StatusIcon className="h-5 w-5" />
          </button>
          <div className="min-w-0">
            <p
              className={cn(
                'font-medium leading-snug text-pretty',
                done && 'text-muted-foreground line-through',
              )}
            >
              {task.title}
            </p>
            {task.comment && (
              <p className="mt-1 line-clamp-1 text-xs text-muted-foreground">
                {task.comment}
              </p>
            )}
          </div>
        </div>

        <button
          onClick={onTogglePin}
          className={cn(
            'shrink-0 transition-colors',
            task.pinned
              ? 'text-primary'
              : 'text-muted-foreground/50 hover:text-foreground',
          )}
          aria-label={task.pinned ? 'ピン留めを解除' : 'ピン留め'}
          title={task.pinned ? 'ピン留めを解除' : 'ピン留め'}
        >
          <Pin className={cn('h-4 w-4', task.pinned && 'fill-current')} />
        </button>
      </div>

      <div className="flex items-center justify-between gap-2 pl-2">
        <div className="flex flex-wrap items-center gap-2 text-xs">
          <span
            className={cn(
              'inline-flex items-center gap-1.5 border border-border px-2 py-0.5 font-medium',
              prio.text,
            )}
          >
            <span className={cn('h-1.5 w-1.5', prio.dot)} />
            優先度 {prio.label}
          </span>
          <span
            className={cn(
              'inline-flex items-center gap-1 border border-border px-2 py-0.5 font-medium',
              STATUS_COLOR[task.status],
            )}
          >
            {STATUS_META[task.status].label}
          </span>
          <span className="inline-flex items-center gap-1 text-muted-foreground">
            <Clock className="h-3.5 w-3.5" />
            {formatDue(task.due)}
          </span>
          {due && (
            <span className={cn('font-semibold', due.cls)}>{due.text}</span>
          )}
        </div>

        <div className="flex items-center gap-1 opacity-0 transition-opacity group-hover:opacity-100 focus-within:opacity-100">
          <button
            onClick={onEdit}
            className="grid h-7 w-7 place-items-center text-muted-foreground hover:bg-secondary hover:text-foreground"
            aria-label="編集"
          >
            <Pencil className="h-3.5 w-3.5" />
          </button>
          <button
            onClick={onDelete}
            className="grid h-7 w-7 place-items-center text-muted-foreground hover:bg-destructive/15 hover:text-destructive"
            aria-label="削除"
          >
            <Trash2 className="h-3.5 w-3.5" />
          </button>
        </div>
      </div>
    </div>
  )
}
