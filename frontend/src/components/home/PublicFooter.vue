<template>
  <footer class="relative z-10 mx-auto w-full max-w-[1920px] px-4 pb-10 lg:px-9">
    <div class="rounded-lg border border-gray-200 bg-white px-6 py-10 dark:border-dark-800 dark:bg-dark-900 lg:px-10 lg:py-12">
      <div class="grid grid-cols-2 gap-8 lg:grid-cols-12 lg:gap-6">
        <div class="col-span-2 flex flex-col gap-3 lg:col-span-4">
          <div class="flex items-center gap-2.5">
            <span class="block h-6 w-6 overflow-hidden rounded-sm">
              <img :src="resolvedSiteLogo || '/logo.svg'" alt="" class="h-full w-full object-contain" />
            </span>
            <span class="font-mono text-[13px] uppercase tracking-[0.14em] text-gray-900 dark:text-gray-100">{{ resolvedSiteName }}</span>
          </div>
          <p class="max-w-[36ch] font-mono text-[12px] uppercase tracking-[0.12em] text-gray-500 dark:text-dark-400">
            {{ t('home.footer2.tagline') }}
          </p>
          <a
            href="https://siliconbase.link"
            target="_blank"
            rel="noopener noreferrer"
            class="font-mono text-[12px] uppercase tracking-[0.12em] text-primary-500 transition-colors duration-200 hover:text-primary-400"
          >{{ t('home.footer2.domain') }}</a>
        </div>

        <div class="flex flex-col gap-3 lg:col-span-3 lg:col-start-7">
          <p class="eyebrow">{{ t('home.footer2.columns.resources') }}</p>
          <a
            v-if="resolvedDocUrl"
            :href="resolvedDocUrl"
            target="_blank"
            rel="noopener noreferrer"
            class="footer-link"
          >{{ t('home.footer2.links.docs') }}</a>
          <router-link to="/status" class="footer-link">{{ t('home.footer2.links.status') }}</router-link>
          <router-link to="/models" class="footer-link">{{ t('home.footer2.links.modelMarketplace') }}</router-link>
        </div>

        <div class="flex flex-col gap-3 lg:col-span-3">
          <p class="eyebrow">{{ t('home.footer2.columns.legal') }}</p>
          <router-link to="/legal/terms" class="footer-link">{{ t('home.footer2.links.serviceTerms') }}</router-link>
          <router-link to="/legal/usage-policy" class="footer-link">{{ t('home.footer2.links.usagePolicy') }}</router-link>
          <router-link to="/legal/supported-countries" class="footer-link">{{ t('home.footer2.links.supportedCountries') }}</router-link>
          <router-link to="/legal/service-specific-terms" class="footer-link">{{ t('home.footer2.links.serviceSpecificTerms') }}</router-link>
        </div>
      </div>

      <div class="mt-10 flex flex-col items-start justify-between gap-3 border-t border-dashed border-gray-200 pt-6 dark:border-dark-800 md:flex-row md:items-center">
        <p class="font-mono text-[11px] uppercase tracking-[0.12em] text-gray-500 dark:text-dark-400">
          {{ t('home.footer2.copyright') }}
        </p>
        <p class="font-mono text-[11px] uppercase tracking-[0.12em] text-gray-400 dark:text-dark-400">
          {{ t('home.footer2.domain') }}
        </p>
      </div>
    </div>
  </footer>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAppStore } from '@/stores'

const props = defineProps<{
  docUrl?: string
  siteLogo?: string
}>()

const { t } = useI18n()
const appStore = useAppStore()

const resolvedDocUrl = computed(() => props.docUrl || appStore.cachedPublicSettings?.doc_url || appStore.docUrl || '')
const resolvedSiteLogo = computed(() => props.siteLogo || appStore.cachedPublicSettings?.site_logo || appStore.siteLogo || '')
const resolvedSiteName = computed(() => appStore.cachedPublicSettings?.site_name?.trim() || appStore.siteName.trim() || 'Sub2API')
</script>

<style scoped>
.footer-link {
  @apply font-mono text-[13px] tracking-[-0.01em] text-gray-600 transition-colors duration-200 hover:text-primary-500 dark:text-dark-300;
}
</style>
