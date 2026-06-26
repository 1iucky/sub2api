import { mount } from '@vue/test-utils'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { nextTick } from 'vue'
import Select from '../Select.vue'

vi.mock('vue-i18n', () => ({
  useI18n: () => ({
    t: (_key: string, fallback?: string) => fallback || _key,
  }),
}))

describe('Select', () => {
  afterEach(() => {
    document.body.innerHTML = ''
  })

  it('emits search text from the dropdown search input', async () => {
    const wrapper = mount(Select, {
      attachTo: document.body,
      props: {
        modelValue: null,
        options: [],
        searchable: true,
      },
    })

    await wrapper.get('button').trigger('click')
    const input = document.body.querySelector('.select-search-input') as HTMLInputElement
    input.value = 'glm'
    input.dispatchEvent(new Event('input'))
    await nextTick()

    expect(wrapper.emitted('search')).toEqual([['glm']])
    wrapper.unmount()
  })
})
