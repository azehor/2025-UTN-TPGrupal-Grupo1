import { type FC } from "react";
import { NavLink } from "react-router-dom";
import logo from "../assets/logo.png";

export const Header: FC = () => {
  return (
    <header className="bg-[#101c22]/80 backdrop-blur-sm sticky top-0 z-50 border-b border-gray-700/50">
      <nav className="container mx-auto px-6 py-4 flex items-center justify-between">
        <NavLink to="/" className="flex items-center gap-2 md:gap-3">
          <img src={logo} alt="QuéPC Logo" className="h-6 w-6 md:h-8 md:w-8 object-contain" />
          <span className="text-lg md:text-2xl font-bold text-white">QuéPC</span>
        </NavLink>
        
        <div className="flex items-center gap-4 md:gap-8">
          <NavLink 
            to="/busqueda-carrera" 
            className={({ isActive }) =>
              isActive
                ? "text-[#13a4ec] transition-colors text-sm md:text-base"
                : "text-gray-400 hover:text-[#13a4ec] transition-colors text-sm md:text-base"
            }
          >
            Por Carrera
          </NavLink>
          <NavLink 
            to="/busqueda-juego" 
            className={({ isActive }) =>
              isActive
                ? "text-[#13a4ec] transition-colors text-sm md:text-base"
                : "text-gray-400 hover:text-[#13a4ec] transition-colors text-sm md:text-base"
            }
          >
            Por Juego
          </NavLink>
          <NavLink 
            to="/busqueda-software" 
            className={({ isActive }) =>
              isActive
                ? "text-[#13a4ec] transition-colors text-sm md:text-base"
                : "text-gray-400 hover:text-[#13a4ec] transition-colors text-sm md:text-base"
            }
          >
            Por Software
          </NavLink>
        </div>
        
        <div className="flex items-center">
          <button className="p-2 rounded-full text-gray-400 hover:bg-gray-800 transition-colors">
            <span className="material-symbols-outlined text-lg">light_mode</span>
          </button>
        </div>
      </nav>
    </header>
  );
};
