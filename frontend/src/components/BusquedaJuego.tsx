import type { FC } from "react";
import { useState, useEffect, useCallback, useRef } from "react";
import { apiService } from "../lib/api";
import { useBusqueda } from "../context/BusquedaContext";
import { useNavigate } from "react-router-dom";
import "./BusquedaJuego.css";

export const BusquedaJuego: FC = () => {
  const [query, setQuery] = useState("");
  const [juegos, setJuegos] = useState<Juego[]>([]);
  const [seleccionados, setSeleccionados] = useState<string[]>([]);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const mountedRef = useRef(true);

  type Juego = {
    id: string;
    nombre: string;
    imageURL: string;
    ordenGrafica?: number;
  };

  const GameTile: FC<{ j: Juego; seleccionado: boolean; onClick: () => void }> = ({ j, seleccionado, onClick }) => {
    const [imgError, setImgError] = useState(false);
    return (
      <div
        onClick={onClick}
        className={`group relative cursor-pointer rounded-xl overflow-hidden shadow-md transition-transform duration-300 hover:scale-105 hover:shadow-lg ${
          seleccionado ? "ring-2 ring-[#13a4ec]" : ""
        }`}
        style={{ width: "160px", height: "200px" }}
      >
        {!imgError ? (
          <img
            src={j.imageURL}
            alt={j.nombre}
            onError={() => setImgError(true)}
            className="w-full h-full object-cover object-center"
            loading="lazy"
          />
        ) : (
          <div className="w-full h-full bg-gray-800 flex items-end">
            <div className="w-full h-full bg-gradient-to-t from-black/70 via-black/0" />
          </div>
        )}
        <div className="absolute inset-0 bg-gradient-to-t from-black/70 via-black/0"></div>
        <p className="absolute bottom-2 left-2 text-white text-xs font-semibold truncate w-[90%]">
          {j.nombre}
        </p>
      </div>
    );
  };

  const fetchConRetry = useCallback(async (intento: number = 1, term: string = ""): Promise<void> => {
    if (!mountedRef.current) return;
    setLoading(true);
    setError(null);
    try {
      // Filtrado por tipo=juego desde el backend
      const softwares = await apiService.list("software", { tipo: "juego", q: term || undefined });
      const juegosApi = softwares.map<Juego>((s) => ({
        id: s.id,
        nombre: s.nombre,
        imageURL: s.image_url || "",
        ordenGrafica: s.orden_grafica,
      }));
      if (mountedRef.current) setJuegos(juegosApi);
    } catch (e) {
      console.error("Fallo fetch juegos intento", intento, e);
      if (intento < 3) {
        const delay = 500 * intento;
        setTimeout(() => fetchConRetry(intento + 1, term), delay);
      } else {
        if (mountedRef.current) setError("Error al cargar los juegos");
      }
    } finally {
      if (mountedRef.current) setLoading(false);
    }
  }, []);

  useEffect(() => {
    mountedRef.current = true;
    fetchConRetry(1, query);
    return () => { mountedRef.current = false; };
  }, [fetchConRetry]);

  // Debounce de la búsqueda al escribir
  useEffect(() => {
    const h = setTimeout(() => {
      fetchConRetry(1, query);
    }, 300);
    return () => clearTimeout(h);
  }, [query, fetchConRetry]);

  const juegosFiltrados = juegos.filter((j) =>
    j.nombre.toLowerCase().includes(query.toLowerCase())
  );

  const toggleSeleccionado = (id: string) => {
    setSeleccionados((prev) =>
      prev.includes(id) ? prev.filter((f) => f !== id) : [...prev, id]
    );
  };

  const { setBusqueda } = useBusqueda();
  const navigate = useNavigate();

  const recomendarPC = async () => {
    if (seleccionados.length === 0) return;
    setBusqueda("software", { ids: seleccionados });
    navigate("/recomendacion");
  };

  // Selecciona juegos "populares": top 6 por ordenGrafica (desc)
  const juegosPopulares = [...juegos]
    .sort((a, b) => (b.ordenGrafica || 0) - (a.ordenGrafica || 0))
    .slice(0, 6);

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
            placeholder="Buscar juego (ej. Cyberpunk 2077)"
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            className="w-full bg-[#101c22] border-2 border-gray-700 text-gray-200 placeholder-gray-500 rounded-lg py-3 pl-12 pr-4 focus:outline-none focus:border-[#13a4ec]"
          />
        </div>
      </div>

      {error && (
        <div className="text-center text-red-400 mb-6">
          <span>{error}</span>
          <button
            onClick={() => fetchConRetry()}
            disabled={loading}
            className="ml-3 inline-flex items-center px-3 py-1 rounded bg-gray-700 text-gray-100 hover:bg-gray-600 disabled:opacity-50"
          >
            Reintentar
          </button>
        </div>
      )}

      {/* Juegos Populares */}
      {loading ? (
        <div className="max-w-7xl mx-auto mb-12">
          <h2 className="text-2xl font-bold text-white mb-6">Juegos Populares</h2>
          <div className="flex flex-wrap justify-center gap-4">
            {Array.from({ length: 6 }).map((_, i) => (
              <div
                key={i}
                className="rounded-xl overflow-hidden shadow-md animate-pulse bg-gray-800"
                style={{ width: "160px", height: "200px" }}
              />
            ))}
          </div>
        </div>
      ) : query === "" && (
        <div className="max-w-7xl mx-auto mb-12">
          <h2 className="text-2xl font-bold text-white mb-6">
            Juegos Populares
          </h2>
          <div className="flex flex-wrap justify-center gap-4">
            {juegosPopulares.map((j) => (
              <GameTile
                key={j.id}
                j={j}
                seleccionado={seleccionados.includes(j.id)}
                onClick={() => toggleSeleccionado(j.id)}
              />
            ))}
          </div>
        </div>
      )}

      {/* Todos los Juegos */}
      <div className="max-w-7xl mx-auto">
        <h2 className="text-2xl font-bold text-white mb-6">Todos los juegos</h2>
        <div className="flex flex-wrap justify-center gap-4">
          {loading
            ? Array.from({ length: 12 }).map((_, i) => (
                <div
                  key={i}
                  className="rounded-xl overflow-hidden shadow-md animate-pulse bg-gray-800"
                  style={{ width: "160px", height: "200px" }}
                />
              ))
            : juegosFiltrados.map((j) => (
                <GameTile
                  key={j.id}
                  j={j}
                  seleccionado={seleccionados.includes(j.id)}
                  onClick={() => toggleSeleccionado(j.id)}
                />
              ))}
        </div>
      </div>

      {/* Botón Recomendar */}
      {seleccionados.length > 0 && (
        <button
          onClick={recomendarPC}
          disabled={loading}
          className="fixed bottom-5 right-5 bg-[#13a4ec] disabled:opacity-60 text-white px-6 py-3 rounded-full shadow-lg hover:bg-sky-500 transition-colors"
        >
          {loading ? "Cargando..." : `Recomendar PC (${seleccionados.length})`}
        </button>
      )}
    </div>
  );
};
