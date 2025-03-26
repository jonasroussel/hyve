import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { useLogin } from '@/hooks/api/auth'
import { useForm } from '@tanstack/react-form'

export function LoginPage() {
	const login = useLogin()

	const form = useForm({
		defaultValues: {
			username_email: '',
			password: '',
		},
		onSubmit: async ({ value }) => {
			await login.mutateAsync({
				username_email: value.username_email,
				password: value.password,
			})
		},
	})

	return (
		<div className="flex flex-col items-center justify-center min-h-screen">
			<h1 className="text-5xl font-bold mb-6">Hyve</h1>
			<div className="w-full max-w-md p-6">
				<form
					onSubmit={(e) => {
						e.preventDefault()
						e.stopPropagation()
						form.handleSubmit()
					}}
				>
					<div className="flex flex-col gap-4">
						<div className="grid gap-2">
							<form.Field
								name="username_email"
								validators={{
									onBlur: ({ value }) => {
										if (!value) return 'Username or email is required'
									},
								}}
							>
								{(field) => (
									<>
										<Label htmlFor={field.name}>
											Username / email <span className="text-red-500">*</span>
										</Label>
										<Input
											id={field.name}
											type="text"
											name={field.name}
											value={field.state.value}
											onBlur={field.handleBlur}
											onChange={(e) => field.handleChange(e.target.value)}
											error={field.state.meta.errors?.join(', ')}
										/>
									</>
								)}
							</form.Field>
						</div>
						<div className="grid gap-2">
							<form.Field
								name="password"
								validators={{
									onBlur: ({ value }) => {
										if (!value) return 'Password is required'
									},
								}}
							>
								{(field) => (
									<>
										<Label htmlFor={field.name}>
											Password <span className="text-red-500">*</span>
										</Label>
										<Input
											id={field.name}
											type="password"
											name={field.name}
											value={field.state.value}
											onBlur={field.handleBlur}
											onChange={(e) => field.handleChange(e.target.value)}
											error={field.state.meta.errors?.join(', ')}
										/>
									</>
								)}
							</form.Field>
						</div>
						<Button type="submit" className="w-full mt-6">
							Login
						</Button>
					</div>
				</form>
			</div>
		</div>
	)
}
