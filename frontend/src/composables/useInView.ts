/**
 * useInView — IntersectionObserver composable (factory.ai homepage motion).
 *
 * Mirrors factory's reveal pattern: a `ref` + an IntersectionObserver that sets
 * `inView = true` once the element enters the viewport, then disconnects. The
 * reveal itself is driven purely by CSS keyframes gated on the `data-ready`
 * attribute (attribute-toggle pattern), so the JS layer stays tiny.
 *
 * Reduced-motion: short-circuits to `inView = true` immediately so all reveal
 * sections render fully visible with no animation.
 */
import { onMounted, onBeforeUnmount, readonly, ref, type Ref } from 'vue'

const REVEAL_ROOT_MARGIN = '0px 0px -80px 0px'
const REVEAL_THRESHOLD = 0.3

function prefersReducedMotion(): boolean {
  return (
    typeof window !== 'undefined' &&
    typeof window.matchMedia === 'function' &&
    window.matchMedia('(prefers-reduced-motion: reduce)').matches
  )
}

export interface UseInViewOptions {
  /** Once revealed, stop observing (default: true). */
  once?: boolean
  /** IntersectionObserver rootMargin (default: factory's bottom bias). */
  rootMargin?: string
  /** IntersectionObserver threshold (default: 0.3). */
  threshold?: number
}

export function useInView<T extends Element = Element>(options: UseInViewOptions = {}): {
  target: Ref<T | null>
  inView: Readonly<Ref<boolean>>
  /** Imperatively mark as in-view (used for SSR / reduced-motion fallbacks). */
  reveal: () => void
} {
  const { once = true, rootMargin = REVEAL_ROOT_MARGIN, threshold = REVEAL_THRESHOLD } = options
  const target = ref<T | null>(null) as Ref<T | null>
  const inView = ref(false)
  let observer: IntersectionObserver | null = null

  function reveal() {
    if (inView.value) return
    inView.value = true
    if (target.value) target.value.setAttribute('data-ready', 'true')
  }

  function cleanup() {
    if (observer) {
      observer.disconnect()
      observer = null
    }
  }

  onMounted(() => {
    // Reduced-motion: never animate — show everything immediately.
    if (prefersReducedMotion()) {
      reveal()
      return
    }

    const el = target.value
    if (!el || typeof IntersectionObserver === 'undefined') {
      reveal()
      return
    }

    observer = new IntersectionObserver(
      (entries) => {
        for (const entry of entries) {
          if (entry.isIntersecting) {
            reveal()
            if (once) cleanup()
          }
        }
      },
      { rootMargin, threshold }
    )
    observer.observe(el)
  })

  onBeforeUnmount(cleanup)

  return { target, inView: readonly(inView), reveal }
}

/**
 * v-reveal — directive form of {@link useInView}.
 *
 * Slides `data-ready="true"` onto the host element when it enters the viewport.
 * Pair with CSS keyframes gated under `[data-ready="true"]` for a JS-light,
 * attribute-toggle reveal (factory's pattern). Respects reduced-motion by
 * setting the attribute on mount.
 *
 * Register once in main.ts:  `app.directive('reveal', vReveal)`
 */
export const vReveal = {
  mounted(el: Element & { __revealObserver?: IntersectionObserver | null }) {
    // Tag the host as a reveal target so the CSS `[data-reveal]` rules apply
    // (factory.ai attribute-toggle pattern). `data-ready` is flipped to "true"
    // when the element enters the viewport.
    el.setAttribute('data-reveal', 'true')
    if (prefersReducedMotion() || typeof IntersectionObserver === 'undefined') {
      el.setAttribute('data-ready', 'true')
      return
    }
    const observer = new IntersectionObserver(
      (entries, obs) => {
        for (const entry of entries) {
          if (entry.isIntersecting) {
            entry.target.setAttribute('data-ready', 'true')
            obs.disconnect()
          }
        }
      },
      { rootMargin: REVEAL_ROOT_MARGIN, threshold: REVEAL_THRESHOLD }
    )
    ;(el as Element & { __revealObserver?: IntersectionObserver | null }).__revealObserver = observer
    observer.observe(el)
  },
  unmounted(el: Element & { __revealObserver?: IntersectionObserver | null }) {
    el.__revealObserver?.disconnect()
    el.__revealObserver = null
  }
}
