import apiClient from './client'

export interface Lab {
    id: number
    course_id: number
    number: number
    title: string
    description: string
    deadline: string | null
    max_score: number
    github_url: string | null
    created_at: string
    updated_at: string
}

export const labsApi = {
    getAll: (courseId?: number) =>
        apiClient.get<Lab[]>('/labs', {
            params: courseId ? { course_id: courseId } : {}
        }),
    getById: (id: number) =>
        apiClient.get<Lab>(`/labs/${id}`),

    create: (data: { course_id: number; number: number; title: string; description: string; max_score: number; github_url?: string; deadline?: string }) =>
        apiClient.post<Lab>('/admin/labs', data),
    update: (id: number, data: { course_id: number; number: number; title: string; description: string; max_score: number; github_url?: string; deadline?: string }) =>
        apiClient.put(`/admin/labs/${id}`, data),
    delete: (id: number) =>
        apiClient.delete(`/admin/labs/${id}`),
}