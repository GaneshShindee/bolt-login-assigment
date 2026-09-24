import type { Address } from './api'
import { LIMITS } from './validation'

export function digitsOnly(value: string, max: number): string {
  return value.replace(/\D/g, '').slice(0, max)
}

/**
 * Normalizes typed or pasted input to a 10-digit mobile number: keeps digits only and drops a
 * pasted +91 country code or leading 0 ("+91 98765 43210" → "9876543210").
 */
export function toMobileDigits(value: string): string {
  const digits = value.replace(/\D/g, '')
  if (digits.length === 12 && digits.startsWith('91')) return digits.slice(2)
  if (digits.length === 11 && digits.startsWith('0')) return digits.slice(1)
  return digits.slice(0, LIMITS.mobile)
}

export function labelIcon(label: string): string {
  if (label === 'Home') return '🏠'
  if (label === 'Work') return '💼'
  return '📍'
}

/** "12 MG Road, Near Metro, Bengaluru, Karnataka 560001" */
export function formatAddress(a: Address): string {
  return [a.line1, a.line2, `${a.city}, ${a.state} ${a.pincode}`].filter(Boolean).join(', ')
}
