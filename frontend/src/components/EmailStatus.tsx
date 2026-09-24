import type { Recognition } from '../hooks/useEmailRecognition'

type Props = {
  email: string
  valid: boolean
  recognition: Recognition
  loggedIn: boolean
  error?: string
  onLogin: () => void
}

export default function EmailStatus({ email, valid, recognition, loggedIn, error, onLogin }: Props) {
  if (!valid) {
    if (error) return <span className="hint bad">{error}</span>
    return email.includes('@') ? <span className="hint bad">Keep typing a complete email…</span> : null
  }
  if (loggedIn) return <span className="hint good">✓ Valid email</span>

  switch (recognition) {
    case 'checking':
      return (
        <span className="hint">
          <span className="spinner" aria-hidden /> Checking for an account…
        </span>
      )
    case 'recognized':
      return (
        <span className="hint good">
          ✓ Account found ·{' '}
          <button type="button" className="btn link small inline" onClick={onLogin}>
            Log in with your code
          </button>
        </span>
      )
    case 'new':
      return <span className="hint">✓ New customer · checking out as a guest</span>
    case 'error':
      return <span className="hint">✓ Valid email · couldn't check for an account</span>
    case 'idle':
      return <span className="hint good">✓ Valid email</span>
  }
}
