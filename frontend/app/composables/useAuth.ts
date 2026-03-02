interface User {
  id: string
  email: string | null
  name: string | null
  avatar_url: string | null
  provider: string
  created_at: string | null
}

export const useAuth = () => {
  const token = useCookie('whoami_token', { maxAge: 60 * 60 * 24 * 7 })
  const user = useState<User | null>('auth_user', () => null)
  const loading = useState('auth_loading', () => false)
  const config = useRuntimeConfig()

  const isLoggedIn = computed(() => !!token.value && !!user.value)

  const fetchUser = async () => {
    if (!token.value) return
    loading.value = true
    try {
      const data = await $fetch<User>(`${config.public.apiBase}/auth/me`, {
        headers: { Authorization: `Bearer ${token.value}` },
      })
      user.value = data
    } catch {
      token.value = null
      user.value = null
    } finally {
      loading.value = false
    }
  }

  const login = (provider: 'github' | 'google', redirect?: string) => {
    const params = new URLSearchParams()
    if (redirect) params.set('redirect', redirect)
    window.location.href = `${config.public.apiBase}/auth/${provider}/login?${params.toString()}`
  }

  const setToken = (newToken: string) => {
    token.value = newToken
  }

  const logout = async () => {
    const localePath = useLocalePath()
    try {
      await $fetch(`${config.public.apiBase}/auth/logout`, {
        method: 'POST',
        headers: { Authorization: `Bearer ${token.value}` },
      })
    } catch {
      // ignore
    }
    token.value = null
    user.value = null
    navigateTo(localePath('/'))
  }

  return {
    token: readonly(token),
    user: readonly(user),
    loading: readonly(loading),
    isLoggedIn,
    fetchUser,
    login,
    setToken,
    logout,
  }
}
