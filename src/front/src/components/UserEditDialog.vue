<template>
  <Dialog
      v-model:visible="visible"
      modal
      :header="user ? 'Редактировать пользователя' : 'Создать пользователя'"
      :style="{ width: '500px' }"
      @hide="onHide"
  >
    <!-- Предупреждение для owner -->
    <div v-if="user?.role === 'owner'" class="alert mb-4">
      <i class="pi pi-exclamation-triangle" aria-hidden="true"></i>
      <span>Учётную запись владельца нельзя редактировать</span>
    </div>

    <form @submit.prevent="handleSubmit" class="space-y-4">
      <div>
        <label class="field-label">
          Username
        </label>
        <input
            v-model="form.username"
            type="text"
            required
            :disabled="!!user"
            class="field"
        />
      </div>

      <div>
        <label class="field-label">
          Email
        </label>
        <input
            v-model="form.email"
            type="email"
            required
            :disabled="isOwner"
            class="field"
        />
      </div>

      <div v-if="!user">
        <label class="field-label">
          Пароль
        </label>
        <input
            v-model="form.password"
            type="password"
            required
            minlength="8"
            placeholder="Минимум 8 символов"
            class="field"
        />
      </div>

      <div>
        <label class="field-label">
          Роль
        </label>
        <select
            v-model="form.role"
            required
            :disabled="isOwner"
            class="field"
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
            class="w-4 h-4 rounded border-border text-action disabled:cursor-not-allowed"
        />
        <label for="is_active" class="text-sm font-medium text-ink-muted">
          Активен
        </label>
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
            v-if="!isOwner"
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
