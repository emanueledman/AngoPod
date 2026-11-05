import { useEffect, useState } from "react"
import { Link } from "react-router-dom"
import Icon from "../../components/common/Icon"
import { useAudio } from "../../context/AudioContext"
import { getEpisodes, getPodcasts } from "../../services/podcastService"
import SiteFooter from "../../SiteFooter"
import { Episode, Podcast } from "../../types"

const weeks = ["24–30 JUN", "17–23 JUN", "10–16 JUN", "03–09 JUN"]
const categories = ["Todas", "Sociedade", "Cultura", "Tecnologia", "Negócios", "Música"]

// Apresentadores em destaque — imagens reais adicionadas pelo back-end via /presenters
const hosts = [
  {
    name: "Flay Squad",
    slug: "flay-squad",
    role: "Jornalista & Rapper",
    followers: "42 mil",
    image: "https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=400&q=80",
    quote: "Pergunto o que toda a gente quer saber.",
  },
  {
    name: "Euclides Francisco",
    slug: "dinheiro-limpo",
    role: "Investigador cultural",
    followers: "31 mil",
    image: "https://images.unsplash.com/photo-1500648767791-00dcc994a43e?w=400&q=80",
    quote: "O futuro também se constrói com memória.",
  },
  {
    name: "Elisa Manuel",
    slug: "kota-me-conta",
    role: "Contadora de histórias",
    followers: "28 mil",
    image: "https://images.unsplash.com/photo-1573496359142-b8d87734a5a2?w=400&q=80",
    quote: "As histórias dos nossos kotas nunca morrem.",
  },
  {
    name: "Lina Santos",
    slug: "batidas-da-banda",
    role: "Curadora musical",
    followers: "24 mil",
    image: "https://images.unsplash.com/photo-1531746020798-e6953c6e8e04?w=400&q=80",
    quote: "A banda tem ritmo. Eu encontro-o.",
  },
]

export default function RankingPage() {
  const [period, setPeriod] = useState("Semanal")
  const { playEpisode, isPlaying, currentEpisode, togglePlay } = useAudio()
  const [following, setFollowing] = useState<string[]>([])
  const [rankingType, setRankingType] = useState<"podcasts" | "episodios">("episodios")
  const [showAll, setShowAll] = useState(false)
  const [weekIndex, setWeekIndex] = useState(0)
  const [category, setCategory] = useState("Todas")
  const [query, setQuery] = useState("")
  const [shared, setShared] = useState(false)

  const [podcastList, setPodcastList] = useState<Podcast[]>([])
  const [episodeList, setEpisodeList] = useState<Episode[]>([])

  useEffect(() => {
    let isMounted = true

    async function loadRankingData() {
      try {
        const [pods, eps] = await Promise.all([getPodcasts(), getEpisodes()])
        if (isMounted) {
          setPodcastList(pods)
          setEpisodeList(eps)
        }
      } catch (err) {
        console.error("Erro ao carregar ranking da API:", err)
      }
    }

    loadRankingData()

    return () => {
      isMounted = false
    }
  }, [])

  const follow = (name: string) =>
    setFollowing((current) =>
      current.includes(name)
        ? current.filter((item) => item !== name)
        : [...current, name]
    )

  const activeRanking =
    rankingType === "episodios"
      ? episodeList.map((item) => ({
          title: item.title,
          show: item.podcastTitle,
          category: item.tag || "Cultura",
          time: item.duration || item.time || "45 min",
          listens: item.listens || "25K",
          trend: item.trend || "+10%",
          image: item.image || (podcastList[0] ? podcastList[0].image : ""),
          episodeObj: item,
        }))
      : podcastList.map((pod) => ({
          title: pod.title,
          show: pod.host,
          category: pod.category,
          time: pod.province,
          listens: pod.listeners || "50K",
          trend: `+${pod.popularity % 10 + 1}%`,
          image: pod.image,
          episodeObj: episodeList.find((e) => e.podcastSlug === pod.slug) || episodeList[0] || null,
        }))

  const filteredRanking = activeRanking.filter(
    (item) =>
      (category === "Todas" || item.category === category) &&
      `${item.title} ${item.show}`.toLowerCase().includes(query.toLowerCase())
  )

  const visibleRanking = showAll ? filteredRanking : filteredRanking.slice(0, 5)

  const shareRanking = async () => {
    try {
      await navigator.clipboard?.writeText(window.location.href)
      setShared(true)
      window.setTimeout(() => setShared(false), 1800)
    } catch {}
  }

  const p1 = podcastList[1] ?? null
  const p0 = podcastList[0] ?? null
  const p2 = podcastList[2] ?? null

  const podiumShows = [p1, p0, p2]
    .map((p, i) => p ? { ...p, place: [2, 1, 3][i], color: ["violet", "orange", "lime"][i] } : null)
    .filter(Boolean) as Array<Podcast & { place: number; color: string }>

  return (
    <div className="rank-page">
      <header className="rank-header">
        <Link className="brand rank-brand" to="/">
          <span className="brand-mark">
            <span />
            <span />
            <span />
          </span>
          <span>
            pod<span>angola</span>
          </span>
        </Link>
        <nav>
          <Link to="/">Início</Link>
          <Link to="/explorar">Explorar</Link>
          <Link className="active" to="/ranking">
            Ranking
          </Link>
          <a href="/#episodios">Episódios</a>
        </nav>
        <span className="rank-live">
          <i /> Atualizado há 2 horas
        </span>
      </header>

      <main>
        <section className="rank-hero">
          <div>
            <span className="rank-overline">
              ↗ O termómetro das vozes angolanas
            </span>
            <h1>
              O que Angola<br />
              não para de <em>ouvir.</em>
            </h1>
            <p>
              Os podcasts, episódios e vozes que dominaram a semana. Um retrato
              vivo do que está a mexer com o país.
            </p>
            <div className="period-switch">
              {["Semanal", "Mensal", "Do ano"].map((item) => (
                <button
                  className={period === item ? "active" : ""}
                  onClick={() => setPeriod(item)}
                  type="button"
                  key={item}
                >
                  {item}
                </button>
              ))}
            </div>
            <div className="rank-week-nav">
              <button
                disabled={weekIndex === weeks.length - 1}
                onClick={() => setWeekIndex(Math.min(weeks.length - 1, weekIndex + 1))}
                type="button"
              >
                ← Semana anterior
              </button>
              <strong>{weeks[weekIndex]}</strong>
              <button
                disabled={weekIndex === 0}
                onClick={() => setWeekIndex(Math.max(0, weekIndex - 1))}
                type="button"
              >
                Semana seguinte →
              </button>
            </div>
          </div>
          <div className="rank-stamp">
            <span>TOP</span>
            <strong>10</strong>
            <small>
              {period}
              <br />
              {weeks[weekIndex]}
            </small>
          </div>
        </section>

        <section className="rank-dashboard">
          <div>
            <small>Reproduções analisadas</small>
            <strong>1.8M</strong>
            <span>+14% nesta semana</span>
          </div>
          <div>
            <small>Podcasts monitorizados</small>
            <strong>233</strong>
            <span>18 províncias</span>
          </div>
          <div>
            <small>Nova entrada mais alta</small>
            <strong>#04</strong>
            <span>Sem Filtro AO</span>
          </div>
          <div>
            <small>Maior subida</small>
            <strong>+7</strong>
            <span>Ideias do Sul</span>
          </div>
          <button type="button" onClick={shareRanking}>
            {shared ? "Link copiado" : "Partilhar ranking ↗"}
          </button>
        </section>

        <section className="podium-section">
          <div className="podium-heading">
            <span>01</span>
            <div>
              <small>Podcasts da semana</small>
              <h2>O pódio nacional</h2>
            </div>
            <p>Classificação baseada em reproduções, partilhas e crescimento.</p>
          </div>
          <div className="podium">
            {podiumShows.map((show) => {
              const ep = episodeList.find((e) => e.podcastSlug === show.slug) || episodeList[0] || null
              const isThisPlaying = isPlaying && ep && currentEpisode?.id === ep.id

              return (
                <article
                  className={`podium-card place-${show.place} ${show.color}`}
                  key={show.title}
                >
                  <span className="podium-number">{show.place}</span>
                  <div className="podium-photo">
                    <img src={show.image} alt={`Apresentador de ${show.title}`} referrerPolicy="no-referrer" />
                    <button
                      onClick={() => (isThisPlaying ? togglePlay() : ep && playEpisode(ep))}
                      type="button"
                      aria-label={`Ouvir ${show.title}`}
                    >
                      <Icon name={isThisPlaying ? "pause" : "play"} size={20} />
                    </button>
                  </div>
                  <span className="podium-category">{show.category}</span>
                  <h3>
                    <Link to={`/podcast/${show.slug}`}>{show.title}</Link>
                  </h3>
                  <p>
                    com{" "}
                    {show.slug === "fly-podcast" ? (
                      <Link to="/apresentador/flay-squad">{show.host}</Link>
                    ) : (
                      show.host
                    )}
                  </p>
                  <div className="podium-stats">
                    <span>◉ {show.listeners} ouvintes</span>
                    <strong>+{show.popularity % 10 + 2}%</strong>
                  </div>
                </article>
              )
            })}
          </div>
        </section>

        <section className="weekly-section">
          <div className="weekly-title">
            <span className="section-kicker">Ranking completo</span>
            <h2>
              {rankingType === "episodios" ? "Episódios que" : "Podcasts que"}
              <br />
              marcaram a semana
            </h2>
            <p>
              Alterna entre podcasts e episódios para descobrir todas as posições
              do Top 10.
            </p>
            <div className="ranking-type-switch">
              <button
                className={rankingType === "podcasts" ? "active" : ""}
                onClick={() => {
                  setRankingType("podcasts")
                  setShowAll(false)
                }}
                type="button"
              >
                Podcasts
              </button>
              <button
                className={rankingType === "episodios" ? "active" : ""}
                onClick={() => {
                  setRankingType("episodios")
                  setShowAll(false)
                }}
                type="button"
              >
                Episódios
              </button>
            </div>
            <div className="rank-search">
              <span>⌕</span>
              <input
                aria-label="Pesquisar no ranking"
                placeholder="Pesquisar no ranking"
                value={query}
                onChange={(event) => setQuery(event.target.value)}
              />
            </div>
            <div className="rank-category-filter">
              {categories.map((item) => (
                <button
                  className={category === item ? "active" : ""}
                  key={item}
                  onClick={() => {
                    setCategory(item)
                    setShowAll(false)
                  }}
                  type="button"
                >
                  {item}
                </button>
              ))}
            </div>
          </div>
          <div className="weekly-list">
            <div className="ranking-table-head">
              <span>Posição</span>
              <span>Conteúdo</span>
              <span>Audiência</span>
              <span>Movimento</span>
            </div>
            {visibleRanking.map((item, index) => {
              const isThisPlaying =
                isPlaying && currentEpisode?.id === item.episodeObj.id

              return (
                <article
                  className={isThisPlaying ? "is-playing" : ""}
                  key={item.title}
                >
                  <span className="weekly-rank">0{index + 1}</span>
                  <div className="weekly-image">
                    <img src={item.image} alt="" referrerPolicy="no-referrer" loading="lazy" />
                    <button
                      onClick={() =>
                        isThisPlaying ? togglePlay() : playEpisode(item.episodeObj)
                      }
                      type="button"
                      aria-label={`Reproduzir ${item.title}`}
                    >
                      <Icon name={isThisPlaying ? "pause" : "play"} size={18} />
                    </button>
                  </div>
                  <div className="weekly-info">
                    <span>{item.show}</span>
                    <h3>{item.title}</h3>
                    <small>{item.time}</small>
                  </div>
                  <div className="sound-bars">
                    <i />
                    <i />
                    <i />
                    <i />
                    <i />
                    <i />
                  </div>
                  <div className="weekly-listens">
                    <strong>{item.listens}</strong>
                    <span>reproduções</span>
                  </div>
                  <span className="weekly-trend">{item.trend}</span>
                </article>
              )
            })}
            {visibleRanking.length > 0 ? (
              <button
                className="show-ranking-button"
                type="button"
                onClick={() => setShowAll(!showAll)}
              >
                {showAll
                  ? "Mostrar apenas o Top 5"
                  : `Ver Top 10 de ${rankingType}`}{" "}
                <span>→</span>
              </button>
            ) : (
              <div className="ranking-empty">
                <strong>Nenhum resultado nesta categoria.</strong>
                <button
                  type="button"
                  onClick={() => {
                    setCategory("Todas")
                    setQuery("")
                  }}
                >
                  Limpar filtros
                </button>
              </div>
            )}
          </div>
        </section>

        <section className="hosts-section">
          <div className="hosts-heading">
            <div>
              <span className="section-kicker">Por trás do microfone</span>
              <h2>
                Vozes que<br />
                <em>prendem.</em>
              </h2>
            </div>
            <p>
              Não são apenas apresentadores. São as pessoas que dão contexto,
              provocam conversas e fazem o país carregar no play.
            </p>
          </div>
          <div className="host-grid">
            {hosts.map((host, index) => (
              <article className="host-card" key={host.name}>
                <div className="host-image">
                  <img src={host.image} alt={`Retrato de ${host.name}`} />
                  <span>0{index + 1}</span>
                </div>
                <div className="host-copy">
                  <span className="host-role">{host.role}</span>
                  <h3>
                    {host.slug ? (
                      <Link to={`/apresentador/${host.slug}`}>{host.name}</Link>
                    ) : (
                      host.name
                    )}
                  </h3>
                  <p>“{host.quote}”</p>
                  <div>
                    <span>
                      <strong>{host.followers}</strong> seguidores
                    </span>
                    <button
                      className={following.includes(host.name) ? "following" : ""}
                      onClick={() => follow(host.name)}
                      type="button"
                    >
                      {following.includes(host.name) ? "A seguir" : "Seguir"}{" "}
                      <b>+</b>
                    </button>
                  </div>
                </div>
              </article>
            ))}
          </div>
        </section>

        <section className="rank-cta">
          <div>
            <span>O teu podcast está a subir?</span>
            <h2>
              Partilha a tua posição<br />
              com o mundo.
            </h2>
          </div>
          <Link to="/explorar">Explorar todos os podcasts →</Link>
        </section>

        <section className="rank-methodology">
          <div>
            <span className="section-kicker">Ranking com contexto</span>
            <h2>
              Como calculamos<br />
              cada posição?
            </h2>
          </div>
          <div className="method-steps">
            <article>
              <span>01</span>
              <h3>Audiência</h3>
              <p>Reproduções válidas e crescimento de ouvintes únicos durante o período.</p>
            </article>
            <article>
              <span>02</span>
              <h3>Envolvimento</h3>
              <p>Partilhas, favoritos, avaliações e conclusão média dos episódios.</p>
            </article>
            <article>
              <span>03</span>
              <h3>Descoberta</h3>
              <p>Novos seguidores e alcance fora da audiência habitual do podcast.</p>
            </article>
          </div>
          <p className="method-note">
            Atualizado todas as segundas-feiras. Dados demonstrativos nesta versão do produto.
          </p>
        </section>

        <SiteFooter dark />
      </main>
    </div>
  )
}
