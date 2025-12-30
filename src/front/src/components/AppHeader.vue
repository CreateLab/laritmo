<template>
  <!-- Header только для залогиненных -->
  <header v-if="authStore.isAuthenticated" class="bg-white/80 backdrop-blur-sm border-b border-gray-200/50 sticky top-0 z-10">
    <div class="max-w-7xl mx-auto px-6 py-3 flex items-center justify-end">
      <div class="flex items-center gap-3">
        <span class="text-sm text-gray-600">
          {{ authStore.user?.username }}
          <span v-if="authStore.isAdmin" class="text-forest-green font-semibold">(админ)</span>
        </span>
        <button
            @click="handleLogout"
            class="px-4 py-2 text-sm bg-gray-100 hover:bg-gray-200 rounded-lg transition-colors"
        >
          Выйти
        </button>
      </div>
    </div>
  </header>

  <!-- Плавающая кнопка для незалогиненных -->
  <button
      v-else
      @click="showLoginDialog = true"
      class="fixed top-4 right-4 z-50 px-4 py-2 text-sm bg-forest-green text-white hover:bg-forest-dark rounded-lg shadow-lg transition-colors"
  >
    Войти
  </button>

  <!-- Диалог входа -->
  <Dialog
      v-model:visible="showLoginDialog"
      modal
      header="Вход в систему"
      :style="{ width: '400px' }"
  >
    <form @submit.prevent="handleLogin" class="space-y-4">
      <div>
        <label class="block text-sm font-medium text-gray-700 mb-1">Username</label>
        <input
            v-model="loginForm.username"
            type="text"
            required
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

      <div v-if="loginError" class="text-red-600 text-sm">
        {{ loginError }}
      </div>

      <div class="flex justify-end gap-2">
        <button
            type="button"
            @click="showLoginDialog = false"
            class="px-4 py-2 text-sm bg-gray-100 hover:bg-gray-200 rounded-lg"
        >
          Отмена
        </button>
        <button
            type="submit"
            :disabled="loading"
            class="px-4 py-2 text-sm bg-forest-green text-white hover:bg-forest-dark rounded-lg disabled:opacity-50"
        >
          {{ loading ? 'Вход...' : 'Войти' }}
        </button>
      </div>
    </form>
  </Dialog>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useRouter } from 'vue-router'
import Dialog from 'primevue/dialog'

const authStore = useAuthStore()
const router = useRouter()

const showLoginDialog = ref(false)
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
    showLoginDialog.value = false
    loginForm.value = { username: '', password: '' }
  } catch (error: any) {
    loginError.value = error.response?.data?.error || 'Ошибка входа'
  } finally {
    loading.value = false
  }
}

const handleLogout = () => {
  authStore.logout()
  router.push('/')
}
</script>