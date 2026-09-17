import { readFileSync } from 'node:fs'
import { dirname, resolve } from 'node:path'
import { fileURLToPath } from 'node:url'

import { describe, expect, it } from 'vitest'

const componentPath = resolve(dirname(fileURLToPath(import.meta.url)), '../AnnouncementBell.vue')
const componentSource = readFileSync(componentPath, 'utf8')

describe('AnnouncementBell drawer presentation', () => {
  it('uses a right-side drawer instead of the legacy centered modal list', () => {
    expect(componentSource).toContain('data-testid="announcement-drawer"')
    expect(componentSource).toContain('openDrawer')
    expect(componentSource).toContain('closeDrawer')
    expect(componentSource).not.toContain('openModal')
  })

  it('renders announcement content inline in the drawer and supports long-content expansion', () => {
    expect(componentSource).toContain('data-testid="announcement-timeline"')
    expect(componentSource).toContain('shouldCollapseAnnouncement')
    expect(componentSource).toContain('toggleAnnouncementExpanded')
    expect(componentSource).toContain('renderMarkdown(item.content)')
  })

  it('does not depend on local-only announcement category fields', () => {
    expect(componentSource).not.toContain('AnnouncementCategory')
    expect(componentSource).not.toContain('item.category')
    expect(componentSource).not.toContain('categoryDefinitions')
  })
})
