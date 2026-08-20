'use client'

import { useRouter, useSearchParams } from 'next/navigation'
import { Suspense } from 'react'
import { AuthShell } from '@/components/layouts/auth-shell'
import { OtpPanel } from '@/components/auth/otp-panel'

function ResetOtpContent() {
  const router = useRouter()
  const params = useSearchParams()
  const email = params.get('email') || 'taro.yamada@example.com'

  return (
    <AuthShell
      step="パスワードリセット 2 / 3"
      title="認証コードを入力"
      subtitle="送信された認証コードを入力してください。認証に成功すると新しいパスワードを設定できます。"
    >
      <OtpPanel
        email={email}
        backLabel="メールアドレス入力に戻る"
        onBack={() => router.push('/reset-password')}
        onExpire={() => router.push('/reset-password')}
        onSubmit={() => router.push('/reset-password/new')}
      />
    </AuthShell>
  )
}

export default function ResetOtpPage() {
  return (
    <Suspense fallback={null}>
      <ResetOtpContent />
    </Suspense>
  )
}
