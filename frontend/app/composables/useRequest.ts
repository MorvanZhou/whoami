export const useRequest = () => {
  const config = useRuntimeConfig()

  const fetchData = <T>(path: string, options: Record<string, any> = {}) => {
    return $fetch<T>(`${config.public.apiBase}${path}`, {
      ...options,
      credentials: 'include',
    })
  }

  return { fetchData }
}
