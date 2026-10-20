<template>
  <p v-if="columns.length <= 1" class="text-body text-ink-muted py-4">Данные журнала пока пусты</p>

  <div v-else class="grade-wrap relative rounded-lg border border-border bg-surface overflow-hidden">
    <div
        ref="scroller"
        class="grade-scroll overflow-auto max-h-[calc(100vh-14rem)]"
        :class="{ 'is-scrolled': canLeft }"
        @scroll.passive="updateShadows"
    >
      <table class="grade-table">
        <thead>
          <tr>
            <th
                v-for="(col, ci) in columns"
                :key="col.index"
                :class="['text-label', ci === 0 ? 'col-first' : 'col-num', col.total ? 'font-semibold' : '']"
                :title="ci === 0 ? col.header : undefined"
            >
              <span :class="ci === 0 ? 'name-cell' : ''">{{ col.header }}</span>
            </th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="(row, ri) in rows" :key="ri">
            <td
                v-for="(col, ci) in columns"
                :key="col.index"
                :class="[ci === 0 ? 'col-first' : 'col-num', col.total ? 'font-semibold' : '']"
                :title="ci === 0 ? cell(row, col.index) : undefined"
            >
              <span :class="ci === 0 ? 'name-cell' : ''">{{ cell(row, col.index) }}</span>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
    <div v-show="canRight" class="shadow-right" aria-hidden="true"></div>
  </div>
</template>

<script setup lang="ts">
import { computed, nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'

const props = defineProps<{ table: string[][] }>()

const scroller = ref<HTMLElement | null>(null)
const canLeft = ref(false)
const canRight = ref(false)

const norm = (v: unknown) => String(v ?? '').trim()
const rows = computed(() => props.table.slice(1))
const cell = (row: string[], i: number) => norm(row[i])

// фильтрация колонок: исходные данные не мутируем
const columns = computed(() => {
  const head = props.table[0] ?? []
  const width = Math.max(head.length, ...rows.value.map((r) => r.length), 0)
  const result: { index: number; header: string; total: boolean }[] = []
  for (let i = 0; i < width; i++) {
    if (i > 0) {
      const values = rows.value.map((r) => cell(r, i))
      if (values.every((v) => v === '')) continue
      if (rows.value.every((r) => cell(r, i) === cell(r, 0))) continue
    }
    const header = norm(head[i])
    result.push({
      index: i,
      header: header || (i > 0 ? '—' : ''),
      total: /^(экз|итог|сумм)/i.test(header),
    })
  }
  return result
})

const updateShadows = () => {
  const el = scroller.value
  if (!el) return
  canLeft.value = el.scrollLeft > 1
  canRight.value = el.scrollLeft + el.clientWidth < el.scrollWidth - 1
}

let ro: ResizeObserver | null = null
onMounted(() => {
  updateShadows()
  if (typeof ResizeObserver !== 'undefined' && scroller.value) {
    ro = new ResizeObserver(updateShadows)
    ro.observe(scroller.value)
  }
})
onBeforeUnmount(() => ro?.disconnect())
watch(columns, () => nextTick(updateShadows))
</script>

<style scoped>
.grade-scroll { scrollbar-width: thin; scrollbar-color: var(--border-strong) transparent; }
.grade-table {
  border-collapse: separate;
  border-spacing: 0;
  width: 100%;
  font-size: 0.875rem;
  line-height: 1.25rem;
  font-variant-numeric: tabular-nums;
}
.grade-table th,
.grade-table td {
  padding: 0.625rem 0.75rem;
  white-space: nowrap;
  border-top: 1px solid var(--border);
  background: var(--surface);
}
.grade-table th {
  position: sticky;
  top: 0;
  z-index: 2;
  background: var(--sunken);
  border-top: 0;
  font-weight: 500;
}
.grade-table tbody tr:hover td { background: var(--sunken); }
.grade-table .col-num {
  min-width: 4.5rem;
  text-align: center;
  border-left: 1px solid var(--border);
}
.grade-table .col-first {
  position: sticky;
  left: 0;
  z-index: 1;
  text-align: left;
  border-right: 1px solid var(--border-strong);
}
.grade-table th.col-first { z-index: 3; }
.is-scrolled .col-first { box-shadow: 4px 0 8px -4px rgb(0 0 0 / 0.25); }
.name-cell {
  display: block;
  min-width: 10rem;
  max-width: 18rem;
  overflow: hidden;
  text-overflow: ellipsis;
}
@media (min-width: 640px) {
  .name-cell { min-width: 14rem; }
}
.shadow-right {
  position: absolute;
  top: 0;
  right: 0;
  bottom: 0;
  width: 1.5rem;
  pointer-events: none;
  z-index: 4;
  background: linear-gradient(to left, rgb(0 0 0 / 0.14), transparent);
}
</style>
