<template>
  <div class="flex flex-col items-center text-center gap-3 py-12 px-4">
    <BrandMark :size="32" />
    <h2 class="text-h3">{{ title }}</h2>
    <p v-if="text" class="text-body-sm text-ink-muted max-w-md">{{ text }}</p>
    <div v-if="$slots.default" class="mt-2 flex flex-wrap justify-center gap-2">
      <slot />
    </div>
  </div>
</template>

<script setup lang="ts">
import BrandMark from './BrandMark.vue'

defineProps<{ title: string; text?: string }>()
</script>
