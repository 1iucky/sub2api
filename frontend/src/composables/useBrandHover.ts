/**
 * useBrandHover — shared brand-block micro-interaction (factory.ai aesthetic).
 *
 * Drives the logo spin / reverse-rewind + the per-letter 3D flip from a single
 * reactive `phase` so all three brand surfaces (AppSidebar, HomeView nav,
 * AuthLayout) share identical behavior without duplicating state.
 *
 * Lifecycle:
 *   idle  → (mouseenter) →  spin   (forward infinite, ~6s linear)
 *   spin  → (mouseleave) →  rewind (one-shot reverse unwind to 0°, ~1.2s ease-out)
 *   rewind → (ends)       →  idle
 *
 * `rewind` only ever runs AFTER the first hover — initial state is `idle`, so
 * page-load shows the logo static at 0°. The phase is exposed to the template
 * via a `data-phase` attribute on the `.brand-hover` root, which the CSS in
 * `theme-override.css` keys off.
 *
 * The composable itself is pure state — all motion lives in CSS, gated under
 * `prefers-reduced-motion: no-preference`. Under reduce the phase is forced to
 * `idle` (handlers short-circuit) so nothing animates.
 */
import { onBeforeUnmount, ref, readonly, type Ref } from 'vue'

export type BrandPhase = 'idle' | 'spin' | 'rewind'

function prefersReducedMotion(): boolean {
  return (
    typeof window !== 'undefined' &&
    typeof window.matchMedia === 'function' &&
    window.matchMedia('(prefers-reduced-motion: reduce)').matches
  )
}

export function useBrandHover(): {
  phase: Readonly<Ref<BrandPhase>>
  onEnter: () => void
  onLeave: () => void
  /** Bind onto the `.brand-hover` root element: <div v-bind="brandProps"> */
  brandProps: { 'data-phase': () => BrandPhase }
} {
  const phase = ref<BrandPhase>('idle')
  let rewindTimer: ReturnType<typeof setTimeout> | null = null

  function clearTimer() {
    if (rewindTimer !== null) {
      clearTimeout(rewindTimer)
      rewindTimer = null
    }
  }

  function onEnter() {
    if (prefersReducedMotion()) return
    clearTimer()
    phase.value = 'spin'
  }

  function onLeave() {
    if (prefersReducedMotion()) {
      phase.value = 'idle'
      return
    }
    // Only rewind if we were actually spinning — never fire rewind on page load.
    if (phase.value !== 'spin') return
    phase.value = 'rewind'
    clearTimer()
    // Match the CSS `brand-spin-back` duration (1200ms). After the unwind
    // finishes, drop back to idle so the next hover starts cleanly from 0°.
    rewindTimer = setTimeout(() => {
      phase.value = 'idle'
      rewindTimer = null
    }, 1200)
  }

  onBeforeUnmount(clearTimer)

  return {
    phase: readonly(phase),
    onEnter,
    onLeave,
    brandProps: {
      // Vue 3 supports function-form attribute values for reactivity —
      // re-renders when `phase.value` changes.
      'data-phase': () => phase.value,
    },
  }
}
