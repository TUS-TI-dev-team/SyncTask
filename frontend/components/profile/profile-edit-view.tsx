'use client'

import { useRouter } from 'next/navigation'
import { useState } from 'react'
import { toast } from 'sonner'
import { useApp } from '@/lib/store'
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

export function ProfileEditView() {
  const { profile, updateProfile } = useApp()
  const router = useRouter()
  const [username, setUsername] = useState(profile.username)
  const [email, setEmail] = useState(profile.email)
  const [confirmOpen, setConfirmOpen] = useState(false)

  const emailChanged = email.trim() !== profile.email

  const save = () => {
    updateProfile({ username: username.trim(), email: profile.email })
    toast.success('プロフィールを更新しました')
    router.push('/profile')
  }

  const onSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (!emailChanged && username.trim() === profile.username) {
      toast.error('現在のユーザー名と同じです')
      return
    }
    if (emailChanged) {
      setConfirmOpen(true)
    } else {
      save()
    }
  }

  return (
    <div className="flex flex-col gap-6">
      <div className="flex flex-col gap-1">
        <p className="inline-flex items-center gap-2 font-mono text-xs uppercase tracking-widest text-primary">
          <span className="h-1.5 w-1.5 bg-primary" />
          プロフィール編集
        </p>
        <h1 className="font-heading text-3xl font-bold tracking-tight">
          アカウント情報を変更
        </h1>
      </div>

      <form
        className="flex flex-col gap-5 border border-border bg-card p-6"
        onSubmit={onSubmit}
      >
        <div className="grid gap-2">
          <Label htmlFor="username">ユーザー名</Label>
          <Input
            id="username"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            required
          />
        </div>
        <div className="grid gap-2">
          <Label htmlFor="email">メールアドレス</Label>
          <Input
            id="email"
            type="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required
          />
          {emailChanged && (
            <p className="text-xs text-chart-3">
              メールアドレスを変更すると確認コードの入力が必要です。
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
          <Button type="submit" className="flex-1 font-semibold">
            決定
          </Button>
        </div>
      </form>

      <AlertDialog open={confirmOpen} onOpenChange={setConfirmOpen}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle className="font-heading">
              確認メールを送信しますか？
            </AlertDialogTitle>
            <AlertDialogDescription className="whitespace-pre-line">
              {email} に確認コードを送信します。{"\n"}よろしいですか？{"\n"}またメールアドレス変更後は再度ログインが求められます。
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>キャンセル</AlertDialogCancel>
            <AlertDialogAction
              onClick={() => {
                const params = new URLSearchParams({
                  email,
                  username: username.trim(),
                })
                router.push(`/profile/otp?${params.toString()}`)
              }}
            >
              送信する
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </div>
  )
}
