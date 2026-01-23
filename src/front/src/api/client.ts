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
        // Добавляем JWT токен из localStorage, если он есть
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

        // При 401 (Unauthorized) — автоматический логаут
        if (error.response?.status === 401) {
            const token = localStorage.getItem('token')
            // Логаут только если был токен (т.е. пользователь был залогинен)
            if (token) {
                console.warn('Token expired or invalid, logging out...')
                localStorage.removeItem('token')
                localStorage.removeItem('user')
                // Редирект на главную с перезагрузкой для очистки состояния
                window.location.href = '/'
            }
        }

        return Promise.reject(error)
    }
)

export default apiClient