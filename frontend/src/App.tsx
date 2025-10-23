import type { FC } from 'react'
import './App.css'
import Home from "./pages/Home"
import PanelDeAdmin from './pages/PanelDeAdmin'
import { Route, Routes, useLocation, Navigate } from 'react-router-dom'
import Login from './pages/Login'
import { isStaffAuthenticated } from './lib/auth'
import Footer from './components/Footer'
import { BusquedaSoftwarePage } from './pages/BusquedaSoftwarePage'
import { BusquedaCarreraPage } from './pages/BusquedaCarreraPage'
import { BusquedaJuegoPage } from './pages/BusquedaJuegoPage'
import { RecomendacionPage } from './pages/RecomendacionPage'
import { Header } from './components/Header'


const App: FC = () => {
  const location = useLocation();

  return (
    <div className="app-surface min-h-screen bg-gray-100">

      <Header />

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
