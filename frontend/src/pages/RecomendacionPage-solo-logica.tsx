import { useEffect, useState } from "react";
import { useBusqueda } from "../context/BusquedaContext"; // ajusta la ruta
import { useNavigate } from "react-router-dom";

export interface GenericItem {
  id: string;
  nombre: string;

  // Software fields
  tipo?: string;
  empresa?: string;
  image_url?: string;
  orden_grafica?: number;
  orden_procesador?: number;
  orden_ram?: number;

  // Component fields
  fabricante?: string;
  msrp?: number;

  // Gabinete
  form_factor?: string;
  max_largo_gpu_float?: number;
  socket?: string;

  // Procesador
  consumo?: number;
  generacion?: string;
  nucleos?: number;

  // GPU
  largo?: number;
  modelo?: string;
  vram?: number;

  // Disco
  capacidad?: number;
  tipo_almacenamiento?: string;

  // Motherboard  
  socket_procesador?: string;
  socket_ram?: string;

  // RAM
  // capacidad ya definida arriba
  // generacion ya definida arriba

  // Otros campos de respuesta
  precio?: string;
  orden?: number;
}

const camposVisibles = ["id", "nombre", "tipo", "precio"];

const filtrarCamposVisibles = (item: GenericItem) => {
  const result: Partial<GenericItem> = {};
  camposVisibles.forEach((campo) => {
    if (item[campo as keyof GenericItem] !== undefined) {
      result[campo as keyof GenericItem] = item[campo as keyof GenericItem];
    }
  });
  return result;
};

export const useRecomendaciones = () => {
  const { tipo, datos } = useBusqueda();
  const [recomendaciones, setRecomendaciones] = useState<GenericItem[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const navigate = useNavigate();

  useEffect(() => {
    if (!tipo || !datos) {
      navigate("/busqueda-carrera");
      return;
    }

    const fetchRecomendaciones = async () => {
      setLoading(true);
      try {
        const res = await fetch(`http://localhost:8080/v1/recomendaciones-carrera/${datos.id}`);
        if (!res.ok) throw new Error(`Error HTTP: ${res.status}`);
        const data: GenericItem[] = await res.json();

        // Filtrar solo campos visibles
        const filtrados = data.map(filtrarCamposVisibles) as GenericItem[];
        setRecomendaciones(filtrados);
      } catch (err: any) {
        setError(err.message || "Error al obtener recomendaciones");
        console.error(err);
      } finally {
        setLoading(false);
      }
    };

    fetchRecomendaciones();
  }, [tipo, datos, navigate]);

  return { recomendaciones, loading, error, tipo, datos };
};


