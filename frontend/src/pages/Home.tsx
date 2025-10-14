import type { FC } from "react";
import Boton from "../components/Boton";
import "./Home.css";
import { Link } from 'react-router-dom'

const Home: FC = () => {
  return (
    <div className="container">
      <Link to="/panel" className="admin-top-button">
        <Boton>Panel</Boton>
      </Link>
      <h1 className="titulo">QuePC</h1>
      <h2 className="subtitulo">Seleccione debajo la forma por la cual buscará su PC</h2>

      <div className="modulos">
        <Link to="/busqueda-software">
          <Boton>Buscar por Software</Boton>
        </Link>
        <Link to="/busqueda-carrera">
        <Boton>Buscar por Carrera</Boton>
        </Link>
      </div>
    </div>
  )
} 

export default Home
