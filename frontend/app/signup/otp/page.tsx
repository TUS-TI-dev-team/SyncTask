'use client'

import { useRouter, useSearchParams } from 'next/navigation'
import { Suspense } from 'react'
import { toast } from 'sonner'
import { AuthShell } from '@/components/layouts/auth-shell'
import { OtpPanel } from '@/components/auth/otp-panel'

function SignupOtpContent() {
  const router = useRouter()
  const params = useSearchParams()
  const email = params.get('email') || 'your@email.com'

  return (
    <AuthShell
      step="アカウント作成 2 / 2"
      title="メールを確認"
      subtitle="送信された確認コードを入力してアカウントを有効化します。認証に成功すると自動的にログインします。"
    >
      <OtpPanel
        email={email}
        onBack={() => router.push('/signup')}
        onExpire={() => router.push('/signup')}
        onSubmit={() => {
          toast.success('アカウントが作成されました')
          router.push('/home')
        }}
      />
      <p className="mt-6 text-center text-xs text-muted-foreground">
        コードが届かない場合は迷惑メールフォルダをご確認ください。
      </p>
    </AuthShell>
  )
}

export default function SignupOtpPage() {
  return (
    <Suspense fallback={null}>
      <SignupOtpContent />
    </Suspense>
  )
}
