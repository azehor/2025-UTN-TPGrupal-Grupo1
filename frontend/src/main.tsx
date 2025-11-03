import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { BrowserRouter } from 'react-router-dom'
import './index.css'
// import './styles/theme.css'
import App from './App.tsx'
import { BusquedaProvider } from "./context/BusquedaContext"

async function loadRuntimeConfig() {
  try {
    const res = await fetch('/config.json', { cache: 'no-store' });
    if (!res.ok) return;
    const cfg = await res.json();
    (window as any).__ENV = cfg;
  } catch (e) {
    // no runtime config available
  }
}

(async () => {
  await loadRuntimeConfig();
  createRoot(document.getElementById('root')!).render(
    <StrictMode>
      <BusquedaProvider>
        <BrowserRouter>
          <App />
        </BrowserRouter>
      </BusquedaProvider>
    </StrictMode>,
  )
})()
