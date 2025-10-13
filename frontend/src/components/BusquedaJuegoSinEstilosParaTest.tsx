import type { FC } from "react";
import { useState, useEffect } from "react";
import type { Juego } from "../mock/juegosMock";
import { juegosHardCodeados } from "../mock/juegosMock";

export const BusquedaJuegoSinEstilos: FC = () => {
  const [query, setQuery] = useState("");
  const [juegos, setJuegos] = useState<Juego[]>([]);
  const [seleccionados, setSeleccionados] = useState<string[]>([]);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  // Carga inicial de juegos
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

  // Filtrado según la búsqueda
  const juegosFiltrados = juegos.filter((j) =>
    j.nombre.toLowerCase().includes(query.toLowerCase())
  );

  // Alternar selección de un juego
  const toggleSeleccionado = (id: string) => {
    setSeleccionados((prev) =>
      prev.includes(id) ? prev.filter((f) => f !== id) : [...prev, id]
    );
  };

  const enviarSeleccionados = () => {
    // Aquí envías la lista al back
    console.log("Juegos seleccionados:", seleccionados);
    alert("Juegos enviados al backend: " + seleccionados.join(", "));
  };

  // Todo: DEFINIR LOGICA 
  const juegosPopulares = juegos.slice(0, 4);

  return (
    <div style={{ padding: "1rem", fontFamily: "sans-serif" }}>
      <h1>Busqueda de Juegos</h1>

      <input
        type="text"
        placeholder="Buscar un juego..."
        value={query}
        onChange={(e) => setQuery(e.target.value)}
        style={{ width: "100%", padding: "0.5rem", marginBottom: "1rem" }}
      />

      {loading && <p>Cargando juegos...</p>}
      {error && <p>{error}</p>}

      {/* Juegos Populares solo si no hay búsqueda */}
      {query === "" && (
        <>
          <h2>Juegos Populares</h2>
          <div style={{ display: "flex", gap: "1rem", flexWrap: "wrap" }}>
            {juegosPopulares.map((j) => {
              const seleccionado = seleccionados.includes(j.id);
              return (
                <div
                  key={j.id}
                  onClick={() => toggleSeleccionado(j.id)}
                  style={{
                    border: seleccionado ? "2px solid blue" : "1px solid #ccc",
                    padding: "0.5rem",
                    cursor: "pointer",
                    textAlign: "center",
                    width: 150,
                  }}
                >
                  <img src={j.imageURL} alt={j.nombre} width={150} />
                  <p>{j.nombre}</p>
                </div>
              );
            })}
          </div>
        </>
      )}

      {/* Todos los juegos */}
      <h2>Todos los Juegos</h2>
      <div
        style={{
          maxHeight: "400px",
          overflowY: "auto",
          display: "flex",
          flexWrap: "wrap",
          gap: "1rem",
        }}
      >
        {juegosFiltrados.map((j) => {
          const seleccionado = seleccionados.includes(j.id);
          return (
            <div
              key={j.id}
              onClick={() => toggleSeleccionado(j.id)}
              style={{
                border: seleccionado ? "2px solid green" : "1px solid #ccc",
                padding: "0.5rem",
                cursor: "pointer",
                textAlign: "center",
                width: 150,
              }}
            >
              <img src={j.imageURL} alt={j.nombre} width={150} />
              <p>{j.nombre}</p>
            </div>
          );
        })}
      </div>

      {/* Mensaje si no hay resultados */}
      {!loading && juegosFiltrados.length === 0 && <p>No se encontraron juegos.</p>}

      {/* Botón flotante */}
      {seleccionados.length > 0 && (
        <button
          onClick={enviarSeleccionados}
          style={{
            position: "fixed",
            bottom: 20,
            right: 20,
            backgroundColor: "blue",
            color: "white",
            padding: "1rem 2rem",
            borderRadius: "50px",
            border: "none",
            cursor: "pointer",
          }}
        >
          Enviar {seleccionados.length}
        </button>
      )}
    </div>
  );
};




