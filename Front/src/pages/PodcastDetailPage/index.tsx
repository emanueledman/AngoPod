import { useEffect, useState } from "react"
import { Link, useParams } from "react-router-dom"
import Navbar from "../../components/layout/Navbar"
import { useAudio } from "../../context/AudioContext"
import { getEpisodes, getPodcastBySlug } from "../../services/podcastService"
import SiteFooter from "../../SiteFooter"
import { Episode, Podcast } from "../../types"

export default function PodcastDetailPage() {
  const { slug } = useParams<{ slug: string }>()
  const { playEpisode, isPlaying, currentEpisode } = useAudio()
  const [followed, setFollowed] = useState(false)
  const [podcast, setPodcast] = useState<Podcast | null>(null)
  const [episodeList, setEpisodeList] = useState<Episode[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    let isMounted = true
    const currentSlug = slug || "fly-podcast"

    async function loadPodcastData() {
      setLoading(true)
      try {
        const [pod, eps] = await Promise.all([
          getPodcastBySlug(currentSlug),
          getEpisodes(currentSlug),
        ])

        if (isMounted) {
          if (pod) setPodcast(pod)
          setEpisodeList(eps)
          setLoading(false)
        }
      } catch (err) {
        console.error("Erro ao carregar podcast da API:", err)
        if (isMounted) setLoading(false)
      }
    }

    loadPodcastData()

    return () => {
      isMounted = false
    }
  }, [slug])

  if (loading && !podcast) {
    return (
      <div className="detail-page">
        <Navbar variant="detail" backLink={{ to: "/explorar", label: "← Voltar a explorar" }} />
        <main style={{ padding: "4rem 2rem", textAlign: "center" }}>
          <h2>A carregar podcast...</h2>
        </main>
        <SiteFooter />
      </div>
    )
  }

  if (!podcast) {
    return (
      <div className="detail-page">
        <Navbar variant="detail" backLink={{ to: "/explorar", label: "← Voltar a explorar" }} />
        <main style={{ padding: "4rem 2rem", textAlign: "center" }}>
          <h2>Podcast não encontrado</h2>
          <Link to="/explorar">Voltar ao mapa de podcasts</Link>
        </main>
        <SiteFooter />
      </div>
    )
  }

  const handlePlay = (episodeToPlay: Episode) => {
    playEpisode(episodeToPlay)
  }

  const presenterInitials = podcast.host
    .split(" ")
    .map((name) => name[0])
    .join("")
    .slice(0, 2)

  return (
    <div className="detail-page">
      <Navbar
        variant="detail"
        backLink={{ to: "/explorar", label: "← Voltar a explorar" }}
      />

      <main>
        <section className="podcast-detail-hero">
          <div className="podcast-art">
            <img src={podcast.image} alt={`Capa do ${podcast.title}`} referrerPolicy="no-referrer" />
            <span>#{podcast.rank} esta semana</span>
          </div>
          <div className="podcast-detail-copy">
            <span className="section-kicker">
              {podcast.category} · {podcast.province}
            </span>
            <h1>
              {podcast.title.split(" ").slice(0, -1).join(" ") || podcast.title}
              <br />
              {podcast.title.split(" ").slice(-1)}
            </h1>
            <p>{podcast.description}</p>

            {podcast.slug === "fly-podcast" ? (
              <Link className="presenter-inline" to="/apresentador/flay-squad">
                <span>FS</span>
                <div>
                  <small>Apresentado por</small>
                  <strong>{podcast.host} →</strong>
                </div>
              </Link>
            ) : (
              <div className="presenter-inline">
                <span>{presenterInitials}</span>
                <div>
                  <small>Apresentado por</small>
                  <strong>{podcast.host}</strong>
                </div>
              </div>
            )}

            <div className="detail-actions">
              {episodeList.length > 0 && (
                <button
                  type="button"
                  onClick={() => handlePlay(episodeList[0])}
                >
                  ▶ Ouvir último episódio
                </button>
              )}
              <button
                className={followed ? "followed" : ""}
                onClick={() => setFollowed(!followed)}
                type="button"
              >
                {followed ? "A seguir" : "+ Seguir podcast"}
              </button>
            </div>
            <div className="podcast-numbers">
              <div>
                <strong>{podcast.listeners}</strong>
                <span>ouvintes mensais</span>
              </div>
              <div>
                <strong>{podcast.episodesCount}</strong>
                <span>episódios</span>
              </div>
              <div>
                <strong>{podcast.rating}</strong>
                <span>avaliação</span>
              </div>
            </div>
          </div>
        </section>

        <section className="podcast-episodes">
          <div className="detail-section-title">
            <div>
              <span className="section-kicker">Todos os episódios</span>
              <h2>Conversas que ficam.</h2>
            </div>
            <span>Atualizado semanalmente</span>
          </div>
          <div className="detail-episode-list">
            {episodeList.map((episode) => {
              const isThisPlaying = isPlaying && currentEpisode?.id === episode.id

              return (
                <article
                  className={isThisPlaying ? "playing" : ""}
                  key={episode.id}
                >
                  <span>{episode.number || "EP."}</span>
                  <button
                    type="button"
                    onClick={() => handlePlay(episode)}
                  >
                    {isThisPlaying ? "⏸" : "▶"}
                  </button>
                  <div>
                    <h3>{episode.title}</h3>
                    <p>{episode.guest}</p>
                  </div>
                  <small>{episode.date}</small>
                  <small>{episode.duration || episode.time}</small>
                </article>
              )
            })}
          </div>
        </section>
      </main>

      <SiteFooter />
    </div>
  )
}
