import type { FC } from "react";
import "./RecomendacionesPage.css";

export const RecomendacionPage: FC = () => {
    //Pc hardcodeada!
    const pc = {
        nombre: "PC Gamer Básica",
        procesador: "Intel Core i5 10400F",
        ram: "16GB DDR4",
        almacenamiento: "512GB SSD",
        grafica: "NVIDIA GTX 1660 Super",
        precio: "$2500000"
    };

    return (
        <>

            <div className="resultados-container">
                <h1>Resultados del Recomendador</h1>

                <div className="pc-card">
                    <h2>{pc.nombre}</h2>
                    <ul>
                        <li><strong>Procesador:</strong> {pc.procesador}</li>
                        <li><strong>Memoria RAM:</strong> {pc.ram}</li>
                        <li><strong>Almacenamiento:</strong> {pc.almacenamiento}</li>
                        <li><strong>Gráfica:</strong> {pc.grafica}</li>
                        <li><strong>Precio:</strong> {pc.precio}</li>
                    </ul>
                </div>
            </div>
        </>

    );
};