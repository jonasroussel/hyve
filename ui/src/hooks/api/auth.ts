import { API } from '@/lib/api'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { useNavigate } from 'react-router'

export function useLogout() {
	const navigate = useNavigate()

	return useMutation<void, Error, { noRedirect?: boolean } | void>({
		mutationFn: async ({ noRedirect } = {}) => {
			await API.post('/auth/logout').catch(() => {})
			if (!noRedirect) await navigate('/login')
		},
	})
}

export function useLogin() {
	const queryClient = useQueryClient()

	return useMutation<void, Error, { username_email: string; password: string }>({
		mutationFn: async (data) => {
			await API.post('/auth/login', data)

			await queryClient.invalidateQueries({ queryKey: ['connected-user'] })
		},
	})
}
