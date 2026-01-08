<template>
  <div class="min-h-screen flex items-center justify-center p-6 relative">
    <!-- Кнопка переключения темы в правом верхнем углу -->
    <div class="absolute top-4 right-4">
      <ThemeToggle />
    </div>
    
    <div class="w-full max-w-md">
      <div class="text-center mb-8">
        <h1 class="text-4xl font-bold text-forest-dark dark:text-dark-text flex items-center justify-center gap-2 mb-2 transition-colors duration-300">
          <span class="frog-animation">🐸</span> Laritmo
        </h1>
        <p class="text-sm text-gray-600 dark:text-dark-text-secondary transition-colors duration-300">Forest Academy</p>
      </div>

      <div class="bg-white dark:bg-dark-surface rounded-xl shadow-lg p-8 transition-colors duration-300">
        <h2 class="text-2xl font-semibold mb-6 text-center text-forest-dark dark:text-dark-text transition-colors duration-300">
          Вход в систему
        </h2>

        <form @submit.prevent="handleLogin" class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-dark-text-secondary mb-1 transition-colors duration-300">Username</label>
            <input
                v-model="loginForm.username"
                type="text"
                required
                autofocus
                class="w-full px-3 py-2 border border-gray-300 dark:border-dark-border rounded-lg focus:outline-none focus:ring-2 focus:ring-forest-green dark:focus:ring-forest-green-dark bg-white dark:bg-dark-bg text-gray-900 dark:text-dark-text transition-colors duration-300"
            />
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-dark-text-secondary mb-1 transition-colors duration-300">Password</label>
            <input
                v-model="loginForm.password"
                type="password"
                required
                class="w-full px-3 py-2 border border-gray-300 dark:border-dark-border rounded-lg focus:outline-none focus:ring-2 focus:ring-forest-green dark:focus:ring-forest-green-dark bg-white dark:bg-dark-bg text-gray-900 dark:text-dark-text transition-colors duration-300"
            />
          </div>

          <div v-if="loginError" class="text-red-600 dark:text-red-400 text-sm bg-red-50 dark:bg-red-900/30 p-3 rounded-lg transition-colors duration-300">
            {{ loginError }}
          </div>

          <button
              type="submit"
              :disabled="loading"
              class="w-full px-4 py-2 bg-forest-green dark:bg-forest-green-dark text-white hover:bg-forest-dark dark:hover:bg-forest-green rounded-lg disabled:opacity-50 transition-colors duration-300 font-medium"
          >
            {{ loading ? 'Вход...' : 'Войти' }}
          </button>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useRouter } from 'vue-router'
import ThemeToggle from '@/components/ThemeToggle.vue'

const authStore = useAuthStore()
const router = useRouter()

const loginForm = ref({
  username: '',
  password: '',
})
const loginError = ref('')
const loading = ref(false)

const handleLogin = async () => {
  loginError.value = ''
  loading.value = true

  try {
    await authStore.login(loginForm.value.username, loginForm.value.password)
    router.push('/')
  } catch (error: any) {
    loginError.value = error.response?.data?.error || 'Ошибка входа'
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  if (authStore.isAuthenticated) {
    router.push('/')
  }
})
</script>
