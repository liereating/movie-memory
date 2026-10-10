import react from '@vitejs/plugin-react'
import { defineConfig } from 'vite'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    proxy: {
      // 开发时把 /api 请求转发到后端
      '/api': 'http://localhost:8080',
    },
  },
})
