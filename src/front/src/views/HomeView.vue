<template>
  <div class="min-h-screen p-6">
    <header class="mb-8">
      <h1 class="text-4xl font-bold text-forest-dark flex items-center gap-2">
        🐸 Laritmo
        <span class="text-sm font-normal text-gray-600">Forest Academy</span>
      </h1>
    </header>

    <main>
      <h2 class="text-2xl font-semibold mb-6 flex items-center gap-2">
        🍄 Мои курсы
      </h2>

      <button
          v-if="authStore.isAdmin"
          @click="showCreateDialog = true"
          class="flex gap-6 mb-8 px-4 py-2 bg-forest-green text-white rounded-lg hover:bg-forest-dark transition-colors"
      >
        ➕ Добавить курс
      </button>

      <div v-if="loading" class="text-center py-12">
        <p class="text-gray-600">Загрузка курсов...</p>
      </div>

      <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <div
            v-for="course in courses"
            :key="course.id"
            @click="goToCourse(course.id)"
            class="bg-white rounded-xl shadow-md p-6 hover:shadow-lg transition-shadow cursor-pointer border-2 border-transparent hover:border-forest-mint"
        >
          <div class="text-4xl mb-4">🍄</div>
          <h3 class="text-xl font-semibold mb-2 text-forest-dark">
            {{ course.name }}
          </h3>
          <p class="text-sm text-gray-600 mb-4">{{ course.semester }}</p>
          <p class="text-sm text-gray-700">{{ course.description }}</p>
        </div>
      </div>
    </main>
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