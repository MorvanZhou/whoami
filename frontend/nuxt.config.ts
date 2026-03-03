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
      meta: [
        { charset: 'utf-8' },
        { name: 'viewport', content: 'width=device-width, initial-scale=1' },
        { name: 'description', content: 'whoami gives every AI agent a persistent identity profile about you — synced across Cursor, Claude, Codex and more. One profile, every AI. Stop repeating yourself.' },
        { name: 'msvalidate.01', content: 'A25A16143F8A52E0CC3FD76FE0E5F03B' },
        // Open Graph
        { property: 'og:type', content: 'website' },
        { property: 'og:site_name', content: 'whoami' },
        { property: 'og:title', content: 'whoami — Your Persistent Identity for Every AI Agent' },
        { property: 'og:description', content: 'One profile that every AI agent can read. Works with Cursor, Claude, Codex. Stop repeating yourself.' },
        { property: 'og:image', content: 'https://whoamiagent.com/og-image.png' },
        { property: 'og:url', content: 'https://whoamiagent.com' },
        // Twitter Card
        { name: 'twitter:card', content: 'summary_large_image' },
        { name: 'twitter:title', content: 'whoami — Persistent AI Identity Profile' },
        { name: 'twitter:description', content: 'Let every AI know who you are. One profile, synced across all your AI agents.' },
        { name: 'twitter:image', content: 'https://whoamiagent.com/og-image.png' },
        // SEO keywords
        { name: 'keywords', content: 'AI agent memory, persistent user profile, cross-agent identity, AI identity sync, Cursor memory, Claude context, Codex memory, MCP skill, whoami' },
      ],
      link: [
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
            '@type': 'SoftwareApplication',
            'name': 'whoami',
            'description': 'Persistent AI identity profile that syncs across Cursor, Claude, Codex and all MCP-compatible AI agents.',
            'applicationCategory': 'ProductivityApplication',
            'operatingSystem': 'Web',
            'url': 'https://whoamiagent.com',
            'offers': {
              '@type': 'Offer',
              'price': '0',
              'priceCurrency': 'USD',
            },
          }),
        },
      ],
    },
  },
})