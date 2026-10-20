<template>
  <Dialog
      v-model:visible="visible"
      modal
      :header="gradeSheet ? 'Изменить ссылку на журнал' : 'Добавить ссылку на журнал'"
      :style="{ width: '500px' }"
      @hide="onHide"
  >
    <form @submit.prevent="handleSubmit" class="space-y-4">
      <div>
        <label class="field-label">URL журнала (Google Sheets)</label>
        <input
            v-model="form.sheet_url"
            type="url"
            required
            placeholder="https://docs.google.com/spreadsheets/d/..."
            class="field"
        />
      </div>

      <div>
        <label class="field-label">Описание (опционально)</label>
        <input
            v-model="form.description"
            type="text"
            placeholder="Например: Основной журнал группы"
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
          {{ loading ? 'Сохранение...' : 'Сохранить' }}
        </button>
      </div>
    </form>
  </Dialog>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import Dialog from 'primevue/dialog'
import { gradeSheetsApi, type GradeSheet } from '@/api/gradesheets'

const props = defineProps<{
  modelValue: boolean
  gradeSheet?: GradeSheet | null
  courseId: number
}>()

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  'saved': []
}>()

const visible = ref(props.modelValue)
const form = ref({
  sheet_url: '',
  description: '',
})
const error = ref('')
const loading = ref(false)

watch(() => props.modelValue, (val) => {
  visible.value = val
  if (val && props.gradeSheet) {
    form.value = {
      sheet_url: props.gradeSheet.sheet_url,
      description: props.gradeSheet.description || '',
    }
  } else if (val) {
    form.value = {
      sheet_url: '',
      description: '',
    }
  }
})

watch(visible, (val) => {
  emit('update:modelValue', val)
})

const handleSubmit = async () => {
  error.value = ''
  loading.value = true

  try {
    const data = {
      course_id: props.courseId,
      sheet_url: form.value.sheet_url,
      description: form.value.description || undefined,
    }

    if (props.gradeSheet) {
      await gradeSheetsApi.update(props.gradeSheet.id, data)
    } else {
      await gradeSheetsApi.create(data)
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
  form.value = { sheet_url: '', description: '' }
  error.value = ''
}
</script>
