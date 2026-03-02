<script setup lang="ts">
import type { ThemeMode } from '~/composables/useTheme'

const { t } = useI18n()
const { mode, setMode } = useTheme()
const showMenu = ref(false)
const wrapper = ref<HTMLElement>()

const options: { value: ThemeMode; icon: string }[] = [
  { value: 'light', icon: 'sun' },
  { value: 'dark', icon: 'moon' },
  { value: 'system', icon: 'monitor' },
]

const currentIcon = computed(() => {
  return options.find(o => o.value === mode.value)?.icon || 'monitor'
})

const select = (value: ThemeMode) => {
  setMode(value)
  showMenu.value = false
}

const onClickOutside = (e: MouseEvent) => {
  if (wrapper.value && !wrapper.value.contains(e.target as Node)) {
    showMenu.value = false
  }
}

onMounted(() => document.addEventListener('click', onClickOutside))
onUnmounted(() => document.removeEventListener('click', onClickOutside))
</script>

<template>
  <div ref="wrapper" class="relative">
    <button
      class="flex items-center justify-center w-8 h-8 rounded-lg text-th-text-s hover:text-th-text hover:bg-th-bg-t transition-colors"
      :title="t('theme.title')"
      @click="showMenu = !showMenu"
    >
      <!-- Sun -->
      <svg v-if="currentIcon === 'sun'" class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z" />
      </svg>
      <!-- Moon -->
      <svg v-else-if="currentIcon === 'moon'" class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z" />
      </svg>
      <!-- Monitor / System -->
      <svg v-else class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.75 17L9 20l-1 1h8l-1-1-.75-3M3 13h18M5 17h14a2 2 0 002-2V5a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
      </svg>
    </button>

    <Transition
      enter-active-class="transition duration-150 ease-out"
      enter-from-class="opacity-0 scale-95 -translate-y-1"
      enter-to-class="opacity-100 scale-100 translate-y-0"
      leave-active-class="transition duration-100 ease-in"
      leave-from-class="opacity-100 scale-100 translate-y-0"
      leave-to-class="opacity-0 scale-95 -translate-y-1"
    >
      <div
        v-if="showMenu"
        class="absolute right-0 mt-2 w-36 rounded-xl bg-th-card border border-th-text/[0.08] shadow-2xl overflow-hidden z-50"
      >
        <button
          v-for="opt in options"
          :key="opt.value"
          class="flex items-center gap-2.5 w-full px-3.5 py-2.5 text-xs transition-colors"
          :class="mode === opt.value
            ? 'text-accent bg-accent/10 font-semibold'
            : 'text-th-text-s hover:text-th-text hover:bg-th-bg-t'"
          @click="select(opt.value)"
        >
          <!-- Sun -->
          <svg v-if="opt.icon === 'sun'" class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z" />
          </svg>
          <!-- Moon -->
          <svg v-else-if="opt.icon === 'moon'" class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z" />
          </svg>
          <!-- Monitor -->
          <svg v-else class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.75 17L9 20l-1 1h8l-1-1-.75-3M3 13h18M5 17h14a2 2 0 002-2V5a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
          </svg>
          {{ t(`theme.${opt.value}`) }}
        </button>
      </div>
    </Transition>
  </div>
</template>
