<script setup lang="ts">
definePageMeta({
  middleware: 'auth',
})

const { t } = useI18n()
const { fetchData } = useRequest()

useSeoMeta({
  title: () => t('seo.dashboard.title'),
  description: () => t('seo.dashboard.description'),
  robots: 'noindex, nofollow',
})

interface AccessKeyItem {
  id: string
  key_prefix: string
  key_suffix: string
  label: string | null
  created_at: string
  last_used: string | null
}

const keys = ref<AccessKeyItem[]>([])
const storeUrl = ref<string | null>(null)
const showPromptGuide = ref(false)
const setupCompleted = ref(false)
const loading = ref(false)
const showCreate = ref(false)
const newLabel = ref('')
const labelError = ref(false)
const deleteTarget = ref<string | null>(null)
const isInitialLoad = ref(true)
const showProfile = ref(false)
const countdown = ref(0)
const pendingKeyId = ref<string | null>(null)

let countdownTimer: ReturnType<typeof setInterval> | null = null
let pollTimer: ReturnType<typeof setInterval> | null = null

const formattedCountdown = computed(() => {
  const m = Math.floor(countdown.value / 60)
  const s = countdown.value % 60
  return `${m}:${s.toString().padStart(2, '0')}`
})

const startCountdownAndPolling = (expiresIn: number, keyId: string) => {
  stopCountdownAndPolling()
  countdown.value = expiresIn
  pendingKeyId.value = keyId
  setupCompleted.value = false

  countdownTimer = setInterval(() => {
    countdown.value--
    if (countdown.value <= 0) {
      stopCountdownAndPolling()
    }
  }, 1000)

  pollTimer = setInterval(async () => {
    if (!pendingKeyId.value) return
    try {
      const res = await fetchData<{ configured: boolean }>(`/keys/${pendingKeyId.value}/status`)
      if (res.configured) {
        stopCountdownAndPolling()
        showPromptGuide.value = false
        setupCompleted.value = true
      }
    } catch {
      // ignore polling errors
    }
  }, 3000)
}

const stopCountdownAndPolling = () => {
  if (countdownTimer) { clearInterval(countdownTimer); countdownTimer = null }
  if (pollTimer) { clearInterval(pollTimer); pollTimer = null }
  pendingKeyId.value = null
}

const fetchKeys = async () => {
  loading.value = true
  try {
    keys.value = await fetchData<AccessKeyItem[]>('/keys')
    if (isInitialLoad.value && keys.value.length === 0) {
      await autoCreateKey()
    }
    isInitialLoad.value = false
  } catch {
    // handle error
  } finally {
    loading.value = false
  }
}

const autoCreateKey = async () => {
  try {
    const result = await fetchData<{ id: string; store_url: string; expires_in: number }>('/keys', {
      method: 'POST',
      body: { label: 'My First Agent' },
    })
    showPromptGuide.value = true
    storeUrl.value = result.store_url
    startCountdownAndPolling(result.expires_in, result.id)
    keys.value = await fetchData<AccessKeyItem[]>('/keys')
  } catch {
    // handle error
  }
}

const createKey = async () => {
  if (!newLabel.value.trim()) {
    labelError.value = true
    return
  }
  labelError.value = false
  try {
    const result = await fetchData<{ id: string; store_url: string; expires_in: number }>('/keys', {
      method: 'POST',
      body: { label: newLabel.value.trim() },
    })
    showPromptGuide.value = true
    storeUrl.value = result.store_url
    startCountdownAndPolling(result.expires_in, result.id)
    showCreate.value = false
    newLabel.value = ''
    await fetchKeys()
  } catch {
    // handle error
  }
}

const deleteKey = (id: string) => {
  deleteTarget.value = id
}

const confirmDelete = async () => {
  if (!deleteTarget.value) return
  try {
    await fetchData(`/keys/${deleteTarget.value}`, { method: 'DELETE' })
    await fetchKeys()
    if (showPromptGuide.value) {
      showPromptGuide.value = false
      storeUrl.value = null
      stopCountdownAndPolling()
    }
    setupCompleted.value = false
  } catch {
    // handle error
  } finally {
    deleteTarget.value = null
  }
}

onMounted(fetchKeys)
onUnmounted(stopCountdownAndPolling)

const promptText = computed(() => {
  return t('dashboard.promptText', { storeKeyUrl: storeUrl.value || '' })
})
</script>

<template>
  <div class="min-h-[calc(100vh-4rem)] py-10 px-6">
    <div class="max-w-2xl mx-auto space-y-10">

      <!-- Page header -->
      <header>
        <h1 class="text-2xl font-bold text-th-text tracking-tight">{{ t('dashboard.title') }}</h1>
        <p class="text-base text-th-text-t mt-1.5">{{ t('dashboard.subtitle') }}</p>
      </header>

      <!-- ═══════ Section: Identity Profile ═══════ -->
      <section>
        <div class="flex items-center gap-2.5 mb-4">
          <div class="w-1 h-5 rounded-full bg-accent" />
          <h2 class="text-sm font-semibold uppercase tracking-widest text-th-text-s">{{ t('profile.sectionTitle') }}</h2>
        </div>

        <button
          class="w-full flex items-center justify-between gap-4 px-5 py-4 rounded-xl bg-th-bg-s border border-th-text/[var(--color-card-border-opacity)] hover:border-accent/25 transition-all duration-300 group"
          @click="showProfile = true"
        >
          <div class="flex items-center gap-4">
            <div class="w-10 h-10 rounded-lg bg-gradient-to-br from-accent/15 to-accent/5 flex items-center justify-center ring-1 ring-accent/10 group-hover:ring-accent/25 transition-all">
              <svg class="w-5 h-5 text-accent" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
              </svg>
            </div>
            <div class="text-left">
              <h3 class="text-base font-medium text-th-text">{{ t('profile.myProfile') }}</h3>
              <p class="text-sm text-th-text-t mt-0.5 leading-relaxed">{{ t('profile.myProfileDesc') }}</p>
            </div>
          </div>
          <svg class="w-5 h-5 text-th-text-m group-hover:text-accent/70 group-hover:translate-x-0.5 transition-all duration-200" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
          </svg>
        </button>
      </section>

      <!-- ═══════ Section: Agents ═══════ -->
      <section>
        <div class="flex items-center justify-between mb-4">
          <div class="flex items-center gap-2.5">
            <div class="w-1 h-5 rounded-full bg-highlight" />
            <h2 class="text-sm font-semibold uppercase tracking-widest text-th-text-s">{{ t('dashboard.agentsSection') }}</h2>
          </div>

          <button
            v-if="!showCreate"
            class="inline-flex items-center gap-1.5 px-4 py-2 rounded-lg text-sm font-semibold text-accent border border-accent/20 hover:bg-accent/10 hover:border-accent/40 transition-all duration-200"
            @click="showCreate = true"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
            </svg>
            {{ t('dashboard.addAgent') }}
          </button>
        </div>

        <!-- Setup complete banner -->
        <Transition
          enter-active-class="transition duration-300 ease-out"
          enter-from-class="opacity-0 -translate-y-2"
          enter-to-class="opacity-100 translate-y-0"
          leave-active-class="transition duration-200 ease-in"
          leave-from-class="opacity-100"
          leave-to-class="opacity-0"
        >
          <div v-if="setupCompleted" class="mb-5">
            <div class="p-5 rounded-xl bg-success/[0.06] border border-success/20">
              <div class="flex items-center gap-2 mb-1">
                <svg class="w-5 h-5 text-success" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                <h3 class="text-base font-semibold text-th-text">{{ t('dashboard.setupCompleteTitle') }}</h3>
              </div>
              <p class="text-sm text-th-text-s ml-7">{{ t('dashboard.setupComplete') }}</p>
            </div>
          </div>
        </Transition>

        <!-- Agent prompt guide -->
        <Transition
          enter-active-class="transition duration-300 ease-out"
          enter-from-class="opacity-0 -translate-y-2"
          enter-to-class="opacity-100 translate-y-0"
          leave-active-class="transition duration-200 ease-in"
          leave-from-class="opacity-100"
          leave-to-class="opacity-0"
        >
          <div v-if="showPromptGuide" class="mb-5">
            <div class="p-5 rounded-xl bg-accent/[0.04] border border-accent/15">
              <div class="flex items-center gap-2 mb-1">
                <svg class="w-5 h-5 text-accent" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
                </svg>
                <h3 class="text-base font-semibold text-th-text">{{ t('dashboard.promptTitle') }}</h3>
              </div>
              <p v-if="countdown > 0" class="text-sm text-th-text-s mb-3 ml-7">
                {{ t('dashboard.promptDesc', { time: formattedCountdown }) }}
              </p>
              <p v-else class="text-sm text-danger mb-3 ml-7">
                {{ t('dashboard.tokenExpired') }}
              </p>
              <div class="ml-7">
                <CopyBlock :text="promptText" label="AGENT PROMPT" prominent />
              </div>
            </div>
          </div>
        </Transition>

        <!-- Create key form -->
        <Transition
          enter-active-class="transition duration-200 ease-out"
          enter-from-class="opacity-0 -translate-y-1"
          enter-to-class="opacity-100 translate-y-0"
          leave-active-class="transition duration-150 ease-in"
          leave-from-class="opacity-100"
          leave-to-class="opacity-0"
        >
          <div v-if="showCreate" class="mb-5 p-5 rounded-xl bg-th-bg-s border border-th-text/[var(--color-card-border-opacity)]">
            <label class="block text-sm font-medium text-th-text-s mb-2">
              {{ t('dashboard.addAgentLabel') }}
            </label>
            <div class="flex gap-3">
              <div class="flex-1 relative">
                <input
                  v-model="newLabel"
                  type="text"
                  required
                  maxlength="15"
                  :placeholder="t('dashboard.addAgentPlaceholder')"
                  :class="['w-full px-4 py-2.5 pr-14 rounded-lg bg-th-input border text-th-text text-base placeholder-th-text-m focus:outline-none transition-all', labelError ? 'border-danger focus:border-danger focus:ring-1 focus:ring-danger/20' : 'border-th-text/[var(--color-input-border-opacity)] focus:border-accent/40 focus:ring-1 focus:ring-accent/20']"
                  @input="labelError = false"
                  @keyup.enter="createKey"
                >
                <span class="absolute right-3 top-1/2 -translate-y-1/2 text-xs text-th-text-m tabular-nums pointer-events-none">{{ newLabel.length }}/15</span>
              </div>
              <button
                class="px-5 py-2.5 rounded-lg text-sm font-semibold text-white bg-accent hover:bg-accent-light transition-colors"
                @click="createKey"
              >
                {{ t('dashboard.addAgent') }}
              </button>
            </div>
            <p v-if="labelError" class="mt-2 text-xs text-danger">{{ t('dashboard.agentNameRequired') }}</p>
          </div>
        </Transition>

        <!-- Keys list -->
        <div class="space-y-2.5">
          <div v-if="loading" class="flex items-center justify-center py-10">
            <div class="w-5 h-5 border-2 border-accent border-t-transparent rounded-full animate-spin" />
          </div>
          <div v-else-if="keys.length === 0" class="text-center py-10 text-th-text-t text-base">
            {{ t('dashboard.noAgents') }}
          </div>
          <template v-else>
            <AccessKeyCard
              v-for="key in keys"
              :key="key.id"
              :id="key.id"
              :key-prefix="key.key_prefix"
              :key-suffix="key.key_suffix"
              :label="key.label"
              :created-at="key.created_at"
              :last-used="key.last_used"
              @delete="deleteKey"
            />
          </template>
        </div>
      </section>

    </div>
  </div>

  <!-- Unlink confirm dialog -->
  <ConfirmDialog
    :open="!!deleteTarget"
    :title="t('dashboard.unlinkTitle')"
    :message="t('dashboard.unlinkConfirm')"
    :confirm-text="t('dashboard.unlinkAction')"
    :cancel-text="t('common.cancel')"
    variant="danger"
    @confirm="confirmDelete"
    @cancel="deleteTarget = null"
  />

  <!-- Profile dialog -->
  <ProfileDialog
    :open="showProfile"
    @close="showProfile = false"
  />
</template>
