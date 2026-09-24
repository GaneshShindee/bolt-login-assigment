import { NavLink, Navigate, Route, Routes } from 'react-router-dom'
import { useAuth } from './auth/AuthContext'
import Register from './pages/Register'
import Checkout from './pages/Checkout'

export default function App() {
  const { user } = useAuth()

  return (
    <>
      <header className="topbar">
        <span className="brand">⚡ BoltShop</span>
        <nav>
          <NavLink to="/" end>
            Checkout
          </NavLink>
          {!user && <NavLink to="/register">Register</NavLink>}
        </nav>
      </header>
      <main className="container">
        <Routes>
          <Route path="/" element={<Checkout />} />
          <Route path="/register" element={<Register />} />
          <Route path="*" element={<Navigate to="/" replace />} />
        </Routes>
      </main>
    </>
  )
}
