<template>
  <div class="max-w-7xl mx-auto px-4 sm:px-8 py-8 sm:py-12">
    <header class="mb-8">
      <button
          @click="router.push('/')"
          class="text-action hover:text-action-hover text-sm mb-4 flex items-center gap-2 transition-colors duration-150"
      >
        ← На главную
      </button>

      <div v-if="loading" class="animate-pulse">
        <div class="h-8 bg-sunken rounded-sm w-3/4 mb-2"></div>
        <div class="h-4 bg-sunken rounded-sm w-1/4"></div>
      </div>

      <div v-else-if="course" class="flex flex-wrap items-start justify-between gap-4">
        <div class="min-w-0">
          <div class="flex flex-wrap items-center gap-3 mb-2">
            <h1 class="text-h1 text-ink">{{ course.name }}</h1>
            <AppChip>{{ course.semester }}</AppChip>
          </div>
          <p class="text-ink-muted">{{ course.description }}</p>
        </div>

        <div v-if="authStore.isAdmin" class="flex items-center gap-2">
          <AppButton icon="pi-pencil" @click="editCourse">Редактировать</AppButton>
          <AppButton variant="danger" icon="pi-trash" @click="deleteCourse">Удалить</AppButton>
        </div>
      </div>
    </header>

    <main>
      <TabView v-if="!loading">
        <TabPanel header="Лекции" value="lectures">

          <div class="mb-4 flex gap-2 flex-wrap">
            <AppButton
                v-if="authStore.isAdmin"
                variant="primary"
                icon="pi-plus"
                @click="addLecture"
            >
              Добавить лекцию
            </AppButton>
            <AppButton
                v-if="authStore.isAuthenticated && lectures.length > 0"
                icon="pi-download"
                @click="downloadLecturesMarkdown"
                :disabled="exportingLectures"
            >
              {{ exportingLectures ? 'Скачивание...' : 'Скачать все в Markdown' }}
            </AppButton>
          </div>

          <div v-if="lecturesLoading" class="text-center py-8">
            <p class="text-ink-muted">Загрузка лекций...</p>
          </div>

          <EmptyState v-else-if="lectures.length === 0" title="Лекции пока не добавлены" text="Когда преподаватель добавит материалы, они появятся здесь." />

          <div v-else class="space-y-3">
            <div
                v-for="lecture in lectures"
                :key="lecture.id"
                @click="goToLecture(lecture.id)"
                class="row-item cursor-pointer"
            >
              <AppChip>Неделя {{ lecture.week }}</AppChip>
              <h3 class="text-body font-medium text-ink">{{ lecture.title }}</h3>
            </div>
          </div>
        </TabPanel>

        <TabPanel header="Лабораторные" value="labs">
          <div v-if="authStore.isAdmin" class="mb-4">
            <AppButton variant="primary" icon="pi-plus" @click="addLab">Добавить лабораторную</AppButton>
          </div>

          <div v-if="labsLoading" class="text-center py-8">
            <p class="text-ink-muted">Загрузка лаб...</p>
          </div>

          <EmptyState v-else-if="labs.length === 0" title="Лабораторные пока не добавлены" text="Задания по курсу появятся здесь." />

          <div v-else class="space-y-3">
            <div
                v-for="lab in labs"
                :key="lab.id"
                @click="goToLab(lab.id)"
                class="row-item cursor-pointer items-start"
            >
              <div class="flex-1 min-w-0">
                <div class="flex flex-wrap items-center gap-2 mb-2">
                  <AppChip>Лаба #{{ lab.number }}</AppChip>
                  <AppChip>Макс: {{ lab.max_score }} баллов</AppChip>
                  <span
                      v-if="lab.deadline"
                      :class="isOverdue(lab.deadline) ? 'bg-danger-bg text-danger-text' : 'bg-sunken text-ink'"
                      class="chip"
                  >
                    Дедлайн: {{ formatDate(lab.deadline) }}
                  </span>
                </div>
                <h3 class="text-body font-medium text-ink mb-1">{{ lab.title }}</h3>
                <p class="text-body-sm text-ink-muted line-clamp-2">
                  {{ lab.description.substring(0, 150) }}...
                </p>
              </div>
            </div>
          </div>
        </TabPanel>

        <TabPanel header="Журнал" value="grades">
          <div v-if="authStore.isAdmin" class="mb-4">
            <AppButton
                variant="primary"
                :icon="gradeSheets.length > 0 ? 'pi-pencil' : 'pi-plus'"
                @click="addOrEditGradeSheet"
            >
              {{ gradeSheets.length > 0 ? 'Изменить ссылку на журнал' : 'Добавить ссылку на журнал' }}
            </AppButton>
          </div>

          <div v-if="gradeSheetsLoading" class="text-center py-8">
            <p class="text-ink-muted">Загрузка...</p>
          </div>

          <EmptyState v-else-if="!gradeSheets || gradeSheets.length === 0" title="Журнал пока не добавлен" text="Ведомость с оценками появится здесь." />

          <div v-else class="space-y-3">
            <template v-for="sheet in gradeSheets" :key="sheet.id">
            <div v-if="sheet.table && sheet.table.length" class="space-y-2">
              <div class="flex items-center justify-between gap-3">
                <h3 class="text-h3 text-ink">
                  {{ sheet.description || 'Google Sheets журнал' }}
                </h3>
                <a
                    :href="sheet.sheet_url"
                    target="_blank"
                    rel="noopener"
                    class="text-body-sm text-action hover:text-action-hover hover:underline"
                >Открыть в новой вкладке ↗</a>
              </div>
              <div class="table-wrap max-h-[70vh] overflow-auto">
                <table class="data-table">
                  <thead class="sticky top-0">
                    <tr>
                      <th
                          v-for="(cell, ci) in sheet.table[0]"
                          :key="ci"
                          class="text-label"
                      >{{ cell }}</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr
                        v-for="(row, ri) in sheet.table.slice(1)"
                        :key="ri"
                    >
                      <td v-for="(cell, ci) in row" :key="ci" class="whitespace-nowrap">{{ cell }}</td>
                    </tr>
                  </tbody>
                </table>
              </div>
            </div>
            <a
                v-else
                :href="sheet.sheet_url"
                target="_blank"
                class="row-item"
            >
              <div class="flex-1">
                <h3 class="text-body font-medium text-ink">
                  {{ sheet.description || 'Google Sheets журнал' }}
                </h3>
                <p class="text-body-sm text-action">Открыть журнал →</p>
              </div>
            </a>
            </template>
          </div>
        </TabPanel>

        <TabPanel header="Вопросы к экзамену" value="exam">
          <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 mb-4">
            <div v-if="authStore.isAdmin" class="flex gap-2 flex-wrap">
              <AppButton variant="primary" icon="pi-plus" @click="addExamQuestion">Добавить вопрос</AppButton>
              <AppButton icon="pi-upload" @click="showBulkUpload = true">Массовая загрузка</AppButton>
              <AppButton
                  v-if="examQuestions.length > 0"
                  :icon="selectionMode ? 'pi-times' : 'pi-check-square'"
                  @click="toggleSelectionMode"
              >
                {{ selectionMode ? 'Отмена' : 'Выбрать для удаления' }}
              </AppButton>
            </div>
            <AppButton icon="pi-ticket" @click="showTicketGenerator = true">Сгенерировать билет</AppButton>
          </div>

          <!-- Панель действий при выборе -->
          <div v-if="selectionMode && authStore.isAdmin" class="mb-4 p-4 bg-sunken rounded-md border border-border-strong">
            <div class="flex flex-wrap items-center gap-3">
              <span class="text-label text-ink">
                Выбрано: {{ selectedQuestions.size }} из {{ examQuestions.length }}
              </span>
              <AppButton @click="allSelected ? deselectAll() : selectAll()">
                {{ allSelected ? 'Снять все' : 'Выбрать все' }}
              </AppButton>
              <AppButton
                  v-if="selectedQuestions.size > 0"
                  variant="danger"
                  icon="pi-trash"
                  @click="bulkDeleteQuestions"
              >
                Удалить выбранные ({{ selectedQuestions.size }})
              </AppButton>
            </div>
          </div>

          <div v-if="examQuestionsLoading" class="text-center py-8">
            <p class="text-ink-muted">Загрузка вопросов...</p>
          </div>

          <EmptyState v-else-if="examQuestions.length === 0" title="Вопросы пока не добавлены" text="Список вопросов к экзамену появится здесь." />

          <div v-else class="space-y-6">
            <div v-for="group in groupedQuestions" :key="group.section">
              <h3 class="text-h3 text-ink mb-3">{{ group.section }}</h3>
              <div class="space-y-2">
                <div
                    v-for="q in group.questions"
                    :key="q.id"
                    :class="[
                      'row-item items-start justify-between',
                      selectionMode && selectedQuestions.has(q.id) ? 'outline outline-2 outline-action' : ''
                    ]"
                >
                  <div class="flex gap-3 flex-1 items-start">
                    <!-- Чекбокс в режиме выбора -->
                    <label v-if="selectionMode && authStore.isAdmin" class="flex items-center cursor-pointer mt-1">
                      <input
                          type="checkbox"
                          :checked="selectedQuestions.has(q.id)"
                          @change="toggleQuestion(q.id)"
                          class="w-5 h-5 cursor-pointer accent-[var(--action)]"
                      />
                    </label>
                    <span class="text-body font-bold text-action">{{ q.number }}.</span>
                    <p class="text-body text-ink">{{ q.question }}</p>
                  </div>

                  <div v-if="authStore.isAdmin && !selectionMode" class="flex gap-2 ml-4">
                    <AppButton icon="pi-pencil" aria-label="Редактировать" @click="editExamQuestion(q)" />
                    <AppButton variant="danger" icon="pi-trash" aria-label="Удалить" @click="deleteExamQuestion(q.id)" />
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
import EmptyState from '@/components/ui/EmptyState.vue'
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { coursesApi, type Course } from '@/api/courses'
import { lecturesApi, type Lecture } from '@/api/lectures'
import { labsApi, type Lab } from '@/api/labs'
import { gradeSheetsApi, type GradeSheet } from '@/api/gradesheets'
import { examQuestionsApi, type ExamQuestion } from '@/api/examquestions'
import TabView from 'primevue/tabview'
import AppButton from '@/components/ui/AppButton.vue'
import AppChip from '@/components/ui/AppChip.vue'
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

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleDateString('ru-RU', {
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  })
}

const isOverdue = (dateString: string) => {
  return new Date(dateString) < new Date()
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