<template>
  <header
      v-if="authStore.isAuthenticated"
      class="backdrop-blur border-b border-border sticky top-0 z-10"
      style="background: color-mix(in srgb, var(--surface) 90%, transparent)"
  >
    <div class="page py-3 flex flex-wrap items-center justify-between gap-3 max-sm:flex-nowrap">
      <RouterLink to="/" class="flex items-center gap-2 text-ink no-underline">
        <BrandMark stem :size="22" />
        <span class="text-h3">Laritmo</span>
      </RouterLink>

      <!-- sm и шире: кнопки в строку -->
      <div class="hidden sm:flex flex-wrap items-center gap-2">
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

      <!-- ниже sm: кнопка-меню -->
      <div ref="menuRoot" class="sm:hidden relative">
        <button
            type="button"
            class="w-11 h-11 rounded-md bg-action-soft hover:bg-action-soft-hover text-ink transition-colors duration-150 flex items-center justify-center"
            aria-label="Меню"
            aria-haspopup="true"
            :aria-expanded="menuOpen"
            aria-controls="mobile-menu"
            @click="menuOpen = !menuOpen"
        >
          <i :class="['pi', menuOpen ? 'pi-times' : 'pi-bars']" aria-hidden="true"></i>
        </button>

        <div
            v-if="menuOpen"
            id="mobile-menu"
            class="absolute right-0 top-full mt-2 w-64 bg-surface border border-border shadow-md rounded-lg p-2 flex flex-col gap-1"
        >
          <div class="flex items-center gap-2 px-3 min-h-11 text-body-sm text-ink-muted">
            <span class="truncate">{{ authStore.user?.username }}</span>
            <AppChip v-if="authStore.isAdmin">админ</AppChip>
          </div>
          <button
              v-if="authStore.isOwner"
              type="button"
              class="menu-item"
              @click="runAction(() => router.push('/admin/users'))"
          >
            <i class="pi pi-users" aria-hidden="true"></i>
            Пользователи
          </button>
          <button type="button" class="menu-item" @click="runAction(() => (showChangePasswordDialog = true))">
            <i class="pi pi-key" aria-hidden="true"></i>
            Пароль
          </button>
          <button type="button" class="menu-item" @click="runAction(handleLogout)">
            <i class="pi pi-sign-out" aria-hidden="true"></i>
            Выйти
          </button>
          <div class="flex items-center justify-between gap-2 px-3 min-h-11 border-t border-border mt-1 pt-1 text-label">
            <span>Тема</span>
            <ThemeToggle />
          </div>
        </div>
      </div>
    </div>
  </header>

  <ChangePasswordDialog
      v-model="showChangePasswordDialog"
      @changed="handlePasswordChanged"
  />
</template>

<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { RouterLink, useRoute, useRouter } from 'vue-router'
import ThemeToggle from './ThemeToggle.vue'
import ChangePasswordDialog from './ChangePasswordDialog.vue'
import AppButton from './ui/AppButton.vue'
import AppChip from './ui/AppChip.vue'
import BrandMark from './ui/BrandMark.vue'

const authStore = useAuthStore()
const router = useRouter()
const route = useRoute()

const showChangePasswordDialog = ref(false)
const menuOpen = ref(false)
const menuRoot = ref<HTMLElement | null>(null)

const handleLogout = () => {
  authStore.logout()
  router.push('/')
}

const handlePasswordChanged = () => {
  // Password changed successfully
}

const runAction = (action: () => unknown) => {
  menuOpen.value = false
  action()
}

const onDocumentPointerDown = (e: Event) => {
  if (menuOpen.value && menuRoot.value && !menuRoot.value.contains(e.target as Node)) {
    menuOpen.value = false
  }
}

const onKeydown = (e: KeyboardEvent) => {
  if (e.key === 'Escape') menuOpen.value = false
}

watch(() => route.fullPath, () => {
  menuOpen.value = false
})

onMounted(() => {
  document.addEventListener('pointerdown', onDocumentPointerDown)
  document.addEventListener('keydown', onKeydown)
})

onBeforeUnmount(() => {
  document.removeEventListener('pointerdown', onDocumentPointerDown)
  document.removeEventListener('keydown', onKeydown)
})
</script>

<style scoped>
.menu-item {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  width: 100%;
  min-height: 2.75rem;
  padding: 0 0.75rem;
  border-radius: var(--radius-md);
  background: transparent;
  color: var(--ink);
  font-size: 0.875rem;
  line-height: 1.25rem;
  font-weight: 500;
  text-align: left;
  transition: background-color 0.15s ease;
}
.menu-item:hover {
  background: var(--sunken);
}
</style>
