<template>
  <div>
    <button
      @click="openDrawer"
      class="relative flex h-9 w-9 items-center justify-center rounded-lg text-gray-600 transition-all hover:scale-105 hover:bg-gray-100 dark:text-gray-400 dark:hover:bg-dark-800"
      :class="{ 'text-blue-600 dark:text-blue-400': unreadCount > 0 }"
      :aria-label="t('announcements.title')"
    >
      <Icon name="bell" size="md" />
      <span
        v-if="unreadCount > 0"
        class="absolute right-1 top-1 flex h-2 w-2"
      >
        <span class="absolute inline-flex h-full w-full animate-ping rounded-full bg-red-500 opacity-75"></span>
        <span class="relative inline-flex h-2 w-2 rounded-full bg-red-500"></span>
      </span>
    </button>

    <Teleport to="body">
      <Transition name="drawer">
        <div
          v-if="isDrawerOpen"
          class="fixed inset-0 z-[100] bg-black/45 backdrop-blur-sm"
          @click="closeDrawer"
        >
          <aside
            data-testid="announcement-drawer"
            class="fixed right-0 top-0 flex h-full w-full max-w-[520px] flex-col border-l border-gray-200 bg-white shadow-2xl dark:border-dark-700 dark:bg-dark-900 sm:w-[520px]"
            @click.stop
          >
            <header class="border-b border-gray-200 px-5 py-4 dark:border-dark-700">
              <div class="flex items-start justify-between gap-4">
                <div class="min-w-0">
                  <div class="flex items-center gap-2">
                    <span class="flex h-9 w-9 items-center justify-center rounded-lg bg-gray-100 text-gray-700 dark:bg-dark-800 dark:text-gray-200">
                      <Icon name="bell" size="sm" />
                    </span>
                    <div>
                      <h2 class="text-base font-semibold text-gray-900 dark:text-white">
                        {{ t('announcements.title') }}
                      </h2>
                      <p class="mt-0.5 text-xs text-gray-500 dark:text-gray-400">
                        <span v-if="unreadCount > 0">{{ unreadCount }} {{ t('announcements.unread') }}</span>
                        <span v-else>{{ t('announcements.emptyDescription') }}</span>
                      </p>
                    </div>
                  </div>
                </div>

                <div class="flex items-center gap-2">
                  <button
                    v-if="unreadCount > 0"
                    @click="markAllAsRead"
                    :disabled="loading"
                    class="rounded-lg border border-gray-200 px-3 py-2 text-xs font-medium text-gray-700 transition-colors hover:bg-gray-50 disabled:cursor-not-allowed disabled:opacity-50 dark:border-dark-600 dark:text-gray-200 dark:hover:bg-dark-800"
                  >
                    {{ t('announcements.markAllRead') }}
                  </button>
                  <button
                    @click="closeDrawer"
                    class="flex h-9 w-9 items-center justify-center rounded-lg text-gray-500 transition-colors hover:bg-gray-100 hover:text-gray-700 dark:text-gray-400 dark:hover:bg-dark-800 dark:hover:text-gray-200"
                    :aria-label="t('common.close')"
                  >
                    <Icon name="x" size="sm" />
                  </button>
                </div>
              </div>
            </header>

            <div class="flex-1 overflow-y-auto px-5 py-5">
              <div v-if="loading" class="flex h-full min-h-[240px] items-center justify-center">
                <div class="h-10 w-10 animate-spin rounded-full border-2 border-gray-200 border-t-blue-600 dark:border-dark-700 dark:border-t-blue-400"></div>
              </div>

              <div
                v-else-if="announcements.length > 0"
                data-testid="announcement-timeline"
                class="relative space-y-0"
              >
                <article
                  v-for="item in announcements"
                  :key="item.id"
                  class="group relative border-l border-gray-200 pb-8 pl-7 last:border-l-transparent last:pb-0 dark:border-dark-700"
                >
                  <button
                    class="absolute -left-[7px] top-[5px] flex h-3.5 w-3.5 rounded-full border-2 transition-colors"
                    :class="item.read_at
                      ? 'border-blue-500 bg-blue-500'
                      : 'border-blue-500 bg-blue-500 shadow-[0_0_0_4px_rgba(59,130,246,0.18)]'"
                    :aria-label="item.read_at ? t('announcements.read') : t('announcements.unread')"
                    @click="markAsRead(item)"
                  ></button>

                  <div class="mb-2 flex min-h-[24px] items-center gap-2">
                    <time class="font-mono text-xs text-gray-500 dark:text-gray-500">
                      {{ formatRelativeWithDateTime(item.created_at) }}
                    </time>
                    <span
                      class="rounded-md px-2 py-1 text-[11px] font-medium"
                      :class="item.read_at
                        ? 'bg-gray-100 text-gray-500 dark:bg-dark-800 dark:text-gray-400'
                        : 'bg-blue-100 text-blue-700 dark:bg-blue-500/20 dark:text-blue-300'"
                    >
                      {{ item.read_at ? t('announcements.read') : t('announcements.unread') }}
                    </span>
                  </div>

                  <button
                    class="-ml-3 w-[calc(100%+0.75rem)] rounded-lg px-3 py-2 text-left transition-colors hover:bg-gray-50 dark:hover:bg-dark-800/70"
                    :class="{ 'bg-blue-50/40 dark:bg-blue-500/5': !item.read_at }"
                    @click="markAsRead(item)"
                  >
                    <h3 class="break-words text-base font-bold leading-7 text-gray-900 dark:text-white">
                      {{ item.title }}
                    </h3>

                    <div class="relative mt-2">
                      <div
                        :data-testid="`announcement-content-${item.id}`"
                        class="markdown-body overflow-hidden text-[15px] leading-7 text-gray-600 transition-[max-height] duration-300 ease-out dark:text-gray-300"
                        :class="shouldCollapseAnnouncement(item) && !isAnnouncementExpanded(item.id) ? 'max-h-40' : 'max-h-none'"
                        v-html="renderMarkdown(item.content)"
                      ></div>
                      <div
                        v-if="shouldCollapseAnnouncement(item) && !isAnnouncementExpanded(item.id)"
                        class="pointer-events-none absolute inset-x-0 bottom-0 h-12 bg-gradient-to-t from-white via-white/90 to-transparent dark:from-dark-900 dark:via-dark-900/90"
                      ></div>
                    </div>
                  </button>

                  <button
                    v-if="shouldCollapseAnnouncement(item)"
                    :data-testid="`announcement-toggle-${item.id}`"
                    class="mt-2 inline-flex items-center rounded-md px-2 py-1 text-xs font-medium text-blue-600 transition-colors hover:bg-blue-50 hover:text-blue-700 dark:text-blue-300 dark:hover:bg-blue-500/10"
                    @click.stop="toggleAnnouncementExpanded(item.id)"
                  >
                    {{ isAnnouncementExpanded(item.id) ? t('common.collapse') : t('common.expand') }}
                  </button>
                </article>
              </div>

              <div v-else class="flex h-full min-h-[260px] flex-col items-center justify-center text-center">
                <div class="flex h-14 w-14 items-center justify-center rounded-full bg-gray-100 text-gray-400 dark:bg-dark-800 dark:text-gray-500">
                  <Icon name="inbox" size="lg" />
                </div>
                <p class="mt-4 text-sm font-medium text-gray-900 dark:text-white">
                  {{ t('announcements.empty') }}
                </p>
                <p class="mt-1 max-w-[280px] text-xs text-gray-500 dark:text-gray-400">
                  {{ t('announcements.emptyDescription') }}
                </p>
              </div>
            </div>
          </aside>
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { storeToRefs } from 'pinia'
import { marked } from 'marked'
import DOMPurify from 'dompurify'
import { useAppStore } from '@/stores/app'
import { useAnnouncementStore } from '@/stores/announcements'
import { formatRelativeWithDateTime } from '@/utils/format'
import type { UserAnnouncement } from '@/types'
import Icon from '@/components/icons/Icon.vue'
import '@/styles/announcement-markdown.css'

const { t } = useI18n()
const appStore = useAppStore()
const announcementStore = useAnnouncementStore()

marked.setOptions({
  breaks: true,
  gfm: true,
})

const { announcements, loading } = storeToRefs(announcementStore)
const unreadCount = computed(() => announcementStore.unreadCount)

const isDrawerOpen = ref(false)
const expandedAnnouncementIds = ref<Set<number>>(new Set())
const COLLAPSIBLE_CONTENT_LENGTH = 360

function renderMarkdown(content: string): string {
  if (!content) return ''
  const html = marked.parse(content) as string
  return DOMPurify.sanitize(html)
}

function shouldCollapseAnnouncement(announcement: UserAnnouncement): boolean {
  return announcement.content.trim().length > COLLAPSIBLE_CONTENT_LENGTH
}

function isAnnouncementExpanded(id: number): boolean {
  return expandedAnnouncementIds.value.has(id)
}

function toggleAnnouncementExpanded(id: number) {
  const next = new Set(expandedAnnouncementIds.value)
  if (next.has(id)) {
    next.delete(id)
  } else {
    next.add(id)
  }
  expandedAnnouncementIds.value = next
}

function openDrawer() {
  isDrawerOpen.value = true
}

function closeDrawer() {
  isDrawerOpen.value = false
  expandedAnnouncementIds.value = new Set()
}

async function markAsRead(announcement: UserAnnouncement) {
  if (announcement.read_at) return

  try {
    await announcementStore.markAsRead(announcement.id)
  } catch (err: any) {
    appStore.showError(err?.message || t('common.unknownError'))
  }
}

async function markAllAsRead() {
  try {
    await announcementStore.markAllAsRead()
    appStore.showSuccess(t('announcements.allMarkedAsRead'))
  } catch (err: any) {
    appStore.showError(err?.message || t('common.unknownError'))
  }
}

function handleEscape(e: KeyboardEvent) {
  if (e.key === 'Escape' && isDrawerOpen.value) {
    closeDrawer()
  }
}

onMounted(() => {
  document.addEventListener('keydown', handleEscape)
})

onBeforeUnmount(() => {
  document.removeEventListener('keydown', handleEscape)
  document.body.style.overflow = ''
})

watch(
  [isDrawerOpen, () => announcementStore.currentPopup],
  ([drawer, popup]) => {
    document.body.style.overflow = (drawer || popup) ? 'hidden' : ''
  }
)
</script>

<style scoped>
.drawer-enter-active,
.drawer-leave-active {
  transition: opacity 0.36s ease;
}

.drawer-enter-active aside,
.drawer-leave-active aside {
  transition: transform 0.42s cubic-bezier(0.22, 1, 0.36, 1);
}

.drawer-enter-from,
.drawer-leave-to {
  opacity: 0;
}

.drawer-enter-from aside,
.drawer-leave-to aside {
  transform: translateX(100%);
}

.overflow-y-auto::-webkit-scrollbar {
  width: 8px;
}

.overflow-y-auto::-webkit-scrollbar-track {
  background: transparent;
}

.overflow-y-auto::-webkit-scrollbar-thumb {
  background: #cbd5e1;
  border-radius: 4px;
}

.dark .overflow-y-auto::-webkit-scrollbar-thumb {
  background: #374151;
}
</style>
