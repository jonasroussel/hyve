import { API } from '@/lib/api'
import { useMutation, useQuery } from '@tanstack/react-query'
import { AxiosError } from 'axios'
import { useLocation, useNavigate } from 'react-router'
import { useLogout } from './auth'

interface User {
	id: string
	username: string
	email?: string
}

export const useQueryConnectedUser = () => {
	const navigate = useNavigate()
	const location = useLocation()
	const logout = useLogout()

	const query = useQuery<unknown, Error, User>({
		queryKey: ['connected-user'],
		queryFn: async () => {
			try {
				const res = await API.get('/accounts/me')

				if (['/register', '/login'].includes(location.pathname)) {
					await navigate('/', { replace: true })
				}

				return res.data
			} catch (err) {
				if (!(err instanceof AxiosError)) throw err

				if (err.response?.data.error === 'NO_ROOT_ACCOUNT') {
					if (location.pathname !== '/register') {
						await navigate('/register')
					}
				} else {
					await logout.mutateAsync()
				}

				throw err['response']
			}
		},
		retry: false,
	})

	return query
}

export const useConnectedUser = () => useQueryConnectedUser().data!

export const useRootRegister = () => {
	return useMutation<User, Error, { username: string; email: string; password: string }>({
		mutationFn: async (data) => {
			const res = await API.post('/accounts/root', data)
			return res.data
		},
	})
}
