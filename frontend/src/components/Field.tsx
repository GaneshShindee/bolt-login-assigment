import type { ReactNode } from 'react'

type Props = {
  label: ReactNode
  error?: string
  /** The input, plus anything shown under it (e.g. a live status hint). */
  children: ReactNode
}

export default function Field({ label, error, children }: Props) {
  return (
    <label>
      {label}
      {children}
      {error && <span className="hint bad">{error}</span>}
    </label>
  )
}
