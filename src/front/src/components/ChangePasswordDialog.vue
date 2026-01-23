<template>
  <Dialog
      v-model:visible="visible"
      modal
      header="Изменить пароль"
      :style="{ width: '400px' }"
      @hide="onHide"
  >
    <form @submit.prevent="handleSubmit" class="space-y-4">
      <div>
        <label class="block text-sm font-medium text-gray-700 dark:text-dark-text-secondary mb-1 transition-colors duration-300">
          Текущий пароль
        </label>
        <input
            v-model="form.oldPassword"
            type="password"
            required
            class="w-full px-3 py-2 border border-gray-300 dark:border-dark-border rounded-lg focus:outline-none focus:ring-2 focus:ring-forest-green dark:focus:ring-forest-green-dark bg-white dark:bg-dark-bg text-gray-900 dark:text-dark-text transition-colors duration-300"
        />
      </div>

      <div>
        <label class="block text-sm font-medium text-gray-700 dark:text-dark-text-secondary mb-1 transition-colors duration-300">
          Новый пароль
        </label>
        <input
            v-model="form.newPassword"
            type="password"
            required
            minlength="8"
            placeholder="Минимум 8 символов"
            class="w-full px-3 py-2 border border-gray-300 dark:border-dark-border rounded-lg focus:outline-none focus:ring-2 focus:ring-forest-green dark:focus:ring-forest-green-dark bg-white dark:bg-dark-bg text-gray-900 dark:text-dark-text transition-colors duration-300"
        />
      </div>

      <div>
        <label class="block text-sm font-medium text-gray-700 dark:text-dark-text-secondary mb-1 transition-colors duration-300">
          Подтвердите новый пароль
        </label>
        <input
            v-model="form.confirmPassword"
            type="password"
            required
            class="w-full px-3 py-2 border border-gray-300 dark:border-dark-border rounded-lg focus:outline-none focus:ring-2 focus:ring-forest-green dark:focus:ring-forest-green-dark bg-white dark:bg-dark-bg text-gray-900 dark:text-dark-text transition-colors duration-300"
        />
      </div>

      <div v-if="error" class="text-red-600 dark:text-red-400 text-sm bg-red-50 dark:bg-red-900/30 p-3 rounded-lg transition-colors duration-300">
        {{ error }}
      </div>

      <div v-if="success" class="text-green-600 dark:text-green-400 text-sm bg-green-50 dark:bg-green-900/30 p-3 rounded-lg transition-colors duration-300">
        {{ success }}
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
import { ref, watch } from 'vue'
import Dialog from 'primevue/dialog'
import { useAuthStore } from '@/stores/auth'

const props = defineProps<{
  modelValue: boolean
}>()

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  'changed': []
}>()

const authStore = useAuthStore()

const visible = ref(props.modelValue)
const form = ref({
  oldPassword: '',
  newPassword: '',
  confirmPassword: '',
})
const error = ref('')
const success = ref('')
const loading = ref(false)

watch(() => props.modelValue, (val) => {
  visible.value = val
})

watch(visible, (val) => {
  emit('update:modelValue', val)
})

const handleSubmit = async () => {
  error.value = ''
  success.value = ''

  if (form.value.newPassword !== form.value.confirmPassword) {
    error.value = 'Пароли не совпадают'
    return
  }

  if (form.value.newPassword.length < 8) {
    error.value = 'Новый пароль должен содержать минимум 8 символов'
    return
  }

  loading.value = true

  try {
    await authStore.changeMyPassword(form.value.oldPassword, form.value.newPassword)
    success.value = 'Пароль успешно изменён'
    
    setTimeout(() => {
      visible.value = false
      emit('changed')
    }, 1500)
  } catch (err: any) {
    error.value = err.response?.data?.error || 'Ошибка смены пароля'
  } finally {
    loading.value = false
  }
}

const onHide = () => {
  form.value = {
    oldPassword: '',
    newPassword: '',
    confirmPassword: '',
  }
  error.value = ''
  success.value = ''
}
</script>
