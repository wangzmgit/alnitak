import { fileURLToPath, URL } from 'node:url';
import { defineConfig } from 'vite';
import vue from '@vitejs/plugin-vue';
import { globalConfig } from "./src/utils/global-config";


export default defineConfig({
  plugins: [
    vue(),
  ],
  base: `/${globalConfig.baseUrl}`,
  build: {
    chunkSizeWarningLimit: 2048,
    outDir: 'admin', //指定输出路径
  },
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    }
  }
})
