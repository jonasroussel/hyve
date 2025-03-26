import { RootPage } from '@/pages/root'
import { createBrowserRouter } from 'react-router'

export const router = createBrowserRouter([
	{
		path: '/',
		element: <RootPage />,
		children: [
			{
				path: '',
				lazy: async () => ({ Component: (await import('@/pages/dashboard')).DashboardPage }),
			},
			{
				path: 'register',
				lazy: async () => ({ Component: (await import('@/pages/register')).RegisterPage }),
			},
			{
				path: 'login',
				lazy: async () => ({ Component: (await import('@/pages/login')).LoginPage }),
			},
			{
				path: 'onboarding',
				lazy: async () => ({ Component: (await import('@/pages/onboarding')).OnboardingPage }),
			},
		],
	},
])
