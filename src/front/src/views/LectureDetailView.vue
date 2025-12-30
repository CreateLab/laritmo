<template>
  <div class="min-h-screen p-6">
    <header class="mb-8">
      <button
          @click="goBack()"
          class="text-forest-green hover:text-forest-dark mb-4 flex items-center gap-2"
      >
        ← Назад к курсу
      </button>

      <div v-if="loading" class="animate-pulse">
        <div class="h-8 bg-gray-200 rounded w-3/4 mb-4"></div>
        <div class="h-4 bg-gray-200 rounded w-1/4"></div>
      </div>

      <div v-else-if="lecture">
        <div class="flex items-center gap-2 mb-2">
          <span class="px-3 py-1 bg-forest-green text-white rounded-full text-sm">
            Неделя {{ lecture.week }}
          </span>
        </div>
        <h1 class="text-3xl font-bold text-forest-dark mb-2">{{ lecture.title }}</h1>
        <a
            v-if="lecture.github_url"
            :href="lecture.github_url"
            target="_blank"
            class="text-sm text-forest-green hover:underline flex items-center gap-1"
        >
          <i class="pi pi-github"></i>
          Открыть на GitHub
        </a>
      </div>

      <div v-if="authStore.isAdmin" class="flex gap-4">
        <button
            @click="editLecture"
            class="px-4 py-2 bg-forest-green text-white rounded-lg hover:bg-forest-dark transition-colors"
        >
          ✏️ Редактировать
        </button>
        <button
            @click="deleteLecture"
            class="px-4 py-2 bg-red-500 text-white rounded-lg hover:bg-red-600 transition-colors"
        >
          🗑️ Удалить
        </button>
      </div>
    </header>

    <main>
      <div v-if="loading" class="text-center py-12">
        <p class="text-gray-600">Загрузка лекции...</p>
      </div>

      <div v-else-if="!lecture" class="text-center py-12">
        <p class="text-red-600">Лекция не найдена</p>
      </div>

      <div
          v-else
          class="bg-white rounded-xl shadow-md p-8 prose prose-slate max-w-none markdown-content"
          v-html="renderedContent"
      ></div>
    </main>

    <LectureEditDialog
        v-model="showEditDialog"
        :lecture="lecture"
        :course-id="Number(courseId)"
        @saved="handleLectureSaved"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { lecturesApi, type Lecture } from '@/api/lectures'
import { marked, Renderer } from 'marked'
import hljs from 'highlight.js'
import 'highlight.js/styles/github.css'

import { useAuthStore } from '@/stores/auth'
import LectureEditDialog from '@/components/LectureEditDialog.vue'

const authStore = useAuthStore()
const showEditDialog = ref(false)
const route = useRoute()
const router = useRouter()
const lecture = ref<Lecture | null>(null)
const loading = ref(true)
const lectureId = Number(route.params.id)

const courseId = route.params.courseId

const goBack = () => {
  router.push(`/courses/${courseId}`)
}

const editLecture = () => {
  showEditDialog.value = true
}

const deleteLecture = async () => {
  if (!confirm('Удалить лекцию?')) return

  try {
    await lecturesApi.delete(lectureId)
    router.push(`/courses/${courseId}`)
  } catch (error) {
    console.error('Failed to delete:', error)
    alert('Ошибка удаления лекции')
  }
}

const handleLectureSaved = async () => {
  try {
    const { data } = await lecturesApi.getById(lectureId)
    lecture.value = data
  } catch (error) {
    console.error('Failed to load:', error)
  }
}

const renderer = new Renderer()
renderer.code = function({ text, lang }: { text: string; lang?: string }) {
  const language = lang || 'plaintext'
  if (language && hljs.getLanguage(language)) {
    return `<pre><code class="hljs language-${language}">${hljs.highlight(text, { language }).value}</code></pre>`
  }
  return `<pre><code class="hljs">${hljs.highlightAuto(text).value}</code></pre>`
}

marked.use({
  breaks: true,
  gfm: true,
  renderer: renderer,
})

const renderedContent = computed(() => {
  if (!lecture.value) return ''
  return marked(lecture.value.content)
})

onMounted(async () => {
  try {
    const id = Number(route.params.id)
    const { data } = await lecturesApi.getById(id)
    lecture.value = data
  } catch (error) {
    console.error('Failed to load lecture:', error)
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.markdown-content {
  line-height: 1.7;
}

.markdown-content :deep(h1) {
  @apply text-3xl font-bold text-forest-dark mt-8 mb-4;
}

.markdown-content :deep(h2) {
  @apply text-2xl font-semibold text-forest-dark mt-6 mb-3;
}

.markdown-content :deep(h3) {
  @apply text-xl font-semibold text-forest-dark mt-4 mb-2;
}

.markdown-content :deep(p) {
  @apply mb-4 text-gray-700;
}

.markdown-content :deep(ul),
.markdown-content :deep(ol) {
  @apply mb-4 ml-6;
}

.markdown-content :deep(li) {
  @apply mb-2;
}

.markdown-content :deep(code) {
  @apply bg-gray-100 px-2 py-1 rounded text-sm font-mono;
}

.markdown-content :deep(pre) {
  @apply bg-gray-900 text-gray-100 p-4 rounded-lg overflow-x-auto mb-4;
}

.markdown-content :deep(pre code) {
  @apply bg-transparent p-0;
}

.markdown-content :deep(a) {
  @apply text-forest-green hover:underline;
}

.markdown-content :deep(blockquote) {
  @apply border-l-4 border-forest-green pl-4 italic my-4 text-gray-600;
}

.markdown-content :deep(table) {
  @apply w-full mb-4 border-collapse;
}

.markdown-content :deep(th) {
  @apply bg-forest-green text-white p-2 text-left;
}

.markdown-content :deep(td) {
  @apply border border-gray-300 p-2;
}
</style>