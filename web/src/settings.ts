// User-facing display settings: theme, font, density, size. Applied by setting
// data-* attributes on <html>; the CSS in index.css reacts to them. Persisted
// in localStorage. No web fonts are loaded (system stacks only) — both for
// speed and to avoid sending user IPs to a foreign font CDN (RU data rules).

import { useEffect, useState } from 'react'

export type ThemeId = 'warm' | 'neon' | 'brutal'
export type FontId = 'sans' | 'serif' | 'mono' | 'rounded'
export type Density = 'comfortable' | 'compact'
export type SizeId = 's' | 'm' | 'l'

export interface Settings {
  theme: ThemeId
  font: FontId
  density: Density
  size: SizeId
}

export const DEFAULT_SETTINGS: Settings = {
  theme: 'warm', // тёплая журнальная — оранж + слива
  font: 'sans',
  density: 'comfortable',
  size: 'm',
}

export const THEMES: { id: ThemeId; label: string; hint: string }[] = [
  { id: 'warm', label: 'Тёплая журнальная', hint: 'крем · жжёный оранж · слива' },
  { id: 'neon', label: 'Тёмная неоновая', hint: 'фиолет · неон-оранж' },
  { id: 'brutal', label: 'Необрутализм', hint: 'контраст · моно · рамки' },
]

export const FONTS: { id: FontId; label: string }[] = [
  { id: 'sans', label: 'Гротеск' },
  { id: 'serif', label: 'Сериф' },
  { id: 'rounded', label: 'Округлый' },
  { id: 'mono', label: 'Моно' },
]

export const SIZES: { id: SizeId; label: string }[] = [
  { id: 's', label: 'S' },
  { id: 'm', label: 'M' },
  { id: 'l', label: 'L' },
]

const KEY = 'fm.settings'

export function loadSettings(): Settings {
  try {
    const raw = localStorage.getItem(KEY)
    if (raw) return { ...DEFAULT_SETTINGS, ...(JSON.parse(raw) as Partial<Settings>) }
  } catch {
    // ignore corrupt storage
  }
  return DEFAULT_SETTINGS
}

export function applySettings(s: Settings): void {
  const r = document.documentElement
  r.dataset.theme = s.theme
  r.dataset.font = s.font
  r.dataset.density = s.density
  r.dataset.size = s.size
}

// useSettings keeps state in sync with <html> and localStorage.
export function useSettings(): [Settings, (next: Partial<Settings>) => void] {
  const [settings, setSettings] = useState<Settings>(loadSettings)

  useEffect(() => {
    applySettings(settings)
    try {
      localStorage.setItem(KEY, JSON.stringify(settings))
    } catch {
      // ignore storage write failures (private mode etc.)
    }
  }, [settings])

  const update = (next: Partial<Settings>) => setSettings((s) => ({ ...s, ...next }))
  return [settings, update]
}
