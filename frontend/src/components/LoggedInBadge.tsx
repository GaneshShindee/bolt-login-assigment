import type { User } from '../api'

type Props = {
  user: User
  onLogout: () => void
}

export default function LoggedInBadge({ user, onLogout }: Props) {
  return (
    <div className="welcome">
      <span className="avatar" aria-hidden>
        {user.first_name[0]}
        {user.last_name[0]}
      </span>
      <span>
        Logged in as{' '}
        <strong>
          {user.first_name} {user.last_name}
        </strong>
      </span>
      <button type="button" className="btn link small" onClick={onLogout}>
        Log out
      </button>
    </div>
  )
}
