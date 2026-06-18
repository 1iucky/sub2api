<template>
  <!-- CSS-only bar chart motif: multi-upstream coverage. -->
  <div class="flex h-24 w-full items-end gap-2" aria-hidden="true">
    <div
      v-for="(bar, idx) in bars"
      :key="idx"
      class="group/bar relative flex flex-1 flex-col items-center justify-end gap-1.5"
    >
      <div
        class="w-full rounded-sm transition-colors duration-150"
        :class="bar.active
          ? 'bg-primary-500'
          : 'bg-gray-200 dark:bg-dark-700 group-hover/bar:bg-primary-500/40'"
        :style="{ height: `${bar.height}px` }"
      ></div>
      <span
        class="font-mono text-[9px] uppercase tracking-[0.1em] text-gray-400 dark:text-dark-400"
      >{{ bar.label }}</span>
    </div>
  </div>
</template>

<script setup lang="ts">
// Static motif — no live data; the values are tuned to read as "many upstreams covered".
const bars = [
  { label: 'CL', height: 72, active: true },
  { label: 'GP', height: 58, active: true },
  { label: 'GE', height: 48, active: true },
  { label: 'BR', height: 36, active: true },
  { label: 'AG', height: 28, active: true },
  { label: '+', height: 16, active: false }
]
</script>
