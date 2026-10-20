<template>
  <Dialog
      v-model:visible="visible"
      modal
      header="Сбросить пароль"
      :style="{ width: '400px' }"
      @hide="onHide"
  >
    <div class="alert mb-4">
      <i class="pi pi-exclamation-triangle" aria-hidden="true"></i>
      <span>Новый пароль будет установлен пользователю. Убедитесь, что передадите его безопасным способом.</span>
    </div>

    <form @submit.prevent="handleSubmit" class="space-y-4">
      <div>
        <label class="field-label">
          Новый пароль
        </label>
        <input
            v-model="form.newPassword"
            type="password"
            required
            minlength="8"
            placeholder="Минимум 8 символов"
            class="field"
        />
      </div>

      <div>
        <label class="field-label">
          Подтвердите пароль
        </label>
        <input
            v-model="form.confirmPassword"
            type="password"
            required
            class="field"
        />
      </div>

      <div v-if="error" class="alert">
        <i class="pi pi-exclamation-triangle" aria-hidden="true"></i>
        <span>{{ error }}</span>
      </div>

      <div class="flex justify-end gap-2 pt-4">
        <button
            type="button"
            @click="visible = false"
            class="btn btn-secondary"
        >
          Отмена
        </button>
        <button
            type="submit"
            :disabled="loading"
            class="btn btn-primary"
        >
          {{ loading ? 'Сохранение...' : 'Сбросить пароль' }}
        </button>
      </div>
    </form>
  </Dialog>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import Dialog from 'primevue/dialog'
import { usersApi } from '@/api/users'

const props = defineProps<{
  modelValue: boolean
  userId: number | null
}>()

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  'reset': []
}>()

const visible = ref(props.modelValue)
const form = ref({
  newPassword: '',
  confirmPassword: '',
})
const error = ref('')
const loading = ref(false)

watch(() => props.modelValue, (val) => {
  visible.value = val
})

watch(visible, (val) => {
  emit('update:modelValue', val)
})

const handleSubmit = async () => {
  error.value = ''

  if (!props.userId) {
    error.value = 'Пользователь не выбран'
    return
  }

  if (form.value.newPassword !== form.value.confirmPassword) {
    error.value = 'Пароли не совпадают'
    return
  }

  if (form.value.newPassword.length < 8) {
    error.value = 'Пароль должен содержать минимум 8 символов'
    return
  }

  loading.value = true

  try {
    await usersApi.resetPassword(props.userId, form.value.newPassword)
    visible.value = false
    emit('reset')
  } catch (err: any) {
    error.value = err.response?.data?.error || 'Ошибка сброса пароля'
  } finally {
    loading.value = false
  }
}

const onHide = () => {
  form.value = {
    newPassword: '',
    confirmPassword: '',
  }
  error.value = ''
}
</script>
