import React, { createContext, useContext, useEffect, useRef, useState } from "react"
import { Episode } from "../types"

interface AudioContextType {
  currentEpisode: Episode | null
  isPlaying: boolean
  currentTime: number
  duration: number
  playEpisode: (episode: Episode) => void
  togglePlay: () => void
  pause: () => void
  seek: (time: number) => void
}

const AudioContext = createContext<AudioContextType | undefined>(undefined)

export function AudioProvider({ children }: { children: React.ReactNode }) {
  const [currentEpisode, setCurrentEpisode] = useState<Episode | null>(null)
  const [isPlaying, setIsPlaying] = useState<boolean>(false)
  const [currentTime, setCurrentTime] = useState<number>(1122) // 18:42 default
  const [duration, setDuration] = useState<number>(2890) // 48:10 default
  const audioRef = useRef<HTMLAudioElement | null>(null)

  useEffect(() => {
    const audio = new Audio()
    audioRef.current = audio

    const updateTime = () => setCurrentTime(audio.currentTime)
    const updateDuration = () => setDuration(audio.duration || 2890)
    const onEnded = () => setIsPlaying(false)

    audio.addEventListener("timeupdate", updateTime)
    audio.addEventListener("loadedmetadata", updateDuration)
    audio.addEventListener("ended", onEnded)

    return () => {
      audio.removeEventListener("timeupdate", updateTime)
      audio.removeEventListener("loadedmetadata", updateDuration)
      audio.removeEventListener("ended", onEnded)
      audio.pause()
    }
  }, [])

  const playEpisode = (episode: Episode) => {
    setCurrentEpisode(episode)
    setIsPlaying(true)

    if (audioRef.current) {
      if (episode.audioUrl && audioRef.current.src !== episode.audioUrl) {
        audioRef.current.src = episode.audioUrl
        audioRef.current.load()
      }
      audioRef.current.play().catch(() => {
        // Ignora política de autoplay se bloqueado pelo browser
      })
    }
  }

  const togglePlay = () => {
    if (!currentEpisode) return

    if (isPlaying) {
      pause()
    } else {
      setIsPlaying(true)
      audioRef.current?.play().catch(() => {})
    }
  }

  const pause = () => {
    setIsPlaying(false)
    audioRef.current?.pause()
  }

  const seek = (time: number) => {
    setCurrentTime(time)
    if (audioRef.current) {
      audioRef.current.currentTime = time
    }
  }

  return (
    <AudioContext.Provider
      value={{
        currentEpisode,
        isPlaying,
        currentTime,
        duration,
        playEpisode,
        togglePlay,
        pause,
        seek,
      }}
    >
      {children}
    </AudioContext.Provider>
  )
}

export function useAudio() {
  const context = useContext(AudioContext)
  if (!context) {
    throw new Error("useAudio deve ser usado dentro de um AudioProvider")
  }
  return context
}

