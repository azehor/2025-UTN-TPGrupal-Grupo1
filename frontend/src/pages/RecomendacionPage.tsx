import React, { useEffect, useState } from "react";
import { useBusqueda } from "../context/BusquedaContext"; // ajusta la ruta
import { useNavigate } from "react-router-dom";

interface Carrera {
  id: string;
  nombre: string;
  imageURL: string;
}

interface Recomendacion {
  id: string;
  nombre: string;
  descripcion: string;
  imageURL?: string;
}

export const RecomendacionPage: React.FC = () => {
  const { tipo, datos } = useBusqueda(); // "carrera" y el objeto seleccionado
  const [recomendaciones, setRecomendaciones] = useState<Recomendacion[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const navigate = useNavigate();

  useEffect(() => {
    // Si no hay datos, volvemos a la búsqueda
    if (!tipo || !datos) {
      navigate("/busqueda-carrera"); // ruta por defecto
      return;
    }

    const fetchRecomendaciones = async () => {
      setLoading(true);
      try {
        // Ejemplo: GET con id de carrera
        const res = await fetch(`http://localhost:8080/v1/recomendaciones-carrera/${datos.id}`);
        if (!res.ok) throw new Error(`Error HTTP: ${res.status}`);
        const data: Recomendacion[] = await res.json();
        setRecomendaciones(data);
      } catch (err: any) {
        setError(err.message || "Error al obtener recomendaciones");
        console.error(err);
      } finally {
        setLoading(false);
      }
    };

    fetchRecomendaciones();
  }, [tipo, datos, navigate]);

  if (!tipo || !datos) {
    return <p>Redirigiendo a búsqueda...</p>;
  }

  return (
    <div className="min-h-screen bg-[#101c22] text-gray-300 font-['Space_Grotesk',sans-serif] py-8">
      <div className="container mx-auto px-6 max-w-4xl">
        <h2 className="text-3xl md:text-4xl font-bold text-white mb-8 text-center">
          Recomendaciones para: {(datos as Carrera).nombre}
        </h2>

        {loading && <p className="text-gray-400 text-center">Cargando recomendaciones...</p>}
        {error && <p className="text-red-400 text-center">{error}</p>}

        <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-6">
          {recomendaciones.map((r) => (
            <div key={r.id} className="bg-[#1a2831] border border-gray-600 rounded-lg p-4">
              {r.imageURL && (
                <img
                  src={r.imageURL}
                  alt={r.nombre}
                  className="w-full h-32 object-cover rounded-lg mb-4"
                />
              )}
              <h3 className="font-semibold text-white text-lg">{r.nombre}</h3>
              <p className="text-gray-400">{r.descripcion}</p>
            </div>
          ))}
        </div>

        {!loading && recomendaciones.length === 0 && (
          <p className="text-gray-400 text-center mt-6">
            No hay recomendaciones disponibles para esta carrera.
          </p>
        )}
      </div>
    </div>
  );
};
