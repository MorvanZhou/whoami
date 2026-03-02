<script setup lang="ts">
const { t } = useI18n()
const localePath = useLocalePath()
</script>

<template>
  <section class="relative min-h-[85vh] flex items-center justify-center overflow-hidden">
    <!-- Animated background grid -->
    <div class="absolute inset-0">
      <div class="absolute inset-0 bg-[radial-gradient(ellipse_at_center,_rgba(14,165,233,0.08)_0%,_transparent_70%)]" />
      <svg class="absolute inset-0 w-full h-full opacity-[0.03] text-th-text" xmlns="http://www.w3.org/2000/svg">
        <defs>
          <pattern id="grid" width="60" height="60" patternUnits="userSpaceOnUse">
            <path d="M 60 0 L 0 0 0 60" fill="none" stroke="currentColor" stroke-width="0.5" />
          </pattern>
        </defs>
        <rect width="100%" height="100%" fill="url(#grid)" />
      </svg>
      <!-- Floating glow orbs -->
      <div class="absolute top-1/4 left-1/4 w-96 h-96 bg-accent/20 rounded-full animate-glow-pulse" />
      <div class="absolute bottom-1/4 right-1/4 w-72 h-72 bg-accent-light/15 rounded-full animate-glow-pulse" style="animation-delay: 1.5s" />
    </div>

    <!-- Content -->
    <div class="relative z-10 text-center px-6 max-w-4xl mx-auto">

      <h1 class="text-5xl sm:text-6xl lg:text-7xl font-extrabold tracking-tight leading-[1.1] mb-6 animate-fade-in-up" style="animation-delay: 0.1s">
        <span class="text-th-text">{{ t('hero.title') }}</span>
        <br>
        <span class="bg-gradient-to-r from-accent via-accent-light to-accent-glow bg-clip-text text-transparent">
          {{ t('hero.titleHighlight') }}
        </span>
      </h1>

      <p class="text-lg sm:text-xl text-th-text-s max-w-2xl mx-auto mb-10 animate-fade-in-up" style="animation-delay: 0.2s">
        {{ t('hero.subtitle') }}
      </p>

      <div class="animate-fade-in-up" style="animation-delay: 0.3s">
        <NuxtLink
          :to="`${localePath('/login')}?redirect=dashboard`"
          class="relative inline-flex items-center gap-2 px-8 py-4 rounded-xl text-base font-bold text-white overflow-hidden group"
        >
          <span class="absolute inset-0 bg-gradient-to-r from-accent to-accent-light" />
          <span class="absolute inset-0 bg-gradient-to-r from-accent-light to-accent opacity-0 group-hover:opacity-100 transition-opacity duration-300" />
          <span class="absolute inset-0 opacity-0 group-hover:opacity-100 transition-opacity duration-500 bg-[radial-gradient(circle_at_center,_rgba(255,255,255,0.15)_0%,_transparent_70%)]" />
          <span class="relative flex items-center gap-2">
            {{ t('hero.cta') }}
            <svg class="w-5 h-5 group-hover:translate-x-0.5 transition-transform" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7l5 5m0 0l-5 5m5-5H6" />
            </svg>
          </span>
        </NuxtLink>
      </div>
    </div>
  </section>
</template>
