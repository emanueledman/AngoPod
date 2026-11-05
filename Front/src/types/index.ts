export type IconName =
  | "arrow"
  | "bookmark"
  | "chevron"
  | "clock"
  | "headphones"
  | "menu"
  | "pause"
  | "play"
  | "search"
  | "sparkle"
  | "trending"
  | "volume"

export interface Episode {
  id: string
  number?: string
  title: string
  podcastTitle: string
  podcastSlug: string
  guest?: string
  show?: string
  time?: string
  duration?: string
  date?: string
  listens?: string
  trend?: string
  image?: string
  tag?: string
  audioUrl?: string
}

export interface Podcast {
  id: string
  slug: string
  rank: number
  title: string
  host: string
  category: string
  province: string
  language: string
  platform: string
  format: string
  popularity: number
  publishedAt: string
  status: "Atualizado" | "Novo" | string
  image: string
  color: "coral" | "purple" | "yellow" | "green" | "orange" | "violet" | "lime" | string
  description?: string
  listeners?: string
  episodesCount?: number
  rating?: number
}

export interface Presenter {
  id: string
  slug: string
  name: string
  role: string
  quote: string
  bio: string
  story1: string
  story2: string
  photo: string
  badgeNumber: string
  followers: string
  listeners: string
  conversations: number
  podcastSlug: string
  podcastTitle: string
}

export interface Province {
  name: string
  count: number
  x: number
  y: number
  description?: string
}

export interface PulseHighlight {
  label: string
  title: string
  href: string
}

export interface MoverItem {
  type: "up" | "new" | "share"
  badgeLabel: string
  value: string
  title: string
  subtitle: string
  image: string
  linkHref: string
  linkLabel: string
}

