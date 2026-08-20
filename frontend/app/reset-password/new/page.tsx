'use client'

import { useRouter } from 'next/navigation'
import { useEffect, useState } from 'react'
import { AlertTriangle } from 'lucide-react'
import { toast } from 'sonner'
import { AuthShell } from '@/components/layouts/auth-shell'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'

const SESSION_SECONDS = 15 * 60

export default function NewPasswordPage() {
  const router = useRouter()
  const [password, setPassword] = useState('')
  const [confirm, setConfirm] = useState('')
  const [remaining, setRemaining] = useState(SESSION_SECONDS)

  useEffect(() => {
    const timer = setInterval(() => {
      setRemaining((s) => Math.max(0, s - 1))
    }, 1000)
    return () => clearInterval(timer)
  }, [])

  const expired = remaining === 0
  const mm = String(Math.floor(remaining / 60)).padStart(2, '0')
  const ss = String(remaining % 60).padStart(2, '0')

  const onSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (password !== confirm) {
      toast.error('パスワードが一致しません')
      return
    }
    toast.success('パスワードを再設定しました')
    router.push('/login')
  }

  return (
    <AuthShell
      step="パスワードリセット 3 / 3"
      title="新しいパスワード"
      subtitle="新しいパスワードを入力して再設定を完了します。"
    >
      {expired ? (
        <div className="flex flex-col gap-5">
          <div className="flex items-start gap-3 border border-destructive/40 bg-destructive/10 p-4">
            <AlertTriangle className="mt-0.5 h-5 w-5 shrink-0 text-destructive" />
            <div className="text-sm">
              <p className="font-semibold text-destructive">セッションの有効期限が切れました</p>
              <p className="mt-1 text-muted-foreground leading-relaxed">
                15分以内に操作が完了しなかったため、最初からやり直してください。
              </p>
            </div>
          </div>
          <Button
            className="w-full font-semibold"
            onClick={() => router.push('/reset-password')}
          >
            メールアドレス入力に戻る
          </Button>
        </div>
      ) : (
        <form className="flex flex-col gap-5" onSubmit={onSubmit}>
          <div className="flex items-center justify-between border border-border bg-card px-3 py-2 font-mono text-xs">
            <span className="uppercase tracking-widest text-muted-foreground">
              残り時間
            </span>
            <span className="font-bold text-primary">
              {mm}:{ss}
            </span>
          </div>
          <div className="grid gap-2">
            <Label htmlFor="password">新しいパスワード</Label>
            <Input
              id="password"
              type="password"
              placeholder="••••••••"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
            />
          </div>
          <div className="grid gap-2">
            <Label htmlFor="confirm">パスワード（確認）</Label>
            <Input
              id="confirm"
              type="password"
              placeholder="••••••••"
              value={confirm}
              onChange={(e) => setConfirm(e.target.value)}
              required
            />
          </div>
          <Button type="submit" className="w-full font-semibold">
            決定
          </Button>
        </form>
      )}
    </AuthShell>
  )
}
