<template>
  <main class="page flex-1 flex items-center justify-center py-12">
    <EmptyState title="Страница не найдена" text="Такого адреса нет. Возможно, ссылка устарела или в ней опечатка.">
      <AppButton variant="primary" icon="pi-home" to="/">На главную</AppButton>
    </EmptyState>
  </main>
</template>

<script setup lang="ts">
import EmptyState from '@/components/ui/EmptyState.vue'
import AppButton from '@/components/ui/AppButton.vue'
</script>
