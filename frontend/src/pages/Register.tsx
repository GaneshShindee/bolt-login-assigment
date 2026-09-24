import { useState } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { Link } from 'react-router-dom'
import { api, ApiError } from '../api'
import { useAuth } from '../auth/AuthContext'
import Field from '../components/Field'
import { registerSchema, type RegisterValues } from '../validation'

export default function Register() {
  const { user } = useAuth()
  const {
    register,
    handleSubmit,
    setError,
    formState: { errors, isSubmitting },
  } = useForm<RegisterValues>({
    resolver: zodResolver(registerSchema),
    mode: 'onSubmit',
    reValidateMode: 'onChange',
    defaultValues: { email: '', firstName: '', lastName: '' },
  })
  const [result, setResult] = useState<{ code: string; firstName: string } | null>(null)
  const [copied, setCopied] = useState(false)

  const onSubmit = async (values: RegisterValues) => {
    try {
      const res = await api.register(values.email, values.firstName, values.lastName)
      setResult({ code: res.code, firstName: res.user.first_name })
    } catch (err) {
      setError('root', { message: err instanceof ApiError ? err.message : 'Registration failed.' })
    }
  }

  const copy = async () => {
    if (!result) return
    try {
      await navigator.clipboard.writeText(result.code)
      setCopied(true)
      setTimeout(() => setCopied(false), 2000)
    } catch {
      /* clipboard blocked; the code is still visible */
    }
  }

  if (user) {
    return (
      <section className="card narrow center">
        <h1>You're already logged in</h1>
        <p className="muted">
          You're logged in as{' '}
          <strong>
            {user.first_name} {user.last_name}
          </strong>
          . Log out on the checkout page to register a different account.
        </p>
        <Link to="/" className="btn">
          Go to checkout →
        </Link>
      </section>
    )
  }

  if (result) {
    return (
      <section className="card narrow center">
        <h1>You're registered, {result.firstName}!</h1>
        <p className="muted">This is your login code. Save it now, because it won't be shown again.</p>
        <div className="code-display" aria-label="Your login code">
          {result.code}
        </div>
        <div className="row">
          <button type="button" className="btn secondary" onClick={copy}>
            {copied ? 'Copied ✓' : 'Copy code'}
          </button>
          <Link to="/" className="btn">
            Go to checkout →
          </Link>
        </div>
      </section>
    )
  }

  return (
    <section className="card narrow">
      <h1>Create an account</h1>
      <p className="muted">Register once, then use your 6-digit code to log in at checkout.</p>
      <form onSubmit={handleSubmit(onSubmit)} noValidate>
        <Field label="Email address" error={errors.email?.message}>
          <input type="email" autoComplete="email" aria-invalid={!!errors.email} {...register('email')} />
        </Field>
        <div className="grid-2">
          <Field label="First name" error={errors.firstName?.message}>
            <input autoComplete="given-name" aria-invalid={!!errors.firstName} {...register('firstName')} />
          </Field>
          <Field label="Last name" error={errors.lastName?.message}>
            <input autoComplete="family-name" aria-invalid={!!errors.lastName} {...register('lastName')} />
          </Field>
        </div>
        {errors.root && (
          <p className="error" role="alert">
            {errors.root.message}
          </p>
        )}
        <button className="btn full" type="submit" disabled={isSubmitting}>
          {isSubmitting ? 'Registering…' : 'Register'}
        </button>
      </form>
    </section>
  )
}
