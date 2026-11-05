import { useState } from "react"
import { Link } from "react-router-dom"
import Icon from "../../../components/common/Icon"
import { useAudio } from "../../../context/AudioContext"
import { Episode, Podcast } from "../../../types"

const categories = ["Todos", "Cultura", "Música", "Negócios", "Tecnologia", "Sociedade", "Notícias"]

interface DiscoverySectionProps {
  query: string
  setQuery: (q: string) => void
  activeCategory: string
  setActiveCategory: (cat: string) => void
  podcasts?: Podcast[]
  episodes?: Episode[]
}

export default function DiscoverySection({
  query,
  setQuery,
  activeCategory,
  setActiveCategory,
  podcasts = [],
  episodes = [],
}: DiscoverySectionProps) {
  const { playEpisode } = useAudio()
  const [saved, setSaved] = useState<string[]>(() => {
    try {
      const stored = localStorage.getItem("angopod_saved")
      return stored ? JSON.parse(stored) : []
    } catch {
      return []
    }
  })

  const toggleSaved = (title: string) => {
    setSaved((current) => {
      const updated = current.includes(title)
        ? current.filter((item) => item !== title)
        : [...current, title]
      try {
        localStorage.setItem("angopod_saved", JSON.stringify(updated))
      } catch {}
      return updated
    })
  }

  const filteredPodcasts = podcasts.slice(0, 4).filter((podcast) => {
    const matchesCategory =
      activeCategory === "Todos" || podcast.category === activeCategory
    const text =
      `${podcast.title} ${podcast.host} ${podcast.category} ${podcast.province}`.toLowerCase()
    return matchesCategory && text.includes(query.toLowerCase())
  })

  const handlePlayPodcast = (podcastSlug: string) => {
    const ep = episodes.find((e) => e.podcastSlug === podcastSlug) || episodes[0]
    if (ep) playEpisode(ep)
  }

  return (
    <section className="discovery-section" id="descobrir">
      <div className="section-heading">
        <div>
          <span className="section-kicker">Feito para descobrir</span>
          <h2>Encontra a tua próxima voz favorita</h2>
        </div>
        <div className="category-tabs" role="list" aria-label="Filtrar por tema">
          {categories.map((category) => (
            <button
              className={activeCategory === category ? "active" : ""}
              key={category}
              onClick={() => setActiveCategory(category)}
              type="button"
            >
              {category}
            </button>
          ))}
        </div>
        <Link className="section-all-link" to="/explorar">
          Ver todos <Icon name="arrow" size={17} />
        </Link>
      </div>

      <div className="podcast-grid">
        {filteredPodcasts.map((podcast) => (
          <article className="podcast-card" key={podcast.title}>
            <div className={`cover ${podcast.color}`}>
              <img alt="" src={podcast.image} referrerPolicy="no-referrer" loading="lazy" />
              <span className="rank">#{podcast.rank}</span>
              <button
                className={saved.includes(podcast.title) ? "save saved" : "save"}
                type="button"
                aria-label={`${
                  saved.includes(podcast.title) ? "Remover" : "Guardar"
                } ${podcast.title}`}
                onClick={() => toggleSaved(podcast.title)}
              >
                <Icon name="bookmark" size={18} />
              </button>
              <button
                className="card-play"
                type="button"
                aria-label={`Ouvir ${podcast.title}`}
                onClick={() => handlePlayPodcast(podcast.slug)}
              >
                <Icon name="play" size={22} />
              </button>
            </div>
            <div className="card-meta">
              <span>{podcast.category}</span>
              <span>{podcast.province}</span>
            </div>
            <h3>
              <Link to={`/podcast/${podcast.slug}`}>{podcast.title}</Link>
            </h3>
            <p>com {podcast.host}</p>
          </article>
        ))}
        {filteredPodcasts.length === 0 && (
          <div className="empty-state">
            <Icon name="search" size={30} />
            <h3>Nenhum podcast encontrado</h3>
            <p>Tenta pesquisar por outro tema, criador ou província.</p>
            <button
              type="button"
              onClick={() => {
                setQuery("")
                setActiveCategory("Todos")
              }}
            >
              Limpar pesquisa
            </button>
          </div>
        )}
      </div>
    </section>
  )
}
