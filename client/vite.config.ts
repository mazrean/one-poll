import { defineConfig } from 'vite'
import path from 'path'
import vue from '@vitejs/plugin-vue'
import brotli from 'rollup-plugin-brotli'

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
        target: 'http://localhost:3000',
        changeOrigin: true
      }
    }
  },
  plugins: [vue()],
  build: {
    rollupOptions: {
      plugins: [
        brotli({
          test: /\.(js|css|html|svg)$/
        })
      ]
    }
  }
})
