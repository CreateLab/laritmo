<template>
  <div class="min-h-screen p-6">
    <header class="mb-8">
      <h1 class="text-4xl font-bold text-forest-dark dark:text-dark-text flex items-center gap-2 transition-colors duration-300">
        <span class="frog-animation">🐸</span> Laritmo
        <span class="text-sm font-normal text-gray-600 dark:text-dark-text-secondary transition-colors duration-300">Forest Academy</span>
      </h1>
    </header>

    <main>
      <h2 class="text-2xl font-semibold mb-6 flex items-center gap-2 text-forest-dark dark:text-dark-text transition-colors duration-300">
        🍄 Мои курсы
      </h2>

      <button
          v-if="authStore.isAdmin"
          @click="showCreateDialog = true"
          class="flex gap-6 mb-8 px-4 py-2 bg-forest-green dark:bg-forest-green-dark text-white rounded-lg hover:bg-forest-dark dark:hover:bg-forest-green transition-colors duration-300"
      >
        ➕ Добавить курс
      </button>

      <div v-if="loading" class="text-center py-12">
        <p class="text-gray-600 dark:text-dark-text-secondary transition-colors duration-300">Загрузка курсов...</p>
      </div>

      <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <div
            v-for="course in courses"
            :key="course.id"
            @click="goToCourse(course.id)"
            class="bg-white dark:bg-dark-surface rounded-xl shadow-md dark:shadow-lg p-6 hover:shadow-lg dark:hover:shadow-xl transition-all duration-300 cursor-pointer border-2 border-transparent hover:border-forest-mint dark:hover:border-forest-mint-dark"
        >
          <div class="text-4xl mb-4">🍄</div>
          <h3 class="text-xl font-semibold mb-2 text-forest-dark dark:text-dark-text transition-colors duration-300">
            {{ course.name }}
          </h3>
          <p class="text-sm text-gray-600 dark:text-dark-text-secondary mb-4 transition-colors duration-300">{{ course.semester }}</p>
          <p class="text-sm text-gray-700 dark:text-dark-text-secondary transition-colors duration-300">{{ course.description }}</p>
        </div>
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