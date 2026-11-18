<template>
  <!-- CSS-only routing diagram: client -> gateway pool -> upstreams. -->
  <div class="relative flex h-24 w-full items-center" aria-hidden="true">
    <!-- source -->
    <div class="flex flex-col items-center gap-1">
      <div
        class="flex h-9 w-9 items-center justify-center rounded-sm border border-gray-300 bg-white font-mono text-[10px] uppercase tracking-wide text-gray-500 dark:border-dark-700 dark:bg-dark-900 dark:text-dark-300"
      >client</div>
    </div>

    <!-- dashed line to gateway -->
    <div class="h-px flex-1 border-t border-dashed border-gray-300 dark:border-dark-700"></div>

    <!-- gateway -->
    <div class="flex flex-col items-center gap-1">
      <div
        class="flex h-9 w-9 items-center justify-center rounded-sm border border-primary-500/40 bg-primary-500/10 font-mono text-[10px] font-semibold uppercase tracking-wide text-primary-500"
      >gw</div>
      <span class="font-mono text-[8px] uppercase tracking-wide text-gray-400 dark:text-dark-400">pool</span>
    </div>

    <!-- dashed line to upstreams -->
    <div class="h-px flex-1 border-t border-dashed border-gray-300 dark:border-dark-700"></div>

    <!-- upstreams (stacked) -->
    <div class="flex flex-col gap-1">
      <div
        v-for="(u, i) in upstreams"
        :key="i"
        class="flex items-center gap-1.5 font-mono text-[9px] uppercase tracking-wide text-gray-500 dark:text-dark-400"
      >
        <span class="inline-block h-1.5 w-1.5 rounded-full" :class="u.live ? 'bg-primary-500' : 'bg-gray-300 dark:bg-dark-600'"></span>
        {{ u.name }}
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const upstreams = [
  { name: 'claude', live: true },
  { name: 'openai', live: true },
  { name: 'gemini', live: false },
  { name: 'bedrock', live: true }
]
</script>
