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
    <!-- Nav -->
    <nav aria-label="Main navigation" class="fixed top-0 inset-x-0 z-50 border-b border-th-text/[var(--color-border-opacity)] bg-th-nav/80 backdrop-blur-xl">
      <div class="mx-auto max-w-6xl px-3 sm:px-6 h-14 sm:h-16 flex items-center justify-between">
        <!-- Logo -->
        <NuxtLink :to="localePath('/')" class="flex items-center gap-2.5 group">
          <img src="/logo.png" alt="whoami" class="w-7 h-7" />
          <span class="text-accent font-mono font-bold text-lg tracking-tight group-hover:text-accent-light transition-colors">
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
                class="flex items-center gap-1.5 sm:gap-2 rounded-lg px-2 sm:px-3 py-1.5 hover:bg-th-bg-t transition-colors"
                @click="showUserMenu = !showUserMenu"
              >
                <img
                  v-if="user?.avatar_url"
                  :src="user.avatar_url"
                  :alt="user.name || ''"
                  class="w-7 h-7 rounded-full ring-2 ring-accent/30"
                >
                <span class="text-sm text-th-text-s hidden sm:inline">{{ user?.name }}</span>
                <svg class="w-4 h-4 text-th-text-t" fill="none" stroke="currentColor" viewBox="0 0 24 24">
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
                  class="absolute right-0 mt-2 w-48 rounded-xl bg-th-card border border-th-text/[0.08] shadow-2xl overflow-hidden"
                >
                  <NuxtLink
                    :to="localePath('/dashboard')"
                    class="block px-4 py-2.5 text-sm text-th-text-s hover:bg-th-bg-t hover:text-th-text transition-colors"
                    @click="closeMenu"
                  >
                    {{ t('nav.dashboard') }}
                  </NuxtLink>
                  <button
                    class="block w-full text-left px-4 py-2.5 text-sm text-th-text-s hover:bg-th-bg-t hover:text-th-text transition-colors"
                    @click="logout"
                  >
                    {{ t('nav.logout') }}
                  </button>
                </div>
              </Transition>
            </div>
          </template>

          <template v-else>
            <NuxtLink
              :to="localePath('/login')"
              class="relative inline-flex items-center px-3 sm:px-5 py-1.5 sm:py-2 rounded-lg text-xs sm:text-sm font-semibold text-white overflow-hidden group"
            >
              <span class="absolute inset-0 bg-gradient-to-r from-accent to-accent-light group-hover:from-accent-light group-hover:to-accent transition-all duration-300" />
              <span class="relative">{{ t('nav.getStarted') }}</span>
            </NuxtLink>
          </template>
        </div>
      </div>
    </nav>

    <!-- Main -->
    <main class="pt-14 sm:pt-16">
      <slot />
    </main>

    <!-- Footer -->
    <footer aria-label="Site footer" class="border-t border-th-text/[var(--color-footer-border-opacity)] py-8 mt-20">
      <div class="mx-auto max-w-6xl px-6 flex flex-col sm:flex-row items-center justify-between gap-4 text-sm text-th-text-t">
        <span>{{ t('footer.rights', { year: new Date().getFullYear() }) }}</span>
        <div class="flex items-center gap-6">
          <a href="https://github.com/MorvanZhou/whoami" target="_blank" rel="noopener noreferrer" class="hover:text-th-text-s transition-colors">
            {{ t('footer.github') }}
          </a>
        </div>
      </div>
    </footer>
  </div>
</template>
