import { useAudio } from "../../context/AudioContext"
import Icon from "../common/Icon"

function formatTime(seconds: number): string {
  const mins = Math.floor(seconds / 60)
  const secs = Math.floor(seconds % 60)
  return `${mins.toString().padStart(2, "0")}:${secs.toString().padStart(2, "0")}`
}

export default function GlobalAudioPlayer() {
  const { currentEpisode, isPlaying, togglePlay, currentTime, duration } = useAudio()

  if (!currentEpisode) return null

  const progressPercent = duration > 0 ? (currentTime / duration) * 100 : 35

  return (
    <div className={isPlaying ? "mini-player visible" : "mini-player"}>
      <img alt="" src={currentEpisode.image} />
      <div>
        <span>A reproduzir</span>
        <strong>{currentEpisode.title}</strong>
      </div>
      <button
        type="button"
        aria-label={isPlaying ? "Pausar" : "Reproduzir"}
        onClick={togglePlay}
      >
        <Icon name={isPlaying ? "pause" : "play"} size={20} />
      </button>
      <div className="progress">
        <span style={{ width: `${progressPercent}%` }} />
      </div>
      <span className="player-time">
        {formatTime(currentTime)} / {formatTime(duration)}
      </span>
      <Icon name="volume" size={19} />
    </div>
  )
}

