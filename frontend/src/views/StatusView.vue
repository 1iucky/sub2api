<template>
  <div class="min-h-screen bg-gray-50 text-gray-900 dark:bg-dark-950 dark:text-gray-100">
    <PublicTopNav />
    <main class="mx-auto max-w-[1600px] px-4 py-8 lg:px-8">
      <ChannelStatusV1View v-if="isV1" is-public />
      <ChannelStatusV2View v-else is-public />
    </main>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useAppStore, useAuthStore } from '@/stores'
import PublicTopNav from '@/components/home/PublicTopNav.vue'
import ChannelStatusV1View from '@/views/user/ChannelStatusV1View.vue'
import ChannelStatusV2View from '@/views/user/ChannelStatusV2View.vue'
import { useTheme } from '@/composables/useTheme'
import { isChannelMonitorV1Mode } from '@/utils/featureFlags'

// Standalone public /status page: renders the same channel-status view as the
// console (v1 or v2 per channel_monitor_mode) without AppLayout, fed by the
// unauthenticated public APIs.
const appStore = useAppStore()
const authStore = useAuthStore()
const { syncThemeFromDocument } = useTheme()

const isV1 = computed(() => isChannelMonitorV1Mode())

onMounted(() => {
  syncThemeFromDocument()
  authStore.checkAuth()
  if (!appStore.publicSettingsLoaded) {
    appStore.fetchPublicSettings()
  }
})
</script>
