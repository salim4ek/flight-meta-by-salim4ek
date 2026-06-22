// Small display formatters. Times from the API are UTC instants; until
// per-airport timezones land (later phase) we render them in UTC to stay
// deterministic and avoid implying a wrong local time.

export function formatPrice(minor: number, currency: string): string {
  try {
    return new Intl.NumberFormat('ru-RU', {
      style: 'currency',
      currency: currency || 'RUB',
      maximumFractionDigits: 0,
    }).format(minor / 100)
  } catch {
    return `${Math.round(minor / 100)} ${currency}`
  }
}

export function formatTimeUTC(iso: string): string {
  const d = new Date(iso)
  const hh = String(d.getUTCHours()).padStart(2, '0')
  const mm = String(d.getUTCMinutes()).padStart(2, '0')
  return `${hh}:${mm}`
}

export function formatDayUTC(iso: string): string {
  const d = new Date(iso)
  return new Intl.DateTimeFormat('ru-RU', {
    day: '2-digit',
    month: 'short',
    timeZone: 'UTC',
  }).format(d)
}

export function humanDurationMs(ms: number): string {
  const totalMin = Math.max(0, Math.round(ms / 60000))
  const h = Math.floor(totalMin / 60)
  const m = totalMin % 60
  return h === 0 ? `${m}м` : `${h}ч ${String(m).padStart(2, '0')}м`
}

export function humanDurationNs(ns: number): string {
  return humanDurationMs(ns / 1e6)
}
