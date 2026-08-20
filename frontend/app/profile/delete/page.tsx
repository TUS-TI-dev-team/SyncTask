'use client'

import { useRouter } from 'next/navigation'
import { useState } from 'react'
import { AlertTriangle } from 'lucide-react'
import { toast } from 'sonner'
import { AppShell } from '@/components/layouts/app-shell'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from '@/components/ui/alert-dialog'

const MAX_ATTEMPTS = 5
// Mock: any non-empty password other than this "correct" one counts as a failure.
const MOCK_PASSWORD = 'password123'

export default function AccountDeletePage() {
  const router = useRouter()
  const [password, setPassword] = useState('')
  const [attempts, setAttempts] = useState(0)
  const [confirmOpen, setConfirmOpen] = useState(false)

  const onSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (!password) return
    setConfirmOpen(true)
  }

  const executeDelete = () => {
    if (password === MOCK_PASSWORD) {
      toast.success('アカウントを削除しました')
      router.push('/login')
      return
    }

    const next = attempts + 1
    setAttempts(next)
    setPassword('')
    if (next >= MAX_ATTEMPTS) {
      toast.error('認証に5回失敗しました。再ログインが必要です')
      router.push('/login')
      return
    }
    toast.error(`パスワードが正しくありません（${next}/${MAX_ATTEMPTS}）`)
  }

  const remaining = MAX_ATTEMPTS - attempts

  return (
    <AppShell>
      <div className="mx-auto max-w-md">
        <div className="mb-6 flex flex-col gap-1">
          <p className="inline-flex items-center gap-2 font-mono text-xs uppercase tracking-widest text-destructive">
            <span className="h-1.5 w-1.5 bg-destructive" />
            アカウント削除
          </p>
          <h1 className="font-heading text-3xl font-bold tracking-tight">
            本人確認
          </h1>
          <p className="mt-2 text-sm text-muted-foreground leading-relaxed">
            アカウントを削除するには、現在のパスワードを入力してください。
          </p>
        </div>

        <div className="mb-5 flex items-start gap-3 border border-destructive/40 bg-destructive/10 p-4">
          <AlertTriangle className="mt-0.5 h-5 w-5 shrink-0 text-destructive" />
          <div className="text-sm">
            <p className="font-semibold text-destructive">
              この操作は取り消せません
            </p>
            <p className="mt-1 text-muted-foreground leading-relaxed">
              アカウントと所有するすべてのタスクデータが完全に削除され、全セッションが無効化されます。
            </p>
          </div>
        </div>

        <form
          className="flex flex-col gap-5 border border-border bg-card p-6"
          onSubmit={onSubmit}
        >
          <div className="grid gap-2">
            <Label htmlFor="password">現在のパスワード</Label>
            <Input
              id="password"
              type="password"
              placeholder="••••••••"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
            />
            {attempts > 0 && (
              <p className="text-xs text-destructive">
                認証失敗。残り {remaining} 回で強制ログアウトされます。
              </p>
            )}
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
            <Button
              type="submit"
              variant="destructive"
              className="flex-1 font-semibold"
            >
              削除
            </Button>
          </div>
        </form>

        <AlertDialog open={confirmOpen} onOpenChange={setConfirmOpen}>
          <AlertDialogContent>
            <AlertDialogHeader>
              <AlertDialogTitle className="font-heading">
                本当にアカウントを削除しますか？
              </AlertDialogTitle>
              <AlertDialogDescription>
                この操作は取り消せません。アカウントと所有するすべてのタスクデータが完全に削除されます。
              </AlertDialogDescription>
            </AlertDialogHeader>
            <AlertDialogFooter>
              <AlertDialogCancel>キャンセル</AlertDialogCancel>
              <AlertDialogAction
                onClick={executeDelete}
                className="bg-destructive text-white hover:bg-destructive/90"
              >
                削除する
              </AlertDialogAction>
            </AlertDialogFooter>
          </AlertDialogContent>
        </AlertDialog>

        <p className="mt-4 text-center text-xs text-muted-foreground">
          デモ用パスワード: <span className="font-mono">password123</span>
        </p>
      </div>
    </AppShell>
  )
}
