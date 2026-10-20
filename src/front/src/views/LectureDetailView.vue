<template>
  <div class="min-h-screen max-w-7xl mx-auto px-4 sm:px-8 py-12">
    <header class="mb-8">
      <a
          href="#"
          @click.prevent="goBack()"
          class="text-action hover:text-action-hover text-sm no-underline inline-block mb-4"
      >
        ← Назад к курсу
      </a>

      <nav v-if="allLectures.length > 1" class="mb-6 mx-auto flex items-center gap-2" style="max-width: 760px">
        <AppButton @click="goToLecture(prevLecture)" :disabled="!prevLecture">← Предыдущая</AppButton>
        <select
            :value="lectureId"
            @change="goToLecture(allLectures.find(l => l.id === Number(($event.target as HTMLSelectElement).value)) ?? null)"
            class="field flex-1 min-w-0"
        >
          <option v-for="l in allLectures" :key="l.id" :value="l.id">
            Неделя {{ l.week }} — {{ l.title }}
          </option>
        </select>
        <AppButton @click="goToLecture(nextLecture)" :disabled="!nextLecture">Следующая →</AppButton>
      </nav>
    </header>

    <main class="mx-auto" style="max-width: 72ch">
      <div v-if="loading" class="text-center py-12">
        <p class="text-ink-muted">Загрузка лекции...</p>
      </div>

      <div v-else-if="!lecture" class="text-center py-12">
        <p class="text-danger-text">Лекция не найдена</p>
      </div>

      <div v-else class="bg-surface border border-border rounded-lg p-8">
        <AppChip>Неделя {{ lecture.week }}</AppChip>
        <h1 class="text-h1 mt-3 mb-4">{{ lecture.title }}</h1>
        <div class="flex flex-wrap gap-2 mb-6">
          <a v-if="lecture.github_url" :href="lecture.github_url" target="_blank" rel="noopener" class="btn btn-secondary">
            <i class="pi pi-github" aria-hidden="true"></i>
            Открыть на GitHub
          </a>
          <template v-if="authStore.isAdmin">
            <AppButton icon="pi-pencil" @click="editLecture">Редактировать</AppButton>
            <AppButton variant="danger" icon="pi-trash" @click="deleteLecture">Удалить</AppButton>
          </template>
        </div>
        <div class="lecture-prose" v-html="renderedContent"></div>
      </div>
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

import { useAuthStore } from '@/stores/auth'
import LectureEditDialog from '@/components/LectureEditDialog.vue'
import AppButton from '@/components/ui/AppButton.vue'
import AppChip from '@/components/ui/AppChip.vue'

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
