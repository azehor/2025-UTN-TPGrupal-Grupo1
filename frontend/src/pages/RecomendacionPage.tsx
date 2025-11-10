import { useEffect, useState } from "react";
import { useBusqueda } from "../context/BusquedaContext";
import { useNavigate } from "react-router-dom";
import getApiBase from "../lib/env";

export interface GenericItem {
  id: string;
  nombre: string;
  tipo?: string;
  empresa?: string;
  image_url?: string;
  orden_grafica?: number;
  orden_procesador?: number;
  orden_ram?: number;
  fabricante?: string;
  msrp?: number;
  form_factor?: string;
  max_largo_gpu_float?: number;
  socket?: string;
  consumo?: number;
  generacion?: string;
  nucleos?: number;
  largo?: number;
  modelo?: string;
  vram?: number;
  capacidad?: number;
  tipo_almacenamiento?: string;
  socket_procesador?: string;
  socket_ram?: string;
  precio?: number;
}

export const RecomendacionPage = () => {
  const { tipo, datos } = useBusqueda();
  const [items, setItems] = useState<GenericItem[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const navigate = useNavigate();

  useEffect(() => {
    if (!tipo || !datos) {
      navigate("/busqueda-carrera");
      return;
    }

    const fetchItems = async () => {
      setLoading(true);
      try {
        let res: Response;
        const API_BASE = getApiBase();

        if (tipo === "carrera") {
          res = await fetch(`${API_BASE}/recomendaciones-carrera/${datos.id}`);
        } else if (tipo === "software") {
          res = await fetch(`${API_BASE}/recomendaciones-softwares`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ ids: datos.ids as string[] }),
          });
        } else {
          throw new Error("Tipo de búsqueda no soportado");
        }

        if (!res.ok) throw new Error(`Error HTTP: ${res.status}`);
        const data: GenericItem[] = await res.json();
        setItems(data);
      } catch (err: any) {
        setError(err.message || "Error al obtener recomendaciones");
      } finally {
        setLoading(false);
      }
    };

    fetchItems();
  }, [tipo, datos, navigate]);

  const iconMap: Record<string, string> = {
    Procesador: "🧠",
    GPU: "🎮",
    Disco: "💾",
    RAM: "📦",
    Motherboard: "🖥️",
    PSU: "⚡",
    Gabinete: "🧱",
  };

  return (
    <main className="min-h-screen bg-[#101c22] text-gray-300 font-['Space_Grotesk',sans-serif] py-16">
      <div className="max-w-4xl mx-auto px-6">
        {/* Título principal */}
        <div className="text-center mb-12">
          <h2 className="text-4xl md:text-5xl font-bold text-white">
            Tu configuración de PC recomendada
          </h2>
          <p className="mt-4 text-lg text-gray-400">
            Esta es la configuración óptima de PC según tus necesidades.
          </p>
        </div>

        {/* Tarjeta principal */}
        <div className="bg-[#1a2831] rounded-2xl shadow-xl p-8 space-y-8">
          {loading && (
            <p className="text-center text-gray-400">Cargando recomendaciones...</p>
          )}
          {error && <p className="text-center text-red-500">{error}</p>}

          {/* Lista de componentes */}
          <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
            {items.map((item) => (
              <div
                key={item.id}
                className="flex flex-col bg-[#101c22]/80 rounded-xl p-5 hover:bg-[#101c22] transition"
              >
                <div className="flex items-center gap-2 mb-1">
                  <span className="text-2xl">
                    {iconMap[item.tipo || "Processor"] || "💻"}
                  </span>
                  <h3 className="text-lg font-semibold text-[#13a4ec]">
                    {item.tipo || "Componente"}
                  </h3>
                </div>
                <p className="font-medium text-white">{item.nombre}</p>
                <p className="text-sm text-gray-400 mt-1">
                  {item.empresa ? `Fabricante: ${item.empresa}` : "—"}
                </p>
                {item.precio && (
                  <p className="text-sm text-gray-500 mt-1">Precio: ${item.precio}</p>
                )}
              </div>
            ))}

            {!loading && items.length === 0 && !error && (
              <p className="col-span-2 text-center text-gray-400">
                No se encontraron recomendaciones.
              </p>
            )}
          </div>

          {/* Botones inferiores */}
          <div className="flex flex-col sm:flex-row justify-center gap-4 pt-8 border-t border-gray-700">
            <button className="bg-[#13a4ec] text-white font-semibold py-3 px-6 rounded-lg shadow-md hover:bg-[#0d7ec1] transition-all transform hover:scale-105 flex items-center gap-2 justify-center">
              Guardar recomendación
            </button>
            <button className="bg-[#24333f] text-white font-semibold py-3 px-6 rounded-lg shadow-md hover:bg-[#2e4455] transition-all flex items-center gap-2 justify-center">
              Imprimir
            </button>
          </div>
        </div>
      </div>
    </main>
  );
};