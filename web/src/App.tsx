import { useState } from 'react'
import { search, type SearchParams, type SearchResult } from './api'
import { SearchForm } from './components/SearchForm'
import { Results } from './components/Results'

export default function App() {
  const [result, setResult] = useState<SearchResult | null>(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  async function handleSearch(p: SearchParams) {
    setLoading(true)
    setError(null)
    try {
      setResult(await search(p))
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Не удалось выполнить поиск')
      setResult(null)
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="app">
      <header className="hero">
        <h1 className="hero__title">flight-meta</h1>
        <p className="hero__sub">
          Находим маршруты из отдельных билетов, которых нет у агрегаторов — через любой хаб,
          с учётом транзитных виз и безопасного времени стыковки. Без наценки: показываем и ведём покупать напрямую.
        </p>
      </header>

      <SearchForm onSearch={handleSearch} loading={loading} />

      {error && <p className="error">{error}</p>}
      {result && <Results result={result} />}
      {!result && !error && (
        <p className="hint">
          Введите маршрут и нажмите «Найти». Данные сейчас демонстрационные (mock-источник) —
          реальные источники (Kiwi и др.) подключаются в следующих фазах.
        </p>
      )}

      <footer className="foot">Фаза 2 · UI на mock-бэкенде</footer>
    </div>
  )
}
