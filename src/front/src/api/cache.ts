import type { AxiosResponse, InternalAxiosRequestConfig } from 'axios'

// Session cache for rarely-changing GET data (courses, lectures, labs, ...).
// In-memory Map mirrored to sessionStorage so it survives F5 within the tab.

const STORAGE_KEY = 'api-cache'
const TTL_MS = 30 * 60 * 1000

const CACHEABLE = /^\/(courses(\/\d+)?|lectures(\/\d+)?|labs(\/\d+)?|grade-sheets(\/\d+)?|exam-questions(\/\d+)?)$/

interface Entry {
    data: unknown
    status: number
    ts: number
}

const store = new Map<string, Entry>()

try {
    const raw = sessionStorage.getItem(STORAGE_KEY)
    if (raw) {
        for (const [k, v] of Object.entries(JSON.parse(raw) as Record<string, Entry>)) {
            store.set(k, v)
        }
    }
} catch {
    // sessionStorage unavailable or corrupted: start empty
}

const persist = () => {
    try {
        sessionStorage.setItem(STORAGE_KEY, JSON.stringify(Object.fromEntries(store)))
    } catch {
        // quota exceeded or storage blocked: in-memory only
    }
}

export const isCacheable = (config: InternalAxiosRequestConfig): boolean =>
    (config.method ?? 'get').toLowerCase() === 'get' &&
    config.responseType !== 'blob' &&
    CACHEABLE.test(config.url ?? '')

export const cacheKey = (config: InternalAxiosRequestConfig): string =>
    `${config.url}?${JSON.stringify(config.params ?? {})}`

export const getCached = (key: string): Entry | undefined => {
    const entry = store.get(key)
    if (!entry) return undefined
    if (Date.now() - entry.ts > TTL_MS) {
        store.delete(key)
        return undefined
    }
    return entry
}

export const setCached = (key: string, response: AxiosResponse) => {
    store.set(key, { data: response.data, status: response.status, ts: Date.now() })
    persist()
}

export const clearCache = () => {
    store.clear()
    try {
        sessionStorage.removeItem(STORAGE_KEY)
    } catch {
        // ignore
    }
}
