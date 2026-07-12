import { type ClassValue, clsx } from 'clsx'
import { twMerge } from 'tailwind-merge'

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export function getApiBaseUrl(): string {
  const base = import.meta.env.VITE_API_BASE_URL ?? ''
  return base.replace(/\/$/, '')
}

export function apiUrl(path: string): string {
  const base = getApiBaseUrl()
  return `${base}${path.startsWith('/') ? path : `/${path}`}`
}
