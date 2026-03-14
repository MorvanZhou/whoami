<script setup lang="ts">
import { marked } from 'marked'

const props = defineProps<{
  open: boolean
}>()

const emit = defineEmits<{
  close: []
}>()

const { t } = useI18n()
const { fetchData } = useRequest()
const { trackProfileSaved } = useTracking()

const profileContent = ref('')
const originalContent = ref('')
const loading = ref(false)
const saving = ref(false)
const isEditing = ref(false)
const saveMessage = ref('')
const saveError = ref('')

const hasChanges = computed(() => profileContent.value !== originalContent.value)
const charCount = computed(() => profileContent.value.length)
const isOverLimit = computed(() => charCount.value > 5000)

const renderedHtml = computed(() => {
  if (!profileContent.value) return ''
  return marked.parse(profileContent.value, { async: false }) as string
})

const fetchProfile = async () => {
  loading.value = true
  saveMessage.value = ''
  saveError.value = ''
  try {
    const result = await fetchData<{ content: string | null }>('/web/profile')
    profileContent.value = result.content || ''
    originalContent.value = profileContent.value
  } catch {
    saveError.value = t('profile.fetchError')
  } finally {
    loading.value = false
  }
}

const saveProfile = async () => {
  if (isOverLimit.value) return
  saving.value = true
  saveMessage.value = ''
  saveError.value = ''
  try {
    await fetchData('/web/profile', {
      method: 'POST',
      body: { content: profileContent.value },
    })
    originalContent.value = profileContent.value
    trackProfileSaved()
    saveMessage.value = t('profile.saveSuccess')
    isEditing.value = false
    setTimeout(() => { saveMessage.value = '' }, 3000)
  } catch {
    saveError.value = t('profile.saveError')
  } finally {
    saving.value = false
  }
}

const startEditing = () => {
  isEditing.value = true
  saveMessage.value = ''
  saveError.value = ''
}

const cancelEditing = () => {
  profileContent.value = originalContent.value
  isEditing.value = false
  saveMessage.value = ''
  saveError.value = ''
}

const onBackdrop = (e: MouseEvent) => {
  if (e.target === e.currentTarget) {
    handleClose()
  }
}

const handleClose = () => {
  if (isEditing.value && hasChanges.value) {
    if (!confirm(t('profile.unsavedWarning'))) return
  }
  isEditing.value = false
  saveMessage.value = ''
  saveError.value = ''
  emit('close')
}

watch(() => props.open, (val) => {
  if (val) {
    fetchProfile()
  }
})
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
        class="fixed inset-0 z-[100] flex items-center justify-center px-4 bg-th-overlay/60 backdrop-blur-sm"
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
            class="w-full max-w-2xl max-h-[85vh] rounded-2xl bg-th-card border border-th-text/[0.08] shadow-2xl flex flex-col overflow-hidden"
          >
            <!-- Header -->
            <div class="flex items-center justify-between px-6 py-4 border-b border-th-text/[var(--color-card-border-opacity)]">
              <div class="flex items-center gap-3">
                <div class="w-9 h-9 rounded-lg bg-gradient-to-br from-accent/15 to-accent/5 flex items-center justify-center ring-1 ring-accent/10">
                  <svg class="w-4.5 h-4.5 text-accent" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                  </svg>
                </div>
                <div>
                  <h2 class="text-lg font-semibold text-th-text">{{ t('profile.title') }}</h2>
                  <p v-if="!isEditing && profileContent" class="text-xs text-th-text-t mt-0.5">
                    {{ charCount }} {{ t('profile.characters') }}
                  </p>
                </div>
              </div>
              <button
                class="w-8 h-8 rounded-lg flex items-center justify-center text-th-text-t hover:text-th-text-s hover:bg-th-bg-t transition-colors"
                @click="handleClose"
              >
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>

            <!-- Body -->
            <div class="flex-1 overflow-y-auto px-6 py-5">
              <!-- Loading -->
              <div v-if="loading" class="flex items-center justify-center py-16">
                <div class="w-5 h-5 border-2 border-accent border-t-transparent rounded-full animate-spin" />
              </div>

              <!-- Content -->
              <template v-else>
                <!-- View mode -->
                <div v-if="!isEditing">
                  <div v-if="profileContent" class="relative">
                    <div
                      class="prose prose-sm max-w-none dark:prose-invert prose-headings:text-th-text prose-headings:font-semibold prose-h1:text-lg prose-h1:border-b prose-h1:border-th-text/[0.06] prose-h1:pb-2 prose-h2:text-base prose-h2:mt-6 prose-h2:mb-2 prose-p:text-th-text-s prose-p:leading-relaxed prose-strong:text-accent-glow prose-strong:font-semibold prose-li:text-th-text-s prose-li:marker:text-accent/50 prose-a:text-accent prose-a:no-underline hover:prose-a:underline prose-code:text-accent-glow prose-code:bg-accent/10 prose-code:px-1.5 prose-code:py-0.5 prose-code:rounded prose-code:text-xs prose-code:font-mono prose-code:before:content-none prose-code:after:content-none prose-hr:border-th-text/[0.06] selection:bg-accent/20"
                      v-html="renderedHtml"
                    />
                  </div>

                  <!-- Empty state -->
                  <div v-else class="flex flex-col items-center justify-center py-14">
                    <div class="w-14 h-14 rounded-2xl bg-th-bg-t flex items-center justify-center mb-5 ring-1 ring-th-text/[0.06]">
                      <svg class="w-6 h-6 text-th-text-m" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                      </svg>
                    </div>
                    <p class="text-base text-th-text-s font-medium">{{ t('profile.empty') }}</p>
                    <p class="text-sm text-th-text-m mt-1">{{ t('profile.emptyHint') }}</p>
                  </div>
                </div>

                <!-- Edit mode -->
                <div v-else>
                  <textarea
                    v-model="profileContent"
                    class="w-full h-80 px-4 py-3.5 rounded-xl bg-th-input border border-th-text/[var(--color-input-border-opacity)] text-base text-th-text font-mono leading-relaxed placeholder-th-text-m focus:outline-none focus:border-accent/30 focus:ring-1 focus:ring-accent/15 transition-all resize-y"
                    :placeholder="t('profile.editPlaceholder')"
                    spellcheck="false"
                  />
                  <div class="flex items-center justify-between mt-2 text-xs">
                    <span :class="isOverLimit ? 'text-danger' : 'text-th-text-m'">
                      {{ charCount }} / 5,000 {{ t('profile.characters') }}
                    </span>
                    <span v-if="isOverLimit" class="text-danger font-medium">
                      {{ t('profile.overLimit') }}
                    </span>
                  </div>
                </div>
              </template>

              <!-- Messages -->
              <Transition
                enter-active-class="transition duration-200 ease-out"
                enter-from-class="opacity-0 -translate-y-1"
                enter-to-class="opacity-100 translate-y-0"
              >
                <div v-if="saveMessage" class="mt-4 px-4 py-2.5 rounded-lg bg-success/8 border border-success/15 text-sm text-success font-medium">
                  {{ saveMessage }}
                </div>
              </Transition>
              <div v-if="saveError" class="mt-4 px-4 py-2.5 rounded-lg bg-danger/8 border border-danger/15 text-sm text-danger font-medium">
                {{ saveError }}
              </div>
            </div>

            <!-- Footer -->
            <div class="flex items-center justify-end gap-3 px-6 py-4 border-t border-th-text/[var(--color-card-border-opacity)] bg-th-bg/30">
              <template v-if="!isEditing">
                <button
                  class="px-4 py-2 rounded-lg text-sm font-medium text-th-text-t hover:text-th-text-s hover:bg-th-bg-t transition-colors"
                  @click="handleClose"
                >
                  {{ t('profile.close') }}
                </button>
                <button
                  class="inline-flex items-center gap-1.5 px-5 py-2 rounded-lg text-sm font-semibold text-white bg-accent hover:bg-accent-light transition-colors"
                  @click="startEditing"
                >
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                  </svg>
                  {{ t('profile.edit') }}
                </button>
              </template>
              <template v-else>
                <button
                  class="px-4 py-2 rounded-lg text-sm font-medium text-th-text-t hover:text-th-text-s hover:bg-th-bg-t transition-colors"
                  @click="cancelEditing"
                >
                  {{ t('common.cancel') }}
                </button>
                <button
                  :disabled="!hasChanges || isOverLimit || saving"
                  class="inline-flex items-center gap-1.5 px-5 py-2 rounded-lg text-sm font-semibold text-white bg-accent hover:bg-accent-light transition-colors disabled:opacity-30 disabled:cursor-not-allowed"
                  @click="saveProfile"
                >
                  <div v-if="saving" class="w-4 h-4 border-2 border-white/60 border-t-transparent rounded-full animate-spin" />
                  <svg v-else class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                  </svg>
                  {{ t('profile.save') }}
                </button>
              </template>
            </div>
          </div>
        </Transition>
      </div>
    </Transition>
  </Teleport>
</template>
