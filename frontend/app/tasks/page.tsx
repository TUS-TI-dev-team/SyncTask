import { AppShell } from '@/components/layouts/app-shell'
import { TasksView } from '@/components/tasks/tasks-view'

export default function TasksPage() {
  return (
    <AppShell>
      <TasksView />
    </AppShell>
  )
}
