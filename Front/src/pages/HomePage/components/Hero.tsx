import React from "react"
import { Link } from "react-router-dom"
import Icon from "../../../components/common/Icon"
import { useAudio } from "../../../context/AudioContext"
import { Episode, Podcast } from "../../../types"

interface HeroProps {
  query: string
  setQuery: (q: string) => void
  onSearch: (e: React.FormEvent<HTMLFormElement>) => void
  featuredPodcast?: Podcast
  featuredEpisode?: Episode
}

export default function Hero({
  query,
  setQuery,
  onSearch,
  featuredPodcast,
  featuredEpisode,
}: HeroProps) {
  const { isPlaying, togglePlay, playEpisode, currentEpisode } = useAudio()

  if (!featuredPodcast) return null

  const isCurrentPlaying = isPlaying && currentEpisode?.podcastSlug === featuredPodcast.slug

  const handleToggleFeatured = () => {
    if (!featuredEpisode) return
    if (currentEpisode?.id === featuredEpisode.id) {
      togglePlay()
    } else {
      playEpisode(featuredEpisode)
    }
  }

  return (
    <section className="hero">
      <div className="hero-copy">
        <div className="eyebrow">
          <span /> TUDO SOBRE PODCASTS EM ANGOLA
        </div>
        <h1>
          Pod<em>Angola</em>
        </h1>
        <p>
          A plataforma que junta as conversas, ideias e histórias que estão
          a transformar Angola — de Cabinda ao Cunene.
        </p>
        <form className="search-box" onSubmit={onSearch}>
          <Icon name="search" size={22} />
          <input
            aria-label="Pesquisar podcasts"
            onChange={(event) => setQuery(event.target.value)}
            placeholder="Pesquisar podcasts, temas ou criadores"
            value={query}
          />
          <button type="submit">Pesquisar</button>
        </form>

        <div className="hero-proof" aria-label="PodAngola em números">
          <div>
            <strong>233+</strong>
            <span>podcasts</span>
          </div>
          <div>
            <strong>18</strong>
            <span>províncias</span>
          </div>
          <div>
            <strong>6</strong>
            <span>idiomas</span>
          </div>
          <Link to="/explorar">
            Explorar tudo <Icon name="arrow" size={16} />
          </Link>
        </div>
      </div>

      <div className="hero-visual" aria-label="Podcast em destaque">
        <div className="sun-shape" />
        <div className="feature-image">
          <img alt={`Apresentador do ${featuredPodcast.title}`} src={featuredPodcast.image} referrerPolicy="no-referrer" />
          <div className="image-label">Em destaque</div>
          <div className="feature-caption">
            <div>
              <span>Novo episódio</span>
              <strong>
                <Link to={`/podcast/${featuredPodcast.slug}`}>{featuredPodcast.title}</Link>
              </strong>
            </div>
            <button
              type="button"
              aria-label={`Reproduzir ${featuredPodcast.title}`}
              onClick={handleToggleFeatured}
            >
              <Icon name={isCurrentPlaying ? "pause" : "play"} size={23} />
            </button>
          </div>
        </div>
        <div className="listener-pill">
          <span className="listener-avatars">
            <i />
            <i />
            <i />
          </span>
          <span>
            <strong>+12 mil</strong> ouvintes esta semana
          </span>
        </div>
        <Link className="hero-rank-card" to="/ranking">
          <span>Ranking semanal</span>
          <strong>
            <i>#1</i> {featuredPodcast.title}
          </strong>
          <small>
            Subiu 2 posições <Icon name="arrow" size={14} />
          </small>
        </Link>
      </div>
    </section>
  )
}
