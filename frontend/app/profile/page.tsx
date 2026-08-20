import { AppShell } from '@/components/layouts/app-shell'
import { ProfileView } from '@/components/profile/profile-view'

export default function ProfilePage() {
  return (
    <AppShell>
      <div className="mx-auto max-w-2xl">
        <ProfileView />
      </div>
    </AppShell>
  )
}
