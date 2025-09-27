import React, { useEffect, useState } from "react";

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

  useEffect(() => {
    const obtenerListadoCarreraApi = async () => {
      setLoading(true);
      try {
        const res = await fetch("http://localhost:8080/v1/carreras"); // Endpoint de api en go local, ver luego en .env de manejar estas urls
        if (!res.ok) {
          throw new Error(`Error HTTP: ${res.status}`);
        }
        const data: Carrera[] = await res.json();
        setCarreraList(data);
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
    // aca se puede llegar a llamar otra api futura
    console.log("Botón clickeado:", carrera);
  };

  return (
    <div style={{ fontFamily: "sans-serif", padding: "1rem", maxWidth: "900px", margin: "0 auto" }}>
      <h2>Buscar por carrera</h2>

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

      {loading && <p>Cargando...</p>}
      {error && <p style={{ color: "red" }}>{error}</p>}

      <div
        style={{
          display: "grid",
          gridTemplateColumns: "repeat(auto-fill, minmax(180px, 1fr))",
          gap: "12px"
        }}
      >
        {filteredList.length > 0 ? (
          filteredList.map((carrera) => (
            <div
              key={carrera.id}
              style={{
                border: "1px solid #ddd",
                borderRadius: "8px",
                padding: "8px",
                display: "flex",
                flexDirection: "column",
                alignItems: "stretch",
                background: "#fff",
                boxShadow: selectedId === carrera.id ? "0 0 0 3px rgba(0,123,255,0.12)" : undefined
              }}
            >
              {carrera.imageURL ? (
                <img
                  src={carrera.imageURL}
                  alt={carrera.nombre}
                  style={{ width: "100%", height: "110px", objectFit: "cover", borderRadius: "6px", marginBottom: "8px" }}
                />
              ) : (
                <div style={{ width: "100%", height: "110px", background: "#f4f4f4", borderRadius: "6px", marginBottom: "8px" }} />
              )}

              <div style={{ flex: 1, marginBottom: "8px", fontWeight: 600 }}>{carrera.nombre}</div>

              <button
                onClick={() => handleAccion(carrera)}
                style={{
                  padding: "8px",
                  borderRadius: "6px",
                  border: "none",
                  background: selectedId === carrera.id ? "#0056b3" : "#007bff",
                  color: "#fff",
                  cursor: "pointer"
                }}
              >
                Seleccionar
              </button>
            </div>
          ))
        ) : (
          <div style={{ gridColumn: "1/-1" }}>
            <p>No se encontró carrera</p>
          </div>
        )}
      </div>
    </div>
  );
};