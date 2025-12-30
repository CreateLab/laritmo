<template>
  <Dialog
      v-model:visible="visible"
      modal
      :header="course ? 'Редактировать курс' : 'Создать курс'"
      :style="{ width: '500px' }"
      @hide="onHide"
  >
    <form @submit.prevent="handleSubmit" class="space-y-4">
      <div>
        <label class="block text-sm font-medium text-gray-700 mb-1">Название курса</label>
        <input
            v-model="form.name"
            type="text"
            required
            class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-forest-green"
        />
      </div>

      <div>
        <label class="block text-sm font-medium text-gray-700 mb-1">Семестр</label>
        <input
            v-model="form.semester"
            type="text"
            required
            placeholder="2024-2025"
            class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-forest-green"
        />
      </div>

      <div>
        <label class="block text-sm font-medium text-gray-700 mb-1">Описание</label>
        <textarea
            v-model="form.description"
            rows="4"
            class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-forest-green"
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
import { ref, watch } from 'vue'
import Dialog from 'primevue/dialog'
import { coursesApi, type Course } from '@/api/courses'

const props = defineProps<{
  modelValue: boolean
  course?: Course | null
}>()

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  'saved': []
}>()

const visible = ref(props.modelValue)
const form = ref({
  name: '',
  semester: '',
  description: '',
})
const error = ref('')
const loading = ref(false)

watch(() => props.modelValue, (val) => {
  visible.value = val
  if (val && props.course) {
    form.value = {
      name: props.course.name,
      semester: props.course.semester,
      description: props.course.description,
    }
  } else if (val) {
    form.value = {
      name: '',
      semester: '',
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
    if (props.course) {
      await coursesApi.update(props.course.id, form.value)
    } else {
      await coursesApi.create(form.value)
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
  form.value = { name: '', semester: '', description: '' }
  error.value = ''
}
</script>