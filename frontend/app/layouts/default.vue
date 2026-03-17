<script setup lang="ts">
const { t } = useI18n()
const localePath = useLocalePath()
const { isLoggedIn, user, logout, fetchUser } = useAuth()
const { initSystemListener } = useTheme()
const showUserMenu = ref(false)

onMounted(async () => {
  initSystemListener()
  if (!user.value) {
    await fetchUser()
  }
})

const closeMenu = () => {
  showUserMenu.value = false
}
</script>

<template>
  <div class="min-h-screen bg-th-bg text-th-text font-sans" @click="closeMenu">
    <!-- Skip link for keyboard navigation (a11y) -->
    <a href="#main-content" class="skip-link">
      Skip to main content
    </a>

    <!-- Nav -->
    <nav aria-label="Main navigation" class="fixed top-0 inset-x-0 z-50 border-b border-th-text/[var(--color-border-opacity)] bg-th-nav/80 backdrop-blur-xl">
      <div class="mx-auto max-w-6xl px-3 sm:px-6 h-14 sm:h-16 flex items-center justify-between">
        <!-- Logo -->
        <NuxtLink :to="localePath('/')" class="flex items-center gap-2.5 group" aria-label="whoami home">
          <img src="/logo.png" alt="" class="w-7 h-7" aria-hidden="true" />
          <span class="text-accent font-mono font-bold text-lg tracking-tight group-hover:text-accent-light transition-colors duration-200">
            whoami
          </span>
        </NuxtLink>

        <!-- Right -->
        <div class="flex items-center gap-1.5 sm:gap-3">
          <ThemeSwitcher />
          <LanguageSwitcher />

          <template v-if="isLoggedIn">
            <div class="relative" @click.stop>
              <button
                class="flex items-center gap-1.5 sm:gap-2 rounded-lg px-2 sm:px-3 py-1.5 hover:bg-th-bg-t transition-colors duration-200"
                :aria-expanded="showUserMenu"
                aria-haspopup="true"
                @click="showUserMenu = !showUserMenu"
              >
                <img
                  v-if="user?.avatar_url"
                  :src="user.avatar_url"
                  :alt="user.name || 'User avatar'"
                  class="w-7 h-7 rounded-full ring-2 ring-accent/30"
                >
                <span class="text-sm text-th-text-s hidden sm:inline">{{ user?.name }}</span>
                <svg class="w-4 h-4 text-th-text-t transition-transform duration-200" :class="{ 'rotate-180': showUserMenu }" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
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
                  v-if="showUserMenu"
                  role="menu"
                  class="absolute right-0 mt-2 w-48 rounded-xl bg-th-card border border-th-text/[0.08] shadow-2xl overflow-hidden"
                >
                  <NuxtLink
                    :to="localePath('/dashboard')"
                    role="menuitem"
                    class="flex items-center gap-2.5 px-4 py-2.5 text-sm text-th-text-s hover:bg-th-bg-t hover:text-th-text transition-colors duration-150"
                    @click="closeMenu"
                  >
                    <svg class="w-4 h-4 opacity-50" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6zM14 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2v-2zM14 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z" />
                    </svg>
                    {{ t('nav.dashboard') }}
                  </NuxtLink>
                  <div class="h-px bg-th-text/[0.04] mx-2" />
                  <button
                    role="menuitem"
                    class="flex items-center gap-2.5 w-full text-left px-4 py-2.5 text-sm text-th-text-s hover:bg-th-bg-t hover:text-th-text transition-colors duration-150"
                    @click="logout"
                  >
                    <svg class="w-4 h-4 opacity-50" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
                    </svg>
                    {{ t('nav.logout') }}
                  </button>
                </div>
              </Transition>
            </div>
          </template>

          <template v-else>
            <NuxtLink
              :to="localePath('/login')"
              class="relative inline-flex items-center px-3 sm:px-5 py-1.5 sm:py-2 rounded-lg text-xs sm:text-sm font-semibold text-white overflow-hidden group active:scale-[0.97] transition-transform duration-150"
            >
              <span class="absolute inset-0 bg-gradient-to-r from-accent to-accent-light group-hover:from-accent-light group-hover:to-accent transition-all duration-300" />
              <span class="relative">{{ t('nav.getStarted') }}</span>
            </NuxtLink>
          </template>
        </div>
      </div>
    </nav>

    <!-- Main -->
    <main id="main-content" class="pt-14 sm:pt-16">
      <slot />
    </main>

    <!-- Footer -->
    <footer aria-label="Site footer" class="border-t border-th-text/[var(--color-footer-border-opacity)] py-8 mt-20">
      <div class="mx-auto max-w-6xl px-6 flex flex-col sm:flex-row items-center justify-between gap-4 text-sm text-th-text-t">
        <span>{{ t('footer.rights', { year: new Date().getFullYear() }) }}</span>
        <nav aria-label="Footer navigation" class="flex items-center gap-6">
          <NuxtLink :to="localePath('/terms')" class="hover:text-th-text-s transition-colors duration-200">
            {{ t('footer.terms') }}
          </NuxtLink>
          <NuxtLink :to="localePath('/privacy')" class="hover:text-th-text-s transition-colors duration-200">
            {{ t('footer.privacy') }}
          </NuxtLink>
          <a href="https://github.com/MorvanZhou/whoami" target="_blank" rel="noopener noreferrer" class="hover:text-th-text-s transition-colors duration-200 inline-flex items-center gap-1.5">
            <svg class="w-4 h-4" viewBox="0 0 24 24" fill="currentColor" aria-hidden="true">
              <path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z" />
            </svg>
            {{ t('footer.github') }}
          </a>
        </nav>
      </div>
    </footer>
  </div>
</template>
