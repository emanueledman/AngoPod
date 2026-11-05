import { useState } from "react"
import { Link, useLocation } from "react-router-dom"
import Icon from "../common/Icon"

interface NavbarProps {
  variant?: "home" | "detail" | "map"
  backLink?: { to: string; label: string }
}

export default function Navbar({ variant = "home", backLink }: NavbarProps) {
  const [mobileOpen, setMobileOpen] = useState(false)
  const location = useLocation()

  if (variant === "detail") {
    return (
      <header className="detail-header">
        <Link className="brand" to="/">
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
          <Link to="/" className={location.pathname === "/" ? "active" : ""}>
            Início
          </Link>
          <Link to="/explorar" className={location.pathname === "/explorar" ? "active" : ""}>
            Explorar
          </Link>
          <Link to="/ranking" className={location.pathname === "/ranking" ? "active" : ""}>
            Ranking
          </Link>
        </nav>
        {backLink && (
          <Link className="back-link" to={backLink.to}>
            {backLink.label}
          </Link>
        )}
      </header>
    )
  }

  if (variant === "map") {
    return (
      <header className="explore-header map-header">
        <Link className="brand" to="/">
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
          <Link className="active" to="/explorar">
            Explorar
          </Link>
          <Link to="/ranking">Ranking</Link>
          <a href="/#episodios">Episódios</a>
        </nav>
        <span className="map-header-note">18 províncias · uma só voz</span>
      </header>
    )
  }

  return (
    <header className="topbar">
      <Link className="brand" to="/" aria-label="PodAngola, página inicial">
        <span className="brand-mark">
          <span />
          <span />
          <span />
        </span>
        <span>
          pod<span>angola</span>
        </span>
      </Link>
      <nav
        className={mobileOpen ? "nav-links open" : "nav-links"}
        aria-label="Navegação principal"
      >
        <Link to="/explorar" className={location.pathname === "/explorar" ? "active" : ""}>
          Explorar
        </Link>
        <Link to="/ranking" className={location.pathname === "/ranking" ? "active" : ""}>
          Ranking
        </Link>
        <a href="/#episodios">Episódios</a>
      </nav>
      <div className="header-actions">
        <button className="text-button" type="button">
          Entrar
        </button>
        <button className="creator-button" type="button">
          Para criadores <Icon name="arrow" size={17} />
        </button>
        <button
          className="menu-button"
          type="button"
          aria-label="Abrir menu"
          aria-expanded={mobileOpen}
          onClick={() => setMobileOpen(!mobileOpen)}
        >
          <Icon name="menu" />
        </button>
      </div>
    </header>
  )
}

