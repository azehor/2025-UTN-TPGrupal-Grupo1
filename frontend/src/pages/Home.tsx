import type { FC } from "react";
import { Link } from 'react-router-dom';
import fondoProyecto from '../assets/fondo proyecto.avif';
import { GraduationCap, Gamepad2, Cpu } from "lucide-react";

const Home: FC = () => {
  return (
    <div className="bg-[#101c22] text-gray-300 font-['Space_Grotesk',sans-serif] flex flex-col min-h-screen overflow-x-hidden" style={{ overscrollBehavior: 'none' }}>
      <main className="flex-grow relative">
        {/* Hero Section */}
        <section className="relative py-20 md:py-32">
          {/* Admin Panel Button */}
          <Link 
            to="/panel" 
            className="absolute top-4 right-4 z-20 bg-[#13a4ec] hover:bg-[#0d7ec1] text-white font-bold py-3 px-6 rounded-lg text-sm shadow-lg border-2 border-[#13a4ec] hover:border-white transition-all duration-300 hover:scale-105"
          >
            Panel Admin
          </Link>
          
          <div 
            className="absolute inset-0 bg-cover bg-center"
            style={{
              backgroundImage: `linear-gradient(to bottom, rgba(16, 28, 34, 0.1), rgba(16, 28, 34, 1)), url("${fondoProyecto}")`
            }}
          ></div>
          
          <div className="container mx-auto px-6 text-center relative z-10">
            <h1 className="text-4xl md:text-6xl font-bold text-white mb-4">
              Encuentra la PC ideal para tus necesidades académicas, profesionales o de ocio 
            </h1>
            <p className="text-lg md:text-xl text-gray-400 max-w-3xl mx-auto mb-8">
              QuePC analiza tus requerimientos, valida las opciones del mercado y selecciona los componentes adecuados para recomendarte la mejor opción y ayudarte a rendir al máximo
            </p>
          </div>
        </section>

        {/* Features Section */}
        <section className="py-16 bg-[#101c22]">
          <div className="container mx-auto px-6">
            <div className="grid grid-cols-1 md:grid-cols-3 gap-8 text-center">

              {/* === Por Carrera === */}
              <Link 
                to="/busqueda-carrera"
                className="block bg-[#1a2831] p-8 rounded-xl shadow-sm hover:shadow-lg hover:shadow-[#13a4ec]/20 transition-all duration-300 border border-transparent hover:border-[#13a4ec]/30 hover:scale-[1.02]"
              >
                <div className="flex justify-center mb-4 text-[#13a4ec]">
                  <GraduationCap size={48} strokeWidth={1.5} />
                </div>
                <h3 className="text-2xl font-bold text-white mb-2">Por Carrera</h3>
                <p className="text-gray-400 mb-4">
                  Recomendaciones basadas en las demandas de tus programas académicos
                </p>
                <span className="text-[#13a4ec] font-bold hover:underline">
                  Explorar Carreras
                </span>
              </Link>

              {/* === Por Juego === */}
              <Link 
                to="/busqueda-juego"
                className="block bg-[#1a2831] p-8 rounded-xl shadow-sm hover:shadow-lg hover:shadow-[#13a4ec]/20 transition-all duration-300 border border-transparent hover:border-[#13a4ec]/30 hover:scale-[1.02]"
              >
                <div className="flex justify-center mb-4 text-[#13a4ec]">
                  <Gamepad2 size={48} strokeWidth={1.5} />
                </div>                
                <h3 className="text-2xl font-bold text-white mb-2">Por Juego</h3>
                <p className="text-gray-400 mb-4">
                  Encuentra una PC que pueda manejar tus juegos favoritos después de clases
                </p>
                <span className="text-[#13a4ec] font-bold hover:underline">
                  Explorar Juegos
                </span>
              </Link>

              {/* === Por Software === */}
              <Link 
                to="/busqueda-software"
                className="block bg-[#1a2831] p-8 rounded-xl shadow-sm hover:shadow-lg hover:shadow-[#13a4ec]/20 transition-all duration-300 border border-transparent hover:border-[#13a4ec]/30 hover:scale-[1.02]"
              >
                <div className="flex justify-center mb-4 text-[#13a4ec]">
                  <Cpu size={48} strokeWidth={1.5} />
                </div>               
                <h3 className="text-2xl font-bold text-white mb-2">Por Software</h3>
                <p className="text-gray-400 mb-4">
                  Asegúrate de que tu máquina ejecute el software específico que necesitas
                </p>
                <span className="text-[#13a4ec] font-bold hover:underline">
                  Explorar Software
                </span>
              </Link>

            </div>
          </div>
        </section>
      </main>
    </div>
  );
};

export default Home;
