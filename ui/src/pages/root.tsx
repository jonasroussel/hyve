import { Spinner } from '@/components/ui/spinner'
import { useQueryConnectedUser } from '@/hooks/api/users'
import { Outlet } from 'react-router'

export function RootPage() {
	const connectedUserQuery = useQueryConnectedUser()

	if (connectedUserQuery.isLoading) {
		return (
			<div className="flex h-screen items-center justify-center">
				<Spinner className="mb-[20px]" size="lg" />
			</div>
		)
	}

	return <Outlet />
}
