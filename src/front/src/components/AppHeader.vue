<template>
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
</template>

<script setup lang="ts">
import { useAuthStore } from '@/stores/auth'
import { useRouter } from 'vue-router'

const authStore = useAuthStore()
const router = useRouter()

const handleLogout = () => {
  authStore.logout()
  router.push('/')
}
</script>