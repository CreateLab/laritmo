<template>
  <button
      @click="handleToggle"
      class="w-10 h-10 rounded-full bg-action-soft hover:bg-action-soft-hover text-ink transition-colors duration-150 flex items-center justify-center"
      :title="label"
      :aria-label="label"
  >
    <i :class="['pi', themeStore.effectiveTheme === 'dark' ? 'pi-sun' : 'pi-moon']" aria-hidden="true"></i>
  </button>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useThemeStore } from '@/stores/theme'

const themeStore = useThemeStore()

const label = computed(() =>
  themeStore.effectiveTheme === 'dark' ? 'Переключить на светлую тему' : 'Переключить на темную тему'
)

const handleToggle = () => {
  themeStore.toggleTheme()
}
</script>
