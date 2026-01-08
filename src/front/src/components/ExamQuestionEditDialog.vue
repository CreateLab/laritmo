<template>
  <Dialog
      v-model:visible="visible"
      modal
      :header="question ? 'Редактировать вопрос' : 'Создать вопрос'"
      :style="{ width: '700px', maxHeight: '90vh' }"
      @hide="onHide"
  >
    <form @submit.prevent="handleSubmit" class="space-y-4">
      <div class="grid grid-cols-2 gap-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 dark:text-dark-text-secondary mb-1 transition-colors duration-300">Номер</label>
          <input
              v-model.number="form.number"
              type="number"
              required
              min="1"
              class="w-full px-3 py-2 border border-gray-300 dark:border-dark-border rounded-lg focus:outline-none focus:ring-2 focus:ring-forest-green dark:focus:ring-forest-green-dark bg-white dark:bg-dark-bg text-gray-900 dark:text-dark-text transition-colors duration-300"
          />
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 dark:text-dark-text-secondary mb-1 transition-colors duration-300">Секция</label>
          <input
              v-model="form.section"
              type="text"
              required
              placeholder="Например: Основы ASP.NET Core"
              class="w-full px-3 py-2 border border-gray-300 dark:border-dark-border rounded-lg focus:outline-none focus:ring-2 focus:ring-forest-green dark:focus:ring-forest-green-dark bg-white dark:bg-dark-bg text-gray-900 dark:text-dark-text transition-colors duration-300"
          />
        </div>
      </div>

      <div>
        <label class="block text-sm font-medium text-gray-700 dark:text-dark-text-secondary mb-1 transition-colors duration-300">Вопрос</label>
        <textarea
            v-model="form.question"
            required
            rows="6"
            placeholder="Введите текст вопроса..."
            class="w-full px-3 py-2 border border-gray-300 dark:border-dark-border rounded-lg focus:outline-none focus:ring-2 focus:ring-forest-green dark:focus:ring-forest-green-dark bg-white dark:bg-dark-bg text-gray-900 dark:text-dark-text resize-none transition-colors duration-300"
        ></textarea>
      </div>

      <div v-if="error" class="text-red-600 dark:text-red-400 text-sm transition-colors duration-300">
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
import { examQuestionsApi, type ExamQuestion } from '@/api/examquestions'

const props = defineProps<{
  modelValue: boolean
  question?: ExamQuestion | null
  courseId: number
}>()

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  'saved': []
}>()

const visible = ref(props.modelValue)

const form = ref({
  number: 1,
  section: '',
  question: '',
})
const error = ref('')
const loading = ref(false)

watch(() => props.modelValue, (val) => {
  visible.value = val
  if (val) {
    if (props.question) {
      form.value = {
        number: props.question.number,
        section: props.question.section,
        question: props.question.question,
      }
    } else {
      form.value = {
        number: 1,
        section: '',
        question: '',
      }
    }
    error.value = ''
  }
})

watch(visible, (val) => {
  emit('update:modelValue', val)
})

const handleSubmit = async () => {
  error.value = ''
  loading.value = true

  try {
    if (props.question) {
      await examQuestionsApi.update(props.question.id, {
        number: form.value.number,
        section: form.value.section,
        question: form.value.question,
      })
    } else {
      await examQuestionsApi.create({
        course_id: props.courseId,
        number: form.value.number,
        section: form.value.section,
        question: form.value.question,
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
  form.value = { number: 1, section: '', question: '' }
  error.value = ''
}
</script>
