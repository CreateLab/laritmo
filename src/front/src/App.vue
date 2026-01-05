<template>
  <AppHeader v-if="!isAuthPage" />
  <RouterView />
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { RouterView, useRoute } from 'vue-router'
import AppHeader from './components/AppHeader.vue'

const route = useRoute()
const isAuthPage = computed(() => route.name === 'auth')
</script>