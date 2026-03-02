<script setup lang="ts">
const { t } = useI18n()

const features = computed(() => [
  {
    icon: 'sync',
    title: t('features.crossSync.title'),
    desc: t('features.crossSync.desc'),
    gradient: 'from-accent to-accent-light',
  },
  {
    icon: 'shield',
    title: t('features.privacy.title'),
    desc: t('features.privacy.desc'),
    gradient: 'from-success to-emerald-400',
  },
  {
    icon: 'code',
    title: t('features.devFriendly.title'),
    desc: t('features.devFriendly.desc'),
    gradient: 'from-info to-cyan-400',
  },
])
</script>

<template>
  <section class="py-24 px-6">
    <div class="max-w-5xl mx-auto">
      <h2 class="text-3xl sm:text-4xl font-bold text-center text-white mb-16">
        {{ t('features.title') }}
      </h2>

      <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div
          v-for="(feature, i) in features"
          :key="i"
          class="group relative p-6 rounded-2xl bg-surface-800/50 border border-white/5 hover:border-white/15 backdrop-blur-sm transition-all duration-300 hover:-translate-y-1"
        >
          <!-- Hover glow -->
          <div class="absolute inset-0 rounded-2xl opacity-0 group-hover:opacity-100 transition-opacity duration-500 bg-gradient-to-br" :class="feature.gradient" style="filter: blur(40px); transform: scale(0.8); z-index: -1; opacity: 0.05" />

          <!-- Icon -->
          <div class="w-12 h-12 rounded-xl flex items-center justify-center mb-4 bg-gradient-to-br" :class="feature.gradient" style="opacity: 0.15">
            <svg v-if="feature.icon === 'sync'" class="w-6 h-6 text-accent-glow" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h5M20 20v-5h-5M4 9a8 8 0 0113.5-4.5L20 7M20 15a8 8 0 01-13.5 4.5L4 17" />
            </svg>
            <svg v-else-if="feature.icon === 'shield'" class="w-6 h-6 text-success" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
            </svg>
            <svg v-else class="w-6 h-6 text-info" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 20l4-16m4 4l4 4-4 4M6 16l-4-4 4-4" />
            </svg>
          </div>

          <h3 class="text-lg font-semibold text-white mb-2">{{ feature.title }}</h3>
          <p class="text-sm text-gray-400 leading-relaxed">{{ feature.desc }}</p>
        </div>
      </div>
    </div>
  </section>
</template>
