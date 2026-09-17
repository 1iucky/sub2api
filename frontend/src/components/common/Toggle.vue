<template>
  <button
    type="button"
    @click="toggle"
    class="relative inline-flex h-6 w-11 flex-shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-150 ease-in-out focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-500"
    :class="[modelValue ? 'bg-primary-500' : 'bg-gray-300 dark:bg-dark-700']"
    role="switch"
    :aria-checked="modelValue"
  >
    <span
      class="pointer-events-none inline-block h-5 w-5 transform rounded-full bg-white transition-transform duration-150 ease-in-out dark:bg-gray-100"
      :class="[modelValue ? 'translate-x-5' : 'translate-x-0']"
    />
  </button>
</template>

<script setup lang="ts">
const props = defineProps<{
  modelValue: boolean
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: boolean): void
}>()

function toggle() {
  emit('update:modelValue', !props.modelValue)
}
</script>
