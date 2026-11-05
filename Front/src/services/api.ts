// Client HTTP para integração futura com o Back-end
export const API_BASE_URL = import.meta.env.VITE_API_URL || "http://localhost:8080/api"

export async function apiFetch<T>(endpoint: string, options?: RequestInit): Promise<T> {
  const response = await fetch(`${API_BASE_URL}${endpoint}`, {
    headers: {
      "Content-Type": "application/json",
      ...options?.headers,
    },
    ...options,
  })

  if (!response.ok) {
    throw new Error(`Erro na API (${response.status}): ${response.statusText}`)
  }

  return response.json()
}

