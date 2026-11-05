import { Episode, MoverItem, Podcast, Presenter, Province, PulseHighlight } from "../types"
import { apiFetch } from "./api"

// ─── Podcasts ────────────────────────────────────────────────────────────────

export async function getPodcasts(filters?: {
  category?: string
  province?: string
  language?: string
  platform?: string
  search?: string
}): Promise<Podcast[]> {
  const query = new URLSearchParams()
  if (filters?.category) query.set("category", filters.category)
  if (filters?.province) query.set("province", filters.province)
  if (filters?.language) query.set("language", filters.language)
  if (filters?.platform) query.set("platform", filters.platform)
  if (filters?.search) query.set("search", filters.search)

  const queryString = query.toString() ? `?${query.toString()}` : ""
  return apiFetch<Podcast[]>(`/podcasts${queryString}`)
}

export async function getPodcastBySlug(slug: string): Promise<Podcast | undefined> {
  return apiFetch<Podcast>(`/podcasts/${slug}`)
}

// ─── Episódios ────────────────────────────────────────────────────────────────

export async function getEpisodes(podcastSlug?: string): Promise<Episode[]> {
  const endpoint = podcastSlug ? `/podcasts/${podcastSlug}/episodes` : "/episodes"
  return apiFetch<Episode[]>(endpoint)
}

export async function getEpisodeByID(id: string): Promise<Episode | undefined> {
  return apiFetch<Episode>(`/episodes/${id}`)
}

// ─── Apresentadores ───────────────────────────────────────────────────────────

export async function getAllPresenters(): Promise<Presenter[]> {
  return apiFetch<Presenter[]>("/presenters")
}

export async function getPresenterBySlug(slug: string): Promise<Presenter | undefined> {
  return apiFetch<Presenter>(`/presenters/${slug}`)
}

// ─── Províncias ───────────────────────────────────────────────────────────────

export async function getProvinces(): Promise<Province[]> {
  return apiFetch<Province[]>("/provinces")
}

// ─── Destaques e Ranking ──────────────────────────────────────────────────────

export async function getHighlights(): Promise<{
  pulseHighlights: PulseHighlight[]
  movers: MoverItem[]
}> {
  return apiFetch<{ pulseHighlights: PulseHighlight[]; movers: MoverItem[] }>("/highlights")
}

export async function getRanking(): Promise<Podcast[]> {
  return apiFetch<Podcast[]>("/ranking")
}
