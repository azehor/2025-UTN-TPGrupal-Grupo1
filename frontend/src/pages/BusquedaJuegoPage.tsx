import type { FC } from "react"
import { BusquedaJuego } from "../components/BusquedaJuego"
import { BusquedaJuegoSinEstilos } from "../components/BusquedaJuegoSinEstilosParaTest"

export const BusquedaJuegoPage: FC = () => {
    return <>
         <BusquedaJuego />
        { /*<BusquedaJuegoSinEstilos / >*/}
    </>
}