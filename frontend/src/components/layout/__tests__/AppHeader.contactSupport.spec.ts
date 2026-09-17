import { readFileSync } from 'node:fs'
import { dirname, resolve } from 'node:path'
import { fileURLToPath } from 'node:url'

import { describe, expect, it } from 'vitest'

const componentPath = resolve(dirname(fileURLToPath(import.meta.url)), '../AppHeader.vue')
const componentSource = readFileSync(componentPath, 'utf8')

describe('AppHeader contact support entry', () => {
  it('renders a top-bar contact support icon before announcements', () => {
    const contactIndex = componentSource.indexOf('data-testid="header-contact-support"')
    const announcementIndex = componentSource.indexOf('<AnnouncementBell')

    expect(contactIndex).toBeGreaterThan(-1)
    expect(announcementIndex).toBeGreaterThan(-1)
    expect(contactIndex).toBeLessThan(announcementIndex)
  })

  it('uses existing public contact info for a popover without requiring backend settings changes', () => {
    expect(componentSource).toContain('contactSupportPopoverOpen')
    expect(componentSource).toContain('showContactSupportEntry')
    expect(componentSource).toContain('data-testid="header-contact-support-popover"')
    expect(componentSource).toContain('contactInfo')
    expect(componentSource).not.toContain('customerServiceInviteURL')
  })
})
