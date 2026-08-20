import { AppShell } from '@/components/layouts/app-shell'
import { ProfileEditView } from '@/components/profile/profile-edit-view'

export default function ProfileEditPage() {
  return (
    <AppShell>
      <div className="mx-auto max-w-2xl">
        <ProfileEditView />
      </div>
    </AppShell>
  )
}
