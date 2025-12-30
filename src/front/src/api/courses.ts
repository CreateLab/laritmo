import apiClient from './client'

export interface Course {
    id: number
    name: string
    semester: string
    description: string
    created_at: string
    updated_at: string
}

export const coursesApi = {
    getAll: () => apiClient.get<Course[]>('/courses'),
    getById: (id: number) => apiClient.get<Course>(`/courses/${id}`),

    // Admin methods
    create: (data: { name: string; semester: string; description: string }) =>
        apiClient.post<Course>('/admin/courses', data),
    update: (id: number, data: { name: string; semester: string; description: string }) =>
        apiClient.put(`/admin/courses/${id}`, data),
    delete: (id: number) =>
        apiClient.delete(`/admin/courses/${id}`),
}