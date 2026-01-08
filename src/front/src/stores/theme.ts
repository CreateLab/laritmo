import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

type Theme = 'light' | 'dark' | 'system'

export const useThemeStore = defineStore('theme', () => {
    const currentTheme = ref<Theme>(
        (localStorage.getItem('laritmo-theme') as Theme) || 'system'
    )

    const effectiveTheme = computed<'light' | 'dark'>(() => {
        if (currentTheme.value === 'system') {
            return window.matchMedia('(prefers-color-scheme: dark)').matches
                ? 'dark'
                : 'light'
        }
        return currentTheme.value
    })

    const applyTheme = (theme: 'light' | 'dark') => {
        console.log('Applying theme:', theme) // Для отладки
        const root = document.documentElement
        if (theme === 'dark') {
            root.classList.add('dark')
        } else {
            root.classList.remove('dark')
        }
    }

    const setTheme = (theme: Theme) => {
        currentTheme.value = theme
        localStorage.setItem('laritmo-theme', theme)
        
        // Явно вычисляем эффективную тему для надежности
        let effective: 'light' | 'dark'
        if (theme === 'system') {
            effective = window.matchMedia('(prefers-color-scheme: dark)').matches
                ? 'dark'
                : 'light'
        } else {
            effective = theme
        }
        
        applyTheme(effective)
    }

    const toggleTheme = () => {
        // Переключаем на основе текущей эффективной темы
        const currentEffective = effectiveTheme.value
        const newTheme = currentEffective === 'dark' ? 'light' : 'dark'
        setTheme(newTheme)
    }

    const initTheme = () => {
        // Применяем текущую тему
        // Явно вычисляем эффективную тему для надежности
        let effective: 'light' | 'dark'
        if (currentTheme.value === 'system') {
            effective = window.matchMedia('(prefers-color-scheme: dark)').matches
                ? 'dark'
                : 'light'
        } else {
            effective = currentTheme.value
        }
        applyTheme(effective)

        // Слушаем изменения системной темы, если выбрана 'system'
        if (currentTheme.value === 'system') {
            const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)')
            const handleChange = (e: MediaQueryListEvent) => {
                if (currentTheme.value === 'system') {
                    applyTheme(e.matches ? 'dark' : 'light')
                }
            }
            
            // Современный способ
            if (mediaQuery.addEventListener) {
                mediaQuery.addEventListener('change', handleChange)
            } else {
                // Fallback для старых браузеров
                mediaQuery.addListener(handleChange)
            }
        }
    }

    return {
        currentTheme,
        effectiveTheme,
        setTheme,
        toggleTheme,
        initTheme,
    }
})
