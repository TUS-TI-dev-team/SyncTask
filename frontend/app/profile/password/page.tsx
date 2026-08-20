'use client'

import { useRouter } from 'next/navigation'
import { useState } from 'react'
import { toast } from 'sonner'
import { AppShell } from '@/components/layouts/app-shell'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'

const MAX_ATTEMPTS = 5
// Mock: this is treated as the "correct" current password.
const MOCK_PASSWORD = 'password123'

export default function PasswordChangePage() {
  const router = useRouter()
  const [current, setCurrent] = useState('')
  const [next, setNext] = useState('')
  const [confirm, setConfirm] = useState('')
  const [attempts, setAttempts] = useState(0)

  const onSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (current !== MOCK_PASSWORD) {
      const n = attempts + 1
      setAttempts(n)
      setCurrent('')
      if (n >= MAX_ATTEMPTS) {
        toast.error('現在のパスワード認証に5回失敗しました。再ログインが必要です')
        router.push('/login')
        return
      }
      toast.error(`現在のパスワードが正しくありません（${n}/${MAX_ATTEMPTS}）`)
      return
    }
    if (next !== confirm) {
      toast.error('新しいパスワードが一致しません')
      return
    }
    toast.success('パスワードを変更しました。再ログインしてください')
    router.push('/login')
  }

  return (
    <AppShell>
      <div className="mx-auto max-w-md">
        <div className="mb-6 flex flex-col gap-1">
          <p className="inline-flex items-center gap-2 font-mono text-xs uppercase tracking-widest text-primary">
            <span className="h-1.5 w-1.5 bg-primary" />
            セキュリティ
          </p>
          <h1 className="font-heading text-3xl font-bold tracking-tight">
            パスワード変更
          </h1>
          <p className="mt-2 text-sm text-muted-foreground leading-relaxed">
            現在のパスワードと新しいパスワードを入力してください。変更後は全セッションが無効化され、再ログインが必要です。
          </p>
        </div>

        <form
          className="flex flex-col gap-5 border border-border bg-card p-6"
          onSubmit={onSubmit}
        >
          <div className="grid gap-2">
            <Label htmlFor="current">現在のパスワード</Label>
            <Input
              id="current"
              type="password"
              placeholder="••••••••"
              value={current}
              onChange={(e) => setCurrent(e.target.value)}
              required
            />
            {attempts > 0 && (
              <p className="text-xs text-destructive">
                認証失敗。残り {MAX_ATTEMPTS - attempts} 回で強制ログアウトされます。
              </p>
            )}
          </div>
          <div className="grid gap-2">
            <Label htmlFor="next">新しいパスワード</Label>
            <Input
              id="next"
              type="password"
              placeholder="••••••••"
              value={next}
              onChange={(e) => setNext(e.target.value)}
              required
            />
          </div>
          <div className="grid gap-2">
            <Label htmlFor="confirm">新しいパスワード（確認）</Label>
            <Input
              id="confirm"
              type="password"
              placeholder="••••••••"
              value={confirm}
              onChange={(e) => setConfirm(e.target.value)}
              required
            />
          </div>
          <div className="flex gap-3">
            <Button
              type="button"
              variant="outline"
              className="flex-1"
              onClick={() => router.push('/profile')}
            >
              キャンセル
            </Button>
            <Button type="submit" className="flex-1 font-semibold">
              決定
            </Button>
          </div>
        </form>
        <p className="mt-4 text-center text-xs text-muted-foreground">
          デモ用の現在のパスワード: <span className="font-mono">password123</span>
        </p>
      </div>
    </AppShell>
  )
}
