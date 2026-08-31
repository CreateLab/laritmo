<template>
  <div class="min-h-screen p-6">
    <header class="mb-8">
      <button
          @click="router.back()"
          class="text-forest-green hover:text-forest-dark mb-4 flex items-center gap-2"
      >
        ← Назад
      </button>
      <div class="flex items-center justify-between gap-4 flex-wrap">
        <h1 class="text-3xl font-bold text-forest-dark">📚 Лекции</h1>
        <button
            v-if="authStore.isAuthenticated"
            @click="downloadMarkdown"
            :disabled="exporting"
            class="px-4 py-2 bg-forest-green text-white rounded-lg hover:bg-forest-dark transition-colors disabled:opacity-50"
        >
          {{ exporting ? 'Скачивание...' : '⬇ Скачать в Markdown' }}
        </button>
      </div>
      <p v-if="!loading && lectures.length > 0" class="text-gray-600 mt-2">
        Курс: {{ lectures[0]?.course_id ?? 'N/A' }}
      </p>
    </header>

    <main>
      <div v-if="loading" class="text-center py-12">
        <p class="text-gray-600">Загрузка лекций...</p>
      </div>

      <div v-else-if="lectures.length === 0" class="text-center py-12">
        <p class="text-gray-600">Лекции пока не добавлены</p>
      </div>

      <div v-else class="space-y-4">
        <div
            v-for="lecture in lectures"
            :key="lecture.id"
            @click="router.push(`/lectures/${lecture.id}`)"
            class="bg-white rounded-xl shadow-md p-6 hover:shadow-lg transition-shadow cursor-pointer border-2 border-transparent hover:border-forest-mint"
        >
          <div class="flex items-start gap-4">
            <div class="text-3xl">🍂</div>
            <div class="flex-1">
              <div class="flex items-center gap-2 mb-2">
                <span class="px-3 py-1 bg-forest-green text-white rounded-full text-sm">
                  Неделя {{ lecture.week }}
                </span>
              </div>
              <h3 class="text-xl font-semibold text-forest-dark mb-2">
                {{ lecture.title }}
              </h3>
              <p class="text-sm text-gray-600">
                {{ lecture.content.substring(0, 150) }}...
              </p>
            </div>
          </div>
        </div>
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { lecturesApi, type Lecture } from '@/api/lectures'
import { useAuthStore } from '@/stores/auth'

const route = useRoute()
const lectures = ref<Lecture[]>([])
const loading = ref(true)
const exporting = ref(false)
const router = useRouter()
const authStore = useAuthStore()

const downloadMarkdown = async () => {
  exporting.value = true
  try {
    const courseId = route.query.course_id ? Number(route.query.course_id) : undefined
    const response = await lecturesApi.exportMarkdown(courseId)
    const blob = new Blob([response.data], { type: 'text/markdown;charset=utf-8' })
    const url = URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = courseId ? `lectures-course-${courseId}.md` : 'lectures.md'
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    URL.revokeObjectURL(url)
  } catch (error) {
    console.error('Failed to export lectures:', error)
  } finally {
    exporting.value = false
  }
}


onMounted(async () => {
  try {
    const courseId = route.query.course_id ? Number(route.query.course_id) : undefined
    const { data } = await lecturesApi.getAll(courseId)
    lectures.value = data
  } catch (error) {
    console.error('Failed to load lectures:', error)
  } finally {
    loading.value = false
  }
})
</script>