<template>
  <Dialog
      v-model:visible="visible"
      modal
      header="Массовая загрузка вопросов"
      :style="{ width: '800px', maxHeight: '90vh' }"
      @hide="onHide"
  >
    <TabView>
      <TabPanel header="📤 Загрузить файл" value="upload">
        <div class="space-y-4">
          <div class="bg-blue-50 dark:bg-blue-900/30 border border-blue-200 dark:border-blue-800 rounded-lg p-4 transition-colors duration-300">
            <h4 class="font-semibold text-blue-900 dark:text-blue-200 mb-2 transition-colors duration-300">Поддерживаемые форматы:</h4>
            <ul class="text-sm text-blue-800 dark:text-blue-200 space-y-1 list-disc list-inside transition-colors duration-300">
              <li><strong>JSON</strong> (.json) - структурированный формат</li>
              <li><strong>CSV</strong> (.csv) - табличный формат</li>
            </ul>
          </div>

          <div class="bg-gray-50 dark:bg-dark-surface border border-gray-200 dark:border-dark-border rounded-lg p-4 transition-colors duration-300">
            <h4 class="font-semibold text-gray-900 dark:text-dark-text mb-2 transition-colors duration-300">Пример JSON:</h4>
            <pre class="text-xs bg-white dark:bg-dark-bg p-3 rounded border dark:border-dark-border overflow-x-auto text-gray-900 dark:text-dark-text transition-colors duration-300">{{ jsonExample }}</pre>
          </div>

          <div class="bg-gray-50 dark:bg-dark-surface border border-gray-200 dark:border-dark-border rounded-lg p-4 transition-colors duration-300">
            <h4 class="font-semibold text-gray-900 dark:text-dark-text mb-2 transition-colors duration-300">Пример CSV:</h4>
            <pre class="text-xs bg-white dark:bg-dark-bg p-3 rounded border dark:border-dark-border overflow-x-auto text-gray-900 dark:text-dark-text transition-colors duration-300">{{ csvExample }}</pre>
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-dark-text-secondary mb-2 transition-colors duration-300">Выберите файл</label>
            <input
                ref="fileInput"
                type="file"
                accept=".json,.csv"
                @change="handleFileSelect"
                class="w-full px-3 py-2 border border-gray-300 dark:border-dark-border rounded-lg focus:outline-none focus:ring-2 focus:ring-forest-green dark:focus:ring-forest-green-dark bg-white dark:bg-dark-bg text-gray-900 dark:text-dark-text transition-colors duration-300"
            />
            <p v-if="selectedFile" class="mt-2 text-sm text-gray-600 dark:text-dark-text-secondary transition-colors duration-300">
              Выбран: {{ selectedFile.name }}
            </p>
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
                type="button"
                @click="handleUploadFile"
                :disabled="loading || !selectedFile"
                class="px-4 py-2 text-sm bg-forest-green dark:bg-forest-green-dark text-white hover:bg-forest-dark dark:hover:bg-forest-green rounded-lg disabled:opacity-50 transition-colors duration-300"
            >
              {{ loading ? 'Загрузка...' : 'Загрузить' }}
            </button>
          </div>
        </div>
      </TabPanel>

      <TabPanel header="📋 Вставить JSON" value="json">
        <div class="space-y-4">
          <div class="bg-blue-50 dark:bg-blue-900/30 border border-blue-200 dark:border-blue-800 rounded-lg p-4 transition-colors duration-300">
            <h4 class="font-semibold text-blue-900 dark:text-blue-200 mb-2 transition-colors duration-300">Формат JSON:</h4>
            <pre class="text-xs bg-white dark:bg-dark-bg p-3 rounded border dark:border-dark-border overflow-x-auto text-gray-900 dark:text-dark-text transition-colors duration-300">{{ jsonExample }}</pre>
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-dark-text-secondary mb-2 transition-colors duration-300">Вставьте JSON</label>
            <textarea
                v-model="jsonInput"
                rows="12"
                placeholder='{"questions": [{"number": 1, "section": "Основы", "question": "Вопрос 1"}]}'
                class="w-full px-3 py-2 border border-gray-300 dark:border-dark-border rounded-lg focus:outline-none focus:ring-2 focus:ring-forest-green dark:focus:ring-forest-green-dark font-mono text-sm bg-white dark:bg-dark-bg text-gray-900 dark:text-dark-text transition-colors duration-300"
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
                type="button"
                @click="handleImportJSON"
                :disabled="loading || !jsonInput.trim()"
                class="px-4 py-2 text-sm bg-forest-green dark:bg-forest-green-dark text-white hover:bg-forest-dark dark:hover:bg-forest-green rounded-lg disabled:opacity-50 transition-colors duration-300"
            >
              {{ loading ? 'Импорт...' : 'Импортировать' }}
            </button>
          </div>
        </div>
      </TabPanel>
    </TabView>
  </Dialog>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import Dialog from 'primevue/dialog'
import TabView from 'primevue/tabview'
import TabPanel from 'primevue/tabpanel'
import { examQuestionsApi } from '@/api/examquestions'

const props = defineProps<{
  modelValue: boolean
  courseId: number
}>()

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  'saved': []
}>()

const visible = ref(props.modelValue)
const fileInput = ref<HTMLInputElement>()
const selectedFile = ref<File | null>(null)
const jsonInput = ref('')
const error = ref('')
const loading = ref(false)

const jsonExample = `{
  "questions": [
    {"number": 1, "section": "Основы", "question": "Вопрос 1"},
    {"number": 2, "section": "Основы", "question": "Вопрос 2"}
  ]
}`

const csvExample = `number,section,question
1,Основы,Вопрос 1
2,Основы,Вопрос 2`

watch(() => props.modelValue, (val) => {
  visible.value = val
  if (!val) {
    selectedFile.value = null
    jsonInput.value = ''
    error.value = ''
    if (fileInput.value) {
      fileInput.value.value = ''
    }
  }
})

watch(visible, (val) => {
  emit('update:modelValue', val)
})

const handleFileSelect = (event: Event) => {
  const target = event.target as HTMLInputElement
  const file = target.files?.[0]
  if (file) {
    selectedFile.value = file
    error.value = ''
  }
}

const handleUploadFile = async () => {
  if (!selectedFile.value) return

  error.value = ''
  loading.value = true

  try {
    await examQuestionsApi.uploadFile(props.courseId, selectedFile.value)
    visible.value = false
    emit('saved')
  } catch (err: any) {
    error.value = err.response?.data?.error || 'Ошибка загрузки файла'
  } finally {
    loading.value = false
  }
}

const handleImportJSON = async () => {
  if (!jsonInput.value.trim()) return

  error.value = ''
  loading.value = true

  try {
    const data = JSON.parse(jsonInput.value)
    
    if (!data.questions || !Array.isArray(data.questions)) {
      error.value = 'JSON должен содержать массив "questions"'
      return
    }

    if (data.questions.length === 0) {
      error.value = 'Массив questions не может быть пустым'
      return
    }

    await examQuestionsApi.bulkCreateJSON({
      course_id: props.courseId,
      questions: data.questions,
    })

    visible.value = false
    emit('saved')
  } catch (err: any) {
    if (err instanceof SyntaxError) {
      error.value = 'Неверный формат JSON: ' + err.message
    } else {
      error.value = err.response?.data?.error || 'Ошибка импорта'
    }
  } finally {
    loading.value = false
  }
}

const onHide = () => {
  selectedFile.value = null
  jsonInput.value = ''
  error.value = ''
  if (fileInput.value) {
    fileInput.value.value = ''
  }
}
</script>
