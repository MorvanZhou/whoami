<script setup lang="ts">
const { t } = useI18n()
const visible = ref(false)
const container = ref<HTMLElement>()

onMounted(() => {
  const observer = new IntersectionObserver(
    ([entry]) => {
      if (entry?.isIntersecting) visible.value = true
    },
    { threshold: 0.15 }
  )
  if (container.value) observer.observe(container.value)
  onUnmounted(() => observer.disconnect())
})

const solutions = computed(() => [
  {
    key: 'syncAll',
    title: t('solutions.syncAll.title'),
    desc: t('solutions.syncAll.desc'),
    color: '#0EA5E9',
    colorLight: '#22D3EE',
  },
  {
    key: 'persist',
    title: t('solutions.persist.title'),
    desc: t('solutions.persist.desc'),
    color: '#10B981',
    colorLight: '#34D399',
  },
  {
    key: 'control',
    title: t('solutions.control.title'),
    desc: t('solutions.control.desc'),
    color: '#F59E0B',
    colorLight: '#FBBF24',
  },
])
</script>

<template>
  <section ref="container" class="py-24 px-6">
    <div class="max-w-5xl mx-auto">
      <h2 class="text-3xl sm:text-4xl font-bold text-center text-th-text mb-4">
        {{ t('solutions.title') }}
      </h2>
      <p class="text-center text-th-text-s mb-16 max-w-2xl mx-auto">
        {{ t('solutions.subtitle') }}
      </p>

      <div class="grid grid-cols-1 md:grid-cols-3 gap-8">
        <div
          v-for="(s, i) in solutions"
          :key="i"
          class="group relative rounded-2xl bg-th-bg-s/50 border border-th-text/[0.05] hover:border-th-text/[0.15] backdrop-blur-sm transition-all duration-300 hover:-translate-y-1 overflow-hidden"
          :class="visible ? 'animate-fade-in-up opacity-100' : 'opacity-0'"
          :style="{ animationDelay: `${i * 0.15}s` }"
        >
          <!-- SVG Illustration -->
          <div class="relative h-48 flex items-center justify-center overflow-hidden bg-th-bg-t/30">

            <!-- Solution 1: One profile syncs to all AIs -->
            <svg v-if="s.key === 'syncAll'" viewBox="0 0 220 130" class="w-48 h-auto" fill="none" xmlns="http://www.w3.org/2000/svg">
              <!-- Center whoami profile node -->
              <rect x="75" y="30" width="70" height="44" rx="10" :fill="s.color" fill-opacity="0.12" :stroke="s.color" stroke-width="1.5" />
              <text x="110" y="48" text-anchor="middle" :fill="s.color" font-size="9" font-weight="700" font-family="Space Mono, monospace">whoami</text>
              <text x="110" y="62" text-anchor="middle" :fill="s.colorLight" font-size="8" font-family="Satoshi, sans-serif">Your Profile</text>

              <!-- AI-A top-left -->
              <rect x="10" y="92" width="48" height="28" rx="6" fill-opacity="0.08" :fill="s.color" :stroke="s.color" stroke-width="1" />
              <text x="34" y="110" text-anchor="middle" :fill="s.color" font-size="8" font-weight="600" font-family="Space Mono, monospace">Claude</text>

              <!-- AI-B center -->
              <rect x="86" y="92" width="48" height="28" rx="6" fill-opacity="0.08" :fill="s.color" :stroke="s.color" stroke-width="1" />
              <text x="110" y="110" text-anchor="middle" :fill="s.color" font-size="8" font-weight="600" font-family="Space Mono, monospace">GPT</text>

              <!-- AI-C top-right -->
              <rect x="162" y="92" width="48" height="28" rx="6" fill-opacity="0.08" :fill="s.color" :stroke="s.color" stroke-width="1" />
              <text x="186" y="110" text-anchor="middle" :fill="s.color" font-size="8" font-weight="600" font-family="Space Mono, monospace">Cursor</text>

              <!-- Connecting lines from center to AIs -->
              <line x1="90" y1="74" x2="40" y2="92" :stroke="s.color" stroke-width="1.2" stroke-dasharray="4 3">
                <animate attributeName="stroke-dashoffset" values="0;-14" dur="2s" repeatCount="indefinite" />
              </line>
              <line x1="110" y1="74" x2="110" y2="92" :stroke="s.color" stroke-width="1.2" stroke-dasharray="4 3">
                <animate attributeName="stroke-dashoffset" values="0;-14" dur="2s" repeatCount="indefinite" begin="0.3s" />
              </line>
              <line x1="130" y1="74" x2="180" y2="92" :stroke="s.color" stroke-width="1.2" stroke-dasharray="4 3">
                <animate attributeName="stroke-dashoffset" values="0;-14" dur="2s" repeatCount="indefinite" begin="0.6s" />
              </line>

              <!-- Flow particles -->
              <circle r="2.5" :fill="s.colorLight" opacity="0.9">
                <animateMotion dur="2s" repeatCount="indefinite" path="M90,74 L40,92" />
              </circle>
              <circle r="2.5" :fill="s.colorLight" opacity="0.9">
                <animateMotion dur="2s" repeatCount="indefinite" path="M110,74 L110,92" begin="0.3s" />
              </circle>
              <circle r="2.5" :fill="s.colorLight" opacity="0.9">
                <animateMotion dur="2s" repeatCount="indefinite" path="M130,74 L180,92" begin="0.6s" />
              </circle>

              <!-- Checkmarks on each AI -->
              <g v-if="visible">
                <circle cx="50" cy="95" r="5" fill="#10B981" fill-opacity="0.2" />
                <path d="M48 95 l1.5 1.5 3-3" stroke="#10B981" stroke-width="1.2" stroke-linecap="round" stroke-linejoin="round" fill="none" />
                <circle cx="126" cy="95" r="5" fill="#10B981" fill-opacity="0.2" />
                <path d="M124 95 l1.5 1.5 3-3" stroke="#10B981" stroke-width="1.2" stroke-linecap="round" stroke-linejoin="round" fill="none" />
                <circle cx="202" cy="95" r="5" fill="#10B981" fill-opacity="0.2" />
                <path d="M200 95 l1.5 1.5 3-3" stroke="#10B981" stroke-width="1.2" stroke-linecap="round" stroke-linejoin="round" fill="none" />
              </g>
            </svg>

            <!-- Solution 2: Profile persists across sessions -->
            <svg v-else-if="s.key === 'persist'" viewBox="0 0 220 130" class="w-48 h-auto" fill="none" xmlns="http://www.w3.org/2000/svg">
              <!-- Cloud storage symbol -->
              <path d="M80 38 Q80 22 96 22 Q104 14 116 20 Q130 16 136 28 Q148 28 148 40 Q148 52 136 52 L80 52 Q68 52 68 40 Q68 32 80 38Z" :fill="s.color" fill-opacity="0.1" :stroke="s.color" stroke-width="1.5" />
              <text x="108" y="42" text-anchor="middle" :fill="s.color" font-size="8" font-weight="600" font-family="Space Mono, monospace">whoami</text>

              <!-- Multiple chat sessions below -->
              <!-- Session 1 -->
              <rect x="18" y="72" width="52" height="36" rx="6" fill-opacity="0.08" :fill="s.color" :stroke="s.color" stroke-width="1" />
              <rect x="24" y="80" width="28" height="4" rx="2" :fill="s.color" fill-opacity="0.25" />
              <rect x="24" y="88" width="20" height="4" rx="2" :fill="s.color" fill-opacity="0.15" />
              <text x="44" y="102" text-anchor="middle" class="fill-th-text-t" font-size="7" font-family="Space Mono, monospace">Chat 1</text>

              <!-- Session 2 -->
              <rect x="84" y="72" width="52" height="36" rx="6" fill-opacity="0.08" :fill="s.color" :stroke="s.color" stroke-width="1" />
              <rect x="90" y="80" width="28" height="4" rx="2" :fill="s.color" fill-opacity="0.25" />
              <rect x="90" y="88" width="20" height="4" rx="2" :fill="s.color" fill-opacity="0.15" />
              <text x="110" y="102" text-anchor="middle" class="fill-th-text-t" font-size="7" font-family="Space Mono, monospace">Chat 2</text>

              <!-- Session 3 -->
              <rect x="150" y="72" width="52" height="36" rx="6" fill-opacity="0.08" :fill="s.color" :stroke="s.color" stroke-width="1" />
              <rect x="156" y="80" width="28" height="4" rx="2" :fill="s.color" fill-opacity="0.25" />
              <rect x="156" y="88" width="20" height="4" rx="2" :fill="s.color" fill-opacity="0.15" />
              <text x="176" y="102" text-anchor="middle" class="fill-th-text-t" font-size="7" font-family="Space Mono, monospace">Chat 3</text>

              <!-- Down arrows from cloud to sessions -->
              <line x1="44" y1="52" x2="44" y2="72" :stroke="s.color" stroke-width="1.2" stroke-dasharray="3 2">
                <animate attributeName="stroke-dashoffset" values="0;-10" dur="1.5s" repeatCount="indefinite" />
              </line>
              <line x1="110" y1="52" x2="110" y2="72" :stroke="s.color" stroke-width="1.2" stroke-dasharray="3 2">
                <animate attributeName="stroke-dashoffset" values="0;-10" dur="1.5s" repeatCount="indefinite" begin="0.2s" />
              </line>
              <line x1="176" y1="52" x2="176" y2="72" :stroke="s.color" stroke-width="1.2" stroke-dasharray="3 2">
                <animate attributeName="stroke-dashoffset" values="0;-10" dur="1.5s" repeatCount="indefinite" begin="0.4s" />
              </line>

              <!-- Particles -->
              <circle r="2" :fill="s.colorLight" opacity="0.9">
                <animateMotion dur="1.5s" repeatCount="indefinite" path="M44,52 L44,72" />
              </circle>
              <circle r="2" :fill="s.colorLight" opacity="0.9">
                <animateMotion dur="1.5s" repeatCount="indefinite" path="M110,52 L110,72" begin="0.2s" />
              </circle>
              <circle r="2" :fill="s.colorLight" opacity="0.9">
                <animateMotion dur="1.5s" repeatCount="indefinite" path="M176,52 L176,72" begin="0.4s" />
              </circle>

              <!-- Label -->
              <text x="110" y="122" text-anchor="middle" class="fill-th-text-t" font-size="8" font-family="Satoshi, sans-serif">{{ t('solutions.persist.label') }}</text>
            </svg>

            <!-- Solution 3: Full control to edit profile -->
            <svg v-else viewBox="0 0 220 130" class="w-48 h-auto" fill="none" xmlns="http://www.w3.org/2000/svg">
              <!-- Profile card -->
              <rect x="50" y="14" width="120" height="72" rx="10" :fill="s.color" fill-opacity="0.08" :stroke="s.color" stroke-width="1.5" />

              <!-- Profile content lines -->
              <circle cx="78" cy="36" r="10" :stroke="s.color" stroke-width="1.2" fill="none" />
              <rect x="96" y="30" width="56" height="5" rx="2.5" :fill="s.color" fill-opacity="0.3" />
              <rect x="96" y="40" width="38" height="4" rx="2" :fill="s.color" fill-opacity="0.15" />
              <rect x="62" y="56" width="96" height="4" rx="2" :fill="s.color" fill-opacity="0.12" />
              <rect x="62" y="64" width="72" height="4" rx="2" :fill="s.color" fill-opacity="0.08" />
              <rect x="62" y="72" width="84" height="4" rx="2" :fill="s.color" fill-opacity="0.06" />

              <!-- Edit pencil icon (animated) -->
              <g transform="translate(148, 12)">
                <circle cx="12" cy="12" r="14" :fill="s.color" fill-opacity="0.15" />
                <g>
                  <animateTransform attributeName="transform" type="rotate" values="0 12 12;-8 12 12;0 12 12;8 12 12;0 12 12" dur="3s" repeatCount="indefinite" />
                  <path d="M7 17 L5 19 L7.5 18.5 L17 9 L15 7 Z" :fill="s.color" fill-opacity="0.6" />
                  <path d="M15 7 L17 9 M5 19 L7 17" :stroke="s.color" stroke-width="1.5" stroke-linecap="round" fill="none" />
                </g>
              </g>

              <!-- Animated typing cursor on a line -->
              <rect x="146" y="56" width="2" height="10" rx="1" :fill="s.colorLight">
                <animate attributeName="opacity" values="1;0;1" dur="1s" repeatCount="indefinite" />
              </rect>

              <!-- "Save" button -->
              <rect x="104" y="96" width="56" height="20" rx="6" :fill="s.color" fill-opacity="0.2" :stroke="s.color" stroke-width="1" />
              <text x="132" y="109" text-anchor="middle" :fill="s.color" font-size="8" font-weight="600" font-family="Space Mono, monospace">{{ t('solutions.control.save') }}</text>

              <!-- Sync arrow from save to cloud -->
              <path d="M60 106 Q40 106 40 94 Q40 82 52 82" :stroke="s.color" stroke-width="1" stroke-dasharray="3 2" fill="none">
                <animate attributeName="stroke-dashoffset" values="0;-10" dur="2s" repeatCount="indefinite" />
              </path>
              <!-- Cloud mini icon -->
              <path d="M28 82 Q28 74 36 74 Q40 70 46 73 Q52 71 54 76 Q60 76 60 82 Q60 88 54 88 L28 88 Q22 88 22 82 Q22 78 28 82Z" :fill="s.color" fill-opacity="0.1" :stroke="s.color" stroke-width="1" />

              <!-- Label -->
              <text x="110" y="126" text-anchor="middle" class="fill-th-text-t" font-size="8" font-family="Satoshi, sans-serif">{{ t('solutions.control.label') }}</text>
            </svg>
          </div>

          <!-- Text Content -->
          <div class="p-6 pt-5">
            <h3 class="text-lg font-semibold text-th-text mb-2">{{ s.title }}</h3>
            <p class="text-sm text-th-text-s leading-relaxed">{{ s.desc }}</p>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>
