<template>
  <div class="min-h-screen p-6">
    <header class="mb-8">
      <button
          @click="router.push('/')"
          class="text-forest-green hover:text-forest-dark mb-4 flex items-center gap-2"
      >
        ← На главную
      </button>

      <div v-if="loading" class="animate-pulse">
        <div class="h-8 bg-gray-200 rounded w-3/4 mb-2"></div>
        <div class="h-4 bg-gray-200 rounded w-1/4"></div>
      </div>

      <div v-else-if="course" class="flex items-start justify-between">
        <div>
          <h1 class="text-3xl font-bold text-forest-dark mb-2">{{ course.name }}</h1>
          <p class="text-gray-600">{{ course.semester }}</p>
          <p class="text-gray-700 mt-2">{{ course.description }}</p>
        </div>

        <!-- Кнопка админа для курса -->
        <div v-if="authStore.isAdmin">
          <button
              @click="editCourse"
              class="px-4 py-2 bg-forest-green text-white rounded-lg hover:bg-forest-dark transition-colors"
          >
            ✏️ Редактировать
          </button>
          <button
              @click="deleteCourse"
              class="flex gap-6 mb-8 px-4 py-2 bg-red-500 text-white rounded-lg hover:bg-red-600 transition-colors"
          >
            🗑️ Удалить
          </button>
        </div>
      </div>
    </header>

    <main>
      <TabView v-if="!loading">
        <!-- Таб: Лекции -->

        <TabPanel header="📚 Лекции" value="lectures">

          <div v-if="authStore.isAdmin" class="mb-4">
            <button
                @click="addLecture"
                class="px-4 py-2 bg-forest-green text-white rounded-lg hover:bg-forest-dark transition-colors"
            >
              ➕ Добавить лекцию
            </button>
          </div>

          <div v-if="lecturesLoading" class="text-center py-8">
            <p class="text-gray-600">Загрузка лекций...</p>
          </div>

          <div v-else-if="lectures.length === 0" class="text-center py-8">
            <p class="text-gray-600">Лекции пока не добавлены</p>
          </div>

          <div v-else class="space-y-3">
            <div
                v-for="lecture in lectures"
                :key="lecture.id"
                @click="goToLecture(lecture.id)"
                class="bg-white rounded-lg shadow p-4 hover:shadow-md transition-shadow cursor-pointer border-2 border-transparent hover:border-forest-mint"
            >
              <div class="flex items-start gap-3">
                <div class="text-2xl">🍂</div>
                <div class="flex-1">
                  <div class="flex items-center gap-2 mb-1">
                    <span class="px-2 py-1 bg-forest-green text-white rounded-full text-xs">
                      Неделя {{ lecture.week }}
                    </span>
                  </div>
                  <h3 class="font-semibold text-forest-dark">{{ lecture.title }}</h3>
                </div>
              </div>
            </div>
          </div>
        </TabPanel>

        <!-- Таб: Лабораторные -->
        <TabPanel header="🔬 Лабораторные" value="labs">
          <div v-if="authStore.isAdmin" class="mb-4">
            <button
                @click="addLab"
                class="px-4 py-2 bg-forest-green text-white rounded-lg hover:bg-forest-dark transition-colors"
            >
              ➕ Добавить лабораторную
            </button>
          </div>

          <div v-if="labsLoading" class="text-center py-8">
            <p class="text-gray-600">Загрузка лаб...</p>
          </div>

          <div v-else-if="labs.length === 0" class="text-center py-8">
            <p class="text-gray-600">Лабораторные пока не добавлены</p>
          </div>

          <div v-else class="space-y-3">
            <div
                v-for="lab in labs"
                :key="lab.id"
                @click="goToLab(lab.id)"
                class="bg-white rounded-lg shadow p-4 hover:shadow-md transition-shadow border-2 border-transparent hover:border-forest-mint"
            >
              <div class="flex items-start gap-3">
                <div class="text-2xl">🧪</div>
                <div class="flex-1">
                  <div class="flex items-center gap-2 mb-1">
                    <span class="px-2 py-1 bg-forest-green text-white rounded-full text-xs">
                      Лаба #{{ lab.number }}
                    </span>
                    <span class="text-xs text-gray-600">Макс: {{ lab.max_score }} баллов</span>
                  </div>
                  <h3 class="font-semibold text-forest-dark mb-2">{{ lab.title }}</h3>
                  <p class="text-sm text-gray-600 line-clamp-2">
                    {{ lab.description.substring(0, 150) }}...
                  </p>
                </div>
              </div>
            </div>
          </div>
        </TabPanel>

        <!-- Таб: Журнал -->
        <TabPanel header="📊 Журнал" value="grades">
          <div v-if="authStore.isAdmin" class="mb-4">
            <button
                @click="addOrEditGradeSheet"
                class="px-4 py-2 bg-forest-green text-white rounded-lg hover:bg-forest-dark transition-colors"
            >
              {{ gradeSheets.length > 0 ? '✏️ Изменить ссылку на журнал' : '➕ Добавить ссылку на журнал' }}
            </button>
          </div>

          <div v-if="gradeSheetsLoading" class="text-center py-8">
            <p class="text-gray-600">Загрузка...</p>
          </div>

          <div v-else-if="!gradeSheets || gradeSheets.length === 0" class="text-center py-8">
            <p class="text-gray-600">Журнал пока не добавлен</p>
          </div>

          <div v-else class="space-y-3">
            <a
                v-for="sheet in gradeSheets"
                :key="sheet.id"
                :href="sheet.sheet_url"
                target="_blank"
                class="block bg-white rounded-lg shadow p-4 hover:shadow-md transition-shadow border-2 border-transparent hover:border-forest-mint"
            >
              <div class="flex items-center gap-3">
                <div class="text-2xl">📊</div>
                <div class="flex-1">
                  <h3 class="font-semibold text-forest-dark">
                    {{ sheet.description || 'Google Sheets журнал' }}
                  </h3>
                  <p class="text-sm text-forest-green">Открыть журнал →</p>
                </div>
              </div>
            </a>
          </div>
        </TabPanel>

        <!-- Таб: Экзамен -->
        <TabPanel header="📝 Вопросы к экзамену" value="exam">
          <!-- Кнопки админа -->
          <div v-if="authStore.isAdmin" class="flex gap-4 mb-4">
            <button
                @click="addExamQuestion"
                class="px-4 py-2 bg-forest-green text-white rounded-lg hover:bg-forest-dark transition-colors"
            >
              ➕ Добавить вопрос
            </button>
            <button
                @click="showBulkUpload = true"
                class="px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition-colors"
            >
              📤 Массовая загрузка
            </button>
          </div>

          <div v-if="examQuestionsLoading" class="text-center py-8">
            <p class="text-gray-600">Загрузка вопросов...</p>
          </div>

          <div v-else-if="examQuestions.length === 0" class="text-center py-8">
            <p class="text-gray-600">Вопросы пока не добавлены</p>
          </div>

          <!-- Группировка по секциям -->
          <div v-else class="space-y-6">
            <div v-for="(questions, section) in groupedQuestions" :key="section">
              <h3 class="text-lg font-semibold text-forest-dark mb-3">{{ section }}</h3>
              <div class="space-y-2">
                <div
                    v-for="q in questions"
                    :key="q.id"
                    class="bg-white rounded-lg shadow p-4 flex items-start justify-between"
                >
                  <div class="flex gap-3 flex-1">
                    <span class="font-semibold text-forest-green">{{ q.number }}.</span>
                    <p class="text-gray-700">{{ q.question }}</p>
                  </div>

                  <div v-if="authStore.isAdmin" class="flex gap-2 ml-4">
                    <button
                        @click="editExamQuestion(q)"
                        class="px-3 py-1 text-sm bg-gray-100 hover:bg-gray-200 rounded transition-colors"
                    >
                      ✏️
                    </button>
                    <button
                        @click="deleteExamQuestion(q.id)"
                        class="px-3 py-1 text-sm bg-red-100 hover:bg-red-200 text-red-600 rounded transition-colors"
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


const route = useRoute()
const router = useRouter()

const course = ref<Course | null>(null)
const lectures = ref<Lecture[]>([])
const labs = ref<Lab[]>([])
const gradeSheets = ref<GradeSheet[]>([])
const examQuestions = ref<ExamQuestion[]>([])

const loading = ref(true)
const lecturesLoading = ref(true)
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
const editingExamQuestion = ref<ExamQuestion | null>(null)

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
  return examQuestions.value.reduce((acc, q) => {
    if (acc[q.section]) {
      acc[q.section]!.push(q)
    } else {
      acc[q.section] = [q]
    }
    return acc
  }, {} as Record<string, ExamQuestion[]>)
})

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