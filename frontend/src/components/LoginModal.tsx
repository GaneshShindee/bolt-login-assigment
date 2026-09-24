import { useEffect } from 'react'
import { useForm, useWatch } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { api, ApiError, type User } from '../api'
import { digitsOnly } from '../format'
import { LIMITS, loginCodeSchema, type LoginCodeValues } from '../validation'

type Props = {
  email: string
  onSuccess: (user: User, token: string) => void
  onSkip: (email: string) => void
}

export default function LoginModal({ email, onSuccess, onSkip }: Props) {
  const {
    register,
    handleSubmit,
    setError,
    clearErrors,
    setFocus,
    resetField,
    control,
    formState: { errors, isSubmitting },
  } = useForm<LoginCodeValues>({
    resolver: zodResolver(loginCodeSchema),
    reValidateMode: 'onSubmit',
    defaultValues: { code: '' },
  })
  const code = useWatch({ control, name: 'code' })

  useEffect(() => {
    setFocus('code')
    const onKey = (e: KeyboardEvent) => e.key === 'Escape' && onSkip(email)
    window.addEventListener('keydown', onKey)
    return () => window.removeEventListener('keydown', onKey)
  }, [email, onSkip, setFocus])

  const onSubmit = async ({ code }: LoginCodeValues) => {
    try {
      const res = await api.login(email, code)
      onSuccess(res.user, res.token)
    } catch (err) {
      resetField('code')
      setError('code', { message: err instanceof ApiError ? err.message : 'Login failed.' }, { shouldFocus: true })
    }
  }

  const { onChange: onCodeChange, ...codeField } = register('code')

  return (
    <div className="modal-backdrop" onMouseDown={(e) => e.target === e.currentTarget && onSkip(email)}>
      <div className="modal" role="dialog" aria-modal="true" aria-labelledby="login-title">
        <h2 id="login-title">Welcome back!</h2>
        <p className="muted">
          We recognized <strong>{email}</strong>. Enter the 6-digit code you got when you registered.
        </p>
        <form onSubmit={handleSubmit(onSubmit)} noValidate>
          <input
            className="code-input"
            inputMode="numeric"
            autoComplete="one-time-code"
            maxLength={LIMITS.code}
            placeholder="••••••"
            aria-label="6-digit login code"
            aria-invalid={!!errors.code}
            {...codeField}
            onChange={(e) => {
              e.target.value = digitsOnly(e.target.value, LIMITS.code)
              clearErrors('code')
              return onCodeChange(e)
            }}
          />
          {errors.code && (
            <p className="error" role="alert">
              {errors.code.message}
            </p>
          )}
          <button className="btn full" type="submit" disabled={isSubmitting || code.length !== LIMITS.code}>
            {isSubmitting ? 'Verifying…' : 'Log in'}
          </button>
          <button className="btn link full" type="button" onClick={() => onSkip(email)}>
            Skip and continue as guest
          </button>
        </form>
      </div>
    </div>
  )
}
