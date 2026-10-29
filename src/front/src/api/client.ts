import axios from 'axios'
import { cacheKey, clearCache, getCached, isCacheable, setCached } from './cache'

const defaultAdapter = axios.getAdapter(axios.defaults.adapter)

const apiClient = axios.create({
    baseURL: import.meta.env.VITE_API_BASE_URL || '/api',
    timeout: 10000,
    withCredentials: true,
    headers: {
        'Content-Type': 'application/json',
    },
    // Session cache for GET /courses, /lectures, /labs, ... (cleared on admin writes)
    adapter: async (config) => {
        if (!isCacheable(config)) return defaultAdapter(config)

        const key = cacheKey(config)
        const hit = getCached(key)
        if (hit) {
            return { data: structuredClone(hit.data), status: hit.status, statusText: 'OK', headers: {}, config }
        }

        const response = await defaultAdapter(config)
        setCached(key, response)
        return response
    },
})

apiClient.interceptors.request.use(
    (config) => {
        // Add JWT token from localStorage if available
        const token = localStorage.getItem('token')
        if (token) {
            config.headers.Authorization = `Bearer ${token}`
        }
        console.log('API Request:', config.method?.toUpperCase(), config.url)
        return config
    },
    (error) => {
        return Promise.reject(error)
    }
)

apiClient.interceptors.response.use(
    (response) => {
        const method = response.config.method?.toLowerCase()
        if (method && method !== 'get' && response.config.url?.startsWith('/admin/')) {
            clearCache()
        }
        return response
    },
    (error) => {
        console.error('API Error:', error.response?.status, error.message)

        // On 401 (Unauthorized) - automatic logout
        if (error.response?.status === 401) {
            const token = localStorage.getItem('token')
            // Logout only if token existed (i.e., user was logged in)
            if (token) {
                console.warn('Token expired or invalid, logging out...')
                localStorage.removeItem('token')
                localStorage.removeItem('user')
                clearCache()
                // Redirect to home with reload to clear state
                window.location.href = '/'
            }
        }

        return Promise.reject(error)
    }
)

export default apiClient