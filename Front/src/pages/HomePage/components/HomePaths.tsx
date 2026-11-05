import { Link } from "react-router-dom"

export default function HomePaths() {
  return (
    <section className="home-paths">
      <div className="home-paths-heading">
        <span className="section-kicker">Começa por onde quiseres</span>
        <h2>
          Há mais de uma forma
          <br />
          de ouvir Angola.
        </h2>
      </div>
      <div className="home-path-cards">
        <Link to="/explorar">
          <span>01</span>
          <small>Mapa sonoro</small>
          <h3>
            Descobre por
            <br />
            província
          </h3>
          <p>Encontra vozes perto de ti num mapa vivo do país.</p>
          <i>→</i>
        </Link>
        <Link to="/ranking">
          <span>02</span>
          <small>O pulso da semana</small>
          <h3>
            Segue o<br />
            ranking nacional
          </h3>
          <p>Vê quem está a subir, quem chegou e o que todos partilham.</p>
          <i>→</i>
        </Link>
        <Link to="/apresentador/flay-squad">
          <span>03</span>
          <small>Por trás da voz</small>
          <h3>
            Conhece os
            <br />
            apresentadores
          </h3>
          <p>Histórias, percursos e pessoas que dão vida aos microfones.</p>
          <i>→</i>
        </Link>
      </div>
    </section>
  )
}

