import type { FC } from 'react'
import './App.css'
import Home from "./pages/Home"
import PanelDeAdmin from './pages/PanelDeAdmin'
import { Link, Route, Routes, useLocation, Navigate } from 'react-router-dom'
import Login from './pages/Login'
import { isStaffAuthenticated } from './lib/auth'
import Footer from './components/Footer'
import { BusquedaSoftwarePage } from './pages/BusquedaSoftwarePage'
import { BusquedaCarreraPage } from './pages/BusquedaCarreraPage'
import { BusquedaJuegoPage } from './pages/BusquedaJuegoPage'
import { RecomendacionPage } from './pages/RecomendacionPage'


const App: FC = () => {
  const location = useLocation();
  const isHome = location.pathname === "/";

  return (
    <div className="app-surface">
      {!isHome && (

        <nav
          style={{
            display: "flex",
            gap: "1rem",
            padding: "1rem",
            borderBottom: "1px solid #ccc",
            marginBottom: "1rem",
            justifyContent: "center"
          }}>
          <Link to="/">Inicio</Link>
          <Link to="/busqueda-software">Buscar por Software</Link>
          <Link to="/busqueda-carrera">Buscar por Carrera</Link>
          <Link to="/busqueda-juego">Buscar por Juego</Link>
        </nav>
      )}

      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/busqueda-software" element={<BusquedaSoftwarePage />} />
        <Route path="/busqueda-carrera" element={<BusquedaCarreraPage />} />
        <Route path="/busqueda-juego" element={<BusquedaJuegoPage />} />
        <Route path="/recomendacion" element={<RecomendacionPage />} />
        <Route path="/panel" element={isStaffAuthenticated() ? <PanelDeAdmin /> : <Navigate to="/login" replace />} />
        <Route path="/login" element={<Login />} />
      </Routes>
      {location.pathname !== '/panel' && <Footer />}
    </div>
  )
}


export default App
