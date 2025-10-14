import type { FC, FormEvent } from 'react'
import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import './Login.css'

// Temporary credentials - in a real app use a secure backend
const STAFF_USER = 'admin'
const STAFF_PASS = 'admin123'

export const Login: FC = () => {
  const [user, setUser] = useState('')
  const [pass, setPass] = useState('')
  const [error, setError] = useState<string | null>(null)
  const navigate = useNavigate()

  function submit(e: FormEvent) {
    e.preventDefault()
    if (user === STAFF_USER && pass === STAFF_PASS) {
      sessionStorage.setItem('isStaff', '1')
      navigate('/panel')
      return
    }
    setError('Credenciales inválidas')
  }

  return (
    <div className="login-root">
      <form className="login-box" onSubmit={submit}>
        <h2>Acceso al Panel</h2>
        {error && <div className="login-error">{error}</div>}
        <label>Usuario<input value={user} onChange={(e) => setUser(e.target.value)} /></label>
        <label>Contraseña<input type="password" value={pass} onChange={(e) => setPass(e.target.value)} /></label>
        <div style={{ display: 'flex', gap: 8 }}>
          <button className="btn-login" type="submit">Entrar</button>
        </div>
      </form>
    </div>
  )
}

export default Login
