import { afterEach, describe, expect, it } from 'vitest'
import { enableAutoUnmount, mount } from '@vue/test-utils'
import { createI18n } from 'vue-i18n'
import GatewayDashboard from '../GatewayDashboard.vue'

enableAutoUnmount(afterEach)

const i18n = createI18n({
  legacy: false,
  locale: 'en',
  messages: {
    en: {
      home: {
        hero2: {
          demo: {
            title: () => 'Gateway dashboard',
            sidebarGateway: () => 'Gateway',
            sidebarPools: () => 'Pools',
            sidebarKeys: () => 'Keys',
            sidebarSettings: () => 'Settings',
            kpiRpm: () => 'RPM',
            kpiRpmValue: () => '1.2k',
            kpiP95: () => 'P95',
            kpiP95Value: () => '240ms',
            kpiUptime: () => 'Uptime',
            kpiUptimeValue: () => '99.9%',
            kpiKeys: () => 'Keys',
            kpiKeysValue: () => '348',
            streamTitle: () => 'Request stream',
            streamRoute: () => 'Route',
            streamUpstream: () => 'Upstream',
            streamStatus: () => 'OK',
            sparkClaude: () => 'Claude throughput',
            sparkOpenai: () => 'OpenAI throughput',
            sparkGemini: () => 'Gemini throughput'
          }
        }
      }
    }
  }
})

function mountDashboard() {
  return mount(GatewayDashboard, {
    attachTo: document.body,
    props: { t: i18n.global.t }
  })
}

afterEach(() => {
  document.body.innerHTML = ''
})

describe('GatewayDashboard personalized presentation', () => {
  it('renders an uninterrupted theme-aware panel surface without a grid overlay', () => {
    const wrapper = mountDashboard()
    const panel = wrapper.get('svg > rect')

    expect(panel.classes()).toContain('fill-white')
    expect(panel.classes()).toContain('dark:fill-dark-900')
    expect(wrapper.find('pattern').exists()).toBe(false)
    expect(wrapper.find('rect[fill^="url("]').exists()).toBe(false)
  })

  it.each([
    { x: 206, label: 'Claude throughput', current: '428', delta: '+18%', color: '#ef6f2e' },
    { x: 460, label: 'OpenAI throughput', current: '312', delta: '+11%', color: '#10b981' },
    { x: 714, label: 'Gemini throughput', current: '186', delta: '+7%', color: '#38bdf8' }
  ])('renders the $label chart with its rate, growth and distinct curve', (card) => {
    const wrapper = mountDashboard()
    const chart = wrapper.get(`g[transform="translate(${card.x},404)"]`)

    expect(chart.text()).toContain(card.label)
    expect(chart.text()).toContain(`${card.current}/min`)
    expect(chart.text()).toContain(card.delta)
    expect(chart.get('text tspan').text()).toBe('/min')
    expect(chart.get('text tspan').element.namespaceURI).toBe('http://www.w3.org/2000/svg')
    expect(chart.get('path[fill="none"]').attributes('stroke')).toBe(card.color)
    const axisLabels = chart.findAll('text[y="170"]').map((label) => label.text())
    expect(axisLabels).toEqual(['09:00', '12:00', '16:00'])
  })

  it('plots the original OpenAI series including its dips and filled area', () => {
    const wrapper = mountDashboard()
    const chart = wrapper.get('g[transform="translate(460,404)"]')

    expect(chart.get('path[fill="none"]').attributes('d')).toBe(
      'M18.0 143.6 L47.1 152.0 L76.3 124.6 L105.4 137.2 L134.6 101.4 L163.7 85.6 L192.9 93.0 L222.0 74.0'
    )
    expect(chart.get('path[opacity="0.48"]').attributes('d')).toContain('L222.0 156 L18.0 156 Z')
    expect(chart.get('path[opacity="0.48"]').attributes('fill')).toBe('rgba(16,185,129,0.18)')
  })

  it('exposes all hourly samples with matching tooltip values and crosshairs', () => {
    const wrapper = mountDashboard()
    const chart = wrapper.get('g[transform="translate(206,404)"]')
    const points = chart.findAll('.sf-mini-chart-point')

    expect(points.map((point) => point.get('title').text())).toEqual([
      'Claude 09:00: 314/min',
      'Claude 10:00: 332/min',
      'Claude 11:00: 349/min',
      'Claude 12:00: 371/min',
      'Claude 13:00: 366/min',
      'Claude 14:00: 397/min',
      'Claude 15:00: 412/min',
      'Claude 16:00: 428/min'
    ])
    expect(wrapper.findAll('.sf-mini-chart-point')).toHaveLength(24)
    expect(points[0].get('.sf-mini-chart-tooltip').text()).toBe('09:00314/min')
    expect(points[7].get('.sf-mini-chart-tooltip').text()).toBe('16:00428/min')
    expect(points[0].get('.sf-mini-chart-crosshair').element.tagName).toBe('line')
    expect(points[0].get('.sf-mini-chart-hit').attributes('r')).toBe('8')
    expect(points[0].get('.sf-mini-chart-tooltip').attributes('transform')).toBe('translate(10,-42)')
    expect(points[7].get('.sf-mini-chart-tooltip').attributes('transform')).toBe('translate(-70,-42)')
  })

  it('makes chart samples keyboard-focusable for the existing focus-visible tooltip styles', () => {
    const wrapper = mountDashboard()
    const points = wrapper.findAll('.sf-mini-chart-point')

    expect(points).toHaveLength(24)
    for (const point of points) {
      expect(point.attributes('tabindex')).toBe('0')
    }
    const point = points[0].element as SVGElement
    point.focus()
    expect(document.activeElement).toBe(point)
  })
})
