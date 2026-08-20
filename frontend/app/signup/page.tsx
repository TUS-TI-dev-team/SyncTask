'use client'

import Link from 'next/link'
import { useRouter } from 'next/navigation'
import { useState } from 'react'
import { toast } from 'sonner'
import { AuthShell } from '@/components/layouts/auth-shell'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'

export default function SignupPage() {
  const router = useRouter()
  const [form, setForm] = useState({
    username: '',
    password: '',
    confirm: '',
    email: '',
  })

  const onSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (form.password !== form.confirm) {
      toast.error('パスワードが一致しません')
      return
    }
    const params = new URLSearchParams({ email: form.email })
    router.push(`/signup/otp?${params.toString()}`)
  }

  return (
    <AuthShell
      step="アカウント作成 1 / 2"
      title="アカウントを作成"
      subtitle="基本情報を入力して新しいワークスペースを開始します。"
    >
      <form className="flex flex-col gap-4" onSubmit={onSubmit}>
        <div className="grid gap-2">
          <Label htmlFor="username">ユーザー名</Label>
          <Input
            id="username"
            placeholder="taro_yamada"
            value={form.username}
            onChange={(e) => setForm({ ...form, username: e.target.value })}
            required
          />
        </div>
        <div className="grid gap-2">
          <Label htmlFor="email">メールアドレス</Label>
          <Input
            id="email"
            type="email"
            placeholder="taro@example.com"
            value={form.email}
            onChange={(e) => setForm({ ...form, email: e.target.value })}
            required
          />
        </div>
        <div className="grid gap-2">
          <Label htmlFor="password">パスワード</Label>
          <Input
            id="password"
            type="password"
            placeholder="••••••••"
            value={form.password}
            onChange={(e) => setForm({ ...form, password: e.target.value })}
            required
          />
        </div>
        <div className="grid gap-2">
          <Label htmlFor="confirm">パスワード（確認）</Label>
          <Input
            id="confirm"
            type="password"
            placeholder="••••••••"
            value={form.confirm}
            onChange={(e) => setForm({ ...form, confirm: e.target.value })}
            required
          />
        </div>
        <Button type="submit" className="mt-2 w-full font-semibold">
          続行
        </Button>
      </form>

      <p className="mt-8 text-center text-sm text-muted-foreground">
        すでにアカウントをお持ちですか？{' '}
        <Link href="/login" className="font-semibold text-primary hover:underline">
          ログイン
        </Link>
      </p>
    </AuthShell>
  )
}
