import apiClient from './client'

export interface User {
    id: number
    username: string
    email: string
    role: 'owner' | 'admin' | 'student'
    is_active: boolean
    created_at: string
    updated_at: string
    last_login_at?: string
}

export interface CreateUserRequest {
    username: string
    email: string
    password: string
    role: 'admin' | 'student'
}

export interface UpdateUserRequest {
    email?: string
    role?: 'owner' | 'admin' | 'student'
    is_active?: boolean
}

export const usersApi = {
    getAll: () => apiClient.get<User[]>('/admin/users'),
    getById: (id: number) => apiClient.get<User>(`/admin/users/${id}`),
    create: (data: CreateUserRequest) => apiClient.post<User>('/admin/users', data),
    update: (id: number, data: UpdateUserRequest) => apiClient.put(`/admin/users/${id}`, data),
    delete: (id: number) => apiClient.delete(`/admin/users/${id}`),
    deactivate: (id: number) => apiClient.put(`/admin/users/${id}/deactivate`),
    activate: (id: number) => apiClient.put(`/admin/users/${id}/activate`),
    resetPassword: (id: number, newPassword: string) =>
        apiClient.put(`/admin/users/${id}/password`, { new_password: newPassword }),
}
