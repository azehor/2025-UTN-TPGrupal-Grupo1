import type { FC } from 'react'
import './App.css'
import Home from "./pages/Home"
import { Link, Route, Routes, useLocation } from 'react-router-dom'
import { BusquedaSoftwarePage } from './pages/BusquedaSoftwarePage'
import { BusquedaCarreraPage } from './pages/BusquedaCarreraPage'
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
            marginBottom: "1rem"
          }}>
          <Link to="/">Inicio</Link>
          <Link to="/busqueda-software">Buscar por Software</Link>
          <Link to="/busqueda-carrera">Buscar por Carrera</Link>
        </nav>
      )}

      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/busqueda-software" element={<BusquedaSoftwarePage />} />
        <Route path="/busqueda-carrera" element={<BusquedaCarreraPage />} />
        <Route path="/recomendacion" element={<RecomendacionPage />} />
      </Routes>
    </div>
  )
}


export default App
