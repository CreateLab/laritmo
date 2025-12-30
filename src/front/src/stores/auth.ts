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

    const isAuthenticated = computed(() => !!token.value)
    const isAdmin = computed(() => user.value?.role === 'admin')

    const login = async (username: string, password: string) => {
        const { data } = await axios.post('/auth/login', { username, password })

        token.value = data.token
        user.value = data.user
        localStorage.setItem('token', data.token)

        axios.defaults.headers.common['Authorization'] = `Bearer ${data.token}`
    }

    const logout = () => {
        token.value = null
        user.value = null
        localStorage.removeItem('token')
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