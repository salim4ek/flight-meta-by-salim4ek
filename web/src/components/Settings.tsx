import type { Settings } from '../settings'
import { THEMES, FONTS, SIZES } from '../settings'

export function SettingsPanel({
  settings,
  update,
  onClose,
}: {
  settings: Settings
  update: (next: Partial<Settings>) => void
  onClose: () => void
}) {
  return (
    <div className="sheet" role="dialog" aria-label="Настройки" onClick={onClose}>
      <div className="sheet__panel" onClick={(e) => e.stopPropagation()}>
        <div className="sheet__head">
          <h2>Настройки</h2>
          <button className="icon-btn" onClick={onClose} aria-label="Закрыть">
            ✕
          </button>
        </div>

        <section className="set">
          <span className="set__label">Тема</span>
          <div className="set__opts set__opts--theme">
            {THEMES.map((t) => (
              <button
                key={t.id}
                className={`swatch swatch--${t.id}${settings.theme === t.id ? ' is-on' : ''}`}
                onClick={() => update({ theme: t.id })}
              >
                <span className="swatch__dots" aria-hidden="true">
                  <i className="d1" />
                  <i className="d2" />
                  <i className="d3" />
                </span>
                <span className="swatch__name">{t.label}</span>
                <span className="swatch__hint">{t.hint}</span>
              </button>
            ))}
          </div>
        </section>

        <section className="set">
          <span className="set__label">Шрифт</span>
          <div className="seg">
            {FONTS.map((f) => (
              <button
                key={f.id}
                className={`seg__btn${settings.font === f.id ? ' seg__btn--on' : ''}`}
                onClick={() => update({ font: f.id })}
              >
                {f.label}
              </button>
            ))}
          </div>
        </section>

        <section className="set">
          <span className="set__label">Размер</span>
          <div className="seg">
            {SIZES.map((s) => (
              <button
                key={s.id}
                className={`seg__btn${settings.size === s.id ? ' seg__btn--on' : ''}`}
                onClick={() => update({ size: s.id })}
              >
                {s.label}
              </button>
            ))}
          </div>
        </section>

        <section className="set">
          <span className="set__label">Плотность</span>
          <div className="seg">
            <button
              className={`seg__btn${settings.density === 'comfortable' ? ' seg__btn--on' : ''}`}
              onClick={() => update({ density: 'comfortable' })}
            >
              Просторно
            </button>
            <button
              className={`seg__btn${settings.density === 'compact' ? ' seg__btn--on' : ''}`}
              onClick={() => update({ density: 'compact' })}
            >
              Компактно
            </button>
          </div>
        </section>
      </div>
    </div>
  )
}
