export default defineNuxtRouteMiddleware(async (to) => {
  // Skip auth check during SSR — HttpOnly cookies are not forwarded
  // and devProxy does not handle server-side requests.
  if (import.meta.server) return

  const { isLoggedIn, fetchUser, user } = useAuth()
  const localePath = useLocalePath()

  // If user state is not loaded yet, try fetching (browser will send HttpOnly cookie)
  if (!user.value) {
    await fetchUser()
  }

  if (!isLoggedIn.value) {
    return navigateTo(`${localePath('/login')}?redirect=${encodeURIComponent(to.fullPath)}`)
  }
})
