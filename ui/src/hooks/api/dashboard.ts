import { API } from '@/lib/api'
import { useQuery } from '@tanstack/react-query'
import { AxiosError } from 'axios'
import { useNavigate } from 'react-router'

export function useDashboard() {
	const navigate = useNavigate()

	const query = useQuery<unknown, Error, any>({
		queryKey: ['dashboard'],
		queryFn: async () => {
			try {
				const res = await API.get('/dashboard')
				return res.data
			} catch (err) {
				if (!(err instanceof AxiosError)) throw err

				if (err.response?.data.error === 'NO_STORE') {
					await navigate('/onboarding')
				}

				throw err['response']
			}
		},
		retry: false,
	})

	return query
}
