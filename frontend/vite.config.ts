import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    react(), 
    tailwindcss(),
  ],
  server: {
    // Bind to all interfaces inside the container
    host: true,
    port: 5173,
    strictPort: true,
    // Ensure HMR websocket points to host:5173 when accessed from the browser on the host
    hmr: {
      host: 'localhost',
      clientPort: 5173,
      protocol: 'ws',
    },
    // On Windows + Docker bind mounts, file system events often don't propagate; use polling
    watch: {
      usePolling: true,
      interval: 100,
    },
  },
})
