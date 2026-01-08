<template>
  <Dialog
    v-model:visible="visible"
    modal
    :header="isAuthenticated ? 'Генерация билетов' : 'Сгенерировать билет'"
    :style="{ width: '90vw', maxWidth: '800px', maxHeight: '90vh' }"
    @hide="onHide"
    :draggable="false"
    :closable="true"
  >
    <div v-if="loading" class="flex flex-col items-center justify-center py-8">
      <ProgressSpinner />
      <p class="mt-4 text-gray-600 dark:text-dark-text-secondary transition-colors duration-300">Генерация билетов...</p>
    </div>

    <div v-else-if="generatedTicket && !isAuthenticated" class="space-y-4">
      <!-- Отображение билета для неавторизованных -->
      <div class="bg-white dark:bg-dark-surface rounded-lg border border-gray-200 dark:border-dark-border p-6 transition-colors duration-300">
        <h3 class="text-xl font-bold text-forest-dark dark:text-dark-text mb-4 transition-colors duration-300">
          Билет № {{ generatedTicket.number }}
        </h3>
        <div class="space-y-4">
          <div
            v-for="(question, index) in generatedTicket.questions"
            :key="index"
            class="border-l-4 border-forest-green dark:border-forest-green-dark pl-4 py-2 transition-colors duration-300"
          >
            <p class="text-gray-800 dark:text-dark-text font-medium mb-1 transition-colors duration-300">
              {{ index + 1 }}. {{ question.question }}
            </p>
            <p class="text-sm text-gray-600 dark:text-dark-text-secondary transition-colors duration-300">(раздел: {{ question.section }})</p>
          </div>
        </div>
      </div>
      <div class="flex justify-end">
        <Button
          label="Закрыть"
          @click="visible = false"
          outlined
        />
      </div>
    </div>

    <form v-else @submit.prevent="handleSubmit" class="space-y-6">
      <!-- Поле для количества вопросов -->
      <div>
        <label class="block text-sm font-medium text-gray-700 dark:text-dark-text-secondary mb-2 transition-colors duration-300">
          Количество вопросов в билете
        </label>
        <InputNumber
          v-model="form.questionsPerTicket"
          :min="1"
          :max="50"
          :disabled="loading"
          class="w-full"
          required
        />
        <p class="mt-1 text-xs text-gray-500 dark:text-dark-text-secondary transition-colors duration-300">
          От 1 до 50 вопросов
        </p>
      </div>

      <!-- Поле для количества билетов (только для авторизованных) -->
      <div v-if="isAuthenticated">
        <label class="block text-sm font-medium text-gray-700 dark:text-dark-text-secondary mb-2 transition-colors duration-300">
          Количество билетов
        </label>
        <InputNumber
          v-model="form.ticketCount"
          :min="1"
          :max="100"
          :disabled="loading"
          class="w-full"
          required
        />
        <p class="mt-1 text-xs text-gray-500 dark:text-dark-text-secondary transition-colors duration-300">
          От 1 до 100 билетов
        </p>
      </div>

      <div v-if="error" class="bg-red-50 dark:bg-red-900/30 border border-red-200 dark:border-red-800 rounded-lg p-4 transition-colors duration-300">
        <p class="text-red-800 dark:text-red-300 text-sm transition-colors duration-300">{{ error }}</p>
      </div>

      <div class="flex justify-end gap-3 pt-4">
        <Button
          type="button"
          label="Отмена"
          @click="visible = false"
          :disabled="loading"
          outlined
        />
        <Button
          type="submit"
          :label="isAuthenticated ? 'Скачать документ' : 'Сгенерировать билет'"
          :disabled="loading || !isFormValid"
          :loading="loading"
          icon="pi pi-file"
        />
      </div>
    </form>
  </Dialog>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import Dialog from 'primevue/dialog'
import Button from 'primevue/button'
import InputNumber from 'primevue/inputnumber'
import ProgressSpinner from 'primevue/progressspinner'
import { useToast } from 'primevue/usetoast'
import {
  generateRandomTicket,
  generateTicketsDocument,
  downloadBlob,
  type Ticket,
} from '@/api/tickets'

interface Props {
  visible: boolean
  courseId: number
  courseName: string
  isAuthenticated: boolean
}

const props = defineProps<Props>()

const emit = defineEmits<{
  'update:visible': [value: boolean]
}>()

// Используем useToast - ToastService должен быть подключен в main.ts
const toast = useToast()

const visible = ref(props.visible)
const loading = ref(false)
const error = ref('')
const generatedTicket = ref<Ticket | null>(null)

const form = ref({
  questionsPerTicket: 10,
  ticketCount: 20,
})

const isFormValid = computed(() => {
  if (props.isAuthenticated) {
    return (
      form.value.questionsPerTicket >= 1 &&
      form.value.questionsPerTicket <= 50 &&
      form.value.ticketCount >= 1 &&
      form.value.ticketCount <= 100
    )
  }
  return (
    form.value.questionsPerTicket >= 1 &&
    form.value.questionsPerTicket <= 50
  )
})

watch(() => props.visible, (val) => {
  visible.value = val
  if (!val) {
    // Сброс состояния при закрытии
    error.value = ''
    generatedTicket.value = null
    form.value = {
      questionsPerTicket: 10,
      ticketCount: 20,
    }
  }
}, { immediate: true })

watch(visible, (val) => {
  emit('update:visible', val)
})

const handleSubmit = async () => {
  error.value = ''
  loading.value = true

  try {
    if (props.isAuthenticated) {
      // Для авторизованных: генерируем и скачиваем документ
      const blob = await generateTicketsDocument(props.courseId, {
        questionsPerTicket: form.value.questionsPerTicket,
        ticketCount: form.value.ticketCount,
      })

      // Создаем имя файла из названия курса
      const courseSlug = props.courseName
        .toLowerCase()
        .replace(/\s+/g, '_')
        .replace(/\//g, '_')
      const filename = `tickets_${courseSlug}.txt`

      downloadBlob(blob, filename)

      toast.add({
        severity: 'success',
        summary: 'Успешно',
        detail: `Сгенерировано ${form.value.ticketCount} билетов`,
        life: 3000,
      })

      visible.value = false
    } else {
      // Для неавторизованных: генерируем один билет и показываем его
      const ticket = await generateRandomTicket(
        props.courseId,
        form.value.questionsPerTicket
      )
      generatedTicket.value = ticket
    }
  } catch (err: any) {
    console.error('Failed to generate tickets:', err)
    const errorMessage =
      err.response?.data?.error ||
      err.message ||
      'Ошибка при генерации билетов'
    error.value = errorMessage

    toast.add({
      severity: 'error',
      summary: 'Ошибка',
      detail: errorMessage,
      life: 5000,
    })
  } finally {
    loading.value = false
  }
}

const onHide = () => {
  error.value = ''
  generatedTicket.value = null
  form.value = {
    questionsPerTicket: 10,
    ticketCount: 20,
  }
}
</script>
