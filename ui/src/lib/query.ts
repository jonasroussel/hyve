import { QueryClient } from '@tanstack/react-query'

export const queryClient = new QueryClient({
	defaultOptions: {
		queries: {
			staleTime: 5 * 60_000,
			retry: (failureCount) => failureCount < 3,
			retryDelay: (attemptIndex) => Math.min(1000 * (attemptIndex + 1), 5_000),
		},
	},
})
