<template>
  <Dialog
      v-model:visible="visible"
      modal
      :header="lab ? 'Редактировать лабораторную работу' : 'Создать лабораторную работу'"
      :style="{ width: '900px', maxHeight: '90vh' }"
      @hide="onHide"
  >
    <form @submit.prevent="handleSubmit" class="space-y-4">
      <div class="grid grid-cols-3 gap-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Номер</label>
          <input
              v-model.number="form.number"
              type="number"
              required
              min="1"
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-forest-green"
          />
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Максимальный балл</label>
          <input
              v-model.number="form.max_score"
              type="number"
              required
              min="1"
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-forest-green"
          />
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Дедлайн (опционально)</label>
          <input
              v-model="form.deadline"
              type="date"
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-forest-green"
          />
        </div>
      </div>

      <div>
        <label class="block text-sm font-medium text-gray-700 mb-1">Название лабораторной</label>
        <input
            v-model="form.title"
            type="text"
            required
            class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-forest-green"
        />
      </div>

      <div>
        <label class="block text-sm font-medium text-gray-700 mb-1">GitHub URL (опционально)</label>
        <input
            v-model="form.github_url"
            type="url"
            placeholder="https://github.com/..."
            class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-forest-green"
        />
      </div>

      <div>
        <label class="block text-sm font-medium text-gray-700 mb-1">Описание (Markdown)</label>
        <textarea
            ref="descriptionTextarea"
            v-model="form.description"
            class="w-full"
        ></textarea>
      </div>

      <div v-if="error" class="text-red-600 text-sm">
        {{ error }}
      </div>

      <div class="flex justify-end gap-2 pt-4">
        <button
            type="button"
            @click="visible = false"
            class="px-4 py-2 text-sm bg-gray-100 hover:bg-gray-200 rounded-lg"
        >
          Отмена
        </button>
        <button
            type="submit"
            :disabled="loading"
            class="px-4 py-2 text-sm bg-forest-green text-white hover:bg-forest-dark rounded-lg disabled:opacity-50"
        >
          {{ loading ? 'Сохранение...' : 'Сохранить' }}
        </button>
      </div>
    </form>
  </Dialog>
</template>

<script setup lang="ts">
import { ref, watch, nextTick } from 'vue'
import Dialog from 'primevue/dialog'
import { labsApi, type Lab } from '@/api/labs'
import EasyMDE from 'easymde'
import 'easymde/dist/easymde.min.css'

const props = defineProps<{
  modelValue: boolean
  lab?: Lab | null
  courseId: number
}>()

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  'saved': []
}>()

const visible = ref(props.modelValue)
const descriptionTextarea = ref<HTMLTextAreaElement>()
let editor: EasyMDE | null = null

const form = ref({
  number: 1,
  title: '',
  description: '',
  max_score: 10,
  github_url: '',
  deadline: '',
})
const error = ref('')
const loading = ref(false)

watch(() => props.modelValue, async (val) => {
  visible.value = val
  if (val) {
    if (props.lab) {
      form.value = {
        number: props.lab.number,
        title: props.lab.title,
        description: props.lab.description,
        max_score: props.lab.max_score,
        github_url: props.lab.github_url || '',
        deadline: (props.lab.deadline ? props.lab.deadline.split('T')[0] : '') as string,
      }
    } else {
      form.value = {
        number: 1,
        title: '',
        description: '',
        max_score: 10,
        github_url: '',
        deadline: '',
      }
    }

    await nextTick()
    initEditor()
  } else {
    destroyEditor()
  }
})

watch(visible, (val) => {
  emit('update:modelValue', val)
})

const initEditor = () => {
  if (!descriptionTextarea.value || editor) return

  editor = new EasyMDE({
    element: descriptionTextarea.value,
    spellChecker: false,
    toolbar: [
      'bold', 'italic', 'heading', '|',
      'quote', 'unordered-list', 'ordered-list', '|',
      'link', 'image', 'code', 'table', '|',
      'preview', 'side-by-side', 'fullscreen', '|',
      'guide'
    ],
    placeholder: 'Введите описание лабораторной в формате Markdown...',
    initialValue: form.value.description,
  })

  editor.codemirror.on('change', () => {
    form.value.description = editor?.value() || ''
  })
}

const destroyEditor = () => {
  if (editor) {
    editor.toTextArea()
    editor = null
  }
}

const handleSubmit = async () => {
  error.value = ''
  loading.value = true

  try {
    const data = {
      course_id: props.courseId,
      number: form.value.number,
      title: form.value.title,
      description: form.value.description,
      max_score: form.value.max_score,
      github_url: form.value.github_url || undefined,
      deadline: form.value.deadline || undefined,
    }

    if (props.lab) {
      await labsApi.update(props.lab.id, data)
    } else {
      await labsApi.create(data)
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
  form.value = { number: 1, title: '', description: '', max_score: 10, github_url: '', deadline: '' }
  error.value = ''
  destroyEditor()
}
</script>

<style>
.CodeMirror {
  height: 400px;
  border: 1px solid #d1d5db;
  border-radius: 0.5rem;
}

.editor-toolbar {
  border: 1px solid #d1d5db;
  border-bottom: none;
  border-radius: 0.5rem 0.5rem 0 0;
}
</style>
