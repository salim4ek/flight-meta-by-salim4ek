import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css'
import App from './App'
import { applySettings, loadSettings } from './settings'

// Apply persisted theme/font/size before first paint to avoid a flash.
applySettings(loadSettings())

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <App />
  </StrictMode>,
)
