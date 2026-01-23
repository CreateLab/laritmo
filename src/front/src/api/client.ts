import axios from 'axios'

const apiClient = axios.create({
    baseURL: import.meta.env.VITE_API_BASE_URL || '/api',
    timeout: 10000,
    withCredentials: true,
    headers: {
        'Content-Type': 'application/json',
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
    (response) => response,
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
                // Redirect to home with reload to clear state
                window.location.href = '/'
            }
        }

        return Promise.reject(error)
    }
)

export default apiClient