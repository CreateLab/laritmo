<template>
  <Dialog
      v-model:visible="visible"
      modal
      :header="lab ? 'Редактировать лабораторную работу' : 'Создать лабораторную работу'"
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

      <div class="grid grid-cols-3 gap-4">
        <div>
          <label class="field-label">Номер</label>
          <input
              v-model.number="form.number"
              type="number"
              required
              min="1"
              class="field"
          />
        </div>

        <div>
          <label class="field-label">Максимальный балл</label>
          <input
              v-model.number="form.max_score"
              type="number"
              required
              min="1"
              class="field"
          />
        </div>

        <div>
          <label class="field-label">Дедлайн (опционально)</label>
          <input
              v-model="form.deadline"
              type="date"
              class="field"
          />
        </div>
      </div>

      <div>
        <label class="field-label">Название лабораторной</label>
        <input
            v-model="form.title"
            type="text"
            required
            class="field"
        />
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
        <label class="field-label">Описание (Markdown)</label>
        <textarea
            ref="descriptionTextarea"
            v-model="form.description"
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
}
</style>
