import type { FC } from "react";
import Boton from "../components/boton";
import "../App.css"; //para los estilos globales

const Home: FC = () => {
  return (
    <div className="container">
      <h1 className="titulo">QuePC</h1>

      <div className="modulos">
        <Boton>Buscar por Software</Boton>
        <Boton>Buscar por Carrera</Boton>
      </div>
    </div>
  )
}

export default Home
