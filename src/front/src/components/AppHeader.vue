<template>
  <header v-if="authStore.isAuthenticated" class="bg-white/80 dark:bg-dark-surface/80 backdrop-blur-sm border-b border-gray-200/50 dark:border-dark-border/50 sticky top-0 z-10 transition-colors duration-300">
    <div class="max-w-7xl mx-auto px-6 py-3 flex items-center justify-end">
      <div class="flex items-center gap-3">
        <ThemeToggle />
        <span class="text-sm text-gray-600 dark:text-dark-text-secondary">
          {{ authStore.user?.username }}
          <span v-if="authStore.isAdmin" class="text-forest-green dark:text-forest-green-dark font-semibold">(админ)</span>
        </span>
        <button
            @click="handleLogout"
            class="px-4 py-2 text-sm bg-gray-100 dark:bg-dark-surface hover:bg-gray-200 dark:hover:bg-dark-border rounded-lg transition-colors duration-300 text-gray-700 dark:text-dark-text"
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
import ThemeToggle from './ThemeToggle.vue'

const authStore = useAuthStore()
const router = useRouter()

const handleLogout = () => {
  authStore.logout()
  router.push('/')
}
</script>