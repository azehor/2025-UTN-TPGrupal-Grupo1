import { getApiBase } from './env'

// Token handling
export const TOKEN_KEY = 'auth_token'

export function getToken(): string | null {
  return localStorage.getItem(TOKEN_KEY)
}

export function setToken(token: string) {
  localStorage.setItem(TOKEN_KEY, token)
}

export function clearToken() {
  localStorage.removeItem(TOKEN_KEY)
}

export function isAuthenticated(): boolean {
  return !!getToken()
}

// Backward-compat (legacy staff flags) — keep until UI migrates fully
export function isStaffAuthenticated(): boolean {
  return sessionStorage.getItem('isStaff') === '1'
}

export function logoutStaff() {
  sessionStorage.removeItem('isStaff')
}

// Auth API helpers
type LoginResponse = {
  token: string
  expires_at?: string
  user?: { id: string; email: string; name?: string; role?: string }
} | { token: string; expires: string; user?: { id: string; email: string; name?: string; role?: string } }

export async function login(email: string, password: string): Promise<LoginResponse> {
  const res = await fetch(`${getApiBase()}/auth/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email, password })
  })
  if (!res.ok) {
    const text = await res.text()
    throw new Error(text || `Login failed (${res.status})`)
  }
  const data = await res.json()
  const token = (data?.token as string) || ''
  if (token) setToken(token)
  return data
}

export async function me(): Promise<any> {
  const token = getToken()
  if (!token) throw new Error('Not authenticated')
  const res = await fetch(`${getApiBase()}/auth/me`, {
    headers: { Authorization: `Bearer ${token}` }
  })
  if (!res.ok) throw new Error(`HTTP ${res.status}`)
  return res.json()
}

export function logout() {
  clearToken()
}
