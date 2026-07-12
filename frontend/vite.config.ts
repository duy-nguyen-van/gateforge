import path from 'node:path'
import tailwindcss from '@tailwindcss/vite'
import react from '@vitejs/plugin-react'
import { defineConfig } from 'vite'

const apiTarget = process.env.VITE_API_PROXY_TARGET ?? 'http://localhost:3000'

export default defineConfig({
  plugins: [react(), tailwindcss()],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
  },
  server: {
    port: 5173,
    proxy: {
      '/api/v1': { target: apiTarget, changeOrigin: true },
      '/oidc': { target: apiTarget, changeOrigin: true },
      '/authorize': { target: apiTarget, changeOrigin: true },
      '/token': { target: apiTarget, changeOrigin: true },
      '/userinfo': { target: apiTarget, changeOrigin: true },
      '/.well-known': { target: apiTarget, changeOrigin: true },
    },
  },
})
