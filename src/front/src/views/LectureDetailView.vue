<template>
  <div class="min-h-screen p-6">
    <header class="mb-8">
      <button
          @click="goBack()"
          class="text-forest-green dark:text-forest-green-dark hover:text-forest-dark dark:hover:text-forest-green mb-4 flex items-center gap-2 transition-colors duration-300"
      >
        ← Назад к курсу
      </button>

      <nav v-if="allLectures.length > 1" class="mb-4 flex items-center gap-2">
        <button
            @click="goToLecture(prevLecture)"
            :disabled="!prevLecture"
            class="px-3 py-2 rounded-lg bg-white dark:bg-dark-surface text-forest-green dark:text-forest-green-dark shadow disabled:opacity-40 disabled:cursor-not-allowed transition-colors duration-300"
        >
          ← Предыдущая
        </button>
        <select
            :value="lectureId"
            @change="goToLecture(allLectures.find(l => l.id === Number(($event.target as HTMLSelectElement).value)) ?? null)"
            class="flex-1 min-w-0 px-3 py-2 rounded-lg bg-white dark:bg-dark-surface text-forest-dark dark:text-dark-text shadow transition-colors duration-300"
        >
          <option v-for="l in allLectures" :key="l.id" :value="l.id">
            Неделя {{ l.week }} — {{ l.title }}
          </option>
        </select>
        <button
            @click="goToLecture(nextLecture)"
            :disabled="!nextLecture"
            class="px-3 py-2 rounded-lg bg-white dark:bg-dark-surface text-forest-green dark:text-forest-green-dark shadow disabled:opacity-40 disabled:cursor-not-allowed transition-colors duration-300"
        >
          Следующая →
        </button>
      </nav>

      <div v-if="loading" class="animate-pulse">
        <div class="h-8 bg-gray-200 dark:bg-dark-surface rounded w-3/4 mb-4"></div>
        <div class="h-4 bg-gray-200 dark:bg-dark-surface rounded w-1/4"></div>
      </div>

      <div v-else-if="lecture">
        <div class="flex items-center gap-2 mb-2">
          <span class="px-3 py-1 bg-forest-green dark:bg-forest-green-dark text-white rounded-full text-sm transition-colors duration-300">
            Неделя {{ lecture.week }}
          </span>
        </div>
        <h1 class="text-3xl font-bold text-forest-dark dark:text-dark-text mb-2 transition-colors duration-300">{{ lecture.title }}</h1>
        <a
            v-if="lecture.github_url"
            :href="lecture.github_url"
            target="_blank"
            class="text-sm text-forest-green dark:text-forest-green-dark hover:underline flex items-center gap-1 transition-colors duration-300"
        >
          <i class="pi pi-github"></i>
          Открыть на GitHub
        </a>
      </div>

      <div v-if="authStore.isAdmin" class="flex gap-4">
        <button
            @click="editLecture"
            class="px-4 py-2 bg-forest-green dark:bg-forest-green-dark text-white rounded-lg hover:bg-forest-dark dark:hover:bg-forest-green transition-colors duration-300"
        >
          ✏️ Редактировать
        </button>
        <button
            @click="deleteLecture"
            class="px-4 py-2 bg-red-500 dark:bg-red-600 text-white rounded-lg hover:bg-red-600 dark:hover:bg-red-700 transition-colors duration-300"
        >
          🗑️ Удалить
        </button>
      </div>
    </header>

    <main>
      <div v-if="loading" class="text-center py-12">
        <p class="text-gray-600 dark:text-dark-text-secondary transition-colors duration-300">Загрузка лекции...</p>
      </div>

      <div v-else-if="!lecture" class="text-center py-12">
        <p class="text-red-600 dark:text-red-400 transition-colors duration-300">Лекция не найдена</p>
      </div>

      <div
          v-else
          class="bg-white dark:bg-dark-surface rounded-xl shadow-md dark:shadow-lg p-8 prose prose-slate max-w-none markdown-content transition-colors duration-300"
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
import { ref, onMounted, computed, watch } from 'vue'
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
const lectureId = computed(() => Number(route.params.id))
const allLectures = ref<Lecture[]>([])

const courseId = route.params.courseId

const currentIndex = computed(() => allLectures.value.findIndex(l => l.id === lectureId.value))
const prevLecture = computed(() => (currentIndex.value > 0 ? allLectures.value[currentIndex.value - 1] : null) ?? null)
const nextLecture = computed(() => (currentIndex.value >= 0 ? allLectures.value[currentIndex.value + 1] : null) ?? null)

const goToLecture = (target: Lecture | null) => {
  if (!target) return
  router.push(`/courses/${courseId}/lectures/${target.id}`)
}

const goBack = () => {
  router.push(`/courses/${courseId}`)
}

const editLecture = () => {
  showEditDialog.value = true
}

const deleteLecture = async () => {
  if (!confirm('Удалить лекцию?')) return

  try {
    await lecturesApi.delete(lectureId.value)
    router.push(`/courses/${courseId}`)
  } catch (error) {
    console.error('Failed to delete:', error)
    alert('Ошибка удаления лекции')
  }
}

const handleLectureSaved = async () => {
  try {
    const { data } = await lecturesApi.getById(lectureId.value)
    lecture.value = data
    await loadLectureList()
  } catch (error) {
    console.error('Failed to load:', error)
  }
}

const loadLectureList = async () => {
  try {
    const { data } = await lecturesApi.getAll(Number(courseId))
    allLectures.value = [...data].sort((a, b) => a.week - b.week || a.id - b.id)
  } catch (error) {
    console.error('Failed to load lecture list:', error)
  }
}

const loadLecture = async () => {
  loading.value = true
  try {
    const { data } = await lecturesApi.getById(lectureId.value)
    lecture.value = data
  } catch (error) {
    lecture.value = null
    console.error('Failed to load lecture:', error)
  } finally {
    loading.value = false
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

onMounted(() => {
  loadLecture()
  loadLectureList()
})

watch(lectureId, () => {
  loadLecture()
  window.scrollTo({ top: 0 })
})
</script>

<style scoped>
.markdown-content {
  line-height: 1.7;
}

.markdown-content :deep(h1) {
  @apply text-3xl font-bold text-forest-dark dark:text-dark-text mt-8 mb-4;
}

.markdown-content :deep(h2) {
  @apply text-2xl font-semibold text-forest-dark dark:text-dark-text mt-6 mb-3;
}

.markdown-content :deep(h3) {
  @apply text-xl font-semibold text-forest-dark dark:text-dark-text mt-4 mb-2;
}

.markdown-content :deep(p) {
  @apply mb-4 text-gray-700 dark:text-dark-text-secondary;
}

.markdown-content :deep(ul),
.markdown-content :deep(ol) {
  @apply mb-4 ml-6;
}

.markdown-content :deep(li) {
  @apply mb-2;
}

.markdown-content :deep(code) {
  @apply bg-gray-100 dark:bg-dark-bg px-2 py-1 rounded text-sm font-mono;
}

.markdown-content :deep(pre) {
  @apply bg-gray-50 dark:bg-dark-bg border border-gray-200 dark:border-dark-border text-gray-900 dark:text-dark-text p-4 rounded-lg overflow-x-auto mb-4 shadow-sm;
}

.markdown-content :deep(pre code) {
  @apply bg-transparent p-0 text-gray-900 dark:text-dark-text;
}

.markdown-content :deep(a) {
  @apply text-forest-green dark:text-forest-green-dark hover:underline;
}

.markdown-content :deep(blockquote) {
  @apply border-l-4 border-forest-green dark:border-forest-green-dark pl-4 italic my-4 text-gray-600 dark:text-dark-text-secondary;
}

.markdown-content :deep(table) {
  @apply w-full mb-4 border-collapse;
}

.markdown-content :deep(th) {
  @apply bg-forest-green dark:bg-forest-green-dark text-white p-2 text-left;
}

.markdown-content :deep(td) {
  @apply border border-gray-300 dark:border-dark-border p-2;
}
</style>