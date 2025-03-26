const EMAIL_REGEX = /^[a-zA-Z0–9._%+-]+@[a-zA-Z0–9.-]+\.[a-zA-Z]{2,}$/

export function isEmail(value: string) {
	return EMAIL_REGEX.test(value)
}
