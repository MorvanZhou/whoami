export type ThemeMode = 'dark' | 'light' | 'system'

const systemDark = ref(false)

export const useTheme = () => {
  const mode = useCookie<ThemeMode>('whoami_theme', {
    default: () => 'system',
    maxAge: 60 * 60 * 24 * 365,
  })

  const isDark = computed(() => {
    if (mode.value === 'system') return systemDark.value
    return mode.value === 'dark'
  })

  const applyTheme = () => {
    if (import.meta.server) return
    document.documentElement.classList.toggle('dark', isDark.value)
  }

  watch(isDark, () => applyTheme())

  const setMode = (newMode: ThemeMode) => {
    mode.value = newMode
  }

  const initSystemListener = () => {
    if (import.meta.server) return
    const mq = window.matchMedia('(prefers-color-scheme: dark)')
    systemDark.value = mq.matches
    mq.addEventListener('change', (e) => {
      systemDark.value = e.matches
    })
    applyTheme()
  }

  return {
    mode: readonly(mode),
    isDark,
    setMode,
    initSystemListener,
  }
}
