<script setup lang="ts">
const props = withDefaults(defineProps<{
  open: boolean
  title?: string
  message?: string
  confirmText?: string
  cancelText?: string
  variant?: 'danger' | 'warning' | 'info'
}>(), {
  title: '',
  message: '',
  confirmText: 'Confirm',
  cancelText: 'Cancel',
  variant: 'danger',
})

const emit = defineEmits<{
  confirm: []
  cancel: []
}>()

const variantClasses = computed(() => {
  switch (props.variant) {
    case 'danger':
      return {
        icon: 'text-danger',
        iconBg: 'bg-danger/10',
        button: 'bg-danger hover:bg-danger/80',
      }
    case 'warning':
      return {
        icon: 'text-warning',
        iconBg: 'bg-warning/10',
        button: 'bg-warning hover:bg-warning/80 text-black',
      }
    case 'info':
      return {
        icon: 'text-accent',
        iconBg: 'bg-accent/10',
        button: 'bg-accent hover:bg-accent-light',
      }
  }
})

const onBackdrop = (e: MouseEvent) => {
  if (e.target === e.currentTarget) {
    emit('cancel')
  }
}
</script>

<template>
  <Teleport to="body">
    <Transition
      enter-active-class="transition duration-200 ease-out"
      enter-from-class="opacity-0"
      enter-to-class="opacity-100"
      leave-active-class="transition duration-150 ease-in"
      leave-from-class="opacity-100"
      leave-to-class="opacity-0"
    >
      <div
        v-if="open"
        class="fixed inset-0 z-[100] flex items-center justify-center px-4 bg-black/60 backdrop-blur-sm"
        @mousedown="onBackdrop"
      >
        <Transition
          enter-active-class="transition duration-200 ease-out"
          enter-from-class="opacity-0 scale-95 translate-y-2"
          enter-to-class="opacity-100 scale-100 translate-y-0"
          leave-active-class="transition duration-150 ease-in"
          leave-from-class="opacity-100 scale-100 translate-y-0"
          leave-to-class="opacity-0 scale-95 translate-y-2"
        >
          <div
            v-if="open"
            class="w-full max-w-sm rounded-2xl bg-surface-800 border border-white/10 shadow-2xl overflow-hidden"
          >
            <div class="p-6">
              <!-- Icon + Title -->
              <div class="flex items-start gap-4">
                <div :class="['shrink-0 w-10 h-10 rounded-xl flex items-center justify-center', variantClasses?.iconBg]">
                  <!-- Danger: trash icon -->
                  <svg v-if="variant === 'danger'" :class="['w-5 h-5', variantClasses?.icon]" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                  </svg>
                  <!-- Warning: exclamation -->
                  <svg v-else-if="variant === 'warning'" :class="['w-5 h-5', variantClasses?.icon]" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4.5c-.77-.833-2.694-.833-3.464 0L3.34 16.5c-.77.833.192 2.5 1.732 2.5z" />
                  </svg>
                  <!-- Info: info circle -->
                  <svg v-else :class="['w-5 h-5', variantClasses?.icon]" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                </div>
                <div class="flex-1 min-w-0">
                  <h3 class="text-base font-semibold text-white">{{ title }}</h3>
                  <p class="text-sm text-gray-400 mt-1">{{ message }}</p>
                </div>
              </div>
            </div>

            <!-- Actions -->
            <div class="flex items-center justify-end gap-3 px-6 py-4 bg-surface-900/50 border-t border-white/5">
              <button
                class="px-4 py-2 rounded-lg text-sm font-medium text-gray-300 hover:text-white hover:bg-white/5 transition-colors"
                @click="emit('cancel')"
              >
                {{ cancelText }}
              </button>
              <button
                :class="['px-4 py-2 rounded-lg text-sm font-semibold text-white transition-colors', variantClasses?.button]"
                @click="emit('confirm')"
              >
                {{ confirmText }}
              </button>
            </div>
          </div>
        </Transition>
      </div>
    </Transition>
  </Teleport>
</template>
