import type { ReactNode } from 'react'
import React, { createContext, useState, useContext } from "react";

type TipoBusqueda = "carrera" | "software" | "juego" | null;

interface BusquedaContextType {
  tipo: TipoBusqueda;
  datos: any; // object para carrera, array para softwares/juegos
  setBusqueda: (tipo: TipoBusqueda, datos: any) => void;
}

const BusquedaContext = createContext<BusquedaContextType | undefined>(undefined);

export const BusquedaProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
  const [tipo, setTipo] = useState<TipoBusqueda>(null);
  const [datos, setDatos] = useState<any>(null);

  const setBusqueda = (t: TipoBusqueda, d: any) => {
    setTipo(t);
    setDatos(d);
  };

  return (
    <BusquedaContext.Provider value={{ tipo, datos, setBusqueda }}>
      {children}
    </BusquedaContext.Provider>
  );
};

export const useBusqueda = () => {
  const context = useContext(BusquedaContext);
  if (!context) throw new Error("useBusqueda debe usarse dentro de BusquedaProvider");
  return context;
};
