<template>
  <Dialog
      v-model:visible="visible"
      modal
      :header="lecture ? 'Редактировать лекцию' : 'Создать лекцию'"
      :style="{ width: '900px', maxHeight: '90vh' }"
      @hide="onHide"
  >
    <form @submit.prevent="handleSubmit" class="space-y-4">
      <div class="flex justify-end gap-2 pb-4 border-b border-border">
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

      <div class="grid grid-cols-2 gap-4">
        <div>
          <label class="field-label">Неделя</label>
          <input
              v-model.number="form.week"
              type="number"
              required
              min="1"
              class="field"
          />
        </div>

        <div>
          <label class="field-label">Название лекции</label>
          <input
              v-model="form.title"
              type="text"
              required
              class="field"
          />
        </div>
      </div>

      <div>
        <label class="field-label">GitHub URL (опционально)</label>
        <input
            v-model="form.github_url"
            type="url"
            placeholder="https://github.com/..."
            class="field"
        />
      </div>

      <div>
        <label class="field-label">Контент (Markdown)</label>
        <textarea
            ref="contentTextarea"
            v-model="form.content"
            class="w-full"
        ></textarea>
      </div>

      <div v-if="error" class="alert">
        <i class="pi pi-exclamation-triangle" aria-hidden="true"></i>
        <span>{{ error }}</span>
      </div>
    </form>
  </Dialog>
</template>

<script setup lang="ts">
import { ref, watch, onMounted, nextTick } from 'vue'
import Dialog from 'primevue/dialog'
import { lecturesApi, type Lecture } from '@/api/lectures'
import EasyMDE from 'easymde'
import 'easymde/dist/easymde.min.css'

const props = defineProps<{
  modelValue: boolean
  lecture?: Lecture | null
  courseId: number
}>()

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  'saved': []
}>()

const visible = ref(props.modelValue)
const contentTextarea = ref<HTMLTextAreaElement>()
let editor: EasyMDE | null = null

const form = ref({
  week: 1,
  title: '',
  content: '',
  github_url: '',
})
const error = ref('')
const loading = ref(false)

watch(() => props.modelValue, async (val) => {
  visible.value = val
  if (val) {
    if (props.lecture) {
      form.value = {
        week: props.lecture.week,
        title: props.lecture.title,
        content: props.lecture.content,
        github_url: props.lecture.github_url || '',
      }
    } else {
      form.value = {
        week: 1,
        title: '',
        content: '',
        github_url: '',
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
  if (!contentTextarea.value || editor) return

  editor = new EasyMDE({
    element: contentTextarea.value,
    spellChecker: false,
    toolbar: [
      'bold', 'italic', 'heading', '|',
      'quote', 'unordered-list', 'ordered-list', '|',
      'link', 'image', 'code', 'table', '|',
      'preview', 'side-by-side', 'fullscreen', '|',
      'guide'
    ],
    placeholder: 'Введите текст лекции в формате Markdown...',
    initialValue: form.value.content,
  })

  editor.codemirror.on('change', () => {
    form.value.content = editor?.value() || ''
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
      week: form.value.week,
      title: form.value.title,
      content: form.value.content,
      github_url: form.value.github_url || undefined,
    }

    if (props.lecture) {
      await lecturesApi.update(props.lecture.id, data)
    } else {
      await lecturesApi.create(data)
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
  form.value = { week: 1, title: '', content: '', github_url: '' }
  error.value = ''
  destroyEditor()
}
</script>

<style>
.CodeMirror {
  height: 400px;
}
</style>
