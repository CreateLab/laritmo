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
        <label class="field-label">
          Текущий пароль
        </label>
        <input
            v-model="form.oldPassword"
            type="password"
            required
            class="field"
        />
      </div>

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
          Подтвердите новый пароль
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

      <div v-if="success" class="flex items-start gap-2 p-3 bg-sunken text-ink rounded-md text-body-sm">
        <i class="pi pi-check text-action mt-0.5" aria-hidden="true"></i>
        <span>{{ success }}</span>
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
