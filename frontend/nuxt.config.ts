export default defineNuxtConfig({
  compatibilityDate: '2025-01-01',
  devtools: { enabled: false },
  ssr: true,

  css: ['~/assets/css/main.css'],

  modules: ['@nuxtjs/tailwindcss', '@nuxtjs/i18n', 'nuxt-gtag', '@nuxtjs/sitemap'],

  gtag: {
    id: 'G-NCNN3SLNCT',
    enabled: process.env.NODE_ENV === 'production',
  },

  site: {
    url: 'https://whoamiagent.com',
  },

  sitemap: {
    sitemaps: false,
    exclude: [
      '/dashboard', '/zh/dashboard',
      '/login', '/zh/login',
      '/auth/**', '/zh/auth/**',
    ],
    autoLastmod: true,
    discoverImages: false,
    xslColumns: [
      { label: 'URL', width: '65%' },
      { label: 'Last Modified', select: 'sitemap:lastmod', width: '25%' },
    ],
  },

  i18n: {
    locales: [
      { code: 'en', language: 'en-US', file: 'en.json', name: 'English', isCatchallLocale: true },
      { code: 'zh', language: 'zh-CN', file: 'zh.json', name: '中文' },
    ],
    defaultLocale: 'en',
    langDir: '../locales',
    strategy: 'prefix_except_default',
    baseUrl: 'https://whoamiagent.com',
    detectBrowserLanguage: {
      useCookie: true,
      cookieKey: 'i18n_redirected',
      redirectOn: 'root',
    },
  },

  tailwindcss: {
    configPath: 'tailwind.config.ts',
  },

  runtimeConfig: {
    public: {
      apiBase: process.env.NUXT_PUBLIC_API_BASE || '/api',
    },
  },

  nitro: {
    devProxy: {
      '/api': {
        target: 'http://localhost:8002/api',
        changeOrigin: true,
      },
    },
  },

  app: {
    head: {
      title: 'whoami — Persistent AI Identity Profile | Works with Cursor, Claude & Codex',
      meta: [
        { charset: 'utf-8' },
        { name: 'viewport', content: 'width=device-width, initial-scale=1' },
        { name: 'description', content: 'whoami gives every AI agent a persistent identity profile about you — synced across Cursor, Claude, Codex and more. One profile, every AI. Stop repeating yourself.' },
        { name: 'msvalidate.01', content: 'A25A16143F8A52E0CC3FD76FE0E5F03B' },
        { name: 'robots', content: 'index, follow, max-image-preview:large, max-snippet:-1' },
        { name: 'referrer', content: 'no-referrer-when-downgrade' },
        { property: 'og:type', content: 'website' },
        { property: 'og:site_name', content: 'whoami' },
        { property: 'og:image', content: 'https://whoamiagent.com/og-image.png' },
        { name: 'twitter:card', content: 'summary_large_image' },
        { name: 'twitter:image', content: 'https://whoamiagent.com/og-image.png' },
        { name: 'keywords', content: 'whoami, AI agent memory, persistent AI identity, cross-AI user profile, AI identity sync, Cursor AI memory, Claude whoami, Claude context memory, Codex memory, MCP skill, AI personalization, cross-platform AI profile, AI agent skill, AI remembers you' },
      ],
      link: [
        { rel: 'icon', type: 'image/png', sizes: '32x32', href: '/favicon-32x32.png' },
        { rel: 'icon', type: 'image/png', sizes: '16x16', href: '/favicon-16x16.png' },
        { rel: 'apple-touch-icon', sizes: '180x180', href: '/apple-touch-icon.png' },
        { rel: 'preload', as: 'font', type: 'font/woff2', href: '/fonts/satoshi-400.woff2', crossorigin: '' },
        { rel: 'preload', as: 'font', type: 'font/woff2', href: '/fonts/space-mono-400.woff2', crossorigin: '' },
      ],
      script: [
        {
          type: 'application/ld+json',
          innerHTML: JSON.stringify({
            '@context': 'https://schema.org',
            '@graph': [
              {
                '@type': 'WebSite',
                'name': 'whoami',
                'alternateName': ['whoami Agent', 'whoamiagent'],
                'url': 'https://whoamiagent.com',
                'description': 'Persistent AI identity profile that syncs across Cursor, Claude, Codex and all MCP-compatible AI agents.',
                'inLanguage': ['en-US', 'zh-CN'],
              },
              {
                '@type': 'SoftwareApplication',
                'name': 'whoami',
                'description': 'whoami gives every AI agent a persistent identity profile about you — synced across Cursor, Claude, Codex and more. One profile, every AI. Stop repeating yourself.',
                'applicationCategory': 'DeveloperApplication',
                'applicationSubCategory': 'AI Developer Tool',
                'operatingSystem': 'Cross-platform',
                'url': 'https://whoamiagent.com',
                'downloadUrl': 'https://www.npmjs.com/package/skills',
                'featureList': [
                  'Cross-AI identity sync across Cursor, Claude, Codex, Windsurf, Cline and more',
                  'Persistent user profile that survives new chat sessions',
                  'User-controlled identity editing from dashboard or via AI commands',
                  'MCP-compatible agent skill architecture',
                  'Free to use with cloud-synced storage',
                ],
                'offers': {
                  '@type': 'Offer',
                  'price': '0',
                  'priceCurrency': 'USD',
                  'availability': 'https://schema.org/InStock',
                },
                'aggregateRating': {
                  '@type': 'AggregateRating',
                  'ratingValue': '5',
                  'ratingCount': '1',
                  'bestRating': '5',
                },
              },
              {
                '@type': 'FAQPage',
                'mainEntity': [
                  {
                    '@type': 'Question',
                    'name': 'What is whoami?',
                    'acceptedAnswer': {
                      '@type': 'Answer',
                      'text': 'whoami is a persistent AI identity profile service that syncs your personal context across all AI agents — including Cursor, Claude, Codex, Windsurf and Cline. Instead of repeating yourself every session, you write your profile once and every AI reads the same you.',
                    },
                  },
                  {
                    '@type': 'Question',
                    'name': 'How does Claude whoami work?',
                    'acceptedAnswer': {
                      '@type': 'Answer',
                      'text': 'Claude whoami works as an MCP-compatible agent skill. Install it via "npx skills add MorvanZhou/whoami", sign in at whoamiagent.com, and Claude automatically reads and updates your identity profile. Your profile is stored in the cloud and synced instantly to every connected agent including Claude, Cursor, and Codex.',
                    },
                  },
                  {
                    '@type': 'Question',
                    'name': 'Which AI tools does whoami support?',
                    'acceptedAnswer': {
                      '@type': 'Answer',
                      'text': 'whoami supports all MCP-compatible AI agents, including Cursor, Claude Code, OpenAI Codex, Windsurf, Cline, OpenClaw, and any tool that supports the MCP agent skill protocol.',
                    },
                  },
                  {
                    '@type': 'Question',
                    'name': 'Is whoami free to use?',
                    'acceptedAnswer': {
                      '@type': 'Answer',
                      'text': 'Yes, whoami is completely free to use. Sign up with GitHub or Google, create your profile, and start syncing your identity across all your AI agents at no cost.',
                    },
                  },
                  {
                    '@type': 'Question',
                    'name': 'How do I get started with whoami?',
                    'acceptedAnswer': {
                      '@type': 'Answer',
                      'text': 'Getting started takes seconds: paste the install prompt to any AI agent, or run "npx skills add MorvanZhou/whoami" in your terminal. Then sign in at whoamiagent.com/login to create your identity profile.',
                    },
                  },
                  {
                    '@type': 'Question',
                    'name': 'How to make AI remember who I am?',
                    'acceptedAnswer': {
                      '@type': 'Answer',
                      'text': 'Use whoami to create a persistent identity profile. Install whoami as an agent skill, write your profile once, and every AI agent — Cursor, Claude, Codex — reads the same context about you. No more repeating yourself in new sessions.',
                    },
                  },
                ],
              },
              {
                '@type': 'Organization',
                'name': 'whoami',
                'url': 'https://whoamiagent.com',
                'logo': 'https://whoamiagent.com/logo.png',
                'sameAs': [
                  'https://github.com/MorvanZhou/whoami',
                ],
              },
              {
                '@type': 'WebPage',
                'name': 'whoami — Persistent AI Identity Profile',
                'description': 'whoami gives every AI agent a persistent identity profile about you — synced across Cursor, Claude, Codex and more.',
                'url': 'https://whoamiagent.com',
                'speakable': {
                  '@type': 'SpeakableSpecification',
                  'cssSelector': ['h1', '.hero-subtitle', '.features-title'],
                },
              },
            ],
          }),
        },
      ],
    },
  },
})