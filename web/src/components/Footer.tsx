import { ORG } from '../legal'

export function Footer() {
  const year = new Date().getFullYear()
  return (
    <footer className="foot">
      <nav className="foot__links">
        <a href="#/legal/privacy">Политика обработки ПД</a>
        <a href="#/legal/consent">Согласие на обработку ПД</a>
        <a href="#/legal/terms">Пользовательское соглашение</a>
      </nav>
      <div className="foot__org">
        {ORG.name} · ИНН {ORG.inn} · ОГРН {ORG.ogrn} ·{' '}
        <a href={`mailto:${ORG.email}`}>{ORG.email}</a>
      </div>
      <div className="foot__note">
        Сервис — метапоиск: билеты не продаёт, наценки нет. Покупка — напрямую у авиакомпаний и
        агентств. Партнёрские ссылки помечены «Реклама». © {year} {ORG.brand}
      </div>
    </footer>
  )
}
