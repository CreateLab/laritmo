<template>
  <div class="min-h-screen p-6">
    <header class="mb-8">
      <button
          @click="router.push('/')"
          class="text-forest-green dark:text-forest-green-dark hover:text-forest-dark dark:hover:text-forest-green mb-4 flex items-center gap-2 transition-colors duration-300"
      >
        ← На главную
      </button>

      <div v-if="loading" class="animate-pulse">
        <div class="h-8 bg-gray-200 dark:bg-dark-surface rounded w-3/4 mb-2"></div>
        <div class="h-4 bg-gray-200 dark:bg-dark-surface rounded w-1/4"></div>
      </div>

      <div v-else-if="course" class="flex items-start justify-between">
        <div>
          <h1 class="text-3xl font-bold text-forest-dark dark:text-dark-text mb-2 transition-colors duration-300">{{ course.name }}</h1>
          <p class="text-gray-600 dark:text-dark-text-secondary transition-colors duration-300">{{ course.semester }}</p>
          <p class="text-gray-700 dark:text-dark-text-secondary mt-2 transition-colors duration-300">{{ course.description }}</p>
        </div>

        <div v-if="authStore.isAdmin">
          <button
              @click="editCourse"
              class="px-4 py-2 bg-forest-green dark:bg-forest-green-dark text-white rounded-lg hover:bg-forest-dark dark:hover:bg-forest-green transition-colors duration-300"
          >
            ✏️ Редактировать
          </button>
          <button
              @click="deleteCourse"
              class="flex gap-6 mb-8 px-4 py-2 bg-red-500 dark:bg-red-600 text-white rounded-lg hover:bg-red-600 dark:hover:bg-red-700 transition-colors duration-300"
          >
            🗑️ Удалить
          </button>
        </div>
      </div>
    </header>

    <main>
      <TabView v-if="!loading">
        <TabPanel header="📚 Лекции" value="lectures">

          <div class="mb-4 flex gap-2 flex-wrap">
            <button
                v-if="authStore.isAdmin"
                @click="addLecture"
                class="px-4 py-2 bg-forest-green dark:bg-forest-green-dark text-white rounded-lg hover:bg-forest-dark dark:hover:bg-forest-green transition-colors duration-300"
            >
              ➕ Добавить лекцию
            </button>
            <button
                v-if="authStore.isAuthenticated && lectures.length > 0"
                @click="downloadLecturesMarkdown"
                :disabled="exportingLectures"
                class="px-4 py-2 bg-blue-500 dark:bg-blue-600 text-white rounded-lg hover:bg-blue-600 dark:hover:bg-blue-700 transition-colors duration-300 disabled:opacity-50"
            >
              {{ exportingLectures ? 'Скачивание...' : '⬇ Скачать все в Markdown' }}
            </button>
          </div>

          <div v-if="lecturesLoading" class="text-center py-8">
            <p class="text-gray-600 dark:text-dark-text-secondary transition-colors duration-300">Загрузка лекций...</p>
          </div>

          <div v-else-if="lectures.length === 0" class="text-center py-8">
            <p class="text-gray-600 dark:text-dark-text-secondary transition-colors duration-300">Лекции пока не добавлены</p>
          </div>

          <div v-else class="space-y-3">
            <div
                v-for="lecture in lectures"
                :key="lecture.id"
                @click="goToLecture(lecture.id)"
                class="bg-white dark:bg-dark-surface rounded-lg shadow dark:shadow-lg p-4 hover:shadow-md dark:hover:shadow-xl transition-all duration-300 cursor-pointer border-2 border-transparent hover:border-forest-mint dark:hover:border-forest-mint-dark"
            >
              <div class="flex items-start gap-3">
                <div class="text-2xl">🍂</div>
                <div class="flex-1">
                  <div class="flex items-center gap-2 mb-1">
                    <span class="px-2 py-1 bg-forest-green dark:bg-forest-green-dark text-white rounded-full text-xs">
                      Неделя {{ lecture.week }}
                    </span>
                  </div>
                  <h3 class="font-semibold text-forest-dark dark:text-dark-text transition-colors duration-300">{{ lecture.title }}</h3>
                </div>
              </div>
            </div>
          </div>
        </TabPanel>

        <TabPanel header="🔬 Лабораторные" value="labs">
          <div v-if="authStore.isAdmin" class="mb-4">
            <button
                @click="addLab"
                class="px-4 py-2 bg-forest-green dark:bg-forest-green-dark text-white rounded-lg hover:bg-forest-dark dark:hover:bg-forest-green transition-colors duration-300"
            >
              ➕ Добавить лабораторную
            </button>
          </div>

          <div v-if="labsLoading" class="text-center py-8">
            <p class="text-gray-600 dark:text-dark-text-secondary transition-colors duration-300">Загрузка лаб...</p>
          </div>

          <div v-else-if="labs.length === 0" class="text-center py-8">
            <p class="text-gray-600 dark:text-dark-text-secondary transition-colors duration-300">Лабораторные пока не добавлены</p>
          </div>

          <div v-else class="space-y-3">
            <div
                v-for="lab in labs"
                :key="lab.id"
                @click="goToLab(lab.id)"
                class="bg-white dark:bg-dark-surface rounded-lg shadow dark:shadow-lg p-4 hover:shadow-md dark:hover:shadow-xl transition-all duration-300 border-2 border-transparent hover:border-forest-mint dark:hover:border-forest-mint-dark cursor-pointer"
            >
              <div class="flex items-start gap-3">
                <div class="text-2xl">🧪</div>
                <div class="flex-1">
                  <div class="flex items-center gap-2 mb-1">
                    <span class="px-2 py-1 bg-forest-green dark:bg-forest-green-dark text-white rounded-full text-xs">
                      Лаба #{{ lab.number }}
                    </span>
                    <span class="text-xs text-gray-600 dark:text-dark-text-secondary transition-colors duration-300">Макс: {{ lab.max_score }} баллов</span>
                  </div>
                  <h3 class="font-semibold text-forest-dark dark:text-dark-text mb-2 transition-colors duration-300">{{ lab.title }}</h3>
                  <p class="text-sm text-gray-600 dark:text-dark-text-secondary line-clamp-2 transition-colors duration-300">
                    {{ lab.description.substring(0, 150) }}...
                  </p>
                </div>
              </div>
            </div>
          </div>
        </TabPanel>

        <TabPanel header="📊 Журнал" value="grades">
          <div v-if="authStore.isAdmin" class="mb-4">
            <button
                @click="addOrEditGradeSheet"
                class="px-4 py-2 bg-forest-green dark:bg-forest-green-dark text-white rounded-lg hover:bg-forest-dark dark:hover:bg-forest-green transition-colors duration-300"
            >
              {{ gradeSheets.length > 0 ? '✏️ Изменить ссылку на журнал' : '➕ Добавить ссылку на журнал' }}
            </button>
          </div>

          <div v-if="gradeSheetsLoading" class="text-center py-8">
            <p class="text-gray-600 dark:text-dark-text-secondary transition-colors duration-300">Загрузка...</p>
          </div>

          <div v-else-if="!gradeSheets || gradeSheets.length === 0" class="text-center py-8">
            <p class="text-gray-600 dark:text-dark-text-secondary transition-colors duration-300">Журнал пока не добавлен</p>
          </div>

          <div v-else class="space-y-3">
            <a
                v-for="sheet in gradeSheets"
                :key="sheet.id"
                :href="sheet.sheet_url"
                target="_blank"
                class="block bg-white dark:bg-dark-surface rounded-lg shadow dark:shadow-lg p-4 hover:shadow-md dark:hover:shadow-xl transition-all duration-300 border-2 border-transparent hover:border-forest-mint dark:hover:border-forest-mint-dark"
            >
              <div class="flex items-center gap-3">
                <div class="text-2xl">📊</div>
                <div class="flex-1">
                  <h3 class="font-semibold text-forest-dark dark:text-dark-text transition-colors duration-300">
                    {{ sheet.description || 'Google Sheets журнал' }}
                  </h3>
                  <p class="text-sm text-forest-green dark:text-forest-green-dark transition-colors duration-300">Открыть журнал →</p>
                </div>
              </div>
            </a>
          </div>
        </TabPanel>

        <TabPanel header="📝 Вопросы к экзамену" value="exam">
          <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 mb-4">
            <div v-if="authStore.isAdmin" class="flex gap-2 flex-wrap">
              <button
                  @click="addExamQuestion"
                  class="px-4 py-2 bg-forest-green dark:bg-forest-green-dark text-white rounded-lg hover:bg-forest-dark dark:hover:bg-forest-green transition-colors duration-300"
              >
                ➕ Добавить вопрос
              </button>
              <button
                  @click="showBulkUpload = true"
                  class="px-4 py-2 bg-blue-500 dark:bg-blue-600 text-white rounded-lg hover:bg-blue-600 dark:hover:bg-blue-700 transition-colors duration-300"
              >
                📤 Массовая загрузка
              </button>
              <button
                  v-if="examQuestions.length > 0"
                  @click="toggleSelectionMode"
                  :class="selectionMode
                    ? 'bg-gray-500 dark:bg-gray-600 hover:bg-gray-600 dark:hover:bg-gray-700'
                    : 'bg-orange-500 dark:bg-orange-600 hover:bg-orange-600 dark:hover:bg-orange-700'"
                  class="px-4 py-2 text-white rounded-lg transition-colors duration-300"
              >
                {{ selectionMode ? '✕ Отмена' : '☑️ Выбрать для удаления' }}
              </button>
            </div>
            <button
                @click="showTicketGenerator = true"
                class="px-4 py-2 bg-purple-500 dark:bg-purple-600 text-white rounded-lg hover:bg-purple-600 dark:hover:bg-purple-700 transition-colors duration-300 whitespace-nowrap"
            >
              🎫 Сгенерировать билет
            </button>
          </div>

          <!-- Панель действий при выборе -->
          <div v-if="selectionMode && authStore.isAdmin" class="mb-4 p-4 bg-gray-100 dark:bg-dark-surface rounded-lg border-2 border-orange-300 dark:border-orange-600 transition-colors duration-300">
            <div class="flex flex-wrap items-center gap-3">
              <span class="text-gray-700 dark:text-dark-text font-medium">
                Выбрано: {{ selectedQuestions.size }} из {{ examQuestions.length }}
              </span>
              <button
                  @click="allSelected ? deselectAll() : selectAll()"
                  class="px-3 py-1 bg-gray-200 dark:bg-dark-border text-gray-700 dark:text-dark-text rounded hover:bg-gray-300 dark:hover:bg-gray-600 transition-colors duration-300"
              >
                {{ allSelected ? 'Снять все' : 'Выбрать все' }}
              </button>
              <button
                  v-if="selectedQuestions.size > 0"
                  @click="bulkDeleteQuestions"
                  class="px-4 py-1 bg-red-500 dark:bg-red-600 text-white rounded hover:bg-red-600 dark:hover:bg-red-700 transition-colors duration-300"
              >
                🗑️ Удалить выбранные ({{ selectedQuestions.size }})
              </button>
            </div>
          </div>

          <div v-if="examQuestionsLoading" class="text-center py-8">
            <p class="text-gray-600 dark:text-dark-text-secondary transition-colors duration-300">Загрузка вопросов...</p>
          </div>

          <div v-else-if="examQuestions.length === 0" class="text-center py-8">
            <p class="text-gray-600 dark:text-dark-text-secondary transition-colors duration-300">Вопросы пока не добавлены</p>
          </div>

          <div v-else class="space-y-6">
            <div v-for="group in groupedQuestions" :key="group.section">
              <h3 class="text-lg font-semibold text-forest-dark dark:text-dark-text mb-3 transition-colors duration-300">{{ group.section }}</h3>
              <div class="space-y-2">
                <div
                    v-for="q in group.questions"
                    :key="q.id"
                    :class="[
                      'bg-white dark:bg-dark-surface rounded-lg shadow dark:shadow-lg p-4 flex items-start justify-between transition-all duration-300',
                      selectionMode && selectedQuestions.has(q.id) ? 'ring-2 ring-orange-400 dark:ring-orange-500' : ''
                    ]"
                >
                  <div class="flex gap-3 flex-1 items-start">
                    <!-- Чекбокс в режиме выбора -->
                    <label v-if="selectionMode && authStore.isAdmin" class="flex items-center cursor-pointer mt-1">
                      <input
                          type="checkbox"
                          :checked="selectedQuestions.has(q.id)"
                          @change="toggleQuestion(q.id)"
                          class="w-5 h-5 rounded border-gray-300 dark:border-dark-border text-orange-500 focus:ring-orange-500 cursor-pointer"
                      />
                    </label>
                    <span class="font-semibold text-forest-green dark:text-forest-green-dark transition-colors duration-300">{{ q.number }}.</span>
                    <p class="text-gray-700 dark:text-dark-text-secondary transition-colors duration-300">{{ q.question }}</p>
                  </div>

                  <div v-if="authStore.isAdmin && !selectionMode" class="flex gap-2 ml-4">
                    <button
                        @click="editExamQuestion(q)"
                        class="px-3 py-1 text-sm bg-gray-100 dark:bg-dark-surface hover:bg-gray-200 dark:hover:bg-dark-border rounded transition-colors duration-300 text-gray-700 dark:text-dark-text"
                    >
                      ✏️
                    </button>
                    <button
                        @click="deleteExamQuestion(q.id)"
                        class="px-3 py-1 text-sm bg-red-100 dark:bg-red-900/30 hover:bg-red-200 dark:hover:bg-red-900/50 text-red-600 dark:text-red-400 rounded transition-colors duration-300"
                    >
                      🗑️
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </TabPanel>
      </TabView>
    </main>

    <CourseEditDialog
        v-model="showEditDialog"
        :course="course"
        @saved="handleCourseSaved"
    />

    <LectureEditDialog
        v-model="showLectureDialog"
        :lecture="editingLecture"
        :course-id="courseId"
        @saved="handleLectureSaved"
    />

    <LabEditDialog
        v-model="showLabDialog"
        :lab="editingLab"
        :course-id="courseId"
        @saved="handleLabSaved"
    />

    <ExamQuestionEditDialog
        v-model="showExamQuestionDialog"
        :question="editingExamQuestion"
        :course-id="courseId"
        @saved="handleExamQuestionSaved"
    />

    <ExamQuestionBulkUpload
        v-model="showBulkUpload"
        :course-id="courseId"
        @saved="handleExamQuestionSaved"
    />

    <GradeSheetEditDialog
        v-model="showGradeSheetDialog"
        :grade-sheet="editingGradeSheet"
        :course-id="courseId"
        @saved="handleGradeSheetSaved"
    />

    <TicketGeneratorDialog
        v-if="course"
        v-model:visible="showTicketGenerator"
        :course-id="courseId"
        :course-name="course.name"
        :is-authenticated="authStore.isAuthenticated"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { coursesApi, type Course } from '@/api/courses'
import { lecturesApi, type Lecture } from '@/api/lectures'
import { labsApi, type Lab } from '@/api/labs'
import { gradeSheetsApi, type GradeSheet } from '@/api/gradesheets'
import { examQuestionsApi, type ExamQuestion } from '@/api/examquestions'
import TabView from 'primevue/tabview'
import TabPanel from 'primevue/tabpanel'
import CourseEditDialog from '@/components/CourseEditDialog.vue'
import { useAuthStore } from '@/stores/auth'
import LectureEditDialog from '@/components/LectureEditDialog.vue'
import LabEditDialog from '@/components/LabEditDialog.vue'
import GradeSheetEditDialog from '@/components/GradeSheetEditDialog.vue'
import ExamQuestionEditDialog from '@/components/ExamQuestionEditDialog.vue'
import ExamQuestionBulkUpload from '@/components/ExamQuestionBulkUpload.vue'
import TicketGeneratorDialog from '@/components/TicketGeneratorDialog.vue'


const route = useRoute()
const router = useRouter()

const course = ref<Course | null>(null)
const lectures = ref<Lecture[]>([])
const labs = ref<Lab[]>([])
const gradeSheets = ref<GradeSheet[]>([])
const examQuestions = ref<ExamQuestion[]>([])

const loading = ref(true)
const lecturesLoading = ref(true)
const exportingLectures = ref(false)
const labsLoading = ref(true)
const gradeSheetsLoading = ref(true)
const examQuestionsLoading = ref(true)

const courseId = Number(route.params.id)

const showEditDialog = ref(false)
const authStore = useAuthStore()

const showLectureDialog = ref(false)
const editingLecture = ref<Lecture | null>(null)

const showLabDialog = ref(false)
const editingLab = ref<Lab | null>(null)

const showGradeSheetDialog = ref(false)
const editingGradeSheet = ref<GradeSheet | null>(null)

const showExamQuestionDialog = ref(false)
const showBulkUpload = ref(false)
const showTicketGenerator = ref(false)
const editingExamQuestion = ref<ExamQuestion | null>(null)

const selectedQuestions = ref<Set<number>>(new Set())
const selectionMode = ref(false)

const editCourse = () => {
  showEditDialog.value = true
}

const deleteCourse = async () => {
  if (!confirm(`Удалить курс "${course.value?.name}"? Это действие нельзя отменить.`)) return

  try {
    await coursesApi.delete(courseId)
    router.push('/')
  } catch (error) {
    console.error('Failed to delete course:', error)
    alert('Ошибка удаления курса')
  }
}

const handleCourseSaved = async () => {
  try {
    const { data } = await coursesApi.getById(courseId)
    course.value = data
  } catch (error) {
    console.error('Failed to load course:', error)
  }
}

const goToLecture = (lectureId: number) => {
  router.push({
    path: `/courses/${courseId}/lectures/${lectureId}`,
  })
}

const addLecture = () => {
  editingLecture.value = null
  showLectureDialog.value = true
}

const downloadLecturesMarkdown = async () => {
  exportingLectures.value = true
  try {
    const response = await lecturesApi.exportMarkdown(courseId)
    const blob = new Blob([response.data], { type: 'text/markdown;charset=utf-8' })
    const url = URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = `lectures-course-${courseId}.md`
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    URL.revokeObjectURL(url)
  } catch (error) {
    console.error('Failed to export lectures:', error)
    alert('Ошибка скачивания лекций')
  } finally {
    exportingLectures.value = false
  }
}

const goToLab = (labId: number) => {
  router.push({
    path: `/courses/${courseId}/labs/${labId}`,
  })
}

const addLab = () => {
  editingLab.value = null
  showLabDialog.value = true
}

const handleLabSaved = async () => {
  try {
    const { data } = await labsApi.getAll(courseId)
    labs.value = data.sort((a, b) => a.number - b.number)
  } catch (error) {
    console.error('Failed to load labs:', error)
  }
}

const addOrEditGradeSheet = () => {
  editingGradeSheet.value = gradeSheets.value.length > 0 ? (gradeSheets.value[0] ?? null) : null
  showGradeSheetDialog.value = true
}

const handleGradeSheetSaved = async () => {
  try {
    const { data } = await gradeSheetsApi.getAll(courseId)
    gradeSheets.value = data
  } catch (error) {
    console.error('Failed to load grade sheets:', error)
  }
}

const handleLectureSaved = async () => {
  try {
    const { data } = await lecturesApi.getAll(courseId)
    lectures.value = data.sort((a, b) => a.week - b.week)
  } catch (error) {
    console.error('Failed to load lectures:', error)
  }
}

const groupedQuestions = computed(() => {
  const grouped = examQuestions.value.reduce((acc, q) => {
    if (acc[q.section]) {
      acc[q.section]!.push(q)
    } else {
      acc[q.section] = [q]
    }
    return acc
  }, {} as Record<string, ExamQuestion[]>)

  const groups = Object.entries(grouped).map(([section, questions]) => {
    const sortedQuestions = [...questions].sort((a, b) => {
      const numA = parseInt(String(a.number), 10)
      const numB = parseInt(String(b.number), 10)
      return numA - numB
    })

    const minNumber = Math.min(...sortedQuestions.map(q => parseInt(String(q.number), 10)))

    return {
      section,
      questions: sortedQuestions,
      minNumber
    }
  })

  groups.sort((a, b) => {
    return a.minNumber - b.minNumber
  })


  return groups
})

const allSelected = computed(() =>
  examQuestions.value.length > 0 &&
  selectedQuestions.value.size === examQuestions.value.length
)

const addExamQuestion = () => {
  editingExamQuestion.value = null
  showExamQuestionDialog.value = true
}

const editExamQuestion = (question: ExamQuestion) => {
  editingExamQuestion.value = question
  showExamQuestionDialog.value = true
}

const deleteExamQuestion = async (id: number) => {
  if (!confirm('Удалить вопрос?')) return
  try {
    await examQuestionsApi.delete(id)
    await loadExamQuestions()
  } catch (error) {
    console.error('Failed to delete question:', error)
    alert('Ошибка удаления вопроса')
  }
}

const toggleSelectionMode = () => {
  selectionMode.value = !selectionMode.value
  if (!selectionMode.value) {
    selectedQuestions.value.clear()
  }
}

const toggleQuestion = (id: number) => {
  if (selectedQuestions.value.has(id)) {
    selectedQuestions.value.delete(id)
  } else {
    selectedQuestions.value.add(id)
  }
}

const selectAll = () => {
  examQuestions.value.forEach(q => selectedQuestions.value.add(q.id))
}

const deselectAll = () => {
  selectedQuestions.value.clear()
}

const bulkDeleteQuestions = async () => {
  if (selectedQuestions.value.size === 0) return
  if (!confirm(`Удалить ${selectedQuestions.value.size} вопросов? Это действие нельзя отменить.`)) return

  try {
    await examQuestionsApi.bulkDelete([...selectedQuestions.value])
    selectedQuestions.value.clear()
    selectionMode.value = false
    await loadExamQuestions()
  } catch (error) {
    console.error('Failed to bulk delete questions:', error)
    alert('Ошибка удаления вопросов')
  }
}

const loadExamQuestions = async () => {
  try {
    examQuestionsLoading.value = true
    const { data } = await examQuestionsApi.getAll(courseId)
    examQuestions.value = data
  } catch (error) {
    console.error('Failed to load questions:', error)
  } finally {
    examQuestionsLoading.value = false
  }
}

const handleExamQuestionSaved = async () => {
  await loadExamQuestions()
}

onMounted(async () => {
  try {
    const { data } = await coursesApi.getById(courseId)
    course.value = data
  } catch (error) {
    console.error('Failed to load course:', error)
  } finally {
    loading.value = false
  }



  Promise.all([
    lecturesApi.getAll(courseId).then(({ data }) => {
      lectures.value = data.sort((a, b) => a.week - b.week)
      lecturesLoading.value = false
    }),
    labsApi.getAll(courseId).then(({ data }) => {
      labs.value = data.sort((a, b) => a.number - b.number)
      labsLoading.value = false
    }),
    gradeSheetsApi.getAll(courseId).then(({ data }) => {
      gradeSheets.value = data
      gradeSheetsLoading.value = false
    }),
    loadExamQuestions(),
  ]).catch(error => {
    console.error('Failed to load data:', error)
  })
})
</script>