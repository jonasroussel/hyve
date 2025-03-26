import axios, { AxiosError } from 'axios'

export type APIError = AxiosError<{ code: string } & { [key: string]: any }>['response']

export const API = axios.create({
	baseURL: '/api',
	timeout: 30_000,
})
