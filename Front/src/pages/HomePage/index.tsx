import React, { useEffect, useState } from "react"
import Navbar from "../../components/layout/Navbar"
import { getEpisodes, getHighlights, getPodcasts, getProvinces } from "../../services/podcastService"
import SiteFooter from "../../SiteFooter"
import { Episode, MoverItem, Podcast, Province, PulseHighlight } from "../../types"
import DiscoverySection from "./components/DiscoverySection"
import Hero from "./components/Hero"
import HomePaths from "./components/HomePaths"
import PulseBar from "./components/PulseBar"
import RankingMapSection from "./components/RankingMapSection"
import RankMovers from "./components/RankMovers"
import RecentEpisodesSection from "./components/RecentEpisodesSection"

export default function HomePage() {
  const [query, setQuery] = useState("")
  const [activeCategory, setActiveCategory] = useState("Todos")

  // Estados com dados do Back-end (Go)
  const [podcasts, setPodcasts] = useState<Podcast[]>([])
  const [episodes, setEpisodes] = useState<Episode[]>([])
  const [provinces, setProvinces] = useState<Province[]>([])
  const [highlights, setHighlights] = useState<PulseHighlight[]>([])
  const [movers, setMovers] = useState<MoverItem[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [apiError, setApiError] = useState(false)

  useEffect(() => {
    let isMounted = true

    async function loadData() {
      setApiError(false)
      try {
        const [podList, epList, provList, highData] = await Promise.all([
          getPodcasts(),
          getEpisodes(),
          getProvinces(),
          getHighlights(),
        ])

        if (isMounted) {
          setPodcasts(podList)
          setEpisodes(epList)
          setProvinces(provList)
          setHighlights(highData.pulseHighlights)
          setMovers(highData.movers)
          setIsLoading(false)
        }
      } catch (err) {
        console.error("Back-end inacessível:", err)
        if (isMounted) {
          setApiError(true)
          setIsLoading(false)
        }
      }
    }

    loadData()

    return () => {
      isMounted = false
    }
  }, [])

  useEffect(() => {
    const sections = document.querySelectorAll<HTMLElement>(".app-shell main > section")
    document.documentElement.classList.add("reveal-ready")

    const observer = new IntersectionObserver(
      (entries) => {
        entries.forEach((entry) => {
          if (entry.isIntersecting) {
            entry.target.classList.add("section-in-view")
            observer.unobserve(entry.target)
          }
        })
      },
      { rootMargin: "0px 0px -12% 0px", threshold: 0.08 }
    )

    sections.forEach((section) => observer.observe(section))

    return () => {
      observer.disconnect()
      document.documentElement.classList.remove("reveal-ready")
    }
  }, [isLoading])

  const handleSearch = (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault()
    document.querySelector("#descobrir")?.scrollIntoView({ behavior: "smooth" })
  }

  const featuredPodcast = podcasts.length > 0 ? podcasts[0] : undefined
  const featuredEpisode = episodes.length > 0 ? episodes[0] : undefined

  return (
    <div className="app-shell">
      <Navbar variant="home" />

      {apiError && (
        <div style={{
          background: "#1a0a0a",
          borderBottom: "1px solid #ff4444",
          padding: "0.75rem 2rem",
          textAlign: "center",
          fontSize: "0.875rem",
          color: "#ff8888",
        }}>
          ⚠️ Servidor indisponível — inicia o back-end Go para ver os conteúdos reais.{" "}
          <code style={{ color: "#ffbbbb" }}>cd Back && ./bin/api</code>
        </div>
      )}

      <main id="inicio">
        <Hero
          query={query}
          setQuery={setQuery}
          onSearch={handleSearch}
          featuredPodcast={featuredPodcast}
          featuredEpisode={featuredEpisode}
        />
        <PulseBar highlights={highlights} />
        <HomePaths />
        <DiscoverySection
          query={query}
          setQuery={setQuery}
          activeCategory={activeCategory}
          setActiveCategory={setActiveCategory}
          podcasts={podcasts}
          episodes={episodes}
        />
        <RankMovers movers={movers} />
        <RankingMapSection
          podcasts={podcasts}
          provinces={provinces}
          episodes={episodes}
        />
        <RecentEpisodesSection episodes={episodes} />
      </main>

      <SiteFooter />
    </div>
  )
}
