<template>
  <div class="flex-1">
    <main class="page py-12">
      <div class="flex flex-wrap items-center justify-between gap-4 mb-8">
        <h1 class="text-h1">Мои курсы</h1>
        <AppButton v-if="authStore.isAdmin" variant="primary" icon="pi-plus" @click="showCreateDialog = true">
          Добавить курс
        </AppButton>
      </div>

      <div v-if="loading" class="text-center py-12">
        <p class="text-ink-muted">Загрузка курсов...</p>
      </div>

      <EmptyState v-else-if="loadError" title="Не удалось загрузить курсы" text="Сервер недоступен или вернул неожиданный ответ. Проверьте, что API запущен.">
        <AppButton variant="primary" icon="pi-refresh" @click="handleCourseSaved">Попробовать снова</AppButton>
      </EmptyState>

      <EmptyState v-else-if="courses.length === 0" title="Курсов пока нет" :text="authStore.isAdmin ? 'Создайте первый курс кнопкой «Добавить курс».' : 'Когда курсы появятся, они будут здесь.'" />

      <div v-else class="grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-3 2xl:grid-cols-4 gap-6">
        <a
            v-for="course in courses"
            :key="course.id"
            :href="`/courses/${course.id}`"
            @click.prevent="goToCourse(course.id)"
            class="card no-underline text-ink flex flex-col min-h-[11rem]"
        >
          <div class="flex items-center gap-2">
            <BrandMark />
            <AppChip>{{ course.semester }}</AppChip>
          </div>
          <h3 class="text-h3 text-xl mt-3">{{ course.name }}</h3>
          <p class="text-body text-ink-muted mt-2 line-clamp-3">{{ course.description }}</p>
          <span class="mt-auto pt-4 text-label text-action flex items-center gap-1">
            Открыть <i class="pi pi-arrow-right" aria-hidden="true"></i>
          </span>
        </a>
      </div>
    </main>

    <!-- Floating кнопка переключения темы для неавторизованных пользователей -->
    <div v-if="!authStore.isAuthenticated" class="fixed bottom-6 right-6 z-50">
      <ThemeToggle />
    </div>
    
    <CourseEditDialog
        v-model="showCreateDialog"
        @saved="handleCourseSaved"
    />
  </div>
</template>

<script setup lang="ts">
import EmptyState from '@/components/ui/EmptyState.vue'
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { coursesApi, type Course } from '@/api/courses'
import { useAuthStore } from '@/stores/auth'
import CourseEditDialog from '@/components/CourseEditDialog.vue'
import ThemeToggle from '@/components/ThemeToggle.vue'
import AppButton from '@/components/ui/AppButton.vue'
import AppChip from '@/components/ui/AppChip.vue'
import BrandMark from '@/components/ui/BrandMark.vue'


const router = useRouter()
const courses = ref<Course[]>([])
const loading = ref(true)
const loadError = ref(false)
const authStore = useAuthStore()
const showCreateDialog = ref(false)

const goToCourse = (courseId: number) => {
  router.push({
    path: `/courses/${courseId}`,
  })
}

const loadCourses = async () => {
  loading.value = true
  loadError.value = false
  try {
    const { data } = await coursesApi.getAll()
    if (!Array.isArray(data)) throw new Error('Unexpected response')
    courses.value = data
  } catch (error) {
    console.error('Failed to load courses:', error)
    courses.value = []
    loadError.value = true
  } finally {
    loading.value = false
  }
}

const handleCourseSaved = loadCourses

onMounted(loadCourses)
</script>