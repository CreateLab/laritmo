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

      <nav v-if="allLabs.length > 1" class="mb-6 mx-auto flex items-center gap-2" style="max-width: 760px">
        <AppButton @click="goToLab(prevLab)" :disabled="!prevLab">← Предыдущая</AppButton>
        <select
            :value="labId"
            @change="goToLab(allLabs.find(l => l.id === Number(($event.target as HTMLSelectElement).value)) ?? null)"
            class="field flex-1 min-w-0"
        >
          <option v-for="l in allLabs" :key="l.id" :value="l.id">
            Лаба #{{ l.number }} — {{ l.title }}
          </option>
        </select>
        <AppButton @click="goToLab(nextLab)" :disabled="!nextLab">Следующая →</AppButton>
      </nav>
    </header>

    <main class="mx-auto" style="max-width: 72ch">
      <div v-if="loading" class="text-center py-12">
        <p class="text-ink-muted">Загрузка лабораторной...</p>
      </div>

      <EmptyState v-else-if="!lab" title="Лабораторная работа не найдена" text="Возможно, она была удалена или ссылка устарела." />

      <div v-else class="bg-surface border border-border rounded-lg p-8">
        <div class="flex flex-wrap items-center gap-2">
          <AppChip>Лаба #{{ lab.number }}</AppChip>
          <AppChip>Макс. {{ lab.max_score }} баллов</AppChip>
          <AppChip v-if="lab.deadline">Дедлайн: {{ formatDate(lab.deadline) }}</AppChip>
        </div>
        <h1 class="text-h1 mt-3 mb-4">{{ lab.title }}</h1>
        <div class="flex flex-wrap gap-2 mb-6">
          <a v-if="lab.github_url" :href="lab.github_url" target="_blank" rel="noopener" class="btn btn-secondary">
            <i class="pi pi-github" aria-hidden="true"></i>
            Открыть на GitHub
          </a>
          <template v-if="authStore.isAdmin">
            <AppButton icon="pi-pencil" @click="editLab">Редактировать</AppButton>
            <AppButton variant="danger" icon="pi-trash" @click="deleteLab">Удалить</AppButton>
          </template>
        </div>
        <div class="lecture-prose" v-html="renderedContent"></div>
      </div>
    </main>

    <LabEditDialog
        v-model="showEditDialog"
        :lab="lab"
        :course-id="Number(courseId)"
        @saved="handleLabSaved"
    />
  </div>
</template>

<script setup lang="ts">
import EmptyState from '@/components/ui/EmptyState.vue'
import { ref, onMounted, computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { labsApi, type Lab } from '@/api/labs'
import { marked, Renderer } from 'marked'
import hljs from 'highlight.js'
import { useAuthStore } from '@/stores/auth'
import LabEditDialog from '@/components/LabEditDialog.vue'
import AppButton from '@/components/ui/AppButton.vue'
import AppChip from '@/components/ui/AppChip.vue'

const route = useRoute()
const router = useRouter()
const lab = ref<Lab | null>(null)
const loading = ref(true)
const authStore = useAuthStore()
const showEditDialog = ref(false)

const courseId = route.params.courseId
const labId = computed(() => Number(route.params.id))
const allLabs = ref<Lab[]>([])

const currentIndex = computed(() => allLabs.value.findIndex(l => l.id === labId.value))
const prevLab = computed(() => (currentIndex.value > 0 ? allLabs.value[currentIndex.value - 1] : null) ?? null)
const nextLab = computed(() => (currentIndex.value >= 0 ? allLabs.value[currentIndex.value + 1] : null) ?? null)

const goToLab = (target: Lab | null) => {
  if (!target) return
  router.push(`/courses/${courseId}/labs/${target.id}`)
}

const goBack = () => {
  router.push(`/courses/${courseId}`)
}

const editLab = () => {
  showEditDialog.value = true
}

const deleteLab = async () => {
  if (!confirm('Удалить лабораторную работу?')) return

  try {
    await labsApi.delete(Number(route.params.id))
    router.push(`/courses/${courseId}`)
  } catch (error) {
    console.error('Failed to delete:', error)
    alert('Ошибка удаления лабораторной работы')
  }
}

const handleLabSaved = async () => {
  try {
    const { data } = await labsApi.getById(Number(route.params.id))
    lab.value = data
  } catch (error) {
    console.error('Failed to load:', error)
  }
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

const loadLab = async () => {
  loading.value = true
  try {
    const { data } = await labsApi.getById(labId.value)
    lab.value = data
  } catch (error) {
    lab.value = null
    console.error('Failed to load lab:', error)
  } finally {
    loading.value = false
  }
}

const loadLabList = async () => {
  try {
    const { data } = await labsApi.getAll(Number(courseId))
    allLabs.value = [...data].sort((a, b) => a.number - b.number)
  } catch (error) {
    console.error('Failed to load lab list:', error)
  }
}

onMounted(() => {
  loadLab()
  loadLabList()
})

watch(labId, () => {
  loadLab()
  window.scrollTo({ top: 0 })
})
</script>
