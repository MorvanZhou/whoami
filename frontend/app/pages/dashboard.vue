<script setup lang="ts">
definePageMeta({
  middleware: 'auth',
})

const { t } = useI18n()
const { apiFetch } = useApi()

useSeoMeta({
  title: () => t('seo.dashboard.title'),
  description: () => t('seo.dashboard.description'),
  robots: 'noindex, nofollow',
})

interface ApiKeyItem {
  id: string
  key_prefix: string
  key_suffix: string
  label: string | null
  created_at: string
  last_used: string | null
}

const keys = ref<ApiKeyItem[]>([])
const newKey = ref<string | null>(null)
const storeUrl = ref<string | null>(null)
const loading = ref(false)
const showCreate = ref(false)
const newLabel = ref('')
const deleteTarget = ref<string | null>(null)
const isInitialLoad = ref(true)
const showProfile = ref(false)

const fetchKeys = async () => {
  loading.value = true
  try {
    keys.value = await apiFetch<ApiKeyItem[]>('/keys')
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
    const result = await apiFetch<{ id: string; key: string; key_prefix: string; key_suffix: string; label: string | null; created_at: string; store_url: string }>('/keys', {
      method: 'POST',
      body: { label: 'First Key' },
    })
    newKey.value = result.key
    storeUrl.value = result.store_url
    keys.value = await apiFetch<ApiKeyItem[]>('/keys')
  } catch {
    // handle error
  }
}

const createKey = async () => {
  try {
    const result = await apiFetch<{ id: string; key: string; key_prefix: string; key_suffix: string; label: string | null; created_at: string; store_url: string }>('/keys', {
      method: 'POST',
      body: { label: newLabel.value || null },
    })
    newKey.value = result.key
    storeUrl.value = result.store_url
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
    await apiFetch(`/keys/${deleteTarget.value}`, { method: 'DELETE' })
    await fetchKeys()
    if (newKey.value) {
      newKey.value = null
      storeUrl.value = null
    }
  } catch {
    // handle error
  } finally {
    deleteTarget.value = null
  }
}

onMounted(fetchKeys)

const promptText = computed(() => {
  return t('dashboard.promptText', { storeApiUrl: storeUrl.value || '' })
})
</script>

<template>
  <div class="min-h-[calc(100vh-4rem)] py-10 px-6">
    <div class="max-w-2xl mx-auto space-y-10">

      <!-- Page header -->
      <header>
        <h1 class="text-2xl font-bold text-th-text tracking-tight">{{ t('dashboard.title') }}</h1>
        <p class="text-sm text-th-text-t mt-1.5">{{ t('dashboard.subtitle') }}</p>
      </header>

      <!-- ═══════ Section: Identity Profile ═══════ -->
      <section>
        <div class="flex items-center gap-2.5 mb-4">
          <div class="w-1 h-4 rounded-full bg-accent" />
          <h2 class="text-xs font-semibold uppercase tracking-widest text-th-text-s">{{ t('profile.sectionTitle') }}</h2>
        </div>

        <button
          class="w-full flex items-center justify-between gap-4 px-5 py-4 rounded-xl bg-th-bg-s border border-th-text/[var(--color-card-border-opacity)] hover:border-accent/25 transition-all duration-300 group"
          @click="showProfile = true"
        >
          <div class="flex items-center gap-4">
            <div class="w-9 h-9 rounded-lg bg-gradient-to-br from-accent/15 to-accent/5 flex items-center justify-center ring-1 ring-accent/10 group-hover:ring-accent/25 transition-all">
              <svg class="w-4 h-4 text-accent" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
              </svg>
            </div>
            <div class="text-left">
              <h3 class="text-sm font-medium text-th-text">{{ t('profile.myProfile') }}</h3>
              <p class="text-xs text-th-text-t mt-0.5 leading-relaxed">{{ t('profile.myProfileDesc') }}</p>
            </div>
          </div>
          <svg class="w-4 h-4 text-th-text-m group-hover:text-accent/70 group-hover:translate-x-0.5 transition-all duration-200" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
          </svg>
        </button>
      </section>

      <!-- ═══════ Section: API Keys ═══════ -->
      <section>
        <div class="flex items-center justify-between mb-4">
          <div class="flex items-center gap-2.5">
            <div class="w-1 h-4 rounded-full bg-highlight" />
            <h2 class="text-xs font-semibold uppercase tracking-widest text-th-text-s">{{ t('dashboard.keysSection') }}</h2>
          </div>

          <button
            v-if="!showCreate"
            class="inline-flex items-center gap-1.5 px-3.5 py-1.5 rounded-lg text-xs font-semibold text-accent border border-accent/20 hover:bg-accent/10 hover:border-accent/40 transition-all duration-200"
            @click="showCreate = true"
          >
            <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
            </svg>
            {{ t('dashboard.createKey') }}
          </button>
        </div>

        <!-- New key alert + Agent prompt guide -->
        <Transition
          enter-active-class="transition duration-300 ease-out"
          enter-from-class="opacity-0 -translate-y-2"
          enter-to-class="opacity-100 translate-y-0"
          leave-active-class="transition duration-200 ease-in"
          leave-from-class="opacity-100"
          leave-to-class="opacity-0"
        >
          <div v-if="newKey" class="mb-5 space-y-3">
            <!-- Agent prompt guide -->
            <div class="p-5 rounded-xl bg-accent/[0.04] border border-accent/15">
              <div class="flex items-center gap-2 mb-1">
                <svg class="w-4 h-4 text-accent" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
                </svg>
                <h3 class="text-sm font-semibold text-th-text">{{ t('dashboard.promptTitle') }}</h3>
              </div>
              <p class="text-xs text-th-text-s mb-3 ml-6">{{ t('dashboard.promptDesc') }}</p>
              <div class="ml-6">
                <CopyBlock :text="promptText" label="AGENT PROMPT" prominent />
              </div>
            </div>

            <!-- Raw key display -->
            <div class="p-4 rounded-xl bg-th-bg-s border border-th-text/[var(--color-card-border-opacity)]">
              <div class="flex items-center gap-2 mb-2.5">
                <svg class="w-3.5 h-3.5 text-warning/70" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
                </svg>
                <span class="text-xs font-medium text-warning/70">{{ t('dashboard.newKeyWarning') }}</span>
              </div>
              <CopyBlock :text="newKey" :mono="true" />
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
            <label class="block text-xs font-medium text-th-text-s mb-2">
              {{ t('dashboard.createKeyLabel') }}
            </label>
            <div class="flex gap-3">
              <input
                v-model="newLabel"
                type="text"
                :placeholder="t('dashboard.createKeyPlaceholder')"
                class="flex-1 px-4 py-2.5 rounded-lg bg-th-input border border-th-text/[var(--color-input-border-opacity)] text-th-text text-sm placeholder-th-text-m focus:outline-none focus:border-accent/40 focus:ring-1 focus:ring-accent/20 transition-all"
                @keyup.enter="createKey"
              >
              <button
                class="px-5 py-2.5 rounded-lg text-sm font-semibold text-white bg-accent hover:bg-accent-light transition-colors"
                @click="createKey"
              >
                {{ t('dashboard.createKey') }}
              </button>
            </div>
          </div>
        </Transition>

        <!-- Keys list -->
        <div class="space-y-2">
          <div v-if="loading" class="flex items-center justify-center py-10">
            <div class="w-5 h-5 border-2 border-accent border-t-transparent rounded-full animate-spin" />
          </div>
          <div v-else-if="keys.length === 0" class="text-center py-10 text-th-text-t text-sm">
            {{ t('dashboard.noKeys') }}
          </div>
          <template v-else>
            <ApiKeyCard
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

  <!-- Delete confirm dialog -->
  <ConfirmDialog
    :open="!!deleteTarget"
    :title="t('dashboard.deleteTitle')"
    :message="t('dashboard.deleteConfirm')"
    :confirm-text="t('dashboard.deleteKey')"
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
