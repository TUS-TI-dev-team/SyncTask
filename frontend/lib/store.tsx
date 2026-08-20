'use client'

import {
  createContext,
  useContext,
  useMemo,
  useState,
  type ReactNode,
} from 'react'

export type Priority = 'high' | 'medium' | 'low'
export type Status = 'todo' | 'in_progress' | 'done'

export type Task = {
  id: string
  title: string
  priority: Priority
  due: string | null // ISO date string, null = no deadline
  comment: string
  status: Status
  pinned: boolean
  createdAt: string // ISO
}

export type Profile = {
  username: string
  email: string
}

function iso(daysFromNow: number, hour = 18, minute = 0) {
  const d = new Date()
  d.setDate(d.getDate() + daysFromNow)
  d.setHours(hour, minute, 0, 0)
  return d.toISOString()
}

function createdIso(minutesAgo: number) {
  return new Date(Date.now() - minutesAgo * 60000).toISOString()
}

const SEED_TASKS: Task[] = [
  {
    id: 't1',
    title: 'Q3 ローンチ計画の最終レビュー',
    priority: 'high',
    due: iso(0, 17),
    comment: 'ステークホルダー全員の承認が必要。',
    status: 'in_progress',
    pinned: true,
    createdAt: createdIso(600),
  },
  {
    id: 't2',
    title: '認証フローのセキュリティ監査',
    priority: 'high',
    due: iso(1),
    comment: 'OTP とパスワードリセットの経路を重点的に確認する。',
    status: 'todo',
    pinned: true,
    createdAt: createdIso(540),
  },
  {
    id: 't3',
    title: 'デザインシステムのトークン整理',
    priority: 'medium',
    due: iso(2),
    comment: '幾何学テーマに合わせて余白とラディウスを統一。',
    status: 'todo',
    pinned: false,
    createdAt: createdIso(480),
  },
  {
    id: 't4',
    title: '請求ダッシュボードのバグ修正',
    priority: 'high',
    due: iso(-1),
    comment: '締切超過。最優先で対応する。',
    status: 'todo',
    pinned: false,
    createdAt: createdIso(420),
  },
  {
    id: 't5',
    title: '週次レポートの作成',
    priority: 'low',
    due: iso(4),
    comment: '',
    status: 'todo',
    pinned: false,
    createdAt: createdIso(360),
  },
  {
    id: 't6',
    title: 'オンボーディング資料の更新',
    priority: 'medium',
    due: iso(3),
    comment: '新メンバー向けに手順を簡略化。',
    status: 'done',
    pinned: false,
    createdAt: createdIso(300),
  },
  {
    id: 't7',
    title: 'API レート制限の検証',
    priority: 'medium',
    due: iso(2, 10),
    comment: 'ピーク時の挙動を計測する。',
    status: 'in_progress',
    pinned: true,
    createdAt: createdIso(240),
  },
  {
    id: 't8',
    title: 'マーケティングLPのコピー調整',
    priority: 'low',
    due: iso(7),
    comment: '',
    status: 'todo',
    pinned: false,
    createdAt: createdIso(180),
  },
  {
    id: 't9',
    title: 'メモ: アイデア整理（締切なし）',
    priority: 'low',
    due: null,
    comment: '思いついたことを書き留める。',
    status: 'todo',
    pinned: false,
    createdAt: createdIso(120),
  },
]

export type NewTask = {
  title: string
  priority: Priority
  due: string | null
  comment: string
  status?: Status
}

let idCounter = 0
function nextId() {
  idCounter += 1
  return `t${Date.now()}_${idCounter}`
}

type AppState = {
  profile: Profile
  tasks: Task[]
  addTask: (t: NewTask) => void
  addTasks: (list: NewTask[]) => void
  updateTask: (id: string, t: NewTask) => void
  deleteTask: (id: string) => void
  setStatus: (id: string, status: Status) => void
  cycleStatus: (id: string) => void
  togglePin: (id: string) => void
  updateProfile: (p: Profile) => void
}

const AppContext = createContext<AppState | null>(null)

export function AppProvider({ children }: { children: ReactNode }) {
  const [tasks, setTasks] = useState<Task[]>(SEED_TASKS)
  const [profile, setProfile] = useState<Profile>({
    username: 'taro_yamada',
    email: 'taro.yamada@example.com',
  })

  const value = useMemo<AppState>(
    () => ({
      profile,
      tasks,
      addTask: (t) =>
        setTasks((prev) => [
          {
            ...t,
            id: nextId(),
            status: 'todo',
            pinned: false,
            createdAt: new Date().toISOString(),
          },
          ...prev,
        ]),
      addTasks: (list) =>
        setTasks((prev) => [
          ...list.map((t) => ({
            ...t,
            id: nextId(),
            status: 'todo' as Status,
            pinned: false,
            createdAt: new Date().toISOString(),
          })),
          ...prev,
        ]),
      updateTask: (id, t) =>
        setTasks((prev) =>
          prev.map((task) => (task.id === id ? { ...task, ...t } : task)),
        ),
      deleteTask: (id) =>
        setTasks((prev) => prev.filter((task) => task.id !== id)),
      setStatus: (id, status) =>
        setTasks((prev) =>
          prev.map((task) => (task.id === id ? { ...task, status } : task)),
        ),
      cycleStatus: (id) =>
        setTasks((prev) =>
          prev.map((task) => {
            if (task.id !== id) return task
            const next: Status =
              task.status === 'todo'
                ? 'in_progress'
                : task.status === 'in_progress'
                  ? 'done'
                  : 'todo'
            return { ...task, status: next }
          }),
        ),
      togglePin: (id) =>
        setTasks((prev) =>
          prev.map((task) =>
            task.id === id ? { ...task, pinned: !task.pinned } : task,
          ),
        ),
      updateProfile: (p) => setProfile(p),
    }),
    [tasks, profile],
  )

  return <AppContext.Provider value={value}>{children}</AppContext.Provider>
}

export function useApp() {
  const ctx = useContext(AppContext)
  if (!ctx) throw new Error('useApp must be used within AppProvider')
  return ctx
}

export const StoreProvider = AppProvider
export const useStore = useApp


// ---- display helpers ----

export const PRIORITY_META: Record<
  Priority,
  { label: string; dot: string; text: string }
> = {
  high: { label: '高', dot: 'bg-destructive', text: 'text-destructive' },
  medium: { label: '中', dot: 'bg-chart-3', text: 'text-chart-3' },
  low: { label: '低', dot: 'bg-chart-2', text: 'text-chart-2' },
}

export const STATUS_META: Record<Status, { label: string }> = {
  todo: { label: '未着手' },
  in_progress: { label: '進行中' },
  done: { label: '完了' },
}

export function formatDue(value: string | null) {
  if (!value) return '締切なし'
  const d = new Date(value)
  return d.toLocaleString('ja-JP', {
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
  })
}

export function daysUntil(value: string | null) {
  if (!value) return Infinity
  const now = new Date()
  const d = new Date(value)
  const diff = Math.ceil(
    (d.getTime() - now.setHours(0, 0, 0, 0)) / (1000 * 60 * 60 * 24),
  )
  return diff
}

/** Hours until deadline (negative = overdue). Infinity when no deadline. */
export function hoursUntil(value: string | null) {
  if (!value) return Infinity
  return (new Date(value).getTime() - Date.now()) / (1000 * 60 * 60)
}

/**
 * Default sort: pinned first, then earliest deadline (nulls last),
 * then newest createdAt.
 */
export function defaultSort(a: Task, b: Task) {
  if (a.pinned !== b.pinned) return a.pinned ? -1 : 1
  if (a.due && b.due) {
    const diff = +new Date(a.due) - +new Date(b.due)
    if (diff !== 0) return diff
  } else if (a.due && !b.due) {
    return -1
  } else if (!a.due && b.due) {
    return 1
  }
  return +new Date(b.createdAt) - +new Date(a.createdAt)
}

/** Deadline sort without pin priority (for priority/deadline views). */
export function deadlineSort(a: Task, b: Task) {
  if (a.due && b.due) {
    const diff = +new Date(a.due) - +new Date(b.due)
    if (diff !== 0) return diff
  } else if (a.due && !b.due) {
    return -1
  } else if (!a.due && b.due) {
    return 1
  }
  return +new Date(b.createdAt) - +new Date(a.createdAt)
}
