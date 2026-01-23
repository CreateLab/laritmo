<template>
  <div class="min-h-screen p-6">
    <header class="mb-8">
      <button
          @click="router.push('/')"
          class="text-forest-green dark:text-forest-green-dark hover:text-forest-dark dark:hover:text-forest-green mb-4 flex items-center gap-2 transition-colors duration-300"
      >
        ← На главную
      </button>

      <div class="flex justify-between items-center">
        <h1 class="text-3xl font-bold text-forest-dark dark:text-dark-text transition-colors duration-300">
          👥 Управление пользователями
        </h1>
        <button
            @click="handleCreate"
            class="px-4 py-2 bg-forest-green dark:bg-forest-green-dark text-white rounded-lg hover:bg-forest-dark dark:hover:bg-forest-green transition-colors duration-300"
        >
          ➕ Добавить пользователя
        </button>
      </div>
    </header>

    <main>
      <div v-if="loading" class="text-center py-12">
        <p class="text-gray-600 dark:text-dark-text-secondary transition-colors duration-300">Загрузка пользователей...</p>
      </div>

      <div v-else-if="error" class="text-center py-12">
        <p class="text-red-600 dark:text-red-400 transition-colors duration-300">{{ error }}</p>
        <button
            @click="loadUsers"
            class="mt-4 px-4 py-2 bg-forest-green dark:bg-forest-green-dark text-white rounded-lg hover:bg-forest-dark dark:hover:bg-forest-green transition-colors duration-300"
        >
          Попробовать снова
        </button>
      </div>

      <div v-else-if="users.length === 0" class="text-center py-12">
        <p class="text-gray-600 dark:text-dark-text-secondary transition-colors duration-300">Пользователи не найдены</p>
      </div>

      <div v-else class="space-y-4">
        <div
            v-for="user in users"
            :key="user.id"
            class="bg-white dark:bg-dark-surface rounded-xl shadow-md dark:shadow-lg p-6 transition-all duration-300 border-2 border-transparent hover:border-forest-mint dark:hover:border-forest-mint-dark"
        >
          <div class="flex justify-between items-start">
            <div class="flex-1">
              <div class="flex items-center gap-3 mb-2">
                <h3 class="text-xl font-semibold text-forest-dark dark:text-dark-text transition-colors duration-300">
                  {{ user.username }}
                </h3>
                <span
                    class="px-2 py-1 rounded-full text-xs font-medium"
                    :class="getRoleBadgeClass(user.role)"
                >
                  {{ getRoleLabel(user.role) }}
                </span>
                <span
                    class="px-2 py-1 rounded-full text-xs font-medium"
                    :class="user.is_active 
                      ? 'bg-green-100 dark:bg-green-900/30 text-green-800 dark:text-green-300' 
                      : 'bg-red-100 dark:bg-red-900/30 text-red-800 dark:text-red-300'"
                >
                  {{ user.is_active ? 'Активен' : 'Неактивен' }}
                </span>
              </div>

              <p class="text-gray-600 dark:text-dark-text-secondary transition-colors duration-300">
                {{ user.email }}
              </p>

              <p v-if="user.last_login_at" class="text-sm text-gray-500 dark:text-gray-400 mt-2 transition-colors duration-300">
                Последний вход: {{ formatDate(user.last_login_at) }}
              </p>
              <p v-else class="text-sm text-gray-500 dark:text-gray-400 mt-2 transition-colors duration-300">
                Ещё не входил
              </p>
            </div>

            <div v-if="user.role !== 'owner'" class="flex gap-2">
              <button
                  @click="handleEdit(user)"
                  title="Редактировать"
                  class="px-3 py-2 text-sm bg-gray-100 dark:bg-dark-border hover:bg-gray-200 dark:hover:bg-gray-600 rounded-lg transition-colors duration-300"
              >
                ✏️
              </button>
              <button
                  @click="handleResetPassword(user)"
                  title="Сбросить пароль"
                  class="px-3 py-2 text-sm bg-blue-100 dark:bg-blue-900/30 hover:bg-blue-200 dark:hover:bg-blue-900/50 text-blue-600 dark:text-blue-400 rounded-lg transition-colors duration-300"
              >
                🔑
              </button>
              <button
                  @click="handleToggleActive(user)"
                  :title="user.is_active ? 'Деактивировать' : 'Активировать'"
                  class="px-3 py-2 text-sm rounded-lg transition-colors duration-300"
                  :class="user.is_active 
                    ? 'bg-yellow-100 dark:bg-yellow-900/30 hover:bg-yellow-200 dark:hover:bg-yellow-900/50 text-yellow-600 dark:text-yellow-400'
                    : 'bg-green-100 dark:bg-green-900/30 hover:bg-green-200 dark:hover:bg-green-900/50 text-green-600 dark:text-green-400'"
              >
                {{ user.is_active ? '🔒' : '🔓' }}
              </button>
              <button
                  @click="handleDelete(user)"
                  title="Удалить"
                  class="px-3 py-2 text-sm bg-red-100 dark:bg-red-900/30 hover:bg-red-200 dark:hover:bg-red-900/50 text-red-600 dark:text-red-400 rounded-lg transition-colors duration-300"
              >
                🗑️
              </button>
            </div>

            <div v-else class="text-sm text-gray-500 dark:text-gray-400 italic">
              👑 Владелец (защищён)
            </div>
          </div>
        </div>
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
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { usersApi, type User } from '@/api/users'
import UserEditDialog from '@/components/UserEditDialog.vue'
import ResetPasswordDialog from '@/components/ResetPasswordDialog.vue'

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

const getRoleBadgeClass = (role: string): string => {
  const classes: Record<string, string> = {
    owner: 'bg-purple-100 dark:bg-purple-900/30 text-purple-800 dark:text-purple-300',
    admin: 'bg-blue-100 dark:bg-blue-900/30 text-blue-800 dark:text-blue-300',
    student: 'bg-gray-100 dark:bg-gray-700 text-gray-800 dark:text-gray-300',
  }
  return classes[role] || 'bg-gray-100 dark:bg-gray-700 text-gray-800 dark:text-gray-300'
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
