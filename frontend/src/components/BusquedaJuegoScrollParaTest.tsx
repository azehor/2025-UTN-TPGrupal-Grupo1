import type { FC } from "react";
import { useState, useEffect } from "react";
import type { Juego } from "../mock/juegosMock";
import { juegosHardCodeados } from "../mock/juegosMock";
import "./BusquedaJuego.css";

export const BusquedaJuegoCopy: FC = () => {
  const [query, setQuery] = useState("");
  const [juegos, setJuegos] = useState<Juego[]>([]);
  const [seleccionados, setSeleccionados] = useState<string[]>([]);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    setLoading(true);
    setTimeout(() => {
      try {
        setJuegos(juegosHardCodeados);
      } catch {
        setError("Error al cargar los juegos");
      } finally {
        setLoading(false);
      }
    }, 500);
  }, []);

  const juegosFiltrados = juegos.filter((j) =>
    j.nombre.toLowerCase().includes(query.toLowerCase())
  );

  const toggleSeleccionado = (id: string) => {
    setSeleccionados((prev) =>
      prev.includes(id) ? prev.filter((f) => f !== id) : [...prev, id]
    );
  };

  const enviarSeleccionados = () => {
    console.log("Juegos seleccionados:", seleccionados);
    alert("Juegos enviados al backend: " + seleccionados.join(", "));
  };

  const juegosPopulares = juegos.slice(0, 6);

  return (
    <div className="min-h-screen bg-[#101c22] text-gray-300 font-[Space Grotesk] p-6">
      {/* Buscador */}
      <div className="max-w-5xl mx-auto mb-10">
        <div className="relative">
          <div className="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none">
            <svg
              className="h-5 w-5 text-gray-500"
              fill="currentColor"
              viewBox="0 0 20 20"
            >
              <path
                clipRule="evenodd"
                fillRule="evenodd"
                d="M9 3.5a5.5 5.5 0 100 11 5.5 5.5 0 000-11zM2 9a7 7 0 1112.452 4.391l3.328 3.329a.75.75 0 11-1.06 1.06l-3.329-3.328A7 7 0 012 9z"
              />
            </svg>
          </div>
          <input
            type="text"
            placeholder="Search for a game (e.g. Cyberpunk 2077)"
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            className="w-full bg-[#101c22] border-2 border-gray-700 text-gray-200 placeholder-gray-500 rounded-lg py-3 pl-12 pr-4 focus:outline-none focus:border-[#13a4ec]"
          />
        </div>
      </div>

      {loading && <p className="text-center text-gray-400">Cargando juegos...</p>}
      {error && <p className="text-center text-red-400">{error}</p>}

      {/* Juegos Populares */}
      {query === "" && (
        <div className="max-w-7xl mx-auto mb-12">
          <h2 className="text-2xl font-bold text-white mb-6">
            Popular Games
          </h2>
          <div className="flex flex-wrap justify-center gap-4">
            {juegosPopulares.map((j) => {
              const seleccionado = seleccionados.includes(j.id);
              return (
                <div
                  key={j.id}
                  onClick={() => toggleSeleccionado(j.id)}
                  className={`group relative cursor-pointer rounded-xl overflow-hidden shadow-md transition-transform duration-300 hover:scale-105 hover:shadow-lg ${
                    seleccionado ? "ring-2 ring-[#13a4ec]" : ""
                  }`}
                  style={{
                    width: "160px",
                    height: "200px",
                  }}
                >
                  <div
                    className="w-full h-full bg-cover bg-center"
                    style={{ backgroundImage: `url(${j.imageURL})` }}
                  />
                  <div className="absolute inset-0 bg-gradient-to-t from-black/70 via-black/0"></div>
                  <p className="absolute bottom-2 left-2 text-white text-xs font-semibold truncate w-[90%]">
                    {j.nombre}
                  </p>
                </div>
              );
            })}
          </div>
        </div>
      )}

      {/* Todos los Juegos */}
      <div className="max-w-7xl mx-auto">
        <h2 className="text-2xl font-bold text-white mb-6">All Games</h2>
        <div className="flex flex-wrap justify-center gap-4 max-h-[600px] overflow-y-auto pr-1">
          {juegosFiltrados.map((j) => {
            const seleccionado = seleccionados.includes(j.id);
            return (
              <div
                key={j.id}
                onClick={() => toggleSeleccionado(j.id)}
                className={`flex flex-col gap-1 cursor-pointer rounded-xl overflow-hidden shadow-md transition-transform duration-300 hover:scale-105 hover:shadow-lg ${
                  seleccionado ? "ring-2 ring-green-500" : ""
                }`}
                style={{
                  width: "160px",
                  height: "200px",
                }}
              >
                <div
                  className="w-full h-full bg-cover bg-center"
                  style={{ backgroundImage: `url(${j.imageURL})` }}
                />
                <p className="text-xs font-medium text-gray-300 text-center truncate px-2 mt-1">
                  {j.nombre}
                </p>
              </div>
            );
          })}
        </div>
      </div>

      {/* Botón Enviar */}
      {seleccionados.length > 0 && (
        <button
          onClick={enviarSeleccionados}
          className="fixed bottom-5 right-5 bg-[#13a4ec] text-white px-6 py-3 rounded-full shadow-lg hover:bg-sky-500 transition-colors"
        >
          Enviar {seleccionados.length}
        </button>
      )}
    </div>
  );
};

