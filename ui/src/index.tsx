import './index.css'

import { queryClient } from '@/lib/query'
import { QueryClientProvider } from '@tanstack/react-query'
import { lazy } from 'react'
import { createRoot } from 'react-dom/client'
import { RouterProvider } from 'react-router'
import { router } from './router'

const root = document.getElementById('root')
if (!root) throw new Error('Root element not found')

const ReactQueryDevtools = lazy(() =>
	import('@tanstack/react-query-devtools').then((d) => ({
		default: d.ReactQueryDevtools,
	}))
)

createRoot(root).render(
	<QueryClientProvider client={queryClient}>
		<RouterProvider router={router} />
		{process.env.NODE_ENV === 'development' && <ReactQueryDevtools />}
	</QueryClientProvider>
)
