import {defineConfig} from 'vite'
import react from '@vitejs/plugin-react'

const srcPath = new URL('./src', import.meta.url).pathname
// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  resolve: {
      alias: {
          "@": srcPath,
      }
  }
})
