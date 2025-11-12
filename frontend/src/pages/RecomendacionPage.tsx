import { useEffect, useMemo, useState } from "react";
import { useBusqueda } from "../context/BusquedaContext";
import { useNavigate } from "react-router-dom";
import getApiBase from "../lib/env";
import { Range } from "react-range";

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

  const [minBudgetStr, setMinBudgetStr] = useState<string>("");
  const [maxBudgetStr, setMaxBudgetStr] = useState<string>("");

  // Gasto total (suma de item.precio)
  const precioTotal = useMemo(() => {
    return items.reduce((acc, it) => acc + (it.precio ?? 0), 0);
  }, [items]);

  const formatCurrency = (v: number) => `$${v.toFixed(2)}`;

  useEffect(() => {
    if (!tipo || !datos) {
      navigate("/busqueda-carrera");
      return;
    }

    const fetchItems = async (minBudget?: number, maxBudget?: number) => {
      setLoading(true);
      setError(null);
      try {
        let res: Response;
        const API_BASE = getApiBase();

        if (tipo === "carrera") {
          const params = new URLSearchParams();
          if (minBudget !== undefined) params.append("minBudget", String(minBudget));
          if (maxBudget !== undefined) params.append("maxBudget", String(maxBudget));
          const qs = params.toString();
          const url = qs
            ? `${API_BASE}/recomendaciones-carrera/${datos.id}?${qs}`
            : `${API_BASE}/recomendaciones-carrera/${datos.id}`;
          res = await fetch(url);
        } else if (tipo === "software") {
          const payload: Record<string, any> = { ids: datos.ids as string[] };
          if (minBudget !== undefined) payload.minBudget = minBudget;
          if (maxBudget !== undefined) payload.maxBudget = maxBudget;
          res = await fetch(`${API_BASE}/recomendaciones-softwares`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(payload),
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

  const parseBudget = (s: string): number | undefined => {
    const v = s.trim();
    if (!v) return undefined;
    const n = Number(v);
    return Number.isFinite(n) ? n : undefined;
  };

  const applyBudgetFilter = async () => {
    setError(null);
    setItems([]);
    const min = parseBudget(minBudgetStr);
    const max = parseBudget(maxBudgetStr);

    setLoading(true);
    try {
      let res: Response;
      const API_BASE = getApiBase();

      if (tipo === "carrera") {
        const params = new URLSearchParams();
        if (min !== undefined) params.append("minBudget", String(min));
        if (max !== undefined) params.append("maxBudget", String(max));
        const qs = params.toString();
        const url = qs
          ? `${API_BASE}/recomendaciones-carrera/${datos.id}?${qs}`
          : `${API_BASE}/recomendaciones-carrera/${datos.id}`;
        res = await fetch(url);
      } else if (tipo === "software") {
        const payload: Record<string, any> = { ids: datos.ids as string[] };
        if (min !== undefined) payload.minBudget = min;
        if (max !== undefined) payload.maxBudget = max;
        res = await fetch(`${API_BASE}/recomendaciones-softwares`, {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify(payload),
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

        {/*Nuevo filtro con slider */}
        <div className="max-w-md mx-auto mb-6 bg-[#101c22]/80 p-4 rounded-lg border border-gray-700">
          <h4 className="text-sm text-gray-300 mb-2 text-center">
            Filtrar por presupuesto (USD)
          </h4>

          <Range
            step={10}
            min={500}
            max={10000}
            values={[
              parseFloat(minBudgetStr) || 500,
              parseFloat(maxBudgetStr) || 10000,
            ]}
            onChange={(values) => {
              setMinBudgetStr(values[0].toString());
              setMaxBudgetStr(values[1].toString());
            }}
            renderTrack={({ props, children }) => {
              const [minVal, maxVal] = [
                parseFloat(minBudgetStr) || 500,
                parseFloat(maxBudgetStr) || 10000,
              ];
              const min = 500;
              const max = 10000;

              return (
                <div
                  {...props}
                  className="h-2 bg-gray-700 rounded-lg relative"
                  style={props.style}
                >
                  <div
                    className="absolute h-2 bg-blue-500 rounded-lg"
                    style={{
                      left: `${((minVal - min) / (max - min)) * 100}%`,
                      right: `${100 - ((maxVal - min) / (max - min)) * 100}%`,
                    }}
                  />
                  {children}
                </div>
              );
            }}
            renderThumb={({ props }) => (
              <div
                {...props}
                className="w-5 h-5 bg-blue-500 rounded-full shadow-md cursor-pointer hover:scale-110 transition-transform"
              />
            )}
          />

          <p className="text-center text-gray-400 mt-3">
            ${minBudgetStr || 500} — ${maxBudgetStr || 2000}
          </p>

          <div className="mt-3 flex justify-center">
            <button
              onClick={applyBudgetFilter}
              className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
            >
              Actualizar recomendaciones
            </button>
          </div>
        </div>

        {/* Tarjeta principal */}
        <div className="bg-[#1a2831] rounded-2xl shadow-xl p-8 space-y-8">
          {loading && (
            <p className="text-center text-gray-400">
              Cargando recomendaciones...
            </p>
          )}
          {error && <p className="text-center text-red-500">{error}</p>}

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
                  <p className="text-sm text-gray-500 mt-1">
                    Precio: ${item.precio}
                  </p>
                )}
              </div>
            ))}

            {!loading && items.length === 0 && !error && (
              <p className="col-span-2 text-center text-gray-400">
                No se encontraron recomendaciones.
              </p>
            )}
          </div>

          {/* Gasto total */}
          <div className="flex items-center justify-end text-white font-semibold mt-4">
            <span className="text-sm text-gray-300 mr-3">Gasto total:</span>
            <span className="text-xl">{formatCurrency(precioTotal)}</span>
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