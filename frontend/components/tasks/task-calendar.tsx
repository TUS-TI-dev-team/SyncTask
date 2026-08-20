'use client'

import { useEffect, useLayoutEffect, useMemo, useRef, useState } from 'react'
import {
  CalendarRange,
  ChevronLeft,
  ChevronRight,
  CircleCheckBig,
  CircleDashed,
  CircleDot,
  Rows3,
} from 'lucide-react'
import { STATUS_META, type Status, type Task } from '@/lib/store'
import { cn } from '@/lib/utils'
import { DateDetailDialog } from '@/components/tasks/date-detail-dialog'

const WEEKDAYS = ['日', '月', '火', '水', '木', '金', '土']

type CalMode = 'all' | 'week'

// How many weeks to render before / after the anchor week in 全体表示.
const RANGE_BACK = 26
const RANGE_FWD = 78

function dayKey(d: Date) {
  return `${d.getFullYear()}-${d.getMonth()}-${d.getDate()}`
}

function startOfWeek(d: Date) {
  const s = new Date(d)
  s.setDate(d.getDate() - d.getDay())
  s.setHours(0, 0, 0, 0)
  return s
}

function addDays(d: Date, n: number) {
  const r = new Date(d)
  r.setDate(d.getDate() + n)
  return r
}

type Week = {
  start: Date
  days: Date[]
  /** Representative month/year (Thursday of the week) used for labels & shading. */
  repMonth: number
  repYear: number
}

export function TaskCalendar({
  tasks,
  onSelectTask,
  onCreateForDate,
  onSetStatus,
  onTogglePin,
}: {
  tasks: Task[]
  onSelectTask: (task: Task) => void
  onCreateForDate: (date: Date) => void
  onSetStatus: (id: string, status: Status) => void
  onTogglePin: (id: string) => void
}) {
  const [mode, setMode] = useState<CalMode>('all')
  // Anchor for 全体表示 (fixed at mount) and cursor for 週表示 navigation.
  const [anchor] = useState(() => startOfWeek(new Date()))
  const [cursor, setCursor] = useState(() => new Date())
  const [selectedDate, setSelectedDate] = useState<Date | null>(null)

  const scrollRef = useRef<HTMLDivElement>(null)
  const rowRefs = useRef<Map<number, HTMLDivElement>>(new Map())
  const focusedIndexRef = useRef(0)

  const tasksByDay = useMemo(() => {
    const map = new Map<string, Task[]>()
    for (const t of tasks) {
      if (!t.due) continue
      const d = new Date(t.due)
      const key = dayKey(d)
      if (!map.has(key)) map.set(key, [])
      map.get(key)!.push(t)
    }
    for (const list of map.values()) {
      list.sort((a, b) => {
        if (a.pinned !== b.pinned) return a.pinned ? -1 : 1
        return +new Date(a.due!) - +new Date(b.due!)
      })
    }
    return map
  }, [tasks])

  // Continuous list of weeks for 全体表示.
  const weeks = useMemo<Week[]>(() => {
    const list: Week[] = []
    for (let w = -RANGE_BACK; w <= RANGE_FWD; w++) {
      const start = addDays(anchor, w * 7)
      const days: Date[] = []
      for (let i = 0; i < 7; i++) days.push(addDays(start, i))
      const thu = days[4]
      list.push({
        start,
        days,
        repMonth: thu.getMonth(),
        repYear: thu.getFullYear(),
      })
    }
    return list
  }, [anchor])

  const todayWeekIndex = useMemo(
    () => weeks.findIndex((w) => dayKey(w.start) === dayKey(startOfWeek(new Date()))),
    [weeks],
  )

  // Which week is at the top of the scroll viewport (drives the header label).
  const [focusedIndex, setFocusedIndex] = useState(
    todayWeekIndex >= 0 ? todayWeekIndex : RANGE_BACK,
  )
  useEffect(() => {
    focusedIndexRef.current = focusedIndex
  }, [focusedIndex])

  const scrollToWeek = (index: number, behavior: ScrollBehavior = 'smooth') => {
    const c = scrollRef.current
    const row = rowRefs.current.get(index)
    if (c && row) c.scrollTo({ top: row.offsetTop, behavior })
  }

  // On first paint of 全体表示, jump to today's week.
  useLayoutEffect(() => {
    if (mode !== 'all') return
    if (todayWeekIndex >= 0) scrollToWeek(todayWeekIndex, 'auto')
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [mode])

  const handleScroll = () => {
    const c = scrollRef.current
    if (!c) return
    const st = c.scrollTop
    let current = 0
    for (const [idx, el] of rowRefs.current) {
      if (el.offsetTop <= st + 8) current = Math.max(current, idx)
    }
    if (current !== focusedIndexRef.current) {
      focusedIndexRef.current = current
      setFocusedIndex(current)
    }
  }

  const navigateAll = (dir: -1 | 1) => {
    const focused = weeks[focusedIndexRef.current] ?? weeks[RANGE_BACK]
    let nextMonth = focused.repMonth + dir
    let nextYear = focused.repYear
    if (nextMonth < 0) {
      nextMonth = 11
      nextYear -= 1
    } else if (nextMonth > 11) {
      nextMonth = 0
      nextYear += 1
    }
    const idx = weeks.findIndex(
      (w) => w.repYear === nextYear && w.repMonth === nextMonth,
    )
    if (idx >= 0) {
      focusedIndexRef.current = idx
      setFocusedIndex(idx)
      scrollToWeek(idx)
    }
  }

  // ---- 週表示 (single week) ----
  const weekDays = useMemo(() => {
    const start = startOfWeek(cursor)
    return Array.from({ length: 7 }, (_, i) => addDays(start, i))
  }, [cursor])

  const today = new Date()
  const isToday = (d: Date) => dayKey(d) === dayKey(today)

  const label =
    mode === 'all'
      ? `${weeks[focusedIndex]?.repYear ?? today.getFullYear()}年 ${(weeks[focusedIndex]?.repMonth ?? today.getMonth()) + 1}月`
      : `${startOfWeek(cursor).getFullYear()}年 ${startOfWeek(cursor).getMonth() + 1}/${startOfWeek(cursor).getDate()} – ${weekDays[6].getMonth() + 1}/${weekDays[6].getDate()}`

  const navigateWeek = (dir: -1 | 1) => setCursor(addDays(cursor, dir * 7))

  const goToday = () => {
    if (mode === 'all') {
      if (todayWeekIndex >= 0) {
        focusedIndexRef.current = todayWeekIndex
        setFocusedIndex(todayWeekIndex)
        scrollToWeek(todayWeekIndex)
      }
    } else {
      setCursor(new Date())
    }
  }

  const selectedTasks = selectedDate
    ? (tasksByDay.get(dayKey(selectedDate)) ?? [])
    : []

  const maxBars = mode === 'all' ? 3 : 12

  const renderTaskBar = (t: Task) => (
    <span
      key={t.id}
      role="button"
      tabIndex={0}
      onClick={(e) => {
        e.stopPropagation()
        onSelectTask(t)
      }}
      onKeyDown={(e) => {
        if (e.key === 'Enter') {
          e.stopPropagation()
          onSelectTask(t)
        }
      }}
      className={cn(
        'flex items-center gap-1 truncate border-l-2 bg-secondary/60 px-1 py-0.5 text-left text-[11px] hover:bg-secondary',
        t.status === 'done' && 'text-muted-foreground line-through opacity-60',
      )}
      style={{
        borderColor: `var(--${
          t.priority === 'high'
            ? 'destructive'
            : t.priority === 'medium'
              ? 'chart-3'
              : 'chart-2'
        })`,
      }}
    >
      <span
        onClick={(e) => {
          e.stopPropagation()
          onSetStatus(
            t.id,
            t.status === 'todo'
              ? 'in_progress'
              : t.status === 'in_progress'
                ? 'done'
                : 'todo',
          )
        }}
        className={cn(
          'shrink-0 cursor-pointer',
          t.status === 'todo'
            ? 'text-muted-foreground'
            : t.status === 'in_progress'
              ? 'text-chart-3'
              : 'text-primary',
        )}
        title={`${STATUS_META[t.status].label}（クリックで変更）`}
      >
        {t.status === 'todo' ? (
          <CircleDashed className="h-3 w-3" />
        ) : t.status === 'in_progress' ? (
          <CircleDot className="h-3 w-3" />
        ) : (
          <CircleCheckBig className="h-3 w-3" />
        )}
      </span>
      <span className="truncate">{t.title}</span>
    </span>
  )

  const renderDayCell = (
    d: Date,
    opts: { inFocusMonth: boolean; heightClass: string; colIndex: number },
  ) => {
    const dayTasks = tasksByDay.get(dayKey(d)) ?? []
    const firstOfMonth = d.getDate() === 1
    return (
      <button
        key={dayKey(d)}
        onClick={() => setSelectedDate(d)}
        className={cn(
          'flex flex-col gap-1 border-b border-r border-border p-1.5 text-left transition-colors hover:bg-secondary/40',
          opts.colIndex === 0 && 'border-l',
          opts.heightClass,
        )}
      >
        <div className="flex items-center gap-1">
          <div
            className={cn(
              'flex h-6 min-w-6 items-center justify-center px-1 text-xs',
              isToday(d)
                ? 'bg-primary font-bold text-primary-foreground'
                : opts.inFocusMonth
                  ? 'text-foreground'
                  : 'text-muted-foreground/50',
            )}
          >
            {d.getDate()}
          </div>
          {firstOfMonth && (
            <span className="bg-secondary px-1 text-[10px] font-medium text-muted-foreground">
              {d.getMonth() + 1}月
            </span>
          )}
        </div>
        <div className="flex flex-col gap-1">
          {dayTasks.slice(0, maxBars).map(renderTaskBar)}
          {dayTasks.length > maxBars && (
            <span className="px-1 text-[10px] text-muted-foreground">
              +{dayTasks.length - maxBars} 件
            </span>
          )}
        </div>
      </button>
    )
  }

  return (
    <div className="border border-border bg-card">
      {/* Header */}
      <div className="flex flex-col gap-3 border-b border-border p-4 sm:flex-row sm:items-center sm:justify-between">
        <h2 className="font-heading text-lg font-bold">{label}</h2>
        <div className="flex items-center gap-2">
          {/* mode toggle */}
          <div className="flex border border-border">
            <button
              onClick={() => setMode('all')}
              className={cn(
                'flex items-center gap-1.5 px-2.5 py-1.5 text-xs font-medium transition-colors',
                mode === 'all'
                  ? 'bg-primary text-primary-foreground'
                  : 'text-muted-foreground hover:text-foreground',
              )}
            >
              <Rows3 className="h-3.5 w-3.5" />
              全体
            </button>
            <button
              onClick={() => setMode('week')}
              className={cn(
                'flex items-center gap-1.5 px-2.5 py-1.5 text-xs font-medium transition-colors',
                mode === 'week'
                  ? 'bg-primary text-primary-foreground'
                  : 'text-muted-foreground hover:text-foreground',
              )}
            >
              <CalendarRange className="h-3.5 w-3.5" />
              週
            </button>
          </div>
          {/* nav */}
          <div className="flex items-center gap-1">
            <button
              onClick={() => (mode === 'all' ? navigateAll(-1) : navigateWeek(-1))}
              className="grid h-8 w-8 place-items-center border border-border text-muted-foreground hover:bg-secondary hover:text-foreground"
              aria-label={mode === 'all' ? '前の月へ' : '前週'}
            >
              <ChevronLeft className="h-4 w-4" />
            </button>
            <button
              onClick={goToday}
              className="border border-border px-3 py-1.5 text-xs font-medium text-muted-foreground hover:bg-secondary hover:text-foreground"
            >
              今日
            </button>
            <button
              onClick={() => (mode === 'all' ? navigateAll(1) : navigateWeek(1))}
              className="grid h-8 w-8 place-items-center border border-border text-muted-foreground hover:bg-secondary hover:text-foreground"
              aria-label={mode === 'all' ? '次の月へ' : '翌週'}
            >
              <ChevronRight className="h-4 w-4" />
            </button>
          </div>
        </div>
      </div>

      {/* Weekday headers */}
      <div className="grid grid-cols-7 border-b border-border">
        {WEEKDAYS.map((w, i) => (
          <div
            key={w}
            className={cn(
              'py-2 text-center text-xs font-medium',
              i === 0 && 'text-destructive',
              i === 6 && 'text-chart-2',
              i !== 0 && i !== 6 && 'text-muted-foreground',
            )}
          >
            {w}
          </div>
        ))}
      </div>

      {mode === 'all' ? (
        // Continuous, vertically scrollable list of week rows.
        <div
          ref={scrollRef}
          onScroll={handleScroll}
          className="relative h-[26rem] overflow-y-auto"
          style={{ scrollbarGutter: 'stable' }}
        >
          {weeks.map((week, wi) => (
            <div
              key={wi}
              ref={(el) => {
                if (el) rowRefs.current.set(wi, el)
                else rowRefs.current.delete(wi)
              }}
              className="grid grid-cols-7"
            >
              {week.days.map((d, ci) =>
                renderDayCell(d, {
                  inFocusMonth: d.getMonth() === week.repMonth,
                  heightClass: 'min-h-24',
                  colIndex: ci,
                }),
              )}
            </div>
          ))}
        </div>
      ) : (
        // Single week (tall cells).
        <div className="grid grid-cols-7">
          {weekDays.map((d, ci) =>
            renderDayCell(d, {
              inFocusMonth: true,
              heightClass: 'min-h-56',
              colIndex: ci,
            }),
          )}
        </div>
      )}

      <DateDetailDialog
        date={selectedDate}
        tasks={selectedTasks}
        onOpenChange={(o) => !o && setSelectedDate(null)}
        onSelectTask={(t) => {
          onSelectTask(t)
        }}
        onCreate={(date) => {
          onCreateForDate(date)
        }}
        onSetStatus={onSetStatus}
        onTogglePin={onTogglePin}
      />
    </div>
  )
}
