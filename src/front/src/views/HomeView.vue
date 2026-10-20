<template>
  <div class="min-h-screen">
    <main class="max-w-7xl mx-auto px-4 sm:px-8 py-12">
      <div class="flex flex-wrap items-center justify-between gap-4 mb-8">
        <h1 class="text-h1">Мои курсы</h1>
        <AppButton v-if="authStore.isAdmin" variant="primary" icon="pi-plus" @click="showCreateDialog = true">
          Добавить курс
        </AppButton>
      </div>

      <div v-if="loading" class="text-center py-12">
        <p class="text-ink-muted">Загрузка курсов...</p>
      </div>

      <div v-else class="grid gap-6" style="grid-template-columns: repeat(auto-fill, minmax(min(300px, 100%), 1fr))">
        <a
            v-for="course in courses"
            :key="course.id"
            :href="`/courses/${course.id}`"
            @click.prevent="goToCourse(course.id)"
            class="card no-underline text-ink"
        >
          <div class="flex items-center gap-2">
            <BrandMark />
            <AppChip>{{ course.semester }}</AppChip>
          </div>
          <h3 class="text-h3 mt-3">{{ course.name }}</h3>
          <p class="text-body-sm text-ink-muted">{{ course.description }}</p>
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
const authStore = useAuthStore()
const showCreateDialog = ref(false)

const goToCourse = (courseId: number) => {
  router.push({
    path: `/courses/${courseId}`,
  })
}

const handleCourseSaved = async () => {
  loading.value = true
  try {
    const { data } = await coursesApi.getAll()
    courses.value = data
  } catch (error) {
    console.error('Failed to load courses:', error)
  } finally {
    loading.value = false
  }
}

onMounted(async () => {
  try {
    const { data } = await coursesApi.getAll()
    courses.value = data
  } catch (error) {
    console.error('Failed to load courses:', error)
  } finally {
    loading.value = false
  }
})
</script>