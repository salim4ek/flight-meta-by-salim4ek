// API client + types mirroring the Go backend's JSON (internal/offer, search).

export type Connection = 'official' | 'self_transfer'

export interface Segment {
  from: string
  to: string
  departUtc: string
  arriveUtc: string
  marketingCarrier: string
  flightNumber: string
  cabin?: string
}

export type VisaStatus = 'no_visa' | 'twov' | 'visa_required' | 'unknown'
export type ConnRisk = 'safe' | 'risky' | 'infeasible'

export interface Layover {
  airport: string
  durationNs: number // Go time.Duration marshals as int64 nanoseconds
  selfTransfer: boolean
  risk?: ConnRisk
  visaStatus?: VisaStatus
  transitNote?: string
}

export interface Offer {
  id: string
  source: string
  segments: Segment[]
  layovers?: Layover[]
  connection: Connection
  priceMinor: number
  currency: string
  deepLink: string
  carriers: string[]
  unique: boolean
}

export interface SourceStat {
  name: string
  count: number
  error?: string
  tookMs: number
}

export interface SearchResult {
  offers: Offer[]
  sources: SourceStat[]
  visaDisclaimer?: string
}

export type StopsMode = '' | 'direct' | 'one' | 'one_plus'

export interface SearchParams {
  origin: string
  destination: string
  depart: string
  ret?: string
  passengers?: number
  passport?: string
  currency?: string
  stops?: StopsMode
  airlines?: string
  excludeAirlines?: string
  selfTransfer?: boolean
  visaFreeTransit?: boolean
  hideInfeasible?: boolean
}

export interface CalendarDay {
  date: string // YYYY-MM-DD
  priceMinor: number
  currency: string
  hasOffers: boolean
  cheapest: boolean
}

export interface CalendarResult {
  days: CalendarDay[]
}

const API_BASE = (import.meta.env.VITE_API_BASE as string | undefined) ?? '/api'

function buildParams(p: SearchParams): URLSearchParams {
  const q = new URLSearchParams()
  q.set('origin', p.origin)
  q.set('destination', p.destination)
  q.set('depart', p.depart)
  if (p.ret) q.set('return', p.ret)
  if (p.passengers) q.set('passengers', String(p.passengers))
  if (p.passport) q.set('passport', p.passport)
  if (p.currency) q.set('currency', p.currency)
  if (p.stops) q.set('stops', p.stops)
  if (p.airlines) q.set('airlines', p.airlines)
  if (p.excludeAirlines) q.set('exclude_airlines', p.excludeAirlines)
  if (p.selfTransfer === false) q.set('self_transfer', 'false')
  if (p.visaFreeTransit) q.set('visa_free_transit', 'true')
  if (p.hideInfeasible) q.set('hide_infeasible', 'true')
  return q
}

async function getJSON<T>(path: string, signal?: AbortSignal): Promise<T> {
  const res = await fetch(`${API_BASE}${path}`, { signal })
  if (!res.ok) {
    let msg = `Ошибка ${res.status}`
    try {
      const body = (await res.json()) as { error?: string }
      if (body?.error) msg = body.error
    } catch {
      // non-JSON error body — keep the status message
    }
    throw new Error(msg)
  }
  return (await res.json()) as T
}

export function search(p: SearchParams, signal?: AbortSignal): Promise<SearchResult> {
  return getJSON<SearchResult>(`/search?${buildParams(p).toString()}`, signal)
}

export function calendar(p: SearchParams, window: number, signal?: AbortSignal): Promise<CalendarResult> {
  const q = buildParams(p)
  q.set('window', String(window))
  return getJSON<CalendarResult>(`/calendar?${q.toString()}`, signal)
}
