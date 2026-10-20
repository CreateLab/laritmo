import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import path from 'path'

export default defineConfig({
    plugins: [vue()],
    resolve: {
        alias: {
            '@': path.resolve(__dirname, './src'),
        },
    },
    server: {
        proxy: {
            '/api': {
                target: process.env.API_PROXY_TARGET || 'https://localhost:8443',
                changeOrigin: true,
                secure: false,
            },
        },
    },
    build: {
        outDir: process.env.DOCKER_BUILD ? './dist' : '../back/web',
        emptyOutDir: true,
    },
})