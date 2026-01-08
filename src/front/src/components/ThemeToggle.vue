<template>
  <button
      @click="handleToggle"
      class="w-10 h-10 rounded-lg bg-gray-100 dark:bg-dark-surface hover:bg-gray-200 dark:hover:bg-dark-border transition-colors duration-300 flex items-center justify-center text-lg"
      :title="themeStore.effectiveTheme === 'dark' ? 'Переключить на светлую тему' : 'Переключить на темную тему'"
  >
    <span v-if="themeStore.effectiveTheme === 'dark'">☀️</span>
    <span v-else>🌙</span>
  </button>
</template>

<script setup lang="ts">
import { useThemeStore } from '@/stores/theme'

const themeStore = useThemeStore()

const handleToggle = () => {
  themeStore.toggleTheme()
}
</script>
