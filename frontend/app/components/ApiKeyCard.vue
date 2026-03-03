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

const displayName = computed(() => props.label || t('dashboard.unnamedAgent'))
</script>

<template>
  <div class="group rounded-xl bg-th-bg-s border border-th-text/[var(--color-card-border-opacity)] hover:border-th-text/[0.12] transition-all duration-200">
    <!-- Desktop: horizontal layout -->
    <div class="hidden sm:flex items-center justify-between gap-4 px-5 py-4">
      <div class="flex items-center gap-4 min-w-0">
        <div class="shrink-0 w-10 h-10 rounded-lg bg-highlight/8 flex items-center justify-center">
          <svg class="w-4 h-4 text-highlight/60" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.75 17L9 20l-1 1h8l-1-1-.75-3M3 13h18M5 17h14a2 2 0 002-2V5a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
          </svg>
        </div>
        <div class="min-w-0">
          <div class="flex items-center gap-2.5 mb-1">
            <span class="text-base font-medium text-th-text truncate">{{ displayName }}</span>
            <span
              class="shrink-0 flex items-center gap-1.5 text-xs px-2 py-0.5 rounded-full font-medium"
              :class="lastUsed ? 'bg-success/10 text-success' : 'bg-th-text/[0.04] text-th-text-m'"
            >
              <span class="w-1.5 h-1.5 rounded-full" :class="lastUsed ? 'bg-success' : 'bg-th-text-m'" />
              {{ lastUsed ? t('dashboard.active') : t('dashboard.never') }}
            </span>
          </div>
          <div class="flex items-center gap-3 text-xs text-th-text-t">
            <span>{{ t('dashboard.created') }} {{ formatDate(createdAt) }}</span>
            <template v-if="lastUsed">
              <span class="text-th-text-m">·</span>
              <span>{{ t('dashboard.lastUsed') }} {{ formatDate(lastUsed) }}</span>
            </template>
            <span class="text-th-text-m">·</span>
            <code class="text-xs font-mono text-th-text-m">{{ keyPrefix }}••••{{ keySuffix }}</code>
          </div>
        </div>
      </div>
      <button
        class="shrink-0 px-3 py-1.5 rounded-md text-xs font-medium text-th-text-t hover:text-danger hover:bg-danger/8 opacity-0 group-hover:opacity-100 transition-all duration-200"
        @click="emit('delete', id)"
      >
        {{ t('dashboard.unlink') }}
      </button>
    </div>

    <!-- Mobile: compact vertical layout -->
    <div class="sm:hidden px-4 py-3.5">
      <!-- Row 1: name + status badge -->
      <div class="flex items-center gap-2 mb-2">
        <svg class="shrink-0 w-4 h-4 text-highlight/60" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.75 17L9 20l-1 1h8l-1-1-.75-3M3 13h18M5 17h14a2 2 0 002-2V5a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
        </svg>
        <span class="text-sm font-medium text-th-text truncate">{{ displayName }}</span>
        <span
          class="shrink-0 flex items-center gap-1 text-[11px] px-1.5 py-0.5 rounded-full font-medium leading-none"
          :class="lastUsed ? 'bg-success/10 text-success' : 'bg-th-text/[0.04] text-th-text-m'"
        >
          <span class="w-1.5 h-1.5 rounded-full" :class="lastUsed ? 'bg-success' : 'bg-th-text-m'" />
          {{ lastUsed ? t('dashboard.active') : t('dashboard.never') }}
        </span>
      </div>

      <!-- Row 2: metadata + unlink -->
      <div class="flex items-center justify-between">
        <div class="flex items-center gap-1.5 text-xs text-th-text-t flex-wrap">
          <span class="whitespace-nowrap">{{ formatDate(createdAt) }}</span>
          <template v-if="lastUsed">
            <span class="text-th-text-m">·</span>
            <span class="whitespace-nowrap">{{ t('dashboard.lastUsed') }} {{ formatDate(lastUsed) }}</span>
          </template>
          <span class="text-th-text-m">·</span>
          <code class="text-xs font-mono text-th-text-m whitespace-nowrap">{{ keyPrefix }}••••{{ keySuffix }}</code>
        </div>
        <button
          class="shrink-0 ml-3 text-xs font-medium text-th-text-t hover:text-danger transition-colors"
          @click="emit('delete', id)"
        >
          {{ t('dashboard.unlink') }}
        </button>
      </div>
    </div>
  </div>
</template>
