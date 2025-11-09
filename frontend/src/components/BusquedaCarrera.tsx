import React, { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { useBusqueda } from  "../context/BusquedaContext";
import { getApiBase } from "../lib/env";

interface CarreraApiDTO {
  id: string;
  nombre: string;
  image_url?: string;
}

interface Carrera {
  id: string;
  nombre: string;
  imageURL: string;
}

export const BusquedaCarrera: React.FC = () => {
  const [query, setQuery] = useState("");
  const [carreraList, setCarreraList] = useState<Carrera[]>([]);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const [selectedId, setSelectedId] = useState<string | null>(null);

  const filteredList = carreraList.filter(soft =>
    soft.nombre.toLowerCase().includes(query.toLowerCase())
  );

  const { setBusqueda } = useBusqueda();
  const navigate = useNavigate();

  useEffect(() => {
    const obtenerListadoCarreraApi = async () => {
      setLoading(true);
      try {
        const base = getApiBase();
        const url = `${base.replace(/\/+$/, '')}/carreras`;
        const res = await fetch(url); // URL tomada de config (runtime/build) sin hardcodear
        if (!res.ok) {
          throw new Error(`Error HTTP: ${res.status}`);
        }
        const data: CarreraApiDTO[] = await res.json();
        const mapped: Carrera[] = data.map((c) => ({
          id: c.id,
          nombre: c.nombre,
          imageURL: c.image_url || "",
        }));
        setCarreraList(mapped);
      } catch (err: any) {
        setError(err.message || "Error al obtener listado de carreras");
        console.error(err);
      } finally {
        setLoading(false);
      }
    };

    obtenerListadoCarreraApi();
  }, []);

  const handleAccion = (carrera: Carrera) => {
    setSelectedId(carrera.id);
    setBusqueda("carrera", carrera);
    navigate("/recomendacion")
  };

  return (
    <div className="min-h-screen bg-[#101c22] text-gray-300 font-['Space_Grotesk',sans-serif] py-8">
      <div className="container mx-auto px-6 max-w-6xl">
        <h2 className="text-3xl md:text-4xl font-bold text-white mb-8 text-center">Buscar por Carrera</h2>
        
        <div className="mb-8">
          <input
            type="text"
            placeholder="Busca tu carrera..."
            value={query}
            onChange={e => setQuery(e.target.value)}
            className="w-full px-4 py-3 bg-[#1a2831] border border-gray-600 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:border-[#13a4ec] focus:ring-2 focus:ring-[#13a4ec]/20 transition-all"
          />
        </div>

        {loading && (
          <div className="text-center py-8">
            <div className="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-[#13a4ec]"></div>
            <p className="mt-4 text-gray-400">Cargando carreras...</p>
          </div>
        )}
        
        {error && (
          <div className="bg-red-900/20 border border-red-500/50 text-red-300 px-4 py-3 rounded-lg mb-6">
            {error}
          </div>
        )}

        <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">
          {filteredList.length > 0 ? (
            filteredList.map((carrera) => (
              <div
                key={carrera.id}
                className={`bg-[#1a2831] border rounded-xl p-4 flex flex-col transition-all duration-300 hover:shadow-lg hover:shadow-[#13a4ec]/20 ${
                  selectedId === carrera.id 
                    ? 'border-[#13a4ec] shadow-lg shadow-[#13a4ec]/20' 
                    : 'border-gray-600 hover:border-[#13a4ec]/50'
                }`}
              >
                {carrera.imageURL ? (
                  <img
                    src={carrera.imageURL}
                    alt={carrera.nombre}
                    className="w-full h-32 object-cover rounded-lg mb-4"
                  />
                ) : (
                  <div className="w-full h-32 bg-gray-700 rounded-lg mb-4 flex items-center justify-center">
                    <span className="text-gray-500 text-sm">Sin imagen</span>
                  </div>
                )}

                <div className="flex-1 mb-4">
                  <h3 className="font-semibold text-white text-lg leading-tight">{carrera.nombre}</h3>
                </div>
                
                {/* <Link to="/recomendacion"
                  className="w-full"
                  state={{ carreraSeleccionada: carrera }}> */}
                  <button
                    onClick={() => handleAccion(carrera)}
                    className={`w-full py-2 px-4 rounded-lg font-medium transition-all ${
                      selectedId === carrera.id 
                        ? 'bg-[#0d7ec1] text-white' 
                        : 'bg-[#13a4ec] text-white hover:opacity-90'
                    }`}
                  >
                    Seleccionar
                  </button>
                {/* </Link> */}
              </div>
            ))
          ) : !loading && (
            <div className="col-span-full text-center py-12">
              <p className="text-gray-400 text-lg">
                {query.trim() ? 'No se encontraron carreras que coincidan con tu búsqueda' : 'No hay carreras disponibles'}
              </p>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};