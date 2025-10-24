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
    <div className="min-h-screen bg-[#101c22] text-gray-300 font-['Space_Grotesk',sans-serif] py-8">
      <div className="container mx-auto px-6 max-w-4xl">
        <h2 className="text-3xl md:text-4xl font-bold text-white mb-8 text-center">Buscar por Software</h2>

        <div className="relative mb-6">
          <input
            ref={inputRef}
            type="text"
            value={q}
            placeholder="Buscar software..."
            onChange={(e) => setQ(e.target.value)}
            onKeyDown={onKeyDown}
            className="w-full px-4 py-3 bg-[#1a2831] border border-gray-600 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:border-[#13a4ec] focus:ring-2 focus:ring-[#13a4ec]/20 transition-all"
          />
          {filtered.length > 0 && (
            <ul className="absolute left-0 right-0 mt-2 bg-[#1a2831] border border-gray-600 rounded-lg shadow-lg z-50 max-h-64 overflow-y-auto">
              {filtered.map((f, i) => (
                <li
                  key={f.id}
                  onMouseDown={(ev) => { ev.preventDefault(); addToStack(f); }}
                  onMouseEnter={() => setHighlight(i)}
                  className={`px-4 py-3 cursor-pointer transition-colors border-b border-gray-700 last:border-b-0 ${
                    i === highlight ? 'bg-[#13a4ec] text-white' : 'hover:bg-gray-700'
                  }`}
                >
                  <span className="font-medium">{f.nombre}</span>
                  <span className="text-gray-400 ml-2">– {f.empresa}</span>
                </li>
              ))}
            </ul>
          )}
        </div>

        {error && (
          <div className="bg-red-900/20 border border-red-500/50 text-red-300 px-4 py-3 rounded-lg mb-6">
            {error}
          </div>
        )}

        {stack.length > 0 && (
          <div className="mb-8">
            <h3 className="text-xl font-semibold text-white mb-4">Software Seleccionado:</h3>
            <div className="space-y-3">
              {stack.map(item => (
                <div key={item.id} className="flex items-center justify-between bg-[#1a2831] border border-gray-600 rounded-lg px-4 py-3">
                  <div>
                    <span className="font-medium text-white">{item.nombre}</span>
                    <span className="text-gray-400 ml-2">– {item.empresa}</span>
                  </div>
                  <button
                    onClick={() => removeFromStack(item.id)}
                    className="ml-4 text-gray-400 hover:text-red-400 transition-colors text-xl leading-none"
                  >
                    ✕
                  </button>
                </div>
              ))}
            </div>
          </div>
        )}

        <Link to="/recomendacion">
          <button
            disabled={stack.length === 0}
            className={`w-full py-3 px-8 rounded-lg font-bold text-lg transition-all ${
              stack.length === 0
                ? 'bg-gray-600 text-gray-400 cursor-not-allowed'
                : 'bg-[#13a4ec] text-white hover:opacity-90'
            }`}
          >
            Recomendar PC
          </button>
        </Link>
      </div>
    </div>
  );
};

