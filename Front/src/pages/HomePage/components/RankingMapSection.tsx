import { useEffect, useState } from "react"
import { Link } from "react-router-dom"
import Icon from "../../../components/common/Icon"
import { useAudio } from "../../../context/AudioContext"
import { Episode, Podcast, Province } from "../../../types"

interface RankingMapProps {
  podcasts?: Podcast[]
  provinces?: Province[]
  episodes?: Episode[]
}

export default function RankingMapSection({
  podcasts = [],
  provinces = [],
  episodes = [],
}: RankingMapProps) {
  const [selectedProvince, setSelectedProvince] = useState<Province | null>(null)
  const { playEpisode } = useAudio()

  useEffect(() => {
    if (provinces.length > 0) {
      setSelectedProvince(provinces[0])
    }
  }, [provinces])

  const handlePlayPodcast = (podcastSlug: string) => {
    const ep = episodes.find((e) => e.podcastSlug === podcastSlug) || episodes[0]
    if (ep) playEpisode(ep)
  }

  return (
    <section className="ranking-map">
      <div className="ranking-panel" id="ranking">
        <div className="panel-heading">
          <div>
            <span className="section-kicker light">
              <Icon name="trending" size={16} /> Ranking nacional
            </span>
            <h2>
              O que Angola
              <br />
              está a ouvir
            </h2>
          </div>
          <span className="week-badge">Semana 24</span>
        </div>
        <ol className="ranking-list">
          {podcasts.slice(0, 3).map((podcast, index) => (
            <li key={podcast.title}>
              <span className="ranking-number">0{index + 1}</span>
              <img alt="" src={podcast.image} referrerPolicy="no-referrer" loading="lazy" />
              <div>
                <strong>
                  <Link to={`/podcast/${podcast.slug}`}>{podcast.title}</Link>
                </strong>
                <span>
                  {podcast.category} · {podcast.host}
                </span>
              </div>
              <span className="change">+{12 - index * 3}%</span>
              <button
                type="button"
                aria-label={`Ouvir ${podcast.title}`}
                onClick={() => handlePlayPodcast(podcast.slug)}
              >
                <Icon name="play" size={18} />
              </button>
            </li>
          ))}
        </ol>
        <Link className="panel-link" to="/ranking">
          Ver ranking completo <Icon name="arrow" size={18} />
        </Link>
      </div>

      <div className="map-panel" id="mapa">
        <div className="map-copy">
          <span className="section-kicker">
            <span className="pin-dot" /> Vozes de todo o país
          </span>
          <h2>
            Explora por
            <br />
            província
          </h2>
          <p>
            Cada lugar tem uma história diferente. Encontra quem está a
            criar perto de ti.
          </p>
          <div className="province-result">
            <span>{selectedProvince?.count ?? 0}</span>
            <div>
              <strong>podcasts</strong>
              <small>em {selectedProvince?.name ?? "Angola"}</small>
            </div>
          </div>
          <Link className="outline-button" to="/explorar">
            Explorar {selectedProvince?.name ?? "Angola"} <Icon name="arrow" size={18} />
          </Link>
        </div>
        <div className="angola-map" aria-label="Mapa interativo de Angola">
          <svg viewBox="0 0 360 420" role="img" aria-label="Silhueta estilizada de Angola">
            <path
              className="map-shadow"
              d="M112 25 248 22l34 71-19 72 45 66-13 128-89 38-64-25-54 12-37-79 25-79-23-55 45-53Z"
            />
            <path
              className="map-shape"
              d="M105 18 241 15l34 71-19 72 45 66-13 128-89 38-64-25-54 12-37-79 25-79-23-55 45-53Z"
            />
            <path
              className="map-lines"
              d="M100 120l156 38M76 219l218 5M81 298l207 54M142 365l15-203M206 390l-1-232M105 18l37 102-66 99M241 15l15 143 45 66"
            />
          </svg>
          {provinces.map((province) => (
            <button
              className={
                selectedProvince?.name === province.name
                  ? "map-marker selected"
                  : "map-marker"
              }
              key={province.name}
              onClick={() => setSelectedProvince(province)}
              style={{ left: `${province.x}%`, top: `${province.y}%` }}
              type="button"
              aria-label={`Selecionar ${province.name}, ${province.count} podcasts`}
            >
              <span />
              <label>{province.name}</label>
            </button>
          ))}
        </div>
      </div>
    </section>
  )
}
