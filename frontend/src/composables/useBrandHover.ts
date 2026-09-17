/**
 * Shared brand-block micro-interaction for the Factory-style visual system.
 *
 * The composable only owns phase state. CSS owns all motion so the same logo
 * spin / reverse rewind and per-letter flip can be reused by the homepage,
 * login page, and console sidebar without duplicating animation logic.
 */
import { computed, onBeforeUnmount, readonly, ref, type ComputedRef, type Ref } from 'vue'

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
  brandClass: ComputedRef<Record<string, boolean>>
} {
  const phase = ref<BrandPhase>('idle')
  let rewindTimer: ReturnType<typeof setTimeout> | null = null

  const brandClass = computed(() => ({
    'is-spinning': phase.value === 'spin',
    'is-rewinding': phase.value === 'rewind',
  }))

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
    if (phase.value !== 'spin') return
    phase.value = 'rewind'
    clearTimer()
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
    brandClass,
  }
}
