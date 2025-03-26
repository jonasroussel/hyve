import { useDashboard } from '@/hooks/api/dashboard'

export function DashboardPage() {
	useDashboard()

	return <div>Dashboard</div>
}
