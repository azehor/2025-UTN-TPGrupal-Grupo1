import React, { useEffect, useState } from "react";

interface Carrera {
  id: string;
  nombre: string;
  empresa: string;
  imageURL: string;
}

export const BusquedaCarrera: React.FC = () => {
  const [query, setQuery] = useState("");
  const [carreraList, setCarreraList] = useState<Carrera[]>([]);
  const [error, setError] = useState<string | null>(null);

  const filteredList = carreraList.filter(soft =>
    soft.nombre.toLowerCase().includes(query.toLowerCase())
  );

  useEffect(() => {
    const obtenerListadoCarreraApi = async () => {
      try {
        const res = await fetch("http://localhost:8080/v1/carreras"); // Endpoint de api en go local, ver luego en .env de manejar estas urls 
        if (!res.ok) {
          throw new Error(`Error HTTP: ${res.status}`);
        }
        const data: Carrera[] = await res.json();
        setCarreraList(data);
      } catch (err: any) {
        setError(err.message || "Error al obtener listado de carreras");
        console.error(error);
      } 
    };

    obtenerListadoCarreraApi();
  }, []);
  
  return (
    <div style={{ fontFamily: "sans-serif", padding: "1rem", maxWidth: "400px" }}>
      <h2>Busqueda Carrera</h2>

      <input
        type="text"
        placeholder="Type to search..."
        value={query}
        onChange={e => setQuery(e.target.value)}
        style={{
          width: "100%",
          padding: "0.5rem",
          marginBottom: "1rem",
          border: "1px solid #ccc",
          borderRadius: "6px"
        }}
      />

      <ul>
        {filteredList.length > 0 ? (
          filteredList.map((carrera, index) => <li key={index}>{carrera.nombre}</li>)
        ) : (
          <li>No se encontro carrera </li>
        )}
      </ul>
    </div>
  );
};