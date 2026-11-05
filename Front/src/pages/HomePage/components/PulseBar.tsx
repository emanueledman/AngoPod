import { Link } from "react-router-dom"
import { PulseHighlight } from "../../../types"

interface PulseBarProps {
  highlights?: PulseHighlight[]
}

export default function PulseBar({ highlights = [] }: PulseBarProps) {
  return (
    <section className="home-pulse-bar" aria-label="Destaques da semana">
      <div className="home-pulse-track">
        {[0, 1].map((copy) => (
          <div aria-hidden={copy === 1} className="home-pulse-group" key={copy}>
            <div className="home-pulse-label">
              <i /> Agora em Angola
            </div>
            {highlights.map((highlight) => (
              <Link to={highlight.href} key={`${copy}-${highlight.label}`}>
                <span>{highlight.label}</span>
                <strong>{highlight.title}</strong>
              </Link>
            ))}
          </div>
        ))}
      </div>
    </section>
  )
}
