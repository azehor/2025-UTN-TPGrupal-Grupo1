import type { FC } from 'react'
import './App.css'
import Home from "./pages/Home"
import { Link, Route, Routes, useLocation } from 'react-router-dom'
import { BusquedaSoftwarePage } from './pages/BusquedaSoftwarePage'
import { BusquedaCarreraPage } from './pages/BusquedaCarreraPage'
import { BusquedaJuegoPage } from './pages/BusquedaJuegoPage'
import { RecomendacionPage } from './pages/RecomendacionPage'


const App: FC = () => {
  const location = useLocation();
  const isHome = location.pathname === "/";

  return (
    <div>
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
      </Routes>
    </div>
  )
}


export default App
