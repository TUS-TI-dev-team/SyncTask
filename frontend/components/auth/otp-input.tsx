'use client'

import { useRef, type KeyboardEvent, type ClipboardEvent } from 'react'
import { cn } from '@/lib/utils'

const SANITIZE = /[^a-zA-Z0-9]/g

export function OtpInput({
  length = 8,
  value,
  onChange,
}: {
  length?: number
  value: string
  onChange: (v: string) => void
}) {
  const refs = useRef<(HTMLInputElement | null)[]>([])

  const setChar = (index: number, char: string) => {
    const chars = value.split('')
    chars[index] = char
    const next = chars.join('').slice(0, length)
    onChange(next)
  }

  const handleChange = (index: number, raw: string) => {
    const char = raw.replace(SANITIZE, '').slice(-1).toUpperCase()
    if (!char) return
    setChar(index, char)
    if (index < length - 1) refs.current[index + 1]?.focus()
  }

  const handleKeyDown = (index: number, e: KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Backspace') {
      if (value[index]) {
        setChar(index, '')
      } else if (index > 0) {
        refs.current[index - 1]?.focus()
        setChar(index - 1, '')
      }
    }
    if (e.key === 'ArrowLeft' && index > 0) refs.current[index - 1]?.focus()
    if (e.key === 'ArrowRight' && index < length - 1)
      refs.current[index + 1]?.focus()
  }

  const handlePaste = (e: ClipboardEvent<HTMLInputElement>) => {
    e.preventDefault()
    const pasted = e.clipboardData
      .getData('text')
      .replace(SANITIZE, '')
      .toUpperCase()
      .slice(0, length)
    if (pasted) {
      onChange(pasted)
      refs.current[Math.min(pasted.length, length - 1)]?.focus()
    }
  }

  return (
    <div className="grid grid-cols-8 gap-1.5 sm:gap-2">
      {Array.from({ length }).map((_, i) => (
        <input
          key={i}
          ref={(el) => {
            refs.current[i] = el
          }}
          inputMode="text"
          autoCapitalize="characters"
          maxLength={1}
          value={value[i] ?? ''}
          onChange={(e) => handleChange(i, e.target.value)}
          onKeyDown={(e) => handleKeyDown(i, e)}
          onPaste={handlePaste}
          aria-label={`認証コード ${i + 1}桁目`}
          className={cn(
            'h-12 w-full min-w-0 border border-input bg-card text-center font-heading text-lg font-bold uppercase outline-none transition-colors sm:h-14 sm:text-xl',
            'focus:border-primary focus:ring-2 focus:ring-primary/30',
            value[i] && 'border-primary',
          )}
        />
      ))}
    </div>
  )
}
