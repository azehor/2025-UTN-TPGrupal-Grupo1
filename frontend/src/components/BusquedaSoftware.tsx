import React, { useState, useEffect, useRef } from "react";
import { Link } from "react-router-dom";

interface Software {
  id: string;
  nombre: string;
  empresa: string;
  imageURL: string;
}

export const BusquedaSoftware: React.FC = () => {
  const [q, setQ] = useState<string>("");
  const [softwares, setSoftwares] = useState<Software[]>([]);
  const [filtered, setFiltered] = useState<Software[]>([]);
  const [highlight, setHighlight] = useState<number>(-1);
  const [stack, setStack] = useState<Software[]>([]);
  const [error, setError] = useState<string | null>(null);
  const inputRef = useRef<HTMLInputElement>(null);

  // Fetch API
  useEffect(() => {
    const fetchSoftwares = async () => {
      try {
        const res = await fetch("http://localhost:8080/v1/softwares");
        if (!res.ok) throw new Error(`Error HTTP: ${res.status}`);
        const data: Software[] = await res.json();
        setSoftwares(data);
      } catch (err: any) {
        setError(err.message || "Error al obtener listado de software");
        console.error(err);
      }
    };
    fetchSoftwares();
  }, []);

  // Filtrado en vivo
  useEffect(() => {
    if (!q) {
      setFiltered([]);
      setHighlight(-1);
      return;
    }
    const ql = q.trim().toLowerCase();
    const f = softwares
      .filter(s => s.nombre.toLowerCase().includes(ql) && !stack.some(st => st.id === s.id))
      .slice(0, 8);
    setFiltered(f);
    setHighlight(f.length ? 0 : -1);
  }, [q, softwares, stack]);

  // Eventos de teclado
  const onKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === "ArrowDown") {
      e.preventDefault();
      setHighlight(h => Math.min(h + 1, filtered.length - 1));
    } else if (e.key === "ArrowUp") {
      e.preventDefault();
      setHighlight(h => Math.max(h - 1, 0));
    } else if (e.key === "Enter") {
      e.preventDefault();
      const pick = filtered[highlight] ?? softwares.find(s => s.nombre.toLowerCase() === q.trim().toLowerCase());
      if (pick) addToStack(pick);
    } else if (e.key === "Escape") {
      setFiltered([]);
      setHighlight(-1);
    }
  };

  const addToStack = (soft: Software) => {
    if (stack.some(s => s.id === soft.id)) {
      resetInput();
      return;
    }
    setStack(s => [soft, ...s]);
    resetInput();
  };

  const removeFromStack = (id: string) => setStack(s => s.filter(x => x.id !== id));

  const resetInput = () => {
    setQ("");
    setFiltered([]);
    setHighlight(-1);
    inputRef.current?.focus();
  };

  return (
    <div style={styles.container}>
      <h2>Buscar por Software</h2>

      <div style={styles.inputWrap}>
        <input
          ref={inputRef}
          type="text"
          value={q}
          placeholder="Buscar..."
          onChange={(e) => setQ(e.target.value)}
          onKeyDown={onKeyDown}
          style={styles.input}
        />
        {filtered.length > 0 && (
          <ul style={styles.suggList}>
            {filtered.map((f, i) => (
              <li
                key={f.id}
                onMouseDown={(ev) => { ev.preventDefault(); addToStack(f); }}
                onMouseEnter={() => setHighlight(i)}
                style={{
                  ...styles.suggItem,
                  ...(i === highlight ? styles.suggItemActive : {})
                }}
              >
                {f.nombre} – {f.empresa}
              </li>
            ))}
          </ul>
        )}
      </div>

      {error && <p style={styles.error}>{error}</p>}

      <ul style={styles.stackList}>
        {stack.map(item => (
          <li key={item.id} style={styles.stackItem}>
            <span>{item.nombre} – {item.empresa}</span>
            <button
              onClick={() => removeFromStack(item.id)}
              style={styles.removeBtn}
            >
              ✕
            </button>
          </li>
        ))}
      </ul>

      {/* Botón al final */}
      <Link to="/recomendacion">
      <button
        disabled={stack.length === 0}
        style={{
          ...styles.actionBtn,
          ...(stack.length === 0 ? styles.actionBtnDisabled : {})
        }}
      >
        Recomendar PC
      </button>
      </Link>
    </div>
  );
};

// Estilos básicos
const styles: Record<string, React.CSSProperties> = {
  container: { width: "100%", maxWidth: 720, fontFamily: "Inter, sans-serif", padding: "1rem", margin: "0 auto" },
  inputWrap: { position: "relative" },
  input: { width: "100%", padding: "10px 12px", fontSize: 14, boxSizing: "border-box" },
  suggList: {
    position: "absolute",
    left: 0,
    right: 0,
    marginTop: 4,
    background: "#fff",
    border: "1px solid #ddd",
    listStyle: "none",
    padding: 0,
    zIndex: 50
  },
  suggItem: { padding: "8px 10px", cursor: "pointer" },
  suggItemActive: { background: "#efefef" },
  stackList: {
    marginTop: 16,
    padding: 0,
    listStyle: "none",
    display: "flex",
    flexDirection: "column",
    gap: 8
  },
  stackItem: {
    display: "flex",
    justifyContent: "space-between",
    alignItems: "center",
    padding: "8px 12px",
    border: "1px solid #ddd",
    borderRadius: 6,
    background: "#fafafa"
  },
  removeBtn: {
    marginLeft: 8,
    border: "none",
    background: "transparent",
    cursor: "pointer",
    fontSize: 16,
    lineHeight: 1
  },
  error: { color: "red", fontSize: 14, marginTop: 8 },
  actionBtn: {
    marginTop: 20,
    width: "100%",
    padding: "12px",
    borderRadius: 6,
    border: "none",
    cursor: "pointer",
    background: "#007bff",
    color: "#fff",
    fontSize: 16
  },
  actionBtnDisabled: {
    background: "#ccc",
    cursor: "not-allowed"
  }
};