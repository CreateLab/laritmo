<template>
  <div class="min-h-screen flex flex-col">
    <AppHeader v-if="!isAuthPage" />
    <HeaderStrip v-if="!isAuthPage" />
    <div class="flex-1 flex flex-col">
      <RouterView />
    </div>
    <footer class="border-t border-border mt-12">
      <div class="page py-6 flex flex-wrap items-center justify-between gap-2 text-body-sm text-ink-muted">
        <span>Laritmo · {{ year }}</span>
        <a href="https://github.com/CreateLab/laritmo" target="_blank" rel="noopener" class="text-ink-muted underline hover:text-ink">
          Репозиторий
        </a>
      </div>
    </footer>
  </div>
  <Toast />
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { RouterView, useRoute } from 'vue-router'
import AppHeader from './components/AppHeader.vue'
import HeaderStrip from './components/ui/HeaderStrip.vue'
import Toast from 'primevue/toast'
import { useThemeStore } from './stores/theme'

const route = useRoute()
const isAuthPage = computed(() => route.name === 'auth')
const year = new Date().getFullYear()

const themeStore = useThemeStore()
onMounted(() => {
  themeStore.initTheme()
})
</script>
