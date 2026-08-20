'use client'

import { useMemo, useState } from 'react'
import { CalendarDays, List, Plus, Search, X } from 'lucide-react'
import { toast } from 'sonner'
import {
  defaultSort,
  useApp,
  type NewTask,
  type Priority,
  type Status,
  type Task,
} from '@/lib/store'
import { cn } from '@/lib/utils'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { TaskCard } from '@/components/tasks/task-card'
import { TaskCalendar } from '@/components/tasks/task-calendar'
import { TaskFormDialog } from '@/components/tasks/task-form-dialog'
import { DeleteTaskDialog } from '@/components/tasks/delete-task-dialog'
import { CompletedToggle, Pagination } from '@/components/common/pagination'

/** Normalize katakana to hiragana for kana-insensitive search. */
function toHiragana(str: string): string {
  return str.replace(/[\u30A1-\u30F6]/g, (ch) =>
    String.fromCharCode(ch.charCodeAt(0) - 0x60),
  )
}

type Mode = 'list' | 'calendar'
const PER_PAGE = 20

export function TasksView() {
  const {
    tasks,
    addTask,
    addTasks,
    updateTask,
    deleteTask,
    cycleStatus,
    setStatus,
    togglePin,
  } = useApp()
  const [mode, setMode] = useState<Mode>('list')
  const [query, setQuery] = useState('')
  const [priority, setPriority] = useState<Priority | 'all'>('all')
  const [status, setStatus_] = useState<Status | 'all'>('all')
  const [dueBefore, setDueBefore] = useState('')
  const [showCompleted, setShowCompleted] = useState(true)
  const [page, setPage] = useState(1)

  const [formOpen, setFormOpen] = useState(false)
  const [editing, setEditing] = useState<Task | null>(null)
  const [presetDue, setPresetDue] = useState<string | null>(null)
  const [deleteTarget, setDeleteTarget] = useState<Task | null>(null)

  const filtered = useMemo(() => {
    const q = toHiragana(query.trim().toLowerCase())
    let list = tasks.filter((t) => {
      if (q && !(toHiragana(t.title.toLowerCase()).includes(q) || toHiragana(t.comment.toLowerCase()).includes(q)))
        return false
      if (priority !== 'all' && t.priority !== priority) return false
      if (status !== 'all') {
        if (t.status !== status) return false
      } else if (!showCompleted && t.status === 'done') {
        return false
      }
      if (dueBefore) {
        if (!t.due) return false
        const limit = new Date(`${dueBefore}T23:59:59`)
        if (new Date(t.due) > limit) return false
      }
      return true
    })
    list = [...list].sort(defaultSort)
    return list
  }, [tasks, query, priority, status, dueBefore, showCompleted])

  const totalPages = Math.max(1, Math.ceil(filtered.length / PER_PAGE))
  const safePage = Math.min(Math.max(1, page), totalPages)
  const pageItems = filtered.slice((safePage - 1) * PER_PAGE, safePage * PER_PAGE)

  const handleQueryChange = (val: string) => {
    setQuery(val)
    setPage(1)
  }
  const handlePriorityChange = (val: Priority | 'all') => {
    setPriority(val)
    setPage(1)
  }
  const handleStatusChange = (val: Status | 'all') => {
    setStatus_(val)
    setPage(1)
  }
  const handleDueBeforeChange = (val: string) => {
    setDueBefore(val)
    setPage(1)
  }
  const handleShowCompletedChange = (val: boolean) => {
    setShowCompleted(val)
    setPage(1)
  }

  const hasFilters =
    query || priority !== 'all' || status !== 'all' || dueBefore
  const clearFilters = () => {
    setQuery('')
    setPriority('all')
    setStatus_('all')
    setDueBefore('')
    setPage(1)
  }

  const openCreate = () => {
    setEditing(null)
    setPresetDue(null)
    setFormOpen(true)
  }
  const openEdit = (task: Task) => {
    setEditing(task)
    setPresetDue(null)
    setFormOpen(true)
  }
  const openCreateForDate = (date: Date) => {
    const y = date.getFullYear()
    const m = String(date.getMonth() + 1).padStart(2, '0')
    const d = String(date.getDate()).padStart(2, '0')
    setEditing(null)
    setPresetDue(new Date(`${y}-${m}-${d}T23:59:00`).toISOString())
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
    <div className="flex flex-col gap-6">
      <div className="flex flex-col gap-1">
        <p className="inline-flex items-center gap-2 font-mono text-xs uppercase tracking-widest text-primary">
          <span className="h-1.5 w-1.5 bg-primary" />
          タスク一覧
        </p>
        <h1 className="font-heading text-3xl font-bold tracking-tight">
          すべてのタスク
        </h1>
      </div>

      {/* Toolbar */}
      <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        <div className="relative flex-1 sm:max-w-xs">
          <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
          <Input
            value={query}
            onChange={(e) => handleQueryChange(e.target.value)}
            placeholder="タスク名・コメントを検索..."
            className="pl-9"
          />
        </div>

        <div className="flex items-center gap-3">
          <div className="flex border border-border">
            <button
              onClick={() => setMode('list')}
              className={cn(
                'flex items-center gap-1.5 px-3 py-2 text-sm font-medium transition-colors',
                mode === 'list'
                  ? 'bg-primary text-primary-foreground'
                  : 'text-muted-foreground hover:text-foreground',
              )}
            >
              <List className="h-4 w-4" />
              リスト
            </button>
            <button
              onClick={() => setMode('calendar')}
              className={cn(
                'flex items-center gap-1.5 px-3 py-2 text-sm font-medium transition-colors',
                mode === 'calendar'
                  ? 'bg-primary text-primary-foreground'
                  : 'text-muted-foreground hover:text-foreground',
              )}
            >
              <CalendarDays className="h-4 w-4" />
              カレンダー
            </button>
          </div>
          <Button onClick={openCreate} className="gap-2 font-semibold">
            <Plus className="h-4 w-4" />
            <span className="hidden sm:inline">タスク作成</span>
          </Button>
        </div>
      </div>

      {/* Filters */}
      <div className="flex flex-col gap-4 border border-border bg-card p-4">
        <div className="grid gap-4 sm:grid-cols-3">
          <div className="grid gap-1.5">
            <Label className="text-xs">優先度</Label>
            <Select
              value={priority}
              onValueChange={(v) => handlePriorityChange(v as Priority | 'all')}
            >
              <SelectTrigger>
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">すべて</SelectItem>
                <SelectItem value="high">高</SelectItem>
                <SelectItem value="medium">中</SelectItem>
                <SelectItem value="low">低</SelectItem>
              </SelectContent>
            </Select>
          </div>
          <div className="grid gap-1.5">
            <Label className="text-xs">ステータス</Label>
            <Select
              value={status}
              onValueChange={(v) => handleStatusChange(v as Status | 'all')}
            >
              <SelectTrigger>
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">すべて</SelectItem>
                <SelectItem value="todo">未着手</SelectItem>
                <SelectItem value="in_progress">進行中</SelectItem>
                <SelectItem value="done">完了</SelectItem>
              </SelectContent>
            </Select>
          </div>
          <div className="grid gap-1.5">
            <Label htmlFor="dueBefore" className="text-xs">
              締切日（この日まで）
            </Label>
            <Input
              id="dueBefore"
              type="date"
              value={dueBefore}
              onChange={(e) => handleDueBeforeChange(e.target.value)}
            />
          </div>
        </div>
        <div className="flex items-center justify-end gap-3">
          {hasFilters && (
            <button
              onClick={clearFilters}
              className="inline-flex items-center gap-1 text-xs font-medium text-muted-foreground hover:text-foreground"
            >
              <X className="h-3.5 w-3.5" />
              条件をクリア
            </button>
          )}
        </div>
      </div>

      {/* Task count + completed toggle (shared for both modes) */}
      <div className="flex items-center justify-between gap-3">
        <p className="text-sm text-muted-foreground">
          <span className="font-mono font-bold text-foreground">
            {filtered.length}
          </span>{' '}
          件のタスク
        </p>
        <CompletedToggle
          checked={showCompleted}
          onChange={handleShowCompletedChange}
          disabled={status !== 'all'}
        />
      </div>

      {mode === 'list' ? (
        <>
          <Pagination page={page} totalPages={totalPages} onChange={setPage} />
          <div className="flex flex-col gap-3">
            {pageItems.length === 0 ? (
              <div className="flex flex-col items-center gap-3 border border-dashed border-border bg-card/50 px-6 py-16 text-center">
                <span className="grid h-12 w-12 place-items-center bg-secondary">
                  <span className="h-4 w-4 rotate-45 border-2 border-muted-foreground" />
                </span>
                <p className="font-medium">
                  {hasFilters ? '該当するタスクがありません' : 'タスクがありません'}
                </p>
                <p className="text-sm text-muted-foreground">
                  {hasFilters
                    ? '検索条件を変更してみてください。'
                    : '右上のボタンから新しいタスクを作成しましょう。'}
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
        </>
      ) : (
        <TaskCalendar
          tasks={filtered}
          onSelectTask={openEdit}
          onCreateForDate={openCreateForDate}
          onSetStatus={setStatus}
          onTogglePin={(id: string) => togglePin(id)}
        />
      )}

      <TaskFormDialog
        open={formOpen}
        onOpenChange={setFormOpen}
        task={editing}
        presetDue={presetDue}
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
