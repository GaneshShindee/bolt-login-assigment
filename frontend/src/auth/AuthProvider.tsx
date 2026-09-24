import { useCallback, useEffect, useMemo, useState, type ReactNode } from 'react'
import { api, ApiError, getToken, setToken, type User } from '../api'
import { AuthContext } from './AuthContext'

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null)

  useEffect(() => {
    if (!getToken()) return
    api
      .me()
      .then((res) => setUser(res.user))
      .catch((err) => {
        // Only drop the token if the server rejected it; a network blip shouldn't log the user out.
        if (err instanceof ApiError && err.status === 401) setToken(null)
      })
  }, [])

  const login = useCallback((u: User, token: string) => {
    setToken(token)
    setUser(u)
  }, [])

  const logout = useCallback(() => {
    setToken(null)
    setUser(null)
  }, [])

  const value = useMemo(() => ({ user, login, logout }), [user, login, logout])
  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
}
