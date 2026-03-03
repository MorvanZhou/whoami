<script setup lang="ts">
const { t } = useI18n()

useSeoMeta({
  title: () => t('seo.home.title'),
  description: () => t('seo.home.description'),
  ogTitle: () => t('seo.home.title'),
  ogDescription: () => t('seo.home.description'),
  twitterTitle: () => t('seo.home.title'),
  twitterDescription: () => t('seo.home.description'),
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
