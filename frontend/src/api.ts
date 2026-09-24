const API_URL = (import.meta.env.VITE_API_URL ?? 'http://localhost:8080').replace(/\/$/, '')

export type User = {
  id: number
  email: string
  first_name: string
  last_name: string
}

export type Address = {
  label: string
  line1: string
  line2: string // optional; '' when not given
  city: string
  state: string
  pincode: string
}

export type Order = {
  id: number
  user_id: number | null // null for guest checkouts
  email: string
  phone: string
  address: Address
  created_at: string
}

export type SavedDetails = {
  phone: string
  address: Address
  last_used_at: string
}

export class ApiError extends Error {
  status: number
  constructor(status: number, message: string) {
    super(message)
    this.status = status
  }
}

const TOKEN_KEY = 'bolt_token'

export function getToken(): string | null {
  try {
    return sessionStorage.getItem(TOKEN_KEY)
  } catch {
    return null
  }
}

export function setToken(token: string | null) {
  try {
    if (token) sessionStorage.setItem(TOKEN_KEY, token)
    else sessionStorage.removeItem(TOKEN_KEY)
  } catch {
    /* storage unavailable: session lasts for this page load only */
  }
}

type RequestOptions = {
  method?: 'GET' | 'POST'
  json?: unknown
  signal?: AbortSignal
}

async function request<T>(path: string, { method = 'GET', json, signal }: RequestOptions = {}): Promise<T> {
  const headers: Record<string, string> = {}
  if (json !== undefined) headers['Content-Type'] = 'application/json'
  const token = getToken()
  if (token) headers['Authorization'] = `Bearer ${token}`

  let res: Response
  try {
    res = await fetch(`${API_URL}${path}`, {
      method,
      headers,
      body: json !== undefined ? JSON.stringify(json) : undefined,
      signal,
    })
  } catch (err) {
    if (err instanceof DOMException && err.name === 'AbortError') throw err
    throw new ApiError(0, 'Could not reach the server. Please check your connection.')
  }

  const data = await res.json().catch(() => ({}))
  if (!res.ok) throw new ApiError(res.status, data.error ?? 'Something went wrong. Please try again.')
  return data as T
}

export const api = {
  register: (email: string, first_name: string, last_name: string) =>
    request<{ user: User; code: string }>('/api/register', {
      method: 'POST',
      json: { email, first_name, last_name },
    }),

  recognize: (email: string, signal?: AbortSignal) =>
    request<{ recognized: boolean }>('/api/recognize', { method: 'POST', json: { email }, signal }),

  login: (email: string, code: string) =>
    request<{ user: User; token: string }>('/api/login', { method: 'POST', json: { email, code } }),

  me: () => request<{ user: User }>('/api/me'),

  checkout: (email: string, phone: string, address: Address) =>
    request<Order>('/api/checkout', { method: 'POST', json: { email, phone, address } }),

  savedDetails: () => request<{ saved: SavedDetails[] }>('/api/me/saved-details'),
}
