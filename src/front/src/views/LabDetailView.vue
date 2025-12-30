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

      <div v-else-if="lab">
        <div class="flex items-center gap-2 mb-2">
          <span class="px-3 py-1 bg-forest-green text-white rounded-full text-sm">
            Лаба #{{ lab.number }}
          </span>
          <span class="px-3 py-1 bg-gray-200 text-gray-700 rounded-full text-sm">
            Макс. {{ lab.max_score }} баллов
          </span>
          <span v-if="lab.deadline" class="px-3 py-1 bg-red-100 text-red-700 rounded-full text-sm">
            Дедлайн: {{ formatDate(lab.deadline) }}
          </span>
        </div>
        <h1 class="text-3xl font-bold text-forest-dark mb-2">{{ lab.title }}</h1>
        <a
            v-if="lab.github_url"
            :href="lab.github_url"
            target="_blank"
            class="text-sm text-forest-green hover:underline flex items-center gap-1"
        >
          <i class="pi pi-github"></i>
          Открыть на GitHub
        </a>
      </div>
    </header>

    <main>
      <div v-if="loading" class="text-center py-12">
        <p class="text-gray-600">Загрузка лабораторной...</p>
      </div>

      <div v-else-if="!lab" class="text-center py-12">
        <p class="text-red-600">Лабораторная работа не найдена</p>
      </div>

      <div
          v-else
          class="bg-white rounded-xl shadow-md p-8 prose prose-slate max-w-none markdown-content"
          v-html="renderedContent"
      ></div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { labsApi, type Lab } from '@/api/labs'
import { marked, Renderer } from 'marked'
import hljs from 'highlight.js'
import 'highlight.js/styles/github.css'

const route = useRoute()
const router = useRouter()
const lab = ref<Lab | null>(null)
const loading = ref(true)

const courseId = route.params.courseId

const goBack = () => {
  router.push(`/courses/${courseId}`)
}

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleDateString('ru-RU', {
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  })
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
  if (!lab.value) return ''
  return marked(lab.value.description)
})

onMounted(async () => {
  try {
    const id = Number(route.params.id)
    const { data } = await labsApi.getById(id)
    lab.value = data
  } catch (error) {
    console.error('Failed to load lab:', error)
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