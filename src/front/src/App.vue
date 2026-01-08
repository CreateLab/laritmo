<template>
  <AppHeader v-if="!isAuthPage" />
  <RouterView />
  <Toast />
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { RouterView, useRoute } from 'vue-router'
import AppHeader from './components/AppHeader.vue'
import Toast from 'primevue/toast'
import { useThemeStore } from './stores/theme'

const route = useRoute()
const isAuthPage = computed(() => route.name === 'auth')

const themeStore = useThemeStore()
onMounted(() => {
  themeStore.initTheme()
})
</script>
