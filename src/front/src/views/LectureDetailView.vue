<template>
  <div class="page px-4 sm:px-6 md:px-10 py-4 sm:py-12">
    <div class="grid gap-8 items-start xl:grid-cols-[14rem_minmax(0,72ch)_18rem] xl:justify-center">
      <a
          href="#"
          @click.prevent="goBack()"
          class="text-action hover:text-action-hover text-body-sm no-underline inline-flex items-center gap-2 min-h-11 xl:col-span-3 justify-self-start"
      >
        <i class="pi pi-arrow-left" aria-hidden="true"></i>
        Назад к курсу
      </a>

      <aside class="hidden xl:block sticky top-24 self-start" aria-label="Оглавление">
        <nav v-if="tocItems.length >= 3" class="border-l border-border">
          <p class="text-label text-ink-muted pl-4 mb-2">Содержание</p>
          <a
              v-for="item in tocItems"
              :key="item.id"
              :href="`#${item.id}`"
              @click.prevent="scrollToSection(item.id)"
              :class="[
                'block py-1 text-body-sm no-underline -ml-px border-l-2 transition-colors duration-150',
                item.level === 3 ? 'pl-8' : 'pl-4',
                item.id === activeTocId ? 'border-action text-ink font-medium' : 'border-transparent text-ink-muted hover:text-ink'
              ]"
          >{{ item.text }}</a>
        </nav>
      </aside>

      <main class="min-w-0 mx-auto w-full" style="max-width: 72ch">
        <div v-if="loading" class="text-center py-12">
          <p class="text-ink-muted">Загрузка лекции...</p>
        </div>

        <EmptyState v-else-if="!lecture" title="Лекция не найдена" text="Возможно, она была удалена или ссылка устарела." />

        <article v-else class="bg-surface border border-border border-x-0 rounded-none -mx-4 sm:mx-0 sm:border-x sm:rounded-lg p-4 sm:p-6 xl:p-8">
          <h1 class="lecture-title mb-3 sm:mb-4">{{ lecture.title }}</h1>
          <div class="xl:hidden flex flex-wrap items-center gap-x-4 gap-y-1 mb-4 sm:mb-6">
            <a v-if="lecture.github_url" :href="lecture.github_url" target="_blank" rel="noopener" class="inline-flex items-center gap-2 min-h-11 text-body-sm text-action hover:text-action-hover no-underline">
              <i class="pi pi-github" aria-hidden="true"></i>
              Открыть на GitHub
            </a>
            <template v-if="authStore.isAdmin">
              <AppButton class="min-h-11" icon="pi-pencil" @click="editLecture">Редактировать</AppButton>
              <AppButton class="min-h-11" variant="danger" icon="pi-trash" @click="deleteLecture">Удалить</AppButton>
            </template>
          </div>
          <div ref="proseRef" class="lecture-prose" v-html="renderedContent"></div>
        </article>
        <nav v-if="allLectures.length > 1" class="xl:hidden grid grid-cols-2 gap-3 mt-4" aria-label="Переключение">
          <button
              v-if="prevLecture"
              type="button"
              class="btn btn-secondary !h-auto min-h-11 py-2 flex-col !items-start !gap-0 min-w-0 text-left"
              @click="goToLecture(prevLecture)"
          >
            <span class="text-label">← Предыдущая</span>
            <span class="text-body-sm text-ink-muted max-w-full truncate" :title="prevLecture.title">{{ prevLecture.title }}</span>
          </button>
          <span v-else></span>
          <button
              v-if="nextLecture"
              type="button"
              class="btn btn-secondary !h-auto min-h-11 py-2 flex-col !items-end !gap-0 min-w-0 text-right"
              @click="goToLecture(nextLecture)"
          >
            <span class="text-label">Следующая →</span>
            <span class="text-body-sm text-ink-muted max-w-full truncate" :title="nextLecture.title">{{ nextLecture.title }}</span>
          </button>
          <span v-else></span>
        </nav>
      </main>

      <aside v-if="lecture" class="hidden xl:flex flex-col gap-3 sticky top-24 self-start">
        <div v-if="allLectures.length > 1" class="flex flex-col gap-2">
          <p class="text-label text-ink-muted">Лекции</p>
          <select
              :value="lectureId"
              @change="goToLecture(allLectures.find(l => l.id === Number(($event.target as HTMLSelectElement).value)) ?? null)"
              class="field w-full min-w-0 truncate" :title="lecture?.title"
          >
            <option v-for="l in allLectures" :key="l.id" :value="l.id" :title="l.title">
              {{ l.title }}
            </option>
          </select>
          <AppButton @click="goToLecture(prevLecture)" :disabled="!prevLecture">← Предыдущая</AppButton>
          <AppButton @click="goToLecture(nextLecture)" :disabled="!nextLecture">Следующая →</AppButton>
        </div>
        <a v-if="lecture.github_url" :href="lecture.github_url" target="_blank" rel="noopener" class="btn btn-secondary">
          <i class="pi pi-github" aria-hidden="true"></i>
          Открыть на GitHub
        </a>
        <template v-if="authStore.isAdmin">
          <AppButton icon="pi-pencil" @click="editLecture">Редактировать</AppButton>
          <AppButton variant="danger" icon="pi-trash" @click="deleteLecture">Удалить</AppButton>
        </template>
      </aside>
    </div>

    <LectureEditDialog
        v-model="showEditDialog"
        :lecture="lecture"
        :course-id="Number(courseId)"
        @saved="handleLectureSaved"
    />
  </div>
</template>

<script setup lang="ts">
import EmptyState from '@/components/ui/EmptyState.vue'
import { ref, onMounted, computed, watch } from 'vue'
import { useToc } from '@/composables/useToc'
import { useRoute, useRouter } from 'vue-router'
import { lecturesApi, type Lecture } from '@/api/lectures'
import { marked, Renderer } from 'marked'
import hljs from 'highlight.js'

import { useAuthStore } from '@/stores/auth'
import LectureEditDialog from '@/components/LectureEditDialog.vue'
import AppButton from '@/components/ui/AppButton.vue'

const authStore = useAuthStore()
const showEditDialog = ref(false)
const route = useRoute()
const router = useRouter()
const lecture = ref<Lecture | null>(null)
const loading = ref(true)
const lectureId = computed(() => Number(route.params.id))
const allLectures = ref<Lecture[]>([])

const courseId = route.params.courseId
const proseRef = ref<HTMLElement | null>(null)

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

const { items: tocItems, activeId: activeTocId, scrollTo: scrollToSection } = useToc(proseRef, () => renderedContent.value)

onMounted(() => {
  loadLecture()
  loadLectureList()
})

watch(lectureId, () => {
  loadLecture()
  window.scrollTo({ top: 0 })
})
</script>

