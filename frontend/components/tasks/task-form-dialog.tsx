'use client'

import { useMemo, useState } from 'react'
import { Repeat } from 'lucide-react'
import { toast } from 'sonner'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import { cn } from '@/lib/utils'
import { STATUS_META, type NewTask, type Priority, type Status, type Task } from '@/lib/store'

function toLocalInput(value: string | null) {
  if (!value) return ''
  const d = new Date(value)
  const off = d.getTimezoneOffset()
  const local = new Date(d.getTime() - off * 60000)
  return local.toISOString().slice(0, 16)
}

function todayInput() {
  const d = new Date()
  const off = d.getTimezoneOffset()
  return new Date(d.getTime() - off * 60000).toISOString().slice(0, 10)
}

// Sun..Sat order; value maps to JS getDay() (0=Sun..6=Sat)
const WEEKDAYS: { label: string; day: number }[] = [
  { label: '日', day: 0 },
  { label: '月', day: 1 },
  { label: '火', day: 2 },
  { label: '水', day: 3 },
  { label: '木', day: 4 },
  { label: '金', day: 5 },
  { label: '土', day: 6 },
]

const MAX_GENERATED = 100
const MAX_RANGE_DAYS = 52 * 7

export function TaskFormDialog({
  open,
  onOpenChange,
  task,
  presetDue,
  onSubmit,
  onSubmitMany,
}: {
  open: boolean
  onOpenChange: (open: boolean) => void
  task?: Task | null
  presetDue?: string | null
  onSubmit: (data: NewTask) => void
  onSubmitMany?: (list: NewTask[]) => void
}) {
  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      {open && (
        <TaskFormContent
          key={task?.id ?? presetDue ?? 'new'}
          task={task}
          presetDue={presetDue}
          onOpenChange={onOpenChange}
          onSubmit={onSubmit}
          onSubmitMany={onSubmitMany}
        />
      )}
    </Dialog>
  )
}

function TaskFormContent({
  task,
  presetDue,
  onOpenChange,
  onSubmit,
  onSubmitMany,
}: {
  task?: Task | null
  presetDue?: string | null
  onOpenChange: (open: boolean) => void
  onSubmit: (data: NewTask) => void
  onSubmitMany?: (list: NewTask[]) => void
}) {
  const editing = Boolean(task)
  const [title, setTitle] = useState(task?.title ?? '')
  const [priority, setPriority] = useState<Priority>(task?.priority ?? 'medium')
  const [due, setDue] = useState(toLocalInput(task?.due ?? presetDue ?? null))
  const [comment, setComment] = useState(task?.comment ?? '')
  const [status, setStatusValue] = useState<Status>(task?.status ?? 'todo')

  // repeat state (create mode only)
  const [repeat, setRepeat] = useState(false)
  const [startDate, setStartDate] = useState(todayInput())
  const [endDate, setEndDate] = useState(todayInput())
  const [time, setTime] = useState('')
  const [days, setDays] = useState<number[]>([])

  const toggleDay = (day: number) =>
    setDays((prev) =>
      prev.includes(day) ? prev.filter((d) => d !== day) : [...prev, day],
    )

  // preview count of generated tasks
  const genCount = useMemo(() => {
    if (!repeat || !startDate || !endDate || days.length === 0) return 0
    const start = new Date(`${startDate}T00:00:00`)
    const end = new Date(`${endDate}T00:00:00`)
    if (start > end) return -1
    const rangeDays = Math.round(
      (Date.UTC(end.getFullYear(), end.getMonth(), end.getDate()) -
        Date.UTC(start.getFullYear(), start.getMonth(), start.getDate())) /
        86_400_000,
    )
    if (rangeDays > MAX_RANGE_DAYS) return -2
    let count = 0
    const cur = new Date(start)
    let guard = 0
    while (cur <= end && guard < 400) {
      if (days.includes(cur.getDay())) count++
      cur.setDate(cur.getDate() + 1)
      guard++
    }
    return count
  }, [repeat, startDate, endDate, days])

  const submitSingle = () => {
    const dueIso = due ? new Date(due).toISOString() : null
    onSubmit({
      title: title.trim(),
      priority,
      due: dueIso,
      comment: comment.trim(),
      ...(editing ? { status } : {}),
    })
    onOpenChange(false)
  }

  const submitRepeat = () => {
    if (days.length === 0) {
      toast.error('曜日を1つ以上選択してください')
      return
    }
    if (genCount === -1) {
      toast.error('開始日は終了日以前に設定してください')
      return
    }
    if (genCount === -2) {
      toast.error('期間は最大52週以内に設定してください')
      return
    }
    if (genCount === 0) {
      toast.error('指定された期間内に該当する曜日が存在しません')
      return
    }
    if (genCount > MAX_GENERATED) {
      toast.error(`生成件数が上限（${MAX_GENERATED}件）を超えています`)
      return
    }
    const hhmm = time || '23:59'
    const start = new Date(`${startDate}T00:00:00`)
    const end = new Date(`${endDate}T00:00:00`)
    const list: NewTask[] = []
    const cur = new Date(start)
    while (cur <= end) {
      if (days.includes(cur.getDay())) {
        const y = cur.getFullYear()
        const m = String(cur.getMonth() + 1).padStart(2, '0')
        const d = String(cur.getDate()).padStart(2, '0')
        list.push({
          title: title.trim(),
          priority,
          due: new Date(`${y}-${m}-${d}T${hhmm}:00`).toISOString(),
          comment: comment.trim(),
        })
      }
      cur.setDate(cur.getDate() + 1)
    }
    onSubmitMany?.(list)
    onOpenChange(false)
  }

  const submit = (e: React.FormEvent) => {
    e.preventDefault()
    if (!title.trim()) {
      toast.error('タスク名を入力してください')
      return
    }
    if (repeat && !editing) submitRepeat()
    else submitSingle()
  }

  return (
    <DialogContent className="max-h-[90vh] overflow-y-auto sm:max-w-md">
        <DialogHeader>
          <DialogTitle className="font-heading">
            {editing ? 'タスクを編集' : 'タスクを作成'}
          </DialogTitle>
          <DialogDescription>
            タスクの属性を入力して{editing ? '更新' : '作成'}します。
          </DialogDescription>
        </DialogHeader>

        <form id="task-form" className="flex flex-col gap-4" onSubmit={submit}>
          <div className="grid gap-2">
            <Label htmlFor="title">タスク名</Label>
            <Input
              id="title"
              placeholder="例: 仕様レビューを完了する"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              maxLength={100}
              required
            />
          </div>

          <div className="grid gap-2">
            <Label htmlFor="priority">優先度</Label>
            <Select
              value={priority}
              onValueChange={(v) => setPriority(v as Priority)}
            >
              <SelectTrigger id="priority">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="high">高</SelectItem>
                <SelectItem value="medium">中</SelectItem>
                <SelectItem value="low">低</SelectItem>
              </SelectContent>
            </Select>
          </div>

          {editing && (
            <div className="grid gap-2">
              <Label htmlFor="status">ステータス</Label>
              <Select
                value={status}
                onValueChange={(v) => setStatusValue(v as Status)}
              >
                <SelectTrigger id="status">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="todo">{STATUS_META.todo.label}</SelectItem>
                  <SelectItem value="in_progress">{STATUS_META.in_progress.label}</SelectItem>
                  <SelectItem value="done">{STATUS_META.done.label}</SelectItem>
                </SelectContent>
              </Select>
            </div>
          )}

          {!repeat && (
            <div className="grid gap-2">
              <Label htmlFor="due">締切り日時（任意）</Label>
              <Input
                id="due"
                type="datetime-local"
                value={due}
                onChange={(e) => setDue(e.target.value)}
              />
            </div>
          )}

          <div className="grid gap-2">
            <Label htmlFor="comment">コメント</Label>
            <Textarea
              id="comment"
              placeholder="補足やメモを入力..."
              rows={3}
              maxLength={1000}
              value={comment}
              onChange={(e) => setComment(e.target.value)}
            />
          </div>

          {/* Repeat section — create only */}
          {!editing && (
            <div className="border border-border">
              <button
                type="button"
                onClick={() => setRepeat((v) => !v)}
                className="flex w-full items-center justify-between gap-3 p-3 text-left"
              >
                <span className="flex items-center gap-2">
                  <Repeat className="h-4 w-4 text-primary" />
                  <span className="text-sm font-medium">繰り返し作成</span>
                </span>
                <span
                  className={cn(
                    'relative h-5 w-9 shrink-0 rounded-full transition-colors',
                    repeat ? 'bg-primary' : 'bg-secondary',
                  )}
                >
                  <span
                    className={cn(
                      'absolute top-0.5 h-4 w-4 rounded-full bg-background transition-all',
                      repeat ? 'left-4' : 'left-0.5',
                    )}
                  />
                </span>
              </button>

              {repeat && (
                <div className="flex flex-col gap-4 border-t border-border p-3">
                  <div className="grid grid-cols-2 gap-3">
                    <div className="grid gap-1.5">
                      <Label htmlFor="start" className="text-xs">
                        開始日
                      </Label>
                      <Input
                        id="start"
                        type="date"
                        value={startDate}
                        onChange={(e) => setStartDate(e.target.value)}
                      />
                    </div>
                    <div className="grid gap-1.5">
                      <Label htmlFor="end" className="text-xs">
                        終了日
                      </Label>
                      <Input
                        id="end"
                        type="date"
                        value={endDate}
                        onChange={(e) => setEndDate(e.target.value)}
                      />
                    </div>
                  </div>

                  <div className="grid gap-1.5">
                    <Label htmlFor="time" className="text-xs">
                      締切時刻（任意・未指定は 23:59）
                    </Label>
                    <Input
                      id="time"
                      type="time"
                      value={time}
                      onChange={(e) => setTime(e.target.value)}
                    />
                  </div>

                  <div className="grid gap-1.5">
                    <Label className="text-xs">曜日</Label>
                    <div className="flex flex-wrap gap-1.5">
                      {WEEKDAYS.map((w) => {
                        const active = days.includes(w.day)
                        return (
                          <button
                            key={w.day}
                            type="button"
                            onClick={() => toggleDay(w.day)}
                            className={cn(
                              'h-9 w-9 border text-sm font-medium transition-colors',
                              active
                                ? 'border-primary bg-primary text-primary-foreground'
                                : 'border-border text-muted-foreground hover:border-primary/40',
                            )}
                          >
                            {w.label}
                          </button>
                        )
                      })}
                    </div>
                  </div>

                  <div
                    className={cn(
                      'flex items-center justify-between border px-3 py-2 text-xs',
                      genCount > MAX_GENERATED || genCount < 0
                        ? 'border-destructive/40 bg-destructive/10 text-destructive'
                        : 'border-border bg-secondary/40 text-muted-foreground',
                    )}
                  >
                    <span>生成予定件数</span>
                    <span className="font-mono font-bold">
                      {genCount === -1
                        ? '日付順エラー'
                        : genCount === -2
                          ? '52週上限超過'
                        : `${genCount} / ${MAX_GENERATED} 件`}
                    </span>
                  </div>
                </div>
              )}
            </div>
          )}
        </form>

        <DialogFooter>
          <Button variant="outline" onClick={() => onOpenChange(false)}>
            キャンセル
          </Button>
          <Button type="submit" form="task-form" className="font-semibold">
            決定
          </Button>
        </DialogFooter>
      </DialogContent>
  )
}
