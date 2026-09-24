import { useEffect, useState } from 'react'
import { api } from '../api'
import { isValidEmail, normalizeEmail } from '../validation'

export type Recognition = 'idle' | 'checking' | 'recognized' | 'new' | 'error'

const DEBOUNCE_MS = 400

export function useEmailRecognition(email: string, enabled: boolean): Recognition {
  const normalized = normalizeEmail(email)
  const shouldCheck = enabled && isValidEmail(email)
  const [result, setResult] = useState<{ email: string; status: Recognition } | null>(null)

  useEffect(() => {
    if (!shouldCheck) return
    const controller = new AbortController()
    const timer = setTimeout(() => {
      setResult({ email: normalized, status: 'checking' })
      api
        .recognize(normalized, controller.signal)
        .then(({ recognized }) => setResult({ email: normalized, status: recognized ? 'recognized' : 'new' }))
        .catch(() => {
          if (!controller.signal.aborted) setResult({ email: normalized, status: 'error' })
        })
    }, DEBOUNCE_MS)

    return () => {
      clearTimeout(timer)
      controller.abort()
    }
  }, [normalized, shouldCheck])

  if (!shouldCheck || result?.email !== normalized) return 'idle'
  return result.status
}
