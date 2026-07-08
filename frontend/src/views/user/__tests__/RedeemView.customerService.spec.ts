import { readFileSync } from 'node:fs'
import { dirname, resolve } from 'node:path'
import { fileURLToPath } from 'node:url'

import { describe, expect, it } from 'vitest'

const viewPath = resolve(dirname(fileURLToPath(import.meta.url)), '../RedeemView.vue')
const source = readFileSync(viewPath, 'utf8')

describe('RedeemView customer service contact', () => {
  it('embeds configured contact info in the customer service invite URL when present', () => {
    expect(source).toContain('const customerServiceInviteURL = ref')
    expect(source).toContain('settings.customer_service_invite_url?.trim()')
    expect(source).toContain('v-if="customerServiceInviteURL"')
    expect(source).toContain(':href="customerServiceInviteURL"')
    expect(source).toContain('{{ contactInfo }}')
  })
})
