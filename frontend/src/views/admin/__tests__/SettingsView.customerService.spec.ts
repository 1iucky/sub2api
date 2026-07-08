import { readFileSync } from 'node:fs'
import { dirname, resolve } from 'node:path'
import { fileURLToPath } from 'node:url'

import { describe, expect, it } from 'vitest'

const viewPath = resolve(dirname(fileURLToPath(import.meta.url)), '../SettingsView.vue')
const source = readFileSync(viewPath, 'utf8')

describe('SettingsView customer service widget settings', () => {
  it('contains the customer service invite URL field next to contact info', () => {
    expect(source).toContain('form.contact_info')
    expect(source).toContain('form.customer_service_invite_url')
    expect(source).toContain('admin.settings.site.customerServiceInviteURL')
    expect(source).not.toContain('form.customer_service_enabled')
    expect(source).not.toContain('form.customer_service_qr_code_url')
  })
})
