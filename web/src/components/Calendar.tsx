import type { CalendarResult } from '../api'
import { formatPrice } from '../format'

function dayLabel(iso: string): { dow: string; dm: string } {
  const d = new Date(iso + 'T00:00:00Z')
  const dow = new Intl.DateTimeFormat('ru-RU', { weekday: 'short', timeZone: 'UTC' }).format(d)
  const dm = new Intl.DateTimeFormat('ru-RU', { day: '2-digit', month: 'short', timeZone: 'UTC' }).format(d)
  return { dow, dm }
}

export function Calendar({
  data,
  selected,
  onPick,
}: {
  data: CalendarResult
  selected: string
  onPick: (date: string) => void
}) {
  return (
    <section className="cal" aria-label="Цены по датам">
      <div className="cal__title">Гибкие даты — рядом может быть дешевле</div>
      <div className="cal__strip">
        {data.days.map((d) => {
          const { dow, dm } = dayLabel(d.date)
          const cls =
            'cal__day' +
            (d.date === selected ? ' is-selected' : '') +
            (d.cheapest ? ' is-cheapest' : '') +
            (!d.hasOffers ? ' is-empty' : '')
          return (
            <button
              key={d.date}
              className={cls}
              onClick={() => onPick(d.date)}
              disabled={!d.hasOffers}
              title={d.cheapest ? 'Самый дешёвый день в окне' : undefined}
            >
              <span className="cal__dow">{dow}</span>
              <span className="cal__dm">{dm}</span>
              <span className="cal__price">
                {d.hasOffers ? formatPrice(d.priceMinor, d.currency) : '—'}
              </span>
              {d.cheapest && <span className="cal__tag">мин</span>}
            </button>
          )
        })}
      </div>
    </section>
  )
}
