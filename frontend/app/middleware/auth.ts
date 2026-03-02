export default defineNuxtRouteMiddleware((to) => {
  const token = useCookie('whoami_token')
  const localePath = useLocalePath()
  if (!token.value) {
    return navigateTo(`${localePath('/login')}?redirect=${encodeURIComponent(to.fullPath)}`)
  }
})
