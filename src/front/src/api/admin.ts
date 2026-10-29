import apiClient from './client'

export const adminApi = {
    dump: () =>
        apiClient.get<Blob>('/admin/dump', {
            responseType: 'blob',
            timeout: 60000,
        }),
}
