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
        console.log('Applying theme:', theme) // For debugging
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
        
        // Explicitly compute effective theme for reliability
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
        // Toggle based on current effective theme
        const currentEffective = effectiveTheme.value
        const newTheme = currentEffective === 'dark' ? 'light' : 'dark'
        setTheme(newTheme)
    }

    const initTheme = () => {
        // Apply current theme
        // Explicitly compute effective theme for reliability
        let effective: 'light' | 'dark'
        if (currentTheme.value === 'system') {
            effective = window.matchMedia('(prefers-color-scheme: dark)').matches
                ? 'dark'
                : 'light'
        } else {
            effective = currentTheme.value
        }
        applyTheme(effective)

        // Listen for system theme changes if 'system' is selected
        if (currentTheme.value === 'system') {
            const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)')
            const handleChange = (e: MediaQueryListEvent) => {
                if (currentTheme.value === 'system') {
                    applyTheme(e.matches ? 'dark' : 'light')
                }
            }
            
            // Modern approach
            if (mediaQuery.addEventListener) {
                mediaQuery.addEventListener('change', handleChange)
            } else {
                // Fallback for older browsers
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
