<template>
  <Dialog
      v-model:visible="visible"
      modal
      header="Массовая загрузка вопросов"
      :style="{ width: '800px', maxHeight: '90vh' }"
      @hide="onHide"
  >
    <TabView>
      <TabPanel header="Загрузить файл" value="upload">
        <div class="space-y-4">
          <div class="bg-sunken border border-border rounded-lg p-4">
            <h4 class="font-bold text-ink mb-2">Поддерживаемые форматы:</h4>
            <ul class="text-sm text-ink space-y-1 list-disc list-inside">
              <li><strong>JSON</strong> (.json) - структурированный формат</li>
              <li><strong>CSV</strong> (.csv) - табличный формат</li>
            </ul>
          </div>

          <div class="bg-sunken border border-border rounded-lg p-4">
            <h4 class="font-bold text-ink mb-2">Пример JSON:</h4>
            <pre class="text-xs bg-surface p-3 rounded border overflow-x-auto text-ink">{{ jsonExample }}</pre>
          </div>

          <div class="bg-sunken border border-border rounded-lg p-4">
            <h4 class="font-bold text-ink mb-2">Пример CSV:</h4>
            <pre class="text-xs bg-surface p-3 rounded border overflow-x-auto text-ink">{{ csvExample }}</pre>
          </div>

          <div>
            <label class="field-label">Выберите файл</label>
            <input
                ref="fileInput"
                type="file"
                accept=".json,.csv"
                @change="handleFileSelect"
                class="field"
            />
            <p v-if="selectedFile" class="mt-2 text-sm text-ink-muted">
              Выбран: {{ selectedFile.name }}
            </p>
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
                type="button"
                @click="handleUploadFile"
                :disabled="loading || !selectedFile"
                class="btn btn-primary"
            >
              {{ loading ? 'Загрузка...' : 'Загрузить' }}
            </button>
          </div>
        </div>
      </TabPanel>

      <TabPanel header="Вставить JSON" value="json">
        <div class="space-y-4">
          <div class="bg-sunken border border-border rounded-lg p-4">
            <h4 class="font-bold text-ink mb-2">Формат JSON:</h4>
            <pre class="text-xs bg-surface p-3 rounded border overflow-x-auto text-ink">{{ jsonExample }}</pre>
          </div>

          <div>
            <label class="field-label">Вставьте JSON</label>
            <textarea
                v-model="jsonInput"
                rows="12"
                placeholder='{"questions": [{"number": 1, "section": "Основы", "question": "Вопрос 1"}]}'
                class="field font-mono"
            ></textarea>
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
                type="button"
                @click="handleImportJSON"
                :disabled="loading || !jsonInput.trim()"
                class="btn btn-primary"
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
