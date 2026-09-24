import { z } from 'zod'

// Field limits. Keep in sync with the API (backend/internal/biz/validation.go) and db/migrations.
export const LIMITS = {
  email: 254,
  name: 100,
  addressLabel: 30,
  addressLine: 200,
  cityOrState: 100,
  pincode: 6,
  mobile: 10,
  code: 6,
} as const

// Mirrors the API's rule: local@domain.tld, TLD of 2+ letters, no leading/trailing hyphens in labels.
const EMAIL_RE =
  /^[A-Za-z0-9._%+-]+@[A-Za-z0-9](?:[A-Za-z0-9-]*[A-Za-z0-9])?(?:\.[A-Za-z0-9](?:[A-Za-z0-9-]*[A-Za-z0-9])?)*\.[A-Za-z]{2,}$/
const MOBILE_RE = /^[6-9][0-9]{9}$/ // Indian mobile numbers: 10 digits, starting 6–9
const PINCODE_RE = /^[1-9][0-9]{5}$/ // Indian PIN codes are 6 digits and never start with 0

const TOO_LONG = 'This is too long.'

const requiredText = (message: string, max: number) => z.string().trim().min(1, message).max(max, TOO_LONG)

const emailSchema = z
  .string()
  .trim()
  .min(1, 'Email address is required.')
  .max(LIMITS.email, 'Email address is too long.')
  .regex(EMAIL_RE, 'Please enter a complete email address.')

export const registerSchema = z.object({
  email: emailSchema,
  firstName: requiredText('First name is required.', LIMITS.name),
  lastName: requiredText('Last name is required.', LIMITS.name),
})
export type RegisterValues = z.infer<typeof registerSchema>

export const checkoutSchema = z.object({
  email: emailSchema,
  phone: z.string().trim().regex(MOBILE_RE, 'Please enter a valid 10-digit mobile number.'),
  label: requiredText('Please give this address a name, like Home or Work.', LIMITS.addressLabel),
  line1: requiredText('Please enter your house / flat and street.', LIMITS.addressLine),
  line2: z.string().trim().max(LIMITS.addressLine, TOO_LONG),
  city: requiredText('Please enter your city.', LIMITS.cityOrState),
  state: requiredText('Please choose your state.', LIMITS.cityOrState),
  pincode: z.string().trim().regex(PINCODE_RE, 'Please enter a valid 6-digit PIN code.'),
})
export type CheckoutValues = z.infer<typeof checkoutSchema>

export const loginCodeSchema = z.object({
  code: z.string().regex(/^\d{6}$/, 'Please enter all 6 digits.'),
})
export type LoginCodeValues = z.infer<typeof loginCodeSchema>

export function isValidEmail(value: string): boolean {
  return emailSchema.safeParse(value).success
}

export function normalizeEmail(value: string): string {
  return value.trim().toLowerCase()
}
