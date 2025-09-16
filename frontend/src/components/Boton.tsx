import type { FC, ReactNode } from "react";
import "./Boton.css";

interface BotonProps {
  children: ReactNode;
}

const Boton: FC<BotonProps> = ({ children }) => {
  return <button className="botonAzul">{children}</button>;
};

export default Boton;