<script setup lang="ts">
const { t } = useI18n()

interface ChatStep {
  type: 'user' | 'agent' | 'skill' | 'typing' | 'divider'
  textKey?: string
  delay: number
}

const steps: ChatStep[] = [
  { type: 'user', textKey: 'hero.demo.user1', delay: 600 },
  { type: 'typing', delay: 800 },
  { type: 'agent', textKey: 'hero.demo.agent1', delay: 0 },
  { type: 'skill', textKey: 'hero.demo.skillRead', delay: 400 },
  { type: 'agent', textKey: 'hero.demo.agent2', delay: 600 },
  { type: 'divider', delay: 400 },
  { type: 'user', textKey: 'hero.demo.user2', delay: 600 },
  { type: 'typing', delay: 800 },
  { type: 'agent', textKey: 'hero.demo.agent3', delay: 0 },
  { type: 'skill', textKey: 'hero.demo.skillWrite', delay: 400 },
  { type: 'agent', textKey: 'hero.demo.agent4', delay: 600 },
]

interface RenderedItem {
  id: number
  type: string
  text: string
}

const renderedItems = ref<RenderedItem[]>([])
const typingVisible = ref(false)
const currentTypingText = ref('')
const isTypingAgent = ref(false)
const chatContainer = ref<HTMLElement | null>(null)
const loadedSkills = reactive(new Set<number>())
let itemId = 0

const scrollToBottom = () => {
  nextTick(() => {
    if (chatContainer.value) {
      chatContainer.value.scrollTop = chatContainer.value.scrollHeight
    }
  })
}

const typeText = (fullText: string, speed: number = 25): Promise<void> => {
  return new Promise((resolve) => {
    typingVisible.value = false
    isTypingAgent.value = true
    currentTypingText.value = ''
    let i = 0
    const interval = setInterval(() => {
      if (i < fullText.length) {
        currentTypingText.value += fullText[i]
        i++
        scrollToBottom()
      } else {
        clearInterval(interval)
        // Don't clear typing state here — caller will handle the transition
        resolve()
      }
    }, speed)
  })
}

const sleep = (ms: number) => new Promise(r => setTimeout(r, ms))

const addItem = (type: string, text: string): number => {
  const id = itemId++
  renderedItems.value.push({ id, type, text })
  scrollToBottom()
  return id
}

const markSkillLoaded = (id: number) => {
  loadedSkills.add(id)
}

const runSequence = async () => {
  renderedItems.value = []
  loadedSkills.clear()
  typingVisible.value = false
  currentTypingText.value = ''
  isTypingAgent.value = false
  itemId = 0
  await sleep(800)

  let lastSkillId: number | null = null

  for (const step of steps) {
    if (step.type === 'typing') {
      typingVisible.value = true
      scrollToBottom()
      await sleep(step.delay)
      // Don't hide dots here — the next step (agent typeText) will replace it seamlessly
      continue
    }

    if (step.type === 'agent' && step.textKey) {
      // Before showing agent reply, mark the previous skill as loaded
      if (lastSkillId !== null) {
        markSkillLoaded(lastSkillId)
        lastSkillId = null
      }
      const fullText = t(step.textKey)
      await typeText(fullText)
      // Push the final item FIRST, then clear typing state in the same tick
      // This prevents the flash caused by the typing bubble disappearing before the new item appears
      addItem('agent', fullText)
      isTypingAgent.value = false
      currentTypingText.value = ''
      await sleep(step.delay)
    } else if (step.type === 'skill') {
      const text = step.textKey ? t(step.textKey) : ''
      lastSkillId = addItem('skill', text)
      await sleep(step.delay)
    } else {
      typingVisible.value = false
      isTypingAgent.value = false
      currentTypingText.value = ''
      const text = step.textKey ? t(step.textKey) : ''
      addItem(step.type, text)
      await sleep(step.delay)
    }
  }

  // Mark any remaining skill as loaded
  if (lastSkillId !== null) {
    markSkillLoaded(lastSkillId)
  }
}

onMounted(() => {
  runSequence()
})
</script>

<template>
  <div class="w-full max-w-md mx-auto lg:mx-0 font-system">
    <!-- Window chrome -->
    <div class="rounded-2xl border border-th-text/[0.08] bg-th-bg-s/80 backdrop-blur-md shadow-2xl shadow-accent/5 overflow-hidden">
      <!-- Title bar -->
      <div class="flex items-center gap-2 px-4 py-3 border-b border-th-text/[0.06] bg-th-bg-t/30">
        <div class="flex gap-1.5">
          <span class="w-3 h-3 rounded-full bg-danger/60" />
          <span class="w-3 h-3 rounded-full bg-warning/60" />
          <span class="w-3 h-3 rounded-full bg-success/60" />
        </div>
        <span class="flex-1 text-center text-xs font-mono text-th-text-m tracking-wide">AI Agent</span>
        <div class="w-12" />
      </div>

      <!-- Chat body -->
      <div
        ref="chatContainer"
        class="px-4 py-4 space-y-3 h-[340px] sm:h-[380px] overflow-y-auto scrollbar-hide"
      >
        <template v-for="item in renderedItems" :key="item.id">
          <!-- User message -->
          <div v-if="item.type === 'user'" class="flex justify-end animate-chat-in">
            <div class="max-w-[85%] px-3.5 py-2.5 rounded-2xl rounded-br-md bg-accent/15 border border-accent/20 text-sm text-th-text leading-relaxed">
              {{ item.text }}
            </div>
          </div>

          <!-- Agent message -->
          <div v-else-if="item.type === 'agent'" class="flex justify-start animate-chat-in">
            <div class="flex gap-2.5 max-w-[90%]">
              <div class="shrink-0 w-7 h-7 rounded-lg bg-gradient-to-br from-accent/20 to-accent-light/10 flex items-center justify-center ring-1 ring-accent/15 mt-0.5">
                <svg class="w-3.5 h-3.5 text-accent" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M9.813 15.904L9 18.75l-.813-2.846a4.5 4.5 0 00-3.09-3.09L2.25 12l2.846-.813a4.5 4.5 0 003.09-3.09L9 5.25l.813 2.846a4.5 4.5 0 003.09 3.09L15.75 12l-2.846.813a4.5 4.5 0 00-3.09 3.09z" />
                </svg>
              </div>
              <div class="px-3.5 py-2.5 rounded-2xl rounded-bl-md bg-th-bg-t/60 border border-th-text/[0.05] text-sm text-th-text-s leading-relaxed">
                {{ item.text }}
              </div>
            </div>
          </div>

          <!-- Skill card -->
          <div v-else-if="item.type === 'skill'" class="flex justify-start animate-chat-in">
            <div class="flex gap-2.5 max-w-[90%]">
              <div class="w-7 shrink-0" />
              <div class="flex items-center gap-2.5 px-3.5 py-2.5 rounded-xl bg-accent/[0.06] border border-accent/15">
                <div class="w-8 h-8 rounded-lg bg-accent/10 flex items-center justify-center ring-1 ring-accent/20">
                  <svg class="w-4 h-4 text-accent" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M15.75 6a3.75 3.75 0 11-7.5 0 3.75 3.75 0 017.5 0zM4.501 20.118a7.5 7.5 0 0114.998 0" />
                  </svg>
                </div>
                <div class="min-w-0">
                  <div class="text-xs font-semibold text-accent font-mono leading-tight">whoami</div>
                  <div class="text-[11px] text-th-text-m leading-tight mt-0.5">{{ item.text }}</div>
                </div>
                <!-- Loading spinner or checkmark -->
                <div v-if="!loadedSkills.has(item.id)" class="shrink-0 w-4 h-4 rounded-full border-2 border-success border-t-transparent animate-spin" />
                <svg v-else class="shrink-0 w-4 h-4 text-success" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M4.5 12.75l6 6 9-13.5" />
                </svg>
              </div>
            </div>
          </div>

          <!-- Divider -->
          <div v-else-if="item.type === 'divider'" class="flex items-center gap-3 py-1 animate-chat-in">
            <div class="flex-1 h-px bg-th-text/[0.06]" />
            <span class="text-[10px] text-th-text-m font-mono tracking-widest uppercase">{{ t('hero.demo.laterLabel') }}</span>
            <div class="flex-1 h-px bg-th-text/[0.06]" />
          </div>
        </template>

        <!-- Typing indicator for agent (inline typing text) -->
        <div v-if="isTypingAgent && currentTypingText" class="flex justify-start">
          <div class="flex gap-2.5 max-w-[90%]">
            <div class="shrink-0 w-7 h-7 rounded-lg bg-gradient-to-br from-accent/20 to-accent-light/10 flex items-center justify-center ring-1 ring-accent/15 mt-0.5">
              <svg class="w-3.5 h-3.5 text-accent" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path stroke-linecap="round" stroke-linejoin="round" d="M9.813 15.904L9 18.75l-.813-2.846a4.5 4.5 0 00-3.09-3.09L2.25 12l2.846-.813a4.5 4.5 0 003.09-3.09L9 5.25l.813 2.846a4.5 4.5 0 003.09 3.09L15.75 12l-2.846.813a4.5 4.5 0 00-3.09 3.09z" />
              </svg>
            </div>
            <div class="px-3.5 py-2.5 rounded-2xl rounded-bl-md bg-th-bg-t/60 border border-th-text/[0.05] text-sm text-th-text-s leading-relaxed">
              {{ currentTypingText }}<span class="inline-block w-0.5 h-4 bg-accent/60 ml-0.5 animate-pulse align-middle" />
            </div>
          </div>
        </div>

        <!-- Dots typing indicator -->
        <div v-if="typingVisible" class="flex justify-start">
          <div class="flex gap-2.5">
            <div class="shrink-0 w-7 h-7 rounded-lg bg-gradient-to-br from-accent/20 to-accent-light/10 flex items-center justify-center ring-1 ring-accent/15 mt-0.5">
              <svg class="w-3.5 h-3.5 text-accent" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path stroke-linecap="round" stroke-linejoin="round" d="M9.813 15.904L9 18.75l-.813-2.846a4.5 4.5 0 00-3.09-3.09L2.25 12l2.846-.813a4.5 4.5 0 003.09-3.09L9 5.25l.813 2.846a4.5 4.5 0 003.09 3.09L15.75 12l-2.846.813a4.5 4.5 0 00-3.09 3.09z" />
              </svg>
            </div>
            <div class="px-4 py-3 rounded-2xl rounded-bl-md bg-th-bg-t/60 border border-th-text/[0.05] flex items-center gap-1">
              <span class="w-1.5 h-1.5 rounded-full bg-th-text-m/50 animate-bounce" style="animation-delay: 0ms" />
              <span class="w-1.5 h-1.5 rounded-full bg-th-text-m/50 animate-bounce" style="animation-delay: 150ms" />
              <span class="w-1.5 h-1.5 rounded-full bg-th-text-m/50 animate-bounce" style="animation-delay: 300ms" />
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.font-system {
  font-family: system-ui, -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
}

.scrollbar-hide {
  -ms-overflow-style: none;
  scrollbar-width: none;
}
.scrollbar-hide::-webkit-scrollbar {
  display: none;
}

@keyframes chat-in {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

.animate-chat-in {
  animation: chat-in 0.25s ease-out both;
}
</style>
