import type { ReactNode } from 'react'
import { labelIcon } from '../format'

const PRESETS = ['Home', 'Work'] as const

type Props = {
  value: string
  onChange: (label: string) => void
  customInput: ReactNode
  error?: string
}

export default function AddressLabelPicker({ value, onChange, customInput, error }: Props) {
  const isPreset = (PRESETS as readonly string[]).includes(value)

  return (
    <div className="label-picker">
      <span className="label-title">Save address as</span>
      <div className="chips" role="radiogroup" aria-label="Save address as">
        {PRESETS.map((p) => (
          <button
            key={p}
            type="button"
            role="radio"
            aria-checked={value === p}
            className={`chip${value === p ? ' selected' : ''}`}
            onClick={() => onChange(p)}
          >
            {labelIcon(p)} {p}
          </button>
        ))}
        <button
          type="button"
          role="radio"
          aria-checked={!isPreset}
          className={`chip${!isPreset ? ' selected' : ''}`}
          onClick={() => {
            if (isPreset) onChange('') // start an empty custom name; keep one already being typed
          }}
        >
          {labelIcon('')} Other
        </button>
        {!isPreset && customInput}
      </div>
      {error && <span className="hint bad">{error}</span>}
    </div>
  )
}
