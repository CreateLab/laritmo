import '@fontsource/roboto/cyrillic-400.css'
import '@fontsource/roboto/cyrillic-500.css'
import '@fontsource/roboto/cyrillic-700.css'
import '@fontsource/roboto/latin-400.css'
import '@fontsource/roboto/latin-500.css'
import '@fontsource/roboto/latin-700.css'
import '@fontsource/roboto-mono/cyrillic-400.css'
import '@fontsource/roboto-mono/latin-400.css'
import './assets/main.css'

import { createApp } from 'vue'
import { createPinia } from 'pinia'

import App from './App.vue'
import router from './router'

import PrimeVue from 'primevue/config'
import Aura from '@primevue/themes/aura'
import { definePreset } from '@primevue/themes'
import ToastService from 'primevue/toastservice'
import 'primeicons/primeicons.css'

import { useAuthStore } from './stores/auth'

const Laritmo = definePreset(Aura, {
    semantic: {
        primary: {
            50: '#eef6f0', 100: '#cde3d3', 200: '#a9cfb3', 300: '#7fb78e',
            400: '#58a06c', 500: '#3d8553', 600: '#2a6b3f', 700: '#1f5230',
            800: '#173f25', 900: '#102c19', 950: '#0a1c10'
        },
        colorScheme: {
            light: {
                primary: { color: '#2a6b3f', contrastColor: '#ffffff', hoverColor: '#1f5230', activeColor: '#173f25' },
                surface: { 0: '#ffffff', 50: '#f5f6f4', 100: '#ebefea', 200: '#e0e4de', 300: '#c9cfc8', 400: '#7b857d', 500: '#55605a', 600: '#3c4540', 700: '#2a312d', 800: '#1d221f', 900: '#151a16', 950: '#0c0f0d' }
            },
            dark: {
                primary: { color: '#8fd19e', contrastColor: '#0c1a10', hoverColor: '#aee0b9', activeColor: '#c7ebcf' },
                surface: { 0: '#ecefec', 50: '#9aa59d', 100: '#6b766e', 200: '#2a2f2c', 300: '#212522', 400: '#181b19', 500: '#181b19', 600: '#181b19', 700: '#181b19', 800: '#181b19', 900: '#101211', 950: '#0a0c0b' }
            }
        }
    }
})

const app = createApp(App)
const pinia = createPinia()

app.use(pinia)
app.use(router)
app.use(PrimeVue, {
    theme: {
        preset: Laritmo,
        options: {
            darkModeSelector: '.dark',
        }
    }
})
app.use(ToastService)

const authStore = useAuthStore()
authStore.initAuth()

app.mount('#app')