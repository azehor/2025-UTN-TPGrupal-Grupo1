import type { FC } from "react";
import Boton from "../components/Boton";
import "../App.css"; //para los estilos globales
import { Link } from 'react-router-dom'

const Home: FC = () => {
  return (
    <div className="container">
      <h1 className="titulo">QuePC</h1>
      <h2 className="subtitulo">Seleccione debajo la forma por la cual buscará su PC</h2>

      <div className="modulos">
        <Link to="/software-search">
          <Boton>Busqueda por Software</Boton>
        </Link>
        <Boton>Buscar por Carrera</Boton>
      </div>
    </div>
  )
} 

export default Home
