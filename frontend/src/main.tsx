import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import { BrowserRouter } from 'react-router-dom';
import './index.css';
import './App.css'; 
import App from './App';
import { BusquedaProvider } from './context/BusquedaContext';
import { ThemeProvider } from './context/ThemeContext';

async function loadRuntimeConfig() {
  try {
    const res = await fetch('/config.json', { cache: 'no-store' });
    if (!res.ok) return;
    const cfg = await res.json();
    (window as any).__ENV = cfg;
  } catch {
    // no runtime config available
  }
}

(async () => {
  await loadRuntimeConfig();

  createRoot(document.getElementById('root')!).render(
    <StrictMode>
      <ThemeProvider>
        <BusquedaProvider>
          <BrowserRouter>
            <App />
          </BrowserRouter>
        </BusquedaProvider>
      </ThemeProvider>
    </StrictMode>
  );
})();