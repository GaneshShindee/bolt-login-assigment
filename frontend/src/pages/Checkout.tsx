import { useCallback, useEffect, useState } from 'react'
import { useForm, useWatch } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { Link } from 'react-router-dom'
import { api, ApiError, type Order } from '../api'
import { useAuth } from '../auth/AuthContext'
import AddressFields from '../components/AddressFields'
import EmailStatus from '../components/EmailStatus'
import Field from '../components/Field'
import LoggedInBadge from '../components/LoggedInBadge'
import LoginModal from '../components/LoginModal'
import OrderReceipt from '../components/OrderReceipt'
import OrderSummary from '../components/OrderSummary'
import SavedDetailsPicker from '../components/SavedDetailsPicker'
import { toMobileDigits } from '../format'
import { EMPTY_DELIVERY, useDeliveryChoice } from '../hooks/useDeliveryChoice'
import { useEmailRecognition } from '../hooks/useEmailRecognition'
import { checkoutSchema, isValidEmail, normalizeEmail, type CheckoutValues } from '../validation'

const EMPTY_FORM: CheckoutValues = { email: '', ...EMPTY_DELIVERY }

export default function Checkout() {
  const { user, login, logout } = useAuth()

  const form = useForm<CheckoutValues>({
    resolver: zodResolver(checkoutSchema),
    mode: 'onSubmit',
    reValidateMode: 'onChange',
    defaultValues: EMPTY_FORM,
  })
  const {
    register,
    handleSubmit,
    control,
    setValue,
    reset,
    setError,
    formState: { errors, isSubmitting },
  } = form

  const email = useWatch({ control, name: 'email' })
  const label = useWatch({ control, name: 'label' })
  const normalized = normalizeEmail(email)
  const emailValid = isValidEmail(email)
  const recognition = useEmailRecognition(email, !user)
  const delivery = useDeliveryChoice(user, form)

  const [dismissed, setDismissed] = useState<ReadonlySet<string>>(new Set())
  const showLogin = !user && recognition === 'recognized' && !dismissed.has(normalized)

  const [order, setOrder] = useState<Order | null>(null)

  useEffect(() => {
    if (user) setValue('email', user.email, { shouldValidate: true })
  }, [user, setValue])

  const dismissLogin = useCallback(
    (forEmail: string) => {
      setDismissed((prev) => new Set(prev).add(forEmail))
    },
    [setDismissed],
  )

  const reopenLogin = () => {
    setDismissed((prev) => {
      const next = new Set(prev)
      next.delete(normalized)
      return next
    })
  }

  const handleLogout = () => {
    dismissLogin(normalized) // don't pop the modal straight back up for the email still in the field
    delivery.clear()
    logout()
  }

  const onSubmit = async ({ email, phone, ...address }: CheckoutValues) => {
    try {
      setOrder(await api.checkout(email, phone, address))
    } catch (err) {
      setError('root', { message: err instanceof ApiError ? err.message : 'Checkout failed.' })
    }
  }

  const startOver = () => {
    setOrder(null)
    reset({ ...EMPTY_FORM, email: user?.email ?? '' })
    delivery.reload()
  }

  const { onChange: onPhoneChange, ...phoneField } = register('phone')

  if (order) return <OrderReceipt order={order} user={user} onDone={startOver} />

  return (
    <>
      <div className="checkout-grid">
        <section className="card">
          <div className="checkout-head">
            <h1>Checkout</h1>
            {user && <LoggedInBadge user={user} onLogout={handleLogout} />}
          </div>
          {!user && (
            <p className="muted checkout-intro">
              Returning customer? Type your email to log in. New here? <Link to="/register">Register</Link>
            </p>
          )}

          <form id="checkout-form" onSubmit={handleSubmit(onSubmit, delivery.edit)} noValidate>
            <fieldset>
              <legend>Contact information</legend>
              <div className="grid-2">
                <Field label="Email address">
                  <input
                    type="email"
                    autoComplete="email"
                    placeholder="you@example.com"
                    readOnly={!!user}
                    aria-invalid={!!errors.email}
                    className={emailValid ? 'valid' : undefined}
                    {...register('email')}
                  />
                  <EmailStatus
                    email={email}
                    valid={emailValid}
                    recognition={recognition}
                    loggedIn={!!user}
                    error={errors.email?.message}
                    onLogin={reopenLogin}
                  />
                </Field>
                {delivery.showFields && (
                  <Field label="Mobile number" error={errors.phone?.message}>
                    <input
                      type="tel"
                      inputMode="numeric"
                      autoComplete="tel-national"
                      placeholder="9876543210"
                      aria-invalid={!!errors.phone}
                      {...phoneField}
                      onChange={(e) => {
                        e.target.value = toMobileDigits(e.target.value)
                        return onPhoneChange(e)
                      }}
                    />
                  </Field>
                )}
              </div>
            </fieldset>

            <fieldset>
              <legend>Delivery details</legend>
              {delivery.savedList.length > 0 && (
                <SavedDetailsPicker
                  saved={delivery.savedList}
                  selected={delivery.choice}
                  editing={delivery.editing}
                  compact={delivery.showFields}
                  newAddress={delivery.newAddress}
                  onSelect={delivery.choose}
                  onEdit={delivery.edit}
                />
              )}
              {delivery.showFields && (
                <AddressFields
                  register={register}
                  errors={errors}
                  label={label}
                  onLabelChange={(l) => setValue('label', l, { shouldValidate: !!errors.label })}
                />
              )}
            </fieldset>
          </form>
        </section>

        <OrderSummary>
          {errors.root && (
            <p className="error" role="alert">
              {errors.root.message}
            </p>
          )}
          <button className="btn full large" type="submit" form="checkout-form" disabled={isSubmitting}>
            {isSubmitting ? 'Placing order…' : 'Place order'}
          </button>
        </OrderSummary>
      </div>

      {showLogin && <LoginModal email={normalized} onSuccess={login} onSkip={dismissLogin} />}
    </>
  )
}
