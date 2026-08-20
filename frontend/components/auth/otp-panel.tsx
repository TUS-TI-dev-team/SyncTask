'use client'

import { useEffect, useRef, useState } from 'react'
import { ArrowLeft, RotateCw } from 'lucide-react'
import { toast } from 'sonner'
import { OtpInput } from '@/components/auth/otp-input'
import { Button } from '@/components/ui/button'
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from '@/components/ui/alert-dialog'

/**
 * Mask an email as: first 4 chars + 10 fixed-width mask chars + domain.
 * e.g. taro.yamada@example.com -> taro••••••••••@example.com
 */
export function maskEmail(email: string) {
  const [local, domain] = email.split('@')
  if (!domain) return email
  const head = local.slice(0, 4)
  return `${head}${'•'.repeat(10)}@${domain}`
}

const OTP_LENGTH = 8
const COOLDOWN = 60
const OTP_SECONDS = 5 * 60
const SESSION_SECONDS = 15 * 60

export function OtpPanel({
  email,
  onBack,
  onSubmit,
  onExpire,
  submitLabel = '決定',
  backLabel = '戻る',
}: {
  email: string
  onBack: () => void
  onSubmit: (code: string) => void
  onExpire: () => void
  submitLabel?: string
  backLabel?: string
}) {
  const [code, setCode] = useState('')
  const [cooldown, setCooldown] = useState(COOLDOWN)
  const [otpRemaining, setOtpRemaining] = useState(OTP_SECONDS)
  const [sessionRemaining, setSessionRemaining] = useState(SESSION_SECONDS)
  const [expiredOpen, setExpiredOpen] = useState(false)
  const timer = useRef<ReturnType<typeof setInterval> | null>(null)
  const sessionExpired = useRef(false)

  useEffect(() => {
    timer.current = setInterval(() => {
      setCooldown((c) => {
        return Math.max(0, c - 1)
      })
      setOtpRemaining((s) => {
        const next = Math.max(0, s - 1)
        if (s > 0 && next === 0) {
          setExpiredOpen(true)
        }
        return next
      })
      setSessionRemaining((s) => {
        const next = Math.max(0, s - 1)
        if (s > 0 && next === 0 && !sessionExpired.current) {
          sessionExpired.current = true
          onExpire()
        }
        return next
      })
    }, 1000)
    return () => {
      if (timer.current) clearInterval(timer.current)
    }
  }, [onExpire])

  const resend = () => {
    if (cooldown > 0) return
    setCooldown(COOLDOWN)
    setOtpRemaining(OTP_SECONDS)
    setExpiredOpen(false)
    toast.success('認証コードを再送信しました')
  }

  const submit = (e: React.FormEvent) => {
    e.preventDefault()
    if (code.length < OTP_LENGTH) {
      toast.error(`${OTP_LENGTH}桁の認証コードを入力してください`)
      return
    }
    if (otpRemaining === 0) {
      setExpiredOpen(true)
      return
    }
    onSubmit(code)
  }

  const sessionMinutes = String(Math.floor(sessionRemaining / 60)).padStart(2, '0')
  const sessionSeconds = String(sessionRemaining % 60).padStart(2, '0')

  return (
    <form className="flex flex-col gap-6" onSubmit={submit}>
      <div className="border border-border bg-secondary/40 px-4 py-3">
        <p className="font-mono text-xs uppercase tracking-widest text-muted-foreground">
          送信先
        </p>
        <p className="mt-1 break-all font-mono text-sm text-foreground">
          {maskEmail(email)}
        </p>
      </div>

      <div className="flex items-center justify-between border border-border bg-card px-3 py-2 font-mono text-xs">
        <span className="uppercase tracking-widest text-muted-foreground">
          セッション残り時間
        </span>
        <span className="font-bold text-primary">
          {sessionMinutes}:{sessionSeconds}
        </span>
      </div>

      <OtpInput length={OTP_LENGTH} value={code} onChange={setCode} />

      <div className="flex items-center justify-between">
        <p className="text-xs text-muted-foreground">
          英数字8桁（大文字・小文字の区別なし）
        </p>
        <button
          type="button"
          onClick={resend}
          disabled={cooldown > 0}
          className="inline-flex items-center gap-1.5 text-xs font-medium text-primary transition-colors hover:underline disabled:cursor-not-allowed disabled:text-muted-foreground disabled:no-underline"
        >
          <RotateCw className="h-3.5 w-3.5" />
          {cooldown > 0 ? `再送信 (${cooldown}s)` : '再送信'}
        </button>
      </div>

      <div className="flex gap-3">
        <Button
          type="button"
          variant="outline"
          className="flex-1 gap-2"
          onClick={onBack}
        >
          <ArrowLeft className="h-4 w-4" />
          {backLabel}
        </Button>
        <Button type="submit" className="flex-1 font-semibold">
          {submitLabel}
        </Button>
      </div>


      <AlertDialog open={expiredOpen} onOpenChange={setExpiredOpen}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>認証コードの有効期限が切れました</AlertDialogTitle>
            <AlertDialogDescription>
              認証コードの発行から5分が経過しました。新しいコードを再送信してください。
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogAction onClick={() => setExpiredOpen(false)}>
              閉じる
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </form>
  )
}
