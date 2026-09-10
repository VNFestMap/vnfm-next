interface ApiEnvelope<T> {
  code: number
  message: string
  data: T
}

type ApiMethod = 'GET' | 'POST' | 'PATCH' | 'PUT' | 'DELETE'

export const apiFetch = async <T>(
  path: string,
  options: { method?: ApiMethod; body?: Record<string, unknown> } = {}
): Promise<T> => {
  const config = useRuntimeConfig()
  const resp = await $fetch<ApiEnvelope<T>>(`${config.public.apiBase}${path}`, {
    credentials: 'include',
    method: options.method || 'GET',
    body: options.body
  })
  if (!resp || resp.code !== 0) {
    throw new Error(resp?.message || '请求失败')
  }
  return resp.data
}
