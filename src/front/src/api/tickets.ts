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
  ticketsPerPage?: number // Number of tickets per page (default 1)
}

interface TicketResponse {
  ticket: Ticket
}

/**
 * Generates a single random ticket for a course
 * @param courseId Course ID
 * @param questionsCount Number of questions per ticket (1-50)
 * @returns Generated ticket
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
 * Generates multiple tickets and returns TXT file for download
 * @param courseId Course ID
 * @param request Generation parameters
 * @returns Blob with TXT file contents
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
    // If server returned JSON error but we expected blob
    if (error.response?.data instanceof Blob) {
      const text = await error.response.data.text()
      try {
        const jsonError = JSON.parse(text)
        throw new Error(jsonError.error || 'Error generating tickets')
      } catch {
        throw new Error('Error generating tickets')
      }
    }
    throw error
  }
}

/**
 * Downloads blob as file
 * @param blob Blob to download
 * @param filename File name
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
