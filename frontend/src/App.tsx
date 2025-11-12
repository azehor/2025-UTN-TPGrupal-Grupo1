import type { FC } from 'react';
import { Route, Routes, useLocation, Navigate } from 'react-router-dom';
import Home from './pages/Home';
import PanelDeAdmin from './pages/PanelDeAdmin';
import Login from './pages/Login';
import { isStaffAuthenticated } from './lib/auth';
import Footer from './components/Footer';
import { BusquedaSoftwarePage } from './pages/BusquedaSoftwarePage';
import { BusquedaCarreraPage } from './pages/BusquedaCarreraPage';
import { BusquedaJuegoPage } from './pages/BusquedaJuegoPage';
import { RecomendacionPage } from './pages/RecomendacionPage';
import { Header } from './components/Header';
import { useTheme } from './context/ThemeContext';

const App: FC = () => {
  const location = useLocation();
  const { theme } = useTheme();

  return (
    <div className={`app-surface flex flex-col ${theme}`}>
      <Header />

      <main className="flex-grow">
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/busqueda-software" element={<BusquedaSoftwarePage />} />
          <Route path="/busqueda-carrera" element={<BusquedaCarreraPage />} />
          <Route path="/busqueda-juego" element={<BusquedaJuegoPage />} />
          <Route path="/recomendacion" element={<RecomendacionPage />} />
          <Route
            path="/panel"
            element={
              isStaffAuthenticated() ? <PanelDeAdmin /> : <Navigate to="/login" replace />
            }
          />
          <Route path="/login" element={<Login />} />
        </Routes>
      </main>

      {location.pathname !== '/panel' && <Footer />}
    </div>
  );
};

export default App;