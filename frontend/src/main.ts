import './style.css'
import App from './App.svelte'

// Disable the browser's native right-click context menu — this is a desktop app.
window.addEventListener('contextmenu', (e) => e.preventDefault())

// Block trackpad pinch-to-zoom (sent as wheel events with ctrlKey = true in webviews).
window.addEventListener('wheel', (e) => { if (e.ctrlKey) e.preventDefault() }, { passive: false })



const app = new App({
  target: document.getElementById('app')
})

export default app
