import React, { useEffect, useMemo, useState } from "react"
import { Link } from "react-router-dom"
import Navbar from "../../components/layout/Navbar"
import { getPodcasts, getProvinces } from "../../services/podcastService"
import SiteFooter from "../../SiteFooter"
import { Podcast, Province } from "../../types"

export default function MapPage() {
  const [provinceList, setProvinceList] = useState<Province[]>([])
  const [podcastList, setPodcastList] = useState<Podcast[]>([])
  const [loading, setLoading] = useState(true)
  const [province, setProvince] = useState("Todas")
  const [category, setCategory] = useState("Todas")
  const [language, setLanguage] = useState("Todos")
  const [platform, setPlatform] = useState("Todas")
  const [query, setQuery] = useState("")
  const [rotationPaused, setRotationPaused] = useState(false)

  useEffect(() => {
    let isMounted = true

    async function loadData() {
      try {
        const [pods, provs] = await Promise.all([getPodcasts(), getProvinces()])
        if (isMounted) {
          setPodcastList(pods)
          setProvinceList(provs)
          setLoading(false)
        }
      } catch (err) {
        console.error("Erro ao carregar mapa do back-end:", err)
        if (isMounted) setLoading(false)
      }
    }

    loadData()

    return () => {
      isMounted = false
    }
  }, [])

  const selected = provinceList.find((item) => item.name === province)
  const results = useMemo(
    () =>
      podcastList.filter(
        (podcast) =>
          (province === "Todas" || podcast.province === province) &&
          (category === "Todas" || podcast.category === category) &&
          (language === "Todos" || podcast.language === language) &&
          (platform === "Todas" || podcast.platform === platform) &&
          `${podcast.title} ${podcast.host} ${podcast.category}`
            .toLowerCase()
            .includes(query.toLowerCase())
      ),
    [category, language, platform, podcastList, province, query]
  )

  return (
    <div className="map-page">
      <Navbar variant="map" />

      <main>
        <section className="map-hero">
          <div className="map-hero-copy">
            <span className="section-kicker">Geografia sonora</span>
            <h1>
              Ouve Angola<br />
              <em>província a província.</em>
            </h1>
            <p>
              Gira o mapa, escolhe um ponto e encontra podcasts feitos perto de ti
              — por tema, língua e comunidade.
            </p>
            <div className="map-search">
              <span>⌕</span>
              <input
                aria-label="Pesquisar podcasts"
                placeholder="Podcast, apresentador ou tema..."
                value={query}
                onChange={(event) => setQuery(event.target.value)}
              />
              <small>⌘ K</small>
            </div>
            <div className="map-filters">
              <label>
                <span>Província</span>
                <select
                  value={province}
                  onChange={(event) => setProvince(event.target.value)}
                >
                  <option>Todas</option>
                  {provinceList.map((item) => (
                    <option key={item.name}>{item.name}</option>
                  ))}
                </select>
              </label>
              <label>
                <span>Categoria</span>
                <select
                  value={category}
                  onChange={(event) => setCategory(event.target.value)}
                >
                  <option>Todas</option>
                  <option>Sociedade</option>
                  <option>Negócios</option>
                  <option>Cultura</option>
                  <option>Tecnologia</option>
                  <option>Notícias</option>
                  <option>Música</option>
                </select>
              </label>
              <label>
                <span>Idioma</span>
                <select
                  value={language}
                  onChange={(event) => setLanguage(event.target.value)}
                >
                  <option>Todos</option>
                  <option>Português</option>
                  <option>Kimbundu</option>
                  <option>Umbundu</option>
                </select>
              </label>
              <label>
                <span>Plataforma</span>
                <select
                  value={platform}
                  onChange={(event) => setPlatform(event.target.value)}
                >
                  <option>Todas</option>
                  <option>Spotify</option>
                  <option>YouTube</option>
                  <option>Apple Podcasts</option>
                </select>
              </label>
            </div>
          </div>

          <div className="map-stage">
            <div
              className={rotationPaused ? "map-3d paused" : "map-3d"}
              onMouseEnter={() => setRotationPaused(true)}
              onMouseLeave={() => setRotationPaused(false)}
            >
              {[5, 4, 3, 2, 1, 0].map((layer) => (
                <svg
                  className="map-layer"
                  key={layer}
                  style={{ "--layer": layer } as React.CSSProperties}
                  viewBox="0 0 360 420"
                  aria-hidden={layer !== 0}
                >
                  <path d="M105 18 241 15l34 71-19 72 45 66-13 128-89 38-64-25-54 12-37-79 25-79-23-55 45-53Z" />
                </svg>
              ))}
              <svg className="map-grid-lines" viewBox="0 0 360 420" aria-hidden="true">
                <path d="M100 120l156 38M76 219l218 5M81 298l207 54M142 365l15-203M206 390l-1-232M105 18l37 102-66 99M241 15l15 143 45 66" />
              </svg>
              {provinceList.map((item) => (
                <button
                  className={
                    province === item.name ? "map-3d-marker active" : "map-3d-marker"
                  }
                  style={{ left: `${item.x}%`, top: `${item.y}%` }}
                  key={item.name}
                  onClick={() => setProvince(item.name)}
                  type="button"
                >
                  <i />
                  <span>
                    {item.name}
                    <small>{item.count} podcasts</small>
                  </span>
                </button>
              ))}
            </div>
            <div className="map-orbit">
              <span>Arrasta com o olhar</span>
              <i />
            </div>
            <div className="map-stat">
              <strong>{selected?.count ?? 233}</strong>
              <span>podcasts<br />mapeados</span>
            </div>
          </div>
        </section>

        <section className="explore-pulse">
          <div className="pulse-intro">
            <span className="section-kicker">Agora em Angola</span>
            <h2>
              Um país inteiro<br />
              a conversar.
            </h2>
          </div>
          <div className="pulse-stat">
            <strong>233</strong>
            <span>podcasts<br />mapeados</span>
          </div>
          <div className="pulse-stat">
            <strong>18</strong>
            <span>províncias<br />representadas</span>
          </div>
          <div className="pulse-topics">
            <small>Temas a crescer</small>
            <div>
              <button type="button" onClick={() => setCategory("Negócios")}>
                # Empreendedorismo
              </button>
              <button type="button" onClick={() => setCategory("Cultura")}>
                # Cultura
              </button>
              <button type="button" onClick={() => setCategory("Tecnologia")}>
                # Tecnologia
              </button>
            </div>
          </div>
        </section>

        <section className="province-spotlight">
          <div className="spotlight-visual">
            <img
              src="https://images.unsplash.com/photo-1743172140481-5c7044dd23d6?crop=entropy&cs=tinysrgb&fit=crop&fm=jpg&q=85&w=1200"
              alt="Vista urbana de Luanda"
            />
            <span>Província da semana</span>
          </div>
          <div className="spotlight-copy">
            <span className="section-kicker">Destaque local · Luanda</span>
            <h2>Onde as conversas começam — e ganham o país.</h2>
            <p>
              124 podcasts, dezenas de formatos e uma nova geração de vozes. Descobre
              o que Luanda está a dizer esta semana.
            </p>
            <button type="button" onClick={() => setProvince("Luanda")}>
              Explorar Luanda <span>→</span>
            </button>
          </div>
          <div className="spotlight-rank">
            <small>Mais ouvido em Luanda</small>
            <strong>01</strong>
            <Link to="/podcast/fly-podcast">Fly Podcast →</Link>
          </div>
        </section>

        <section className="province-results">
          <div className="province-results-head">
            <div>
              <span className="section-kicker">
                {selected ? `Em ${selected.name}` : "De norte a sul"}
              </span>
              <h2>
                {selected ? selected.description : "Podcasts encontrados no mapa"}
              </h2>
            </div>
            <p>
              <strong>{results.length}</strong> resultados
            </p>
          </div>
          <div className="province-grid">
            {results.map((podcast, index) => (
              <article className="province-card" key={podcast.title}>
                <Link className="province-card-cover" to={`/podcast/${podcast.slug}`}>
                  <img src={podcast.image} alt="" referrerPolicy="no-referrer" loading="lazy" />
                  <span className="province-card-number">0{index + 1}</span>
                  <i>▶</i>
                </Link>
                <span>
                  {podcast.province} · {podcast.category}
                </span>
                <Link to={`/podcast/${podcast.slug}`}>
                  <h3>{podcast.title}</h3>
                </Link>
                <p>
                  com{" "}
                  {podcast.slug === "fly-podcast" ? (
                    <Link to="/apresentador/flay-squad">{podcast.host}</Link>
                  ) : (
                    podcast.host
                  )}
                </p>
                <div className="province-card-meta">
                  <small>
                    {podcast.language} · {podcast.platform}
                  </small>
                  <small>{podcast.episodesCount} episódios</small>
                </div>
                <Link className="province-card-action" to={`/podcast/${podcast.slug}`}>
                  Ver podcast <span>→</span>
                </Link>
              </article>
            ))}
          </div>
          {results.length === 0 && (
            <div className="map-empty">
              <strong>Ainda não encontrámos uma voz com estes filtros.</strong>
              <button
                type="button"
                onClick={() => {
                  setProvince("Todas")
                  setCategory("Todas")
                  setLanguage("Todos")
                  setPlatform("Todas")
                  setQuery("")
                }}
              >
                Mostrar todos
              </button>
            </div>
          )}
        </section>
      </main>

      <SiteFooter dark />
    </div>
  )
}
