<template>
  <div class="mt-5 grid grid-cols-2 gap-2">
    <div
      class="rounded-md p-3 bg-gray-50 dark:bg-dark-900/60 border border-gray-200 dark:border-dark-700"
    >
      <div
        class="flex items-center gap-1.5 text-[10px] font-mono uppercase tracking-wide text-gray-500 dark:text-dark-400"
      >
        <Icon :name="primaryIcon" size="xs" />
        <span>{{ primaryLabel }}</span>
      </div>
      <div class="mt-1.5 text-lg font-normal font-mono tabular-nums tracking-tight text-gray-900 dark:text-gray-100">
        {{ primaryValue }}<span class="text-xs font-normal text-gray-400 dark:text-dark-400 ml-0.5">{{ primaryUnit }}</span>
      </div>
    </div>
    <div
      class="rounded-md p-3 bg-gray-50 dark:bg-dark-900/60 border border-gray-200 dark:border-dark-700"
    >
      <div
        class="flex items-center gap-1.5 text-[10px] font-mono uppercase tracking-wide text-gray-500 dark:text-dark-400"
      >
        <Icon :name="secondaryIcon" size="xs" />
        <span>{{ secondaryLabel }}</span>
      </div>
      <div class="mt-1.5 text-lg font-normal font-mono tabular-nums tracking-tight text-gray-900 dark:text-gray-100">
        {{ secondaryValue }}<span class="text-xs font-normal text-gray-400 dark:text-dark-400 ml-0.5">{{ secondaryUnit }}</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import Icon from '@/components/icons/Icon.vue'

defineProps<{
  primaryLabel: string
  primaryValue: string
  primaryUnit: string
  primaryIcon: 'bolt' | 'globe' | 'clock' | 'link'
  secondaryLabel: string
  secondaryValue: string
  secondaryUnit: string
  secondaryIcon: 'bolt' | 'globe' | 'clock' | 'link'
}>()
</script>
