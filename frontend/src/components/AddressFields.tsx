import type { FieldErrors, UseFormRegister } from 'react-hook-form'
import { INDIAN_STATES } from '../data/indianStates'
import { LIMITS, type CheckoutValues } from '../validation'
import AddressLabelPicker from './AddressLabelPicker'
import Field from './Field'

type Props = {
  register: UseFormRegister<CheckoutValues>
  errors: FieldErrors<CheckoutValues>
  label: string
  onLabelChange: (label: string) => void
}

export default function AddressFields({ register, errors, label, onLabelChange }: Props) {
  return (
    <>
      <div className="grid-2">
        <Field label="House / flat, street" error={errors.line1?.message}>
          <input
            autoComplete="address-line1"
            placeholder="Flat 4B, 12 MG Road"
            aria-invalid={!!errors.line1}
            {...register('line1')}
          />
        </Field>
        <Field
          label={
            <span>
              Area, landmark <span className="optional">(optional)</span>
            </span>
          }
          error={errors.line2?.message}
        >
          <input
            autoComplete="address-line2"
            placeholder="Near Metro Station"
            aria-invalid={!!errors.line2}
            {...register('line2')}
          />
        </Field>
      </div>

      <div className="grid-3">
        <Field label="City" error={errors.city?.message}>
          <input autoComplete="address-level2" placeholder="Bengaluru" aria-invalid={!!errors.city} {...register('city')} />
        </Field>
        <Field label="State" error={errors.state?.message}>
          <select autoComplete="address-level1" aria-invalid={!!errors.state} {...register('state')}>
            <option value="">Select state</option>
            {INDIAN_STATES.map((s) => (
              <option key={s} value={s}>
                {s}
              </option>
            ))}
          </select>
        </Field>
        <Field label="PIN code" error={errors.pincode?.message}>
          <input
            inputMode="numeric"
            autoComplete="postal-code"
            maxLength={LIMITS.pincode}
            placeholder="560001"
            aria-invalid={!!errors.pincode}
            {...register('pincode')}
          />
        </Field>
      </div>

      <AddressLabelPicker
        value={label}
        onChange={onLabelChange}
        error={errors.label?.message}
        customInput={
          <input
            autoFocus
            maxLength={LIMITS.addressLabel}
            placeholder="e.g. Mom's place"
            aria-label="Address name"
            aria-invalid={!!errors.label}
            {...register('label')}
          />
        }
      />
    </>
  )
}
