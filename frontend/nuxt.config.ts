export default defineNuxtConfig({
  compatibilityDate: '2025-01-01',
  devtools: { enabled: false },
  ssr: true, // 启用SSR以生成完整HTML

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
  },

  i18n: {
    locales: [
      { code: 'en', file: 'en.json', name: 'English' },
      { code: 'zh', file: 'zh.json', name: '中文' },
    ],
    defaultLocale: 'en',
    langDir: '../locales',
    strategy: 'prefix_except_default',
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
      htmlAttrs: {
        lang: 'en',
      },
      meta: [
        { charset: 'utf-8' },
        { name: 'viewport', content: 'width=device-width, initial-scale=1' },
        { name: 'description', content: 'whoami gives every AI agent a persistent identity profile about you — synced across Cursor, Claude, Codex and more. One profile, every AI. Stop repeating yourself.' },
        { name: 'msvalidate.01', content: 'A25A16143F8A52E0CC3FD76FE0E5F03B' },
        { name: 'robots', content: 'index, follow, max-image-preview:large, max-snippet:-1' },
        { name: 'referrer', content: 'no-referrer-when-downgrade' },
        // Open Graph
        { property: 'og:type', content: 'website' },
        { property: 'og:site_name', content: 'whoami' },
        { property: 'og:title', content: 'whoami — Your Persistent Identity for Every AI Agent' },
        { property: 'og:description', content: 'One profile that every AI agent can read. Works with Cursor, Claude, Codex. Stop repeating yourself.' },
        { property: 'og:image', content: 'https://whoamiagent.com/og-image.png' },
        { property: 'og:url', content: 'https://whoamiagent.com' },
        { property: 'og:locale', content: 'en_US' },
        { property: 'og:locale:alternate', content: 'zh_CN' },
        // Twitter Card
        { name: 'twitter:card', content: 'summary_large_image' },
        { name: 'twitter:title', content: 'whoami — Persistent AI Identity Profile' },
        { name: 'twitter:description', content: 'Let every AI know who you are. One profile, synced across all your AI agents.' },
        { name: 'twitter:image', content: 'https://whoamiagent.com/og-image.png' },
        // SEO keywords
        { name: 'keywords', content: 'AI agent memory, persistent user profile, cross-agent identity, AI identity sync, Cursor memory, Claude context, Codex memory, MCP skill, whoami, AI personalization, cross-platform AI profile' },
      ],
      link: [
        { rel: 'canonical', href: 'https://whoamiagent.com' },
        { rel: 'icon', type: 'image/png', sizes: '32x32', href: '/favicon-32x32.png' },
        { rel: 'icon', type: 'image/png', sizes: '16x16', href: '/favicon-16x16.png' },
        { rel: 'apple-touch-icon', sizes: '180x180', href: '/apple-touch-icon.png' },
        { rel: 'preconnect', href: 'https://fonts.googleapis.com' },
        { rel: 'preconnect', href: 'https://fonts.gstatic.com', crossorigin: '' },
        {
          rel: 'stylesheet',
          href: 'https://fonts.googleapis.com/css2?family=Space+Mono:wght@400;700&display=swap',
        },
        {
          rel: 'stylesheet',
          href: 'https://api.fontshare.com/v2/css?f[]=satoshi@300,400,500,600,700,800&display=swap',
        },
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
                'alternateName': 'whoami Agent',
                'url': 'https://whoamiagent.com',
                'description': 'Persistent AI identity profile that syncs across Cursor, Claude, Codex and all MCP-compatible AI agents.',
                'inLanguage': ['en', 'zh'],
              },
              {
                '@type': 'SoftwareApplication',
                'name': 'whoami',
                'description': 'whoami gives every AI agent a persistent identity profile about you — synced across Cursor, Claude, Codex and more. One profile, every AI. Stop repeating yourself.',
                'applicationCategory': 'ProductivityApplication',
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
                    'name': 'How does whoami work?',
                    'acceptedAnswer': {
                      '@type': 'Answer',
                      'text': 'whoami works as an MCP-compatible agent skill. Install it via "npx skills add MorvanZhou/whoami", sign in at whoamiagent.com, and your AI agents automatically read and update your identity profile. Your profile is stored in the cloud and synced instantly to every connected agent.',
                    },
                  },
                  {
                    '@type': 'Question',
                    'name': 'Which AI tools does whoami support?',
                    'acceptedAnswer': {
                      '@type': 'Answer',
                      'text': 'whoami supports all MCP-compatible AI agents, including Cursor, Claude Code, OpenAI Codex, Windsurf, Cline, and any tool that supports the MCP agent skill protocol.',
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