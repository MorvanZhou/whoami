interface User {
  id: string
  email: string | null
  name: string | null
  avatar_url: string | null
  provider: string
  created_at: string | null
}

export const useAuth = () => {
  const user = useState<User | null>('auth_user', () => null)
  const loading = useState('auth_loading', () => false)
  const config = useRuntimeConfig()

  const isLoggedIn = computed(() => !!user.value)

  let _fetchPromise: Promise<void> | null = null

  const fetchUser = () => {
    if (_fetchPromise) return _fetchPromise
    _fetchPromise = (async () => {
      loading.value = true
      try {
        const data = await $fetch<User>(`${config.public.apiBase}/auth/me`, {
          credentials: 'include',
        })
        user.value = data
      } catch {
        user.value = null
      } finally {
        loading.value = false
        _fetchPromise = null
      }
    })()
    return _fetchPromise
  }

  const login = (provider: 'github' | 'google', redirect?: string) => {
    const params = new URLSearchParams()
    if (redirect) params.set('redirect', redirect)
    window.location.href = `${config.public.apiBase}/auth/${provider}/login?${params.toString()}`
  }

  const logout = async () => {
    const localePath = useLocalePath()
    try {
      await $fetch(`${config.public.apiBase}/auth/logout`, {
        method: 'POST',
        credentials: 'include',
      })
    } catch {
      // ignore
    }
    user.value = null
    navigateTo(localePath('/'))
  }

  return {
    user: readonly(user),
    loading: readonly(loading),
    isLoggedIn,
    fetchUser,
    login,
    logout,
  }
}
