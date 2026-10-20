<template>
  <div class="min-h-screen flex items-center justify-center p-6 relative">
    <div class="absolute top-4 right-4">
      <ThemeToggle />
    </div>

    <div class="w-full max-w-md">
      <div class="flex items-center justify-center gap-3 mb-8">
        <BrandMark stem :size="32" />
        <h1 class="text-display">Laritmo</h1>
      </div>

      <div class="bg-surface border border-border rounded-lg shadow-md p-8">
        <h2 class="text-h2 mb-6 text-center">Вход в систему</h2>

        <form @submit.prevent="handleLogin" class="space-y-4">
          <TextField v-model="loginForm.username" label="Username" type="text" required autofocus />
          <TextField v-model="loginForm.password" label="Password" type="password" required />

          <AppAlert v-if="loginError">{{ loginError }}</AppAlert>

          <AppButton type="submit" variant="primary" :disabled="loading" class="w-full">
            {{ loading ? 'Вход...' : 'Войти' }}
          </AppButton>
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
import AppButton from '@/components/ui/AppButton.vue'
import AppAlert from '@/components/ui/AppAlert.vue'
import BrandMark from '@/components/ui/BrandMark.vue'
import TextField from '@/components/ui/TextField.vue'

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
