<script setup lang="ts">
const { t } = useI18n()
const localePath = useLocalePath()
const route = useRoute()
const { isLoggedIn, fetchUser, token } = useAuth()

const rawRedirect = computed(() => (route.query.redirect as string) || '')
const redirect = computed(() => rawRedirect.value.replace(/^\/+/, ''))
const isFromAgent = computed(() => redirect.value === 'dashboard')

onMounted(async () => {
  if (token.value && !isLoggedIn.value) {
    await fetchUser()
  }
})

watchEffect(() => {
  if (isLoggedIn.value) {
    navigateTo(localePath(`/${redirect.value || 'dashboard'}`))
  }
})
</script>

<template>
  <div class="relative min-h-[calc(100vh-4rem)] flex items-center justify-center px-6">
    <!-- Background -->
    <div class="absolute inset-0">
      <div class="absolute inset-0 bg-[radial-gradient(ellipse_at_center,_rgba(14,165,233,0.06)_0%,_transparent_60%)]" />
      <svg class="absolute inset-0 w-full h-full opacity-[0.02]" xmlns="http://www.w3.org/2000/svg">
        <defs>
          <pattern id="loginGrid" width="40" height="40" patternUnits="userSpaceOnUse">
            <path d="M 40 0 L 0 0 0 40" fill="none" stroke="white" stroke-width="0.5" />
          </pattern>
        </defs>
        <rect width="100%" height="100%" fill="url(#loginGrid)" />
      </svg>
    </div>

    <!-- Login card -->
    <div class="relative z-10 w-full" :class="isFromAgent ? 'max-w-md' : 'max-w-sm'">
      <!-- AI context banner — only when opened by agent -->
      <Transition
        enter-active-class="transition duration-500 ease-out"
        enter-from-class="opacity-0 -translate-y-3"
        enter-to-class="opacity-100 translate-y-0"
      >
        <div v-if="isFromAgent" class="mb-5 p-5 rounded-2xl bg-accent/[0.04] border border-accent/15 backdrop-blur-sm">
          <!-- Badge -->
          <div class="flex items-center gap-2.5 mb-3">
            <span class="relative flex h-2.5 w-2.5">
              <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-accent opacity-60" />
              <span class="relative inline-flex rounded-full h-2.5 w-2.5 bg-accent" />
            </span>
            <span class="text-sm font-bold text-accent tracking-wide">{{ t('login.aiContext.badge') }}</span>
          </div>
          <!-- Explanation -->
          <h3 class="text-sm font-bold text-white/90 mb-1.5">{{ t('login.aiContext.title') }}</h3>
          <p class="text-xs leading-relaxed text-gray-400">{{ t('login.aiContext.desc') }}</p>
          <!-- Safety note -->
          <div class="flex items-center gap-1.5 mt-3 pt-3 border-t border-white/5">
            <svg class="w-3.5 h-3.5 text-success/70 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
            </svg>
            <span class="text-[11px] text-gray-500">{{ t('login.aiContext.safe') }}</span>
          </div>
        </div>
      </Transition>

      <div class="p-8 rounded-2xl bg-surface-800/60 border border-white/5 backdrop-blur-xl shadow-2xl">
        <!-- Logo -->
        <div class="text-center mb-8">
          <img src="/logo.png" alt="whoami" class="w-14 h-14 mx-auto mb-3" />
          <h1 class="text-accent font-mono font-bold text-2xl">{{ t('login.title') }}</h1>
          <p class="text-sm text-gray-400 mt-2">{{ t('login.subtitle') }}</p>
        </div>

        <!-- OAuth buttons -->
        <OAuthButtons :redirect="redirect || 'dashboard'" />

        <!-- Note -->
        <p class="text-center text-xs text-gray-500 mt-6">
          {{ t('login.note') }}
        </p>
      </div>
    </div>
  </div>
</template>
