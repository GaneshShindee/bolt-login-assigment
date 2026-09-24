import { useEffect, useState } from 'react'
import { api } from '../api'
import { isValidEmail, normalizeEmail } from '../validation'

export type Recognition = 'idle' | 'checking' | 'recognized' | 'new' | 'error'

const DEBOUNCE_MS = 400

/**
 * Checks in the background whether `email` belongs to a registered user, once it is well-formed.
 * Requests are debounced while the user types, and stale ones are aborted when the email changes.
 */
export function useEmailRecognition(email: string, enabled: boolean): Recognition {
  const normalized = normalizeEmail(email)
  const shouldCheck = enabled && isValidEmail(email)
  // The status is tagged with the email it belongs to, so a result for an old email is never shown.
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
