<script setup lang="ts">
const props = defineProps<{
  id: string
  keyPrefix: string
  keySuffix: string
  label: string | null
  createdAt: string
  lastUsed: string | null
}>()

const emit = defineEmits<{
  delete: [id: string]
}>()

const { t } = useI18n()

const formatDate = (dateStr: string) => {
  return new Date(dateStr).toLocaleDateString(undefined, {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
  })
}
</script>

<template>
  <div class="group flex items-center justify-between gap-4 px-4 py-3.5 rounded-xl bg-surface-800 border border-white/[0.06] hover:border-white/10 transition-all duration-200">
    <div class="flex items-center gap-3.5 min-w-0">
      <!-- Key icon -->
      <div class="shrink-0 w-8 h-8 rounded-lg bg-highlight/8 flex items-center justify-center">
        <svg class="w-3.5 h-3.5 text-highlight/60" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z" />
        </svg>
      </div>

      <div class="min-w-0">
        <div class="flex items-center gap-2 mb-0.5">
          <code class="text-sm font-mono text-accent-glow tracking-wide">{{ keyPrefix }}••••{{ keySuffix }}</code>
          <span v-if="label" class="text-[10px] px-1.5 py-0.5 rounded bg-white/[0.04] text-gray-500 font-medium">
            {{ label }}
          </span>
        </div>
        <div class="flex items-center gap-3 text-[11px] text-gray-600">
          <span>{{ t('dashboard.created') }} {{ formatDate(createdAt) }}</span>
          <span class="text-gray-700">·</span>
          <span>
            {{ t('dashboard.lastUsed') }}
            {{ lastUsed ? formatDate(lastUsed) : t('dashboard.never') }}
          </span>
        </div>
      </div>
    </div>

    <button
      class="shrink-0 px-2.5 py-1 rounded-md text-[11px] font-medium text-gray-600 hover:text-danger hover:bg-danger/8 opacity-0 group-hover:opacity-100 transition-all duration-200"
      @click="emit('delete', id)"
    >
      {{ t('dashboard.deleteKey') }}
    </button>
  </div>
</template>
