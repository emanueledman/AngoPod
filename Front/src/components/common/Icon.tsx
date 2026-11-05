import React from "react"
import { IconName } from "../../types"

interface IconProps {
  name: IconName
  size?: number
  className?: string
}

export default function Icon({ name, size = 20, className }: IconProps) {
  const paths: Record<IconName, React.ReactNode> = {
    arrow: (
      <>
        <path d="M5 12h14M13 6l6 6-6 6" />
      </>
    ),
    bookmark: (
      <path d="M6 4.8A1.8 1.8 0 0 1 7.8 3h8.4A1.8 1.8 0 0 1 18 4.8V21l-6-3.7L6 21V4.8Z" />
    ),
    chevron: <path d="m9 18 6-6-6-6" />,
    clock: (
      <>
        <circle cx="12" cy="12" r="9" />
        <path d="M12 7v5l3 2" />
      </>
    ),
    headphones: (
      <>
        <path d="M4 14v-2a8 8 0 0 1 16 0v2" />
        <path d="M4 14h3v6H5.5A1.5 1.5 0 0 1 4 18.5V14ZM20 14h-3v6h1.5a1.5 1.5 0 0 0 1.5-1.5V14Z" />
      </>
    ),
    menu: (
      <>
        <path d="M4 7h16M4 12h16M4 17h16" />
      </>
    ),
    pause: (
      <>
        <path d="M9 7v10M15 7v10" />
      </>
    ),
    play: <path d="m9 7 8 5-8 5V7Z" />,
    search: (
      <>
        <circle cx="11" cy="11" r="7" />
        <path d="m16 16 4 4" />
      </>
    ),
    sparkle: (
      <>
        <path d="m12 3 1.3 4.1a5.5 5.5 0 0 0 3.6 3.6L21 12l-4.1 1.3a5.5 5.5 0 0 0-3.6 3.6L12 21l-1.3-4.1a5.5 5.5 0 0 0-3.6-3.6L3 12l4.1-1.3a5.5 5.5 0 0 0 3.6-3.6L12 3Z" />
      </>
    ),
    trending: (
      <>
        <path d="m4 17 6-6 4 4 6-7" />
        <path d="M15 8h5v5" />
      </>
    ),
    volume: (
      <>
        <path d="M5 10v4h3l4 3V7l-4 3H5Z" />
        <path d="M16 9.5a4 4 0 0 1 0 5" />
      </>
    ),
  }

  return (
    <svg
      aria-hidden="true"
      className={className}
      fill="none"
      height={size}
      viewBox="0 0 24 24"
      width={size}
      stroke="currentColor"
      strokeLinecap="round"
      strokeLinejoin="round"
      strokeWidth="1.8"
    >
      {paths[name]}
    </svg>
  )
}

