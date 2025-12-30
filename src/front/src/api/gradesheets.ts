import apiClient from './client'

export interface GradeSheet {
    id: number
    course_id: number
    sheet_url: string
    description: string | null
    created_at: string
    updated_at: string
}

export const gradeSheetsApi = {
    getAll: (courseId?: number) =>
        apiClient.get<GradeSheet[]>('/grade-sheets', {
            params: courseId ? { course_id: courseId } : {}
        }),
}