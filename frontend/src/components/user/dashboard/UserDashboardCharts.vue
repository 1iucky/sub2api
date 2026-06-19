<template>
  <div class="space-y-6">
    <!-- Date Range Filter -->
    <div class="card p-4">
      <div class="flex flex-wrap items-center gap-4">
        <div class="flex items-center gap-2">
          <span class="font-mono text-xs uppercase tracking-wide text-gray-500 dark:text-dark-400">{{ t('dashboard.timeRange') }}:</span>
          <DateRangePicker :start-date="startDate" :end-date="endDate" @update:startDate="$emit('update:startDate', $event)" @update:endDate="$emit('update:endDate', $event)" @change="$emit('dateRangeChange', $event)" />
        </div>
        <button @click="$emit('refresh')" :disabled="loading" class="btn btn-secondary">
          {{ t('common.refresh') }}
        </button>
        <div class="ml-auto flex items-center gap-2">
          <span class="font-mono text-xs uppercase tracking-wide text-gray-500 dark:text-dark-400">{{ t('dashboard.granularity') }}:</span>
          <div class="w-28">
            <Select :model-value="granularity" :options="[{value:'day', label:t('dashboard.day')}, {value:'hour', label:t('dashboard.hour')}]" @update:model-value="$emit('update:granularity', $event)" @change="$emit('granularityChange')" />
          </div>
        </div>
      </div>
    </div>

    <!-- Charts Grid -->
    <div class="grid grid-cols-1 gap-6 lg:grid-cols-2">
      <!-- Model Distribution Chart -->
      <div class="card relative overflow-hidden p-4">
        <div v-if="loading" class="absolute inset-0 z-10 flex items-center justify-center bg-white/60 backdrop-blur-sm dark:bg-dark-900/60">
          <LoadingSpinner size="md" />
        </div>
        <h3 class="eyebrow mb-4 !text-sm !normal-case !tracking-wide text-gray-700 dark:text-gray-300">{{ t('dashboard.modelDistribution') }}</h3>
        <div class="flex items-center gap-6">
          <div class="h-48 w-48">
            <Doughnut v-if="modelData" :data="modelData" :options="doughnutOptions" />
            <div v-else class="flex h-full items-center justify-center text-sm text-gray-500 dark:text-gray-400">{{ t('dashboard.noDataAvailable') }}</div>
          </div>
          <div class="max-h-48 flex-1 overflow-y-auto">
            <table class="w-full text-xs">
              <thead>
                <tr class="border-b border-dashed border-gray-300 font-mono uppercase tracking-wide text-gray-500 dark:border-dark-600 dark:text-dark-400">
                  <th class="pb-2 text-left font-medium">{{ t('dashboard.model') }}</th>
                  <th class="pb-2 text-right font-medium">{{ t('dashboard.requests') }}</th>
                  <th class="pb-2 text-right font-medium">{{ t('dashboard.tokens') }}</th>
                  <th class="pb-2 text-right font-medium">{{ t('dashboard.actual') }}</th>
                  <th class="pb-2 text-right font-medium">{{ t('dashboard.standard') }}</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="model in models" :key="model.model" class="border-b border-dashed border-gray-200 dark:border-dark-700">
                  <td class="max-w-[100px] truncate py-1.5 font-normal text-gray-900 dark:text-white" :title="model.model">{{ model.model }}</td>
                  <td class="py-1.5 text-right font-mono tabular-nums text-gray-600 dark:text-gray-400">{{ formatNumber(model.requests) }}</td>
                  <td class="py-1.5 text-right font-mono tabular-nums text-gray-600 dark:text-gray-400">{{ formatTokens(model.total_tokens) }}</td>
                  <td class="py-1.5 text-right font-mono tabular-nums text-primary-600 dark:text-primary-400">${{ formatCost(model.actual_cost) }}</td>
                  <td class="py-1.5 text-right font-mono tabular-nums text-gray-400 dark:text-gray-500">${{ formatCost(model.cost) }}</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>

      <!-- Token Usage Trend Chart -->
      <TokenUsageTrend :trend-data="trend" :loading="loading" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import LoadingSpinner from '@/components/common/LoadingSpinner.vue'
import DateRangePicker from '@/components/common/DateRangePicker.vue'
import Select from '@/components/common/Select.vue'
import { Doughnut } from 'vue-chartjs'
import TokenUsageTrend from '@/components/charts/TokenUsageTrend.vue'
import type { TrendDataPoint, ModelStat } from '@/types'
import { formatCostFixed as formatCost, formatNumberLocaleString as formatNumber, formatTokensK as formatTokens } from '@/utils/format'
import { Chart as ChartJS, CategoryScale, LinearScale, PointElement, LineElement, ArcElement, Title, Tooltip, Legend, Filler } from 'chart.js'
ChartJS.register(CategoryScale, LinearScale, PointElement, LineElement, ArcElement, Title, Tooltip, Legend, Filler)

const props = defineProps<{ loading: boolean, startDate: string, endDate: string, granularity: string, trend: TrendDataPoint[], models: ModelStat[] }>()
defineEmits(['update:startDate', 'update:endDate', 'update:granularity', 'dateRangeChange', 'granularityChange', 'refresh'])
const { t } = useI18n()

// Factory chart palette: vermillion primary + warm grays + status colors.
// No cool blue/teal/purple.
const CHART_PALETTE = [
  '#ef6f2e', // primary vermillion
  '#a09486', // warm gray-400
  '#7d7264', // warm gray-500
  '#4a4339', // warm gray-700
  '#10b981', // emerald (status)
  '#f59e0b', // amber (status)
  '#ef4444', // red (status)
  '#bfb3a5'  // warm gray-300
]

const modelData = computed(() => !props.models?.length ? null : {
  labels: props.models.map((m: ModelStat) => m.model),
  datasets: [{
    data: props.models.map((m: ModelStat) => m.total_tokens),
    backgroundColor: props.models.map((_, i: number) => CHART_PALETTE[i % CHART_PALETTE.length]),
    borderColor: '#0a0908',
    borderWidth: 1
  }]
})

const doughnutOptions = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: { display: false },
    tooltip: {
      backgroundColor: '#0a0908',
      titleColor: '#fafafa',
      bodyColor: '#e8e3dd',
      borderColor: 'rgba(125,114,100,0.25)',
      borderWidth: 1,
      padding: 10,
      callbacks: {
        label: (context: any) => `${context.label}: ${formatTokens(context.parsed)} tokens`
      }
    }
  }
}
</script>
