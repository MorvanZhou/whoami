<script setup lang="ts">
const { t } = useI18n()
const visible = ref(false)
const container = ref<HTMLElement>()

onMounted(() => {
  const observer = new IntersectionObserver(
    ([entry]) => {
      if (entry.isIntersecting) visible.value = true
    },
    { threshold: 0.2 }
  )
  if (container.value) observer.observe(container.value)
  onUnmounted(() => observer.disconnect())
})

const steps = computed(() => [
  { icon: 'install', label: t('flow.step1'), desc: t('flow.step1Desc'), color: '#0EA5E9' },
  { icon: 'key', label: t('flow.step2'), desc: t('flow.step2Desc'), color: '#06B6D4' },
  { icon: 'read', label: t('flow.step3'), desc: t('flow.step3Desc'), color: '#22D3EE' },
  { icon: 'sync', label: t('flow.step4'), desc: t('flow.step4Desc'), color: '#10B981' },
])
</script>

<template>
  <section ref="container" class="py-24 px-6">
    <div class="max-w-5xl mx-auto">
      <h2 class="text-3xl sm:text-4xl font-bold text-center text-white mb-16">
        {{ t('flow.title') }}
      </h2>

      <!-- SVG Flow Diagram -->
      <div class="relative">
        <svg class="w-full" viewBox="0 0 900 220" fill="none" xmlns="http://www.w3.org/2000/svg">
          <!-- Connecting lines with flow animation -->
          <template v-if="visible">
            <!-- Line 1→2 -->
            <line x1="180" y1="90" x2="290" y2="90" stroke="url(#lineGrad1)" stroke-width="2" stroke-dasharray="6 4" class="animate-flow" />
            <!-- Line 2→3 -->
            <line x1="410" y1="90" x2="520" y2="90" stroke="url(#lineGrad2)" stroke-width="2" stroke-dasharray="6 4" class="animate-flow" style="animation-delay: 0.3s" />
            <!-- Line 3→4 -->
            <line x1="640" y1="90" x2="750" y2="90" stroke="url(#lineGrad3)" stroke-width="2" stroke-dasharray="6 4" class="animate-flow" style="animation-delay: 0.6s" />

            <!-- Flow particles -->
            <circle r="3" fill="#0EA5E9" opacity="0.8">
              <animateMotion dur="2s" repeatCount="indefinite" path="M180,90 L290,90" />
            </circle>
            <circle r="3" fill="#06B6D4" opacity="0.8">
              <animateMotion dur="2s" repeatCount="indefinite" path="M410,90 L520,90" begin="0.5s" />
            </circle>
            <circle r="3" fill="#22D3EE" opacity="0.8">
              <animateMotion dur="2s" repeatCount="indefinite" path="M640,90 L750,90" begin="1s" />
            </circle>
          </template>

          <!-- Gradients -->
          <defs>
            <linearGradient id="lineGrad1" x1="0" y1="0" x2="1" y2="0">
              <stop offset="0%" stop-color="#0EA5E9" />
              <stop offset="100%" stop-color="#06B6D4" />
            </linearGradient>
            <linearGradient id="lineGrad2" x1="0" y1="0" x2="1" y2="0">
              <stop offset="0%" stop-color="#06B6D4" />
              <stop offset="100%" stop-color="#22D3EE" />
            </linearGradient>
            <linearGradient id="lineGrad3" x1="0" y1="0" x2="1" y2="0">
              <stop offset="0%" stop-color="#22D3EE" />
              <stop offset="100%" stop-color="#10B981" />
            </linearGradient>
            <linearGradient id="nodeGrad1">
              <stop offset="0%" stop-color="#0EA5E9" stop-opacity="0.15" />
              <stop offset="100%" stop-color="#0EA5E9" stop-opacity="0.05" />
            </linearGradient>
            <linearGradient id="nodeGrad2">
              <stop offset="0%" stop-color="#06B6D4" stop-opacity="0.15" />
              <stop offset="100%" stop-color="#06B6D4" stop-opacity="0.05" />
            </linearGradient>
            <linearGradient id="nodeGrad3">
              <stop offset="0%" stop-color="#22D3EE" stop-opacity="0.15" />
              <stop offset="100%" stop-color="#22D3EE" stop-opacity="0.05" />
            </linearGradient>
            <linearGradient id="nodeGrad4">
              <stop offset="0%" stop-color="#10B981" stop-opacity="0.15" />
              <stop offset="100%" stop-color="#10B981" stop-opacity="0.05" />
            </linearGradient>
          </defs>

          <!-- Step Nodes -->
          <g v-for="(step, i) in steps" :key="i"
            :class="visible ? 'animate-scale-in' : 'opacity-0'"
            :style="{ animationDelay: `${i * 0.15}s` }"
          >
            <!-- Node background -->
            <rect
              :x="60 + i * 230 - 60" :y="30" width="120" height="120" rx="20"
              :fill="`url(#nodeGrad${i + 1})`"
              :stroke="step.color" stroke-width="1" stroke-opacity="0.3"
            />

            <!-- Icon -->
            <g :transform="`translate(${60 + i * 230 - 12}, 52)`">
              <!-- Install icon -->
              <template v-if="step.icon === 'install'">
                <path d="M12 3v12m0 0l-4-4m4 4l4-4M4 19h16" :stroke="step.color" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" fill="none" />
              </template>
              <!-- Key icon -->
              <template v-else-if="step.icon === 'key'">
                <path d="M15 7a4 4 0 00-8 0 4 4 0 003 3.87V17h2v-2h2v-2h-2v-1.13A4 4 0 0015 7zm-4-1.5a1.5 1.5 0 110 3 1.5 1.5 0 010-3z" :fill="step.color" />
              </template>
              <!-- Read icon -->
              <template v-else-if="step.icon === 'read'">
                <path d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" :stroke="step.color" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" fill="none" />
              </template>
              <!-- Sync icon -->
              <template v-else>
                <path d="M4 4v5h5M20 20v-5h-5M4 9a8 8 0 0113.5-4.5L20 7M20 15a8 8 0 01-13.5 4.5L4 17" :stroke="step.color" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" fill="none" />
              </template>
            </g>

            <!-- Label -->
            <text
              :x="60 + i * 230" y="175"
              text-anchor="middle" fill="white" font-size="14" font-weight="600" font-family="Satoshi, sans-serif"
            >
              {{ step.label }}
            </text>
            <!-- Description -->
            <text
              :x="60 + i * 230" y="195"
              text-anchor="middle" fill="#94A3B8" font-size="11" font-family="Space Mono, monospace"
            >
              {{ step.desc }}
            </text>
          </g>
        </svg>
      </div>
    </div>
  </section>
</template>
