import { useState, type FormEvent } from 'react'
import type { SearchParams, StopsMode } from '../api'

const STOPS: { value: StopsMode; label: string }[] = [
  { value: '', label: 'Любые' },
  { value: 'direct', label: 'Прямой' },
  { value: 'one', label: '1 пересадка' },
  { value: 'one_plus', label: '1+' },
]

export function SearchForm({
  onSearch,
  loading,
}: {
  onSearch: (p: SearchParams) => void
  loading: boolean
}) {
  const [origin, setOrigin] = useState('MOW')
  const [destination, setDestination] = useState('BKK')
  const [depart, setDepart] = useState('2026-04-20')
  const [ret, setRet] = useState('')
  const [passengers, setPassengers] = useState(1)
  const [passport, setPassport] = useState('RU')
  const [stops, setStops] = useState<StopsMode>('')
  const [airlines, setAirlines] = useState('')
  const [excludeAirlines, setExcludeAirlines] = useState('')
  const [selfTransfer, setSelfTransfer] = useState(true)
  const [visaFreeTransit, setVisaFreeTransit] = useState(false)

  function submit(e: FormEvent) {
    e.preventDefault()
    onSearch({
      origin: origin.trim().toUpperCase(),
      destination: destination.trim().toUpperCase(),
      depart,
      ret: ret || undefined,
      passengers,
      passport: passport.trim().toUpperCase() || undefined,
      stops,
      airlines: airlines.trim() || undefined,
      excludeAirlines: excludeAirlines.trim() || undefined,
      selfTransfer,
      visaFreeTransit,
    })
  }

  return (
    <form className="form" onSubmit={submit}>
      <div className="form__row">
        <label className="field">
          <span>Откуда</span>
          <input value={origin} onChange={(e) => setOrigin(e.target.value)} maxLength={3} placeholder="MOW" required />
        </label>
        <label className="field">
          <span>Куда</span>
          <input value={destination} onChange={(e) => setDestination(e.target.value)} maxLength={3} placeholder="BKK" required />
        </label>
        <label className="field">
          <span>Туда</span>
          <input type="date" value={depart} onChange={(e) => setDepart(e.target.value)} required />
        </label>
        <label className="field">
          <span>Обратно</span>
          <input type="date" value={ret} onChange={(e) => setRet(e.target.value)} />
        </label>
        <label className="field field--sm">
          <span>Пасс.</span>
          <input type="number" min={1} max={9} value={passengers} onChange={(e) => setPassengers(Number(e.target.value))} />
        </label>
        <label className="field field--sm">
          <span>Паспорт</span>
          <input value={passport} onChange={(e) => setPassport(e.target.value)} maxLength={2} placeholder="RU" />
        </label>
      </div>

      <div className="form__row form__row--filters">
        <div className="seg" role="group" aria-label="Пересадки">
          {STOPS.map((s) => (
            <button
              key={s.value || 'any'}
              type="button"
              className={`seg__btn${stops === s.value ? ' seg__btn--on' : ''}`}
              onClick={() => setStops(s.value)}
            >
              {s.label}
            </button>
          ))}
        </div>
        <label className="field">
          <span>Только а/к</span>
          <input value={airlines} onChange={(e) => setAirlines(e.target.value)} placeholder="TK, EK" />
        </label>
        <label className="field">
          <span>Исключить а/к</span>
          <input value={excludeAirlines} onChange={(e) => setExcludeAirlines(e.target.value)} placeholder="SU" />
        </label>
        <label className="check">
          <input type="checkbox" checked={selfTransfer} onChange={(e) => setSelfTransfer(e.target.checked)} />
          <span>Само-стыковка</span>
        </label>
        <label className="check">
          <input
            type="checkbox"
            checked={visaFreeTransit}
            onChange={(e) => setVisaFreeTransit(e.target.checked)}
          />
          <span>Только безвиз-транзит</span>
        </label>
        <button className="btn btn--go" type="submit" disabled={loading}>
          {loading ? 'Ищем…' : 'Найти'}
        </button>
      </div>
    </form>
  )
}
