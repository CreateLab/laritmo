import apiClient from './client'

export interface ExamQuestion {
  id: number
  course_id: number
  number: number
  section: string
  question: string
  created_at: string
  updated_at: string
}

export const examQuestionsApi = {
  // Public
  getAll: (courseId?: number) =>
    apiClient.get<ExamQuestion[]>('/exam-questions', {
      params: courseId ? { course_id: courseId } : {}
    }),
  getById: (id: number) =>
    apiClient.get<ExamQuestion>(`/exam-questions/${id}`),

  // Admin
  create: (data: { course_id: number; number: number; section: string; question: string }) =>
    apiClient.post<ExamQuestion>('/admin/exam-questions', data),

  bulkCreateJSON: (data: { course_id: number; questions: Array<{ number: number; section: string; question: string }> }) =>
    apiClient.post('/admin/exam-questions/bulk', data),

  uploadFile: (courseId: number, file: File) => {
    const formData = new FormData()
    formData.append('course_id', courseId.toString())
    formData.append('file', file)
    return apiClient.post('/admin/exam-questions/upload', formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })
  },

  update: (id: number, data: { number: number; section: string; question: string }) =>
    apiClient.put(`/admin/exam-questions/${id}`, data),

  delete: (id: number) =>
    apiClient.delete(`/admin/exam-questions/${id}`),
}
