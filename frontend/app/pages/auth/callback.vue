<script setup lang="ts">
const route = useRoute()
const { fetchUser } = useAuth()
const { t } = useI18n()
const localePath = useLocalePath()

const error = ref(false)

onMounted(async () => {
  const rawRedirect = (route.query.redirect as string) || '/dashboard'
  const redirect = rawRedirect.startsWith('/') ? rawRedirect : `/${rawRedirect}`

  try {
    // JWT is now set as HttpOnly cookie by the backend callback;
    // just fetch user to validate the session.
    await fetchUser()
    navigateTo(localePath(redirect))
  } catch {
    error.value = true
  }
})
</script>

<template>
  <div class="min-h-[calc(100vh-4rem)] flex items-center justify-center">
    <div v-if="error" class="text-center">
      <p class="text-danger text-lg font-medium">{{ t('common.error') }}</p>
      <NuxtLink :to="localePath('/login')" class="text-accent mt-4 inline-block hover:underline">
        {{ t('common.back') }}
      </NuxtLink>
    </div>
    <div v-else class="text-center">
      <div class="w-8 h-8 border-2 border-accent border-t-transparent rounded-full animate-spin mx-auto mb-4" />
      <p class="text-th-text-s">{{ t('common.loading') }}</p>
    </div>
  </div>
</template>
