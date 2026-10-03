import { createContext, useContext } from 'react'

export type ThemePreference = 'light' | 'dark' | 'system'

export const THEME_STORAGE_KEY = 'gateforge-theme'

export type ThemeContextValue = {
  theme: ThemePreference
  setTheme: (theme: ThemePreference) => void
}

export const ThemeContext = createContext<ThemeContextValue | null>(null)

export function readStoredTheme(): ThemePreference {
  try {
    const stored = localStorage.getItem(THEME_STORAGE_KEY)
    if (stored === 'light' || stored === 'dark' || stored === 'system') {
      return stored
    }
  } catch {
    // Storage can be blocked. Fall through to the system preference.
  }
  return 'system'
}

export function applyTheme(theme: ThemePreference) {
  const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches
  const dark = theme === 'dark' || (theme === 'system' && prefersDark)
  document.documentElement.classList.toggle('dark', dark)
}

export function useTheme() {
  const value = useContext(ThemeContext)
  if (!value) {
    throw new Error('useTheme must be used within ThemeProvider')
  }
  return value
}

if (typeof window !== 'undefined') {
  applyTheme(readStoredTheme())
}
