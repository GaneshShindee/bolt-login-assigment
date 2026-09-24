import type { Order, User } from '../api'
import { labelIcon } from '../format'

type Props = {
  order: Order
  user: User | null
  onDone: () => void
}

export default function OrderReceipt({ order, user, onDone }: Props) {
  const placedAt = new Date(order.created_at).toLocaleString()
  const linkedToUser = order.user_id !== null

  return (
    <section className="card narrow center">
      <div className="success-icon" aria-hidden>
        ✓
      </div>
      <h1>Order placed{user ? `, ${user.first_name}` : ''}!</h1>
      <p className="muted">This is the record saved in the database. No payment was taken.</p>

      <dl className="receipt">
        <dt>Order</dt>
        <dd>#{order.id}</dd>
        <dt>Placed</dt>
        <dd>{placedAt}</dd>
        <dt>Account</dt>
        <dd>
          {linkedToUser ? (
            <span className="badge good">
              Linked to {user ? `${user.first_name} ${user.last_name}` : `user #${order.user_id}`}
            </span>
          ) : (
            <span className="badge">Guest checkout</span>
          )}
        </dd>
        <dt>Email</dt>
        <dd>{order.email}</dd>
        <dt>Mobile</dt>
        <dd>{order.phone}</dd>
        <dt>Ship to</dt>
        <dd>
          <strong>
            {labelIcon(order.address.label)} {order.address.label}
          </strong>
          <br />
          {order.address.line1}
          {order.address.line2 && <br />}
          {order.address.line2}
          <br />
          {order.address.city}, {order.address.state} {order.address.pincode}
        </dd>
      </dl>

      <button className="btn" onClick={onDone}>
        Place another order
      </button>
    </section>
  )
}
