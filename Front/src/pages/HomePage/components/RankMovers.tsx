import { Link } from "react-router-dom"
import { MoverItem } from "../../../types"

interface RankMoversProps {
  movers?: MoverItem[]
}

export default function RankMovers({ movers = [] }: RankMoversProps) {
  return (
    <section className="rank-movers home-rank-movers">
      <div className="movers-heading">
        <span className="section-kicker">Movimentos da semana</span>
        <h2>
          Quem está a<br />
          mudar o jogo.
        </h2>
        <p>
          Entradas, subidas e tendências que merecem atenção antes de
          chegarem ao pódio.
        </p>
      </div>
      <div className="mover-cards">
        {movers.map((mover) => (
          <article className={`mover-card mover-${mover.type}`} key={mover.title}>
            <span>{mover.badgeLabel}</span>
            <strong>{mover.value}</strong>
            <div>
              <img src={mover.image} alt={mover.title} />
              <p>
                <b>{mover.title}</b>
                <small>{mover.subtitle}</small>
              </p>
            </div>
            <Link to={mover.linkHref}>{mover.linkLabel}</Link>
          </article>
        ))}
      </div>
    </section>
  )
}
