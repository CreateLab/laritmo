<template>
  <header
      v-if="authStore.isAuthenticated"
      class="backdrop-blur border-b border-border sticky top-0 z-10"
      style="background: color-mix(in srgb, var(--surface) 90%, transparent)"
  >
    <div class="page py-3 flex flex-wrap items-center justify-between gap-3">
      <RouterLink to="/" class="flex items-center gap-2 text-ink no-underline">
        <BrandMark stem :size="22" />
        <span class="text-h3">Laritmo</span>
      </RouterLink>

      <div class="flex flex-wrap items-center gap-2">
        <ThemeToggle />
        <span class="flex items-center gap-2 text-body-sm text-ink-muted">
          {{ authStore.user?.username }}
          <AppChip v-if="authStore.isAdmin">админ</AppChip>
        </span>
        <AppButton
            v-if="authStore.isOwner"
            @click="router.push('/admin/users')"
            title="Управление пользователями"
        >
          Пользователи
        </AppButton>
        <AppButton @click="showChangePasswordDialog = true" title="Изменить пароль">
          Пароль
        </AppButton>
        <AppButton @click="handleLogout">
          Выйти
        </AppButton>
      </div>
    </div>
  </header>

  <ChangePasswordDialog
      v-model="showChangePasswordDialog"
      @changed="handlePasswordChanged"
  />
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { RouterLink, useRouter } from 'vue-router'
import ThemeToggle from './ThemeToggle.vue'
import ChangePasswordDialog from './ChangePasswordDialog.vue'
import AppButton from './ui/AppButton.vue'
import AppChip from './ui/AppChip.vue'
import BrandMark from './ui/BrandMark.vue'

const authStore = useAuthStore()
const router = useRouter()

const showChangePasswordDialog = ref(false)

const handleLogout = () => {
  authStore.logout()
  router.push('/')
}

const handlePasswordChanged = () => {
  // Password changed successfully
}
</script>
