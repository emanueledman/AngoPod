import { Link } from "react-router-dom"

export default function SiteFooter({ dark = false }: { dark?: boolean }) {
  return (
    <footer className={dark ? "site-footer site-footer-dark" : "site-footer"}>
      <div className="site-footer-top">
        <Link className="brand footer-brand" to="/">
          <span className="brand-mark">
            <span />
            <span />
            <span />
          </span>
          <span>
            pod<span>angola</span>
          </span>
        </Link>
        <p>A maior casa das vozes angolanas. Descobre histórias, ideias e pessoas de todas as províncias.</p>
        <div className="site-footer-nav">
          <div>
            <strong>Descobrir</strong>
            <Link to="/explorar">Explorar podcasts</Link>
            <Link to="/ranking">Ranking nacional</Link>
            <a href="/#episodios">Novos episódios</a>
          </div>
          <div>
            <strong>Comunidade</strong>
            <Link to="/apresentador/flay-squad">Apresentadores</Link>
            <Link to="/podcast/fly-podcast">Podcast em destaque</Link>
            <Link to="/">Para criadores</Link>
          </div>
        </div>
      </div>
      <div className="site-footer-bottom">
        <span>© 2025 PodAngola</span>
        <span>Feito para amplificar Angola.</span>
      </div>
    </footer>
  )
}
