import type { Address } from './api'
import { LIMITS } from './validation'

export function digitsOnly(value: string, max: number): string {
  return value.replace(/\D/g, '').slice(0, max)
}

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

export function formatAddress(a: Address): string {
  return [a.line1, a.line2, `${a.city}, ${a.state} ${a.pincode}`].filter(Boolean).join(', ')
}
