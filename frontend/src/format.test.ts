import { describe, expect, it } from 'vitest'
import { formatAddress, labelIcon, toMobileDigits } from './format'

describe('formatAddress', () => {
  const base = { label: 'Home', line1: '12 MG Road', line2: '', city: 'Bengaluru', state: 'Karnataka', pincode: '560001' }

  it('skips an empty line 2', () => {
    expect(formatAddress(base)).toBe('12 MG Road, Bengaluru, Karnataka 560001')
  })
  it('includes line 2 when given', () => {
    expect(formatAddress({ ...base, line2: 'Near Metro' })).toBe('12 MG Road, Near Metro, Bengaluru, Karnataka 560001')
  })
})

describe('toMobileDigits', () => {
  it.each([
    ['9876543210', '9876543210'],
    ['98765 43210', '9876543210'], // spaces removed
    ['+91 98765 43210', '9876543210'], // pasted with country code
    ['09876543210', '9876543210'], // leading 0
    ['98765432109999', '9876543210'], // never more than 10
    ['98ab76', '9876'], // letters dropped while typing
  ])('%j → %j', (input, expected) => {
    expect(toMobileDigits(input)).toBe(expected)
  })
})

describe('labelIcon', () => {
  it('has icons for Home, Work and custom names', () => {
    expect(labelIcon('Home')).toBe('🏠')
    expect(labelIcon('Work')).toBe('💼')
    expect(labelIcon("Mom's place")).toBe('📍')
  })
})
