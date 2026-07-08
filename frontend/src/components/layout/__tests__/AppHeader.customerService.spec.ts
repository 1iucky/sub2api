import { readFileSync } from 'node:fs'
import { dirname, resolve } from 'node:path'
import { fileURLToPath } from 'node:url'

import { describe, expect, it } from 'vitest'

const componentPath = resolve(dirname(fileURLToPath(import.meta.url)), '../AppHeader.vue')
const componentSource = readFileSync(componentPath, 'utf8')

describe('AppHeader customer service entry', () => {
  it('renders a top-bar contact support icon before announcements', () => {
    const contactIndex = componentSource.indexOf('data-testid="header-contact-support"')
    const announcementIndex = componentSource.indexOf('<AnnouncementBell')

    expect(contactIndex).toBeGreaterThan(-1)
    expect(announcementIndex).toBeGreaterThan(-1)
    expect(contactIndex).toBeLessThan(announcementIndex)
  })

  it('uses invite URL for direct links and falls back to a contact-info popover', () => {
    expect(componentSource).toContain('customerServiceInviteURL')
    expect(componentSource).toContain('contactSupportPopoverOpen')
    expect(componentSource).toContain('data-testid="header-contact-support-popover"')
  })
})
