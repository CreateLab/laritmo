<template>
  <div class="min-h-screen max-w-7xl mx-auto px-4 sm:px-8 py-12">
    <header class="mb-8">
      <a
          href="#"
          @click.prevent="router.push('/')"
          class="text-action hover:text-action-hover text-sm no-underline inline-block mb-4"
      >
        ← На главную
      </a>

      <div class="flex flex-wrap justify-between items-center gap-4">
        <h1 class="text-h1">Управление пользователями</h1>
        <AppButton variant="primary" icon="pi-plus" @click="handleCreate">
          Добавить пользователя
        </AppButton>
      </div>
    </header>

    <main>
      <div v-if="loading" class="text-center py-12">
        <p class="text-ink-muted">Загрузка пользователей...</p>
      </div>

      <EmptyState v-else-if="error" title="Не удалось загрузить пользователей" :text="error">
        <AppButton variant="primary" icon="pi-refresh" @click="loadUsers">Попробовать снова</AppButton>
      </EmptyState>

      <EmptyState v-else-if="users.length === 0" title="Пользователи не найдены" />

      <div v-else class="table-wrap">
        <table class="data-table">
          <thead>
            <tr>
              <th>Пользователь</th>
              <th>Роль</th>
              <th>Статус</th>
              <th>Последний вход</th>
              <th><span class="sr-only">Действия</span></th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="user in users" :key="user.id">
              <td>
                <div class="font-medium">{{ user.username }}</div>
                <div class="text-ink-muted">{{ user.email }}</div>
              </td>
              <td><AppChip>{{ getRoleLabel(user.role) }}</AppChip></td>
              <td><AppChip>{{ user.is_active ? 'Активен' : 'Неактивен' }}</AppChip></td>
              <td class="text-ink-muted whitespace-nowrap">
                {{ user.last_login_at ? formatDate(user.last_login_at) : 'Ещё не входил' }}
              </td>
              <td>
                <div v-if="user.role !== 'owner'" class="flex justify-end gap-2">
                  <AppButton icon="pi-pencil" title="Редактировать" aria-label="Редактировать" class="!px-3" @click="handleEdit(user)" />
                  <AppButton icon="pi-key" title="Сбросить пароль" aria-label="Сбросить пароль" class="!px-3" @click="handleResetPassword(user)" />
                  <AppButton
                      :icon="user.is_active ? 'pi-lock' : 'pi-lock-open'"
                      :title="user.is_active ? 'Деактивировать' : 'Активировать'"
                      :aria-label="user.is_active ? 'Деактивировать' : 'Активировать'"
                      class="!px-3"
                      @click="handleToggleActive(user)"
                  />
                  <AppButton variant="danger" icon="pi-trash" title="Удалить" aria-label="Удалить" class="!px-3" @click="handleDelete(user)" />
                </div>
                <div v-else class="text-right text-ink-muted whitespace-nowrap">Владелец (защищён)</div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </main>

    <!-- Edit Dialog -->
    <UserEditDialog
        v-model="showEditDialog"
        :user="selectedUser"
        @saved="onUserSaved"
    />

    <!-- Reset Password Dialog -->
    <ResetPasswordDialog
        v-model="showResetPasswordDialog"
        :user-id="resetPasswordUserId"
        @reset="onPasswordReset"
    />
  </div>
</template>

<script setup lang="ts">
import EmptyState from '@/components/ui/EmptyState.vue'
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { usersApi, type User } from '@/api/users'
import UserEditDialog from '@/components/UserEditDialog.vue'
import ResetPasswordDialog from '@/components/ResetPasswordDialog.vue'
import AppButton from '@/components/ui/AppButton.vue'
import AppChip from '@/components/ui/AppChip.vue'

const router = useRouter()

const users = ref<User[]>([])
const loading = ref(true)
const error = ref('')

const showEditDialog = ref(false)
const selectedUser = ref<User | null>(null)

const showResetPasswordDialog = ref(false)
const resetPasswordUserId = ref<number | null>(null)

const loadUsers = async () => {
  loading.value = true
  error.value = ''

  try {
    const { data } = await usersApi.getAll()
    users.value = data
  } catch (err: any) {
    error.value = err.response?.data?.error || 'Ошибка загрузки пользователей'
  } finally {
    loading.value = false
  }
}

const handleCreate = () => {
  selectedUser.value = null
  showEditDialog.value = true
}

const handleEdit = (user: User) => {
  selectedUser.value = user
  showEditDialog.value = true
}

const handleDelete = async (user: User) => {
  if (!confirm(`Удалить пользователя "${user.username}"? Это действие нельзя отменить.`)) {
    return
  }

  try {
    await usersApi.delete(user.id)
    await loadUsers()
  } catch (err: any) {
    alert(err.response?.data?.error || 'Ошибка удаления пользователя')
  }
}

const handleToggleActive = async (user: User) => {
  const action = user.is_active ? 'деактивировать' : 'активировать'
  if (!confirm(`${user.is_active ? 'Деактивировать' : 'Активировать'} пользователя "${user.username}"?`)) {
    return
  }

  try {
    if (user.is_active) {
      await usersApi.deactivate(user.id)
    } else {
      await usersApi.activate(user.id)
    }
    await loadUsers()
  } catch (err: any) {
    alert(err.response?.data?.error || `Ошибка: не удалось ${action} пользователя`)
  }
}

const handleResetPassword = (user: User) => {
  resetPasswordUserId.value = user.id
  showResetPasswordDialog.value = true
}

const onUserSaved = async () => {
  await loadUsers()
}

const onPasswordReset = () => {
  resetPasswordUserId.value = null
}

const getRoleLabel = (role: string): string => {
  const labels: Record<string, string> = {
    owner: 'Владелец',
    admin: 'Администратор',
    student: 'Студент',
  }
  return labels[role] || role
}

const formatDate = (dateString: string): string => {
  const date = new Date(dateString)
  return date.toLocaleString('ru-RU', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}

onMounted(() => {
  loadUsers()
})
</script>
