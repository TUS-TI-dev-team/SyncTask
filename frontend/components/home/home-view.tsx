'use client'

import { useMemo, useState } from 'react'
import { Flame, Pin, Plus, Timer } from 'lucide-react'
import { toast } from 'sonner'
import {
  deadlineSort,
  hoursUntil,
  useApp,
  type NewTask,
  type Task,
} from '@/lib/store'
import { cn } from '@/lib/utils'
import { Button } from '@/components/ui/button'
import { TaskCard } from '@/components/tasks/task-card'
import { TaskFormDialog } from '@/components/tasks/task-form-dialog'
import { DeleteTaskDialog } from '@/components/tasks/delete-task-dialog'
import { CompletedToggle, Pagination } from '@/components/common/pagination'

type ViewKey = 'priority' | 'deadline' | 'pinned'

const PER_PAGE = 20

const VIEWS: {
  key: ViewKey
  label: string
  icon: typeof Flame
  desc: string
}[] = [
  { key: 'priority', label: '高優先度', icon: Flame, desc: '優先度が高いタスク' },
  { key: 'deadline', label: '締切間近', icon: Timer, desc: '72時間以内・超過分' },
  { key: 'pinned', label: 'ピン留め', icon: Pin, desc: 'ピン留めしたタスク' },
]

export function HomeView() {
  const { tasks, addTask, addTasks, updateTask, deleteTask, cycleStatus, togglePin } =
    useApp()
  const [view, setView] = useState<ViewKey>('priority')
  const [showCompleted, setShowCompleted] = useState(true)
  const [page, setPage] = useState(1)
  const [formOpen, setFormOpen] = useState(false)
  const [editing, setEditing] = useState<Task | null>(null)
  const [deleteTarget, setDeleteTarget] = useState<Task | null>(null)

  // deadline view always excludes completed
  const completedDisabled = view === 'deadline'
  const includeCompleted = showCompleted && !completedDisabled

  const filtered = useMemo(() => {
    if (view === 'priority') {
      return tasks
        .filter((t) => t.priority === 'high')
        .filter((t) => (includeCompleted ? true : t.status !== 'done'))
        .sort(deadlineSort)
    }
    if (view === 'deadline') {
      return tasks
        .filter((t) => t.status !== 'done')
        .filter((t) => t.due && hoursUntil(t.due) <= 72)
        .sort(deadlineSort)
    }
    return tasks
      .filter((t) => t.pinned)
      .filter((t) => (includeCompleted ? true : t.status !== 'done'))
      .sort(deadlineSort)
  }, [tasks, view, includeCompleted])

  const counts = useMemo(
    () => ({
      priority: tasks.filter((t) => t.status !== 'done' && t.priority === 'high')
        .length,
      deadline: tasks.filter(
        (t) => t.status !== 'done' && t.due && hoursUntil(t.due) <= 72,
      ).length,
      pinned: tasks.filter((t) => t.pinned && t.status !== 'done').length,
    }),
    [tasks],
  )

  const totalPages = Math.max(1, Math.ceil(filtered.length / PER_PAGE))
  const safePage = Math.min(Math.max(1, page), totalPages)
  const pageItems = filtered.slice((safePage - 1) * PER_PAGE, safePage * PER_PAGE)

  const handleSetView = (newView: ViewKey) => {
    setView(newView)
    setPage(1)
  }

  const handleSetShowCompleted = (newVal: boolean) => {
    setShowCompleted(newVal)
    setPage(1)
  }

  const openCreate = () => {
    setEditing(null)
    setFormOpen(true)
  }
  const openEdit = (task: Task) => {
    setEditing(task)
    setFormOpen(true)
  }
  const handleSubmit = (data: NewTask) => {
    if (editing) {
      updateTask(editing.id, data)
      toast.success('タスクを更新しました')
    } else {
      addTask(data)
      toast.success('タスクを作成しました')
    }
  }
  const handleSubmitMany = (list: NewTask[]) => {
    addTasks(list)
    toast.success(`${list.length}件のタスクを作成しました`)
  }

  return (
    <div className="flex flex-col gap-8">
      {/* Hero */}
      <div className="relative overflow-hidden border border-border bg-card">
        <div className="bg-grid bg-grid-fade absolute inset-0 opacity-10" />
        <div className="relative flex flex-col gap-6 p-6 sm:flex-row sm:items-end sm:justify-between sm:p-8">
          <div>
            <p className="inline-flex items-center gap-2 font-mono text-xs uppercase tracking-widest text-primary">
              <span className="h-1.5 w-1.5 bg-primary" />
              ダッシュボード
            </p>
            <h1 className="mt-3 font-heading text-3xl font-bold tracking-tight text-balance">
              今日の作業を、構造で捉える。
            </h1>
            <p className="mt-2 max-w-md text-sm text-muted-foreground leading-relaxed">
              3つの軸でタスクを切り替えて表示。優先度・締切・ピン留めから、いま手をつけるべきものを見極めましょう。
            </p>
          </div>
          <Button onClick={openCreate} className="gap-2 font-semibold">
            <Plus className="h-4 w-4" />
            タスクを作成
          </Button>
        </div>
      </div>

      {/* View switcher */}
      <div className="grid gap-3 sm:grid-cols-3">
        {VIEWS.map((v) => {
          const active = view === v.key
          return (
            <button
              key={v.key}
              onClick={() => handleSetView(v.key)}
              className={cn(
                'flex items-center gap-3 border p-4 text-left transition-colors',
                active
                  ? 'border-primary bg-primary/10'
                  : 'border-border bg-card hover:border-primary/40',
              )}
            >
              <span
                className={cn(
                  'grid h-10 w-10 shrink-0 place-items-center',
                  active
                    ? 'bg-primary text-primary-foreground'
                    : 'bg-secondary text-muted-foreground',
                )}
              >
                <v.icon className="h-5 w-5" />
              </span>
              <span className="min-w-0">
                <span className="flex items-center gap-2">
                  <span className="font-heading font-semibold">{v.label}</span>
                  <span
                    className={cn(
                      'font-mono text-xs',
                      active ? 'text-primary' : 'text-muted-foreground',
                    )}
                  >
                    {counts[v.key]}
                  </span>
                </span>
                <span className="block truncate text-xs text-muted-foreground">
                  {v.desc}
                </span>
              </span>
            </button>
          )
        })}
      </div>

      {/* Filter row */}
      <div className="flex items-center justify-between gap-3">
        <p className="text-sm text-muted-foreground">
          <span className="font-mono font-bold text-foreground">
            {filtered.length}
          </span>{' '}
          件のタスク
        </p>
        <CompletedToggle
          checked={includeCompleted}
          onChange={handleSetShowCompleted}
          disabled={completedDisabled}
        />
      </div>

      <Pagination page={page} totalPages={totalPages} onChange={setPage} />

      {/* Task list */}
      <div className="flex flex-col gap-3">
        {pageItems.length === 0 ? (
          <div className="flex flex-col items-center gap-3 border border-dashed border-border bg-card/50 px-6 py-16 text-center">
            <span className="grid h-12 w-12 place-items-center bg-secondary">
              <span className="h-4 w-4 rotate-45 border-2 border-muted-foreground" />
            </span>
            <p className="font-medium">表示するタスクがありません</p>
            <p className="text-sm text-muted-foreground">
              このビューに該当するタスクはまだありません。
            </p>
          </div>
        ) : (
          pageItems.map((task) => (
            <TaskCard
              key={task.id}
              task={task}
              onCycleStatus={() => cycleStatus(task.id)}
              onTogglePin={() => togglePin(task.id)}
              onEdit={() => openEdit(task)}
              onDelete={() => setDeleteTarget(task)}
            />
          ))
        )}
      </div>

      <Pagination page={page} totalPages={totalPages} onChange={setPage} />

      <TaskFormDialog
        open={formOpen}
        onOpenChange={setFormOpen}
        task={editing}
        onSubmit={handleSubmit}
        onSubmitMany={handleSubmitMany}
      />
      <DeleteTaskDialog
        open={Boolean(deleteTarget)}
        onOpenChange={(o) => !o && setDeleteTarget(null)}
        taskTitle={deleteTarget?.title}
        onConfirm={() => {
          if (deleteTarget) {
            deleteTask(deleteTarget.id)
            toast.success('タスクを削除しました')
            setDeleteTarget(null)
          }
        }}
      />
    </div>
  )
}
