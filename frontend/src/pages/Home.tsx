import type { FC } from "react";
import { Link } from 'react-router-dom'
import fondoProyecto from '../assets/fondo proyecto.avif'

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
            🔧 Panel Admin
          </Link>
          
          <div 
            className="absolute inset-0 bg-cover bg-center"
            style={{
              backgroundImage: `linear-gradient(to bottom, rgba(16, 28, 34, 0.1), rgba(16, 28, 34, 1)), url("${fondoProyecto}")`
            }}
          ></div>
          
          <div className="container mx-auto px-6 text-center relative z-10">
            <h1 className="text-4xl md:text-6xl font-bold text-white mb-4">
              Encuentra la PC perfecta para tus estudios
            </h1>
            <p className="text-lg md:text-xl text-gray-400 max-w-3xl mx-auto mb-8">
              Ya seas un estudiante de diseño que necesita una estación de trabajo potente o un estudiante de literatura que busca una laptop confiable, tenemos lo que necesitas. Explora nuestras recomendaciones adaptadas a tus necesidades específicas.
            </p>
          </div>
        </section>

        {/* Features Section */}
        <section className="py-16 bg-[#101c22]">
          <div className="container mx-auto px-6">
            <div className="grid grid-cols-1 md:grid-cols-3 gap-8 text-center">
              <div className="bg-[#1a2831] p-8 rounded-xl shadow-sm hover:shadow-lg hover:shadow-[#13a4ec]/20 transition-all duration-300 border border-transparent hover:border-[#13a4ec]/30">
                <h3 className="text-2xl font-bold text-white mb-2">Por Carrera</h3>
                <p className="text-gray-400 mb-4">
                  Recomendaciones basadas en las demandas de tu programa académico.
                </p>
                <Link 
                  to="/busqueda-carrera"
                  className="text-[#13a4ec] font-bold hover:underline"
                >
                  Explorar Carreras
                </Link>
              </div>
              
              <div className="bg-[#1a2831] p-8 rounded-xl shadow-sm hover:shadow-lg hover:shadow-[#13a4ec]/20 transition-all duration-300 border border-transparent hover:border-[#13a4ec]/30">
                <h3 className="text-2xl font-bold text-white mb-2">Por Juego</h3>
                <p className="text-gray-400 mb-4">
                  Encuentra una PC que pueda manejar tus juegos favoritos después de clases.
                </p>
                <Link 
                  to="/busqueda-juego"
                  className="text-[#13a4ec] font-bold hover:underline"
                >
                  Explorar Juegos
                </Link>
              </div>
              
              <div className="bg-[#1a2831] p-8 rounded-xl shadow-sm hover:shadow-lg hover:shadow-[#13a4ec]/20 transition-all duration-300 border border-transparent hover:border-[#13a4ec]/30">
                <h3 className="text-2xl font-bold text-white mb-2">Por Software</h3>
                <p className="text-gray-400 mb-4">
                  Asegúrate de que tu máquina ejecute el software específico que necesitas.
                </p>
                <Link 
                  to="/busqueda-software"
                  className="text-[#13a4ec] font-bold hover:underline"
                >
                  Explorar Software
                </Link>
              </div>
            </div>
          </div>
        </section>
      </main>
    </div>
  )
} 

export default Home
