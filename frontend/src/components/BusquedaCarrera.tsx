import React, { useState } from "react";

const listaCarrera: string[] = [
  "Arquitectura",
  "Diseño Gráfico",
  "Fotografía",
  "ingeniería de Software",
  "Arquitectura de Software",
  "Desarrollo Web",
  "profesorado de matemática",
  "ingenieria civil",
  "ingenieria industrial",
  "fisica",
  "marketing",
  "Traducción",
  "Humanidades",
  "ingenieria quimica",
  "medicina",
  "Leyes",
];

export const BusquedaCarrera: React.FC = () => {
  const [query, setQuery] = useState("");

  const filteredList = listaCarrera.filter(soft =>
    soft.toLowerCase().includes(query.toLowerCase())
  );

  return (
    <div style={{ fontFamily: "sans-serif", padding: "1rem", maxWidth: "400px" }}>
      <h2>Busqueda por carrera</h2>

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
          filteredList.map((soft, index) => <li key={index}>{soft}</li>)
        ) : (
          <li>No software found</li>
        )}
      </ul>
    </div>
  );
};