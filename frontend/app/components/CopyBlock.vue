<script setup lang="ts">
const props = defineProps<{
  text: string
  label?: string
  mono?: boolean
  prominent?: boolean
}>()

const { t } = useI18n()
const copied = ref(false)

const copy = async () => {
  try {
    await navigator.clipboard.writeText(props.text)
    copied.value = true
    setTimeout(() => { copied.value = false }, 2000)
  } catch {
    // fallback
    const el = document.createElement('textarea')
    el.value = props.text
    document.body.appendChild(el)
    el.select()
    document.execCommand('copy')
    document.body.removeChild(el)
    copied.value = true
    setTimeout(() => { copied.value = false }, 2000)
  }
}
</script>

<template>
  <div class="group relative rounded-xl bg-surface-800 border border-white/5 overflow-hidden">
    <!-- Label row: with prominent copy button inline -->
    <div v-if="label" class="flex items-center justify-between px-4 pt-3">
      <span class="text-xs font-medium text-gray-500 uppercase tracking-wider">{{ label }}</span>
      <button
        v-if="prominent"
        class="flex items-center gap-1.5 px-3 py-1 rounded-lg text-xs font-semibold transition-all duration-200"
        :class="copied
          ? 'bg-success/20 text-success border border-success/30'
          : 'bg-accent/10 text-accent hover:bg-accent/20 border border-accent/30 hover:border-accent/50'"
        @click="copy"
      >
        <svg v-if="!copied" class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
        </svg>
        <svg v-else class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
        </svg>
        {{ copied ? t('common.copied') : t('common.copy') }}
      </button>
    </div>
    <div class="p-4" :class="prominent ? '' : 'flex items-start gap-3'">
      <pre
        class="flex-1 text-sm leading-relaxed whitespace-pre-wrap break-all"
        :class="mono ? 'font-mono text-accent-glow' : 'text-gray-300'"
      >{{ text }}</pre>
      <!-- Default: inline button beside text -->
      <button
        v-if="!prominent"
        class="shrink-0 flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-xs font-medium transition-all duration-200"
        :class="copied
          ? 'bg-success/20 text-success border border-success/30'
          : 'bg-white/5 text-gray-400 hover:text-white hover:bg-white/10 border border-white/10'"
        @click="copy"
      >
        <svg v-if="!copied" class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
        </svg>
        <svg v-else class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
        </svg>
        {{ copied ? t('common.copied') : t('common.copy') }}
      </button>
    </div>
  </div>
</template>
