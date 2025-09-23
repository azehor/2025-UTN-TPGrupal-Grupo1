import type { FC } from 'react'
import './App.css'
import Home from "./pages/Home"
import { Link, Route, Routes } from 'react-router-dom'
import { SoftwareSearchPage } from './pages/SoftwareSearchPage'
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
        <Link to="/software-search">Software Search</Link>
        <Link to="/busqueda-carrera">Busqueda Carrera</Link>
      </nav>

      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/software-search" element={<SoftwareSearchPage />} /> 
        <Route path="/busqueda-carrera" element={<BusquedaCarreraPage />} />
      </Routes>
    </div>
  )
}


export default App
