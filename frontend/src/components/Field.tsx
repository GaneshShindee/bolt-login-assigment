import type { ReactNode } from 'react'

type Props = {
  label: ReactNode
  error?: string
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
