import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { useRootRegister } from '@/hooks/api/users'
import { isEmail } from '@/lib/checkers'
import { useForm } from '@tanstack/react-form'
import { useQueryClient } from '@tanstack/react-query'

export function RegisterPage() {
	const rootRegister = useRootRegister()
	const queryClient = useQueryClient()

	const form = useForm({
		defaultValues: {
			username: '',
			email: '',
			password: '',
			confirmPassword: '',
		},
		onSubmit: async ({ value }) => {
			await rootRegister.mutateAsync({
				username: value.username,
				email: value.email,
				password: value.password,
			})

			await queryClient.invalidateQueries({ queryKey: ['connected-user'] })
		},
	})

	return (
		<div className="flex flex-col items-center justify-center min-h-screen">
			<h1 className="text-5xl font-bold mb-6">Hyve</h1>
			<div className="w-full max-w-md p-6">
				<h3 className="text-2xl font-semibold mb-1">Create an account</h3>
				<p className="text-xs text-primary mb-6">This user will be the root user. Full admin access.</p>
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
								name="username"
								validators={{
									onBlur: ({ value }) => {
										if (!value) return 'Username is required'
										if (value.length < 3) return 'Username must be at least 3 characters'
									},
								}}
							>
								{(field) => (
									<>
										<Label htmlFor={field.name}>
											Username <span className="text-red-500">*</span>
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
								name="email"
								validators={{
									onBlur: ({ value }) => {
										if (!value) return
										if (!isEmail(value)) return 'Email must be a valid email'
									},
								}}
							>
								{(field) => (
									<>
										<Label htmlFor={field.name}>Email</Label>
										<Input
											id={field.name}
											type="email"
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
						<div className="grid gap-2">
							<form.Field
								name="confirmPassword"
								validators={{
									onBlur: ({ value }) => {
										if (!value) return 'Password is required'
										if (value !== form.state.values.password) return 'Passwords do not match'
									},
								}}
							>
								{(field) => (
									<>
										<Label htmlFor={field.name}>
											Password again <span className="text-red-500">*</span>
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
							Register
						</Button>
					</div>
				</form>
			</div>
		</div>
	)
}
