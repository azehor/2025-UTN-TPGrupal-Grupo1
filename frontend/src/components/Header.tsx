import { FC } from 'react';
import { NavLink } from 'react-router-dom';
import logo from '../assets/logo.png';
import { useTheme } from '../context/ThemeContext';

export const Header: FC = () => {
  const { theme, toggleTheme } = useTheme();

  return (
    <header className="site-header sticky top-0 z-50 border-b border-gray-700/50 backdrop-blur-sm">
      <nav className="container mx-auto px-6 py-4 flex items-center justify-between">
        <NavLink to="/" className="flex items-center gap-3">
          <img src={logo} alt="QuéPC Logo" className="h-8 w-8 object-contain" />
          <span className="text-lg md:text-2xl font-bold">QuéPC</span>
        </NavLink>

        <div className="flex items-center gap-6">
          <NavLink
            to="/busqueda-carrera"
            className={({ isActive }) =>
              isActive ? 'nav-link font-medium' : 'nav-link opacity-80 hover:opacity-100'
            }
          >
            Por Carrera
          </NavLink>
          <NavLink
            to="/busqueda-juego"
            className={({ isActive }) =>
              isActive ? 'nav-link font-medium' : 'nav-link opacity-80 hover:opacity-100'
            }
          >
            Por Juego
          </NavLink>
          <NavLink
            to="/busqueda-software"
            className={({ isActive }) =>
              isActive ? 'nav-link font-medium' : 'nav-link opacity-80 hover:opacity-100'
            }
          >
            Por Software
          </NavLink>
        </div>

        <div className="flex items-center">
          <button
            onClick={toggleTheme}
            className="p-2 rounded-full hover:bg-gray-700/30 transition-colors"
            title={theme === 'dark' ? 'Cambiar a modo claro' : 'Cambiar a modo oscuro'}
          >
            <span className="material-symbols-outlined text-lg" style={{ color: 'var(--text-color)' }}>
              {theme === 'dark' ? 'light_mode' : 'dark_mode'}
            </span>
          </button>
        </div>
      </nav>
    </header>
  );
};