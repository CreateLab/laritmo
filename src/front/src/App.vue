<template>
  <AppHeader v-if="!isAuthPage" />
  <HeaderStrip v-if="!isAuthPage && authStore.isAuthenticated" />
  <RouterView />
  <Toast />
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { RouterView, useRoute } from 'vue-router'
import AppHeader from './components/AppHeader.vue'
import HeaderStrip from './components/ui/HeaderStrip.vue'
import Toast from 'primevue/toast'
import { useThemeStore } from './stores/theme'
import { useAuthStore } from './stores/auth'

const route = useRoute()
const authStore = useAuthStore()
const isAuthPage = computed(() => route.name === 'auth')

const themeStore = useThemeStore()
onMounted(() => {
  themeStore.initTheme()
})
</script>
