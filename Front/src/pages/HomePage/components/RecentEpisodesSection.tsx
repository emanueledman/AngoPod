import Icon from "../../../components/common/Icon"
import { useAudio } from "../../../context/AudioContext"
import { Episode } from "../../../types"

interface RecentEpisodesProps {
  episodes?: Episode[]
}

export default function RecentEpisodesSection({
  episodes = [],
}: RecentEpisodesProps) {
  const { playEpisode, isPlaying, currentEpisode } = useAudio()

  return (
    <section className="episodes-section" id="episodios">
      <div className="episodes-intro">
        <span className="section-kicker">
          <Icon name="headphones" size={16} /> Acabados de chegar
        </span>
        <h2>
          Episódios
          <br />
          Recentes
        </h2>
        <p>Novas conversas, publicadas esta semana.</p>
        <a href="#descobrir">
          Ver todos os episódios <Icon name="arrow" size={18} />
        </a>
      </div>
      <div className="episode-list">
        {episodes.slice(0, 5).map((episode, index) => {
          const isThisPlaying = isPlaying && currentEpisode?.id === episode.id

          return (
            <article key={episode.id}>
              <button
                type="button"
                aria-label={`Reproduzir ${episode.title}`}
                onClick={() => playEpisode(episode)}
              >
                <Icon name={isThisPlaying ? "pause" : "play"} size={20} />
              </button>
              <span className="episode-index">0{index + 1}</span>
              <div>
                <span className="episode-tag">{episode.tag}</span>
                <h3>{episode.title}</h3>
                <p>{episode.show || episode.podcastTitle}</p>
              </div>
              <span className="episode-time">
                <Icon name="clock" size={15} /> {episode.time || episode.duration}
              </span>
              <Icon name="chevron" size={20} />
            </article>
          )
        })}
      </div>
    </section>
  )
}
