import type { FC, FormEvent } from 'react'
import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import './Login.css'
import { login as apiLogin } from '../lib/auth'

export const Login: FC = () => {
  const [user, setUser] = useState('')
  const [pass, setPass] = useState('')
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)
  const navigate = useNavigate()

  async function submit(e: FormEvent) {
    e.preventDefault()
    setError(null)
    setLoading(true)
    try {
      await apiLogin(user.trim(), pass)
      // Compat: mantener flag legacy para secciones que aún lo usan
      sessionStorage.setItem('isStaff', '1')
      navigate('/panel')
    } catch (err: any) {
      const msg = typeof err?.message === 'string' ? err.message : 'Credenciales inválidas'
      setError(msg)
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="login-root">
      <form className="login-box" onSubmit={submit}>
        <h2>Acceso al Panel</h2>
        {error && <div className="login-error">{error}</div>}
        <label>Usuario<input value={user} onChange={(e) => setUser(e.target.value)} /></label>
        <label>Contraseña<input type="password" value={pass} onChange={(e) => setPass(e.target.value)} /></label>
        <div style={{ display: 'flex', gap: 8 }}>
          <button className="btn-login" type="submit" disabled={loading}>
            {loading ? 'Ingresando...' : 'Entrar'}
          </button>
        </div>
      </form>
    </div>
  )
}

export default Login
