import type { SearchResult } from '../api'
import { OfferCard } from './OfferCard'

export function Results({ result }: { result: SearchResult }) {
  if (result.offers.length === 0) {
    return <p className="empty">Ничего не найдено по этим фильтрам.</p>
  }
  return (
    <section className="results">
      <div className="results__head">
        <span>
          <strong>{result.offers.length}</strong> вариант(ов) · дешёвые сверху
        </span>
        <span className="results__src">
          {result.sources
            .map((s) => `${s.name}: ${s.count}${s.error ? ' (ошибка)' : ''} · ${s.tookMs}мс`)
            .join('   ')}
        </span>
      </div>
      <div className="results__list">
        {result.offers.map((o) => (
          <OfferCard key={o.id} offer={o} />
        ))}
      </div>
      {result.visaDisclaimer && (
        <p className="disclaimer">ℹ️ {result.visaDisclaimer} Время рейсов указано в UTC.</p>
      )}
    </section>
  )
}
