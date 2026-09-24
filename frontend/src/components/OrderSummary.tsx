import type { ReactNode } from 'react'

const ITEM = { name: 'Bolt Classic Tee', variant: 'Black · M', price: 999 }
const SHIPPING = 0

const inr = (n: number) => `₹${n.toLocaleString('en-IN')}`

export default function OrderSummary({ children }: { children?: ReactNode }) {
  return (
    <aside className="card summary">
      <h2>Order summary</h2>
      <div className="summary-item">
        <div className="summary-thumb" aria-hidden>
          👕
        </div>
        <div className="summary-item-text">
          <strong>{ITEM.name}</strong>
          <span className="muted small">{ITEM.variant} · Qty 1</span>
        </div>
        <span>{inr(ITEM.price)}</span>
      </div>
      <dl className="summary-totals">
        <dt>Subtotal</dt>
        <dd>{inr(ITEM.price)}</dd>
        <dt>Shipping</dt>
        <dd>{SHIPPING === 0 ? 'Free' : inr(SHIPPING)}</dd>
        <dt className="total">Total</dt>
        <dd className="total">{inr(ITEM.price + SHIPPING)}</dd>
      </dl>
      {children}
      <p className="demo-note">Demo store: placing an order only saves your details. No payment is taken.</p>
    </aside>
  )
}
