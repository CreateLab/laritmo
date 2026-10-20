<template>
  <!-- с ножкой: шапка и вход -->
  <svg v-if="stem" :width="size" :height="size" viewBox="0 0 20 20" aria-hidden="true" class="shrink-0">
    <path d="M2 11a8 7.5 0 0 1 16 0z" fill="var(--accent)" />
    <circle cx="7" cy="8" r="1.4" fill="var(--surface)" />
    <circle cx="12.5" cy="6.8" r="1" fill="var(--surface)" />
    <rect x="7.5" y="11" width="5" height="6.5" rx="2" fill="var(--border-strong)" />
  </svg>
  <!-- без ножки: карточки и содержимое -->
  <svg v-else :width="size" :height="size / 2" viewBox="0 0 20 10" aria-hidden="true" class="shrink-0">
    <path d="M1 10a9 9 0 0 1 18 0z" fill="var(--accent)" />
    <circle cx="7" cy="6.2" r="1.5" fill="var(--surface)" />
    <circle cx="13" cy="5" r="1.1" fill="var(--surface)" />
  </svg>
</template>

<script setup lang="ts">
withDefaults(defineProps<{ stem?: boolean; size?: number }>(), { stem: false, size: 16 })
</script>
