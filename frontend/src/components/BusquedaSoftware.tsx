import React, { useEffect, useState } from "react";

interface Software {
  id: string;
  nombre: string;
  empresa: string;
  imageURL: string;
}

export const BusquedaSoftware: React.FC = () => {
  const [query, setQuery] = useState("");
  const [softwareList, setSoftwareList] = useState<Software[]>([]);
  const [error, setError] = useState<string | null>(null);

  const filteredList = softwareList.filter(soft =>
    soft.nombre.toLowerCase().includes(query.toLowerCase())
  );

  useEffect(() => {
    const obtenerListadoSoftwaresApi = async () => {
      try {
        const res = await fetch("http://localhost:8080/v1/softwares"); // Endpoint de api en go local, ver luego en .env de manejar estas urls 
        if (!res.ok) {
          throw new Error(`Error HTTP: ${res.status}`);
        }
        const data: Software[] = await res.json();
        setSoftwareList(data);
      } catch (err: any) {
        setError(err.message || "Error al obtener listado de software");
        console.error(error);
      } 
    };

    obtenerListadoSoftwaresApi();
  }, []);
  
  return (
    <div style={{ fontFamily: "sans-serif", padding: "1rem", maxWidth: "400px" }}>
      <h2>Busqueda por Software</h2>

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
          filteredList.map((soft, index) => <li key={index}>{soft.nombre}</li>)
        ) : (
          <li>Software no encontrado</li>
        )}
      </ul>
    </div>
  );
};
