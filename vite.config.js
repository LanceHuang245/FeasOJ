import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import vuetify from 'vite-plugin-vuetify'
import compression from 'vite-plugin-compression';

export default defineConfig({
  base: './',
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
