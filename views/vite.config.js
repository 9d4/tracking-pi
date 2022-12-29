import { resolve } from 'path'
import { defineConfig } from 'vite'

export default defineConfig({
  build: {
    rollupOptions: {
      input: {
        main: resolve(__dirname, 'index.html'),
        admin: resolve(__dirname, 'admin/index.html'),
        admin_industry: resolve(__dirname, 'admin/industry/index.html'),
        admin_volunteer: resolve(__dirname, 'admin/volunteer/index.html'),
      },
    },
  },
  resolve: {
    alias: {
      '~bootstrap': resolve(__dirname, 'node_modules/bootstrap'),
    }
  }
})
