<template>
  <component
    :is="to ? RouterLink : 'button'"
    :to="to"
    :type="to ? undefined : type"
    :class="['btn', `btn-${variant}`]"
  >
    <i v-if="icon" :class="['pi', icon]" aria-hidden="true"></i>
    <slot />
  </component>
</template>

<script setup lang="ts">
import { RouterLink, type RouteLocationRaw } from 'vue-router'

withDefaults(
  defineProps<{
    variant?: 'primary' | 'secondary' | 'danger'
    icon?: string
    to?: RouteLocationRaw
    type?: 'button' | 'submit'
  }>(),
  { variant: 'secondary', type: 'button' }
)
</script>
