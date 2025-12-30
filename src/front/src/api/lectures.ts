import apiClient from './client'

export interface Lecture {
    id: number
    course_id: number
    week: number
    title: string
    content: string
    github_url: string | null
    created_at: string
    updated_at: string
}

export const lecturesApi = {
    getAll: (courseId?: number) =>
        apiClient.get<Lecture[]>('/lectures', {
            params: courseId ? { course_id: courseId } : {}
        }),
    getById: (id: number) =>
        apiClient.get<Lecture>(`/lectures/${id}`),

    create: (data: { course_id: number; week: number; title: string; content: string; github_url?: string }) =>
        apiClient.post<Lecture>('/admin/lectures', data),
    update: (id: number, data: { course_id: number; week: number; title: string; content: string; github_url?: string }) =>
        apiClient.put(`/admin/lectures/${id}`, data),
    delete: (id: number) =>
        apiClient.delete(`/admin/lectures/${id}`),
}