import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { BrowserRouter } from 'react-router-dom'
import './index.css'
// import './styles/theme.css'
import App from './App.tsx'
import { BusquedaProvider } from "./context/BusquedaContext"

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    < BusquedaProvider>
      <BrowserRouter>
        <App />
      </BrowserRouter>
    </BusquedaProvider>
  </StrictMode>,
)
