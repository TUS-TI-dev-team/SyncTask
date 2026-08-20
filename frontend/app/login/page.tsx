'use client'

import Link from 'next/link'
import { useRouter } from 'next/navigation'
import { useState } from 'react'
import { Eye, EyeOff } from 'lucide-react'
import { AuthShell } from '@/components/layouts/auth-shell'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'

export default function LoginPage() {
  const router = useRouter()
  const [show, setShow] = useState(false)

  return (
    <AuthShell
      step="認証"
      title="おかえりなさい"
      subtitle="メールアドレスとパスワードを入力してワークスペースにアクセスします。"
    >
      <form
        className="flex flex-col gap-5"
        onSubmit={(e) => {
          e.preventDefault()
          router.push('/home')
        }}
      >
        <div className="grid gap-2">
          <Label htmlFor="email">メールアドレス</Label>
          <Input
            id="email"
            type="email"
            placeholder="taro@example.com"
            autoComplete="email"
            required
          />
        </div>

        <div className="grid gap-2">
          <div className="flex items-center justify-between">
            <Label htmlFor="password">パスワード</Label>
            <Link
              href="/reset-password"
              className="text-xs font-medium text-primary hover:underline"
            >
              パスワードを忘れた？
            </Link>
          </div>
          <div className="relative">
            <Input
              id="password"
              type={show ? 'text' : 'password'}
              placeholder="••••••••"
              autoComplete="current-password"
              required
            />
            <button
              type="button"
              onClick={() => setShow((v) => !v)}
              className="absolute right-3 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground"
              aria-label={show ? 'パスワードを隠す' : 'パスワードを表示'}
            >
              {show ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
            </button>
          </div>
        </div>

        <Button type="submit" className="mt-1 w-full font-semibold">
          ログイン
        </Button>
      </form>

      <p className="mt-8 text-center text-sm text-muted-foreground">
        アカウントをお持ちでない方は{' '}
        <Link href="/signup" className="font-semibold text-primary hover:underline">
          アカウント作成
        </Link>
      </p>
    </AuthShell>
  )
}
