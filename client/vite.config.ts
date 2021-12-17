import { defineConfig } from 'vite'
import path from 'path'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  resolve: {
    alias: {
      '/@': path.resolve(__dirname, 'src')
    }
  },
  server: {
    port: 8080,
    proxy: {
      '/api': {
        target: 'https://localhost:3000/api',
        changeOrigin: true
      }
    }
  },
  plugins: [vue()]
})
