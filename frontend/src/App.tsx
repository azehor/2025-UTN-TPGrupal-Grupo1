import type { FC } from 'react'
import './App.css'
import Home from "./pages/Home"
import { Link, Route, Routes } from 'react-router-dom'
import { BusquedaSoftwarePage } from './pages/BusquedaSoftwarePage'
import { BusquedaCarreraPage } from './pages/BusquedaCarreraPage'


const App: FC = () => {
  return (
    <div>
      <nav
        style={{
          display: "flex",
          gap: "1rem",
          padding: "1rem",
          borderBottom: "1px solid #ccc",
          marginBottom: "1rem"
        }}>
        <Link to="/">Inicio</Link>
        <Link to="/busqueda-software">Busqueda Software</Link>
        <Link to="/busqueda-carrera">Busqueda Carrera</Link>
      </nav>

      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/busqueda-software" element={<BusquedaSoftwarePage />} /> 
        <Route path="/busqueda-carrera" element={<BusquedaCarreraPage />} />
      </Routes>
    </div>
  )
}


export default App
