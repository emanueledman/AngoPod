import { BrowserRouter, Navigate, Route, Routes } from "react-router-dom"
import GlobalAudioPlayer from "./components/player/GlobalAudioPlayer"
import { AudioProvider } from "./context/AudioContext"
import HomePage from "./pages/HomePage"
import MapPage from "./pages/MapPage"
import PodcastDetailPage from "./pages/PodcastDetailPage"
import PresenterDetailPage from "./pages/PresenterDetailPage"
import RankingPage from "./pages/RankingPage"

export default function App() {
  return (
    <BrowserRouter>
      <AudioProvider>
        <Routes>
          <Route path="/" element={<HomePage />} />
          <Route path="/explorar" element={<MapPage />} />
          <Route path="/mapa" element={<MapPage />} />
          <Route path="/ranking" element={<RankingPage />} />
          <Route path="/podcast/:slug" element={<PodcastDetailPage />} />
          <Route path="/apresentador/:slug" element={<PresenterDetailPage />} />
          <Route path="*" element={<Navigate to="/" replace />} />
        </Routes>
        <GlobalAudioPlayer />
      </AudioProvider>
    </BrowserRouter>
  )
}
