import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// In dev, calls to /api are proxied to the Go backend so the browser makes
// same-origin requests (no CORS dance). In prod the API base is set via
// VITE_API_BASE. See src/api.ts.
export default defineConfig({
  plugins: [react()],
  server: {
    port: 5173,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
        rewrite: (p) => p.replace(/^\/api/, ''),
      },
    },
  },
})
