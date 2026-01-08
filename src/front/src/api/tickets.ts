import apiClient from './client'

export interface Ticket {
  number: number
  questions: Question[]
}

export interface Question {
  number: number
  section: string
  question: string
}

export interface TicketGenerationRequest {
  questionsPerTicket: number
  ticketCount: number
}

interface TicketResponse {
  ticket: Ticket
}

/**
 * Генерирует один случайный билет для курса
 * @param courseId ID курса
 * @param questionsCount Количество вопросов в билете (1-50)
 * @returns Сгенерированный билет
 */
export async function generateRandomTicket(
  courseId: number,
  questionsCount: number
): Promise<Ticket> {
  const { data } = await apiClient.get<TicketResponse>(
    `/courses/${courseId}/tickets/random`,
    {
      params: { questions: questionsCount },
    }
  )
  return data.ticket
}

/**
 * Генерирует несколько билетов и возвращает TXT файл для скачивания
 * @param courseId ID курса
 * @param request Параметры генерации
 * @returns Blob с содержимым TXT файла
 */
export async function generateTicketsDocument(
  courseId: number,
  request: TicketGenerationRequest
): Promise<Blob> {
  try {
    const { data } = await apiClient.post<Blob>(
      `/admin/courses/${courseId}/tickets/generate`,
      request,
      {
        responseType: 'blob',
      }
    )
    return data
  } catch (error: any) {
    // Если сервер вернул ошибку в виде JSON, но мы ожидали blob
    if (error.response?.data instanceof Blob) {
      const text = await error.response.data.text()
      try {
        const jsonError = JSON.parse(text)
        throw new Error(jsonError.error || 'Ошибка при генерации билетов')
      } catch {
        throw new Error('Ошибка при генерации билетов')
      }
    }
    throw error
  }
}

/**
 * Скачивает blob как файл
 * @param blob Blob для скачивания
 * @param filename Имя файла
 */
export function downloadBlob(blob: Blob, filename: string): void {
  const url = window.URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = filename
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  window.URL.revokeObjectURL(url)
}
