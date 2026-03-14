<script setup lang="ts">
const { t } = useI18n()

const title = computed(() => t('seo.home.title'))
const description = computed(() => t('seo.home.description'))

useSeoMeta({
  title,
  description,
  ogTitle: title,
  ogDescription: description,
  twitterTitle: title,
  twitterDescription: description,
})
</script>

<template>
  <div>
    <HomeHeroSection />
    <HomeFeaturesSection />
    <HomeFlowAnimation />
    <HomeInstallGuide />
  </div>
</template>
