<template>
  <div>
    <label v-if="label" :for="id" class="field-label">{{ label }}</label>
    <textarea
      v-if="multiline"
      :id="id"
      v-model="model"
      class="field"
      v-bind="$attrs"
    ></textarea>
    <input v-else :id="id" v-model="model" class="field" v-bind="$attrs" />
  </div>
</template>

<script setup lang="ts">
defineOptions({ inheritAttrs: false })

const model = defineModel<string | number | undefined>()
defineProps<{ label?: string; multiline?: boolean }>()

const id = `tf-${Math.random().toString(36).slice(2, 9)}`
</script>
