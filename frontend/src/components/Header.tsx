import { FC } from "react";
import { Sun } from "lucide-react";
import { NavLink } from "react-router-dom";
import QuePcLogo from "../assets/logo.png";

export const Header: FC = () => {
  return (
    <header className="bg-[#101c22] text-white flex items-center justify-between px-6 py-3 shadow-md">
      <NavLink to="/" className="flex items-center gap-2">
        <img src={QuePcLogo} alt="Logo QuéPC" className="w-8 h-8 object-contain" />
        <h1 className="text-lg font-semibold">QuéPC</h1>
      </NavLink>

      <nav className="flex gap-6 text-sm text-gray-300">
        <NavLink
            to="/busqueda-carrera"
            className={({ isActive }) =>
            isActive
              ? "text-white border-b-2 border-white pb-1 transition"
              : "text-gray-300 hover:text-white transition"
          }
        >
          Por Carrera
        </NavLink>
        <NavLink to="/busqueda-juego" 
        className={({ isActive }) =>
            isActive
                ? "text-white border-b-2 border-white pb-1 transition"
                : "text-gray-300 hover:text-white transition"
            }
        >
          Por Juego
        </NavLink>
        <NavLink to="/busqueda-software" 
        className={({ isActive }) =>
            isActive
                ? "text-white border-b-2 border-white pb-1 transition"
                : "text-gray-300 hover:text-white transition"
            }
        >
          Por Software
        </NavLink>
      </nav>

      <button className="p-2 hover:bg-slate-800 rounded-full transition">
        <Sun size={18} />
      </button>
    </header>
  );
};
