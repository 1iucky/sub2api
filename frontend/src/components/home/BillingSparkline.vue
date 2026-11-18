<template>
  <!-- CSS-only sparkline + key/usage tiles: quotas & metered billing. -->
  <div class="flex h-24 w-full flex-col gap-2" aria-hidden="true">
    <!-- tiles row -->
    <div class="flex gap-2">
      <div class="flex-1 rounded-sm border border-gray-200 px-2.5 py-1.5 dark:border-dark-700">
        <div class="font-mono text-[8px] uppercase tracking-[0.12em] text-gray-400 dark:text-dark-400">keys</div>
        <div class="font-mono text-sm text-gray-900 dark:text-gray-100">348</div>
      </div>
      <div class="flex-1 rounded-sm border border-gray-200 px-2.5 py-1.5 dark:border-dark-700">
        <div class="font-mono text-[8px] uppercase tracking-[0.12em] text-gray-400 dark:text-dark-400">quota</div>
        <div class="font-mono text-sm text-primary-500">72%</div>
      </div>
    </div>
    <!-- sparkline -->
    <svg viewBox="0 0 200 32" class="w-full" preserveAspectRatio="none">
      <defs>
        <linearGradient :id="gradId" x1="0" y1="0" x2="0" y2="1">
          <stop offset="0%" stop-color="#ef6f2e" stop-opacity="0.35" />
          <stop offset="100%" stop-color="#ef6f2e" stop-opacity="0" />
        </linearGradient>
      </defs>
      <path d="M0 26 L20 22 L40 24 L60 18 L80 20 L100 12 L120 16 L140 8 L160 12 L180 4 L200 8" fill="none" stroke="#ef6f2e" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" />
      <path d="M0 26 L20 22 L40 24 L60 18 L80 20 L100 12 L120 16 L140 8 L160 12 L180 4 L200 8 L200 32 L0 32 Z" :fill="`url(#${gradId})`" />
    </svg>
  </div>
</template>

<script setup lang="ts">
// Unique gradient id so multiple instances on one page do not collide.
import { useId } from 'vue'
const gradId = `sf-billing-grad-${useId()}`
</script>
