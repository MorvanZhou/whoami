<script setup lang="ts">
const { t } = useI18n()

const cases = computed(() => [
  {
    key: 'switchAI',
    title: t('features.switchAI.title'),
    desc: t('features.switchAI.desc'),
    color: '#0EA5E9',
    colorLight: '#22D3EE',
  },
  {
    key: 'newChat',
    title: t('features.newChat.title'),
    desc: t('features.newChat.desc'),
    color: '#F59E0B',
    colorLight: '#FBBF24',
  },
  {
    key: 'locked',
    title: t('features.locked.title'),
    desc: t('features.locked.desc'),
    color: '#F43F5E',
    colorLight: '#FB7185',
  },
])
</script>

<template>
  <section class="py-24 px-6">
    <div class="max-w-5xl mx-auto">
      <h2 class="text-3xl sm:text-4xl font-bold text-center text-th-text mb-4">
        {{ t('features.title') }}
      </h2>
      <p class="text-center text-th-text-s mb-16 max-w-2xl mx-auto">
        {{ t('features.subtitle') }}
      </p>

      <div class="grid grid-cols-1 md:grid-cols-3 gap-8">
        <div
          v-for="(c, i) in cases"
          :key="i"
          class="group relative rounded-2xl bg-th-bg-s/50 border border-th-text/[0.05] hover:border-th-text/[0.15] backdrop-blur-sm transition-all duration-300 hover:-translate-y-1 overflow-hidden"
        >
          <!-- SVG Animation Area -->
          <div class="relative h-48 flex items-center justify-center overflow-hidden bg-th-bg-t/30">
            <!-- Switch AI: two chat bubbles with a broken arrow -->
            <svg v-if="c.key === 'switchAI'" viewBox="0 0 200 120" class="w-44 h-auto" fill="none" xmlns="http://www.w3.org/2000/svg">
              <!-- Left AI bubble -->
              <rect x="10" y="20" width="60" height="40" rx="8" :fill="c.color" fill-opacity="0.15" :stroke="c.color" stroke-width="1.5" />
              <text x="40" y="44" text-anchor="middle" :fill="c.color" font-size="11" font-weight="600" font-family="Space Mono, monospace">AI-A</text>
              <!-- Smile face on left -->
              <circle cx="30" cy="78" r="8" fill="none" :stroke="c.colorLight" stroke-width="1.2" />
              <path d="M27 80 Q30 83 33 80" fill="none" :stroke="c.colorLight" stroke-width="1" stroke-linecap="round" />
              <circle cx="28" cy="77" r="1" :fill="c.colorLight" />
              <circle cx="32" cy="77" r="1" :fill="c.colorLight" />

              <!-- Arrow with break / question mark -->
              <line x1="80" y1="40" x2="110" y2="40" :stroke="c.color" stroke-width="1.5" stroke-dasharray="4 3">
                <animate attributeName="stroke-dashoffset" values="0;-14" dur="1.5s" repeatCount="indefinite" />
              </line>
              <text x="95" y="32" text-anchor="middle" fill="#F43F5E" font-size="16" font-weight="bold">?</text>

              <!-- Right AI bubble -->
              <rect x="120" y="20" width="60" height="40" rx="8" fill-opacity="0.15" :fill="c.color" :stroke="c.color" stroke-width="1.5" stroke-dasharray="4 2" />
              <text x="150" y="44" text-anchor="middle" :fill="c.color" font-size="11" font-weight="600" font-family="Space Mono, monospace">AI-B</text>
              <!-- Confused face on right -->
              <circle cx="160" cy="78" r="8" fill="none" stroke="#94A3B8" stroke-width="1.2" />
              <path d="M157 81 Q160 78 163 81" fill="none" stroke="#94A3B8" stroke-width="1" stroke-linecap="round" />
              <circle cx="158" cy="77" r="1" fill="#94A3B8" />
              <circle cx="162" cy="77" r="1" fill="#94A3B8" />

              <!-- Label -->
              <text x="100" y="108" text-anchor="middle" class="fill-th-text-t" font-size="9" font-family="Satoshi, sans-serif">{{ t('features.switchAI.label') }}</text>
            </svg>

            <!-- New Chat: chat history vanishing -->
            <svg v-else-if="c.key === 'newChat'" viewBox="0 0 200 120" class="w-44 h-auto" fill="none" xmlns="http://www.w3.org/2000/svg">
              <!-- Chat window frame -->
              <rect x="40" y="10" width="120" height="80" rx="10" fill-opacity="0.1" :fill="c.color" :stroke="c.color" stroke-width="1.5" />
              <!-- Title bar dots -->
              <circle cx="54" cy="22" r="3" fill="#F43F5E" fill-opacity="0.6" />
              <circle cx="64" cy="22" r="3" :fill="c.color" fill-opacity="0.6" />
              <circle cx="74" cy="22" r="3" fill="#10B981" fill-opacity="0.6" />

              <!-- Chat lines fading out -->
              <g>
                <rect x="52" y="34" width="50" height="6" rx="3" :fill="c.color" fill-opacity="0.3">
                  <animate attributeName="fill-opacity" values="0.3;0.05;0.3" dur="3s" repeatCount="indefinite" />
                </rect>
                <rect x="52" y="46" width="70" height="6" rx="3" :fill="c.color" fill-opacity="0.2">
                  <animate attributeName="fill-opacity" values="0.2;0.03;0.2" dur="3s" repeatCount="indefinite" begin="0.3s" />
                </rect>
                <rect x="52" y="58" width="40" height="6" rx="3" :fill="c.color" fill-opacity="0.1">
                  <animate attributeName="fill-opacity" values="0.1;0.01;0.1" dur="3s" repeatCount="indefinite" begin="0.6s" />
                </rect>
              </g>

              <!-- "New Chat" badge -->
              <rect x="95" y="70" width="56" height="16" rx="8" fill="#F43F5E" fill-opacity="0.15" stroke="#F43F5E" stroke-width="1" />
              <text x="123" y="81" text-anchor="middle" fill="#F43F5E" font-size="8" font-weight="600" font-family="Space Mono, monospace">NEW CHAT</text>

              <!-- Label -->
              <text x="100" y="108" text-anchor="middle" class="fill-th-text-t" font-size="9" font-family="Satoshi, sans-serif">{{ t('features.newChat.label') }}</text>
            </svg>

            <!-- Locked: padlock on profile -->
            <svg v-else viewBox="0 0 200 120" class="w-44 h-auto" fill="none" xmlns="http://www.w3.org/2000/svg">
              <!-- Profile card -->
              <rect x="30" y="15" width="100" height="70" rx="10" fill-opacity="0.1" :fill="c.color" :stroke="c.color" stroke-width="1.5" />
              <!-- Profile lines -->
              <circle cx="58" cy="38" r="10" :stroke="c.color" stroke-width="1.5" fill="none" />
              <rect x="74" y="32" width="42" height="5" rx="2.5" :fill="c.color" fill-opacity="0.3" />
              <rect x="74" y="42" width="30" height="4" rx="2" :fill="c.color" fill-opacity="0.15" />
              <rect x="42" y="58" width="76" height="4" rx="2" :fill="c.color" fill-opacity="0.1" />
              <rect x="42" y="66" width="56" height="4" rx="2" :fill="c.color" fill-opacity="0.08" />

              <!-- Lock icon overlaid -->
              <g transform="translate(140, 25)">
                <rect x="4" y="16" width="32" height="24" rx="4" :fill="c.color" fill-opacity="0.2" :stroke="c.color" stroke-width="1.5" />
                <path d="M12 16V10a8 8 0 0116 0v6" fill="none" :stroke="c.color" stroke-width="1.5" stroke-linecap="round" />
                <circle cx="20" cy="30" r="3" :fill="c.color">
                  <animate attributeName="fill-opacity" values="0.4;1;0.4" dur="2s" repeatCount="indefinite" />
                </circle>
              </g>

              <!-- X marks for "can't edit" -->
              <g opacity="0.6">
                <line x1="42" y1="56" x2="48" y2="68" stroke="#F43F5E" stroke-width="1.5" stroke-linecap="round">
                  <animate attributeName="opacity" values="0;1;0" dur="2.5s" repeatCount="indefinite" />
                </line>
                <line x1="48" y1="56" x2="42" y2="68" stroke="#F43F5E" stroke-width="1.5" stroke-linecap="round">
                  <animate attributeName="opacity" values="0;1;0" dur="2.5s" repeatCount="indefinite" />
                </line>
              </g>

              <!-- Label -->
              <text x="100" y="108" text-anchor="middle" class="fill-th-text-t" font-size="9" font-family="Satoshi, sans-serif">{{ t('features.locked.label') }}</text>
            </svg>
          </div>

          <!-- Text Content -->
          <div class="p-6 pt-5">
            <h3 class="text-lg font-semibold text-th-text mb-2">{{ c.title }}</h3>
            <p class="text-sm text-th-text-s leading-relaxed">{{ c.desc }}</p>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>
