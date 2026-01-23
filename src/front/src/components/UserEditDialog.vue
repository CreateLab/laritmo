<template>
  <Dialog
      v-model:visible="visible"
      modal
      :header="user ? 'Редактировать пользователя' : 'Создать пользователя'"
      :style="{ width: '500px' }"
      @hide="onHide"
  >
    <!-- Предупреждение для owner -->
    <div v-if="user?.role === 'owner'" class="mb-4 p-3 bg-yellow-50 dark:bg-yellow-900/30 border border-yellow-300 dark:border-yellow-600 rounded-lg">
      <p class="text-yellow-800 dark:text-yellow-200 text-sm">
        ⚠️ Учётную запись владельца нельзя редактировать
      </p>
    </div>

    <form @submit.prevent="handleSubmit" class="space-y-4">
      <div>
        <label class="block text-sm font-medium text-gray-700 dark:text-dark-text-secondary mb-1 transition-colors duration-300">
          Username
        </label>
        <input
            v-model="form.username"
            type="text"
            required
            :disabled="!!user"
            class="w-full px-3 py-2 border border-gray-300 dark:border-dark-border rounded-lg focus:outline-none focus:ring-2 focus:ring-forest-green dark:focus:ring-forest-green-dark bg-white dark:bg-dark-bg text-gray-900 dark:text-dark-text transition-colors duration-300 disabled:bg-gray-100 dark:disabled:bg-dark-surface disabled:cursor-not-allowed"
        />
      </div>

      <div>
        <label class="block text-sm font-medium text-gray-700 dark:text-dark-text-secondary mb-1 transition-colors duration-300">
          Email
        </label>
        <input
            v-model="form.email"
            type="email"
            required
            :disabled="isOwner"
            class="w-full px-3 py-2 border border-gray-300 dark:border-dark-border rounded-lg focus:outline-none focus:ring-2 focus:ring-forest-green dark:focus:ring-forest-green-dark bg-white dark:bg-dark-bg text-gray-900 dark:text-dark-text transition-colors duration-300 disabled:bg-gray-100 dark:disabled:bg-dark-surface disabled:cursor-not-allowed"
        />
      </div>

      <div v-if="!user">
        <label class="block text-sm font-medium text-gray-700 dark:text-dark-text-secondary mb-1 transition-colors duration-300">
          Пароль
        </label>
        <input
            v-model="form.password"
            type="password"
            required
            minlength="8"
            placeholder="Минимум 8 символов"
            class="w-full px-3 py-2 border border-gray-300 dark:border-dark-border rounded-lg focus:outline-none focus:ring-2 focus:ring-forest-green dark:focus:ring-forest-green-dark bg-white dark:bg-dark-bg text-gray-900 dark:text-dark-text transition-colors duration-300"
        />
      </div>

      <div>
        <label class="block text-sm font-medium text-gray-700 dark:text-dark-text-secondary mb-1 transition-colors duration-300">
          Роль
        </label>
        <select
            v-model="form.role"
            required
            :disabled="isOwner"
            class="w-full px-3 py-2 border border-gray-300 dark:border-dark-border rounded-lg focus:outline-none focus:ring-2 focus:ring-forest-green dark:focus:ring-forest-green-dark bg-white dark:bg-dark-bg text-gray-900 dark:text-dark-text transition-colors duration-300 disabled:bg-gray-100 dark:disabled:bg-dark-surface disabled:cursor-not-allowed"
        >
          <option value="admin">Администратор</option>
          <option value="student">Студент</option>
        </select>
      </div>

      <div v-if="user" class="flex items-center gap-2">
        <input
            id="is_active"
            v-model="form.is_active"
            type="checkbox"
            :disabled="isOwner"
            class="w-4 h-4 rounded border-gray-300 dark:border-dark-border text-forest-green focus:ring-forest-green dark:focus:ring-forest-green-dark disabled:cursor-not-allowed"
        />
        <label for="is_active" class="text-sm font-medium text-gray-700 dark:text-dark-text-secondary transition-colors duration-300">
          Активен
        </label>
      </div>

      <div v-if="error" class="text-red-600 dark:text-red-400 text-sm bg-red-50 dark:bg-red-900/30 p-3 rounded-lg transition-colors duration-300">
        {{ error }}
      </div>

      <div class="flex justify-end gap-2 pt-4">
        <button
            type="button"
            @click="visible = false"
            class="px-4 py-2 text-sm bg-gray-100 dark:bg-dark-surface hover:bg-gray-200 dark:hover:bg-dark-border rounded-lg transition-colors duration-300 text-gray-700 dark:text-dark-text"
        >
          Отмена
        </button>
        <button
            v-if="!isOwner"
            type="submit"
            :disabled="loading"
            class="px-4 py-2 text-sm bg-forest-green dark:bg-forest-green-dark text-white hover:bg-forest-dark dark:hover:bg-forest-green rounded-lg disabled:opacity-50 transition-colors duration-300"
        >
          {{ loading ? 'Сохранение...' : 'Сохранить' }}
        </button>
      </div>
    </form>
  </Dialog>
</template>

<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import Dialog from 'primevue/dialog'
import { usersApi, type User } from '@/api/users'

const props = defineProps<{
  modelValue: boolean
  user?: User | null
}>()

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  'saved': []
}>()

const visible = ref(props.modelValue)
const form = ref({
  username: '',
  email: '',
  password: '',
  role: 'admin' as 'admin' | 'student',
  is_active: true,
})
const error = ref('')
const loading = ref(false)

const isOwner = computed(() => props.user?.role === 'owner')

watch(() => props.modelValue, (val) => {
  visible.value = val
  if (val && props.user) {
    form.value = {
      username: props.user.username,
      email: props.user.email,
      password: '',
      role: props.user.role === 'owner' ? 'admin' : props.user.role,
      is_active: props.user.is_active,
    }
  } else if (val) {
    form.value = {
      username: '',
      email: '',
      password: '',
      role: 'admin',
      is_active: true,
    }
  }
})

watch(visible, (val) => {
  emit('update:modelValue', val)
})

const handleSubmit = async () => {
  if (isOwner.value) return

  error.value = ''
  loading.value = true

  try {
    if (props.user) {
      await usersApi.update(props.user.id, {
        email: form.value.email,
        role: form.value.role,
        is_active: form.value.is_active,
      })
    } else {
      await usersApi.create({
        username: form.value.username,
        email: form.value.email,
        password: form.value.password,
        role: form.value.role,
      })
    }

    visible.value = false
    emit('saved')
  } catch (err: any) {
    error.value = err.response?.data?.error || 'Ошибка сохранения'
  } finally {
    loading.value = false
  }
}

const onHide = () => {
  form.value = {
    username: '',
    email: '',
    password: '',
    role: 'admin',
    is_active: true,
  }
  error.value = ''
}
</script>
