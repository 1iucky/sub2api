import { ref } from 'vue'

const isDark = ref(typeof document !== 'undefined' && document.documentElement.classList.contains('dark'))

function syncThemeFromDocument() {
  if (typeof document === 'undefined') return
  isDark.value = document.documentElement.classList.contains('dark')
}

function setTheme(nextDark: boolean) {
  if (typeof document === 'undefined') return
  isDark.value = nextDark
  document.documentElement.classList.toggle('dark', nextDark)
  localStorage.setItem('theme', nextDark ? 'dark' : 'light')
}

function toggleTheme() {
  setTheme(!isDark.value)
}

export function useTheme() {
  syncThemeFromDocument()

  return {
    isDark,
    setTheme,
    toggleTheme,
    syncThemeFromDocument
  }
}
