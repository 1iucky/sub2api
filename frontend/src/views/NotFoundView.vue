<template>
  <div
    class="bg-blueprint relative flex min-h-screen items-center justify-center overflow-hidden bg-gray-50 px-4 dark:bg-dark-950"
  >
    <div class="relative z-10 w-full max-w-xl text-center">
      <!-- Mono eyebrow -->
      <div class="eyebrow mb-6">Error / 404</div>

      <!-- Big 404 code -->
      <div class="mb-8 select-none">
        <span
          class="block font-normal leading-none tracking-[-0.06em] text-gray-200 dark:text-dark-800 text-[10rem] sm:text-[14rem]"
          >404</span
        >
      </div>

      <!-- Title -->
      <h1
        class="mb-3 text-2xl font-normal tracking-tight text-gray-900 dark:text-white sm:text-3xl"
      >
        {{ t('errors.pageNotFound') }}
      </h1>

      <!-- Description -->
      <p class="mx-auto mb-8 max-w-md text-sm text-gray-500 dark:text-dark-400">
        {{ t('notFound.description') }}
      </p>

      <!-- Action buttons -->
      <div class="flex flex-col justify-center gap-3 sm:flex-row">
        <button @click="goBack" class="btn btn-secondary">
          <Icon name="arrowLeft" size="md" class="mr-2" />
          {{ t('notFound.goBack') }}
        </button>
        <router-link to="/dashboard" class="btn btn-primary">
          <Icon name="home" size="md" class="mr-2" />
          {{ t('notFound.goDashboard') }}
        </router-link>
      </div>

      <!-- Brand wordmark -->
      <p class="mt-12 font-mono text-xs uppercase tracking-[0.14em] text-gray-400 dark:text-dark-500">
        SiliconBase &middot; siliconbase.link
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import Icon from '@/components/icons/Icon.vue'

const { t } = useI18n()
const router = useRouter()

function goBack(): void {
  router.back()
}
</script>
