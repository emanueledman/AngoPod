import { useEffect, useState } from "react"
import { Link, useParams } from "react-router-dom"
import Navbar from "../../components/layout/Navbar"
import { getPresenterBySlug } from "../../services/podcastService"
import SiteFooter from "../../SiteFooter"
import { Presenter } from "../../types"

export default function PresenterDetailPage() {
  const { slug } = useParams<{ slug: string }>()
  const [followed, setFollowed] = useState(false)
  const [presenter, setPresenter] = useState<Presenter | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    let isMounted = true
    const currentSlug = slug || "flay-squad"

    async function loadPresenter() {
      setLoading(true)
      try {
        const pres = await getPresenterBySlug(currentSlug)
        if (isMounted) {
          if (pres) setPresenter(pres)
          setLoading(false)
        }
      } catch (err) {
        console.error("Erro ao carregar apresentador da API:", err)
        if (isMounted) setLoading(false)
      }
    }

    loadPresenter()

    return () => {
      isMounted = false
    }
  }, [slug])

  if (loading && !presenter) {
    return (
      <div className="presenter-page">
        <Navbar variant="detail" backLink={{ to: "/ranking", label: "← Voltar ao ranking" }} />
        <main style={{ padding: "4rem 2rem", textAlign: "center" }}>
          <h2>A carregar perfil...</h2>
        </main>
        <SiteFooter dark />
      </div>
    )
  }

  if (!presenter) {
    return (
      <div className="presenter-page">
        <Navbar variant="detail" backLink={{ to: "/ranking", label: "← Voltar ao ranking" }} />
        <main style={{ padding: "4rem 2rem", textAlign: "center" }}>
          <h2>Apresentador não encontrado</h2>
          <Link to="/ranking">Voltar ao ranking</Link>
        </main>
        <SiteFooter dark />
      </div>
    )
  }

  return (
    <div className="presenter-page">
      <Navbar
        variant="detail"
        backLink={{ to: "/ranking", label: "← Voltar ao ranking" }}
      />
      <main>
        <section className="presenter-hero">
          <div className="presenter-photo">
            <img src={presenter.photo} alt={`Retrato de ${presenter.name}`} />
            <span>
              VOZ<br />
              #{presenter.badgeNumber}
            </span>
          </div>
          <div className="presenter-copy">
            <span className="section-kicker">{presenter.role}</span>
            <h1>
              {presenter.name.split(" ")[0]}
              <br />
              <em>{presenter.name.split(" ").slice(1).join(" ")}.</em>
            </h1>
            <blockquote>{presenter.quote}</blockquote>
            <p>{presenter.bio}</p>
            <div className="presenter-actions">
              <button
                className={followed ? "followed" : ""}
                type="button"
                onClick={() => setFollowed(!followed)}
              >
                {followed ? "A seguir" : "+ Seguir apresentador"}
              </button>
              <Link to={`/podcast/${presenter.podcastSlug}`}>
                Ver {presenter.podcastTitle} →
              </Link>
            </div>
            <div className="presenter-stats">
              <div>
                <strong>{presenter.followers}</strong>
                <span>seguidores</span>
              </div>
              <div>
                <strong>{presenter.listeners}</strong>
                <span>ouvintes</span>
              </div>
              <div>
                <strong>{presenter.conversations}</strong>
                <span>conversas</span>
              </div>
            </div>
          </div>
        </section>
        <section className="presenter-story">
          <div>
            <span className="section-kicker">A pessoa por trás da voz</span>
            <h2>
              Curiosidade sem filtro.
              <br />
              Angola sem tradução.
            </h2>
          </div>
          <div>
            <p>{presenter.story1}</p>
            <p>{presenter.story2}</p>
          </div>
        </section>
      </main>
      <SiteFooter dark />
    </div>
  )
}
