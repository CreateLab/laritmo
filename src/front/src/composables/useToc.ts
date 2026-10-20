import { nextTick, onBeforeUnmount, ref, watch, type Ref, type WatchSource } from 'vue'

export interface TocItem {
  id: string
  text: string
  level: number
}

// Строит оглавление из h2/h3 отрендеренного markdown и подсвечивает активный раздел при скролле
export function useToc(container: Ref<HTMLElement | null>, source: WatchSource) {
  const items = ref<TocItem[]>([])
  const activeId = ref('')
  let observer: IntersectionObserver | null = null

  const disconnect = () => {
    observer?.disconnect()
    observer = null
  }

  const build = async () => {
    disconnect()
    await nextTick()
    const root = container.value
    if (!root) {
      items.value = []
      return
    }
    const headings = Array.from(root.querySelectorAll<HTMLElement>('h2, h3'))
    items.value = headings.map((h, i) => {
      if (!h.id) h.id = `section-${i + 1}`
      return { id: h.id, text: h.textContent?.trim() ?? '', level: h.tagName === 'H2' ? 2 : 3 }
    })
    activeId.value = items.value[0]?.id ?? ''
    if (!headings.length || typeof IntersectionObserver === 'undefined') return
    observer = new IntersectionObserver(
      (entries) => {
        const visible = entries.filter((e) => e.isIntersecting)
        if (visible[0]) activeId.value = visible[0].target.id
      },
      { rootMargin: '-96px 0px -70% 0px' },
    )
    headings.forEach((h) => observer!.observe(h))
  }

  const scrollTo = (id: string) => {
    document.getElementById(id)?.scrollIntoView({ behavior: 'smooth', block: 'start' })
    activeId.value = id
  }

  watch(source, build, { flush: 'post' })
  onBeforeUnmount(disconnect)

  return { items, activeId, scrollTo }
}
