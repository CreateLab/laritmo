import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import axios from '@/api/client'

interface User {
    id: number
    username: string
    email: string
    role: string
}

export const useAuthStore = defineStore('auth', () => {
    const token = ref<string | null>(localStorage.getItem('token'))
    const user = ref<User | null>(null)

    const storedUser = localStorage.getItem('user')
    if (storedUser && storedUser !== 'undefined' && storedUser !== 'null') {
        try {
            const parsed = JSON.parse(storedUser)
            if (parsed && typeof parsed === 'object' && parsed.id) {
                user.value = parsed
            } else {
                localStorage.removeItem('user')
            }
        } catch (e) {
            console.error('Failed to parse stored user:', e)
            localStorage.removeItem('user')
        }
    } else {
        localStorage.removeItem('user')
    }

    const isAuthenticated = computed(() => !!token.value && !!user.value)
    const isAdmin = computed(() => user.value?.role === 'admin')

    const login = async (username: string, password: string) => {
        const { data } = await axios.post('/auth/login', { username, password })

        token.value = data.token
        user.value = data.user
        localStorage.setItem('token', data.token)
        localStorage.setItem('user', JSON.stringify(data.user))

        axios.defaults.headers.common['Authorization'] = `Bearer ${data.token}`
    }

    const logout = () => {
        token.value = null
        user.value = null
        localStorage.removeItem('token')
        localStorage.removeItem('user')
        delete axios.defaults.headers.common['Authorization']
    }

    const initAuth = () => {
        if (token.value) {
            axios.defaults.headers.common['Authorization'] = `Bearer ${token.value}`
        }
    }

    return {
        token,
        user,
        isAuthenticated,
        isAdmin,
        login,
        logout,
        initAuth,
    }
})