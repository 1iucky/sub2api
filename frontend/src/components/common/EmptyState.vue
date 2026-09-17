<template>
  <div class="empty-state">
    <!-- Icon -->
    <div
      class="mb-5 flex h-20 w-20 items-center justify-center rounded-2xl bg-gray-100 dark:bg-dark-800"
    >
      <slot name="icon">
        <component v-if="icon" :is="icon" class="empty-state-icon h-10 w-10" aria-hidden="true" />
        <!-- Factory-style empty-state glyph: blueprint grid, empty node slot, and one signal dot. -->
        <svg v-else class="empty-state-icon h-10 w-10" viewBox="0 0 40 40" fill="none" aria-hidden="true">
          <g stroke="currentColor" stroke-width="0.6" opacity="0.35">
            <path d="M8 0v40M16 0v40M24 0v40M32 0v40" />
            <path d="M0 8h40M0 16h40M0 24h40M0 32h40" />
          </g>
          <rect
            x="11.5"
            y="11.5"
            width="17"
            height="17"
            rx="2"
            stroke="currentColor"
            stroke-width="1.1"
            stroke-dasharray="2.5 2.5"
          />
          <path
            d="M20 16.5v7M16.5 20h7"
            stroke="currentColor"
            stroke-width="1"
            stroke-linecap="round"
            opacity="0.7"
          />
          <circle cx="20" cy="20" r="1.8" fill="#ef6f2e" />
        </svg>
      </slot>
    </div>

    <!-- Title -->
    <h3 class="empty-state-title">
      {{ displayTitle }}
    </h3>

    <!-- Description -->
    <p class="empty-state-description">
      {{ description }}
    </p>

    <!-- Action -->
    <div v-if="actionText || $slots.action" class="mt-6">
      <slot name="action">
        <component
          :is="actionTo ? 'RouterLink' : 'button'"
          v-if="actionText"
          :to="actionTo"
          @click="!actionTo && $emit('action')"
          class="btn btn-primary"
        >
          <Icon v-if="actionIcon" name="plus" size="md" class="mr-2" />
          {{ actionText }}
        </component>
      </slot>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import type { Component } from 'vue'
import Icon from '@/components/icons/Icon.vue'

const { t } = useI18n()

interface Props {
  icon?: Component | string
  title?: string
  description?: string
  actionText?: string
  actionTo?: string | object
  actionIcon?: boolean
  message?: string
}

const props = withDefaults(defineProps<Props>(), {
  description: '',
  actionIcon: true
})

const displayTitle = computed(() => props.title || t('common.noData'))

defineEmits(['action'])
</script>
