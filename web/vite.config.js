import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import vuetify from 'vite-plugin-vuetify'
import compression from 'vite-plugin-compression';

export default defineConfig({
  base: './',
  server: {
    cors: true,
    proxy: {
      '/api': {
        target: 'http://127.0.0.1:37882',
        changeOrigin: true,
      },
    },
  },
  plugins: [
    compression({
      algorithm: 'brotliCompress',
      threshold: 10240
    }),
    vue(),
    vuetify({
      autoImport: true,
    }),
  ]
})
