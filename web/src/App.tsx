import { useEffect, useState } from 'react'
import {
  search,
  calendar,
  type SearchParams,
  type SearchResult,
  type CalendarResult,
} from './api'
import { SearchForm } from './components/SearchForm'
import { Results } from './components/Results'
import { Calendar } from './components/Calendar'
import { SettingsPanel } from './components/Settings'
import { ConsentBanner } from './components/ConsentBanner'
import { LegalPage } from './components/Legal'
import { Footer } from './components/Footer'
import { useSettings } from './settings'

const CAL_WINDOW = 3

export default function App() {
  const [settings, update] = useSettings()
  const [showSettings, setShowSettings] = useState(false)

  const [result, setResult] = useState<SearchResult | null>(null)
  const [cal, setCal] = useState<CalendarResult | null>(null)
  const [lastParams, setLastParams] = useState<SearchParams | null>(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const [route, setRoute] = useState(() => window.location.hash)
  useEffect(() => {
    const onHash = () => setRoute(window.location.hash)
    window.addEventListener('hashchange', onHash)
    return () => window.removeEventListener('hashchange', onHash)
  }, [])

  async function runSearch(p: SearchParams) {
    setLastParams(p)
    setLoading(true)
    setError(null)
    try {
      const [r, c] = await Promise.all([search(p), calendar(p, CAL_WINDOW)])
      setResult(r)
      setCal(c)
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Не удалось выполнить поиск')
      setResult(null)
      setCal(null)
    } finally {
      setLoading(false)
    }
  }

  function pickDate(date: string) {
    if (lastParams) void runSearch({ ...lastParams, depart: date })
  }

  if (route.startsWith('#/legal')) {
    return (
      <div className="app">
        <LegalPage route={route} />
        <Footer />
      </div>
    )
  }

  return (
    <div className="app">
      <header className="hero">
        <div className="hero__bar">
          <h1 className="hero__title">flight-meta</h1>
          <button className="icon-btn" onClick={() => setShowSettings(true)} aria-label="Настройки">
            ⚙
          </button>
        </div>
        <p className="hero__sub">
          Маршруты из отдельных билетов, которых нет у агрегаторов — через любой хаб, с учётом
          транзитных виз и безопасного времени стыковки. Без наценки: показываем и ведём покупать
          напрямую.
        </p>
      </header>

      <SearchForm onSearch={runSearch} loading={loading} />

      {error && <p className="error">{error}</p>}
      {cal && lastParams && <Calendar data={cal} selected={lastParams.depart} onPick={pickDate} />}
      {result && <Results result={result} />}
      {!result && !error && (
        <p className="hint">
          Введите маршрут и нажмите «Найти». Данные сейчас демонстрационные (свой комбайнер на
          mock-ценах) — реальный источник цен подключается следующим.
        </p>
      )}

      <Footer />

      {showSettings && (
        <SettingsPanel settings={settings} update={update} onClose={() => setShowSettings(false)} />
      )}
      <ConsentBanner />
    </div>
  )
}
