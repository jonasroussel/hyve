import { Button } from "@/components/ui/button";

export function OnboardingPage() {
	return (
		<div className="flex flex-col items-center justify-center p-10 mx-2 mt-10 lg:p-20">
			<h1 className="text-3xl font-bold lg:text-5xl">Welcome to Hyve</h1>
			<h3 className="py-6 text-center text-muted-foreground text-sm lg:text-xl">Let me help you set up the basics</h3>
			<Button className="py-8 px-16 mt-2">Get Started</Button>
		</div>
	)
}
