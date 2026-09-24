import { describe, expect, it } from 'vitest'
import { checkoutSchema, isValidEmail, normalizeEmail, registerSchema } from './validation'

const validCheckout = {
  email: 'asha@example.com',
  phone: '9876543210',
  label: 'Home',
  line1: 'Flat 4B, 12 MG Road',
  line2: '',
  city: 'Bengaluru',
  state: 'Karnataka',
  pincode: '560001',
}

describe('isValidEmail (real-time check before recognition)', () => {
  it.each(['a@b.co', 'test.user@example.com', 'first+tag@sub.domain.org'])('accepts %s', (e) => {
    expect(isValidEmail(e)).toBe(true)
  })
  it.each(['', 'plain', 'a@', 'a@b', 'a@b.c', 'a b@c.com', 'a@-b.com'])('rejects incomplete %s', (e) => {
    expect(isValidEmail(e)).toBe(false)
  })
})

describe('normalizeEmail', () => {
  it('trims and lowercases', () => {
    expect(normalizeEmail('  Asha@Example.COM ')).toBe('asha@example.com')
  })
})

describe('registerSchema', () => {
  it('requires names', () => {
    const r = registerSchema.safeParse({ email: 'a@b.com', firstName: ' ', lastName: 'Rao' })
    expect(r.success).toBe(false)
  })
})

describe('checkoutSchema', () => {
  it('accepts a complete checkout, with line 2 optional', () => {
    expect(checkoutSchema.safeParse(validCheckout).success).toBe(true)
  })

  it.each([
    ['pincode', '56000'], // 5 digits
    ['pincode', '060001'], // leading 0
    ['pincode', '56A001'],
    ['phone', 'abc'],
    ['phone', '987654321'], // 9 digits
    ['phone', '98765432101'], // 11 digits
    ['phone', '5876543210'], // must start with 6–9
    ['phone', '+919876543210'],
    ['label', '  '],
    ['line1', ''],
    ['city', ''],
    ['state', ''],
  ])('rejects %s = %j', (field, value) => {
    expect(checkoutSchema.safeParse({ ...validCheckout, [field]: value }).success).toBe(false)
  })

  it('trims values', () => {
    const r = checkoutSchema.parse({ ...validCheckout, city: '  Pune  ' })
    expect(r.city).toBe('Pune')
  })
})
