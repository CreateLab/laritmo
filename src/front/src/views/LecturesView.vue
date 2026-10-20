<template>
  <div class="min-h-screen max-w-7xl mx-auto px-4 sm:px-8 py-12">
    <header class="mb-8">
      <a
          href="#"
          @click.prevent="router.back()"
          class="text-action hover:text-action-hover text-sm no-underline inline-block mb-4"
      >
        ← Назад
      </a>
      <h1 class="text-h1">Лекции</h1>
      <p v-if="!loading && lectures.length > 0" class="text-ink-muted mt-2">
        Курс: {{ lectures[0]?.course_id ?? 'N/A' }}
      </p>
    </header>

    <main>
      <div v-if="loading" class="text-center py-12">
        <p class="text-ink-muted">Загрузка лекций...</p>
      </div>

      <EmptyState v-else-if="lectures.length === 0" title="Лекции пока не добавлены" text="Когда преподаватель добавит материалы, они появятся здесь." />

      <div v-else class="space-y-3">
        <a
            v-for="lecture in lectures"
            :key="lecture.id"
            :href="`/lectures/${lecture.id}`"
            @click.prevent="router.push(`/lectures/${lecture.id}`)"
            class="row-item no-underline text-ink"
        >
          <AppChip>Неделя {{ lecture.week }}</AppChip>
          <div class="min-w-0">
            <div class="text-body font-medium">{{ lecture.title }}</div>
            <p class="text-body-sm text-ink-muted">{{ lecture.content.substring(0, 150) }}...</p>
          </div>
        </a>
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import EmptyState from '@/components/ui/EmptyState.vue'
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { lecturesApi, type Lecture } from '@/api/lectures'
import AppChip from '@/components/ui/AppChip.vue'

const route = useRoute()
const lectures = ref<Lecture[]>([])
const loading = ref(true)
const router = useRouter()


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
