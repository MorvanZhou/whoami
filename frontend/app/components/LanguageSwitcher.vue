<script setup lang="ts">
const { locale } = useI18n()
const switchLocalePath = useSwitchLocalePath()

const toggleLocale = () => {
  const next = locale.value === 'en' ? 'zh' : 'en'
  navigateTo(switchLocalePath(next))
}
</script>

<template>
  <button
    class="flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-xs font-mono font-medium text-th-text-s hover:text-th-text hover:bg-th-bg-t border border-th-text/[0.08] transition-all duration-200"
    @click="toggleLocale"
  >
    <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 01-9 9m9-9a9 9 0 00-9-9m9 9H3m9 9a9 9 0 01-9-9m9 9c1.657 0 3-4.03 3-9s-1.343-9-3-9m0 18c-1.657 0-3-4.03-3-9s1.343-9 3-9m-9 9a9 9 0 019-9" />
    </svg>
    {{ locale === 'en' ? '中文' : 'EN' }}
  </button>
</template>
