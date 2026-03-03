<script setup lang="ts">
const { t, locale } = useI18n()

const introImage = computed(() =>
  locale.value === 'zh' ? '/intro-zh.png' : '/intro-en.png'
)
</script>

<template>
  <section class="py-24 px-6">
    <div class="max-w-3xl mx-auto">
      <h2 class="text-3xl sm:text-4xl font-bold text-center text-th-text mb-16">
        {{ t('install.title') }}
      </h2>

      <!-- Option A: Copy prompt for agent -->
      <div class="mb-10">
        <h3 class="text-lg font-semibold text-th-text mb-2 flex items-center gap-2">
          <span class="inline-flex items-center justify-center w-7 h-7 rounded-full bg-accent/20 text-accent text-sm font-bold font-mono">A</span>
          {{ t('install.step1Title') }}
        </h3>
        <p class="text-sm text-th-text-s mb-4 ml-9">{{ t('install.step1Desc') }}</p>
        <div class="ml-9">
          <CopyBlock
            :text="t('install.prompt')"
            label="AGENT PROMPT"
            prominent
          />
        </div>
      </div>

      <!-- Option B: Install via command (for developers) -->
      <div class="mb-16">
        <h3 class="text-lg font-semibold text-th-text mb-2 flex items-center gap-2">
          <span class="inline-flex items-center justify-center w-7 h-7 rounded-full bg-accent-light/20 text-accent-light text-sm font-bold font-mono">B</span>
          {{ t('install.step2Title') }}
        </h3>
        <p class="text-sm text-th-text-s mb-4 ml-9">{{ t('install.step2Desc') }}</p>
        <div class="ml-9">
          <CopyBlock
            text="npx skills add MorvanZhou/whoami"
            label="COMMAND"
            :mono="true"
            prominent
          />
        </div>
      </div>

      <!-- Intro image -->
      <div class="rounded-2xl overflow-hidden border border-th-text/[0.08] shadow-lg">
        <img
          :src="introImage"
          :alt="t('install.introAlt')"
          class="w-full h-auto"
          loading="lazy"
        />
      </div>
    </div>
  </section>
</template>
