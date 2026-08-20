'use client'

import Link from 'next/link'
import { useRouter } from 'next/navigation'
import { useState } from 'react'
import { toast } from 'sonner'
import { AuthShell } from '@/components/layouts/auth-shell'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'

export default function ResetEmailPage() {
  const router = useRouter()
  const [email, setEmail] = useState('')

  const onSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (!email) return
    toast.success('認証メールを送信しました')
    router.push(`/reset-password/otp?email=${encodeURIComponent(email)}`)
  }

  return (
    <AuthShell
      step="パスワードリセット 1 / 3"
      title="パスワードをリセット"
      subtitle="登録済みのメールアドレスを入力すると、認証コードを送信します。"
    >
      <form className="flex flex-col gap-5" onSubmit={onSubmit}>
        <div className="grid gap-2">
          <Label htmlFor="email">メールアドレス</Label>
          <Input
            id="email"
            type="email"
            placeholder="taro@example.com"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required
          />
        </div>
        <Button type="submit" className="w-full font-semibold">
          メール送信
        </Button>
      </form>

      <p className="mt-8 text-center text-sm text-muted-foreground">
        <Link href="/login" className="font-semibold text-primary hover:underline">
          ログイン画面に戻る
        </Link>
      </p>
    </AuthShell>
  )
}
