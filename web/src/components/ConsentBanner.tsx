import { useState } from 'react'

const KEY = 'fm.consent'

// Informs about cookies / personal-data processing and records consent
// (152-ФЗ). No third-party trackers are loaded, so this is consent for our own
// processing of search data and functional cookies only.
export function ConsentBanner() {
  const [accepted, setAccepted] = useState<boolean>(() => {
    try {
      return localStorage.getItem(KEY) === '1'
    } catch {
      return false
    }
  })

  if (accepted) return null

  function accept() {
    try {
      localStorage.setItem(KEY, '1')
    } catch {
      // ignore
    }
    setAccepted(true)
  }

  return (
    <div className="consent" role="dialog" aria-label="Cookie и персональные данные">
      <p className="consent__text">
        Мы используем cookie и обрабатываем данные ваших поисков, чтобы сервис работал. Продолжая,
        вы соглашаетесь с{' '}
        <a href="#/legal/privacy">Политикой обработки персональных данных</a>. Сторонних трекеров
        нет.
      </p>
      <button className="btn btn--go" onClick={accept}>
        Принять
      </button>
    </div>
  )
}
