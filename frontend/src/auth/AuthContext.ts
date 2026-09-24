import { createContext, useContext } from 'react'
import type { User } from '../api'

type AuthState = {
  user: User | null
  login: (user: User, token: string) => void
  logout: () => void
}

export const AuthContext = createContext<AuthState | null>(null)

export function useAuth(): AuthState {
  const ctx = useContext(AuthContext)
  if (!ctx) throw new Error('useAuth must be used inside AuthProvider')
  return ctx
}
