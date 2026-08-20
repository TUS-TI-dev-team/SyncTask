'use client'

import Link from 'next/link'
import { AtSign, KeyRound, Pencil, Trash2, User } from 'lucide-react'
import { useApp } from '@/lib/store'
import { Button } from '@/components/ui/button'

export function ProfileView() {
  const { profile } = useApp()

  return (
    <div className="flex flex-col gap-6">
      <div className="flex flex-col gap-1">
        <p className="inline-flex items-center gap-2 font-mono text-xs uppercase tracking-widest text-primary">
          <span className="h-1.5 w-1.5 bg-primary" />
          アカウント
        </p>
        <h1 className="font-heading text-3xl font-bold tracking-tight">
          プロフィール
        </h1>
      </div>

      {/* Identity banner */}
      <div className="relative overflow-hidden border border-border bg-card">
        <div className="bg-grid absolute inset-0 opacity-10" />
        <div className="relative flex items-center gap-5 p-6">
          <span className="grid h-20 w-20 shrink-0 place-items-center bg-primary font-heading text-3xl font-bold text-primary-foreground">
            {profile.username.charAt(0).toUpperCase()}
          </span>
          <div className="min-w-0">
            <h2 className="font-heading text-2xl font-bold">
              {profile.username}
            </h2>
            <p className="truncate text-sm text-muted-foreground">
              {profile.email}
            </p>
          </div>
        </div>
      </div>

      {/* Detail rows */}
      <div className="grid gap-px border border-border bg-border sm:grid-cols-2">
        <div className="flex items-center gap-3 bg-card p-4">
          <span className="grid h-10 w-10 place-items-center bg-secondary text-muted-foreground">
            <User className="h-5 w-5" />
          </span>
          <div>
            <p className="text-xs uppercase tracking-wide text-muted-foreground">
              ユーザー名
            </p>
            <p className="font-medium">{profile.username}</p>
          </div>
        </div>
        <div className="flex items-center gap-3 bg-card p-4">
          <span className="grid h-10 w-10 place-items-center bg-secondary text-muted-foreground">
            <AtSign className="h-5 w-5" />
          </span>
          <div className="min-w-0">
            <p className="text-xs uppercase tracking-wide text-muted-foreground">
              メールアドレス
            </p>
            <p className="truncate font-medium">{profile.email}</p>
          </div>
        </div>
      </div>

      {/* Actions */}
      <div className="flex flex-col gap-3 sm:flex-row">
        <Button
          nativeButton={false}
          className="gap-2 font-semibold"
          render={<Link href="/profile/edit" />}
        >
          <Pencil className="h-4 w-4" />
          プロフィール変更
        </Button>
        <Button
          nativeButton={false}
          variant="outline"
          className="gap-2"
          render={<Link href="/profile/password" />}
        >
          <KeyRound className="h-4 w-4" />
          パスワード変更
        </Button>
      </div>

      {/* Danger zone */}
      <div className="mt-4 border border-destructive/30 bg-destructive/5 p-5">
        <h3 className="font-heading font-semibold text-destructive">
          アカウントの削除
        </h3>
        <p className="mt-1 text-sm text-muted-foreground leading-relaxed">
          アカウントを削除すると、すべてのタスクとデータが完全に失われます。この操作は取り消せません。
        </p>
        <Button
          nativeButton={false}
          variant="destructive"
          className="mt-4 gap-2"
          render={<Link href="/profile/delete" />}
        >
          <Trash2 className="h-4 w-4" />
          アカウントを削除
        </Button>
      </div>
    </div>
  )
}
