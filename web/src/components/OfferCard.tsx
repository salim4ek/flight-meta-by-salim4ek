import type { Offer, VisaStatus } from '../api'
import {
  formatPrice,
  formatTimeUTC,
  formatDayUTC,
  humanDurationMs,
  humanDurationNs,
} from '../format'

function visaBadge(status?: VisaStatus): { kind: string; label: string } | null {
  switch (status) {
    case 'no_visa':
      return { kind: 'ok', label: 'Транзит без визы' }
    case 'twov':
      return { kind: 'twov', label: 'Безвиз-транзит (TWOV)' }
    case 'visa_required':
      return { kind: 'need', label: 'Нужна транзитная виза' }
    case 'unknown':
      return { kind: 'unknown', label: 'Виза: уточните' }
    default:
      return null
  }
}

// Card layout follows established flight-search conventions (Google Flights /
// Skyscanner): prominent price, a depart→arrive timeline with stop pills, and
// airline + connection badges. Self-transfer offers carry an explicit warning.
export function OfferCard({ offer }: { offer: Offer }) {
  const segs = offer.segments
  const first = segs[0]
  const last = segs[segs.length - 1]
  const stops = Math.max(0, segs.length - 1)
  const durMs = new Date(last.arriveUtc).getTime() - new Date(first.departUtc).getTime()
  const isSelf = offer.connection === 'self_transfer'
  // Defense-in-depth: only follow http(s) deep links. Real partner links come
  // from external sources later; never render javascript:/data: schemes.
  const buyHref = /^https?:\/\//i.test(offer.deepLink) ? offer.deepLink : undefined

  return (
    <article className={`card${isSelf ? ' card--self' : ''}`}>
      <div className="card__main">
        <div className="card__badges">
          {offer.unique && <span className="badge badge--unique">Нет у агрегаторов</span>}
          {isSelf ? (
            <span className="badge badge--self">Само-стыковка · отдельные билеты</span>
          ) : (
            <span className="badge badge--official">Единый билет</span>
          )}
          <span className="badge badge--muted">{stops === 0 ? 'Прямой' : `${stops} пересадка(и)`}</span>
          <span className="badge badge--muted">{offer.carriers.join(' · ')}</span>
        </div>

        <div className="route">
          <div className="route__end">
            <div className="route__time">{formatTimeUTC(first.departUtc)}</div>
            <div className="route__code">{first.from}</div>
            <div className="route__day">{formatDayUTC(first.departUtc)}</div>
          </div>

          <div className="route__mid">
            <div className="route__dur">{humanDurationMs(durMs)} в пути</div>
            <div className="route__line">
              {segs.slice(0, -1).map((s, i) => (
                <span key={i} className="route__stop" title={`Пересадка в ${s.to}`}>
                  {s.to}
                </span>
              ))}
            </div>
            <div className="route__stops">
              {stops === 0 ? 'без пересадок' : `${stops} пересадка(и)`}
            </div>
          </div>

          <div className="route__end">
            <div className="route__time">{formatTimeUTC(last.arriveUtc)}</div>
            <div className="route__code">{last.to}</div>
            <div className="route__day">{formatDayUTC(last.arriveUtc)}</div>
          </div>
        </div>

        {offer.layovers && offer.layovers.length > 0 && (
          <ul className="layovers">
            {offer.layovers.map((l, i) => {
              const v = visaBadge(l.visaStatus)
              return (
                <li key={i} className="layover">
                  <div className="layover__top">
                    <span className="layover__air">{l.airport}</span>
                    <span className="layover__dur">стыковка {humanDurationNs(l.durationNs)}</span>
                    {l.selfTransfer && <span className="layover__tag">смена билета</span>}
                    {v && <span className={`visa visa--${v.kind}`}>{v.label}</span>}
                  </div>
                  {l.transitNote && <div className="layover__note">{l.transitNote}</div>}
                </li>
              )
            })}
          </ul>
        )}

        {isSelf && (
          <p className="warn">
            ⚠️ Раздельные билеты: при опоздании на стыковку перевозчик не отвечает за пересадку,
            багаж нужно получить и сдать заново. Закладывайте запас по времени.
          </p>
        )}
      </div>

      <div className="card__buy">
        <div className="price">{formatPrice(offer.priceMinor, offer.currency)}</div>
        {buyHref ? (
          <a className="btn btn--go" href={buyHref} target="_blank" rel="noopener noreferrer">
            Купить
          </a>
        ) : (
          <span className="btn btn--go" aria-disabled="true" style={{ opacity: 0.6 }}>
            Ссылка недоступна
          </span>
        )}
        <div className="card__source">источник: {offer.source}</div>
      </div>
    </article>
  )
}
