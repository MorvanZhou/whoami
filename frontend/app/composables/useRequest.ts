export const useRequest = () => {
  const config = useRuntimeConfig()
  const token = useCookie('whoami_token')

  const fetchData = <T>(path: string, options: Record<string, any> = {}) => {
    const headers: Record<string, string> = {
      ...((options.headers as Record<string, string>) || {}),
    }
    if (token.value) {
      headers.Authorization = `Bearer ${token.value}`
    }
    return $fetch<T>(`${config.public.apiBase}${path}`, {
      ...options,
      headers,
    })
  }

  return { fetchData }
}
