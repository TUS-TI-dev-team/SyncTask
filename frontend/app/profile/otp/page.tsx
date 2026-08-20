'use client'

import { useRouter, useSearchParams } from 'next/navigation'
import { Suspense } from 'react'
import { toast } from 'sonner'
import { AppShell } from '@/components/layouts/app-shell'
import { OtpPanel } from '@/components/auth/otp-panel'
import { useApp } from '@/lib/store'

function AccountOtpContent() {
  const router = useRouter()
  const params = useSearchParams()
  const { profile } = useApp()
  const email = params.get('email') || profile.email

  return (
    <div className="mx-auto max-w-md">
      <div className="mb-6 flex flex-col gap-1">
        <p className="inline-flex items-center gap-2 font-mono text-xs uppercase tracking-widest text-primary">
          <span className="h-1.5 w-1.5 bg-primary" />
          メールアドレス確認
        </p>
        <h1 className="font-heading text-3xl font-bold tracking-tight">
          認証コードを入力
        </h1>
        <p className="mt-2 text-sm text-muted-foreground leading-relaxed">
          新しいメールアドレス宛の認証コードを入力してください。変更が確定すると全セッションが無効化され、再ログインが必要です。
        </p>
      </div>

      <div className="border border-border bg-card p-6">
        <OtpPanel
          email={email}
          onBack={() => router.push('/profile/edit')}
          onExpire={() => router.push('/profile/edit')}
          onSubmit={() => {
            toast.success('メールアドレスを変更しました。再ログインしてください')
            router.push('/login')
          }}
        />
      </div>
    </div>
  )
}

export default function AccountOtpPage() {
  return (
    <AppShell>
      <Suspense fallback={null}>
        <AccountOtpContent />
      </Suspense>
    </AppShell>
  )
}
