<template>
  <div class="min-h-screen flex items-center justify-center p-6">
    <div class="w-full max-w-md">
      <div class="text-center mb-8">
        <h1 class="text-4xl font-bold text-forest-dark flex items-center justify-center gap-2 mb-2">
          <span class="frog-animation">🐸</span> Laritmo
        </h1>
        <p class="text-sm text-gray-600">Forest Academy</p>
      </div>

      <div class="bg-white rounded-xl shadow-lg p-8">
        <h2 class="text-2xl font-semibold mb-6 text-center text-forest-dark">
          Вход в систему
        </h2>

        <form @submit.prevent="handleLogin" class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Username</label>
            <input
                v-model="loginForm.username"
                type="text"
                required
                autofocus
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-forest-green"
            />
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Password</label>
            <input
                v-model="loginForm.password"
                type="password"
                required
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-forest-green"
            />
          </div>

          <div v-if="loginError" class="text-red-600 text-sm bg-red-50 p-3 rounded-lg">
            {{ loginError }}
          </div>

          <button
              type="submit"
              :disabled="loading"
              class="w-full px-4 py-2 bg-forest-green text-white hover:bg-forest-dark rounded-lg disabled:opacity-50 transition-colors font-medium"
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
    // Редирект на главную страницу после успешного входа
    router.push('/')
  } catch (error: any) {
    loginError.value = error.response?.data?.error || 'Ошибка входа'
  } finally {
    loading.value = false
  }
}

// Если пользователь уже авторизован, редиректим на главную
onMounted(() => {
  if (authStore.isAuthenticated) {
    router.push('/')
  }
})
</script>
